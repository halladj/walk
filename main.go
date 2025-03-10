package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type config struct {
	ext  string
	size int64
	list bool
}

func main() {
	root := flag.String("root", ".", "Root directory to start.")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extentions to filter out")
	size := flag.Int64("size", 0, "Minimum file size")

	flag.Parse()

	c := config{
		list: *list,
		ext:  *ext,
		size: *size,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

}

func run(root string, out io.Writer, c config) error {
	return filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, c.ext, c.size, info) {
				return nil
			}

			if c.list {
				return listFile(path, out)
			}

			return listFile(path, out)
		})
}

func filterOut(
	path,
	ext string,
	minSize int64,
	info os.FileInfo,
) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}
func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
