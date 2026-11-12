CREATE TABLE teachers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    employee_number TEXT NOT NULL UNIQUE,
    date_of_birth TEXT,
    gender TEXT,
    address TEXT,
    
    hire_date TEXT NOT NULL DEFAULT CURRENT_DATE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);