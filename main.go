package main

import (
	"context"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Initialize the Telegram bot API client
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal("Unable to connect to Telegram API: ", err)
	}

	// Initialize the OpenAI API client
	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))

	whiteList := strings.Split(os.Getenv("WHITE_LISTED_USERS"), ",")
	if len(whiteList) == 0 {
		log.Fatal("No users in the white list!")
	}

	// Set up a message handler to respond to incoming messages
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !stringInSlice(update.Message.From.UserName, whiteList) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not my bro! I don't talk to strangers!")
			bot.Send(msg)
			continue
		}

		// Pass the message to the Hugging Face API
		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: update.Message.Text,
					},
				},
			},
		)

		if err != nil {
			log.Println(err)
			continue
		}

		// Send the response back to the user
		respText := response.Choices[0].Message.Content
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, respText)
		bot.Send(msg)
	}
}

// stringInSlice checks if a string is in a slice of strings
func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
