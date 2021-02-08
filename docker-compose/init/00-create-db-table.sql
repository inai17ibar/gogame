create database if not exists gogame_db;

create table if not exists gogame_db.user_table
(
    id INT AUTO_INCREMENT UNIQUE NOT NULL PRIMARY KEY,
    username varchar(32) UNIQUE NOT NULL DEFAULT '',
    password varchar(64) NOT NULL DEFAULT '',
    session_id int(32) DEFAULT 0,
    created datetime DEFAULT CURRENT_TIMESTAMP
);

insert into gogame_db.user_table 
(username, password) VALUES
('Soichiro.Inatani', 'inatani'),
('inai17ibar', '17tani');
