CREATE TABLE profiles
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    display_name VARCHAR NOT NULL,
    birthday DATE NOT NULL,
    gender VARCHAR,
    search_gender VARCHAR,
    location VARCHAR,
    description VARCHAR,
    height INTEGER,
    weight INTEGER,
    looking_for VARCHAR,
    is_deleted BOOL NOT NULL,
    is_blocked BOOL NOT NULL,
    is_premium BOOL NOT NULL,
    is_show_distance BOOL NOT NULL,
    is_invisible BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_online TIMESTAMP NOT NULL
);