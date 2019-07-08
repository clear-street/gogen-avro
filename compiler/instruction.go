package compiler

import (
	"fmt"

	"github.com/clear-street/gogen-avro/vm"
)

type irInstruction interface {
	VMLength() int
	Name() string
	CompileToVM(*irProgram) ([]vm.Instruction, error)
}

type literalIRInstruction struct {
	instruction vm.Instruction
	name        string
}

func (b *literalIRInstruction) VMLength() int {
	return 1
}

func (b *literalIRInstruction) Name() string {
	return fmt.Sprintf("Literal: %v", b.name)
}

func (b *literalIRInstruction) CompileToVM(_ *irProgram) ([]vm.Instruction, error) {
	b.instruction.Name = b.name
	return []vm.Instruction{b.instruction}, nil
}

type methodCallIRInstruction struct {
	method string
}

func (b *methodCallIRInstruction) VMLength() int {
	return 1
}

func (b *methodCallIRInstruction) Name() string {
	return fmt.Sprintf("Method: %v", b.method)
}

func (b *methodCallIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	method, ok := p.methods[b.method]
	if !ok {
		return nil, fmt.Errorf("Unable to call unknown method %q", b.method)
	}
	return []vm.Instruction{vm.Instruction{vm.Call, method.offset, b.Name()}}, nil
}

type blockStartIRInstruction struct {
	blockId int
}

func (b *blockStartIRInstruction) VMLength() int {
	return 8
}

func (b *blockStartIRInstruction) Name() string {
	return "Block start"
}

// At the beginning of a block, read the length into the Long register
// If the block length is 0, jump past the block body because we're done
// If the block length is negative, read the byte count, throw it away, multiply the length by -1
// Once we've figured out the number of iterations, push the loop length onto the loop stack
func (b *blockStartIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.Read, vm.Long, b.Name()},
		vm.Instruction{vm.EvalEqual, 0, b.Name()},
		vm.Instruction{vm.CondJump, block.end + 5, b.Name()},
		vm.Instruction{vm.EvalGreater, 0, b.Name()},
		vm.Instruction{vm.CondJump, block.start + 7, b.Name()},
		vm.Instruction{vm.Read, vm.UnusedLong, b.Name()},
		vm.Instruction{vm.MultLong, -1, b.Name()},
		vm.Instruction{vm.PushLoop, 0, b.Name()},
	}, nil
}

type blockEndIRInstruction struct {
	blockId int
}

func (b *blockEndIRInstruction) Name() string {
	return "Block end"
}

func (b *blockEndIRInstruction) VMLength() int {
	return 5
}

// At the end of a block, pop the loop count and decrement it. If it's zero, go back to the very
// top to read a new block. otherwise jump to start + 7, which pushes the value back on the loop stack
func (b *blockEndIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.PopLoop, 0, b.Name()},
		vm.Instruction{vm.AddLong, -1, b.Name()},
		vm.Instruction{vm.EvalEqual, 0, b.Name()},
		vm.Instruction{vm.CondJump, block.start, b.Name()},
		vm.Instruction{vm.Jump, block.start + 7, b.Name()},
	}, nil
}

type switchStartIRInstruction struct {
	switchId int
	size     int
	errId    int
}

func (s *switchStartIRInstruction) VMLength() int {
	return 2*s.size + 1
}

func (s *switchStartIRInstruction) Name() string {
	return "Switch start"
}

func (s *switchStartIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	body := []vm.Instruction{}
	for value, offset := range sw.cases {
		body = append(body, vm.Instruction{vm.EvalEqual, value, s.Name()})
		body = append(body, vm.Instruction{vm.CondJump, offset + 1, s.Name()})
	}

	body = append(body, vm.Instruction{vm.Halt, s.errId, s.Name()})
	return body, nil
}

type switchCaseIRInstruction struct {
	switchId    int
	writerIndex int
	readerIndex int
}

func (s *switchCaseIRInstruction) VMLength() int {
	return 3
}

func (s *switchCaseIRInstruction) Name() string {
	return "Switch case"
}

func (s *switchCaseIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	return []vm.Instruction{
		vm.Instruction{vm.Jump, sw.end, s.Name()},
		vm.Instruction{vm.AddLong, s.readerIndex - s.writerIndex, s.Name()},
		vm.Instruction{vm.Set, vm.Long, s.Name()},
	}, nil
}

type switchEndIRInstruction struct {
	switchId int
}

func (s *switchEndIRInstruction) VMLength() int {
	return 0
}

func (s *switchEndIRInstruction) Name() string {
	return "Switch end"
}

func (s *switchEndIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	return []vm.Instruction{}, nil
}
