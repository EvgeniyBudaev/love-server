CREATE TABLE profiles
(
    id           BIGSERIAL NOT NULL PRIMARY KEY,
    display_name VARCHAR   NOT NULL UNIQUE
)
