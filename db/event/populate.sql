
INSERT INTO users (uuid, email, first_name, password_hash, last_name, username) VALUES ('dummy', 'dummy@dummy.dummy', 'bobert', '213123123', 'lebowski', 'test');

INSERT INTO organizations (org_name, city, uuid) VALUES ('Big Brewery', 'Homieville', 'test');

INSERT INTO events (event_name, event_location, event_date, organization_fk, uuid) VALUES ('White Claw Sampling', 'Big Man University', '2017-03-14', 1, 'test');

INSERT INTO user_events (user_fk, event_fk) VALUES (1, 1);

INSERT INTO org_users (organization_fk, policies_num) VALUES (1, 1);