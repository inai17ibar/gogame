create database if not exists gogame_db;

create table if not exists gogame_db.user_table
(
    id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    username varchar(32) UNIQUE NOT NULL DEFAULT '',
    password varchar(64) NOT NULL DEFAULT '',
    created_date datetime DEFAULT CURRENT_TIMESTAMP
);

#insert into gogame_db.user_table (username, password) VALUES
#('Soichiro.Inatani', 'inatani'),
#('inai17ibar', '17tani');


create table if not exists gogame_db.gacha_table
(
	id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL DEFAULT '',
    type INT NOT NULL DEFAULT 0,
    rarity1_weight INT DEFAULT 520,
    rarity2_weight INT DEFAULT 300,
    rarity3_weight INT DEFAULT 120,
    rarity4_weight INT DEFAULT 45,
    rarity5_weight INT DEFAULT 15,
    started_date datetime DEFAULT CURRENT_TIMESTAMP,
    end_date datetime DEFAULT NULL
);

create table if not exists gogame_db.gacha_characters_01_table
(
	id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    characterid INT NOT NULL,
    weight INT NOT NULL,
    started_date datetime DEFAULT CURRENT_TIMESTAMP,
    end_date datetime DEFAULT NULL
);

create table if not exists gogame_db.user_characters_table
(
	id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    userid INT NOT NULL,
    characterid INT NOT NULL,
    created_date datetime DEFAULT CURRENT_TIMESTAMP,
    update_date datetime DEFAULT CURRENT_TIMESTAMP
);

create table if not exists gogame_db.character_table
(
    id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL DEFAULT '',
    rarity int DEFAULT 1,
    created_date datetime DEFAULT CURRENT_TIMESTAMP,
    update_date datetime DEFAULT CURRENT_TIMESTAMP
);
 
create table if not exists gogame_db.gacha_rarity_table
(
	id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    gacha_id INT DEFAULT NULL,
    rarity INT DEFAULT 1 NOT NULL,
    rarity_weight INT DEFAULT 520 NOT NULL,
    started_date datetime DEFAULT CURRENT_TIMESTAMP,
    end_date datetime DEFAULT NULL,
    UNIQUE(gacha_id, rarity)
);

create table if not exists gacha_characters_table 
select 
character_table.id as id,
character_table.name as name,
character_table.rarity as rarity,
gacha_rarity_table.rarity_weight as rarity_weight 
from character_table join gacha_rarity_table on character_table.rarity = gacha_rarity_table.rarity;