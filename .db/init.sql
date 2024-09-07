CREATE USER auth_mantle_manager WITH PASSWORD 'dudde';
CREATE SCHEMA IF NOT EXISTS authmantledb;
GRANT ALL ON SCHEMA authmantledb TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.us_share (
    id SERIAL PRIMARY KEY,
    share_name VARCHAR(100) NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.us_share TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.us_share_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.us_role (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL,
    role_description VARCHAR(100) NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.us_role TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.us_role_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.us_country (
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
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.us_country TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.us_country_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.us_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    role_id INTEGER NOT NULL,
    email VARCHAR(150) NOT NULL,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    city VARCHAR(50) NOT NULL,
    state VARCHAR(50),
    country_id INTEGER NOT NULL,
    share_id INTEGER DEFAULT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL,
    FOREIGN KEY (share_id) REFERENCES authmantledb.us_share(id),
    FOREIGN KEY (role_id) REFERENCES authmantledb.us_role(id),
    FOREIGN KEY (country_id) REFERENCES authmantledb.us_country(id)
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.us_user TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.us_user_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.us_session (
    id SERIAL PRIMARY KEY,
    session_id UUID NOT NULL,
    user_ref INTEGER NOT NULL,
    session_data TEXT NOT NULL,
    session_start TIMESTAMP NOT NULL,
    session_end TIMESTAMP NOT NULL,
    is_valid INTEGER NOT NULL DEFAULT 0,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(100) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    registered_by VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_ref) REFERENCES authmantledb.us_user(id)
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.us_session TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.us_session_id_seq TO auth_mantle_manager;

-- TABLES FOR INTERNALS
CREATE TABLE IF NOT EXISTS authmantledb.in_auth_code_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    auth_code UUID NOT NULL DEFAULT gen_random_uuid(),
    valid_until TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 hour',

    FOREIGN KEY (user_id) REFERENCES authmantledb.us_user(id)
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_auth_code_requests TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_auth_code_requests_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_audience (
    id SERIAL PRIMARY KEY,
    audience_name VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_audience TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_audience_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_grant_types (
    id SERIAL PRIMARY KEY,
    grant_type_name VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_grant_types TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_grant_types_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_scopes (
    id SERIAL PRIMARY KEY,
    scope_name VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_scopes TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_scopes_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_claims (
    id SERIAL PRIMARY KEY,
    claim_name VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_claims TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_claims_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_subject_types (
    id SERIAL PRIMARY KEY,
    subject_type_name VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_subject_types TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_subject_types_id_seq TO auth_mantle_manager;

CREATE TABLE IF NOT EXISTS authmantledb.in_supp_auth_allowed_redirects (
    id SERIAL PRIMARY KEY,
    redirect_uri VARCHAR(200) NOT NULL
);
GRANT INSERT, UPDATE, DELETE, SELECT ON authmantledb.in_supp_auth_allowed_redirects TO auth_mantle_manager;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE authmantledb.in_supp_auth_allowed_redirects_id_seq TO auth_mantle_manager;