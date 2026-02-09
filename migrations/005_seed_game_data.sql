-- ═══════════════════════════════════════════════════
-- 005_seed_game_data.sql
-- Seed Bonus, Malus, Sickness, Trait, Event, LifeChoices
-- ═══════════════════════════════════════════════════
-- Effet JSON format: {"stat":"<stat>","op":"<add|multiply|coeff>","value":<number>}
-- Duration: hours (NULL = permanent / until removed)
-- Bonus/Malus columns in other tables reference Bonus/Malus names (comma-separated)

-- ───────────────────────────────────────────────────
-- BONUS (positive effects)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Bonus (Name, `Desc`, Effet, Duration) VALUES
  -- Trait-sourced bonuses
  ('Sharp Mind',         'Increased work efficiency',                     '{"stat":"workSatis","op":"coeff","value":1.2}', NULL),
  ('Strong Body',        'Hunger decays slower',                          '{"stat":"hunger","op":"coeff","value":0.8}', NULL),
  ('Social Butterfly',   'Social satisfaction boost',                     '{"stat":"socialSatis","op":"add","value":10}', NULL),
  ('Iron Will',          'Boredom rises slower',                          '{"stat":"boredom","op":"coeff","value":0.85}', NULL),
  ('Green Thumb',        'Hygiene decays slower',                         '{"stat":"hygiene","op":"coeff","value":0.85}', NULL),
  ('Lucky Star',         'Car accident chance reduced',                   '{"stat":"carAccident","op":"coeff","value":0.5}', NULL),
  ('Penny Pincher',      'Work earns more money',                         '{"stat":"money","op":"multiply","value":1.3}', NULL),
  ('Quick Healer',       'Sickness duration halved',                      '{"stat":"sickness","op":"coeff","value":0.5}', NULL),
  ('Thick Skin',         'All decay rates reduced slightly',              '{"stat":"allDecay","op":"coeff","value":0.9}', NULL),
  ('Born Leader',        'Work satisfaction bonus',                       '{"stat":"workSatis","op":"add","value":15}', NULL),
  -- Event-sourced bonuses
  ('Windfall',           'Surprise money gain',                           '{"stat":"money","op":"add","value":100}', NULL),
  ('Community Spirit',   'Social boost from community event',             '{"stat":"socialSatis","op":"add","value":20}', 72),
  ('Insurance Payout',   'Money from insurance after disaster',           '{"stat":"money","op":"add","value":200}', NULL),
  ('Government Aid',     'Aid package during crisis',                     '{"stat":"money","op":"add","value":150}', NULL),
  ('Morale Boost',       'Temporary happiness lift',                      '{"stat":"personalSatis","op":"add","value":15}', 48),
  -- Life-choice bonuses
  ('Scholarship',        'Education pays off with work satisfaction',     '{"stat":"workSatis","op":"add","value":20}', NULL),
  ('Promotion',          'Career advance boosts money earning',           '{"stat":"money","op":"multiply","value":1.5}', NULL),
  ('Healthy Lifestyle',  'Better hygiene and hunger management',          '{"stat":"allDecay","op":"coeff","value":0.85}', NULL),
  ('True Love',          'Deep social satisfaction from relationship',    '{"stat":"socialSatis","op":"add","value":25}', NULL),
  ('Retirement Fund',    'Late-life money security',                      '{"stat":"money","op":"add","value":300}', NULL);

