package main

import (
	"github.com/XV-521/fileops/public"
	"github.com/XV-521/fileops/zip"
	"log"
)

func main() {
	zt := public.ZipT

	md := &zip.Mode{
		SrcDir: "/Users/xiaoxiao/Desktop/Dir",
		DstDir: "/Users/xiaoxiao/Desktop/dirToT",
		ZT:     zt,
	}
	err := zip.DoBatch(md)
	if err != nil {
		log.Fatal(err)
	}
}
