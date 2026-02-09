-- Migration 006: Change DATE columns to DATETIME for precise timestamps
-- This fixes age calculation bugs caused by timezone differences when
-- only date (no time) was stored.

ALTER TABLE Tama MODIFY Birthday DATETIME;
ALTER TABLE Tama MODIFY DeathDay DATETIME;
