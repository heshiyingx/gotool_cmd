create table feed_store
(
    id         bigint auto_increment,
    user_id    bigint                             not null comment '用户ID',
    feed_num   int      default 0                 not null comment '私聊数量',
    created_at datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint uin_userid
        unique (user_id),
    primary key (id)
)comment '饲料仓库';

create table feed_store_history
(
    id         bigint auto_increment comment '主键' PRIMARY KEY,
    user_id    bigint                                 not null comment '用户ID',
    chicken_id bigint                                 not null comment '小鸡ID',
    op_type    int                                    not null comment '操作类型,1:增加，2：减少',
    value      int          default 0                 not null comment '改变的值',
    op_id      varchar(32)                            not null comment '操作id，用于去重',
    comment    varchar(100) default ''                not null comment '备注说明',
    created_at datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint feed_store_history_uk
        unique (user_id, op_id)
)comment '饲料仓库改变历史';

create table nutrition_store
(
    id         bigint auto_increment PRIMARY KEY,
    user_id    bigint                             not null comment '用户ID',
    num        int      default 0                 not null comment '数量',
    created_at datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint uin_userid
        unique (user_id)
)comment '营养素仓库';

create table nutrition_store_history
(
    id         bigint auto_increment comment '主键' PRIMARY KEY,
    user_id    bigint                                 not null comment '用户ID',
    chicken_id bigint                                 not null comment '小鸡ID',
    op_type    int                                    not null comment '操作类型,1:增加，2：减少',
    value      int          default 0                 not null comment '改变的值',
    op_id      varchar(32)                            not null comment '操作id，用于去重',
    comment    varchar(100) default ''                not null comment '备注说明',
    created_at datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint nutrition_store_history_pk
        unique (user_id, op_id)
) comment '饲料仓库改变历史';

CREATE TABLE `chicken`
(
    `id`              bigint      NOT NULL AUTO_INCREMENT COMMENT '主键ID' PRIMARY KEY,
    `user_id`         bigint      NOT NULL COMMENT '用户ID',
    `feed_slot_value` int         NOT NULL DEFAULT '0' COMMENT '饲料槽中的饲料数量',
    `name`            varchar(50) NOT NULL DEFAULT '' COMMENT '小鸡名字',
    `feed_num`        bigint      NOT NULL DEFAULT '0' COMMENT '喂食次数',
    `op_time`         bigint      NOT NULL COMMENT '上次操作时间',
    `stage`           tinyint     NOT NULL DEFAULT '0' COMMENT '所处阶段，1:孵化期，2.成长期，3.下蛋期',
    `is_die`          tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否死亡,0:正常，1:已死亡',
    `nutrition_value` int         NOT NULL DEFAULT '0' COMMENT '营养值',
    `feed_time`       bigint      NOT NULL DEFAULT '0' COMMENT '最近一次喂饲料的时间(时间戳秒)',
    `process`         int         NOT NULL DEFAULT '0' COMMENT '进度',
    `created_at`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      datetime             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `config`          json                 DEFAULT NULL COMMENT '喂养配置快照',
    KEY               `chicken_user_id_is_die_index` (`user_id`,`is_die`)
);
