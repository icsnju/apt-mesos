package registry

type Container struct{
	Name			string			`json:"name"`
	Image     		string      	`json:"image"`
	Network			string			`json:"network"`
	PortMappings 	[]*PortMapping	`json:"port_mappings"`
	Volumes			[]*Volume 		`json:"volumes"`
	Instances		int64			`json:"instances"`
}
