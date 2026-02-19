-- ═══════════════════════════════════════════════════
-- 012_night_cycle_notifications.sql
-- Day/night cycle + Expo Push Notifications
-- ═══════════════════════════════════════════════════

-- Store user timezone (IANA format, e.g. "Europe/Paris")
ALTER TABLE Users
  ADD COLUMN IF NOT EXISTS Timezone VARCHAR(64) DEFAULT 'Europe/Paris';

-- Lights on/off state for the night cycle
ALTER TABLE Tama_stats
  ADD COLUMN IF NOT EXISTS LightsOff BOOLEAN DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS LightsOffAt DATETIME NULL;

-- Expo push tokens (one user may have multiple devices)
CREATE TABLE IF NOT EXISTS PushToken (
  TokenId    INT AUTO_INCREMENT PRIMARY KEY,
  UserId     INT          NOT NULL,
  Token      VARCHAR(255) NOT NULL,
  Platform   VARCHAR(20)  NOT NULL DEFAULT 'android',  -- android, ios
  CreatedAt  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_user_token (UserId, Token),
  FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Notification log for smart escalation throttling
CREATE TABLE IF NOT EXISTS NotificationLog (
  LogId        INT AUTO_INCREMENT PRIMARY KEY,
  UserId       INT          NOT NULL,
  NotifType    VARCHAR(50)  NOT NULL,  -- 'low_hunger','low_happiness','sickness','friend_request','chat_message','bedtime','wake_up','death','event'
  SentAt       TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  EscalationN  INT          DEFAULT 1, -- 1st, 2nd, 3rd occurrence (for smart escalation)
  FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
  INDEX idx_notiflog_user_type (UserId, NotifType, SentAt)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
