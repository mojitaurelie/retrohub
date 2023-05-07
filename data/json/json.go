package json

import (
	"encoding/json"
	"os"
	"retroHub/data"
)

type FileProvider struct {
	categories []category
}

type fileContent struct {
	Categories []category `json:"categories"`
}

type category struct {
	T  string `json:"title"`
	Ls []link `json:"links"`
}

func (c category) Title() string {
	return c.T
}

func (c category) Links() []data.Link {
	var r []data.Link
	for _, l := range c.Ls {
		r = append(r, data.Link(l))
	}
	return r
}

type link struct {
	T string `json:"title"`
	U string `json:"url"`
	D string `json:"description"`
}

func (l link) Title() string {
	return l.T
}

func (l link) URL() string {
	return l.U
}

func (l link) Description() string {
	return l.D
}

func (jfp *FileProvider) Categories() []data.Category {
	var r []data.Category
	for _, c := range jfp.categories {
		r = append(r, data.Category(c))
	}
	return r
}

func New(fp string) (*FileProvider, error) {
	if _, err := os.Stat(fp); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var c fileContent
	err = json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	jfp := new(FileProvider)
	jfp.categories = c.Categories

	return jfp, nil
}
