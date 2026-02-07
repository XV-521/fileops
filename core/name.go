package core

type Namer interface {
	Next(ei EntryInfo) string
}
