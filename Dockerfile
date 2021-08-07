FROM golang:alpine

RUN mkdir -p /discordBot

WORKDIR /discordBot

COPY . .

RUN go build -o discordBot

ENTRYPOINT ["./discordBot"]