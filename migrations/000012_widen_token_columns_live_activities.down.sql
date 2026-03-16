ALTER TABLE live_activities
    ALTER COLUMN apns_token TYPE character varying(100),
    ALTER COLUMN access_token TYPE character varying(64),
    ALTER COLUMN refresh_token TYPE character varying(64);

ALTER TABLE devices
    ALTER COLUMN apns_token TYPE character varying(100);
