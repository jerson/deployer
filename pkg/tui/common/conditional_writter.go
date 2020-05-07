package common

import "io"

// ConditionalFunction ...
type ConditionalFunction func() (success bool, err error)

// ConditionalWriter ...
type ConditionalWriter struct {
	w         io.Writer
	condition ConditionalFunction
}

// NewConditionalWriter ...
func NewConditionalWriter(condition ConditionalFunction, writer io.Writer) *ConditionalWriter {
	return &ConditionalWriter{
		w:         writer,
		condition: condition,
	}
}
func (c *ConditionalWriter) Write(p []byte) (n int, err error) {

	result, err := c.condition()
	if err != nil {
		return 0, err
	}
	if !result {
		return len(p), nil
	}

	return c.w.Write(p)

}
