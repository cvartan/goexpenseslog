package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cvartan/goconfig"
	"github.com/cvartan/goexpenseslog/internal/model"
)

const startText string = `
<b>Бот учета расходов.</b>

Отправляйте расходы в формате: дата сумма комментарий

Дата должна быть указана в формате: <i>ДД.ММ.ГГГГ</i>
(ДД - день, ММ = номер месяца, ГГГГ - 4 цифры года).
Пример даты: 10.08.2024 (10 августа 2024 года)

Сумма должна быть указана только в рублях, без копеек.
`

type TelegramBotService struct {
	rawMessagesRepo model.RawMessageRepositoryHandler
	expensesRepo    model.ExpensesRepositoryHandler
	configuration   *goconfig.Configuration
}

func NewBot(rawMessageRepo model.RawMessageRepositoryHandler, expensesRepo model.ExpensesRepositoryHandler, configuration *goconfig.Configuration) *TelegramBotService {
	return &TelegramBotService{
		rawMessagesRepo: rawMessageRepo,
		expensesRepo:    expensesRepo,
		configuration:   configuration,
	}
}

func (b *TelegramBotService) HandleStart(req *BotRequest, resp *BotResponse) error {
	resp.Text = startText

	return nil
}

func (b *TelegramBotService) HandleDefault(req *BotRequest, resp *BotResponse) error {
	if req.Message.IsCommand() {
		resp.Text = "Неизвестная команда. Используйте меню."
		return fmt.Errorf("unrecognized command: %s", req.Message.Command())
	}

	log.Printf("Сообщение от %d со значением: %s", req.Message.From.ID, req.Message.Text)

	rm := &model.RawMessage{
		UserId: req.Message.From.ID,
		Value:  req.Message.Text,
	}

	if err := b.rawMessagesRepo.Add(rm); err != nil {
		resp.Text = "Ошибка. Проверьте формат и повторите снова."
		return err
	}

	fields := strings.Split(req.Message.Text, " ")
	ei := &model.ExpenseInfo{
		Id: rm.Id,
	}
	c := 0
	for _, f := range fields {
		if f != "" {
			switch c {
			case 0:
				{
					dt, err := time.Parse("02.01.2006", f)
					if err != nil {
						resp.Text = "Ошибка. Неправильная дата."
						return err
					}
					ei.Date = dt
					c++
				}
			case 1:
				{
					val, err := strconv.ParseInt(f, 10, 32)
					if err != nil {
						resp.Text = "Ошибка. Неправильное значение суммы (должно быть без копеек)"
						return err
					}

					ei.Value = int32(val)
					c++
				}
			case 2:
				{
					if ei.Description != "" {
						ei.Description = ei.Description + " "
					}
					ei.Description = ei.Description + f
				}
			}
		}
	}

	if err := b.expensesRepo.Add(ei); err != nil {
		resp.Text = "Ошибка. Повторите снова."
		return err
	}

	resp.Text = fmt.Sprintf("OK %s", req.Message.From.FirstName)

	return nil
}

func (b *TelegramBotService) HandleList(req *BotRequest, resp *BotResponse) error {
	controlUserId := b.configuration.Get("control.userId").Int()
	if req.Message.From.ID != controlUserId {
		resp.Text = "У Вас нет прав на эту операцию!"
		return fmt.Errorf("unauthorized access for list command for user %d", req.Message.From.ID)
	}

	msgs, err := b.rawMessagesRepo.GetAll()
	if err != nil {
		resp.Text = err.Error()
		return err
	}

	buf := strings.Builder{}

	for _, val := range *msgs {
		buf.WriteString(
			fmt.Sprintf(
				"- %s\n<i> %d : %d : %s</i>\n",
				val.Value,
				val.Id,
				val.UserId,
				val.Created.Format(time.StampMilli),
			),
		)
	}

	resp.Text = buf.String()

	return nil
}

func (b *TelegramBotService) HandleUserData(req *BotRequest, resp *BotResponse) error {
	userId := req.Message.From.ID

	msgs, err := b.rawMessagesRepo.GetUserData(userId)
	if err != nil {
		resp.Text = err.Error()
		return err
	}

	buf := strings.Builder{}

	for _, val := range *msgs {
		buf.WriteString(
			fmt.Sprintf(
				"- %s\n",
				val.Value,
			),
		)
	}

	resp.Text = buf.String()

	return nil
}

func (b *TelegramBotService) HandleMonthSummary(req *BotRequest, resp *BotResponse) error {
	sum, err := b.expensesRepo.GetMonthSummary()
	if err != nil {
		resp.Text = "Ошибка. Повторите запрос позже"
		return err
	}

	resp.Text = fmt.Sprintf("Расходы за этот месяц: %d", sum)
	return nil
}

func (b *TelegramBotService) HandlePrevMonthSummary(req *BotRequest, resp *BotResponse) error {
	sum, err := b.expensesRepo.GetPrevMonthSummary()
	if err != nil {
		resp.Text = "Ошибка. Повторите запрос позже"
		return err
	}

	resp.Text = fmt.Sprintf("Расходы за предыдущий месяц: %d", sum)
	return nil
}
