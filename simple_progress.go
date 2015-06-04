package main

import (
	"fmt"
	"strings"
)

type simpleProgress struct {
	name      string
	width     int
	lastWrote float32
}

func (p *simpleProgress) Init() {
	usableWidth := p.width - 6
	name := p.name
	if len(p.name) > usableWidth {
		name = name[0:usableWidth]
	}
	fmt.Printf("|= %s %s=|\n", name, strings.Repeat("=", usableWidth-len(name)))
}

func (p *simpleProgress) StartPhase(phase string) {
	p.lastWrote = 0.0
	usableWidth := p.width - 6

	if len(phase) > usableWidth {
		phase = phase[0:usableWidth]
	}
	fmt.Printf("|  %s %s |\n|", phase, strings.Repeat(" ", usableWidth-len(phase)))
}

func (p *simpleProgress) Progress(percent float32) {
	// log.Printf("percent: %v\n", percent)
	diff := int((percent - p.lastWrote) * float32(p.width))
	if diff > 1 {
		p.lastWrote = percent
		fmt.Print(strings.Repeat("~", diff))
	}
}

func (p *simpleProgress) Complete() {
	p.Progress(1.0)
	fmt.Println("|")
}
