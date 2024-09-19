FROM alpine:3.20.3

WORKDIR /app

COPY ./bin/loms .
COPY stock-data.json .
COPY .env .

RUN chmod +x ./loms

CMD ["./loms"]
