FROM alpine:latest

WORKDIR /app

COPY ./bin/app .

COPY .env .

RUN chmod +x ./app

CMD ["./app"]