-- ───────────────────────────────────────────────────
-- MALUS (negative effects)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Malus (Name, `Desc`, Effet, Duration) VALUES
  -- Trait-sourced maluses
  ('Slow Learner',       'Work satisfaction reduced',                     '{"stat":"workSatis","op":"coeff","value":0.8}', NULL),
  ('Weak Constitution',  'Hunger decays faster',                          '{"stat":"hunger","op":"coeff","value":1.3}', NULL),
  ('Introvert Penalty',  'Social satisfaction penalty',                   '{"stat":"socialSatis","op":"add","value":-10}', NULL),
  ('Restless Spirit',    'Boredom rises faster',                          '{"stat":"boredom","op":"coeff","value":1.2}', NULL),
  ('Slob',              'Hygiene decays faster',                          '{"stat":"hygiene","op":"coeff","value":1.25}', NULL),
  ('Unlucky',           'Car accident chance increased',                  '{"stat":"carAccident","op":"coeff","value":1.5}', NULL),
  ('Spendthrift',       'Spending more on actions',                       '{"stat":"money","op":"multiply","value":0.7}', NULL),
  ('Fragile Health',    'Sickness duration longer',                       '{"stat":"sickness","op":"coeff","value":1.5}', NULL),
  ('Thin Skin',         'All decay rates increased slightly',             '{"stat":"allDecay","op":"coeff","value":1.15}', NULL),
  ('Pushover',          'Work satisfaction penalty',                      '{"stat":"workSatis","op":"add","value":-15}', NULL),
  -- Sickness maluses
  ('Fever Drain',       'Hunger drains faster while sick',                '{"stat":"hunger","op":"coeff","value":1.5}', NULL),
  ('Brain Fog',         'Boredom rises faster while sick',                '{"stat":"boredom","op":"coeff","value":1.4}', NULL),
  ('Chronic Pain',      'All satisfaction reduced',                       '{"stat":"personalSatis","op":"add","value":-20}', NULL),
  ('Weakened Immune',   'Hygiene decays much faster',                     '{"stat":"hygiene","op":"coeff","value":1.6}', NULL),
  ('Depression',        'Social satisfaction heavily reduced',             '{"stat":"socialSatis","op":"add","value":-25}', NULL),
  -- Event maluses
  ('Quake Damage',      'Money lost to earthquake damage',                '{"stat":"money","op":"add","value":-150}', NULL),
  ('Storm Injury',      'Hygiene and hunger hit by storm',                '{"stat":"hygiene","op":"add","value":-20}', NULL),
  ('Market Crash',      'Work satisfaction tanks during crisis',          '{"stat":"workSatis","op":"add","value":-30}', 168),
  ('Inflation',         'Everything costs more',                          '{"stat":"money","op":"multiply","value":0.6}', 120),
  ('Loneliness',        'Social isolation from life events',              '{"stat":"socialSatis","op":"add","value":-20}', NULL);

-- ───────────────────────────────────────────────────
-- SICKNESS
-- ───────────────────────────────────────────────────

-- Congenital (can be assigned at birth; 5-8 defects)
INSERT IGNORE INTO Sickness (Name, `Desc`, Type, Severity, ExpirationDays, CureCost, Bonus, Malus) VALUES
  ('Weak Lungs',         'Born with reduced respiratory capacity. Tires faster.',                    'congenital', 'mild',     NULL,  50,   NULL,              'Weak Constitution'),
  ('Sensitive Stomach',  'Prone to digestive issues. Hunger decays faster.',                         'congenital', 'mild',     NULL,  50,   NULL,              'Fever Drain'),
  ('Brittle Bones',      'Fractures more easily. Accidents are worse.',                              'congenital', 'moderate', NULL,  120,  NULL,              'Fragile Health,Unlucky'),
  ('Poor Eyesight',      'Difficulty seeing. Work efficiency reduced.',                              'congenital', 'mild',     NULL,  80,   NULL,              'Slow Learner'),
  ('Heart Murmur',       'Minor heart defect. Slightly increased all decay.',                        'congenital', 'moderate', NULL,  200,  NULL,              'Thin Skin'),
  ('Chronic Allergies',  'Persistent allergic reactions. Hygiene matters more.',                      'congenital', 'mild',     NULL,  60,   NULL,              'Slob'),
  ('Low Immunity',       'Immune system is weaker than normal. Gets sick easier.',                   'congenital', 'moderate', NULL,  150,  NULL,              'Fragile Health,Weakened Immune'),
  ('Clumsy Disposition', 'Accident-prone from birth. Higher accident rates.',                        'congenital', 'mild',     NULL,  40,   NULL,              'Unlucky'),

