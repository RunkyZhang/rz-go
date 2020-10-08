package main

import (
	"fmt"
)

func main() {
	page := &Page{}
	var runable Runable
	runable = page

	fmt.Println(page.build(), runable.run(), runable.run())
}

type Runable interface {
	run() int
}

type Base struct {
	Name string
}

func (b *Base) build() string {
	return b.Name
}

type Page struct {
	Base
	Age int
}

func (p *Page) run() int {
	return p.Age
}