CREATE DATABASE users;
USE users;
CREATE TABLE User (
    Id int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    PRIMARY KEY(Id)
);

CREATE TABLE Videos (
    Video varchar(510) NOT NULL,
    Views int DEFAULT 0 NOT NULL,
    PRIMARY KEY(Video)
);

CREATE TABLE User_has_favorite_Videos (
    User_Username varchar(255) NOT NULL,
    Videos_Video varchar(510) NOT NULL,
    PRIMARY KEY(User_Username),
    FOREIGN KEY (User_Username) REFERENCES User(Username),
    FOREIGN KEY (Videos_Video) REFERENCES Videos(Video)
);

