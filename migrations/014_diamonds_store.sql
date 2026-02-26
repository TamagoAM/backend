-- Add diamonds (premium currency) to Users table
ALTER TABLE Users ADD COLUMN IF NOT EXISTS Diamonds INT NOT NULL DEFAULT 0;

-- Remove old USD-based store items and re-seed with diamond items
DELETE FROM UserInventory;
DELETE FROM Payment;
DELETE FROM StoreItem;

-- Update StoreItem to use 'DIA' currency (diamonds)
ALTER TABLE StoreItem MODIFY Currency VARCHAR(3) DEFAULT 'DIA';

-- Seed diamond store items across all categories

-- ═══ Insurance ═══
INSERT INTO StoreItem (Name, Description, Category, Price, Currency, Icon, Effect) VALUES
('Car Insurance',       'Protects against *one* car accident. Consumed on use.', 'insurance', 80, 'DIA', '🚗', '{"type": "insurance", "target": "carAccident", "uses": 1}'),
('Work Insurance',      'Protects against *one* work accident. Consumed on use.', 'insurance', 80, 'DIA', '🏗️', '{"type": "insurance", "target": "workAccident", "uses": 1}'),
('Health Insurance',    'Instantly cures any sickness. Consumed on use.',       'insurance', 120, 'DIA', '🏥', '{"type": "insurance", "target": "sickness", "uses": 1}'),
('Full Coverage Plan',  'Protects against car, work, and sickness (1 each).',  'insurance', 250, 'DIA', '🛡️', '{"type": "insurance", "target": "all", "uses": 1}');

-- ═══ Stat Boosters ═══
INSERT INTO StoreItem (Name, Description, Category, Price, Currency, Icon, Effect) VALUES
('Premium Meal',        'Restores 50 hunger instantly.',                        'booster',  30,  'DIA', '🍖', '{"type": "boost", "stat": "hunger", "value": 50}'),
('Toy Box',             'Reduce boredom by 40 instantly.',                      'booster',  30,  'DIA', '🎮', '{"type": "boost", "stat": "boredom", "value": -40}'),
('Spa Day',             'Restores 60 hygiene instantly.',                       'booster',  35,  'DIA', '🧖', '{"type": "boost", "stat": "hygiene", "value": 60}'),
('Party Pack',          'Boost social satisfaction by 0.3.',                    'booster',  50,  'DIA', '🎉', '{"type": "boost", "stat": "socialSatis", "value": 0.3}'),
('Work Briefcase',      'Boost work satisfaction by 0.2.',                      'booster',  50,  'DIA', '💼', '{"type": "boost", "stat": "workSatis", "value": 0.2}'),
('Meditation Kit',      'Boost personal satisfaction by 0.25.',                 'booster',  45,  'DIA', '🧘', '{"type": "boost", "stat": "personalSatis", "value": 0.25}');

-- ═══ Multipliers ═══
INSERT INTO StoreItem (Name, Description, Category, Price, Currency, Icon, Effect) VALUES
('Money Booster x2',    'Double money earned for 24 hours.',                    'multiplier', 150, 'DIA', '💰', '{"type": "multiplier", "stat": "money", "factor": 2, "durationHours": 24}'),
('XP Booster x2',       'Double XP gained for 24 hours.',                      'multiplier', 150, 'DIA', '🚀', '{"type": "multiplier", "stat": "xp", "factor": 2, "durationHours": 24}'),
('Happiness Aura',      '+0.15 happiness bonus for 24 hours.',                 'multiplier', 100, 'DIA', '✨', '{"type": "multiplier", "stat": "happiness", "factor": 0.15, "durationHours": 24}');

-- ═══ Revival ═══
INSERT INTO StoreItem (Name, Description, Category, Price, Currency, Icon, Effect) VALUES
('Phoenix Feather',     'Revive a dead Tama with 50% stats restored.',         'revival',  500, 'DIA', '🔥', '{"type": "revival", "restore": 0.5}'),
('Second Chance Token', 'Prevents the next death (auto-triggers once).',        'revival',  350, 'DIA', '💫', '{"type": "deathShield", "uses": 1}');

-- ═══ Cosmetics ═══
INSERT INTO StoreItem (Name, Description, Category, Price, Currency, Icon, Effect) VALUES
('Golden Crown',        'A shiny golden crown for your Tama.',                 'cosmetic', 200, 'DIA', '👑', '{"type": "cosmetic", "slot": "head", "asset": "golden_crown"}'),
('Diamond Collar',      'Exclusive diamond collar accessory.',                 'cosmetic', 300, 'DIA', '💎', '{"type": "cosmetic", "slot": "neck", "asset": "diamond_collar"}'),
('Rainbow Wings',       'Majestic rainbow wings.',                             'cosmetic', 400, 'DIA', '🌈', '{"type": "cosmetic", "slot": "back", "asset": "rainbow_wings"}'),
('Pixel Sunglasses',    'Cool retro pixel sunglasses.',                        'cosmetic', 100, 'DIA', '🕶️', '{"type": "cosmetic", "slot": "face", "asset": "pixel_sunglasses"}');
