CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    ledger_count INTEGER NOT NULL,

    CONSTRAINT unique_user_email UNIQUE (email)
);

CREATE TABLE ledgers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID NOT NULL,

    CONSTRAINT fk_ledger_created_by FOREIGN KEY (created_by) REFERENCES users (id)
);


CREATE TABLE ledger_members (
    user_id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID NOT NULL,
    balance INTEGER NOT NULL,

    CONSTRAINT fk_ledger_member_ledger FOREIGN KEY (ledger_id) REFERENCES ledgers (id) ON DELETE CASCADE,
    CONSTRAINT fk_ledger_member_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_ledger_member_created_by FOREIGN KEY (created_by) REFERENCES users (id),

    CONSTRAINT unique_ledger_member UNIQUE (ledger_id, user_id)
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    name TEXT NOT NULL,
    expense_date TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID NOT NULL,

    CONSTRAINT fk_expense_ledger FOREIGN KEY (ledger_id) REFERENCES ledgers (id) ON DELETE CASCADE,
    CONSTRAINT fk_expense_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_expense_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX expense_ledger_id_expense_date_desc ON expenses(ledger_id, expense_date DESC);

CREATE TABLE expense_records (
    id UUID PRIMARY KEY,
    expense_id UUID NOT NULL,

    record_type TEXT NOT NULL CHECK (record_type IN ('debt', 'settlement')),
    amount INTEGER NOT NULL,
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID NOT NULL,

    CONSTRAINT fk_expense_record_expense FOREIGN KEY (expense_id) REFERENCES expenses (id) ON DELETE CASCADE,
    CONSTRAINT fk_expense_record_from_user FOREIGN KEY (from_user_id) REFERENCES users (id),
    CONSTRAINT fk_expense_record_to_user FOREIGN KEY (to_user_id) REFERENCES users (id),
    CONSTRAINT fk_expense_record_created_by FOREIGN KEY (created_by) REFERENCES users (id),

    CONSTRAINT unique_expense_record UNIQUE (id, expense_id)
);

