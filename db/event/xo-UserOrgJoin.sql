SELECT
    OrgUsers.PermissionsCode::INTEGER,
    Organizations.Name::VARCHAR(20),
    Organizations.Description::VARCHAR(255),
    Organizations.OrganizationUuid::VARCHAR(36)
FROM public.Users
INNER JOIN public.OrgUsers ON Users.UserID = OrgUsers.UserFk
INNER JOIN public.Organizations ON OrgUsers.OrgFk = Organizations.OrganizationID
WHERE Users.UserID = %%userId int%%