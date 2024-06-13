CREATE TABLE IF NOT EXISTS Account (
    id SERIAL PRIMARY KEY,
    password varchar(255) NOT NULL,
    display_name varchar(255) NOT NULL UNIQUE,
    profile_url varchar(255),
    active boolean NOT NULL DEFAULT TRUE
);
