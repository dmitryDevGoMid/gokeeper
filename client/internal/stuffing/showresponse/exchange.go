package showresponse

import (
	"fmt"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseExchangeSet() string {

	//exchange := m.data.RequestHTTP["exchange"]
	registerStatus := textGreenStyle.Render("Выполнена успешная передача ключа серверу!")

	tpl := dotStyle + subtleStyle.Render("ctrl+c, esc: quit")
	return fmt.Sprintf("%v\n\n%s", registerStatus, tpl)
}

func (m ShowResponse) ResponseExchangeGet() string {
	var b strings.Builder

	exchange := m.data.RequestHTTP["exchangeget"]

	var registerStatus string

	if exchange.Response.StatusCode == 0 && exchange.Response.Error != "" {
		registerStatus = textRedStyle.Render("Возникли сложности при обращении на сервер!\n")
		registerStatus += textRedStyle.Render("Ошибка:" + exchange.Response.Error)
	} else if exchange.Response.StatusCode != 201 {

		registerStatus = textGreenStyle.Render("Запрос выволнен, однако он не успешный!")
		registerStatus += textGreenStyle.Render(fmt.Sprintf("Код ответа: %d, %s", exchange.Response.StatusCode, exchange.Response.Error))
	} else {
		registerStatus = textGreenStyle.Render(fmt.Sprintf("Выполнен успешный запрос на ключ от севера!\n\nКод ответа: %d\n", exchange.Response.StatusCode))
	}

	tpl := dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Fprintf(&b, "%v\n\n%s", registerStatus, tpl)

	return style.SetStyleBeforeShowMenu(b.String())

	//return fmt.Sprintf("%v\n\n%s", registerStatus, tpl)
}
