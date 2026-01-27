package pdf

import (
	"fmt"
	"rsc.io/pdf"
)

func GetTextByCode(srcPath string) (string, error) {
	r, err := pdf.Open(srcPath)
	if err != nil {
		return "", err
	}

	var text string

	for i := range r.NumPage() {
		fmt.Println(r.NumPage(), i)
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		content := p.Content()
		for _, t := range content.Text {
			text += t.S
		}
	}

	return text, nil
}
