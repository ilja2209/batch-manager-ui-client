docker build -t onef/batch-manager-ui .
docker run -p 8012:8080 --name batch-manager-ui-test onef/batch-manager-ui
