package docker

import (
	"fmt"
	"path"
)

func (dockerfile *Dockerfile) BuildContext() {
	if !dockerfile.HasLocalSources() {
		return
	}

	var workDir = ""
	for _, instruction := range dockerfile.Instructions {
		if instruction.Command == "WORKDIR" {
			workDir = instruction.Arguments[0]
		}
		if instruction.Command == "ADD" || instruction.Command == "COPY" {
			localPath, remotePath := instruction.Arguments[0], instruction.Arguments[1]
			if path.IsAbs(localPath) {
				fmt.Println("absolute path")
			} else {
				if workDir != "" {
					localPath = path.Join(workDir, localPath)
				}
				fmt.Println("relative path")
				fmt.Println(localPath)
				fmt.Println(remotePath)
			}
		}
	}
}
