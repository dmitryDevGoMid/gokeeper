package showresponse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseRegister() string {

	var b strings.Builder

	register := m.data.RequestHTTP["register"]

	registerStatus := textGreenStyle.Render("Регистрация пройдена успешно!")

	tpl := "Результат регистрации: \n\n"
	tpl += "%s\n\n"

	if register.Response.StatusCode != 201 && register.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при регистрация!")
		tpl += fmt.Sprintf("Код ответа: %d %s", register.Response.StatusCode, http.StatusText(register.Response.StatusCode))
	}

	if register.Response.StatusCode == 201 {
		tpl += fmt.Sprintf("Вернитесь в меню и войдите в систему. Меню: login)")
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", register.Response.StatusCode, http.StatusText(register.Response.StatusCode))
	}
	//tpl += "\n\n"
	//tpl += "Тело ответа: %s"

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Fprintf(&b, tpl, registerStatus)

	return style.SetStyleBeforeShowMenu(b.String())

	//return fmt.Sprintf(tpl, registerStatus, register.Response)
}
