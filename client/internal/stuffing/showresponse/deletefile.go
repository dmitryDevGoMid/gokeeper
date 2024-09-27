package showresponse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseDeleteFile() string {

	var b strings.Builder

	file := m.data.RequestHTTP["fileslist"]

	chocedKey, err := strconv.Atoi(file.Choced)
	if err != nil {
		fmt.Println("Error creating:", err)
	}
	filesList := *file.FilesList
	selectedFile := filesList[chocedKey]

	login := m.data.RequestHTTP["deletefile"]

	registerStatus := textGreenStyle.Render("Файл успешно удален!")

	tpl := fmt.Sprintf("Результат удаления файла: %s\n\n", selectedFile.Filename)
	tpl += "%s\n\n"

	if login.Response.StatusCode != 200 && login.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при удалении файла")
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	if login.Response.StatusCode == 200 {
		tpl += "Вернитесь в меню и работайте!"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", login.Response.StatusCode, http.StatusText(login.Response.StatusCode))
	}

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Fprintf(&b, tpl, registerStatus)

	return style.SetStyleBeforeShowMenu(b.String())
}
