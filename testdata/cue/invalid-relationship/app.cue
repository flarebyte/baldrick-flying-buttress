package flyb

source: "cue-invalid-relationship"
name:   "invalid-relationship"

modules: [
  "core",
  "graph",
]

reports: [{
  title:    "Relationship Check"
  filepath: "reports/relationship.md"
  sections: [{
    title: "Main"
    sections: [{
      title: "Links"
      notes: [
        "n.real",
      ]
    }]
  }]
}]

notes: [{
  name:     "n.real"
  title:    "Real Node"
  markdown: "Existing node."
  labels: [
    "graph",
  ]
}]

relationships: [{
  from:   "n.real"
  to:     "n.missing"
  labels: ["depends_on"]
}]
