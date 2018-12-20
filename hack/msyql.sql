SET NAMES utf8mb4;
SET character_set_server =  'utf8mb4';
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';


DROP DATABASE IF EXISTS `sample`;
CREATE DATABASE IF NOT EXISTS `sample`;

GRANT ALL PRIVILEGES ON bwin.* TO 'sampleusr'@'%';
DROP USER 'sampleusr'@'%';

CREATE USER 'sampleusr'@'%' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON bwin.* TO 'sampleusr'@'%';

USE bwin;


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rawobj` json NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;