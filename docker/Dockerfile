FROM debian:bookworm

WORKDIR /app

COPY navReceiveApp .
COPY config.yaml .

RUN chmod +x navReceiveApp

EXPOSE 3088

ENTRYPOINT ["/app/navReceiveApp"]
