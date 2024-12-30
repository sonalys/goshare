CREATE TABLE participant (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP,

    CONSTRAINT participant_unique_email UNIQUE (email)
);