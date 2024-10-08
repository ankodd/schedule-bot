# Schedule-bot bot for check schedule in UrPret!

## Build from source

1. ### Clone repository
  ```bash
  git clone https://github.com/ankodd/schedule-bot.git && cd schedule-bot
  ```
2. ### Install dependencies
  ```bash
  go mod tidy
  go mod download
  ```
3. ### Run
  **Before launching the app you must do env variable BOT_KEY, which equal to your telegram bot token**
  ```bash
  go run ./cmd/schedule-bot/main.go
  ```

## Running docker container

1. ### Clone repository
  ```bash
  git clone https://github.com/ankodd/schedule-bot.git && cd schedule-bot
  ```
2. ### Build docker image
  ```bash
  docker build -t schedule-bot .
  ```
3. ### Run container
  ```bash
  docker run -e BOT_KEY=YOUR_BOT_TOKEN -p 8080:8080 schedule-bot
  ```

**Thanks for read! and visiting my repo :D**
