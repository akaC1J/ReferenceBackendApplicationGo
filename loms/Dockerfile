FROM alpine:latest

WORKDIR /app

COPY ./bin/loms .
COPY stock-data.json .
COPY .env .

RUN chmod +x ./loms

CMD ["./loms"]