-- Acquired (caught during life)
  ('Common Cold',        'A mild cold. Sneezes and sniffles for a few days.',                        'acquired', 'mild',     3,     10,   NULL,              'Fever Drain'),
  ('Pixel Flu',          'Digital influenza. Moderate discomfort.',                                   'acquired', 'mild',     5,     25,   NULL,              'Fever Drain,Brain Fog'),
  ('Stomach Bug',        'Food poisoning. Hunger drains rapidly.',                                   'acquired', 'mild',     2,     15,   NULL,              'Fever Drain'),
  ('Rash',               'Itchy skin rash. Hygiene-related.',                                        'acquired', 'mild',     4,     20,   NULL,              'Weakened Immune'),
  ('Sprained Ankle',     'Twisted ankle from activity. Rest required.',                              'acquired', 'mild',     3,     30,   NULL,              'Chronic Pain'),
  ('Bronchitis',         'Lung infection. Moderate severity.',                                       'acquired', 'moderate', 7,     60,   NULL,              'Fever Drain,Chronic Pain'),
  ('Food Poisoning',     'Severe digestive distress.',                                               'acquired', 'moderate', 4,     40,   NULL,              'Fever Drain,Weakened Immune'),
  ('Broken Bone',        'Fractured limb. Extended recovery.',                                       'acquired', 'moderate', 14,    100,  NULL,              'Chronic Pain'),
  ('Pneumonia',          'Serious lung infection. Dangerous if untreated.',                           'acquired', 'severe',  10,    150,  NULL,              'Fever Drain,Chronic Pain,Weakened Immune'),
  ('Depression',         'Mental health crisis. Deep social and personal impact.',                    'acquired', 'severe',  21,    200,  NULL,              'Depression,Chronic Pain'),
  ('Cancer',             'Severe illness. Very expensive to treat. Long recovery.',                   'acquired', 'severe',  60,    NULL, NULL,              'Chronic Pain,Weakened Immune,Depression'),
  ('Heart Attack',       'Acute cardiac event. Immediate danger.',                                   'acquired', 'severe',  30,    500,  NULL,              'Thin Skin,Chronic Pain');

-- ───────────────────────────────────────────────────
-- TRAITS (20 total: 10 positive, 10 negative)
-- Each has BOTH a bonus AND a malus, but categorised
-- by whether the net effect is positive or negative.
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Trait (Name, `Desc`, Category, Bonus, Malus) VALUES
  -- ── Positive traits (net beneficial) ──
  ('Intelligent',     'Quick thinker but socially awkward.',                  'positive', 'Sharp Mind',        'Introvert Penalty'),
  ('Athletic',        'Strong and healthy but gets bored easily.',            'positive', 'Strong Body',       'Restless Spirit'),
  ('Charming',        'Great with others but lazy about self-care.',          'positive', 'Social Butterfly',  'Slob'),
  ('Disciplined',     'Focused and efficient but rigid and dull.',            'positive', 'Iron Will',         'Introvert Penalty'),
  ('Resourceful',     'Makes money stretch but takes risky shortcuts.',       'positive', 'Penny Pincher',     'Unlucky'),
  ('Resilient',       'Bounces back fast but emotionally distant.',           'positive', 'Quick Healer',      'Introvert Penalty'),
  ('Clean Freak',     'Impeccable hygiene but neurotic about it.',            'positive', 'Green Thumb',       'Restless Spirit'),
  ('Hard Worker',     'Dedicated to work but neglects social life.',          'positive', 'Born Leader',       'Introvert Penalty'),
  ('Lucky',           'Things just go right, but careless with money.',       'positive', 'Lucky Star',        'Spendthrift'),
  ('Tough',           'Thick-skinned and enduring but insensitive.',          'positive', 'Thick Skin',        'Introvert Penalty'),

  -- ── Negative traits (net detrimental) ──
  ('Lazy',            'Low effort but oddly content and relaxed.',            'negative', 'Iron Will',         'Slow Learner'),
  ('Glutton',         'Loves eating but never satisfied.',                    'negative', 'Strong Body',       'Weak Constitution'),
  ('Reckless',        'Bold and exciting but accident-prone.',                'negative', 'Social Butterfly',  'Unlucky'),
  ('Shy',             'Avoids conflict but misses opportunities.',            'negative', 'Green Thumb',       'Introvert Penalty'),
  ('Anxious',         'Hyper-aware but stressed and fragile.',                'negative', 'Quick Healer',      'Thin Skin'),
  ('Selfish',         'Good with money but terrible with people.',            'negative', 'Penny Pincher',     'Introvert Penalty'),
  ('Hot-Headed',      'Passionate but destructive temper.',                   'negative', 'Born Leader',       'Restless Spirit'),
  ('Pessimist',       'Cautious planner but draining to be around.',          'negative', 'Lucky Star',        'Depression'),
  ('Vain',            'Great self-care but narcissistic.',                    'negative', 'Green Thumb',       'Pushover'),
  ('Naive',           'Trusting and kind but easily exploited.',              'negative', 'Social Butterfly',  'Slow Learner');

