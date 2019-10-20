package flaw

import (
	"encoding/xml"
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

type dictionary map[string]interface{}

// MarshalXML marshals the dictionary
func (x dictionary) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	err := encoder.EncodeToken(start)
	if err != nil {
		return err
	}

	type entry struct {
		XMLName xml.Name
		Value   interface{} `xml:",chardata"`
	}

	for key, value := range x {
		encoder.Encode(entry{XMLName: xml.Name{Local: x.pascal(key)}, Value: value})
	}

	return encoder.EncodeToken(start.End())
}

func (x dictionary) pascal(text string) string {
	parts := strings.Split(text, "_")

	for index, part := range parts {
		parts[index] = strings.Title(part)
	}

	text = strings.Join(parts, "")
	return text
}
