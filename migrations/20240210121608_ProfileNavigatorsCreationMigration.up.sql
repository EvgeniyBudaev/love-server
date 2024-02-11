CREATE TABLE profile_navigators (
                                    id BIGSERIAL NOT NULL PRIMARY KEY,
                                    profile_id BIGINT NOT NULL,
                                    latitude VARCHAR,
                                    longitude VARCHAR,
                                    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);