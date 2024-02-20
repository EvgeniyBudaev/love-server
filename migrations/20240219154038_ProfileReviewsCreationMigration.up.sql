CREATE TABLE profile_reviews (
                                 id BIGSERIAL NOT NULL PRIMARY KEY,
                                 profile_id BIGINT NOT NULL,
                                 message TEXT,
                                 rating DECIMAL(3,  1),
                                 has_deleted BOOL NOT NULL,
                                 has_edited BOOL NOT NULL,
                                 created_at TIMESTAMP NOT NULL,
                                 updated_at TIMESTAMP NOT NULL,
                                 CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);