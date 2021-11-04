CREATE TABLE Users (
  UserID SERIAL PRIMARY KEY,
  UserUuid varchar(36) NOT NULL
);

CREATE TABLE OrgUsers (
  OrgUserID SERIAL PRIMARY KEY,
  PermissionsCode int NOT NULL,
  UserFk int NOT NULL,
  OrgFk int NOT NULL
);

CREATE TABLE Organizations (
  OrganizationID SERIAL PRIMARY KEY,
  Name varchar(20) NOT NULL,
  Description varchar(255) NOT NULL,
  OrganizationUuid varchar(36) NOT NULL
);

CREATE TABLE OrganizationEvents (
  OrganizationEventID SERIAL PRIMARY KEY,
  Name varchar(20) NOT NULL,
  Latitude decimal NOT NULL,
  Longitude decimal NOT NULL,
  City varchar(40) NOT NULL,
  State varchar(20) NOT NULL,
  OrgFk int NOT NULL,
  Capacity int NOT NULL,
  Attending int NOT NULL,
  EventUuid varchar(36) NOT NULL
);

CREATE TABLE EventApplications (
  EventApplicationID SERIAL PRIMARY KEY,
  UserFk int NOT NULL,
  DateCreated TIMESTAMP WITH TIME ZONE NOT NULL,
  OrgEventFk int NOT NULL
);

CREATE TABLE UserEvents (
  UserEventID SERIAL PRIMARY KEY NOT NULL,
  OrgEventFk int NOT NULL,
  UserFk int NOT NULL
);

ALTER TABLE OrgUsers ADD FOREIGN KEY (UserFk) REFERENCES Users (UserID) ON DELETE CASCADE;

ALTER TABLE OrgUsers ADD FOREIGN KEY (OrgFk) REFERENCES Organizations (OrganizationID) ON DELETE CASCADE;

ALTER TABLE OrganizationEvents ADD FOREIGN KEY (OrgFk) REFERENCES Organizations (OrganizationID) ON DELETE CASCADE;

ALTER TABLE EventApplications ADD FOREIGN KEY (UserFk) REFERENCES Users (UserID) ON DELETE CASCADE;

ALTER TABLE EventApplications ADD FOREIGN KEY (OrgEventFk) REFERENCES OrganizationEvents (OrganizationEventID) ON DELETE CASCADE;

ALTER TABLE UserEvents ADD FOREIGN KEY (UserFk) REFERENCES Users (UserID) ON DELETE CASCADE;

ALTER TABLE UserEvents ADD FOREIGN KEY (OrgEventFk) REFERENCES OrganizationEvents (OrganizationEventID) ON DELETE CASCADE;

CREATE UNIQUE INDEX UserUuidIndex ON Users (UserUuid);

CREATE UNIQUE INDEX OrganizationEventIndex ON OrganizationEvents (EventUuid);

CREATE INDEX UserEventsIndex ON UserEvents (UserFk);

CREATE INDEX EventApplicationsIndex ON EventApplications (UserFk);