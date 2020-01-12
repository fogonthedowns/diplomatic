USE `diplomacy`;

CREATE TABLE users (
 `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 `email` varchar(255) NOT NULL,
 `user_name`  varchar(255) NOT NULL,
 `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT into users (email, user_name, password) VALUES ("one@user.com", "user one", "12345");
INSERT into users (email, user_name, password) VALUES ("two@user.com", "user two", "12345");
INSERT into users (email, user_name, password) VALUES ("three@user.com", "user three", "12345");
INSERT into users (email, user_name, password) VALUES ("four@user.com", "user four", "12345");
INSERT into users (email, user_name, password) VALUES ("five@user.com", "user five", "12345");
INSERT into users (email, user_name, password) VALUES ("six@user.com", "user six", "12345");
INSERT into users (email, user_name, password) VALUES ("seven@user.com", "user seven", "12345");
