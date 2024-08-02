package gormdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heshiyingx/gotool/dbext/redis_script"
	"github.com/heshiyingx/gotool/strext"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

const (
	notFoundPlaceholder = "*"
	keyUpdatePrefix     = "updating:"
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, tempd1.05] * seconds
	expiryDeviation = 0.05
)

var TypeErr = errors.New("type is err")

type (
	CacheGormDB[P int64 | uint64 | string] struct {
		rdb               redis.UniversalClient
		singleFlight      *singleflight.Group
		notFoundExpireSec int
		cacheExpireSec    int
		randSec           int
		db                *gorm.DB
		antPool           *ants.Pool
		//antFailChan       chan []string
	}
	pkInfoDefine[P int64 | uint64 | string] struct {
		pkCacheKey string
		p          P
	}
)

func MustNewCacheGormDB[P int64 | uint64 | string](c Config) *CacheGormDB[P] {
	gormDB, err := NewCacheGormDB[P](c)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return gormDB
}

func NewCacheGormDB[P int64 | uint64 | string](c Config) (*CacheGormDB[P], error) {
	db, err := gorm.Open(getDialector(c), &c.GormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开的连接数
	sqlDB.SetMaxOpenConns(150)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	_, err = c.Rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	pool, err := ants.NewPool(runtime.NumCPU(), ants.WithExpiryDuration(time.Minute*5))
	if err != nil {
		return nil, err
	}
	cacheGromDB := &CacheGormDB[P]{
		rdb:               c.Rdb,
		singleFlight:      &singleflight.Group{},
		notFoundExpireSec: c.NotFoundExpireSec,
		cacheExpireSec:    c.CacheExpireSec,
		randSec:           c.RandSec,
		db:                db,
		antPool:           pool,
		//antFailChan:       make(chan []string, 20000),
	}
	if c.PreFunc != nil {
		c.PreFunc(cacheGromDB.db)
	}
	return cacheGromDB, nil
}

// QueryByCtx 查询
func (cg *CacheGormDB[P]) QueryByCtx(ctx context.Context, result any, key string, queryDBFun QueryCtxFn) error {
	err := cg.takeCacheCtx(ctx, key, result, queryDBFun, func(resultStr string, waitUpdate bool) error {
		if waitUpdate {
			_, err := cg.rdb.Set(ctx, key, resultStr, time.Second*2).Result()
			return err
		} else {
			isSet, err := cg.rdb.SetNX(ctx, key, resultStr, genDuring(cg.cacheExpireSec, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, key, result, time.Second*2).Result()
				return err
			}
			return nil
		}
	})
	return err
}

// QueryByCustomCacheSecCtx 自定义缓存时间
func (cg *CacheGormDB[P]) QueryByCustomCacheSecCtx(ctx context.Context, result any, key string, cacheSec int, queryDBFun QueryCtxFn) error {
	if cacheSec == 0 {
		cacheSec = 10
	}
	err := cg.takeCacheCtx(ctx, key, result, queryDBFun, func(resultStr string, waitUpdate bool) error {
		if waitUpdate {
			_, err := cg.rdb.Set(ctx, key, resultStr, time.Second*2).Result()
			return err
		} else {
			isSet, err := cg.rdb.SetNX(ctx, key, resultStr, genDuring(cacheSec, cg.randSec)).Result()
			if err != nil {
				return err
			}
			if !isSet {
				_, err = cg.rdb.Set(ctx, key, result, time.Second*2).Result()
				return err
			}
			return nil
		}
	})
	return err
}

// QuerySingleNoCacheCtx 同一时刻，只有一个能进入db查询，已key为判断是否为同一个查询
func (cg *CacheGormDB[P]) QuerySingleNoCacheCtx(ctx context.Context, key string, result any, queryDBFun QueryCtxFn) error {
	err := cg.takeCtx(ctx, key, result, queryDBFun)
	return err
}

// QueryNoCacheCtx 直接进行db查询
func (cg *CacheGormDB[P]) QueryNoCacheCtx(ctx context.Context, result any, queryDBFun QueryCtxFn) error {
	err := queryDBFun(ctx, result, cg.db)
	return err
}

// DelCacheKeys 删除缓存
func (cg *CacheGormDB[P]) DelCacheKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return cg.rdb.Del(ctx, keys...).Err()
}

