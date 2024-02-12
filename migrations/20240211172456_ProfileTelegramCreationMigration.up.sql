CREATE TABLE profile_telegram (
                                  id BIGSERIAL NOT NULL PRIMARY KEY,
                                  profile_id BIGINT,
                                  telegram_id BIGINT,
                                  username VARCHAR,
                                  first_name VARCHAR,
                                  last_name VARCHAR,
                                  language_code VARCHAR,
                                  allows_write_to_pm BOOL,
                                  query_id VARCHAR,
                                  CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);