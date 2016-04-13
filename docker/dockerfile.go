package docker

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path"
	"strings"
)

type Dockerfile struct {
	ID                   string
	Instructions         []*Instruction
	InternalInstructions []*Instruction
	Path                 string
	Registry             string
	Repository           string
}

// Instruction of dockerfile
type Instruction struct {
	Command   string
	Arguments []string
}

// INTERNAL is internal build commands
var INTERNAL = []string{"REGISTRY", "REPOSITORY", "BUILD_CPU", "BUILD_MEM"}

// NewDockerfile return new dockerfile
func NewDockerfile(ID, filePath, registry string) *Dockerfile {
	dockerfile := &Dockerfile{
		ID:                   ID,
		Instructions:         []*Instruction{},
		InternalInstructions: []*Instruction{},
		Registry:             registry,
		Repository:           "",
		Path:                 filePath,
	}
	dockerfile.parse(path.Join(filePath, "Dockerfile"))
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

// Build return full output of dockerfile
func (d *Dockerfile) Build() string {
	buffer := bytes.NewBufferString("")
	for _, i := range d.Instructions {
		buffer.WriteString(i.Command + " " + strings.Join(i.Arguments, " ") + "\n")
	}
	return buffer.String()
}

// HasLocalSources check if dockerfile has local resources
func (d *Dockerfile) HasLocalSources() bool {
	for _, i := range d.Instructions {
		if i.Command == "ADD" || i.Command == "COPY" {
			return true
		}
	}
	return false
}
