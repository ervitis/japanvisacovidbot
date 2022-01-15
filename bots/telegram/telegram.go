package telegram

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/email"
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

var (
	menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
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

func (t *telegramBot) handleSendEmail(msg *tb.Message) {
	links := []string{
		"https://www.mofa.go.jp/j_info/visit/visa/pdfs/application1_e.pdf",
		"https://www.mofa.go.jp/mofaj/files/000124525.pdf",
	}

	if err := t.retrySend(fmt.Sprintf(email.MessageConfirmation, links[0], links[1]), t.user, t.bot.Send, menu); err != nil {
		log.Println(err)
		return
	}

	if err := t.bot.Delete(msg); err != nil {
		log.Println(err)
	}
}

func (t *telegramBot) handleSendEmailToEmbassy(fn *tb.Callback) {
	emailSvc := email.New(&email.Config)

	go func(emailSvc email.IEmail) {
		if err := emailSvc.Send(); err != nil {
			log.Println(err)
			return
		}
	}(emailSvc)

	if err := t.bot.Respond(fn, &tb.CallbackResponse{Text: "Email sent to " + emailSvc.Properties().To + " if not, see logs"}); err != nil {
		log.Println(err)
		return
	}

	if err := t.bot.Delete(fn.Message); err != nil {
		log.Println("Could not delete message to send email", err)
	}
	log.Println("The email was sent with headers", emailSvc.Properties().Headers)
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
	t.bot.Handle("/email", t.handleSendEmail)

	btnSend := menu.Data("Send email", "sendEmail", "email")
	menu.Inline(menu.Row(btnSend))
	t.bot.Handle(&btnSend, t.handleSendEmailToEmbassy)

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
