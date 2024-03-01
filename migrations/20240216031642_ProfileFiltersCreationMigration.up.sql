CREATE TABLE profile_filters (
                                    id BIGSERIAL NOT NULL PRIMARY KEY,
                                    profile_id BIGINT NOT NULL,
                                    search_gender VARCHAR,
                                    looking_for VARCHAR,
                                    age_from INTEGER,
                                    age_to INTEGER,
                                    distance INTEGER,
                                    page INTEGER,
                                    size INTEGER,
                                    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);