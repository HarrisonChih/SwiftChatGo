/*
 Navicat Premium Data Transfer

 Source Server         : Mysql
 Source Server Type    : MySQL
 Source Server Version : 80011
 Source Host           : localhost:3306
 Source Schema         : ginchat

 Target Server Type    : MySQL
 Target Server Version : 80011
 File Encoding         : 65001

 Date: 07/12/2025 13:12:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for communities
-- ----------------------------
DROP TABLE IF EXISTS `communities`;
CREATE TABLE `communities`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `owner_id` bigint(20) UNSIGNED NULL DEFAULT NULL,
  `img` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_communities_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for contacts
-- ----------------------------
DROP TABLE IF EXISTS `contacts`;
CREATE TABLE `contacts`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `owner_id` bigint(20) UNSIGNED NULL DEFAULT NULL,
  `target_id` bigint(20) UNSIGNED NULL DEFAULT NULL,
  `type` bigint(20) NULL DEFAULT NULL,
  `desc` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_contacts_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `from_id` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `target` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `type` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `media` bigint(20) NULL DEFAULT NULL,
  `content` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `pic` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `url` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `desc` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `amount` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_message_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_basic
-- ----------------------------
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `pass_word` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `phone` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `email` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `identity` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `client_ip` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `client_port` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `salt` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `login_time` datetime(3) NULL DEFAULT NULL,
  `heartbeat_time` datetime(3) NULL DEFAULT NULL,
  `login_out_time` datetime(3) NULL DEFAULT NULL,
  `is_log_out` tinyint(1) NULL DEFAULT NULL,
  `device_info` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_basic_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
