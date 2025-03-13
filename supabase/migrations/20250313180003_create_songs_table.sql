CREATE TABLE IF NOT EXISTS songs (
    id UUID PRIMARY KEY,
    pc_id INTEGER,
    admin TEXT,
    author TEXT,
    ccli_number INTEGER,
    copy_right TEXT,
    themes TEXT,
    title TEXT NOT NULL
);