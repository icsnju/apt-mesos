package registry

type Job struct {
	Environment    []*Container    `json: environments`
}


