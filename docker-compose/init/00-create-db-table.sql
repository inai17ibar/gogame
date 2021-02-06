create database gogame_db;

create table gogame_db.user_table
(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    email varchar(30) NOT NULL DEFAULT '',
    username varchar(32) NOT NULL DEFAULT '',
    created datetime DEFAULT CURRENT_TIMESTAMP
);

insert into gogame_db.user_table 
(id, email, username) VALUES
(101, 'foo@example.com', 'Foo bar'),
(212, 'bar@example.com', 'Hoge hoge');