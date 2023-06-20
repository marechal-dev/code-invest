FROM golang:alpine AS build

WORKDIR /usr/app

COPY . .

RUN go mod tidy

RUN make build

CMD ["make", "start"]
