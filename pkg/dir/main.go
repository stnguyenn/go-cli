package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*FileNode
}

func main() {
	root := "/usr/local"
	rootNode := parseDirectory(root)

	// Print the directory tree
	printTree(rootNode, 0)
}

func parseDirectory(root string) *FileNode {
	rootNode := &FileNode{Path: root, Name: filepath.Base(root), IsDir: true}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		node := &FileNode{Name: info.Name(), Path: path, IsDir: info.IsDir()}
		parentPath := filepath.Dir(path)
		parent := findNode(rootNode, parentPath)
		if parent != nil {
			parent.Children = append(parent.Children, node)
		}

		return nil
	})

	return rootNode
}

func findNode(node *FileNode, path string) *FileNode {
	if node.Path == path {
		return node
	}

	for _, child := range node.Children {
		if found := findNode(child, path); found != nil {
			return found
		}
	}

	return nil
}

func printTree(node *FileNode, level int) {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	fmt.Println(indent, node.Name)

	for _, child := range node.Children {
		printTree(child, level+1)
	}
}
