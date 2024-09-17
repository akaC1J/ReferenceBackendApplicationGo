FROM alpine:latest

WORKDIR /app

COPY ./bin/cart .

COPY .env .

RUN chmod +x ./cart

CMD ["./cart"]
