-- 009_add_happiness_to_tama_stats.sql
-- Adds persistent happiness stat to Tama_stats

ALTER TABLE Tama_stats
ADD COLUMN Happiness DOUBLE DEFAULT 100 AFTER PersonalSatis;