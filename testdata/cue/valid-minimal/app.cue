package flyb

source: "cue-valid-minimal"
name:   "minimal-app"

modules: [
  "core",
]

reports: [{
  title:    "Overview"
  filepath: "reports/overview.md"
  sections: [{
    title: "Main"
    sections: [{
      title: "Summary"
      notes: [
        "n.root",
      ]
    }]
  }]
}]

notes: [{
  name:     "n.root"
  title:    "Root Note"
  markdown: "Minimal note body."
  labels: [
    "core",
  ]
}]

relationships: []
