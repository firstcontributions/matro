package commands

import (
	"io"
	"strings"
)

// CommandWriter can be used to write help texts to
// stdout/tcp or any given writer interface
type CommandWriter struct {
	writer io.Writer
}

// NewCommandWriter returns an instance of CommandWriter
func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

// Write writes given string to the registered writer interface
func (c *CommandWriter) Write(msg ...string) error {
	_, err := c.writer.Write([]byte(strings.Join(msg, " ")))
	return err
}
