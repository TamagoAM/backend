CREATE TABLE IF NOT EXISTS Users (
    UserId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    UserName VARCHAR(100) UNIQUE NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    PasswordHash VARCHAR(255) NOT NULL DEFAULT '',
    ClearanceLevel INT NOT NULL DEFAULT 0,
    Verified BOOLEAN NOT NULL DEFAULT FALSE,
    ProfilPicture VARCHAR(255),
    GamingTime INT DEFAULT 0,
    CreationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastConnectionDate DATETIME,
    PRIMARY KEY (UserId),
    INDEX idx_username (UserName),
    INDEX idx_email (Email)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Race (
    RaceId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) UNIQUE NOT NULL,
    `Desc` TEXT,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (RaceId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE Tama_stats (
    TamaStatId INT NOT NULL AUTO_INCREMENT,
    Fed INT DEFAULT 0,
    LastFed DATETIME,
    Played INT DEFAULT 0,
    LastPlayed DATETIME,
    Cleaned INT DEFAULT 0,
    LastCleaned DATETIME,
    Worked INT DEFAULT 0,
    LastWorked DATETIME,
    Hunger INT DEFAULT 0,
    Boredom INT DEFAULT 0,
    Hygiene INT DEFAULT 0,
    Money INT DEFAULT 0,
    CarAccident INT DEFAULT 0,
    WorkAccident INT DEFAULT 0,
    SocialSatis DOUBLE DEFAULT 0,
    WorkSatis DOUBLE DEFAULT 0,
    PersonalSatis DOUBLE DEFAULT 0,
    PRIMARY KEY (TamaStatId)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Tama (
    TamaId INT NOT NULL AUTO_INCREMENT,
    UserId INT NOT NULL,
    TamaStatsID INT NOT NULL,
    Name VARCHAR(100) NOT NULL,
    Sexe BOOLEAN,
    Race VARCHAR(100) NOT NULL,
    Sickness NVARCHAR(255),
    Birthday DATE,
    DeathDay DATE,
    CauseOfDeath NVARCHAR(255),
    Traits NVARCHAR(255),
    PRIMARY KEY (TamaId),
    FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (TamaStatsID) REFERENCES Tama_stats(TamaStatId) ON DELETE CASCADE,
    FOREIGN KEY (Race) REFERENCES Race(Name),
    INDEX idx_userid (UserId),
    INDEX idx_race (Race)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Friends (
    UserID INT NOT NULL,
    FriendID INT NOT NULL,
    DateBecameFriends DATE NOT NULL DEFAULT (CURRENT_DATE),
    PRIMARY KEY (UserID, FriendID),
    FOREIGN KEY (UserID) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (FriendID) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_userid (UserID),
    INDEX idx_friendid (FriendID)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Sponsor (
    SponsorId INT NOT NULL,
    SponsoredId INT NOT NULL,
    DateOfSponsor DATE NOT NULL DEFAULT (CURRENT_DATE),
    PRIMARY KEY (SponsorId, SponsoredId),
    FOREIGN KEY (SponsorId) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (SponsoredId) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_sponsorid (SponsorId),
    INDEX idx_sponsoredid (SponsoredId)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Sickness (
    SicknessId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    ExpirationDays INT,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (SicknessId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Trait (
    TraitId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (TraitId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Bonus (
    BonusId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Effet VARCHAR(255),
    PRIMARY KEY (BonusId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Malus (
    MalusId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Effet VARCHAR(255),
    PRIMARY KEY (MalusId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS Event (
    EventId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (EventId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS LifeChoices (
    LifeChoicesId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Traits NVARCHAR(255),
    PRIMARY KEY (LifeChoicesId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;
