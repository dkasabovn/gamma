DELETE FROM "public"."Organizations";
DELETE FROM "public"."Users";

BEGIN;
WITH "CreatedOrg" AS (
    INSERT INTO "public"."Organizations" ("Name", "Description", "OrganizationUuid")
    VALUES ('SigmaChiV2', 'Grant Kitlowskis version of Sigma Chi.', 'orgGkitt')
    RETURNING *
), createduser AS (
    INSERT INTO "public"."Users" ("UserUuid")
    VALUES ('userGkitt')
    RETURNING "UserID"
)
INSERT INTO "public"."OrgUsers" ("PermissionsCode", "UserFk", "OrgFk")
VALUES (128, (SELECT "UserID" FROM createduser), (SELECT "OrganizationID" FROM "CreatedOrg"))
RETURNING *;
COMMIT;