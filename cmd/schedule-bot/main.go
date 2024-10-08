package main

import (
	"context"
	"database/sql"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"schedule-bot/internal/handlers"
	"schedule-bot/internal/middleware"
	"schedule-bot/internal/storage"
	"schedule-bot/internal/storage/sqlite"
	"sync"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("sqlite3", os.Getenv("STORAGE_PATH"))
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed connect to database: %v", err)
	}

	store, err := sqlite.New(db)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler(store)),
		bot.WithMiddlewares(middleware.SkipIsBot, middleware.Command, middleware.Logging),
	}

	b, err := bot.New(os.Getenv("BOT_KEY"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Start(ctx)
	}()

	log.Printf("bot started\n")
	wg.Wait()
}

func handler(storage storage.Storage) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		h := handlers.Handler{
			Store: storage,
		}

		switch update.Message.Text {
		case "/help":
			h.Help(ctx, b, update)
		case "/schedule":
			h.Schedule(ctx, b, update)
		default:
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Неизвестная команда",
			})
		}
	}
}
