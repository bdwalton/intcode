package intcode

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

type Getter interface {
	Get() int
}

type keyboard struct {
	in *bufio.Reader
}

func NewKeyboard(in io.Reader) *keyboard {
	return &keyboard{bufio.NewReader(in)}
}

func (k *keyboard) Get() int {
	inp, err := k.in.ReadString('\n')
	if err != nil {
		log.Fatalf("Couldn't read input: %v\n", err)
	}

	v, err := strconv.Atoi(inp[0 : len(inp)-1])
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	return v
}

type networkGet struct {
	in chan int
}

func NewNetworkGet(c chan int) *networkGet {
	return &networkGet{c}
}

func (ng *networkGet) Get() int {
	return <-ng.in
}
