
CREATE TABLE exams (
    id TEXT PRIMARY KEY,
    class_id TEXT NOT NULL,
    subject_id TEXT NOT NULL,
    name TEXT NOT NULL,
    exam_date TEXT NOT NULL,
    total_marks REAL NOT NULL,

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (class_id, subject_id, name, exam_date),

    CHECK (total_marks > 0),

    FOREIGN KEY (class_id)
        REFERENCES classes(id)
        ON DELETE RESTRICT,

    FOREIGN KEY (subject_id)
        REFERENCES subjects(id)
        ON DELETE RESTRICT
);
