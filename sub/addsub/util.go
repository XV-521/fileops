package addsub

import (
	"fmt"
	"github.com/XV-521/fileops/v2/internal/util"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func escapeFFmpegFilterValue(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`:`, `\:`,
		`'`, `\'`,
		`,`, `\,`,
		`[`, `\[`,
		`]`, `\]`,
		`;`, `\;`,
	)
	return replacer.Replace(s)
}

func addHardSub(
	videoPath string,
	subPath string,
	dstPath string,
) error {
	subDir, subName := filepath.Dir(subPath), filepath.Base(subPath)
	vf := "subtitles=" + escapeFFmpegFilterValue(subName)

	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath,
		"-vf", vf,
		"-c:a", "copy",
		dstPath,
	)
	cmd.Dir = subDir
	return util.CmdWrap(cmd)
}

func addSoftSub(
	videoPath string,
	subPath string,
	dstPath string,
) error {
	mp4Like := []string{".mp4", ".mov"}
	mkvLike := []string{".mkv"}

	ext := strings.ToLower(filepath.Ext(dstPath))
	r0 := slices.Contains(mp4Like, ext)
	r1 := slices.Contains(mkvLike, ext)

	if !(r0 || r1) {
		return fmt.Errorf("%s is not supported", dstPath)
	}

	args := []string{
		"-i", videoPath,
		"-i", subPath,
		"-map", "0:v",
		"-map", "0:a?",
		"-map", "1:0",
		"-c:v", "copy",
		"-c:a", "copy",
	}

	if r0 {
		args = append(args, "-c:s", "mov_text")
	}
	if r1 {
		args = append(args, "-c:s", "copy")
	}

	args = append(args, dstPath)

	cmd := exec.Command("ffmpeg", args...)
	return util.CmdWrap(cmd)
}
