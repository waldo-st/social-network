CREATE TABLE IF NOT EXISTS "session" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER REFERENCES user(id),
    cookie TEXT,
    ttl TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)