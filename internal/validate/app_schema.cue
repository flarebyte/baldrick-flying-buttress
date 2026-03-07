#NonEmpty: string & !=""

#ReportSection: {
	title: #NonEmpty
	description?: string | null
	arguments?: [...string] | null
	notes?: [...string] | null
	sections?: [...#ReportSection] | null
}

#Report: {
	title: #NonEmpty
	filepath: #NonEmpty
	description?: string | null
	sections: [...#ReportSection]
}

#Note: {
	name: #NonEmpty
	title: #NonEmpty
	markdown?: string | null
	filepath?: string | null
	arguments?: [...string] | null
	labels?: [...string] | null
}

#Relationship: {
	from: #NonEmpty
	to:   #NonEmpty
	label?: string | null
	labels?: [...string] | null
}

#RawApp: {
	source: #NonEmpty
	reports: [...#Report]
	notes: [...#Note]
	relationships: [...#Relationship]

	name?: string | null
	modules?: [...string] | null
	argumentRegistry?: _ | null
	graphIntegrityPolicy?: _ | null
}
