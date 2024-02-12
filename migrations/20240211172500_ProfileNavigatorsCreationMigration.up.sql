CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE profile_navigators (
                                    id BIGSERIAL NOT NULL PRIMARY KEY,
                                    profile_id BIGINT NOT NULL,
                                    location geometry(Point,  4326),
                                    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);