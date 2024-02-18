CREATE TABLE profiles
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id VARCHAR,
    display_name VARCHAR,
    birthday DATE,
    gender VARCHAR,
    location VARCHAR,
    description VARCHAR,
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