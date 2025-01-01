CREATE TABLE ledgers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,

    FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE ledger_participants (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (created_by) REFERENCES users (id),

    CONSTRAINT ledger_participant_unique UNIQUE (ledger_id, user_id)
);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    name TEXT NOT NULL,
    parent_id UUID,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,

    FOREIGN KEY (parent_id) REFERENCES categories (id),
    FOREIGN KEY (ledger_id) REFERENCES ledgers (id),

    CONSTRAINT category_name_unique UNIQUE (id, name)
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY,
    category_id UUID,
    ledger_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    name TEXT NOT NULL,
    expense_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by UUID NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id),
    FOREIGN KEY (category_id) REFERENCES categories (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX expense_ledger_id_expense_date_desc ON expenses(ledger_id, expense_date DESC);

CREATE TABLE expense_payments (
    id UUID PRIMARY KEY,
    expense_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    payment_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    updated_by UUID NOT NULL,

    FOREIGN KEY (expense_id) REFERENCES expenses (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE TABLE ledger_records (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    expense_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by UUID NOT NULL,
    description TEXT NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (expense_id) REFERENCES expenses (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX ledger_record_ledger_id_created_at_desc ON ledger_records(ledger_id, created_at DESC);

CREATE TABLE ledger_participant_balances (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL,
    user_id UUID NOT NULL,
    last_timestamp TIMESTAMP NOT NULL,
    balance INTEGER NOT NULL,

    FOREIGN KEY (ledger_id) REFERENCES ledgers (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    
    CONSTRAINT ledger_participant_balance_unique UNIQUE (ledger_id, user_id)
);