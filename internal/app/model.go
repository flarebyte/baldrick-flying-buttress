package app

type RawApp struct {
	Source string
}

type Report struct {
	ID    string
	Title string
}

type Note struct {
	ID    string
	Label string
}

type Relationship struct {
	FromID string
	ToID   string
	Label  string
}

type ValidatedApp struct {
	Name          string
	Modules       []string
	Reports       []Report
	Notes         []Note
	Relationships []Relationship
}
