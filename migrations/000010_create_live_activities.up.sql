CREATE TABLE live_activities (
    id SERIAL PRIMARY KEY,
    apns_token character varying(100) UNIQUE,
    reddit_account_id character varying(32) DEFAULT ''::character varying,
    access_token character varying(64) DEFAULT ''::character varying,
    refresh_token character varying(64) DEFAULT ''::character varying,
    token_expires_at timestamp without time zone,
    thread_id character varying(32) DEFAULT ''::character varying,
    subreddit character varying(32) DEFAULT ''::character varying,
    next_check_at timestamp without time zone,
    expires_at timestamp without time zone,
    development boolean NOT NULL DEFAULT FALSE
);

CREATE INDEX live_activities_next_check_at_idx ON live_activities(next_check_at);
CREATE INDEX live_activities_expires_at_idx ON live_activities(expires_at);
