-- create table feed_store_history
-- (
--     id         bigint auto_increment comment '主键' PRIMARY KEY,
--     user_id    bigint                                 not null comment '用户ID',
--     chicken_id bigint                                 not null comment '小鸡ID',
--     op_type    int                                    not null comment '操作类型,1:增加，2：减少',
--     value      int          default 0                 not null comment '改变的值',
--     op_id      varchar(32)                            not null comment '操作id，用于去重',
--     comment    varchar(100) default ''                not null comment '备注说明',
--     created_at datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
--     updated_at datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
--     constraint feed_store_history_uk
--         unique (user_id, op_id)
-- )comment '饲料仓库改变历史';
--
-- create table nutrition_store
-- (
--     id         bigint auto_increment PRIMARY KEY,
--     user_id    bigint                             not null comment '用户ID',
--     num        int      default 0                 not null comment '数量',
--     created_at datetime default CURRENT_TIMESTAMP not null comment '创建时间',
--     updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
--     constraint uin_userid
--         unique (user_id)
-- )comment '营养素仓库';

create table nutrition_store_history
(
    id         bigint auto_increment comment '主键' PRIMARY KEY,
    user_id    bigint                                 not null comment '用户ID',
    chicken_id bigint                                 not null comment '小鸡ID',
    op_type    int                                    not null comment '操作类型,1:增加，2：减少',
    value      int          default 0                 not null comment '改变的值',
    op_id      varchar(32)                            not null comment '操作id，用于去重',
    click_id    varchar(36)                            not null comment '唯一ID',
    comment    varchar(100) default ''                not null comment '备注说明',
    created_at datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint nutrition_store_history_pk
        unique (user_id, op_id),
    unique key click_id(click_id)
) comment '饲料仓库改变历史';

