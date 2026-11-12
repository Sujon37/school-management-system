CREATE TABLE attendance (
    id TEXT PRIMARY KEY,

    student_id TEXT NOT NULL,
    attendance_date TEXT NOT NULL,

    status TEXT NOT NULL
        CHECK (
            status IN ('present', 'absent', 'late', 'excused')
        ),

    marked_by TEXT,

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (student_id, attendance_date),

    FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE RESTRICT,

    FOREIGN KEY (marked_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_attendance_date
ON attendance(attendance_date);