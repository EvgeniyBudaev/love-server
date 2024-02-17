CREATE TABLE profile_filters (
                                    id BIGSERIAL NOT NULL PRIMARY KEY,
                                    profile_id BIGINT NOT NULL,
                                    search_gender VARCHAR,
                                    looking_for VARCHAR NOT NULL,
                                    age_from INTEGER NOT NULL,
                                    age_to INTEGER NOT NULL,
                                    distance INTEGER NOT NULL,
                                    page INTEGER NOT NULL,
                                    size INTEGER NOT NULL,
                                    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);