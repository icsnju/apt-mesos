import "github.com/icsnju/apt-mesos/registry"

type DRFElement struct {
	DominantResource      string
	DominantResourceShare float64
	Job                   *registry.Job
}
