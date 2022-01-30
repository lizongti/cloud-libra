package bar

import (
	"bytes"
	"fmt"
)

type Bar struct {
	percent int
	current int
	total   int
	bar     bytes.Buffer
	element string
}

func NewBar(total int) *Bar {
	bar := &Bar{}
	bar.total = total
	bar.element = "="
	return bar
}

func (bar *Bar) Move(n int) {
	bar.current += n
	last := bar.percent
	bar.percent = bar.getPercent()
	for i := last; i < bar.percent; i++ {
		bar.bar.WriteString(bar.element)
	}
	fmt.Printf("\r[%-100s]%3d%%  %8d/%d", bar.bar.String(), bar.percent, bar.current, bar.total)
}

func (bar *Bar) Close() {
	fmt.Println()
}

func (bar *Bar) getPercent() int {
	return int(float32(bar.current) / float32(bar.total) * 100)
}
