-- Seed admin user (djlopez) with ClearanceLevel 2
-- Uses INSERT IGNORE so re-runs are safe (duplicate key on Email/UserName is silently skipped).
INSERT IGNORE INTO Users (Name, LastName, UserName, Email, PasswordHash, ClearanceLevel, Verified)
VALUES ('dj', 'lopez', 'djlopez', 'augustinmqn@gmail.com',
        '$2a$10$AwxFHHE/q3h7e4C38Q7TaOwHEQ9r/.iUTVgtftcMdOPmIAwMF/t6K',
        2, TRUE);

INSERT IGNORE INTO Users (Name, LastName, UserName, Email, PasswordHash, ClearanceLevel, Verified)
VALUES ('dj', 'lopette', 'djlopette', 'augustinmqn9@gmail.com',
        '$2a$10$AwxFHHE/q3h7e4C38Q7TaOwHEQ9r/.iUTVgtftcMdOPmIAwMF/t6K',
        0, TRUE);