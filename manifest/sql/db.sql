CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id',
  `name` varchar(45) DEFAULT NULL COMMENT 'user name',
  `status` tinyint DEFAULT NULL COMMENT 'user status',
  `age` tinyint unsigned DEFAULT NULL COMMENT 'user age',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
