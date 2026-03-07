package flyb

source: "microservice-architecture"
name:   "microservice-architecture"
modules: ["core", "platform"]

reports: [{
	title:       "Microservice Architecture"
	filepath:    "out/microservice.md"
	description: "Service topology and graph view."
	sections: [{
		title:       "Topology"
		description: "Graph-oriented architecture sections."
		sections: [
			{
				title:       "Mermaid Service Graph"
				description: "Graph rendered with Mermaid."
				arguments: [
					"graph-subject-label=service",
					"graph-edge-label=calls",
					"graph-start-node=svc.api",
					"graph-renderer=mermaid",
				]
			},
			{
				title:       "Dependency Tree"
				description: "Tree-style text rendering."
				arguments: [
					"graph-subject-label=service",
					"graph-edge-label=depends_on",
					"graph-start-node=svc.api",
				]
			},
		]
	}]
}]

notes: [
	{
		name:     "svc.api"
		title:    "API Service"
		markdown: "Receives external traffic."
		labels: ["service", "api"]
	},
	{
		name:     "svc.auth"
		title:    "Auth Service"
		markdown: "Handles identity and tokens."
		labels: ["service", "auth"]
	},
	{
		name:     "svc.user"
		title:    "User Service"
		markdown: "Manages user profiles."
		labels: ["service", "user"]
	},
	{
		name:     "svc.orders"
		title:    "Orders Service"
		markdown: "Processes orders."
		labels: ["service", "orders"]
	},
	{
		name:     "infra.db"
		title:    "Shared Database"
		markdown: "Primary data store."
		labels: ["service", "storage"]
	},
]

relationships: [
	{
		from:   "svc.api"
		to:     "svc.auth"
		label:  "calls"
		labels: ["calls"]
	},
	{
		from:   "svc.api"
		to:     "svc.user"
		label:  "calls"
		labels: ["calls"]
	},
	{
		from:   "svc.api"
		to:     "svc.orders"
		label:  "calls"
		labels: ["calls"]
	},
	{
		from:   "svc.auth"
		to:     "infra.db"
		label:  "depends_on"
		labels: ["depends_on"]
	},
	{
		from:   "svc.user"
		to:     "infra.db"
		label:  "depends_on"
		labels: ["depends_on"]
	},
	{
		from:   "svc.orders"
		to:     "infra.db"
		label:  "depends_on"
		labels: ["depends_on"]
	},
]

argumentRegistry: {
	version: "1"
	arguments: [
		{
			name:      "graph-subject-label"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:      "graph-edge-label"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:      "graph-start-node"
			valueType: "string"
			scopes: ["h3-section"]
		},
		{
			name:          "graph-renderer"
			valueType:     "enum"
			scopes: ["h3-section", "note"]
			allowedValues: ["markdown-text", "mermaid"]
			defaultValue:  "markdown-text"
		},
	]
}
