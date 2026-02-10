-- ═══════════════════════════════════════════════════
-- 007_admin_notifications.sql
-- Store admin push notifications for offline users
-- ═══════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS AdminNotification (
  NotificationId INT AUTO_INCREMENT PRIMARY KEY,
  TargetUserId   INT          NOT NULL,
  Type           VARCHAR(50)  NOT NULL,   -- 'money','event','stats','sickness','heal','revive'
  Payload        JSON         NOT NULL,   -- full details of the admin action
  Message        VARCHAR(500) NOT NULL,   -- human-readable summary
  CreatedAt      TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  ReadAt         TIMESTAMP    NULL,       -- NULL = unread
  FOREIGN KEY (TargetUserId) REFERENCES Users(UserId) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
