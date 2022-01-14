docker run -d -p 8013:8080 \
-e DB_HOST=<db_host>:<port> \
-e DB_USER=<db_user> \
-e DB_PASSWORD=<db_password> \
-e DB_NAME=<db_name> \
--name batch-manager-ui-test ilja2209/batch-manager-ui
