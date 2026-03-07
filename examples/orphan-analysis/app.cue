package flyb

source: "orphan-analysis"
name:   "orphan-analysis"
modules: ["analysis"]

reports: [{
	title:       "Orphan Analysis Report"
	filepath:    "out/orphans.md"
	description: "Contextual orphan checks in a report section."
	sections: [{
		title:       "Ingredients"
		description: "Orphan-oriented section."
		sections: [{
			title:       "Missing Ingredient Links"
			description: "Ingredients without qualifying relationships."
			arguments: [
				"orphan-subject-label=ingredient",
				"orphan-edge-label=uses",
				"orphan-counterpart-label=tool",
				"orphan-direction=either",
			]
		}]
	}]
}]

notes: [
	{
		name:     "ingredient.tomato"
		title:    "Tomato"
		markdown: "Fresh tomato."
		labels: ["ingredient", "fresh"]
	},
	{
		name:     "ingredient.salt"
		title:    "Salt"
		markdown: "Sea salt."
		labels: ["ingredient", "dry"]
	},
	{
		name:     "ingredient.pepper"
		title:    "Pepper"
		markdown: "Black pepper."
		labels: ["ingredient", "dry"]
	},
	{
		name:     "tool.knife"
		title:    "Knife"
		markdown: "Cutting tool."
		labels: ["tool"]
	},
]

relationships: [
	{
		from:   "ingredient.tomato"
		to:     "tool.knife"
		label:  "uses"
		labels: ["uses"]
	},
	{
		from:   "tool.knife"
		to:     "ingredient.salt"
		label:  "uses"
		labels: ["uses"]
	},
]

argumentRegistry: {
	version: "1"
	arguments: [
		{
			name:      "orphan-subject-label"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:      "orphan-edge-label"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:      "orphan-counterpart-label"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:          "orphan-direction"
			valueType:     "enum"
			scopes: ["h3-section"]
			allowedValues: ["in", "out", "either"]
			defaultValue:  "either"
		},
	]
}
