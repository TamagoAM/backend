-- ═══════════════════════════════════════════════════
-- 010_stat_history.sql — Periodic stat snapshots for graph display
-- ═══════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS StatHistory (
    HistoryId   INT NOT NULL AUTO_INCREMENT,
    TamaId      INT NOT NULL,
    -- Gauges (0-100)
    Hunger      INT NOT NULL DEFAULT 0,
    Boredom     INT NOT NULL DEFAULT 0,
    Hygiene     INT NOT NULL DEFAULT 0,
    -- Economy
    Money       INT NOT NULL DEFAULT 0,
    -- Satisfaction (0-100 float)
    SocialSatis   DOUBLE NOT NULL DEFAULT 0,
    WorkSatis     DOUBLE NOT NULL DEFAULT 0,
    PersonalSatis DOUBLE NOT NULL DEFAULT 0,
    Happiness     DOUBLE NOT NULL DEFAULT 0,
    -- Cumulative counters
    Fed         INT NOT NULL DEFAULT 0,
    Played      INT NOT NULL DEFAULT 0,
    Cleaned     INT NOT NULL DEFAULT 0,
    Worked      INT NOT NULL DEFAULT 0,
    -- Event counters
    CarAccident  INT NOT NULL DEFAULT 0,
    WorkAccident INT NOT NULL DEFAULT 0,
    -- Trigger context
    Trigger     VARCHAR(50) NOT NULL DEFAULT 'periodic',
    -- Timestamp
    RecordedAt  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (HistoryId),
    FOREIGN KEY (TamaId) REFERENCES Tama(TamaId) ON DELETE CASCADE,
    INDEX idx_tama_time (TamaId, RecordedAt),
    INDEX idx_trigger (Trigger)
) ENGINE=InnoDB;
