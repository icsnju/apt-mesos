package docker

import (
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/utils"
)

var (
	TEMPDIR = "./temp"
)

func (dockerfile *Dockerfile) BuildContext() error {

	if !dockerfile.HasLocalSources() {
		return nil
	}

	// Copy all context to a temp directory
	tempContextDir := path.Join(TEMPDIR, "context-"+dockerfile.ID)
	log.Debugf("Build docker context in path: %v", tempContextDir)

	err := utils.CopyDir(dockerfile.Path, tempContextDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempContextDir)

	tarFile := dockerfile.ID + ".tar"
	log.Debugf("Tar docker context to path: %v", path.Join(TEMPDIR, tarFile))
	err = utils.Tar(tempContextDir, path.Join(TEMPDIR, tarFile), false)
	if err != nil {
		return err
	}
	return nil
}
