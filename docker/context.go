package docker

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/icsnju/apt-mesos/utils"
)

var (
	TEMPDIR = "../temp"
)

func (dockerfile *Dockerfile) BuildContext() error {
	if !dockerfile.HasLocalSources() {
		return nil
	}

	// Copy all context to a temp directory
	contextDir, err := ioutil.TempDir(TEMPDIR, "job_context")
	if err != nil {
		return err
	}
	err = utils.CopyDir(dockerfile.Path, contextDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(contextDir)

	tarFile := dockerfile.ID + ".tar"
	err = utils.Tar(contextDir, path.Join(TEMPDIR, tarFile), false)
	if err != nil {
		return err
	}
	return nil
}
