package showerrors

import (
	"fmt"
	"net/http"
)

func (m ShowErrors) ShowError() string {

	login := m.data.RequestHTTP["login"]

	registerStatus := textGreenStyle.Render("Авторизаци пройдена успешно!")

	tpl := "Результат авторизации: \n\n"
	tpl += "%s\n\n"

	if login.Response.StatusCode != 200 && login.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при авторизации!")
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	if login.Response.StatusCode == 200 {
		tpl += fmt.Sprintf("Вернитесь в меню и работайте!")
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}
	if login.Response.StatusCode == 401 {
		registerStatus = textRedStyle.Render("Авторизация не пройдена!")
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	tpl += "\n\n"
	tpl += "Тело ответа: %s"

	tpl += "\n\n"
	tpl += "Token: %s"

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Println(m.data.User)

	return fmt.Sprintf(tpl, registerStatus, login.Response, m.data.User.Token)
}
