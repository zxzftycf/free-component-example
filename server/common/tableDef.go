package common

const (
	Redis_OnlineUser_Table string = "online"
	Redis_PlayerInfo_Table string = "player"
)

const (
	Mysql_AccountInfo_Table string = "game_account_info"
	Mysql_PlayerInfo_Table  string = "game_player_info"
)

const (
	Mysql_Check_PlayerInfo_Table string = `
		CREATE TABLE IF NOT EXISTS game_player_info (
		auto_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(256) NOT NULL PRIMARY KEY,
		short_id VARCHAR(10) NOT NULL,
		info MEDIUMBLOB NOT NULL,
		update_time BIGINT NOT NULL,
		KEY (update_time),
		UNIQUE KEY (auto_id),
		UNIQUE KEY (short_id)
	) ENGINE=InnoDB;
	`
	Mysql_Check_AccountInfo_Table string = `
		CREATE TABLE IF NOT EXISTS game_account_info (
		auto_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(256) NOT NULL PRIMARY KEY,
		short_id VARCHAR(10) NOT NULL,
		account VARCHAR(20) NOT NULL,
		password VARCHAR(20) NOT NULL,
		update_time BIGINT NOT NULL,
		KEY (update_time),
		UNIQUE KEY (account),
		UNIQUE KEY (auto_id),
		KEY (password),
		UNIQUE KEY (short_id)
	) ENGINE=InnoDB;
	`
)
