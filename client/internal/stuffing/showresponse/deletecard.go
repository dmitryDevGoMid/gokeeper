package showresponse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) ResponseDeleteCard() string {

	var b strings.Builder

	card := m.data.RequestHTTP["cardslist"]

	chocedKey, err := strconv.Atoi(card.Choced)
	if err != nil {
		fmt.Println("Error creating:", err)
	}
	cardsList := *card.CardsList
	selectedCard := cardsList[chocedKey]

	dcard := m.data.RequestHTTP["deletecard"]

	registerStatus := textGreenStyle.Render("Карта успешно удалена!")

	tpl := fmt.Sprintf("Результат удаления карты: %s\n\n", selectedCard.Number)
	tpl += "%s\n\n"

	if dcard.Response.StatusCode != 200 && dcard.Response.StatusCode != 401 {
		registerStatus = textRedStyle.Render("Возникли сложности при удалении карты")
		tpl += fmt.Sprintf("Код ответа: %d %s", dcard.Response.StatusCode, http.StatusText(dcard.Response.StatusCode))
	}

	if dcard.Response.StatusCode == 200 {
		tpl += "Вернитесь в меню и работайте!"
		tpl += "\n\n"
		tpl += fmt.Sprintf("Код ответа: %d %s", dcard.Response.StatusCode, http.StatusText(dcard.Response.StatusCode))
	}

	tpl += "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")

	fmt.Fprintf(&b, tpl, registerStatus)

	return style.SetStyleBeforeShowMenu(b.String())
}
