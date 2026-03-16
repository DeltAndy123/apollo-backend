ALTER TABLE live_activities
    ALTER COLUMN apns_token TYPE text,
    ALTER COLUMN access_token TYPE text,
    ALTER COLUMN refresh_token TYPE text;

ALTER TABLE devices
    ALTER COLUMN apns_token TYPE text;
