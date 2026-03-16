ALTER TABLE accounts
    ALTER COLUMN access_token TYPE character varying(64),
    ALTER COLUMN refresh_token TYPE character varying(64);
