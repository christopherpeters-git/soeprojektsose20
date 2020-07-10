CREATE DATABASE users;
USE users;
CREATE TABLE User (
    Id int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    PRIMARY KEY(Id)
);

CREATE TABLE Favorite_Videos (
    User_Username varchar(255) NOT NULL,
    VideoLink varchar(255) NOT NULL,
    PRIMARY KEY(User_Username),
    FOREIGN KEY (User_Username) REFERENCES User(Username)
);