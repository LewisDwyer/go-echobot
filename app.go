package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot   *tgbotapi.BotAPI
	TOKEN = os.Getenv("BOT_TOKEN")
	URL   = os.Getenv("URL")
)

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
}

func respond(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the update from JSON
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
		return
	}

	chatID := update.Message.Chat.ID
	messageID := update.Message.MessageID
	text := update.Message.Text

	// Handle /start command
	if text == "/start" {
		msg := tgbotapi.NewMessage(chatID, "Hi! I respond by echoing messages. Give it a try!")
		msg.ReplyToMessageID = messageID
		bot.Send(msg)
	} else {
		// Echo the message back
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyToMessageID = messageID
		bot.Send(msg)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func setWebhook(w http.ResponseWriter, r *http.Request) {
	webhookURL := fmt.Sprintf("%s%s", URL, TOKEN)

	// Parse URL
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		log.Printf("URL parse error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "webhook setup failed")
		return
	}

	// Create webhook config
	config := tgbotapi.WebhookConfig{URL: parsedURL}

	_, err = bot.Request(config)
	if err != nil {
		log.Printf("Webhook error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "webhook setup failed")
		return
	}
	fmt.Fprint(w, "webhook setup ok")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, welcome to the telegram bot index page")
}

func main() {
	// Register handlers
	http.HandleFunc(fmt.Sprintf("/%s", TOKEN), respond)
	http.HandleFunc("/setwebhook", setWebhook)
	http.HandleFunc("/", index)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
