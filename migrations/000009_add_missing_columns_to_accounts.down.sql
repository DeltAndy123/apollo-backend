ALTER TABLE accounts
    DROP COLUMN IF EXISTS is_deleted,
    DROP COLUMN IF EXISTS development;
