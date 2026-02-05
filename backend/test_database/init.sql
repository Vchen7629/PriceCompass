CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions {
    id SERIAL PRIMARY KEY,
    username VARCHAR REFERENCES users(username),
    token VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
};

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR UNIQUE NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    last_checked_at TIMESTAMP,
    check_priority INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_watchlist (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    product_id INT NOT NULL REFERENCES products(id),
    added_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);

CREATE TABLE IF NOT EXISTS product_sources (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    platform VARCHAR NOT NULL,
    platform_product_id VARCHAR, -- ASIN, SKU, etc
    product_url VARCHAR,
    UNIQUE(platform, platform_product_id)
);

CREATE TABLE IF NOT EXISTS price_snapshots (
    id SERIAL PRIMARY KEY,
    product_source_id INT REFERENCES product_sources(id),
    price DECIMAL(10, 2),
    currency VARCHAR DEFAULT 'USD',
    in_stock BOOLEAN,
    checked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_price_snapshots_time ON price_snapshots(product_source_id, checked_at DESC);
