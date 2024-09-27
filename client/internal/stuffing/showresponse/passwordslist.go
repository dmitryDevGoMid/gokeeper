package showresponse

import (
	"fmt"
	"net/http"
)

func (m ShowResponse) ResponsePasswordsList() string {

	fmt.Println(m.data.RequestHTTP)

	passwordslist := m.data.RequestHTTP["passwordslist"]

	registerStatus := textGreenStyle.Render("Успешно получили список паролей!")

	tpl := "Результат получения списка паролей: \n\n"
	tpl += "%s\n\n"

	statusCode := passwordslist.Response.StatusCode

	if statusCode != 200 && statusCode != 401 {

		registerStatus = textRedStyle.Render("Возникли сложности при получении списка паролей!")
		tpl += fmt.Sprintf("Код ответа: %d %s", statusCode, http.StatusText(statusCode))
	}

	if statusCode == 200 {
		tpl += fmt.Sprintf(passwordslist.Response.Body)
		tpl += "Вернитесь в меню и продолжите работу. Меню: Password)"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", statusCode, http.StatusText(statusCode))
	}

	tpl += "\n\n"
	tpl += "Тело ответа: %s"
	tpl += "\n\n"
	tpl += "Error: %s"

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	return fmt.Sprintf(tpl, registerStatus, passwordslist.Response.Body, passwordslist.Response.Error)
}
