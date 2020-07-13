CREATE DATABASE userdb;
USE userdb;
CREATE TABLE users (
    Id int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    PRIMARY KEY(Id)
);

CREATE TABLE videos (
    Video varchar(510) NOT NULL,
    Views int DEFAULT 0 NOT NULL,
    PRIMARY KEY(Video)
);

CREATE TABLE user_has_favorite_videos (
    Users_Username varchar(255) NOT NULL,
    Videos_Video varchar(510) NOT NULL,
    PRIMARY KEY(Users_Username),
    FOREIGN KEY (Users_Username) REFERENCES users(Username),
    FOREIGN KEY (Videos_Video) REFERENCES Videos(Video)
);