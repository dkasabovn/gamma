DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS org_users;

DROP TABLE IF EXISTS organizations;

DROP TABLE IF EXISTS events;
 
DROP TABLE IF EXISTS user_events;

DROP TABLE IF EXISTS invites;



ALTER TABLE org_users DROP CONSTRAINT IF EXISTS fk_orgs_users_organization;

ALTER TABLE org_users DROP CONSTRAINT IF EXISTS fk_org_users_user;

ALTER TABLE events DROP CONSTRAINT IF EXISTS fk_events_organization;

ALTER TABLE user_events DROP CONSTRAINT IF EXISTS fk_user_events_user;

ALTER TABLE user_events DROP CONSTRAINT IF EXISTS fk_user_events_event;

ALTER TABLE user_events DROP CONSTRAINT IF EXISTS unique_user_event_fks;

ALTER TABLE invites DROP CONSTRAINT IF EXISTS fk_invites_org_user;

ALTER TABLE invites DROP CONSTRAINT IF EXISTS fk_invites_organization;