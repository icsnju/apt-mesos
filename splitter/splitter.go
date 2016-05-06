package splitter

type Splitter interface {
	Split(input string) ([]string, error)
}
