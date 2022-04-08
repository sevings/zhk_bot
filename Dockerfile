FROM golang:1.18.0-alpine3.15 AS zhk_bot_build

WORKDIR $GOPATH/src/github.com/sevings/zhk_bot

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go .
COPY internal/ internal/
RUN go build -o /opt/zhk_bot/zhk_bot .

FROM alpine:3.15 AS zhk_bot

WORKDIR /opt/zhk_bot
RUN chmod 0777 .

RUN addgroup -S zhk_bot
RUN adduser -S zhk_bot -G zhk_bot
USER zhk_bot

COPY --from=zhk_bot_build /opt/zhk_bot/zhk_bot .

CMD ["./zhk_bot"]
