docker run -d -p 8013:8080 \
-e DB_HOST=<db_host>:<port> \
-e DB_USER=<db_user> \
-e DB_PASSWORD=<db_password> \
-e DB_NAME=<db_name> \
-e BATCH_MANAGER_URL=<batch-manager-service url> \
--name batch-manager-ui-test ghcr.io/ilja2209/batch-manager-ui:<tag>
