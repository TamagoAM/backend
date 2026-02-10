-- ═══════════════════════════════════════════════════
-- 006: Friend Request system + Chat Messages
-- Replaces old Friends table with FriendRequest flow
-- Adds ChatMessage table for persistent message history
-- ═══════════════════════════════════════════════════

-- Drop old Friends table and recreate with request/accept flow
DROP TABLE IF EXISTS ChatMessage;
DROP TABLE IF EXISTS Friends;

-- Friend requests with status: pending → accepted / declined
CREATE TABLE Friends (
    RequestId INT NOT NULL AUTO_INCREMENT,
    SenderID INT NOT NULL,
    ReceiverID INT NOT NULL,
    Status ENUM('pending','accepted','declined') NOT NULL DEFAULT 'pending',
    DateRequested DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DateResponded DATETIME,
    PRIMARY KEY (RequestId),
    UNIQUE KEY uq_sender_receiver (SenderID, ReceiverID),
    FOREIGN KEY (SenderID) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (ReceiverID) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_sender (SenderID),
    INDEX idx_receiver (ReceiverID),
    INDEX idx_status (Status)
) ENGINE=InnoDB;

-- Chat messages stored in MySQL for persistence
-- Redis handles real-time delivery, MySQL is the source of truth
CREATE TABLE ChatMessage (
    MessageId INT NOT NULL AUTO_INCREMENT,
    SenderID INT NOT NULL,
    ReceiverID INT NOT NULL,
    Body TEXT NOT NULL,
    SentAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ReadAt DATETIME,
    PRIMARY KEY (MessageId),
    FOREIGN KEY (SenderID) REFERENCES Users(UserId) ON DELETE CASCADE,
    FOREIGN KEY (ReceiverID) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_conversation (SenderID, ReceiverID, SentAt),
    INDEX idx_receiver_unread (ReceiverID, ReadAt)
) ENGINE=InnoDB;
