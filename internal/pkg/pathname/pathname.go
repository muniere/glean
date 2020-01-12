package pathname

import (
	"errors"
	"path/filepath"
)

type Pathname struct {
	value string
}

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func Base(path string) string {
	return filepath.Base(path)
}

func Dir(path string) string {
	return filepath.Dir(path)
}

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

func Leaf(path string, level int) (string, error) {
	leaf, err := New(path).Leaf(level)
	if err != nil {
		return "", err
	}
	return leaf.String(), nil
}

func MustLeaf(path string, level int) string {
	leaf, err := New(path).Leaf(level)
	if err != nil {
		panic(err)
	}
	return leaf.String()
}

func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func MustAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}

func New(path string) *Pathname {
	return &Pathname{value: path}
}

func (p *Pathname) String() string {
	return p.value
}

func (p *Pathname) Prepend(elem ...string) *Pathname {
	s := filepath.Join(elem...)
	return New(filepath.Join(s, p.value))
}

func (p *Pathname) Append(elem ...string) *Pathname {
	s := filepath.Join(elem...)
	return New(filepath.Join(p.value, s))
}

func (p *Pathname) Base() *Pathname {
	return New(filepath.Base(p.value))
}

func (p *Pathname) Dir() *Pathname {
	return New(filepath.Dir(p.value))
}

func (p *Pathname) Leaf(level int) (*Pathname, error) {
	if level <= 0 {
		return nil, errors.New("level must be greater than 0")
	}

	branch := New(p.value)

	var elems []string

	for i := 0; i < level; i++ {
		elems = append(elems, branch.Base().value)
		branch = branch.Dir()
	}

	return New(filepath.Join(elems...)), nil
}

func (p *Pathname) MustLeaf(level int) *Pathname {
	leaf, err := p.Leaf(level)
	if err != nil {
		panic(err)
	}
	return leaf
}

func (p *Pathname) Abs() (*Pathname, error) {
	abs, err := filepath.Abs(p.value)
	if err != nil {
		return nil, err
	}
	return New(abs), nil
}

func (p *Pathname) MustAbs() *Pathname {
	abs, err := p.Abs()
	if err != nil {
		panic(err)
	}
	return abs
}

func (p *Pathname) IsAbs() bool {
	return filepath.IsAbs(p.value)
}
