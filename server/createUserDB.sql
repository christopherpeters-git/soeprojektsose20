CREATE DATABASE userdb;
USE userdb;
CREATE TABLE users (
    Id int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    Session_Id varchar(255) NOT NULL DEFAULT 0,
    PRIMARY KEY(Id)
);

CREATE TABLE videos (
    VideoTitle varchar(255) NOT NULL,
    Views int DEFAULT 0 NOT NULL,
    PRIMARY KEY(VideoTitle)
);

CREATE TABLE user_has_favorite_videos (
    Users_Username varchar(255) NOT NULL,
    Video varchar(650) NOT NULL,
    PRIMARY KEY(Users_Username,Video),
    FOREIGN KEY (Users_Username) REFERENCES users(Username)
);