-- Add authentication fields to Users table (MySQL 8.x compatible)
-- Columns already present in 001_init.sql for fresh installs.
-- These ALTERs only matter for databases created before auth fields existed.
-- Duplicate-column errors are caught by the migrate runner and logged as warnings.
ALTER TABLE Users ADD COLUMN PasswordHash VARCHAR(255) NOT NULL DEFAULT '' AFTER Email;
ALTER TABLE Users ADD COLUMN ClearanceLevel INT NOT NULL DEFAULT 0 AFTER PasswordHash;
ALTER TABLE Users ADD COLUMN Verified BOOLEAN NOT NULL DEFAULT FALSE AFTER ClearanceLevel;
ALTER TABLE Users ADD INDEX idx_clearance (ClearanceLevel);
