-- +migrate Up
CREATE TABLE IF NOT EXISTS validators (
    id SERIAL PRIMARY KEY,
    pubkey TEXT UNIQUE NOT NULL,
    blockchain TEXT NOT NULL,
    blockchain_network TEXT NOT NULL,
    status TEXT NOT NULL,
    client TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
