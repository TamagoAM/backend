-- ═══════════════════════════════════════════════════
-- 008 — Life Choice Options + Tama Choice History
-- ═══════════════════════════════════════════════════

-- LifeChoiceOption: sub-options for 'pool' type life choices
CREATE TABLE IF NOT EXISTS LifeChoiceOption (
  OptionId       INT AUTO_INCREMENT PRIMARY KEY,
  LifeChoicesId  INT NOT NULL,
  Label          VARCHAR(100) NOT NULL,
  `Desc`         TEXT,
  Traits         NVARCHAR(255),
  Bonus          NVARCHAR(255),
  Malus          NVARCHAR(255),
  FOREIGN KEY (LifeChoicesId) REFERENCES LifeChoices(LifeChoicesId) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- TamaLifeChoiceHistory: persists which tama made which choice + which option
CREATE TABLE IF NOT EXISTS TamaLifeChoiceHistory (
  HistoryId      INT AUTO_INCREMENT PRIMARY KEY,
  TamaId         INT NOT NULL,
  LifeChoicesId  INT NOT NULL,
  ChosenOptionId INT DEFAULT NULL,
  Action         ENUM('accepted','rejected') NOT NULL DEFAULT 'accepted',
  CreatedAt      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (TamaId)         REFERENCES Tama(TamaId)              ON DELETE CASCADE,
  FOREIGN KEY (LifeChoicesId)  REFERENCES LifeChoices(LifeChoicesId) ON DELETE CASCADE,
  FOREIGN KEY (ChosenOptionId) REFERENCES LifeChoiceOption(OptionId) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═══════════════════════════════════════════════════
-- Seed options for all existing pool-type life choices
-- Each pool choice gets 2-3 sub-options with their own bonus/malus/traits
-- The parent choice's bonus/malus become fallback / are split across options
-- ═══════════════════════════════════════════════════

-- Helper: we use sub-selects to look up LifeChoicesId by Name

-- ── Infancy ──────────────────────────────────
-- Comfort Object: "blanket (calm) or ball (active)?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Comfort Object' LIMIT 1),
   '🧸 Blanket', 'A soft cozy blanket — calming and comforting.', NULL, 'Zen Master', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Comfort Object' LIMIT 1),
   '⚽ Ball', 'A bouncy ball — energetic and fun!', NULL, 'Morale Boost', 'Restless Spirit');

-- Feeding Method: "Breastmilk (healthy) or formula (convenient)?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Feeding Method' LIMIT 1),
   '🍼 Breastmilk', 'Natural and healthy — boosts immunity.', NULL, 'Vitamin Boost', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Feeding Method' LIMIT 1),
   '🧴 Formula', 'Convenient and reliable.', NULL, 'Morale Boost', NULL);

-- ── Childhood ────────────────────────────────
-- School Choice: "Public school (social) or homeschool (smart)?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'School Choice' LIMIT 1),
   '🏫 Public School', 'Lots of friends and social activities.', NULL, 'Community Spirit', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'School Choice' LIMIT 1),
   '📚 Homeschool', 'Focused learning at your own pace.', NULL, 'Scholarship', 'Loneliness');

-- First Friend: "Befriend the popular kid or the quiet kid?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'First Friend' LIMIT 1),
   '🌟 Popular Kid', 'Cool and well-known — instant social boost!', NULL, 'Community Spirit', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'First Friend' LIMIT 1),
   '🤫 Quiet Kid', 'Loyal and deep — a true friend.', NULL, 'Best Friend', NULL);

-- Team Sport: "Join a team sport or solo hobby?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Team Sport' LIMIT 1),
   '⚽ Team Sport', 'Teamwork, competition, and exercise!', NULL, 'Community Spirit', 'Chronic Pain'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Team Sport' LIMIT 1),
   '🎨 Solo Hobby', 'Personal expression and calm focus.', NULL, 'Morale Boost', NULL);

-- Candy or Veggies: "Prefer junk food or healthy snacks?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Candy or Veggies' LIMIT 1),
   '🍬 Candy', 'Sweet, sugary, and oh so tempting!', NULL, 'Sugar Rush', 'Picky Eater'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Candy or Veggies' LIMIT 1),
   '🥦 Veggies', 'Healthy and nutritious — body loves it.', NULL, 'Vitamin Boost', NULL);