// DelCacheKeysAndDelay 删除，并且加入延迟删除协程池
func (cg *CacheGormDB[P]) DelCacheKeysAndDelay(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	defer func() {
		err := cg.antPool.Submit(func() {
			deadline, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Second))
			defer cancelFunc()
			select {
			case <-deadline.Done():
			}
			err := cg.rdb.Del(context.Background(), keys...).Err()
			if err != nil {
				log.Printf("ant pool task doing err:%v", err)
				//cg.antFailChan <- keys
			}
		})
		if err != nil {
			log.Printf("ant pool task Submit err:%v", err)
		}
	}()
	return cg.rdb.Del(ctx, keys...).Err()
}
func (cg *CacheGormDB[P]) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn) error {

	singleResult, err, share := cg.singleFlight.Do(key, func() (interface{}, error) {

		err := query(ctx, result, cg.db)
		logx.WithContext(ctx).Debugf("takeCtx->queryFinish   key:%v,  result:%v,  err:%v", key, strext.ToJsonStr(result), err)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, gorm.ErrRecordNotFound
			} else {
				return nil, err
			}
		}

		return result, nil
	})
	if err != nil {
		return err
	}
	if share {
		v := reflect.ValueOf(result).Elem()                     // 获取 result 指针指向的值的 reflect.Value
		singResultValue := reflect.ValueOf(singleResult).Elem() // 获取 singleResult 的值
		if v.Type() != singResultValue.Type() {                 // 检查类型是否匹配
			return fmt.Errorf("unexpected type:%T", singleResult)
		}
		v.Set(singResultValue) // 更新 result 的值
	}
	return err
}
func (cg *CacheGormDB[P]) takeCacheCtx(ctx context.Context, key string, result any, query QueryCtxFn, cacheFn CacheFn) error {

	singleResult, err, share := cg.singleFlight.Do(key, func() (interface{}, error) {
		//fmt.Println("进入redis缓存")
		val, err := cg.rdb.Get(ctx, key).Result()
		logx.WithContext(ctx).Debugf("takeCacheCtx->CacheGet   key:%v,  var:%v,  err:%v", key, val, err)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				err = nil
			} else {
				return nil, err
			}
		}

		if val == notFoundPlaceholder {
			return nil, gorm.ErrRecordNotFound
		}
		if val != "" {
			err = json.Unmarshal([]byte(val), result)
			logx.WithContext(ctx).Debugf("takeCacheCtx->CacheToResult   key:%v,result:%v,jsonStr:%v,err:%v", key, strext.ToJsonStr(result), val, err)
			return result, err

		}

		err = query(ctx, result, cg.db)
		logx.WithContext(ctx).Debugf("takeCacheCtx->queryFinish   key:%v  result:%v  err:%v", key, strext.ToJsonStr(result), err)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = cg.setCacheWithNotFound(ctx, key)
				if cg.db.Logger != nil && err != nil {
					cg.db.Logger.Error(ctx, "setCacheWithNotFound err: %v key:%v", err, key)
				}
				return nil, gorm.ErrRecordNotFound
			} else {
				return nil, err
			}
		}

		resultBytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		isUpdating := true
		_, err = cg.rdb.Get(ctx, keyUpdatePrefix+key).Result()
		if errors.Is(err, redis.Nil) {
			isUpdating = false
			err = nil
		}
		err = cacheFn(string(resultBytes), isUpdating)
		logx.WithContext(ctx).Debugf("takeCacheCtx->cacheFinish  key:%v,result:%v,jsonStr_val:%v,err:%v", key, strext.ToJsonStr(result), string(resultBytes), err)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	if err != nil {
		return err
	}
	if share {
		v := reflect.ValueOf(result).Elem()                     // 获取 result 指针指向的值的 reflect.Value
		singResultValue := reflect.ValueOf(singleResult).Elem() // 获取 singleResult 的值
		if v.Type() != singResultValue.Type() {                 // 检查类型是否匹配
			return fmt.Errorf("unexpected type:%T", singleResult)
		}
		v.Set(singResultValue) // 更新 result 的值
	}
	return err
}

// ExecCtx 执行除查询之外的其他操作
func (cg *CacheGormDB[P]) ExecCtx(ctx context.Context, execFn ExecCtxFn, keys ...string) (int64, error) {
	if len(keys) > 0 {
		err := cg.rdb.Del(ctx, keys...).Err()
		if err != nil {
			return 0, err
		}
	}
	defer func() {
		if len(keys) > 0 {
			err := cg.antPool.Submit(func() {
				deadline, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Second))
				defer cancelFunc()
				select {
				case <-deadline.Done():
				}
				err := cg.rdb.Del(context.Background(), keys...).Err()
				if err != nil {
					log.Printf("ant pool task doing err:%v", err)
					//cg.antFailChan <- keys
				}
			})
			if err != nil {
				log.Printf("ant pool task Submit err:%v", err)
			}
		}
	}()
	for _, key := range keys {
		_, err := redis_script.IncrExpireScript.Run(ctx, cg.rdb, []string{keyUpdatePrefix + key}, 20).Result()
		if err != nil {
			return 0, err
		}
	}
	defer func() {
		for _, key := range keys {
			_, _ = redis_script.DecrZeroDelScript.Run(ctx, cg.rdb, []string{keyUpdatePrefix + key}).Result()

		}
	}()
	result, err := execFn(ctx, cg.db)
	if err != nil {
		return 0, err
	}

	if len(keys) > 0 {
		err = cg.rdb.Del(ctx, keys...).Err()
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (cg *CacheGormDB[P]) setCacheWithNotFound(ctx context.Context, key string) error {
	expire := time.Second*time.Duration(cg.notFoundExpireSec) + genDuring(cg.randSec, cg.notFoundExpireSec)
	_, err := cg.rdb.SetNX(ctx, key, notFoundPlaceholder, expire).Result()
	return err
}
func (cg *CacheGormDB[P]) GetRdb() redis.UniversalClient {
	return cg.rdb
}
func genDuring(oriSec int, randSec int) time.Duration {
	if oriSec == 0 {
		return 0
	}
	if randSec == 0 {
		randSec = 5
	}
	n := rand.Int31n(int32(time.Duration(randSec) * time.Second / time.Millisecond))
	return time.Duration(n)*time.Millisecond + time.Duration(oriSec)*time.Second
}
