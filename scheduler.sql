/*
Navicat MySQL Data Transfer

Source Server         : 计费
Source Server Version : 50520
Source Host           : 192.168.159.39:3306
Source Database       : scheduler

Target Server Type    : MYSQL
Target Server Version : 50520
File Encoding         : 65001

Date: 2018-07-30 13:42:24
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for performance_monitor
-- ----------------------------
DROP TABLE IF EXISTS `performance_monitor`;
CREATE TABLE `performance_monitor` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ip` varchar(20) DEFAULT NULL,
  `use_cpu` varchar(10) DEFAULT NULL,
  `use_memory` varchar(10) DEFAULT NULL,
  `update_time` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27082 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for register_center
-- ----------------------------
DROP TABLE IF EXISTS `register_center`;
CREATE TABLE `register_center` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ip` varchar(20) DEFAULT NULL,
  `register_time` varchar(20) DEFAULT NULL,
  `update_time` varchar(20) DEFAULT NULL,
  `is_death` tinyint(4) DEFAULT '0' COMMENT '是否死亡 0否1是',
  `task_slice_collection` varchar(200) DEFAULT NULL,
  `is_delete` tinyint(4) DEFAULT '0' COMMENT '是否删除  0否1是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_index` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=883269322 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除  0否  1是',
  `task_no` varchar(11) DEFAULT NULL,
  `group_no` varchar(11) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `task_slice` int(11) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `concurrency_num` int(11) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '1 启动 2停止',
  `is_activity` int(11) DEFAULT NULL COMMENT '0 否 1 是',
  `is_record_log` int(11) DEFAULT '0' COMMENT '是否记录日志 0 否 1 是',
  `param` varchar(255) DEFAULT NULL,
  `task_type` int(11) DEFAULT NULL,
  `cron` varchar(255) DEFAULT NULL,
  `create_time` varchar(20) DEFAULT NULL,
  `update_time` varchar(20) DEFAULT NULL,
  `create_user` varchar(20) DEFAULT NULL,
  `update_user` varchar(20) DEFAULT NULL,
  `mail` varchar(255) DEFAULT NULL,
  `request_methods` int(11) DEFAULT '1' COMMENT '1 POST 2 GET',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=863867504 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for task_execute_record
-- ----------------------------
DROP TABLE IF EXISTS `task_execute_record`;
CREATE TABLE `task_execute_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `task_no` varchar(20) DEFAULT NULL,
  `task_name` varchar(20) DEFAULT NULL,
  `group_no` varchar(20) DEFAULT NULL,
  `status` int(11) DEFAULT NULL COMMENT '0 失败 1 成功',
  `error_msg` varchar(500) DEFAULT NULL COMMENT '错误信息',
  `result` varchar(20) DEFAULT NULL,
  `create_time` varchar(20) DEFAULT NULL,
  `is_delete` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=818065808 DEFAULT CHARSET=utf8;
