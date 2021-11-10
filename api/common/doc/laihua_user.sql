-- MySQL dump 10.13  Distrib 8.0.25, for macos11 (x86_64)
--
-- Host: 127.0.0.1    Database: laihua_user
-- ------------------------------------------------------
-- Server version	8.0.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admin_api`
--

DROP TABLE IF EXISTS `admin_api`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '接口名',
  `http_method` varchar(32) NOT NULL COMMENT 'http请求方式',
  `http_path` varchar(32) NOT NULL COMMENT 'http请求路径',
  `platform_id` int unsigned NOT NULL COMMENT '平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`platform_id`,`name`),
  UNIQUE KEY `api_identify` (`platform_id`,`http_method`,`http_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-api表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_permission`
--

DROP TABLE IF EXISTS `admin_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_permission` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `permission_name` varchar(64) NOT NULL COMMENT '权限名称',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `permission_name` (`platform_id`,`permission_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_permission_has_api`
--

DROP TABLE IF EXISTS `admin_permission_has_api`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_permission_has_api` (
  `permission_id` int unsigned NOT NULL COMMENT '权限id',
  `api_id` int unsigned NOT NULL COMMENT 'api_id',
  PRIMARY KEY (`permission_id`,`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-权限与api中间表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_platform`
--

DROP TABLE IF EXISTS `admin_platform`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_platform` (
  `id` int unsigned NOT NULL,
  `platform_en` varchar(64) NOT NULL COMMENT '平台-英文',
  `platform_zh` varchar(64) NOT NULL COMMENT '平台-中文',
  PRIMARY KEY (`id`),
  UNIQUE KEY `platform` (`platform_en`,`platform_zh`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_role`
--

DROP TABLE IF EXISTS `admin_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `role_name` varchar(64) NOT NULL COMMENT '角色名称',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '关联的平台id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_name` (`platform_id`,`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_role_has_permission`
--

DROP TABLE IF EXISTS `admin_role_has_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_role_has_permission` (
  `role_id` int unsigned NOT NULL COMMENT '角色id',
  `permission_id` int unsigned NOT NULL COMMENT '权限id',
  PRIMARY KEY (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-角色与权限中间表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_user`
--

DROP TABLE IF EXISTS `admin_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账户名',
  `password` varchar(255) NOT NULL COMMENT '密码（哈希加盐）',
  `nick_name` varchar(64) NOT NULL DEFAULT '无' COMMENT '昵称',
  `phone` varchar(20) NOT NULL DEFAULT '无' COMMENT '手机号',
  `email` varchar(64) NOT NULL DEFAULT '无' COMMENT '电子邮箱',
  `platform_id` int unsigned NOT NULL DEFAULT '0' COMMENT '平台id',
  `is_super_admin` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否为超级管理员，0为false，1为true',
  `is_ban` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否封禁  0  false  1 true',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `updater` int unsigned NOT NULL DEFAULT '0' COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_UNIQUE` (`platform_id`,`phone`),
  UNIQUE KEY `email_UNIQUE` (`platform_id`,`email`),
  UNIQUE KEY `account_UNIQUE` (`platform_id`,`account`),
  KEY `created_at` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_user_has_role`
--

DROP TABLE IF EXISTS `admin_user_has_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_user_has_role` (
  `user_id` int unsigned NOT NULL COMMENT '用户id',
  `role_id` int unsigned NOT NULL COMMENT '角色id',
  PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端-用户与角色中间表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-11-09 10:31:33
