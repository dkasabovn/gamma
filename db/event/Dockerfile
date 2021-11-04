FROM postgres:latest
ENV POSTGRES_PASSWORD nhPldb98Rt

COPY populate.sql ./
COPY init.sql /docker-entrypoint-initdb.d/
RUN cat populate.sql >> /docker-entrypoint-initdb.d/init.sql
RUN rm populate.sql