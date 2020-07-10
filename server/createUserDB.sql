CREATE DATABASE users;
CREATE TABLE user (
    id int,
    username varchar(20),
    passwordHash varchar(300),
)
CREATE TABLE userinformations (
    id int,
    name varchar(50),
    videos varchar(500),
)