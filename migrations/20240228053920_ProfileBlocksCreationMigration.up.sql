CREATE TABLE profile_blocks(
                                   id BIGSERIAL NOT NULL PRIMARY KEY,
                                   profile_id BIGINT NOT NULL,
                                   blocked_user_id BIGINT NOT NULL,
                                   is_blocked BOOL NOT NULL,
                                   created_at TIMESTAMP NOT NULL,
                                   updated_at TIMESTAMP NOT NULL,
                                   CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);