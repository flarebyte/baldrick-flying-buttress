package app

type RawApp struct {
	Source string
}

type Report struct {
	ID    string
	Title string
}

type ValidatedApp struct {
	Name    string
	Modules []string
	Reports []Report
}
