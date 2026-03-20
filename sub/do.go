package sub

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XV-521/fileops/v2/sub/addsub"
	"github.com/XV-521/fileops/v2/sub/createsub"
	"strings"
)

func DoBatch(md *Mode) error {

	isThisExt := func(ext string) bool {
		clean := func(s string) string {
			return "." + strings.Trim(strings.ToLower(strings.TrimSpace(s)), ".")
		}

		return clean(md.SubExt) == clean(ext)
	}

	cMD := &createsub.Mode{
		SrcDir: md.SrcDir,
		Model:  md.Model,
		Lang:   md.Lang,
		OSrt:   isThisExt(".srt"),
		OVtt:   isThisExt(".vtt"),
		OAss:   isThisExt(".ass"),
		OTxt:   false,
		Rec:    md.Rec,
		Strict: md.Strict,
	}

	err := cMD.Check()
	if err != nil {
		return err
	}

	aMD := &addsub.Mode{
		SrcDir:  md.SrcDir,
		DstDir:  md.DstDir,
		SubExt:  md.SubExt,
		ST:      md.ST,
		Style:   md.Style,
		MoveSub: md.SrcDir != md.DstDir && md.RetainSub,
		Rec:     md.Rec,
		Strict:  md.Strict,
	}

	err = aMD.Check()
	if err != nil {
		return err
	}

	err = createsub.DoBatch(cMD)
	if err != nil && md.Strict {
		return err
	}
	return addsub.DoBatch(aMD)
}

func DoBatchWithFlags(fs *flag.FlagSet, args []string) error {
	srcDir := fs.String(
		"src",
		"",
		"Source directory.",
	)

	dstDir := fs.String(
		"dst",
		"",
		"Destination directory.",
	)

	model := fs.String(
		"model",
		"ggml-large-v3-turbo.bin",
		"Model filename (e.g., ggml-large-v3-turbo.bin).",
	)

	lang := fs.String(
		"lang",
		"en",
		"Language.",
	)

	style := fs.String(
		"style",
		"Fontname=Arial,Fontsize=16,PrimaryColour=&H00ffffff,OutlineColour=&H0,BackColour=&H0,BorderStyle=1,Outline=1,Shadow=0,Alignment=2,MarginV=10",
		"Style of additional substitutions.",
	)

	subExt := fs.String(
		"subext",
		"",
		"Subtitle file extension",
	)

	st := fs.Int(
		"st",
		int(addsub.SubUn),
		fmt.Sprintf("Sub type: { %v: soft sub, %v: hard sub}.", addsub.SubS, addsub.SubH),
	)

	retainSub := fs.Bool(
		"rms",
		false,
		"Retain sub file in destination directory.",
	)

	rec := fs.Bool(
		"rec",
		false,
		"Recursive.",
	)

	strict := fs.Bool(
		"strict",
		false,
		"Stop processing on the first error.",
	)

	err := fs.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	md := &Mode{
		SrcDir:    *srcDir,
		DstDir:    *dstDir,
		Model:     *model,
		Lang:      *lang,
		Style:     addsub.SimpleStyle(*style),
		SubExt:    *subExt,
		ST:        addsub.SubType(*st),
		RetainSub: *retainSub,
		Rec:       *rec,
		Strict:    *strict,
	}

	return DoBatch(md)
}
