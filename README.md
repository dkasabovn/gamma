# What is gamma?

Gamma is IAM but for frats. Currently working on getting this wrapped up and deployed~

----
user sends email -> we get back **User struct**
We cannot give them a jwt yet because we don't know if they know their password but we have their email
We insert a random uuid -> user uuid entry in Redis
We send the base64 encoded uuid as a link to their email
On clicking we check redis if that map exists