package intcode

import (
	"bufio"
	"fmt"
	"io"
)

type Putter interface {
	Put(o int)
}

type display struct {
	out *bufio.Writer
}

func NewDisplay(out io.Writer) *display {
	return &display{bufio.NewWriter(out)}
}

func (d *display) Put(i int) {
	d.out.WriteString(fmt.Sprintf("%d\n", i))
	d.out.Flush()
}

type networkPut struct {
	out chan int
}

func NewNetworkPut(c chan int) *networkPut {
	return &networkPut{c}
}

func (np *networkPut) Put(i int) {
	np.out <- i
}
