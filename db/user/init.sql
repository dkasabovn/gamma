CREATE TABLE org_users (
    id SERIAL PRIMARY KEY,
    perms INTEGER,
    user_uuid VARCHAR(36) UNIQUE,
    organization_id INTEGER
);

CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    org_uuid VARCHAR(36) UNIQUE,
    org_name VARCHAR(30),
    org_description TEXT
);

CREATE TABLE org_events (
    id SERIAL PRIMARY KEY,
    org_id INTEGER,
    event_uuid VARCHAR(36) UNIQUE,
    longitude REAL,
    latitude REAL,
    event_name VARCHAR(30),
    event_time TIMESTAMP WITH TIME ZONE NOT NULL,
    event_capacity SMALLINT,
    event_attending SMALLINT,
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
);

CREATE TABLE event_application (
    id SERIAL PRIMARY KEY,
    user_uuid VARCHAR(36) UNIQUE,
    event_id INTEGER,
    approved BOOLEAN,
    FOREIGN KEY(event_id)
        REFERENCES org_events(id)
);