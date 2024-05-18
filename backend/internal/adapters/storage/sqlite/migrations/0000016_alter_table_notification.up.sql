CREATE TABLE IF NOT EXISTS notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender INTEGER REFERENCES user(id),
    receiver INTEGER REFERENCES user(id),
    status TEXT CHECK( status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    type TEXT,
    message TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)