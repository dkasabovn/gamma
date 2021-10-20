CREATE TABLE "Users" (
  "UserID" SERIAL PRIMARY KEY,
  "UserUuid" varchar(36)
);

CREATE TABLE "OrgUsers" (
  "OrgUserID" SERIAL PRIMARY KEY,
  "PermissionsCode" int,
  "UserFk" int,
  "OrgFk" int
);

CREATE TABLE "Organizations" (
  "OrganizationID" SERIAL PRIMARY KEY,
  "Name" varchar(20),
  "Description" varchar(255),
  "OrganizationUuid" varchar(36)
);

CREATE TABLE "OrganizationEvents" (
  "OrganizationEventID" SERIAL PRIMARY KEY,
  "Name" varchar(20),
  "Latitude" decimal,
  "Longitude" decimal,
  "City" varchar(40),
  "State" varchar(20),
  "OrgFk" int,
  "Capacity" int,
  "Attending" int,
  "EventUuid" varchar(36)
);

CREATE TABLE "EventApplications" (
  "EventApplicationID" SERIAL PRIMARY KEY,
  "UserFk" int,
  "DateCreated" TIMESTAMP WITH TIME ZONE,
  "OrgEventFk" int
);

CREATE TABLE "UserEvents" (
  "UserEventID" SERIAL PRIMARY KEY,
  "OrgEventFk" int,
  "UserFk" int
);

ALTER TABLE "OrgUsers" ADD FOREIGN KEY ("UserFk") REFERENCES "Users" ("UserID");

ALTER TABLE "OrgUsers" ADD FOREIGN KEY ("OrgFk") REFERENCES "Organizations" ("OrganizationID");

ALTER TABLE "OrganizationEvents" ADD FOREIGN KEY ("OrgFk") REFERENCES "Organizations" ("OrganizationID");

ALTER TABLE "EventApplications" ADD FOREIGN KEY ("UserFk") REFERENCES "Users" ("UserID");

ALTER TABLE "EventApplications" ADD FOREIGN KEY ("OrgEventFk") REFERENCES "OrganizationEvents" ("OrganizationEventID");

ALTER TABLE "UserEvents" ADD FOREIGN KEY ("UserFk") REFERENCES "Users" ("UserID");

ALTER TABLE "UserEvents" ADD FOREIGN KEY ("OrgEventFk") REFERENCES "OrganizationEvents" ("OrganizationEventID");

CREATE UNIQUE INDEX "UserUuidIndex" ON "Users" ("UserUuid");