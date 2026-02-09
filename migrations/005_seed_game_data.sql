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

-- ═══════════════════════════════════════════════════
-- EXPANSION WAVE 2 — More variety & depth
-- ═══════════════════════════════════════════════════

-- ───────────────────────────────────────────────────
-- BONUS — Wave 2 (25 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Bonus (Name, `Desc`, Effet, Duration) VALUES
  -- Trait-sourced
  ('Night Owl',           'Thrives at night; boredom ticks slower',                 '{"stat":"boredom","op":"coeff","value":0.8}',      NULL),
  ('Iron Stomach',        'Can eat anything; hunger barely decays',                 '{"stat":"hunger","op":"coeff","value":0.7}',       NULL),
  ('Silver Tongue',       'Smooth talker; social satisfaction+',                    '{"stat":"socialSatis","op":"add","value":12}',     NULL),
  ('Neat Freak',          'Spotless habits; hygiene barely decays',                 '{"stat":"hygiene","op":"coeff","value":0.7}',      NULL),
  ('Adrenaline Junkie',   'Loves danger; accidents have less effect',               '{"stat":"carAccident","op":"coeff","value":0.6}',  NULL),
  ('Bookworm',            'Loves reading; big work satisfaction buff',              '{"stat":"workSatis","op":"add","value":18}',       NULL),
  ('Zen Master',          'Inner calm; personal satisfaction boost',                '{"stat":"personalSatis","op":"add","value":20}',   NULL),
  ('Street Smart',        'Knows how to get by; money multiplier',                  '{"stat":"money","op":"multiply","value":1.2}',     NULL),
  ('Tireless',            'Seemingly infinite energy; global decay down',           '{"stat":"allDecay","op":"coeff","value":0.8}',     NULL),
  ('Immune Boost',        'Sickness heals faster than normal',                      '{"stat":"sickness","op":"coeff","value":0.6}',     NULL),

  -- Event / Situation bonuses
  ('Sugar Rush',          'Short burst of energy from sweets',                      '{"stat":"boredom","op":"add","value":-15}',        6),
  ('Tax Refund',          'Annual tax return money',                                '{"stat":"money","op":"add","value":120}',          NULL),
  ('Good Samaritan',      'Helped someone; social & personal boost',                '{"stat":"socialSatis","op":"add","value":15}',     48),
  ('Spa Day',             'A relaxing spa visit; hygiene & satisfaction',            '{"stat":"hygiene","op":"add","value":25}',         24),
  ('Vitamin Boost',       'Took vitamins; slower hunger decay temporarily',         '{"stat":"hunger","op":"coeff","value":0.6}',       72),
  ('Meditation',          'Cleared the mind; boredom & decay reduced',              '{"stat":"allDecay","op":"coeff","value":0.75}',    48),
  ('Workout High',        'Post-exercise endorphins; personal satis up',            '{"stat":"personalSatis","op":"add","value":12}',   24),
  ('Happy Hour',          'Social drinks; social satisfaction spike',               '{"stat":"socialSatis","op":"add","value":18}',     12),
  ('Nap Time',            'Quick power nap restores energy',                        '{"stat":"boredom","op":"add","value":-20}',        6),
  ('Sunny Day',           'Beautiful weather lifts everyone''s mood',               '{"stat":"personalSatis","op":"add","value":10}',   24),

  -- Life-choice bonuses
  ('Side Hustle',         'Extra income from a side gig',                           '{"stat":"money","op":"add","value":80}',           NULL),
  ('Volunteer Work',      'Helping others gives purpose',                           '{"stat":"socialSatis","op":"add","value":20}',     NULL),
  ('Home Owner',          'Owning property; financial stability',                   '{"stat":"money","op":"multiply","value":1.15}',    NULL),
  ('World Traveler',      'Travel broadens the mind; all satis boost',              '{"stat":"personalSatis","op":"add","value":30}',   NULL),
  ('Best Friend',         'An unshakeable friendship; deep social bond',            '{"stat":"socialSatis","op":"add","value":30}',     NULL);

-- ───────────────────────────────────────────────────
-- MALUS — Wave 2 (25 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Malus (Name, `Desc`, Effet, Duration) VALUES
  -- Trait-sourced
  ('Insomniac',            'Can''t sleep; boredom rises faster',                    '{"stat":"boredom","op":"coeff","value":1.3}',      NULL),
  ('Picky Eater',          'Refuses most food; hunger decays faster',               '{"stat":"hunger","op":"coeff","value":1.4}',       NULL),
  ('Socially Awkward',     'Bad at conversations; social satis penalty',            '{"stat":"socialSatis","op":"add","value":-12}',    NULL),
  ('Messy',                'Can''t keep clean; hygiene decays faster',              '{"stat":"hygiene","op":"coeff","value":1.3}',      NULL),
  ('Accident Prone',       'Everything goes wrong; accident chance up',             '{"stat":"carAccident","op":"coeff","value":1.8}',  NULL),
  ('Procrastinator',       'Delays everything; work satisfaction down',             '{"stat":"workSatis","op":"add","value":-12}',      NULL),
  ('Overthinking',         'Paralysed by thought; personal satis down',             '{"stat":"personalSatis","op":"add","value":-15}',  NULL),
  ('Bad with Money',       'Can''t manage finances; money multiplier down',         '{"stat":"money","op":"multiply","value":0.75}',    NULL),
  ('Burnout',              'Total exhaustion; all decay faster',                    '{"stat":"allDecay","op":"coeff","value":1.25}',    NULL),
  ('Weak Immune',          'Gets sick longer than others',                          '{"stat":"sickness","op":"coeff","value":1.4}',     NULL),

  -- Sickness / Situation maluses
  ('Nausea',               'Feeling queasy; hunger drains',                         '{"stat":"hunger","op":"add","value":-10}',         24),
  ('Migraine',             'Throbbing headache; boredom spikes',                    '{"stat":"boredom","op":"add","value":20}',         12),
  ('Muscle Ache',          'Body is sore; hygiene drops',                           '{"stat":"hygiene","op":"add","value":-15}',        48),
  ('Anxiety Attack',       'Sudden anxiety; personal satis tanks',                  '{"stat":"personalSatis","op":"add","value":-25}',  6),
  ('Food Shortage',        'Can''t find good food; hunger decays fast',             '{"stat":"hunger","op":"coeff","value":1.6}',       72),
  ('Heatwave',             'Extreme heat; hygiene & boredom spike',                 '{"stat":"hygiene","op":"coeff","value":1.4}',      48),
  ('Rent Hike',            'Landlord raised rent; money squeezed',                  '{"stat":"money","op":"add","value":-80}',          NULL),
  ('Backstabbed',          'Betrayed by a friend; social satis destroyed',          '{"stat":"socialSatis","op":"add","value":-30}',    NULL),
  ('Overworked',           'Working too hard; personal satis drops',                '{"stat":"personalSatis","op":"add","value":-18}',  120),
  ('Scandal',              'Public embarrassment; social hit',                      '{"stat":"socialSatis","op":"add","value":-20}',    168),

  -- Life-choice maluses
  ('Student Debt',         'Education loan; money penalty',                         '{"stat":"money","op":"add","value":-100}',         NULL),
  ('Homesick',             'Missing home after moving; social & personal down',     '{"stat":"personalSatis","op":"add","value":-15}',  NULL),
  ('Toxic Relationship',   'Bad partner drags everything down',                     '{"stat":"socialSatis","op":"add","value":-25}',    NULL),
  ('Midlife Crisis',       'Questioning everything; personal satis tanks',          '{"stat":"personalSatis","op":"add","value":-30}',  NULL),
  ('Empty Nest',           'Kids left home; loneliness sets in',                    '{"stat":"socialSatis","op":"add","value":-20}',    NULL);

-- ───────────────────────────────────────────────────
-- SICKNESS — Wave 2 (20 more)
-- ───────────────────────────────────────────────────

-- More congenital (8 more)
INSERT IGNORE INTO Sickness (Name, `Desc`, Type, Severity, ExpirationDays, CureCost, Bonus, Malus) VALUES
  ('Color Blind',          'Cannot distinguish certain colors. Minor work impact.',             'congenital', 'mild',     NULL,  30,   NULL,              'Procrastinator'),
  ('Asthma',               'Chronic breathing difficulty. Activity limited.',                   'congenital', 'moderate', NULL,  100,  NULL,              'Weak Constitution,Burnout'),
  ('Flat Feet',            'Structural foot issue. Tires faster, accident-prone.',              'congenital', 'mild',     NULL,  40,   NULL,              'Accident Prone'),
  ('Dyslexia',             'Reading difficulties. Slower learning but creative.',               'congenital', 'mild',     NULL,  60,   'Zen Master',      'Slow Learner'),
  ('Lactose Intolerant',   'Can''t process dairy. Stomach issues.',                             'congenital', 'mild',     NULL,  20,   NULL,              'Picky Eater'),
  ('Hemophilia',           'Blood doesn''t clot well. Injuries are dangerous.',                 'congenital', 'severe',  NULL,  300,  NULL,              'Fragile Health,Chronic Pain'),
  ('Scoliosis',            'Curved spine. Chronic discomfort.',                                 'congenital', 'moderate', NULL,  150,  NULL,              'Chronic Pain,Burnout'),
  ('Photosensitivity',     'Extreme sensitivity to light. Outdoors is painful.',                'congenital', 'mild',     NULL,  70,   'Night Owl',       'Socially Awkward'),

-- More acquired (12 more)
  ('Migraine',             'Intense recurring headaches. Hard to focus.',                       'acquired', 'mild',     2,     20,   NULL,              'Migraine,Procrastinator'),
  ('Ear Infection',        'Painful ear infection. Hearing affected.',                          'acquired', 'mild',     5,     30,   NULL,              'Brain Fog'),
  ('Pink Eye',             'Contagious eye infection. Hygiene critical.',                       'acquired', 'mild',     4,     15,   NULL,              'Weakened Immune'),
  ('Appendicitis',         'Inflamed appendix. Surgery required.',                              'acquired', 'severe',  7,     250,  NULL,              'Chronic Pain,Fever Drain'),
  ('Diabetes',             'Blood sugar disorder. Ongoing management needed.',                  'acquired', 'moderate', NULL,  180,  NULL,              'Picky Eater,Burnout'),
  ('Insomnia',             'Can''t sleep properly. Exhaustion sets in.',                        'acquired', 'moderate', 10,    80,   NULL,              'Insomniac,Burnout'),
  ('Anxiety Disorder',     'Chronic anxiety. Affects all areas of life.',                       'acquired', 'moderate', 14,    120,  NULL,              'Anxiety Attack,Overthinking'),
  ('Mono',                 'Mononucleosis. Extreme fatigue for weeks.',                         'acquired', 'moderate', 21,    60,   NULL,              'Burnout,Fever Drain'),
  ('Tonsillitis',          'Swollen tonsils. Eating is painful.',                               'acquired', 'mild',     5,     35,   NULL,              'Fever Drain,Nausea'),
  ('Kidney Stones',        'Extremely painful. Requires medical attention.',                    'acquired', 'severe',  7,     200,  NULL,              'Chronic Pain,Nausea'),
  ('Stroke',               'Blood supply to brain blocked. Life-altering.',                     'acquired', 'severe',  45,    600,  NULL,              'Brain Fog,Chronic Pain,Burnout'),
  ('Addiction',             'Substance dependency. Drains money and health.',                    'acquired', 'severe',  NULL,  400,  NULL,              'Burnout,Bad with Money,Depression');

-- ───────────────────────────────────────────────────
-- TRAITS — Wave 2 (16 more: 8 positive, 8 negative)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Trait (Name, `Desc`, Category, Bonus, Malus) VALUES
  -- ── Positive traits ──
  ('Night Owl',       'Active at night but sluggish in daytime.',                   'positive', 'Night Owl',         'Procrastinator'),
  ('Foodie',          'Loves good food and rarely goes hungry, but picky.',         'positive', 'Iron Stomach',      'Picky Eater'),
  ('Smooth Talker',   'Can charm anyone but comes off as fake.',                    'positive', 'Silver Tongue',     'Socially Awkward'),
  ('Tidy',            'Everything in its place, but obsessively so.',               'positive', 'Neat Freak',        'Overthinking'),
  ('Daredevil',       'Lives for thrills but reckless with safety.',                'positive', 'Adrenaline Junkie', 'Accident Prone'),
  ('Bookish',         'Extremely well-read but a bit boring socially.',             'positive', 'Bookworm',          'Socially Awkward'),
  ('Stoic',           'Unshakeable calm but emotionally closed off.',               'positive', 'Zen Master',        'Introvert Penalty'),
  ('Hustler',         'Always making deals but shortcuts catch up.',                'positive', 'Street Smart',      'Unlucky'),

  -- ── Negative traits ──
  ('Sloth',           'Barely moves but weirdly good immune system.',               'negative', 'Immune Boost',      'Burnout'),
  ('Hypochondriac',   'Thinks they''re always sick but stays clean.',               'negative', 'Neat Freak',        'Overthinking'),
  ('People Pleaser',  'Loved by all but stretches too thin.',                       'negative', 'Social Butterfly',  'Burnout'),
  ('Paranoid',        'Always alert but trusts nobody.',                            'negative', 'Adrenaline Junkie', 'Socially Awkward'),
  ('Jealous',         'Driven by envy; works hard but resents others.',             'negative', 'Born Leader',       'Backstabbed'),
  ('Impulsive',       'Acts first, thinks later. Fun but costly.',                  'negative', 'Street Smart',      'Bad with Money'),
  ('Doormat',         'Never says no; loved but exploited.',                        'negative', 'Silver Tongue',     'Pushover'),
  ('Hoarder',         'Saves everything; financially secure but messy.',            'negative', 'Penny Pincher',     'Messy');

-- ───────────────────────────────────────────────────
-- EVENTS — Wave 2 (24 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Event (Name, `Desc`, Severity, Scope, MinStage, Bonus, Malus) VALUES
  -- Individual minor
  ('Good Hair Day',       'Feeling great! Small morale lift.',                                 'minor',        'individual', NULL,           'Morale Boost',     NULL),
  ('Free Lunch',          'Someone bought you lunch! Hunger satisfied.',                       'minor',        'individual', NULL,           'Sugar Rush',       NULL),
  ('Parking Ticket',      'Forgot the meter. Minor expense.',                                  'minor',        'individual', 'teenage',      NULL,               'Rent Hike'),
  ('Random Compliment',   'A stranger said something nice. Mood boost!',                       'minor',        'individual', NULL,           'Good Samaritan',   NULL),
  ('Power Outage',        'Electricity went out. Boredom ensues.',                             'minor',        'individual', NULL,           NULL,               'Migraine'),
  ('Got Ghosted',         'Date didn''t show up. Social ouch.',                                'minor',        'individual', 'teenage',      NULL,               'Loneliness'),
  ('Found Coupon',        'Surprise discount! Small savings.',                                 'minor',        'individual', NULL,           'Tax Refund',       NULL),
  ('Stubbed Toe',         'OUCH. The worst kind of minor injury.',                             'minor',        'individual', NULL,           NULL,               'Chronic Pain'),
  ('Surprise Visit',      'Old friend dropped by! Great social time.',                         'minor',        'individual', 'childhood',    'Happy Hour',       NULL),
  ('WiFi Down',           'Internet died. Pure suffering.',                                    'minor',        'individual', 'childhood',    NULL,               'Migraine'),

  -- Individual major
  ('Wedding',             'Getting married! Huge social & emotional event.',                   'major',        'individual', 'earlyAdult',   'True Love,Morale Boost', 'Spendthrift'),
  ('Divorce',             'Relationship ended. Devastating blow.',                             'major',        'individual', 'midAdult',     NULL,               'Loneliness,Backstabbed,Midlife Crisis'),
  ('New Baby',            'A tama baby is born! Joy and exhaustion.',                          'major',        'individual', 'earlyAdult',   'Morale Boost,Best Friend', 'Burnout,Spendthrift'),
  ('Work Promotion',      'Hard work paid off! Major career advance.',                         'major',        'individual', 'earlyAdult',   'Promotion,Windfall', NULL),
  ('Got Scammed',         'Lost money to a scam. Financial & emotional hit.',                  'major',        'individual', 'teenage',      NULL,               'Market Crash,Backstabbed'),
  ('Graduation',          'Completed education! Doors open.',                                  'major',        'individual', 'teenage',      'Scholarship,Morale Boost', NULL),
  ('Surgery',             'Major surgery needed. Costly but necessary.',                       'major',        'individual', NULL,           'Quick Healer',     'Chronic Pain,Rent Hike'),
  ('Best Friend Moves',   'Your best friend moved far away.',                                  'major',        'individual', 'childhood',    NULL,               'Loneliness,Homesick'),

  -- Global events
  ('Solar Eclipse',       'A rare eclipse! Everyone gets a morale boost.',                     'minor',        'global',     NULL,           'Morale Boost,Meditation', NULL),
  ('Heatwave',            'Scorching temperatures. Hygiene and energy suffer.',                 'major',        'global',     NULL,           NULL,               'Heatwave,Burnout'),
  ('Tech Boom',           'Technology sector explodes. Money flows.',                          'major',        'global',     NULL,           'Tax Refund,Side Hustle',  NULL),
  ('Housing Crisis',      'Rent skyrockets everywhere.',                                       'major',        'global',     NULL,           NULL,               'Rent Hike,Inflation'),
  ('World Cup',           'Global sporting event! Everyone''s excited.',                        'minor',        'global',     NULL,           'Community Spirit,Happy Hour', NULL),
  ('Alien Signal',        'Mysterious signal from space. World is buzzing!',                    'catastrophic', 'global',     NULL,           'Morale Boost,Community Spirit', 'Overthinking');

-- ───────────────────────────────────────────────────
-- LIFE CHOICES — Wave 2 (26 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO LifeChoices (Name, `Desc`, Stage, Rarity, ChoiceType, Traits, Bonus, Malus) VALUES
  -- ── Infancy (2 more) ──
  ('Feeding Method',            'Breastmilk (healthy) or formula (convenient)?',               'infancy',      'common',   'pool',   NULL,            'Vitamin Boost',     NULL),
  ('Crying Fit',                'Your tama won''t stop crying. Soothe or ignore?',             'infancy',      'uncommon', 'yesno',  NULL,            'Zen Master',        'Anxiety Attack'),

  -- ── Childhood (4 more) ──
  ('Pet Adoption',              'Adopt a pet? Social boost but responsibility.',                'childhood',    'common',   'yesno',  NULL,            'Best Friend',       'Spendthrift'),
  ('Team Sport',                'Join a team sport or solo hobby?',                             'childhood',    'common',   'pool',   NULL,            'Community Spirit',  'Chronic Pain'),
  ('Show and Tell',             'Present something at school? Risk embarrassment.',             'childhood',    'uncommon', 'yesno',  NULL,            'Morale Boost,Silver Tongue', 'Anxiety Attack'),
  ('Candy or Veggies',          'Prefer junk food or healthy snacks?',                         'childhood',    'common',   'pool',   NULL,            'Sugar Rush',        'Picky Eater'),

  -- ── Teenage (6 more) ──
  ('Social Media',              'Start posting online? Fame or drama.',                        'teenage',      'common',   'yesno',  NULL,            'Community Spirit',  'Scandal'),
  ('Driver License',            'Learn to drive or rely on others?',                            'teenage',      'common',   'yesno',  NULL,            'Street Smart',      'Accident Prone'),
  ('Study Abroad',              'Spend a semester in another country?',                         'teenage',      'rare',     'yesno',  NULL,            'World Traveler,Scholarship', 'Homesick,Student Debt'),
  ('House Party',               'Throw a huge party? Social boost or disaster.',                'teenage',      'uncommon', 'yesno',  'Reckless',      'Happy Hour,Community Spirit', 'Scandal,Spendthrift'),
  ('Volunteer Camp',            'Spend summer volunteering? Character building.',               'teenage',      'uncommon', 'yesno',  NULL,            'Volunteer Work,Good Samaritan', 'Loneliness'),
  ('Drop Out',                  'Leave school early to work? Risky but rewarding.',             'teenage',      'rare',     'yesno',  'Hustler',       'Street Smart,Side Hustle', 'Student Debt,Procrastinator'),

  -- ── Early Adulthood (6 more) ──
  ('Buy or Rent',               'Buy a house or keep renting?',                                 'earlyAdult',   'common',   'pool',   NULL,            'Home Owner',        'Rent Hike'),
  ('Start a Family',            'Have kids now or wait?',                                       'earlyAdult',   'common',   'yesno',  NULL,            'Best Friend,Morale Boost', 'Burnout,Spendthrift'),
  ('Grad School',               'Go back for a master''s degree?',                              'earlyAdult',   'uncommon', 'yesno',  NULL,            'Scholarship,Bookworm', 'Student Debt'),
  ('Side Gig',                  'Start a side business alongside your job?',                    'earlyAdult',   'common',   'yesno',  NULL,            'Side Hustle',       'Overworked'),
  ('Gym Membership',            'Join a gym? Health vs cost.',                                  'earlyAdult',   'common',   'yesno',  NULL,            'Workout High,Healthy Lifestyle', 'Spendthrift'),
  ('Adopt a Pet',               'Get a furry companion? Responsibility & joy.',                 'earlyAdult',   'common',   'yesno',  NULL,            'Best Friend',       'Spendthrift'),

  -- ── Middle Adulthood (4 more) ──
  ('Midlife Adventure',         'Buy a sports car or go skydiving?',                            'midAdult',     'uncommon', 'pool',   'Daredevil',     'Adrenaline Junkie,Morale Boost', 'Bad with Money'),
  ('Community Leader',          'Run for local office? Influence and stress.',                   'midAdult',     'rare',     'yesno',  NULL,            'Born Leader,Community Spirit', 'Overworked,Scandal'),
  ('Sabbatical',                'Take a year off work? Refreshing but costly.',                  'midAdult',     'uncommon', 'yesno',  NULL,            'Meditation,World Traveler', 'Market Crash'),
  ('Therapy',                   'Start seeing a therapist. Investment in self.',                 'midAdult',     'common',   'yesno',  NULL,            'Zen Master,Meditation', NULL),

  -- ── Late Adulthood (4 more) ──
  ('Downsize Home',             'Move to a smaller place? Save money or feel lonely.',           'lateAdult',    'common',   'pool',   NULL,            'Retirement Fund',   'Loneliness'),
  ('Grandchildren',             'Grandkids arrive! Joy and exhaustion.',                        'lateAdult',    'common',   'yesno',  NULL,            'Best Friend,Morale Boost', 'Burnout'),
  ('Learn New Skill',           'Pick up painting, music, or coding in retirement.',             'lateAdult',    'uncommon', 'yesno',  NULL,            'Bookworm,Morale Boost', NULL),
  ('Bucket List Trip',          'Take that dream vacation you always wanted.',                   'lateAdult',    'rare',     'yesno',  NULL,            'World Traveler,Happy Hour', 'Bad with Money'),

  -- ── Elderly (4 more) ──
  ('Final Will',                'Write your will. Peace of mind or family drama.',               'elderly',      'common',   'yesno',  NULL,            'Zen Master',        'Backstabbed'),
  ('Reunion',                   'Organize a family reunion? Social burst!',                     'elderly',      'uncommon', 'yesno',  NULL,            'Community Spirit,Happy Hour,Best Friend', NULL),
  ('Garden Tending',            'Spend days in the garden. Peaceful.',                          'elderly',      'common',   'yesno',  NULL,            'Green Thumb,Meditation', NULL),
  ('Pass the Torch',            'Train your replacement. Bittersweet.',                         'elderly',      'uncommon', 'yesno',  NULL,            'Born Leader,Volunteer Work', 'Loneliness');

-- ═══════════════════════════════════════════════════
-- EXPANSION WAVE 3 — Even more content
-- ═══════════════════════════════════════════════════

-- ───────────────────────────────────────────────────
-- BONUS — Wave 3 (15 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Bonus (Name, `Desc`, Effet, Duration) VALUES
  ('Coffee Kick',         'Caffeine boost; boredom drops temporarily',              '{"stat":"boredom","op":"add","value":-12}',        8),
  ('Holiday Spirit',      'Festive season lifts everyone''s spirits',               '{"stat":"personalSatis","op":"add","value":15}',   72),
  ('Pet Therapy',         'Cuddling a pet reduces stress',                          '{"stat":"personalSatis","op":"add","value":18}',   48),
  ('Gardening',           'Working with plants; hygiene & personal satis',          '{"stat":"hygiene","op":"add","value":15}',         48),
  ('Cooking Skill',       'Good cook; hunger managed efficiently',                  '{"stat":"hunger","op":"coeff","value":0.75}',      NULL),
  ('Music Lover',         'Listening to music; boredom drops',                      '{"stat":"boredom","op":"coeff","value":0.85}',     NULL),
  ('Early Bird',          'Wakes up early; gets more done',                         '{"stat":"workSatis","op":"add","value":12}',       NULL),
  ('Gratitude Journal',   'Writing daily; personal satisfaction up',                '{"stat":"personalSatis","op":"add","value":10}',   NULL),
  ('Team Player',         'Works well with others; social & work satis',            '{"stat":"workSatis","op":"add","value":10}',       NULL),
  ('Safety Net',          'Emergency fund; less financial stress',                  '{"stat":"money","op":"add","value":50}',           NULL),
  ('Rain Dance',          'Rainy day; cozy indoor vibes',                           '{"stat":"personalSatis","op":"add","value":8}',    24),
  ('Charity Donation',    'Gave to charity; feels good',                            '{"stat":"personalSatis","op":"add","value":12}',   48),
  ('Free Healthcare',     'Government health program; sickness cured faster',       '{"stat":"sickness","op":"coeff","value":0.4}',     168),
  ('Lottery Ticket',      'Won a small prize in the lottery',                       '{"stat":"money","op":"add","value":75}',           NULL),
  ('DIY Fix',             'Fixed something yourself; saved money & feels good',     '{"stat":"money","op":"add","value":40}',           NULL);

-- ───────────────────────────────────────────────────
-- MALUS — Wave 3 (15 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Malus (Name, `Desc`, Effet, Duration) VALUES
  ('Caffeine Crash',      'Coffee wore off; boredom spikes',                        '{"stat":"boredom","op":"add","value":15}',         8),
  ('Holiday Blues',        'Post-holiday depression; satis drops',                   '{"stat":"personalSatis","op":"add","value":-12}',  72),
  ('Allergic Reaction',   'Sudden allergy flare; hygiene & health hit',             '{"stat":"hygiene","op":"add","value":-20}',        24),
  ('Bad Haircut',         'Terrible haircut; personal satis dip',                   '{"stat":"personalSatis","op":"add","value":-8}',   48),
  ('Traffic Jam',         'Stuck in traffic; boredom & frustration',                '{"stat":"boredom","op":"add","value":15}',         6),
  ('Noisy Neighbors',     'Can''t sleep; boredom rises',                            '{"stat":"boredom","op":"coeff","value":1.2}',      72),
  ('Identity Theft',      'Someone stole your info; money lost',                    '{"stat":"money","op":"add","value":-120}',         NULL),
  ('Food Recall',         'Tainted food scare; hunger spikes',                      '{"stat":"hunger","op":"coeff","value":1.3}',       48),
  ('Tax Audit',           'IRS comes knocking; money and stress',                   '{"stat":"money","op":"add","value":-100}',         NULL),
  ('Social Media Drama',  'Online fight; social satisfaction tanks',                '{"stat":"socialSatis","op":"add","value":-15}',    48),
  ('Seasonal Blues',      'Winter doldrums; all satisfaction down',                 '{"stat":"personalSatis","op":"add","value":-10}',  168),
  ('Broken Phone',        'Phone died; social isolation & boredom',                 '{"stat":"socialSatis","op":"add","value":-10}',    24),
  ('Parking Fine',        'Got a parking ticket; minor money hit',                  '{"stat":"money","op":"add","value":-30}',          NULL),
  ('Sunburn',             'Stayed out too long; hygiene & pain',                    '{"stat":"hygiene","op":"add","value":-10}',        48),
  ('Ghosted',             'Someone stopped responding; social hit',                 '{"stat":"socialSatis","op":"add","value":-8}',     24);

-- ───────────────────────────────────────────────────
-- SICKNESS — Wave 3 (10 more acquired)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Sickness (Name, `Desc`, Type, Severity, ExpirationDays, CureCost, Bonus, Malus) VALUES
  ('Vertigo',              'Room won''t stop spinning. Balance disrupted.',                     'acquired', 'mild',     3,     25,   NULL,              'Brain Fog,Nausea'),
  ('Shingles',             'Painful nerve rash. Very uncomfortable.',                          'acquired', 'moderate', 14,    90,   NULL,              'Chronic Pain,Weakened Immune'),
  ('UTI',                  'Urinary tract infection. Painful and distracting.',                 'acquired', 'mild',     5,     25,   NULL,              'Chronic Pain'),
  ('Anemia',               'Low iron. Constant fatigue and weakness.',                         'acquired', 'moderate', 14,    70,   NULL,              'Burnout,Weak Constitution'),
  ('Concussion',           'Head injury. Requires rest and monitoring.',                       'acquired', 'moderate', 7,     120,  NULL,              'Brain Fog,Chronic Pain'),
  ('Carpal Tunnel',        'Wrist pain from repetitive motion. Work affected.',                'acquired', 'mild',     10,    50,   NULL,              'Procrastinator,Chronic Pain'),
  ('Chickenpox',           'Itchy spots everywhere. Very contagious.',                         'acquired', 'mild',     7,     30,   NULL,              'Weakened Immune,Fever Drain'),
  ('Panic Disorder',       'Recurring panic attacks. Debilitating.',                           'acquired', 'severe',  30,    250,  NULL,              'Anxiety Attack,Overthinking,Depression'),
  ('Back Injury',          'Slipped disc. Severe pain, limited mobility.',                     'acquired', 'moderate', 21,    160,  NULL,              'Chronic Pain,Burnout'),
  ('Sepsis',               'Blood infection. Life-threatening emergency.',                     'acquired', 'severe',  14,    700,  NULL,              'Fever Drain,Weakened Immune,Chronic Pain,Burnout');

-- ───────────────────────────────────────────────────
-- EVENTS — Wave 3 (16 more)
-- ───────────────────────────────────────────────────
INSERT IGNORE INTO Event (Name, `Desc`, Severity, Scope, MinStage, Bonus, Malus) VALUES
  -- Individual minor
  ('Rainy Day In',         'Stuck indoors. Cozy but boring.',                                  'minor',        'individual', NULL,           'Rain Dance',       NULL),
  ('Got a Raise',          'Small pay bump at work! Nice.',                                    'minor',        'individual', 'earlyAdult',   'Tax Refund',       NULL),
  ('Wardrobe Malfunction', 'Embarrassing outfit fail in public.',                              'minor',        'individual', 'teenage',      NULL,               'Bad Haircut'),
  ('Missed the Bus',       'Late for everything today.',                                       'minor',        'individual', 'childhood',    NULL,               'Traffic Jam'),
  ('Secret Admirer',       'Someone left a love note! Mood boost.',                            'minor',        'individual', 'teenage',      'Morale Boost',     NULL),
  ('DIY Disaster',         'Tried to fix something. Made it worse.',                           'minor',        'individual', 'earlyAdult',   NULL,               'Rent Hike'),

  -- Individual major
  ('Sued',                 'Legal trouble. Expensive and stressful.',                          'major',        'individual', 'earlyAdult',   NULL,               'Rent Hike,Anxiety Attack,Overworked'),
  ('Won a Prize',          'Entered a contest and won big!',                                   'major',        'individual', NULL,           'Windfall,Lottery Ticket,Morale Boost', NULL),
  ('Pet Died',             'Beloved pet passed away. Heartbreaking.',                          'major',        'individual', 'childhood',    NULL,               'Loneliness,Depression'),
  ('Identity Crisis',      'Who am I? Existential meltdown.',                                  'major',        'individual', 'teenage',      NULL,               'Midlife Crisis,Overthinking'),

  -- Global
  ('Meteor Shower',        'Beautiful celestial event. Morale for all!',                       'minor',        'global',     NULL,           'Morale Boost,Meditation', NULL),
  ('Trade War',            'International trade disputes. Prices rise.',                       'major',        'global',     NULL,           NULL,               'Inflation,Rent Hike'),
  ('Vaccine Breakthrough', 'Medical breakthrough. Sickness cured faster.',                     'major',        'global',     NULL,           'Free Healthcare',  NULL),
  ('Volcanic Eruption',    'A distant volcano erupts. Ash and chaos.',                         'catastrophic', 'global',     NULL,           'Government Aid',   'Storm Injury,Heatwave'),
  ('Cyber Attack',         'Global hack. Financial systems disrupted.',                        'catastrophic', 'global',     NULL,           NULL,               'Identity Theft,Market Crash'),
  ('Golden Age',           'A rare period of peace and prosperity.',                           'minor',        'global',     NULL,           'Holiday Spirit,Community Spirit,Tax Refund', NULL);

-- Update Race table to reference real Bonus/Malus names
UPDATE Race SET Bonus = 'Strong Body', Malus = 'Restless Spirit' WHERE Name = 'bear';
UPDATE Race SET Bonus = 'Sharp Mind', Malus = 'Introvert Penalty' WHERE Name = 'fox';
UPDATE Race SET Bonus = 'Green Thumb', Malus = 'Weak Constitution' WHERE Name = 'frog';