-- ───────────────────────────────────────────────────
-- EVENTS
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Event (Name, `Desc`, Severity, Scope, MinStage, Bonus, Malus) VALUES
  -- Individual minor events
  ('Found Wallet',        'You found a wallet on the street! Lucky day.',                  'minor',        'individual', NULL,           'Windfall',         NULL),
  ('Flat Tire',           'Your tama got a flat tire. Small expense.',                      'minor',        'individual', 'teenage',      NULL,               'Quake Damage'),
  ('Birthday Gift',       'A friend sent a gift! Small morale boost.',                     'minor',        'individual', NULL,           'Morale Boost',     NULL),
  ('Lost Keys',           'Locked out! Frustrating but minor.',                            'minor',        'individual', 'childhood',    NULL,               'Loneliness'),
  ('Food Festival',       'Local food fest! Free meals.',                                  'minor',        'individual', NULL,           'Windfall',         NULL),
  ('Neighbor Dispute',    'Argument with a neighbor. Social stress.',                      'minor',        'individual', 'earlyAdult',   NULL,               'Loneliness'),
  ('Stray Pet',           'Found an adorable stray. Social boost!',                        'minor',        'individual', 'childhood',    'Community Spirit',  NULL),
  ('Minor Theft',         'Pickpocketed! Lost some coins.',                                'minor',        'individual', 'teenage',      NULL,               'Quake Damage'),

  -- Individual major events
  ('Car Accident',        'A serious car accident. Injury and costs.',                     'major',        'individual', 'teenage',      'Insurance Payout', 'Chronic Pain'),
  ('Job Loss',            'Fired from work. Financial and emotional hit.',                 'major',        'individual', 'earlyAdult',   NULL,               'Market Crash,Loneliness'),
  ('Inheritance',         'A distant relative left money behind.',                         'major',        'individual', 'midAdult',     'Windfall,Morale Boost', NULL),
  ('House Fire',          'Fire at home! Major property damage.',                          'major',        'individual', 'earlyAdult',   'Insurance Payout', 'Quake Damage,Storm Injury'),
  ('Lottery Win',         'Small lottery prize! Lucky day.',                               'major',        'individual', 'earlyAdult',   'Windfall,Morale Boost', NULL),
  ('Mugging',             'Attacked on the street. Lost money and health.',                'major',        'individual', 'teenage',      NULL,               'Quake Damage,Chronic Pain'),

  -- Global catastrophic events (admin-triggered or rare auto)
  ('Earthquake',          'A devastating earthquake hits the region.',                     'catastrophic', 'global',     NULL,           'Government Aid',   'Quake Damage,Storm Injury'),
  ('Tempest',             'A violent storm destroys infrastructure.',                       'catastrophic', 'global',     NULL,           'Government Aid',   'Storm Injury,Quake Damage'),
  ('Economic Crisis',     'The markets crash. Everyone feels the pain.',                   'catastrophic', 'global',     NULL,           'Government Aid',   'Market Crash,Inflation'),
  ('Pandemic',            'A disease sweeps through the population.',                      'catastrophic', 'global',     NULL,           'Government Aid',   'Weakened Immune,Depression'),
  ('War',                 'Conflict erupts. Resources become scarce.',                     'catastrophic', 'global',     NULL,           NULL,               'Market Crash,Inflation,Loneliness'),
  ('Famine',              'Crops fail. Food becomes scarce and expensive.',                 'catastrophic', 'global',     NULL,           'Government Aid',   'Inflation,Storm Injury');

