package intcode

import (
	"fmt"
)

// Opcodes
const (
	ADD  = 1  // Add first param value to second and store in third
	MUL  = 2  // Multiply first param value by second and store in third
	INP  = 3  // Read value from input, store in first param value
	OUT  = 4  // Output first param value
	JIT  = 5  // If first param value !=0, set PC to second param value, no-op otherwise
	JIF  = 6  // If first param value ==0, set PC to second param value, no-op otherwise
	LT   = 7  // If first param value < second param value, store 1 in third param value, otherwise store 0
	EQ   = 8  // If first param value == second param value, store 1 in third param value, otherwise store 0
	ARB  = 9  // Increase (or decrease if negative), current RBO by first param value
	HALT = 99 // Stop all processing
)

type opcode int

type mode int // Instruction mode

// Modes
const (
	POSITION  = iota // Value stored in parameter is a memory address - read value from that address
	IMMEDIATE        // Use value in parameter
	RELATIVE         // Value stored is relative to RBO register
)

func (m mode) String() (s string) {
	switch m {
	case POSITION:
		s = "*"
	case IMMEDIATE:
		s = "-"
	case RELATIVE:
		s = "^"
	}

	return
}

type instruction struct {
	op        opcode
	modes     []mode
	nOperands int
}

func (i *instruction) String() string {
	var op string
	switch i.op {
	case ADD:
		op = "ADD"
	case MUL:
		op = "MUL"
	case INP:
		op = "INP"
	case OUT:
		op = "OUT"
	case JIT:
		op = "JIT"
	case JIF:
		op = "JIF"
	case LT:
		op = "LT"
	case EQ:
		op = "EQ"
	case ARB:
		op = "ARB"
	case HALT:
		op = "HALT"
	default:
		op = "?!? UNKNOWN ?!?"
	}

	var m string
	for x := 0; x < i.nOperands; x++ {
		m += i.modes[x].String()
	}

	return fmt.Sprintf("%s %s", op, m)

}

const (
	HALTED = iota
	RUNNING
)

type cpuState int

type Comp struct {
	memory []int
	pc     int // program counter / instruction pointer
	rbo    int // relative base offset
	state  cpuState
	error  bool // True if CPU is halted due to error
	name   string
	in     Getter // Keyboard
	out    Putter // Display
}

func NewComp(name string, memsize int, prog []int, in Getter, out Putter) *Comp {
	m := make([]int, memsize)
	copy(m, prog)
	return &Comp{memory: m, pc: 0, state: RUNNING, name: name, in: in, out: out}
}

func (c *Comp) GetName() string {
	return c.name
}

func (c *Comp) Run() {
	for !c.Halted() {
		c.Step()
	}
}

func (c *Comp) decodeInstruction(i int) *instruction {
	inst := &instruction{op: opcode(i % 1e2), modes: make([]mode, 0)}

	switch inst.op {
	case ADD, MUL, LT, EQ:
		inst.nOperands = 3
	case INP, OUT, ARB:
		inst.nOperands = 1
	case JIT, JIF:
		inst.nOperands = 2
	case HALT:
		inst.nOperands = 0
	}

	// Iterate over the digits in the 3 higher order digits of the
	// instruction and append each one to modes.
	modes := i / 1e2
	for x := 0; x < inst.nOperands; x++ {
		inst.modes = append(inst.modes, mode(modes%10))
		modes = modes / 10
	}

	return inst
}

func (c *Comp) getInst() int {
	return c.GetMemory(c.pc)
}

func (c *Comp) GetMemory(addr int) int {
	return c.memory[addr]
}

// Like getOperand, but used to bypass mode handling for destination
// addresses for writing results.
func (c *Comp) getLocation(i *instruction, n int) int {
	switch i.modes[n-1] {
	case POSITION:
		return c.GetMemory(c.pc + n)
	case RELATIVE:
		return c.GetMemory(c.pc+n) + c.rbo
	}
	return c.memory[c.pc+n]
}

// Retrieve operand n for instruction i.
// Operands are indexed from 1, not 0.
func (c *Comp) getOperand(i *instruction, n int) int {
	switch i.modes[n-1] {
	case IMMEDIATE:
		return c.GetMemory(c.pc + n)
	case POSITION:
		return c.GetMemory(c.GetMemory(c.pc + n))
	case RELATIVE:
		return c.GetMemory(c.GetMemory(c.pc+n) + c.rbo)
	}

	return c.memory[c.pc+n]
}

func (c *Comp) Break() {
	c.error = true
	c.state = HALTED
}

func (c *Comp) Broken() bool {
	return c.error
}

func (c *Comp) Halt() {
	c.state = HALTED
}

func (c *Comp) Halted() bool {
	return c.state == HALTED
}

func (c *Comp) PrintMemory() {
	var cnt int
	for i := 0; i < len(c.memory); i++ {
		var ptr string
		if i == c.pc {
			ptr = "<"
		}
		fmt.Printf("%04d: %08d%s\t", i, c.memory[i], ptr)
		cnt++
		if cnt%10 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n\n")
}

func (c *Comp) SetMemory(addr, val int) {
	c.memory[addr] = val
}

func (c *Comp) Step() {
	var branched bool

	i := c.decodeInstruction(c.getInst())
	switch i.op {
	case JIT:
		if c.getOperand(i, 1) != 0 {
			c.pc = c.getOperand(i, 2)
			branched = true
		}
	case JIF:
		if c.getOperand(i, 1) == 0 {
			c.pc = c.getOperand(i, 2)
			branched = true
		}
	case LT:
		var val int
		if c.getOperand(i, 1) < c.getOperand(i, 2) {
			val = 1
		}
		c.SetMemory(c.getLocation(i, 3), val)
	case EQ:
		var val int
		if c.getOperand(i, 1) == c.getOperand(i, 2) {
			val = 1
		}
		c.SetMemory(c.getLocation(i, 3), val)
	case ADD:
		c.SetMemory(c.getLocation(i, 3), c.getOperand(i, 1)+c.getOperand(i, 2))
	case MUL:
		c.SetMemory(c.getLocation(i, 3), c.getOperand(i, 1)*c.getOperand(i, 2))
	case INP:
		c.SetMemory(c.getLocation(i, 1), c.in.Get())
	case OUT:
		c.out.Put(c.getOperand(i, 1))
	case ARB:
		c.rbo += c.getOperand(i, 1)
	case HALT:
		c.Halt()
	default:
		c.Break()
	}

	if !branched {
		c.pc += i.nOperands + 1 // Always 1 address beyond the number of operands
	}
}
