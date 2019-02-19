CREATE DATABASE  IF NOT EXISTS `diplomacy` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `diplomacy`;

/* Create table for games */

DROP TABLE IF EXISTS games;
CREATE TABLE games (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(50) CHARACTER SET utf8 NOT NULL,
  started_at TIMESTAMP NULL DEFAULT NULL,
  game_year DATE,
  phase int NOT NULL DEFAULT 1, 
  phase_end TIMESTAMP NULL DEFAULT NULL,
  orders_interval int NOT NULL DEFAULT 1,
  PRIMARY KEY (id)
);
