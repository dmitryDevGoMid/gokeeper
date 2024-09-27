package showresponse

import (
	"fmt"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"
)

func (m ShowResponse) DialogSelect() string {
	var b strings.Builder

	buttonCopy := &blurredButtonCopy
	if m.focusIndex == 0 {
		buttonCopy = &focusedButtonCopy
	}

	buttonRewrite := &blurredButtonRewrite
	if m.focusIndex == 1 {
		buttonRewrite = &focusedButtonRewrite
	}

	buttonCancel := &blurredButtonCancel
	if m.focusIndex == 2 {
		buttonCancel = &focusedButtonCancel
	}

	fmt.Fprintf(&b, "%s\n\n", "Вы хотите перезаписать файл или создать копию?")
	fmt.Fprintf(&b, "%s    ", *buttonCopy)
	fmt.Fprintf(&b, "%s    ", *buttonRewrite)
	fmt.Fprintf(&b, "%s\n\n", *buttonCancel)
	tpl := "\n\n"
	tpl += dotStyle + subtleStyle.Render("tab: select button")

	return style.SetStyleBeforeShowMenu(b.String())
}
