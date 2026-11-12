
CREATE TABLE marks (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    exam_id TEXT NOT NULL,
    marks_obtained REAL NOT NULL,

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (student_id, exam_id),

    CHECK (marks_obtained >= 0),

    FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE RESTRICT,

    FOREIGN KEY (exam_id)
        REFERENCES exams(id)
        ON DELETE RESTRICT
);
