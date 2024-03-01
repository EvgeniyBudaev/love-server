CREATE TABLE profile_complaints (
                                    id BIGSERIAL NOT NULL PRIMARY KEY,
                                    profile_id BIGINT NOT NULL,
                                    complaint_user_id BIGINT NOT NULL,
                                    reason VARCHAR,
                                    created_at TIMESTAMP NOT NULL,
                                    updated_at TIMESTAMP NOT NULL,
                                    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);