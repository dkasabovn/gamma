# What is gamma?

Gamma is a an application that seeks to help organize limited capacity impromptu social events by allowing users to apply for events, and for organizers to accept users to events based upon their application / profile.

Hosts will then be able to verify the status of users pre-entry by using a custom QR-like scanning system at the entrance to the event.

# How is it progressing?

The backend is missing API routes and certian custom Queries. However, with the aid of XO they should be done quickly.

# How will this be hosted?

This will be hosted on DigitalOcean Kubernetes Clusters due to the low cost and good performance.

# Why is this public?

There are currently no API Keys / Passwords that will be risky for us to expose. DB password is local, and the image isn't pushed to a container repository. Until we have a frontend, the backend will be public.
