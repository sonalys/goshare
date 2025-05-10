CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP,

    CONSTRAINT users_unique_email UNIQUE (email)
);

CREATE TABLE ledgers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,

    FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE VIEW user_view AS
SELECT 
    u.*,
    COUNT(l.id) AS ledger_count
FROM 
    users u
LEFT JOIN 
    ledgers l ON u.id = l.created_by
GROUP BY 
    u.id;


CREATE TABLE ledger_members (
    user_id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    balance INTEGER NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (created_by) REFERENCES users (id),

    CONSTRAINT ledger_member_unique UNIQUE (ledger_id, user_id)
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    name TEXT NOT NULL,
    expense_date TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by UUID NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX expense_ledger_id_expense_date_desc ON expenses(ledger_id, expense_date DESC);

CREATE TABLE expense_records (
    id UUID PRIMARY KEY,
    expense_id UUID NOT NULL,

    record_type TEXT NOT NULL CHECK (record_type IN ('debt', 'settlement')),
    amount INTEGER NOT NULL,
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,

    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by UUID NOT NULL,

    FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE,
    FOREIGN KEY (from_user_id) REFERENCES users (id),
    FOREIGN KEY (to_user_id) REFERENCES users (id),
    FOREIGN KEY (created_by) REFERENCES users (id),

    CONSTRAINT expense_record_unique UNIQUE (id, expense_id)
);

