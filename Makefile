build:
go build -o server ./cmd/server
go build -o bot ./cmd/tgbot

docker:
sudo docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats

server:
./server

bot:
./bot

run:
docker
server
bot