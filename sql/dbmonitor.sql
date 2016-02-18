create database dbmonitor;
use dbmonitor;

create table db_account (
    id int unsigned not null auto_increment,
    cluster varchar(40) not null default '' comment '业务名',
    addr varchar(25) not null default '' comment 'db address like 127.0.0.1:3306',
    created_time int(10) unsigned not null default 0,
    primary key (id),
    unique key uk_cluster_addr(cluster, addr)
)engine =innodb default charset=utf8;

create table db_status (
    id int unsigned not null default 0,
    qps int unsigned not null default 0,
    com_select int unsigned not null default 0,
    com_insert int unsigned not null default 0,
    com_delete int unsigned not null default 0,
    com_update int unsigned not null default 0,
    byte_received int unsigned not null default 0,
    byte_sent int unsigned not null default 0,
    rows_delete int unsigned not null default 0,
    rows_insert int unsigned not null default 0,
    rows_update int unsigned not null default 0,
    rows_read int unsigned not null default 0,
    thread_created int unsigned not null default 0,
    thread_running int unsigned not null default 0,
    thread_connected int unsigned not null default 0,
    created_time int unsigned not null default 0,
    primary key(id, created_time)
)engine=innodb default charset=utf8;

create table db_shreshold (
    id int unsigned not null default 0,
    qps int unsigned not null default 0,
    com_select int unsigned not null default 0,
    com_insert int unsigned not null default 0,
    com_delete int unsigned not null default 0,
    com_update int unsigned not null default 0,
    byte_received int unsigned not null default 0,
    byte_sent int unsigned not null default 0,
    rows_delete int unsigned not null default 0,
    rows_insert int unsigned not null default 0,
    rows_update int unsigned not null default 0,
    rows_read int unsigned not null default 0,
    thread_created int unsigned not null default 0,
    thread_running int unsigned not null default 0,
    thread_connected int unsigned not null default 0,
    created_time int unsigned not null default 0,
    primary key(id)
)engine=innodb default charset=utf8;
