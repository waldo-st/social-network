CREATE TABLE IF NOT EXISTS chat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId INTEGER REFERENCES user(id),
    username TEXT,
    groupId INTEGER REFERENCES "group"(id),
    content TEXT NOT NULL,
    image TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)