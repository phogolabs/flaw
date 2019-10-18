package flaw

import (
	"go/build"
	"path/filepath"
	"strings"
)

func relative(path string) string {
	if root := build.Default.GOPATH; root != "" {
		if strings.HasPrefix(path, root) {
			if file, err := filepath.Rel(root, path); err == nil {
				const (
					src = "src/"
					pkg = "pkg/mod/"
				)

				switch {
				case strings.HasPrefix(file, src):
					path = strings.TrimPrefix(file, src)
				case strings.HasPrefix(file, pkg):
					path = strings.TrimPrefix(file, pkg)
				}
			}
		}
	}

	return path
}

func function(name string) string {
	withoutPath := name[strings.LastIndex(name, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	name = withoutPackage
	name = strings.Replace(name, "(", "", 1)
	name = strings.Replace(name, "*", "", 1)
	name = strings.Replace(name, ")", "", 1)

	return name
}
