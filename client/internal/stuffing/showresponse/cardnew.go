package showresponse

import (
	"fmt"
	"net/http"
)

func (m ShowResponse) ResponseNewCard() string {

	register := m.data.RequestHTTP["newcard"]

	registerStatus := textGreenStyle.Render("Успешно сохранили новую карту!")

	tpl := "Результат сохранения новой карты: \n\n"
	tpl += "%s\n\n"

	if register.Response.StatusCode != 201 && register.Response.StatusCode != 401 && register.Response.StatusCode != 200 {
		registerStatus = textRedStyle.Render("Возникли сложности при сохранении новой карты!")
		tpl += fmt.Sprintf("Код ответа: %d %s", register.Response.StatusCode, http.StatusText(register.Response.StatusCode))
	}

	if register.Response.StatusCode == 201 {
		tpl += "Вернитесь в меню и продолжите работу. Меню: Card)"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", register.Response.StatusCode, http.StatusText(register.Response.StatusCode))
	}

	if register.Response.StatusCode == 200 {
		tpl += "Вернитесь в меню и продолжите работу. Меню: Card)"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", register.Response.StatusCode, http.StatusText(register.Response.StatusCode))
	}

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	return fmt.Sprintf(tpl, registerStatus)
}
