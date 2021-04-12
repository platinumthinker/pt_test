package processor

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrOperationNotFound = errors.New("operation not found")
)

type Processor interface {
	Process([]byte) ([]byte, error)
}

type Message struct {
	Operation string `json:"operation"`
	Arguments []int  `json:"arguments,omitempty"`
}

type MessageRes struct {
	Operation string `json:"operation"`
	Arguments []int  `json:"arguments,omitempty"`
	Result    int    `json:"result,omitempty"`
	Error     string `json:"error,omitempty"`
}

type JsonProcessor struct {
	commands map[string]Operation
}

func NewJsonProcessor() *JsonProcessor {
	return &JsonProcessor{
		commands: make(map[string]Operation),
	}
}

func (p *JsonProcessor) AddOperation(name string, op Operation) {
	p.commands[name] = op
}

func (p *JsonProcessor) Process(data []byte) ([]byte, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	operation, found := p.commands[msg.Operation]
	if !found {
		msgRes := MessageRes{Operation: msg.Operation,
			Arguments: msg.Arguments,
			Error:     "operation not found",
		}
		dataRes, err := json.Marshal(msgRes)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}

		return dataRes, nil
	}

	result, err := operation.Eval(msg.Arguments)
	if err != nil {
		return nil, fmt.Errorf("eval operation: %w", err)
	}
	msgRes := MessageRes{Operation: msg.Operation,
		Arguments: msg.Arguments,
		Result:    result,
	}
	dataRes, err := json.Marshal(msgRes)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return dataRes, nil
}
