package epub

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/internal"
	"os"
	"path/filepath"
	"strings"
)

type Rootfile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Rootfiles struct {
	RootFile []Rootfile `xml:"rootfile"`
}

type Container struct {
	RootFiles Rootfiles `xml:"rootfiles"`
}

func getOpsBaseName(data []byte) (string, error) {
	var c Container

	err := xml.Unmarshal(data, &c)

	if err != nil {
		return "", err
	}

	var candidates []Rootfile

	for _, rf := range c.RootFiles.RootFile {
		if rf.MediaType == "application/oebps-package+xml" {
			candidates = append(candidates, rf)
		}
	}

	if len(candidates) == 0 {
		candidates = c.RootFiles.RootFile
	}

	for _, rf := range candidates {
		if strings.HasSuffix(rf.FullPath, ".opf") {
			dir := filepath.Dir(rf.FullPath)

			return dir, nil
		}
	}
	return "", fmt.Errorf("could not find the ops dir")
}

func Highlight(md *Mode) error {

	md, err := internal.Prepare(md)
	if err != nil {
		return err
	}

	handler := func(tmpDir string) error {

		metaDirPath := filepath.Join(tmpDir, "META-INF")

		containerPath := filepath.Join(metaDirPath, "container.xml")

		_, err := os.Stat(containerPath)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(containerPath)
		if err != nil {
			return err
		}

		opsBaseName, err := getOpsBaseName(data)
		if err != nil {
			return err
		}

		targetDir := filepath.Join(tmpDir, opsBaseName)

		if err != nil {
			return err
		}

		return HighlightAllHtml(targetDir, md)
	}

	return epubWrapper(md.SrcPath, md.DstDir, handler)
}

func HighlightWithFlags(fs *flag.FlagSet, args []string) error {

	srcPath := fs.String(
		"src",
		"",
		"Source EPUB file or directory.",
	)
	dstDir := fs.String(
		"dst",
		"",
		"Destination directory. If empty, a directory will be created next to the source.",
	)
	strict := fs.Bool(
		"strict",
		false,
		"Stop processing on the first error.",
	)
	strictCodeOnly := fs.Bool(
		"sco",
		false,
		"Highlight only content inside <pre><code> tags.",
	)
	tag := fs.String(
		"tag",
		"",
		"HTML tag name. Required when --sco is enabled.",
	)
	lang := fs.String(
		"lang",
		"",
		"Programming language name (e.g. go, python, javascript).",
	)
	style := fs.String(
		"style",
		"nord",
		"Highlight style.",
	)
	bgColor := fs.String(
		"bg",
		"#303742",
		"Background color (hex).",
	)

	err := fs.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	md := &Mode{
		SrcPath:        *srcPath,
		DstDir:         *dstDir,
		Strict:         *strict,
		StrictCodeOnly: *strictCodeOnly,
		Tag:            *tag,
		Lang:           *lang,
		Style:          *style,
		BgColor:        *bgColor,
	}

	return Highlight(md)
}
