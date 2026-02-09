-- Seed the Race table with the three races supported by the app.
-- Uses INSERT IGNORE so re-runs are safe (duplicate key on Name is silently skipped).
INSERT IGNORE INTO Race (Name, `Desc`, Bonus, Malus) VALUES
  ('bear', 'A sturdy and cuddly bear Tama.', 'High resilience', 'Slow learner'),
  ('fox',  'A clever and swift fox Tama.',   'Quick reflexes',  'Easily bored'),
  ('frog', 'A cheerful and bouncy frog Tama.', 'Great hygiene', 'Low stamina');
