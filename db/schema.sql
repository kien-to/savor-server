CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stores (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    pickup_time VARCHAR(255),
    distance VARCHAR(50),
    price DECIMAL(10,2),
    original_price DECIMAL(10,2),
    background_url TEXT,
    avatar_url TEXT,
    rating DECIMAL(3,1),
    reviews INTEGER,
    address TEXT,
    items_left INTEGER,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    image_url TEXT
);

CREATE TABLE store_highlights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id VARCHAR(36) REFERENCES stores(id) ON DELETE CASCADE,
    highlight TEXT NOT NULL
);

CREATE TABLE saved_stores (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    store_id VARCHAR(36) REFERENCES stores(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, store_id)
);

CREATE INDEX idx_stores_location ON stores(latitude, longitude);
CREATE INDEX idx_saved_stores_user ON saved_stores(user_id);

CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE INDEX idx_stores_title ON stores(title);
CREATE INDEX idx_stores_description ON stores(description);
CREATE INDEX idx_stores_address ON stores(address); 