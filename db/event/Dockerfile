FROM postgres:latest
ENV POSTGRES_PASSWORD nhPldb98Rt

COPY init.sql /docker-entrypoint-initdb.d/
RUN tee -a populate.sql /docker-entrypoint-initdb.d/init.sql