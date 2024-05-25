package module

import (
	"os"
	"time"

	"github.com/avast/retry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var bot *tgbotapi.BotAPI

func SendMessageTelegram(chatID int64, text string) error {
	if bot == nil {
		botTmp, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_KEY"))
		if err != nil {
			logrus.WithError(err).Error("SendMessageTelegram NewBotAPI %v", err)
			return err
		}
		bot = botTmp
	}
	if chatID == 0 {
		// group chat
		chatID = -1002141726230
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msgRes, err := bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("SendMessageTelegram NewBotAPI %v", err)
		return err
	}
	logrus.Debugf("Sent message done %v", msgRes)
	return nil
}

func SendMessageTelegramRetry(chatID int64, text string, attempts uint) error {
	err := retry.Do(
		func() error {
			err := SendMessageTelegram(chatID, text)
			if err != nil {
				logrus.WithError(err).Error("SendMessageTelegramRetry SendMessageTelegram %v", err)
				return err
			}
			return nil
		},
		retry.Attempts(attempts),
		retry.Delay(time.Minute*5),
	)
	if err != nil {
		logrus.Fatalf("SendMessageTelegramRetry err %v", err)
		return err
	}
	return nil
}
