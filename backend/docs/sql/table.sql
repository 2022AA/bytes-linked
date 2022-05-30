create database if not exists `filestore`  default character set utf8mb4 collate utf8mb4_general_ci;

use `filestore`;

-- 创建文件表
DROP TABLE IF EXISTS `file`;
CREATE TABLE `file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` varchar(50) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `img_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '缩略文件存储位置',
  `owner_uid` int(11) DEFAULT '0' COMMENT '文件拥有者',
  `creator_uid` int(11) DEFAULT '0' COMMENT '文件创作者',
  `like_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `ar_tag` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用AR',
  `create_at` datetime default NOW() COMMENT '创建日期',
  `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(2公开发布/1可用/0禁用/-1已删除等状态)',
  `transaction_id` varchar(256) NOT NULL DEFAULT '' COMMENT '区块交易id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_like` (`like_cnt`, `status`),
  KEY `idx_owner` (`owner_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建用户表
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
  `invite_code` varchar(64) NOT NULL COMMENT '邀请码',
  `phone` varchar(64) NOT NULL COMMENT '手机号',
  `avatar_url` varchar(256) NOT NULL DEFAULT '' COMMENT '头像地址',
  `balance` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

-- 创建用户文件表
DROP TABLE IF EXISTS `user_file`;
CREATE TABLE `user_file` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL,
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
  UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 邀请码表
DROP TABLE IF EXISTS `invite_code`;
CREATE TABLE `invite_code` (
    `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `code` varchar(64) NOT NULL COMMENT '邀请码',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '邀请码状态(0未使用，1已使用)',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `last_update` datetime DEFAULT CURRENT_TIMESTAMP
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户secret file
DROP TABLE IF EXISTS `user_secret_file`;
CREATE TABLE `user_secret_file` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL,
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
  `type` int(11) NOT NULL COMMENT '类型(1公钥2私钥)',
  `data` varchar(256) Not NULL COMMENT '文件内容',
  UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 流水
DROP TABLE IF EXISTS `user_order`;
CREATE TABLE `user_order` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `type` int(11) NOT NULL COMMENT '流水类型(1转移流水,点赞流水)',
  `order_sn` varchar(256) NOT NULL COMMENT '流水号',
  `content` varchar(256) NOT NULL COMMENT '流水内容',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `uid` int(11) NOT NULL COMMENT '交易的用户',
  `to_id` int(11) NOT NULL COMMENT '被交易的用户/物品',
  `transaction_id` varchar(256) NOT NULL DEFAULT '' COMMENT '区块交易id',
  UNIQUE KEY `idx_order_number` (`order_sn`),
  KEY `idx_uid` (`uid`),
  KEY `idx_to_id` (`to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
