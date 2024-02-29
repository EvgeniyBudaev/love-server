CREATE TABLE profiles
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    session_id VARCHAR NOT NULL,
    display_name VARCHAR,
    birthday DATE,
    gender VARCHAR,
    location VARCHAR,
    description TEXT,
    height INTEGER NOT NULL,
    weight INTEGER NOT NULL,
    is_deleted BOOL NOT NULL,
    is_blocked BOOL NOT NULL,
    is_premium BOOL NOT NULL,
    is_show_distance BOOL NOT NULL,
    is_invisible BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_online TIMESTAMP NOT NULL
);