FROM golang:1.22-alpine AS builder

WORKDIR /usr/src

COPY . .

EXPOSE 8000

CMD cd . && go run ./app/cmd/app/main.go

RUN apk --no-cache add bash gettext