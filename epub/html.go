package epub

import (
	"bytes"
	"fmt"
	"github.com/XV-521/fileops/internal"
	goHtml "golang.org/x/net/html"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func generateCSS(dstPath string, md Mode) error {

	cmd := exec.Command("pygmentize", "-f", "html", "-S", md.Style)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	cssContent := out.Bytes()

	if md.BgColor != "" {
		bgColorCss := fmt.Sprintf(
			"\n.highlight {background-color: %v !important; padding: 1em 0.3em !important; border-radius: 1px}",
			md.BgColor,
		)
		cssContent = append(cssContent, []byte(bgColorCss)...)
	}

	fixedErrCss := "\n.err { border: none !important; background: none !important } "
	cssContent = append(cssContent, []byte(fixedErrCss)...)

	return os.WriteFile(dstPath, cssContent, 0644)
}

func getCssLinkedHtml(html string, cssPath string) string {

	if strings.Contains(html, cssPath) {
		return html
	}

	link := fmt.Sprintf(`<link rel="stylesheet" href="%v" type="text/css"/>`, cssPath)

	if strings.Contains(html, "</head>") {
		return strings.Replace(html, "</head>", link+"\n</head>", 1)
	}

	re := regexp.MustCompile(`(<html[^>]*>)([\s\S]*)`)
	out := re.ReplaceAllStringFunc(html, func(m string) string {
		parts := re.FindStringSubmatch(m)
		return fmt.Sprintf("%v\n%v", parts[1], link) + parts[2]
	})
	return out
}

func extractText(n *goHtml.Node, buf *strings.Builder) {
	if n.Type == goHtml.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, buf)
	}
}

func getHtmlContent(codeNode string) (string, error) {
	doc, err := goHtml.Parse(strings.NewReader(codeNode))
	if err != nil {
		return codeNode, err
	}

	var buf strings.Builder
	extractText(doc, &buf)

	return buf.String(), nil
}

func getHighlightedHtml(html string, md Mode) string {
	highlightCode := func(code string) (string, error) {
		cmd := exec.Command("pygmentize", "-l", md.Lang, "-f", "html")

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return code, err
		}

		if err = cmd.Start(); err != nil {
			return code, err
		}

		if _, err = io.WriteString(stdin, code); err != nil {
			return code, err
		}
		_ = stdin.Close()

		if err = cmd.Wait(); err != nil {
			return code, err
		}

		return out.String(), nil
	}

	normalize := func(codeHtml string) string {
		codeHtml = strings.Replace(
			codeHtml,
			`<div class="highlight"><pre>`,
			`<pre class="highlight">`,
			1,
		)
		codeHtml = strings.Replace(codeHtml, `</pre></div>`, `</pre>`, 1)
		return codeHtml
	}

	rePreCode := regexp.MustCompile(
		`<pre[^>]*>\s*<code[^>]*>([\s\S]*?)</code>\s*</pre>`,
	)

	html = rePreCode.ReplaceAllStringFunc(html, func(m string) string {
		parts := rePreCode.FindStringSubmatch(m)
		codeHtml := parts[1]

		codeText, err := getHtmlContent(codeHtml)
		if err != nil {
			return m
		}

		codeHtml2, err := highlightCode(codeText)
		if err != nil {
			return m
		}

		return normalize(codeHtml2)
	})

	if !md.StrictCodeOnly {

		rePreOnly := regexp.MustCompile(
			fmt.Sprintf(`<%v([^>]*)>([\s\S]*?)</%v>`, md.Tag, md.Tag),
		)

		html = rePreOnly.ReplaceAllStringFunc(html, func(m string) string {
			parts := rePreOnly.FindStringSubmatch(m)
			preAttrs := parts[1]
			codeHtml := parts[2]

			if strings.Contains(preAttrs, "highlight") {
				return m
			}

			codeText, err := getHtmlContent(codeHtml)
			if err != nil {
				return m
			}

			if !LooksLikeCode(codeText) {
				return m
			}

			codeHtml2, err := highlightCode(codeText)
			if err != nil {
				return m
			}

			return normalize(codeHtml2)
		})
	}

	return html
}

func highlightHtml(
	srcPath string,
	dstPath string,
	cssPath string,
	md Mode,
) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	html := string(data)
	html = getCssLinkedHtml(html, cssPath)
	html = getHighlightedHtml(html, md)

	err = os.WriteFile(dstPath, []byte(html), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getCssLinkedOpf(opf string, cssPath string) string {

	if strings.Contains(opf, cssPath) {
		return opf
	}

	item := fmt.Sprintf(`<item id="highlight-css" href="%v" media-type="text/css"/>`, cssPath)

	if strings.Contains(opf, "</manifest>") {
		return strings.Replace(opf, "</manifest>", item+"\n</manifest>", 1)
	}

	return opf
}

func generateOpf(srcPath string, cssPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	opf := string(data)
	opf = getCssLinkedOpf(opf, cssPath)
	return os.WriteFile(srcPath, []byte(opf), 0644)
}

func highlightAllHtml(
	srcDir string,
	cssPath string,
	md Mode,
) error {

	bm := internal.BatchMode{
		Async:  true,
		Strict: md.Strict,
	}

	filter := func(entry os.DirEntry) bool {
		if entry.IsDir() {
			return false
		}
		return true
	}

	handler := func(entry os.DirEntry) error {

		filename := entry.Name()
		path := filepath.Join(srcDir, filename)

		if internal.IsThisExt(filename, "html") || internal.IsThisExt(filename, "xhtml") {
			return highlightHtml(path, path, cssPath, md)
		}
		if internal.IsThisExt(filename, "opf") {
			return generateOpf(path, cssPath)
		}
		return nil
	}

	return internal.DoBatchWrapper(srcDir, bm, filter, handler)
}

func HighlightAllHtml(
	srcDir string,
	md Mode,
) error {

	cssPath := filepath.Join(srcDir, md.CssBasename)

	_, err := os.Stat(cssPath)
	if os.IsNotExist(err) {
		err = generateCSS(cssPath, md)
		if err != nil {
			return err
		}
	}

	return highlightAllHtml(srcDir, filepath.Base(cssPath), md)
}
