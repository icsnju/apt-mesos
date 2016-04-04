package docker

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

type Dockerfile struct {
	Instructions         []*Instruction
	InternalInstructions []*Instruction
	Registry             string
	Repository           string
}

type Instruction struct {
	Command   string
	Arguments []string
}

var INTERNAL = []string{"REGISTRY", "REPOSITORY", "BUILD_CPU", "BUILD_MEM"}

func NewDockerfile(file string, registry string) *Dockerfile {
	dockerfile := &Dockerfile{
		Instructions:         []*Instruction{},
		InternalInstructions: []*Instruction{},
		Registry:             registry,
		Repository:           "",
	}
	dockerfile.parse(file)
	return dockerfile
}

func (d *Dockerfile) parse(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// filter comments
		trimedLine := strings.TrimSpace(line)
		if len(trimedLine) == 0 || strings.HasPrefix(trimedLine, "#") {
			continue
		}

		lineBuf := strings.TrimLeft(line, " ")
		lineBuf = strings.TrimRight(lineBuf, " ")

		d.addInstruction(lineBuf)
	}
}

func (d *Dockerfile) addInstruction(instruction string) {
	parts := strings.Split(instruction, " ")
	command := strings.ToUpper(parts[0])
	arguments := parts[1:]

	if command == "FROM" && len(d.Registry) > 0 {
		parts = strings.Split(arguments[0], "/")
		if len(parts) <= 2 {
			parts[0] = d.Registry
			arguments[0] = strings.Join(parts, "/")
		}
	}

	d.Instructions = append(d.Instructions, &Instruction{
		Command:   command,
		Arguments: arguments,
	})
}

func (d *Dockerfile) Build() string {
	buffer := bytes.NewBufferString("")
	for _, i := range d.Instructions {
		buffer.WriteString(i.Command + " " + strings.Join(i.Arguments, " ") + "\n")
	}
	return buffer.String()
}

func (d *Dockerfile) HasLocalSources() bool {
	for _, i := range d.Instructions {
		if i.Command == "ADD" || i.Command == "COPY" {
			return true
		}
	}
	return false
}
