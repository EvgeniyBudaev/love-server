CREATE TYPE gender_enum AS ENUM ('man', 'woman');

CREATE TABLE profiles
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    display_name VARCHAR NOT NULL,
    birthday TIMESTAMP NOT NULL,
    gender gender_enum NOT NULL,
    location VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    is_deleted BOOL NOT NULL,
    is_blocked BOOL NOT NULL,
    is_premium BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_online TIMESTAMP NOT NULL
)