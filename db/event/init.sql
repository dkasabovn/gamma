CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    org_user_fk INT
);

CREATE TABLE org_users (
    id SERIAL PRIMARY KEY,
    organization_fk INT
);

CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    org_name TEXT NOT NULL,
    city TEXT NOT NULL,
    uuid TEXT NOT NULL
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_name TEXT NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_location TEXT NOT NULL,
    uuid TEXT NOT NULL,
    organization_fk INT
);

CREATE TABLE user_events (
    id SERIAL PRIMARY KEY,
    user_fk INT,
    event_fk INT
);

CREATE TABLE user_event_invites (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    valid BOOLEAN NOT NULL,
    event_uuid TEXT NOT NULL
);

ALTER TABLE users
    ADD CONSTRAINT fk_users_org_user FOREIGN KEY (org_user_fk) REFERENCES org_users(id) ON DELETE CASCADE;

ALTER TABLE org_users
    ADD CONSTRAINT fk_orgs_users_organization FOREIGN KEY (organization_fk) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE events
    ADD CONSTRAINT fk_events_organization FOREIGN KEY (organization_fk) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE user_events
    ADD CONSTRAINT fk_user_events_user FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_events
    ADD CONSTRAINT fk_user_events_event FOREIGN KEY (event_fk) REFERENCES events(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS index_users_on_uuid ON users USING btree(uuid);

CREATE INDEX IF NOT EXISTS inex_user_invite_on_uuid ON user_event_invites USING btree(uuid);
