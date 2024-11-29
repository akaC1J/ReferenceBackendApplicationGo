FROM alpine:3.20.3

WORKDIR /app

COPY ./bin/notifier .

COPY .env .

RUN chmod +x ./notifier

CMD ["./notifier"]
