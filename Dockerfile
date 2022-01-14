FROM alpine:latest
ENV BM_VERSION=x86-v0.2
RUN apk add --no-cache wget && apk add --no-cache unzip
RUN wget https://github.com/ilja2209/batch-manager-ui-client/releases/download/${BM_VERSION}/batch-manager-ui-client
RUN wget https://github.com/ilja2209/batch-manager-ui-client/releases/download/${BM_VERSION}/static.zip
RUN chmod 777 batch-manager-ui-client
RUN unzip static.zip
RUN rm -f static.zip
ENV SERVICE_PORT=8080
ENV DB_HOST=<db_host:db_port>
ENV DB_USER=<db_user>
ENV DB_PASSWORD=<db_psswd>
ENV DB_NAME=<db_name>
EXPOSE 8080/tcp
RUN echo $PWD
RUN ls -l 
ENTRYPOINT [  "./batch-manager-ui-client" ]
