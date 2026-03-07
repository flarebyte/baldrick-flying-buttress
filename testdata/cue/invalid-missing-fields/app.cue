package flyb

source: "cue-invalid-missing-fields"
name:   "invalid-missing-fields"

modules: [
  "core",
]

reports: [{
  filepath: "reports/invalid.md"
  sections: [{
    title: "Main"
    sections: [{
      title: "Summary"
      notes: [
        "n.missing-name",
      ]
    }]
  }]
}]

notes: [{
  title:    "Missing Name"
  markdown: "This note has no name field."
  labels: [
    "core",
  ]
}]

relationships: []
