# Intcode Computer

Intcode is a computer that runs intcode programs

This computer was defined as part of
https://adventofcode.com/2019. This implementation was used for
solutions to those puzzles.

## Description

The intcode computer has main memory and an instruction pointer (IP)
and relative base offset (RBO) registers. The main memory is stores
integer values and can be sized at creation time. The computer has no
distinction between program and data memory and no write protection,
so programs may write to any memory address. At initialization, the IP
and RBO registers are both initialized to zero, making the computer
read its first instruction from memory address 0. When run, the
computer will execute instructions until it either faults or halts.

### Operations

Code | Name | Description
-----|------|-------------------------------------------------------------------------------------
1    | ADD  | Add first param to second and store in third
2    | MUL  | Multiply first param by second param and store in third
3    | INP  | Read value from input, store in first param
4    | OUT  | Output first param value
5    | JIT  | Jump if true; If first param != 0, set PC to second param, no-op otherwise
6    | JIF  | Jump if false; If first param == 0, set PC to second param, no-op otherwise
7    | LT   | Less than; If first param < second param, store 1 in third param, otherwise store 0
8    | EQ   | Equal; If first param == second param, store 1 in third param, otherwise store 0
9    | ARB  | Increase (or decrease if negative) RBO by first param value
99   | HALT | Stop all processing

Operations are subject to different parameter modes. Parameter modes
are stored as part of the opcode and are 0 if not present.

### Parameter Modes

Code | Name      | Description
-----|-----------|----------------------------------------------------------------------------------------
0    | POSITION  | The parameter value represents a memory address where the value is stored.
1    | IMMEDIATE | The parameter value represents the value to to be used.
2    | RELATIVE  | The parameter value represents an offset to the RBO register where the value is stored.

### Operators

Combining opcodes and parameter modes, operators are stored as single integers as:

```
ABCDE
 1202

DE - two-digit opcode,      02 == opcode 2
 C - mode of 1st parameter,  2 == relative mode
 B - mode of 2nd parameter,  1 == immediate mode
 A - mode of 3rd parameter,  0 == position mode,
                                  omitted due to being a leading zero
```

### Instructions
An instruction is thus comprised of opcodes and a variable number of
parameters which are stored in adjacent memory locations to the opcode.


## Execution

After executing an instruction, the IP is incremented to point to the
memory address beyond the last parameter consumed by the previous
instruction unless branching occurred, at which point IP is set to the
outcome of the branch decision.
