
INSERT INTO users (uuid, email, first_name, last_name) VALUES ('dummy', 'dummy@dummy.dummy', 'bobert', 'lebowski');

INSERT INTO organizations (org_name, city) VALUES ('Big Brewery', 'Homieville');

INSERT INTO events (event_name, event_location, event_date, organization_fk, uuid) VALUES ('White Claw Sampling', 'Big Man University', '2017-03-14', 1, 'test');

INSERT INTO user_events (user_fk, event_fk) VALUES (1, 1);

INSERT INTO org_users (organization_fk) VALUES (1);