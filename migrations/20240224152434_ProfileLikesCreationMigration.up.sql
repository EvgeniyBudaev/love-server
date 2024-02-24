CREATE TABLE profile_likes (
                                 id BIGSERIAL NOT NULL PRIMARY KEY,
                                 profile_id BIGINT NOT NULL,
                                 human_id BIGINT NOT NULL,
                                 is_liked BOOL NOT NULL,
                                 created_at TIMESTAMP NOT NULL,
                                 updated_at TIMESTAMP NOT NULL,
                                 CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);