-- Create user tablef
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP without time zone DEFAULT NOW() NOT NULL
);

CREATE TABLE messages (
    id SERIAL,
    sender_id INT REFERENCES users(id) ON DELETE CASCADE,
    recipient_id INT REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'sent',
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);