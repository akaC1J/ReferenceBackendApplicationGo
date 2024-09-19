FROM alpine:3.20.3

WORKDIR /app

COPY ./bin/cart .

COPY .env .

RUN chmod +x ./cart

CMD ["./cart"]
