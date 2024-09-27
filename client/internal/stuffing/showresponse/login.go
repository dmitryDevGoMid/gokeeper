package showresponse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseLogin() string {

	var b strings.Builder

	login := m.data.RequestHTTP["login"]

	registerStatus := textGreenStyle.Render("Авторизациия пройдена успешно!")

	tpl := "Результат авторизации: \n\n"
	tpl += "%s\n\n"

	if login.Response.StatusCode != 200 && login.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при авторизации!")
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	if login.Response.StatusCode == 200 {
		tpl += "Вернитесь в меню и работайте!"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}
	if login.Response.StatusCode == 401 {
		registerStatus = textRedStyle.Render("Авторизация не пройдена!")
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	/*tpl += "\n\n"
	tpl += "Тело ответа: %s"

	tpl += "\n\n"
	tpl += "Token: %s"*/

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	//fmt.Println(m.data.User)

	fmt.Fprintf(&b, tpl, registerStatus)

	return style.SetStyleBeforeShowMenu(b.String())
	//return fmt.Sprintf(tpl, registerStatus, login.Response, m.data.User.Token)
}
