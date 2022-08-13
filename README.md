# What is gamma?

Gamma is GAMMA

# qr ref

User's qr contains their uuid & smoothstep time by 5 mins

When scanned by org user we send a request to the backend with this information

From org_user {
	policy number
	organization
	event
}

From user {
	uuid
	time
}

Server checks these conditions:

Is user attending event?

Is qr time less than 5 mins old

Is org_user able to scan qr codes

If all are true let user in