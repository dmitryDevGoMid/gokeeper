package showresponse

import (
	"fmt"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) DialogPasswordDeleteOrNot() string {
	var b strings.Builder

	buttonYes := &blurredButtonDeleteFileYes
	if m.focusIndex == 1 {
		buttonYes = &focusedButtonDeleteFileYes
	}

	buttonNo := &blurredButtonDeleteFileNo
	if m.focusIndex == 0 {
		buttonNo = &focusedButtonDeleteFileNo
	}

	fmt.Fprintf(&b, "%s\n\n", "Вы хотите удалить пароль?")
	fmt.Fprintf(&b, "%s    ", *buttonYes)
	fmt.Fprintf(&b, "%s    ", *buttonNo)

	tpl := "\n\n"
	tpl += dotStyle + subtleStyle.Render("ctrl+c, esc: quit")
	fmt.Fprintf(&b, "%s", tpl)

	return style.SetStyleBeforeShowMenu(b.String())
}
