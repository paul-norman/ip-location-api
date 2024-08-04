CREATE TABLE IF NOT EXISTS `ip_city` (
	`ip_range_start`	varchar(45) NOT NULL,
	`ip_range_end`		varchar(45) NOT NULL,
	`ip_number_start`	varbinary(16) NOT NULL,
	`ip_number_end`		varbinary(16) NOT NULL,
	`country_code`		varchar(2) NOT NULL,
	`state1`			varchar(100) NOT NULL,
	`state2`			varchar(100) NOT NULL,
	`city`				varchar(100) NOT NULL,
	`postcode`			varchar(50) NOT NULL,
	`latitude`			decimal(11,8) NOT NULL,
	`longitude`			decimal(11,8) NOT NULL,
	`timezone`			varchar(20) NOT NULL,
	`ip_version`		int(10) NOT NULL DEFAULT 4,
	`db_version`		int(10) NOT NULL DEFAULT 1,

	KEY `ip_number_start` (`ip_number_start`),
	KEY `ip_number_end` (`ip_number_end`),
	KEY `country_code` (`country_code`),
	KEY `ip_version` (`ip_version`),
	KEY `db_version` (`db_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE IF NOT EXISTS `ip_asn` (
	`ip_range_start`	varchar(45) NOT NULL,
	`ip_range_end`		varchar(45) NOT NULL,
	`ip_number_start`	varbinary(16) NOT NULL,
	`ip_number_end`		varbinary(16) NOT NULL,
	`as_number`			int(10) NOT NULL DEFAULT 0,
	`as_organisation`	varchar(100) NOT NULL,
	`ip_version`		int(10) NOT NULL DEFAULT 4,
	`db_version`		int(10) NOT NULL DEFAULT 1,

	KEY `ip_number_start` (`ip_number_start`),
	KEY `ip_number_end` (`ip_number_end`),
	KEY `ip_version` (`ip_version`),
	KEY `db_version` (`db_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE IF NOT EXISTS `ip_country` (
	`ip_range_start`	varchar(45) NOT NULL,
	`ip_range_end`		varchar(45) NOT NULL,
	`ip_number_start`	varbinary(16) NOT NULL,
	`ip_number_end`		varbinary(16) NOT NULL,
	`country_code`		varchar(2) NOT NULL,
	`ip_version`		int(10) NOT NULL DEFAULT 4,
	`db_version`		int(10) NOT NULL DEFAULT 1,

	KEY `ip_number_start` (`ip_number_start`),
	KEY `ip_number_end` (`ip_number_end`),
	KEY `country_code` (`country_code`),
	KEY `ip_version` (`ip_version`),
	KEY `db_version` (`db_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;