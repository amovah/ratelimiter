FROM golang:1.12.6

WORKDIR /app

COPY . /app

RUN go build

CMD ["go", "run", "ratelimiter"]