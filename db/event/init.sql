CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    image_url TEXT NOT NULL
);

-- users that belong to an organization
CREATE TABLE org_users (
    id SERIAL PRIMARY KEY,
    policies_num INT NOT NULL,
    user_fk INT NOT NULL,
    organization_fk INT NOT NULL
);

CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    org_name TEXT NOT NULL,
    city TEXT NOT NULL,
    org_image_url TEXT NOT NULL
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_location TEXT NOT NULL,
    event_image_url TEXT NOT NULL,
    organization_fk INT
);

-- Events that users have been accepted to
CREATE TABLE user_events (
    id SERIAL PRIMARY KEY,
    user_fk INT NOT NULL,
    event_fk INT NOT NULL
);

-- Contains event applications 
CREATE TABLE event_applications (
    id SERIAL PRIMARY KEY,
    user_fk INT NOT NULL,
    event_fk INT NOT NULL
);

-- Contains users that are whitelisted to any organization event
CREATE TABLE organization_whitelist (
    id SERIAL PRIMARY KEY,
    organization_fk INT NOT NULL,
    user_fk INT NOT NULL
);

-- Policy Json
CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    use_limit INT NOT NULL,
    policy_json JSON NOT NULL
);

ALTER TABLE org_users
    ADD CONSTRAINT fk_orgs_users_organization FOREIGN KEY (organization_fk) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE org_users
    ADD CONSTRAINT fk_org_users_user FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE events
    ADD CONSTRAINT fk_events_organization FOREIGN KEY (organization_fk) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE user_events
    ADD CONSTRAINT fk_user_events_user FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_events
    ADD CONSTRAINT fk_user_events_event FOREIGN KEY (event_fk) REFERENCES events(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS index_users_on_uuid ON users USING btree(uuid);

CREATE INDEX IF NOT EXISTS index_organizations_on_uuid ON organizations USING btree(uuid);
