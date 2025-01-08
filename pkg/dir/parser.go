package dir

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Folder interface {
	Prefix() string
	FilePrefix() string
	ElementPath() string
	ElementKey() string
}

type Tree map[string]Group
type Group map[string]map[string]string

type parser struct {
	tree   Tree
	folder Folder
}

func TreeMap(root string, f Folder) (Tree, error) {
	p := &parser{
		tree:   make(Tree),
		folder: f,
	}

	if err := p.walkFolder(root); err != nil {
		return nil, err
	}

	if err := p.walkFile(root); err != nil {
		return nil, err
	}

	return p.tree, nil
}

func (p *parser) walkFolder(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			prefix := p.folder.Prefix()

			if strings.HasPrefix(filepath.Base(path), prefix) {
				folderName := strings.TrimPrefix(path, prefix)

				_, ok := p.tree[folderName]
				if !ok {
					p.tree[folderName] = make(Group)
				}
			}
		}
		return nil
	})
}

func (p *parser) walkFile(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			prefix := p.folder.FilePrefix()

			if strings.HasPrefix(filepath.Base(path), prefix) {
				folderPath := filepath.Dir(path)

				if strings.HasPrefix(filepath.Base(folderPath), p.folder.Prefix()) {
					folderName := strings.TrimPrefix(folderPath, p.folder.Prefix())
					fileName := strings.TrimPrefix(filepath.Base(path), prefix)

					_, ok := p.tree[folderName][fileName]
					if !ok {
						p.tree[folderName][fileName] = p.walkElement(path)
					}
				}
			}
		}
		return nil
	})
}

func (p *parser) walkElement(path string) map[string]string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	element := map[string]any{}
	if err := yaml.Unmarshal(data, &element); err != nil {
		return nil
	}

	elementPath := strings.Split(p.folder.ElementPath(), ".")
	es := map[string]string{}

	log.Println(elementPath)

	for _, k := range elementPath {
		log.Println(element[k])
		element = element[k].(map[string]any)
	}

	for _, v := range element {
		es[v.(string)] = ""
	}

	return es
}
