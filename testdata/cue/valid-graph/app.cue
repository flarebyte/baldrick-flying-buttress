package flyb

source: "cue-valid-graph"
name:   "graph-app"

modules: [
  "core",
  "graph",
]

reports: [{
  title:    "Graph Report"
  filepath: "reports/graph.md"
  sections: [{
    title: "Topology"
    sections: [{
      title: "Nodes"
      notes: [
        "n.a",
        "n.b",
        "n.c",
        "n.d",
      ]
    }]
  }]
}]

notes: [
  {
    name:     "n.a"
    title:    "Node A"
    markdown: "A"
    labels: ["graph", "entry"]
  },
  {
    name:     "n.b"
    title:    "Node B"
    markdown: "B"
    labels: ["graph"]
  },
  {
    name:     "n.c"
    title:    "Node C"
    markdown: "C"
    labels: ["graph"]
  },
  {
    name:     "n.d"
    title:    "Node D"
    markdown: "D"
    labels: ["graph", "leaf"]
  },
]

relationships: [
  {
    from:   "n.a"
    to:     "n.b"
    labels: ["depends_on"]
  },
  {
    from:   "n.a"
    to:     "n.c"
    labels: ["depends_on"]
  },
  {
    from:   "n.c"
    to:     "n.d"
    labels: ["feeds"]
  },
]
