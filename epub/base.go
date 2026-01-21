package epub

import (
	"strings"
)

var middleSet = []string{
	"(", "{", "[",
	"=", ":=", "=>", "::",
	">",
}

var freeSet = []string{
	"<", ".", "//", "#",
}

var superBeginSet = []string{
	"import ", "package ", "type ", "def ", "func ", "class ", "var ", "let ", "const ",
}

var normalBeginSet = []string{
	"function ", "return ", "void ", "assert ", "yield ", "goto ", "defer ", "if ", "for ",
	"else", "while", "do", "elif", "go", "break", "continue", "switch", "case", "default",
	"select",
}

var beginSet = append(normalBeginSet, superBeginSet...)

var normalSet = append(freeSet, append(beginSet, middleSet...)...)

var superSet = []string{"`", `"""`, `'''`, "/*", "<--"}

func hasPrefixOne(content string, set []string) bool {
	for _, s := range set {
		if strings.HasPrefix(content, s) {
			return true
		}
	}
	return false
}

func containsOne(content string, set []string) bool {
	for _, s := range set {
		if strings.Contains(content, s) {
			return true
		}
	}
	return false
}

func mustBeCode(content string) bool {
	has := hasPrefixOne(content, superBeginSet)
	if has {
		return true
	}
	return containsOne(content, superSet)
}

func hasCodeSignal(content string) bool {
	return containsOne(content, normalSet)
}

func conformLineFmt(content string) bool {
	badCount := 0
	linesCount := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		linesCount += 1
		contains := containsOne(line, middleSet)
		if !contains {
			has := hasPrefixOne(line, beginSet) || hasPrefixOne(line, freeSet)
			if !has {
				badCount += 1
				continue
			}
		}
		has := hasPrefixOne(line, middleSet)
		if has {
			badCount += 1
			continue
		}
	}

	return float32(badCount)/float32(linesCount) < 0.8
}

func LooksLikeCode(content string) bool {
	content = strings.TrimSpace(content)

	if content == "" {
		return false
	}

	be := mustBeCode(content)
	if be {
		return true
	}

	has := hasCodeSignal(content)
	if !has {
		return false
	}

	conform := conformLineFmt(content)
	if !conform {
		return false
	}
	return true
}
