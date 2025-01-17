package statedumper

import (
	"fmt"
	"github.com/ethereum-optimism/optimism/l2geth/common"
	"io"
	"os"
	"sync"
)

type StateDumper interface {
	WriteETH(address common.Address)
	WriteMessage(sender common.Address, msg []byte)
}

var DefaultStateDumper StateDumper

func NewStateDumper() StateDumper {
	path := os.Getenv("L2GETH_STATE_DUMP_PATH")
	if path == "" {
		return &noopStateDumper{}
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0o755)
	if err != nil {
		panic(err)
	}

	return &FileStateDumper{
		f: f,
	}
}

type FileStateDumper struct {
	f   io.Writer
	mtx sync.Mutex
}

func (s *FileStateDumper) WriteETH(address common.Address) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, err := s.f.Write([]byte(fmt.Sprintf("ETH|%s\n", address.Hex()))); err != nil {
		panic(err)
	}
}

func (s *FileStateDumper) WriteMessage(sender common.Address, msg []byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, err := s.f.Write([]byte(fmt.Sprintf("MSG|%s|%x\n", sender.Hex(), msg))); err != nil {
		panic(err)
	}
}

type noopStateDumper struct {
}

func (n *noopStateDumper) WriteETH(address common.Address) {
}

func (n *noopStateDumper) WriteMessage(sender common.Address, msg []byte) {
}

func init() {
	DefaultStateDumper = NewStateDumper()
}

func WriteETH(address common.Address) {
	DefaultStateDumper.WriteETH(address)
}

func WriteMessage(sender common.Address, msg []byte) {
	DefaultStateDumper.WriteMessage(sender, msg)
}
