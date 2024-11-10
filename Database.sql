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
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_email ON users (email);
-- Create rules table with ON UPDATE CASCADE
CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    alias_email VARCHAR(255) NOT NULL UNIQUE,
    destination_email VARCHAR(255) NOT NULL REFERENCES destinations(destination_email),
    active BOOLEAN DEFAULT TRUE,
    comment VARCHAR(255),
    name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rules_user_id ON rules (user_id);
CREATE INDEX idx_rules_username ON rules (username);
CREATE INDEX idx_alias_email ON rules (alias_email);

-- Create destinations table with ON UPDATE CASCADE
CREATE TABLE destinations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    destination_email VARCHAR(255) NOT NULL UNIQUE,
    domain VARCHAR(255) NOT NULL,
    cloudflare_destination_id VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_destinations_user_id ON destinations (user_id);
CREATE INDEX idx_destinations_username ON destinations (username);
CREATE INDEX idx_destination_email ON destinations (destination_email);
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
CREATE INDEX idx_social_profiles_username ON social_profiles (username);

-- Create custom_domains table with ON UPDATE CASCADE
CREATE TABLE custom_domains (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    username VARCHAR(15) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX idx_custom_domains_user_id ON custom_domains (user_id);
CREATE INDEX idx_custom_domains_username ON custom_domains (username);

-- Create custom_domains_dns_records table
CREATE TABLE custom_domains_dns_records (
    id SERIAL PRIMARY KEY,
    custom_domain_id INTEGER REFERENCES custom_domains(id),
    cloudflare_id VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    name VARCHAR(20) NOT NULL,
    content VARCHAR(255) NOT NULL,
    ttl INTEGER NOT NULL,
    priority INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_custom_domains_dns_records_custom_domain_id ON custom_domains_dns_records (custom_domain_id);