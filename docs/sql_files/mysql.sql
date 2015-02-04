SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

DROP TABLE IF EXISTS `backups`;
CREATE TABLE IF NOT EXISTS `backups` (
`id` int(10) unsigned NOT NULL,
  `path` text COLLATE utf8_unicode_ci NOT NULL,
  `filename` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `checksum` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `atime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `uid` int(5) NOT NULL,
  `gid` int(5) NOT NULL,
  `perms` int(4) NOT NULL,
  `host_updated` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


ALTER TABLE `backups`
 ADD PRIMARY KEY (`id`), ADD KEY `filename` (`filename`,`host_updated`,`last_update`), ADD FULLTEXT KEY `path` (`path`);


ALTER TABLE `backups`
MODIFY `id` int(10) unsigned NOT NULL AUTO_INCREMENT;