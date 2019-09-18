create database user;

drop table if exists `user`;
CREATE TABLE `user`
(
    `user_id`     INT(10)      NOT NULL AUTO_INCREMENT,
    `username`    VARCHAR(64)  NULL DEFAULT NULL,
    `password`    VARCHAR(128) NULL DEFAULT NULL,
    `age`         INT(10)      NULL DEFAULT NULL,
    `create_time` DATE         NULL DEFAULT NULL,
    `modify_time` DATE         NULL DEFAULT NULL,

    PRIMARY KEY (`user_id`)
);