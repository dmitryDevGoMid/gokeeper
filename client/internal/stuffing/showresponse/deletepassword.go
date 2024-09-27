package showresponse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseDeletePassword() string {

	var b strings.Builder

	card := m.data.RequestHTTP["passwordslist"]

	chocedKey, err := strconv.Atoi(card.Choced)
	if err != nil {
		fmt.Println("Error creating:", err)
	}
	passwordList := *card.PasswordsList
	selectedPassword := passwordList[chocedKey]

	dpassword := m.data.RequestHTTP["deletepassword"]

	registerStatus := textGreenStyle.Render("Пароль успешно удален!")

	tpl := fmt.Sprintf("Результат удаления пароля: %s\n\n", selectedPassword.Description)
	tpl += "%s\n\n"

	if dpassword.Response.StatusCode != 200 && dpassword.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при удалении пароля")
		tpl += fmt.Sprintf("Код ответа: %d %s", dpassword.Response.StatusCode, http.StatusText(dpassword.Response.StatusCode))
	}

	if dpassword.Response.StatusCode == 200 {
		tpl += "Вернитесь в меню и работайте!"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", dpassword.Response.StatusCode, http.StatusText(dpassword.Response.StatusCode))
	}

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Fprintf(&b, tpl, registerStatus)

	return style.SetStyleBeforeShowMenu(b.String())
}
