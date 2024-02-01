CREATE TYPE gender_enum AS ENUM ('man', 'woman');
CREATE TYPE search_gender_enum AS ENUM ('man', 'woman', 'all');
CREATE TYPE looking_for_enum AS ENUM ('chat', 'dates', 'relationship', 'friendship', 'business', 'sex');

CREATE TABLE profiles
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    display_name VARCHAR NOT NULL,
    birthday TIMESTAMP NOT NULL,
    gender gender_enum NOT NULL,
    search_gender search_gender_enum NOT NULL,
    location VARCHAR NOT NULL,
    description VARCHAR,
    height INTEGER,
    weight INTEGER,
    looking_for looking_for_enum,
    is_deleted BOOL NOT NULL,
    is_blocked BOOL NOT NULL,
    is_premium BOOL NOT NULL,
    is_show_distance BOOL NOT NULL,
    is_invisible BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_online TIMESTAMP NOT NULL
)