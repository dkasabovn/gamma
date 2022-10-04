CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT NOT NULL,
    image_url TEXT NOT NULL,
    validated BOOLEAN NOT NULL
);

-- users that belong to an organization
CREATE TABLE org_users (
    id SERIAL PRIMARY KEY,
    policies_num INT NOT NULL,
    user_fk uuid NOT NULL,
    organization_fk uuid NOT NULL
);

CREATE TABLE organizations (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    org_name TEXT NOT NULL,
    city TEXT NOT NULL,
    org_image_url TEXT NOT NULL
);

CREATE TABLE events (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    event_name TEXT NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_location TEXT NOT NULL,
    event_description TEXT NOT NULL,
    event_image_url TEXT NOT NULL,
    organization_fk uuid NOT NULL
);
 
CREATE TABLE user_events (
    id SERIAL PRIMARY KEY,
    user_fk uuid NOT NULL,
    event_fk uuid NOT NULL,
    application_state INT NOT NULL
);

CREATE TABLE invites (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    expiration_date TIMESTAMP NOT NULL,
    capacity INT NOT NULL,
	user_fk uuid NOT NULL,
    org_fk uuid NOT NULL,
	entity_uuid uuid NOT NULL,
	entity_type int NOT NULL
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

ALTER TABLE user_events
    ADD CONSTRAINT unique_user_event_fks UNIQUE (user_fk, event_fk);

ALTER TABLE invites
	ADD CONSTRAINT fk_invites_org_user FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE invites
    ADD CONSTRAINT fk_invites_organization FOREIGN KEY (org_fk) REFERENCES organizations(id) ON DELETE CASCADE;