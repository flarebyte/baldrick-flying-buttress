package app

type RawApp struct {
	Source string
}

type ValidatedApp struct {
	Name    string
	Modules []string
}
