CREATE TABLE sections (
    id TEXT PRIMARY KEY,
    class_id TEXT NOT NULL,

    name TEXT NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 40,

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (class_id, name),

    CHECK (capacity > 0),

    FOREIGN KEY (class_id)
        REFERENCES classes(id)
        ON DELETE CASCADE
);