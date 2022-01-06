package telegram

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot/bots"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type (
	telegramBot struct {
		user *tb.User
		bot  *tb.Bot
	}
)

func New(cfg *ConfigParameters) bots.IBot {
	if cfg.ID == 0 {
		panic("must set telegram user ID")
	}

	bot, err := tb.NewBot(tb.Settings{Token: cfg.Token, Poller: &tb.LongPoller{Timeout: 5 * time.Second}})
	if err != nil {
		panic(err)
	}

	return &telegramBot{
		user: &tb.User{ID: cfg.ID},
		bot:  bot,
	}
}

func (t *telegramBot) handleHealthChecker(_ *tb.Message) {
	log.Println("checking status")
	if err := t.retrySend("I am alive!", t.user, t.bot.Send); err != nil {
		log.Println(err)
	}
}

func (t *telegramBot) SendNotification(msg interface{}) error {
	if err := t.retrySend(msg, t.user, t.bot.Send); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (t *telegramBot) StartServer() error {
	t.bot.Handle("/amialive", t.handleHealthChecker)

	log.Println("starting telegram server")
	t.bot.Start()

	return nil
}

func (t *telegramBot) Close() {
	for i := 0; i < 3; i++ {
		_, err := t.bot.Close()
		if err != nil {
			switch err.(type) {
			case tb.FloodError:
				err := err.(tb.FloodError)
				if err.Code == http.StatusTooManyRequests {
					log.Println("retrying closing telegram bot...")
					time.Sleep(time.Duration(err.RetryAfter) * time.Second)
				}
			default:
				time.Sleep(time.Duration(30) * time.Second)
			}
		} else {
			log.Println("shutting down telegram server")
			break
		}
	}
}

func (t *telegramBot) retrySend(msg interface{}, to tb.Recipient, fnSend func(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error), options ...interface{}) error {
	var err error
	for i := 0; i < 3; i++ {
		if _, err = fnSend(to, msg, options...); err == nil {
			return nil
		}
		switch err.(type) {
		case tb.FloodError:
			err := err.(tb.FloodError)
			if err.Code == http.StatusTooManyRequests {
				log.Println("retrying sending message, time", i)
				time.Sleep(time.Duration(err.RetryAfter) * time.Second)
			}
		default:
			t := rand.Intn(10)
			time.Sleep(time.Duration(t) * time.Second)
		}
	}
	return fmt.Errorf("error retrying sending message: %w", err)
}
