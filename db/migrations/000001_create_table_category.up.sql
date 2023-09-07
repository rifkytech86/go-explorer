create table diary
(
    diary_id         int auto_increment
        primary key,
    user_id          bigint                             null,
    diary_desc       text                               null,
    diary_link_video varchar(255)                       null,
    updated_at       datetime                           null,
    created_at       datetime default CURRENT_TIMESTAMP null,
    diary_is_active  int      default 0                 null
);

create table period_tracker
(
    period_tracker_id bigint auto_increment
        primary key,
    user_id           bigint   null,
    start_period      datetime null,
    end_period        datetime null,
    created_at        datetime null,
    updated_at        datetime null
)
    comment 'period tracker represetnt';

create table user
(
    user_id        bigint auto_increment
        primary key,
    user_name      varchar(255)                       null,
    user_email     varchar(255)                       null,
    user_password  varchar(255)                       null,
    user_dob       datetime                           null,
    user_school    varchar(255)                       null,
    user_is_verify smallint                           null,
    updated_at     datetime                           null,
    created_at     datetime default CURRENT_TIMESTAMP null invisible,
    created_by     int                                null,
    updated_by     int                                null
)
    comment 'represent the user for mobile android or ios ';

create table user_admin
(
    id         bigint auto_increment
        primary key,
    email      varchar(200) null,
    username   varchar(100) null,
    created_at datetime     null,
    updated_at datetime     null,
    password   varchar(200) null,
    name       varchar(200) null,
    is_login   int          null,
    created_by int          null,
    updated_by int          null
);

create table user_mood
(
    id         int auto_increment
        primary key,
    name       varchar(100) null,
    created_at datetime     null,
    updated_at datetime     null,
    created_by int          null,
    updated_by int          null
);

