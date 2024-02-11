CREATE TABLE profile_images (
                                id BIGSERIAL NOT NULL PRIMARY KEY,
                                profile_id BIGINT NOT NULL,
                                name VARCHAR,
                                url VARCHAR,
                                size INTEGER,
                                created_at TIMESTAMP NOT NULL,
                                updated_at TIMESTAMP NOT NULL,
                                is_deleted bool NOT NULL,
                                is_blocked bool NOT NULL,
                                is_primary bool NOT NULL,
                                is_private bool NOT NULL,
                                CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);