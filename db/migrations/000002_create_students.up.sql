CREATE TABLE students (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    admission_number TEXT NOT NULL UNIQUE,
    date_of_birth TEXT,
    gender TEXT,
    address TEXT,
    admission_date TEXT NOT NULL DEFAULT CURRENT_DATE,
    
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
