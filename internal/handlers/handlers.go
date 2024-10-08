package handlers

import (
	"bytes"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"os"
	"schedule-bot/internal/storage"
	"schedule-bot/internal/storage/sqlite"
	"schedule-bot/pkg/downloader"
	appmodels "schedule-bot/pkg/models"
	"schedule-bot/pkg/parse/href"
)

type Handler struct {
	Store storage.Storage
}

func (h *Handler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "/schedule - Проверяет расписание\n" +
			"/help - Показывает все команды",
	})
}

func (h *Handler) Schedule(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Проверяю появилось ли новое расписание...",
	})

	url, err := href.Get()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	filename, err := downloader.DownloadFile(url)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	var message *appmodels.Message

	message, err = h.Store.Fetch(update.Message.Chat.ID)
	if err != nil && err.Error() == sqlite.MessageNotFound {
		msg := &appmodels.Message{
			ChatID:    update.Message.Chat.ID,
			MessageID: 0,
			Filename:  filename,
		}

		if err = h.Store.Add(msg); err != nil {
			log.Printf("handlers.Schedule: %s\n", err.Error())
			return
		}
	}

	if message != nil {
		if message.Filename == filename {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Расписание не обновилось.",
			})

			b.ForwardMessage(ctx, &bot.ForwardMessageParams{
				ChatID:     update.Message.Chat.ID,
				FromChatID: update.Message.Chat.ID,
				MessageID:  message.MessageID,
			})

			err = os.Remove(filename)
			if err != nil {
				log.Printf("Failed to remove file: %v", err)
				return
			}
			return
		}
	} else {
		message, err = h.Store.Fetch(update.Message.Chat.ID)
		if err != nil {
			log.Printf("handlers.Schedule: %v\n", err.Error())
			return
		}
	}

	b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: message.MessageID,
	})

	f, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("handlers.Schedule: %v\n", err.Error())
		return
	}

	msg, _ := b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: update.Message.Chat.ID,
		Document: &models.InputFileUpload{
			Filename: filename,
			Data:     bytes.NewReader(f),
		},
		Caption: "Новое расписание",
	})

	if err := h.Store.Update(&appmodels.Message{
		ChatID:    message.ChatID,
		MessageID: msg.ID,
		Filename:  filename,
	}); err != nil {
		log.Printf("handlers.Schedule: %s\n", err.Error())
		return
	}

	b.PinChatMessage(ctx, &bot.PinChatMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: msg.ID,
	})

	err = os.Remove(filename)
	if err != nil {
		log.Printf("Failed to remove file: %v", err)
		return
	}
}
