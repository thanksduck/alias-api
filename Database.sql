-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) UNIQUE NOT NULL CHECK (username ~ '^[a-z][a-z0-9\-_\.]{3,}$'),
    name VARCHAR(64) NOT NULL CHECK (LENGTH(name) >= 4),
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    alias_count INTEGER DEFAULT 0,
    destination_count INTEGER DEFAULT 0,
    is_premium BOOLEAN DEFAULT FALSE,
    password VARCHAR(255) NOT NULL,
    provider VARCHAR(255),
    avatar VARCHAR(500),
    password_changed_at TIMESTAMP,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_email ON users (email);

-- Create rules table
CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    alias_email VARCHAR(255) UNIQUE NOT NULL,
    destination_email VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    comment VARCHAR(255),
    name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rules_user_id ON rules (user_id);

-- Create user_auth table
CREATE TABLE user_auth (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
     password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_auth_user_id ON user_auth (user_id);

-- Create destinations table with ON UPDATE CASCADE
CREATE TABLE destinations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    destination_email VARCHAR(255) NOT NULL,
    domain VARCHAR(255) NOT NULL,
    cloudflare_destination_id VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_destinations_user_id ON destinations (user_id);
CREATE INDEX idx_domain ON destinations (domain);

-- Create social_profiles table with ON UPDATE CASCADE
CREATE TABLE social_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    github VARCHAR(255),
    google VARCHAR(255),
    facebook VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_social_profiles_user_id ON social_profiles (user_id);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    type VARCHAR(20) DEFAULT 'credit',
    gateway VARCHAR(50) DEFAULT 'phonepe',
    txn_id TEXT UNIQUE NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(10) CHECK (status IN ('success', 'pending', 'failed')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);

CREATE TABLE credits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    balance BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_credits_user_id ON credits(user_id);
CREATE INDEX idx_credits_balance ON credits(balance);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    credit_id INTEGER REFERENCES credits(id) ON DELETE SET NULL ON UPDATE CASCADE,
    plan VARCHAR(10) CHECK (plan IN ('star', 'free', 'galaxy')),
    price INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'paused', 'cancelled'))
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
