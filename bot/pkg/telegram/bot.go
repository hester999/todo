package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"todo/bot/pkg/adapter"
)

// хз правльно или нет так формировать структуру бота
type Bot struct {
	bot         *tg.BotAPI
	taskService adapter.BotAdapter
	states      map[int64]userState
	taskTmp     map[int64]tmpTask
	taskTmpID   map[int64][]tmpTask
}

func NewBot(bot *tg.BotAPI, service adapter.BotAdapter) *Bot {
	return &Bot{
		bot:         bot,
		taskService: service,
		states:      make(map[int64]userState), //для сохранения состояний в диааалоге int64- user_id
		taskTmp:     make(map[int64]tmpTask),   // для хранения структур и передачи даанных в usecases
		taskTmpID:   make(map[int64][]tmpTask), // для хранения структур и передачи даанных в usecases
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdateChannel()
	b.handelUpdate(updates)

	return nil
}

// помогал писать grok)
func (b *Bot) handelUpdate(updates tg.UpdatesChannel) {
	for update := range updates {
		// Обработка нажатия на инлайн-кнопку
		//if update.CallbackQuery != nil {
		//	b.callBackHandler(update)
		//	continue
		//}

		// Проверка, что есть сообщение
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID

		// Обработка команд
		if update.Message.IsCommand() {
			b.commandHandler(update.Message)
			continue
		}

		// Обработка сообщений в диалоге
		state := b.states[userID]
		switch state.command {
		case commandCreate:
			b.createHandler(update.Message)
		case commandGetTask:
			b.getHandler(update.Message)
		default:
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) initUpdateChannel() tg.UpdatesChannel {

	u := tg.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)

}
