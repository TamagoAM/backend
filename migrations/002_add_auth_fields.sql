-- Add authentication fields to Users table (MySQL 8.x compatible)
-- Each ALTER is separate so failures on re-run (duplicate column) are non-fatal.
ALTER TABLE Users ADD COLUMN PasswordHash VARCHAR(255) NOT NULL DEFAULT '' AFTER Email;
ALTER TABLE Users ADD COLUMN ClearanceLevel INT NOT NULL DEFAULT 0 AFTER PasswordHash;
ALTER TABLE Users ADD COLUMN Verified BOOLEAN NOT NULL DEFAULT FALSE AFTER ClearanceLevel;
