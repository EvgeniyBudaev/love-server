CREATE TABLE profile_telegram (
                                  id BIGSERIAL NOT NULL PRIMARY KEY,
                                  profile_id BIGINT NOT NULL,
                                  telegram_id BIGINT NOT NULL UNIQUE,
                                  username VARCHAR NOT NULL UNIQUE,
                                  first_name VARCHAR NOT NULL,
                                  last_name VARCHAR,
                                  language_code VARCHAR NOT NULL,
                                  allows_write_to_pm BOOL NOT NULL,
                                  query_id VARCHAR NOT NULL UNIQUE,
                                  CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);