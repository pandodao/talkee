# Talkee

An Open-source Web3 Commenting and Chat tool with Ethereum Login & Wallets

## Demo

- Comment Demo: https://developers.pando.im/demo/comment-demo.html
- Chat Demo: https://developers.pando.im/demo/chat-demo.html

## Create your own

Please visit [Pando Developers Console](https://developers.pando.im/console/talkee) to create your own Talkee sites.

## Integration

Please visit [this page](https://developers.pando.im/guide/talkee.html) to integrate Talkee into your site.

## Features

- [x] 💻 One-Click Installation
- [x] 🤑 Ethereum Login by Metamask or WalletConnect Wallets
- [x] 💬 Instant Chat APIs (UI in the process)
- [x] 👍 Reply, Like, Share. 
- [x] 🌐 Content on the Blockchain (Arweave right now, IPFS and others in the process)
- [x] 💰 Token AirDrop On-Demand
- [ ] 🔔 Notifications
- [ ] 🔑 Comment/ChatRoom Requirements: crypto requirement for people access your comment threads or chat room
- [ ] 🚫 Anti-Spam Integration & Moderation tools
- [ ] 🧑🏻‍💼 User Profile and Reputation
- [ ] 🤝 Transfer crypto to other users
- [ ] 📥 Import Comments from Disqus, Commento, etc
- [ ] 📤 Export Comments

## Installation

```bash
git clone https://github.com/pandodao/talkee.git
cd talkee
go build
```

## Preqrequisites

You need to have a running postgresql database, a keystore file from [Mixin Developers](https://developers.mixin.one/dashboard).

To enable "Content on the Blockchain" feature, you also need an arweave wallet file from [Arweave](https://docs.arweave.org/info/wallets/arweave-wallet) and put it under `keystores/wallet.json` of working directory.


## Configuration

Create a config file `config.yaml` in the working directory.

```yaml
# database config
db:
  driver: "postgres"
  datasource: "user=foobar dbname=talkee host=localhost password=foobar sslmode=disable"

# auth config
auth:
  # a random string to generate jwt token
  jwt_secret: "112233"
  # please get it from https://developers.mixin.one/dashboard
  mixin_client_secret: ".."

# optional, not implemented yet
aws:
  key: ""
  secret: ""
  region: ""
  bucket: ""

# optional, not implemented yet
sys:
  attachment_base: "http://.."
```

## Run the services manually

run migrate database

```bash
./talkee migrate 
```

run `./talkee help` to see full commands

run workers

```bash
./talkee -f YOUR_KEYSTORE_FILE worker
```

run websocket server

```bash
./talkee -f YOUR_KEYSTORE_FILE wss [port] 
```

run httpd server

```bash
./talkee -f YOUR_KEYSTORE_FILE httpd [port] 
```

## Run the services in docker

build image
```bash
docker build -t talkee:latest .
```

create docker-compose.yml

```bash
version: "3.4"
x-volumes: &default-volumes
  - "./keystores:/app/keystores"
  - "./config.yaml:/app/config.yaml"
services:
  api:
    image: talkee:latest
    entrypoint: ["/app/talkee", "--file","YOUR_KEYSTORE_FILE", "httpd", "80"]
    ports:
      - "8080:80"
    volumes: *default-volumes

  wss:
    image: talkee:latest
    entrypoint: ["/app/talkee", "--file","YOUR_KEYSTORE_FILE", "wss", "80"]
    ports:
      - "8081:80"

  worker:
    image: talkee:latest    
    entrypoint: ["/app/talkee", "--file","YOUR_KEYSTORE_FILE", "worker", "80"]
    ports:
      - "8090:80"
```

run via docker-compose
```bash
docker-compose up -d 
```

run database migration 
```bash
docker run  --rm -ti -v [YOUR_CONFIG_FILE]:/app/config.yaml talkee:latest /app/talkee migrate
```