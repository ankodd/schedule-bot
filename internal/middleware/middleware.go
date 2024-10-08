package middleware

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"schedule-bot/pkg/utils/commands"
)

func Logging(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		log.Printf("@%v send message: %v", update.Message.From.Username, update.Message.Text)
		next(ctx, b, update)
	}
}

func SkipIsBot(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update.Message.From.IsBot == true {
			return
		}

		next(ctx, bot, update)
	}
}

func Command(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if !commands.IsCommand(update.Message.Text) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: fmt.Sprintf(
					"Привет @%v, %v это не команда.\n "+
						"я работаю только с командами.\n /help для отображения всех команд",
					update.Message.From.Username, update.Message.Text,
				),
			})
			return
		}

		next(ctx, b, update)
	}
}