-- ── Teenage ──────────────────────────────────
-- Part-time Job: "Get a part-time job or focus on studies?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Part-time Job' LIMIT 1),
   '💼 Get a Job', 'Earn money and learn responsibility.', 'Hard Worker', 'Penny Pincher', 'Restless Spirit'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Part-time Job' LIMIT 1),
   '📖 Focus on Studies', 'Books first — invest in your future.', NULL, 'Scholarship', NULL);

-- College or Trade: "University track or vocational training?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'College or Trade' LIMIT 1),
   '🎓 University', 'Higher education — opens many doors.', NULL, 'Scholarship', 'Student Debt'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'College or Trade' LIMIT 1),
   '🔧 Vocational', 'Hands-on skills — practical and fast.', NULL, 'Street Smart', NULL);

-- ── Early Adulthood ──────────────────────────
-- Career Path: "Corporate ladder, freelance, or start a business?" (3 options!)
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Career Path' LIMIT 1),
   '🏢 Corporate', 'Climb the ladder — stability and structure.', NULL, 'Promotion', 'Overworked'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Career Path' LIMIT 1),
   '💻 Freelance', 'Be your own boss — freedom with risk.', NULL, 'Side Hustle', 'Market Crash'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Career Path' LIMIT 1),
   '🚀 Start a Business', 'Dream big — high risk, high reward!', 'Hustler', 'Windfall', 'Market Crash');

-- Move Cities: "Stay in your hometown or move to the big city?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Move Cities' LIMIT 1),
   '🏡 Stay Home', 'Comfort and community — roots run deep.', NULL, 'Community Spirit', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Move Cities' LIMIT 1),
   '🌆 Big City', 'New opportunities — exciting but lonely at first.', NULL, 'Windfall', 'Loneliness');

-- Buy or Rent: "Buy a house or keep renting?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Buy or Rent' LIMIT 1),
   '🏠 Buy', 'Invest in a home — long-term security.', NULL, 'Home Owner', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Buy or Rent' LIMIT 1),
   '🏬 Rent', 'Stay flexible — freedom to move.', NULL, 'Morale Boost', 'Rent Hike');

-- ── Middle Adulthood ─────────────────────────
-- Family or Career: "Prioritize family time or chase the promotion?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Family or Career' LIMIT 1),
   '👨‍👩‍👧 Family First', 'Precious moments with loved ones.', NULL, 'True Love', 'Pushover'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Family or Career' LIMIT 1),
   '📈 Chase Promotion', 'Hard work pays off — but at what cost?', NULL, 'Promotion', 'Overworked');

-- Midlife Adventure: "Buy a sports car or go skydiving?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Midlife Adventure' LIMIT 1),
   '🏎️ Sports Car', 'Vroom vroom — feel alive again!', 'Daredevil', 'Adrenaline Junkie', 'Bad with Money'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Midlife Adventure' LIMIT 1),
   '🪂 Skydiving', 'Free-fall thrill — ultimate adrenaline!', 'Daredevil', 'Morale Boost', NULL);

-- ── Late Adulthood ───────────────────────────
-- Retirement Plan: "Retire early or keep working?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Retirement Plan' LIMIT 1),
   '🏖️ Retire Early', 'Enjoy life while you can!', NULL, 'Retirement Fund', 'Loneliness'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Retirement Plan' LIMIT 1),
   '💪 Keep Working', 'Stay active and productive.', NULL, 'Penny Pincher', NULL);

-- Downsize Home: "Move to a smaller place? Save money or feel lonely."
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Downsize Home' LIMIT 1),
   '🏘️ Downsize', 'Smaller space, bigger savings.', NULL, 'Retirement Fund', 'Loneliness'),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Downsize Home' LIMIT 1),
   '🏡 Stay Put', 'Keep the memories — stay in your home.', NULL, 'Morale Boost', NULL);

-- ── Elderly ──────────────────────────────────
-- Legacy Gift: "Leave your wealth to family or charity?"
INSERT IGNORE INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Legacy Gift' LIMIT 1),
   '👨‍👩‍👦 Family', 'Keep it in the family — they deserve it.', NULL, 'True Love', NULL),
  ((SELECT LifeChoicesId FROM LifeChoices WHERE Name = 'Legacy Gift' LIMIT 1),
   '🤝 Charity', 'Give back to the world — leave a lasting mark.', NULL, 'Community Spirit', NULL);
