FROM golang:1.15-buster AS build
WORKDIR /app
ADD . .

RUN go build -o bin/chat  ./cmd/chat

ENV DB_USER=user
ENV DB_PASS=pass
ENV DB_NAME=chat-db
ENV DB_ADDR=db:5432

CMD ["/app/bin/chat"]
