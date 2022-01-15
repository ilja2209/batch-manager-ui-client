FROM alpine:latest
RUN mkdir service
COPY batch-manager-ui-client /service
COPY static /service/static
RUN chmod 777 /service/batch-manager-ui-client
ENV SERVICE_PORT=8080
ENV STATIC_PATH=/service/static/
EXPOSE 8080/tcp
WORKDIR /service
ENTRYPOINT [ "./batch-manager-ui-client" ]
