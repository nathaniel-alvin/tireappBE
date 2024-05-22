CREATE TABLE IF NOT EXISTS Account (
    id SERIAL PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    display_name varchar(255) NOT NULL,
    profile_url varchar(255),
    active boolean NOT NULL DEFAULT TRUE
);