-- ───────────────────────────────────────────────────
-- LIFE CHOICES
-- 1 month real time = full tama life
-- Infancy: 0-1 days (~2 choices)
-- Childhood: 1-4 days (~3 choices)
-- Teenage: 4-8 days (~4 choices)
-- Early Adult: 8-14 days (~4 choices)
-- Mid Adult: 14-20 days (~3 choices)
-- Late Adult: 20-26 days (~2 choices)
-- Elderly: 26-30 days (~2 choices)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO LifeChoices (Name, `Desc`, Stage, Rarity, ChoiceType, Traits, Bonus, Malus) VALUES
  -- ── Infancy (2 choices) ──
  ('First Steps',              'Your tama takes its first steps! Encourage or let be?',           'infancy',      'common',   'yesno',  NULL,            'Morale Boost',     NULL),
  ('Comfort Object',           'Choose a comfort object: blanket (calm) or ball (active)?',       'infancy',      'common',   'pool',   NULL,            'Morale Boost',     NULL),

  -- ── Childhood (3 choices) ──
  ('School Choice',            'Public school (social) or homeschool (smart)?',                   'childhood',    'common',   'pool',   NULL,            'Scholarship',      'Loneliness'),
  ('First Friend',             'Befriend the popular kid or the quiet kid?',                      'childhood',    'common',   'pool',   NULL,            'Community Spirit',  NULL),
  ('Playground Bully',         'Stand up to the bully or walk away?',                             'childhood',    'uncommon', 'yesno',  NULL,            'Morale Boost',     'Chronic Pain'),

  -- ── Teenage (4 choices) ──
  ('Part-time Job',            'Get a part-time job or focus on studies?',                         'teenage',      'common',   'pool',   'Hard Worker',   'Penny Pincher',     'Restless Spirit'),
  ('First Crush',              'Ask them out or admire from afar?',                                'teenage',      'common',   'yesno',  NULL,            'True Love',         'Loneliness'),
  ('Risky Dare',               'Your friends dare you to do something stupid.',                    'teenage',      'uncommon', 'yesno',  'Reckless',      'Community Spirit',  'Chronic Pain'),
  ('College or Trade',         'University track or vocational training?',                         'teenage',      'common',   'pool',   NULL,            'Scholarship',       NULL),

  -- ── Early Adulthood (4 choices) ──
  ('Career Path',              'Corporate ladder, freelance, or start a business?',                'earlyAdult',   'common',   'pool',   NULL,            'Promotion',         NULL),
  ('Move Cities',              'Stay in your hometown or move to the big city?',                   'earlyAdult',   'common',   'pool',   NULL,            'Windfall',          'Loneliness'),
  ('Serious Relationship',     'Commit to a relationship or stay independent?',                    'earlyAdult',   'common',   'yesno',  NULL,            'True Love',         'Spendthrift'),
  ('Investment Opportunity',   'Invest savings in stocks or play it safe?',                        'earlyAdult',   'uncommon', 'yesno',  NULL,            'Windfall',          'Market Crash'),

  -- ── Middle Adulthood (3 choices) ──
  ('Career Change',            'Opportunity for a radical career pivot.',                           'midAdult',     'rare',     'yesno',  NULL,            'Promotion',         'Market Crash'),
  ('Family or Career',         'Prioritize family time or chase the promotion?',                   'midAdult',     'common',   'pool',   NULL,            'True Love',         'Pushover'),
  ('Health Checkup',           'Invest in preventive healthcare? Costs money.',                    'midAdult',     'common',   'yesno',  NULL,            'Healthy Lifestyle',  NULL),

  -- ── Late Adulthood (2 choices) ──
  ('Retirement Plan',          'Retire early or keep working?',                                    'lateAdult',    'common',   'pool',   NULL,            'Retirement Fund',   'Loneliness'),
  ('Mentorship',               'Mentor a young tama? Social boost but time-consuming.',            'lateAdult',    'uncommon', 'yesno',  NULL,            'Community Spirit,Born Leader', NULL),

  -- ── Elderly (2 choices) ──
  ('Write Memoirs',            'Document your life story for future generations.',                  'elderly',      'common',   'yesno',  NULL,            'Morale Boost',      NULL),
  ('Legacy Gift',              'Leave your wealth to family or charity?',                           'elderly',      'common',   'pool',   NULL,            'Community Spirit',  NULL);

-- Update Race table to reference real Bonus/Malus names
UPDATE Race SET Bonus = 'Strong Body', Malus = 'Restless Spirit' WHERE Name = 'bear';
UPDATE Race SET Bonus = 'Sharp Mind', Malus = 'Introvert Penalty' WHERE Name = 'fox';
UPDATE Race SET Bonus = 'Green Thumb', Malus = 'Weak Constitution' WHERE Name = 'frog';
