-- Store items, payments, and user inventory
CREATE TABLE IF NOT EXISTS StoreItem (
    ItemId INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Description TEXT,
    Category VARCHAR(50) NOT NULL,
    Price INT NOT NULL,
    Currency VARCHAR(3) DEFAULT 'USD',
    Icon VARCHAR(10),
    Effect JSON,
    Active BOOLEAN DEFAULT TRUE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Payment (
    PaymentId INT AUTO_INCREMENT PRIMARY KEY,
    UserId INT NOT NULL,
    ItemId INT NOT NULL,
    Amount INT NOT NULL,
    Currency VARCHAR(3) DEFAULT 'USD',
    Status VARCHAR(20) DEFAULT 'pending',
    StripePaymentIntentId VARCHAR(255),
    ErrorMessage TEXT,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (UserId) REFERENCES Users(UserId),
    FOREIGN KEY (ItemId) REFERENCES StoreItem(ItemId)
);

CREATE TABLE IF NOT EXISTS UserInventory (
    InventoryId INT AUTO_INCREMENT PRIMARY KEY,
    UserId INT NOT NULL,
    ItemId INT NOT NULL,
    Quantity INT DEFAULT 1,
    AcquiredAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserId) REFERENCES Users(UserId),
    FOREIGN KEY (ItemId) REFERENCES StoreItem(ItemId),
    UNIQUE KEY unique_user_item (UserId, ItemId)
);

-- Seed store items
INSERT IGNORE INTO StoreItem (Name, Description, Category, Price, Icon, Effect) VALUES
('Premium Food',    'Delicious gourmet meal for your Tama',           'food',      299, '🍖', '{"stat": "hunger", "value": 50}'),
('Golden Toy',      'A shiny golden toy that boosts happiness',       'cosmetic',  499, '🎮', '{"stat": "boredom", "value": -40}'),
('Spa Day',         'Full spa treatment for your Tama',               'boost',     399, '🧖', '{"stat": "hygiene", "value": 60}'),
('Lucky Charm',     'Increases work satisfaction temporarily',        'accessory', 799, '🍀', '{"stat": "workSatis", "value": 0.2}'),
('Party Pack',      'Throw a party! Boosts social satisfaction',      'boost',     599, '🎉', '{"stat": "socialSatis", "value": 0.3}'),
('Energy Drink',    'Quick energy boost for your Tama',               'food',      199, '⚡', '{"stat": "hunger", "value": 30}'),
('Diamond Collar',  'Exclusive diamond collar accessory',             'cosmetic', 1499, '💎', '{"stat": "happiness", "value": 0.1}'),
('Work Briefcase',  'Professional briefcase for better work',         'accessory', 699, '💼', '{"stat": "workSatis", "value": 0.15}'),
('Healing Potion',  'Cures any sickness instantly',                   'medicine',  899, '🧪', '{"stat": "sickness", "value": "cure"}'),
('XP Booster',      'Double experience for 24 hours',                 'boost',     999, '🚀', '{"stat": "xp", "value": 2}');
