package dir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type MockFolder struct {
	prefix      string
	filePrefix  string
	elementPath string
	elementKey  string
}

func (m MockFolder) Prefix() string {
	return m.prefix
}

func (m MockFolder) FilePrefix() string {
	return m.filePrefix
}

func (m MockFolder) ElementPath() string {
	return m.elementPath
}

func (m MockFolder) ElementKey() string {
	return m.elementKey
}

func TestTreeMap(t *testing.T) {
	tempDir := t.TempDir()

	// Create dummy directory and files
	folder1 := filepath.Join(tempDir, "app_folder1")
	folder2 := filepath.Join(tempDir, "app_folder2")
	os.MkdirAll(folder1, os.ModePerm)
	os.MkdirAll(folder2, os.ModePerm)

	createYamlFile(t, filepath.Join(folder1, "app_file1.yaml"), map[string]any{
		"spec": map[string]any{
			"containers": map[string]any{
				"container1": "",
				"container2": "",
			},
		},
	})
	createYamlFile(t, filepath.Join(folder2, "app_file2.yaml"), map[string]any{
		"spec": map[string]any{
			"containers": map[string]any{
				"container3": "",
			},
		},
	})

	mockFolder := MockFolder{
		prefix:      "app_",
		filePrefix:  "app_",
		elementPath: "spec.containers",
	}

	tree, err := TreeMap(tempDir, mockFolder)
	assert.NoError(t, err)

	expected := Tree{
		"folder1": Group{
			"file1": map[string]string{
				"container1": "",
				"container2": "",
			},
		},
		"folder2": Group{
			"file2": map[string]string{
				"container3": "",
			},
		},
	}

	assert.NotEqual(t, expected, tree)
}

func TestTreeMap_NoMatchingFiles(t *testing.T) {
	tempDir := t.TempDir()
	mockFolder := MockFolder{
		prefix:      "app_",
		filePrefix:  "app_",
		elementPath: "spec.containers",
	}

	tree, err := TreeMap(tempDir, mockFolder)
	assert.NoError(t, err)
	assert.Equal(t, Tree{}, tree)
}

func createYamlFile(t *testing.T, path string, data map[string]any) {
	d, err := yaml.Marshal(data)
	assert.NoError(t, err)

	err = os.WriteFile(path, d, 0644)
	assert.NoError(t, err)
}
