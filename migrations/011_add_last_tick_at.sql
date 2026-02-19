-- ═══════════════════════════════════════════════════
-- 011_add_last_tick_at.sql — Track when background ticker last processed each tama
-- ═══════════════════════════════════════════════════

ALTER TABLE Tama_stats
    ADD COLUMN LastTickAt DATETIME DEFAULT NULL;

-- Initialise existing rows: set LastTickAt to the most recent activity timestamp
-- so the first ticker run doesn't apply a huge retroactive decay.
UPDATE Tama_stats ts
JOIN Tama t ON t.TamaStatsID = ts.TamaStatId
SET ts.LastTickAt = GREATEST(
    COALESCE(ts.LastFed,     t.Birthday),
    COALESCE(ts.LastPlayed,  t.Birthday),
    COALESCE(ts.LastCleaned, t.Birthday),
    COALESCE(ts.LastWorked,  t.Birthday)
)
WHERE ts.LastTickAt IS NULL
  AND t.DeathDay IS NULL;
