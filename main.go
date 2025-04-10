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
	del  bool
}

func main() {
	root := flag.String("root", ".", "Root directory to start.")
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete Files")
	ext := flag.String("ext", "", "File extentions to filter out")
	size := flag.Int64("size", 0, "Minimum file size")

	flag.Parse()

	c := config{
		list: *list,
		ext:  *ext,
		size: *size,
		del:  *del,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

			if c.del {
				return delFile(path)
			}

			return listFile(path, out)
		})
}
