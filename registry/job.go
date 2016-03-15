package registry

type Job struct {
	ID          	string   		`json:"id"`		
	Name			string			`json:"name"`
	Environments    	[]*Container    `json: environments`
}


