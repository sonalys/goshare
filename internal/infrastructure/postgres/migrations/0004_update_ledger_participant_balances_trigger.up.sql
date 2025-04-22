-- Function to update ledger_participant_balances
CREATE OR REPLACE FUNCTION update_ledger_participant_balances()
RETURNS TRIGGER AS $$
BEGIN
    -- Update the balance for the user in the ledger_participant_balances table
    UPDATE ledger_participant_balances
    SET 
        balance = balance + NEW.amount,
        last_timestamp = NEW.created_at
    WHERE ledger_id = NEW.ledger_id AND user_id = NEW.user_id;

    -- If no row was updated, insert a new row
    IF NOT FOUND THEN
        INSERT INTO ledger_participant_balances (id, ledger_id, user_id, last_timestamp, balance)
        VALUES (gen_random_uuid(), NEW.ledger_id, NEW.user_id, NEW.created_at, NEW.amount);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function on insert into ledger_records
CREATE TRIGGER trigger_update_ledger_participant_balances
AFTER INSERT ON ledger_records
FOR EACH ROW
EXECUTE FUNCTION update_ledger_participant_balances();
