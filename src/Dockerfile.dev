FROM golang:1.27-alpine

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY src/ /app/src/
RUN cd src && go mod download

CMD sh -c "cd /app/src && air"