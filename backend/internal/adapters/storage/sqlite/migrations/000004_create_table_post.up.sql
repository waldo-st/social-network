CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER REFERENCES user(id),
    groupId INTEGER REFERENCES "group"(id),
    title TEXT,
    content TEXT,
    image TEXT,
    privacy TEXT CHECK(privacy='public' OR privacy='private' OR privacy='almost private'),
    createdAT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)