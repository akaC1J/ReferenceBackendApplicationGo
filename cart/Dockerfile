FROM alpine:latest

WORKDIR /app

COPY ./bin/app/app .

COPY .env .

RUN chmod +x ./app

CMD ["./app"]
