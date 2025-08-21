FROM alpine

WORKDIR /app

COPY chat_go /app/

COPY static /app/static

ENTRYPOINT ["./chat_go"]
