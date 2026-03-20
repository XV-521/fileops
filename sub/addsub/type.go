package addsub

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type SubType int

const (
	SubUn SubType = iota
	SubS
	SubH
)

var UnsupportedCnvTypeErr = errors.New("unsupported sub type")

type SubFn func(videoPath, subPath, dstPath string) error

func getSubFn(st SubType) (SubFn, error) {
	switch st {
	case SubS:
		return addSoftSub, nil
	case SubH:
		return addHardSub, nil
	default:
		return nil, UnsupportedCnvTypeErr
	}
}

type SimpleStyle string

func (s SimpleStyle) check() bool {
	re := regexp.MustCompile(
		`^Fontname=[^,]+,Fontsize=[1-9][0-9]*,PrimaryColour=&H[0-9a-fA-F]+,OutlineColour=&H[0-9a-fA-F]+,BackColour=&H[0-9a-fA-F]+,BorderStyle=[13],Outline=[0-9]+,Shadow=[0-9]+,Alignment=[1-9],MarginV=[1-9][0-9]*$`,
	)
	return re.MatchString(string(s))
}

func (s SimpleStyle) insertToAss(path string) error {
	m := map[string]string{
		"Name":            "Default",
		"SecondaryColour": "&H00ffffff",
		"Bold":            "0",
		"Italic":          "0",
		"Underline":       "0",
		"StrikeOut":       "0",
		"ScaleX":          "100",
		"ScaleY":          "100",
		"Spacing":         "0",
		"Angle":           "0",
		"MarginL":         "10",
		"MarginR":         "10",
		"Encoding":        "1",
	}

	fields := strings.Split(string(s), ",")

	for _, field := range fields {
		kv := strings.Split(field, "=")
		k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		m[k] = v
	}

	formatFields := []string{
		"Name",
		"Fontname",
		"Fontsize",
		"PrimaryColour",
		"SecondaryColour",
		"OutlineColour",
		"BackColour",
		"Bold",
		"Italic",
		"Underline",
		"StrikeOut",
		"ScaleX",
		"ScaleY",
		"Spacing",
		"Angle",
		"BorderStyle",
		"Outline",
		"Shadow",
		"Alignment",
		"MarginL",
		"MarginR",
		"MarginV",
		"Encoding",
	}
	format, style := "Format: ", "Style: "

	for _, i := range formatFields {
		format += fmt.Sprintf("%v, ", i)
		style += fmt.Sprintf("%v, ", m[i])
	}

	format = strings.TrimSuffix(format, ", ")
	style = strings.TrimSuffix(style, ", ")

	styleBlock := fmt.Sprintf("%v\n%v", format, style)

	reAss := regexp.MustCompile(
		`(\[V4\+ Styles])([\s\S]*?)(\[Events])`,
	)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)

	ok := reAss.MatchString(content)
	if !ok {
		return fmt.Errorf("wrong ass: %s", path)
	}

	content = reAss.ReplaceAllStringFunc(content, func(s string) string {
		parts := reAss.FindStringSubmatch(s)
		return parts[1] + fmt.Sprintf("\n%v\n\n", styleBlock) + parts[3]
	})

	return os.WriteFile(path, []byte(content), 0644)
}
