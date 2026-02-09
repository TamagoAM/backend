-- ═══════════════════════════════════════════════════
-- Drop all tables in reverse dependency order
-- so foreign keys don't block the drops.
-- ═══════════════════════════════════════════════════
DROP TABLE IF EXISTS ActiveEvent;
DROP TABLE IF EXISTS LifeChoices;
DROP TABLE IF EXISTS Event;
DROP TABLE IF EXISTS Malus;
DROP TABLE IF EXISTS Bonus;
DROP TABLE IF EXISTS Trait;
DROP TABLE IF EXISTS Sickness;
DROP TABLE IF EXISTS Sponsor;
DROP TABLE IF EXISTS Friends;
DROP TABLE IF EXISTS Tama;
DROP TABLE IF EXISTS Tama_stats;
DROP TABLE IF EXISTS Race;
DROP TABLE IF EXISTS Users;

-- ═══════════════════════════════════════════════════
-- Recreate everything from scratch
-- ═══════════════════════════════════════════════════

CREATE TABLE Users (
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

CREATE TABLE Race (
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

CREATE TABLE Tama (
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

CREATE TABLE Friends (
    UserID INT NOT NULL,
    FriendID INT NOT NULL,
    DateBecameFriends DATE NOT NULL DEFAULT (CURRENT_DATE),
    PRIMARY KEY (UserID, FriendID),
    FOREIGN KEY (UserID) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (FriendID) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_userid (UserID),
    INDEX idx_friendid (FriendID)
) ENGINE=InnoDB;

CREATE TABLE Sponsor (
    SponsorId INT NOT NULL,
    SponsoredId INT NOT NULL,
    DateOfSponsor DATE NOT NULL DEFAULT (CURRENT_DATE),
    PRIMARY KEY (SponsorId, SponsoredId),
    FOREIGN KEY (SponsorId) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (SponsoredId) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_sponsorid (SponsorId),
    INDEX idx_sponsoredid (SponsoredId)
) ENGINE=InnoDB;

CREATE TABLE Sickness (
    SicknessId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Type ENUM('congenital','acquired','both') NOT NULL DEFAULT 'acquired',
    Severity ENUM('mild','moderate','severe') NOT NULL DEFAULT 'mild',
    ExpirationDays INT,
    CureCost INT,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (SicknessId),
    INDEX idx_name (Name),
    INDEX idx_type (Type)
) ENGINE=InnoDB;

CREATE TABLE Trait (
    TraitId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Category ENUM('positive','negative') NOT NULL DEFAULT 'positive',
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (TraitId),
    INDEX idx_name (Name),
    INDEX idx_category (Category)
) ENGINE=InnoDB;

CREATE TABLE Bonus (
    BonusId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Effet VARCHAR(255),
    Duration INT,
    PRIMARY KEY (BonusId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE Malus (
    MalusId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Effet VARCHAR(255),
    Duration INT,
    PRIMARY KEY (MalusId),
    INDEX idx_name (Name)
) ENGINE=InnoDB;

CREATE TABLE Event (
    EventId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Severity ENUM('minor','major','catastrophic') NOT NULL DEFAULT 'minor',
    Scope ENUM('individual','global') NOT NULL DEFAULT 'individual',
    MinStage ENUM('infancy','childhood','teenage','earlyAdult','midAdult','lateAdult','elderly') DEFAULT NULL,
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (EventId),
    INDEX idx_name (Name),
    INDEX idx_scope (Scope)
) ENGINE=InnoDB;

CREATE TABLE LifeChoices (
    LifeChoicesId INT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    `Desc` TEXT,
    Stage ENUM('infancy','childhood','teenage','earlyAdult','midAdult','lateAdult','elderly') NOT NULL DEFAULT 'childhood',
    Rarity ENUM('common','uncommon','rare') NOT NULL DEFAULT 'common',
    ChoiceType ENUM('pool','yesno') NOT NULL DEFAULT 'pool',
    Traits NVARCHAR(255),
    Bonus NVARCHAR(255),
    Malus NVARCHAR(255),
    PRIMARY KEY (LifeChoicesId),
    INDEX idx_name (Name),
    INDEX idx_stage (Stage)
) ENGINE=InnoDB;

-- Active events table: admin-triggered or system-triggered global/individual events
CREATE TABLE ActiveEvent (
    ActiveEventId INT NOT NULL AUTO_INCREMENT,
    EventId INT NOT NULL,
    TargetUserId INT,
    StartDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    EndDate DATETIME,
    TriggeredBy INT,
    IsGlobal BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (ActiveEventId),
    FOREIGN KEY (EventId) REFERENCES Event(EventId) ON DELETE CASCADE,
    FOREIGN KEY (TargetUserId) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (TriggeredBy) REFERENCES Users(UserId) ON DELETE SET NULL,
    INDEX idx_event (EventId),
    INDEX idx_target (TargetUserId),
    INDEX idx_global (IsGlobal),
    INDEX idx_dates (StartDate, EndDate)
) ENGINE=InnoDB;
