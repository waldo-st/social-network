CREATE TABLE IF NOT EXISTS event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    groupId INTEGER REFERENCES groud(id),
    creatorId INTERGER REFERENCES user(id),
    title TEXT,
    description TEXT,
    option TEXT CHECK(option ='going' OR option ='not going'),
    createdAT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)