CREATE USER auth_mantle_manager WITH PASSWORD 'dudde';
CREATE SCHEMA authmantledb;
GRANT ALL ON SCHEMA authmantledb TO auth_mantle_manager;

CREATE TABLE authmantledb."user_share" (
    id SERIAL PRIMARY KEY,
    share_name VARCHAR(100) NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb."user_share" TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.user_share_id_seq TO auth_mantle_manager;

CREATE TABLE authmantledb.role (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL,
    role_description VARCHAR(100) NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.role TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.role_id_seq TO auth_mantle_manager;

CREATE TABLE authmantledb.country (
    id SERIAL PRIMARY KEY,
    country_name VARCHAR(100) NOT NULL,
    country_alpha_name VARCHAR(3) NOT NULL,
    region_name VARCHAR(100) NOT NULL,
    region_alpha_name VARCHAR(3) NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb."country" TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.country_id_seq TO auth_mantle_manager;

CREATE TABLE authmantledb."user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    role INTEGER NOT NULL,
    email VARCHAR(150) NOT NULL,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    city VARCHAR(50) NOT NULL,
    state VARCHAR(50),
    country INTEGER NOT NULL,
    usr_share_id INTEGER DEFAULT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL,
    FOREIGN KEY (usr_share_id) REFERENCES authmantledb."user_share"(id),
    FOREIGN KEY (role) REFERENCES authmantledb."role"(id),
    FOREIGN KEY (country) REFERENCES authmantledb."country"(id)
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb."user" TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.user_id_seq TO auth_mantle_manager;

CREATE TABLE authmantledb.session (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(100) NOT NULL,
    session_data TEXT NOT NULL,
    session_expiry TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb."session" TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.session_id_seq TO auth_mantle_manager;