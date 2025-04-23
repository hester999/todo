package telegram

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const (
	commandStart   = "start"
	commandHelp    = "help"
	commandCreate  = "create"
	commandDelete  = "delete"
	commandUpdate  = "update"
	commandGetTask = "getTask"
)

type tmpTask struct {
	title string
	desc  string
	id    string
}

type userState struct {
	command string
	state   string
}

func (b *Bot) commandHandler(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, "Unknown command")
	switch message.Command() {
	case commandStart:
		msg.Text = "команда старт"
		_, err := b.bot.Send(msg)
		return err

	case commandHelp:
		msg.Text = "команда помощи"
		_, err := b.bot.Send(msg)
		return err

	case commandCreate:
		b.createHandler(message)
		return nil

	case commandDelete:
		msg.Text = "удааляет таск"
		_, err := b.bot.Send(msg)
		return err

	case commandUpdate:
		msg.Text = "обновляет таск"
		_, err := b.bot.Send(msg)
		return err
	case commandGetTask:
		b.getHandler(message)

		return nil
	default:
		_, err := b.bot.Send(msg)
		return err
	}

}

func (b *Bot) handleMessage(message *tg.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tg.NewMessage(message.Chat.ID, "Я тебя не понимаю введи команду из списка")
	b.bot.Send(msg)
}

func (b *Bot) createHandler(message *tg.Message) {
	chatID := message.Chat.ID
	userID := message.From.ID
	text := message.Text

	// этоо нписал grok, на сколько я понял, тут просто в первый раз устанавливается состоояние
	if message.IsCommand() {
		b.states[userID] = userState{commandCreate, "waiting_title"}

		msg := tg.NewMessage(chatID, "Введи название задачи:")
		b.bot.Send(msg)
		return
	}

	switch b.states[userID].state {
	case "waiting_title":

		b.taskTmp[userID] = tmpTask{title: text}
		b.states[userID] = userState{commandCreate, "waiting_description"}
		msg := tg.NewMessage(chatID, "Введи описание задачи:")
		b.bot.Send(msg)

	case "waiting_description":
		task := b.taskTmp[userID]
		task.desc = text
		b.states[userID] = userState{}
		delete(b.taskTmp, userID)
		delete(b.states, userID)

		newTask, err := b.taskService.CreateTask(task.title, task.desc)
		if err != nil {
			msg := tg.NewMessage(chatID, "Ошибка создания задачи.")
			b.bot.Send(msg)
			return
		}

		msg := tg.NewMessage(chatID, "Задача создана: "+newTask.Title+" - "+newTask.Description)
		b.bot.Send(msg)
	}

}

func (b *Bot) getHandler(message *tg.Message) {
	chatID := message.Chat.ID
	userID := message.From.ID

	if message.IsCommand() {
		b.states[userID] = userState{commandGetTask, "choose_task"}
		tasks, err := b.taskService.GetAllTasks()

		if err != nil {
			msg := tg.NewMessage(chatID, "Ошибка получения задач")
			b.bot.Send(msg)
			return
		}
		if len(tasks) == 0 {
			msg := tg.NewMessage(chatID, "У тебя нет задач")
			b.bot.Send(msg)
			return
		}

		var taskList strings.Builder

		keyboard := tg.NewReplyKeyboard()
		var row []tg.KeyboardButton
		for i, task := range tasks {
			taskList.WriteString(fmt.Sprintf("%d. %s\n", i+1, task.Title))
			button := tg.NewKeyboardButton(task.Title)
			b.taskTmpID[chatID] = append(b.taskTmpID[chatID], tmpTask{title: task.Title, id: task.Id})
			ms := tg.NewMessage(chatID, b.taskTmp[chatID].title)
			b.bot.Send(ms)
			row = append(row, button)
			if len(row) == 2 || i == len(tasks)-1 {
				keyboard.Keyboard = append(keyboard.Keyboard, tg.NewKeyboardButtonRow(row...))
				row = nil
			}
		}
		keyboard.ResizeKeyboard = true
		msg := tg.NewMessage(chatID, taskList.String()) // Вывод списка задач
		msg.ReplyMarkup = keyboard
		b.bot.Send(msg)
	} else {
		b.callBackHandler(message.Chat.ID, message.Text)
	}

}

func (b *Bot) callBackHandler(chatID int64, text string) {
	callBackFromButton := text
	ms := tg.NewMessage(chatID, callBackFromButton)
	b.bot.Send(ms)
	var taskID string
	if task, ok := b.taskTmpID[chatID]; ok {
		for _, v := range task {
			if v.title == callBackFromButton {
				//ms = tg.NewMessage(chatID, v.id)
				taskID = v.id
				break
			}
		}

	}

	task, err := b.taskService.GetTaskById(taskID)
	if err != nil {
		msg := tg.NewMessage(chatID, "Ошибка получения задачи.")
		b.bot.Send(msg)
	}
	msg := tg.NewMessage(chatID, fmt.Sprintf("Задача: %s\nОписание: %s\nСтатус: %s", task.Title, task.Description, task.Status))
	msg.ReplyMarkup = tg.NewRemoveKeyboard(true) // Убираем клавиатуру
	b.bot.Send(msg)

	delete(b.states, chatID)
	delete(b.taskTmpID, chatID)

}
