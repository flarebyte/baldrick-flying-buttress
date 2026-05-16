package ghflarebyte

project: {
  org:  "flarebyte"
  repo: "baldrick-flying-buttress"
}

sync: {
  mode: "push"
}

repository: {
  description:   "CLI that turns structured design notes and relationships into living architecture documentation"
  defaultBranch: "main"
  homepage:      "https://github.com/flarebyte/baldrick-flying-buttress"
  visibility:    "public"
  template:      false
  topics: [
    "go",
    "cli",
    "architecture",
    "documentation",
    "graph",
    "cue",
    "flarebyte",
  ]
  labels: [
    {
      name:        "bug"
      color:       "d73a4a"
      description: "Something is not working"
    },
    {
      name:        "enhancement"
      color:       "a2eeef"
      description: "New feature or request"
    },
    {
      name:        "documentation"
      color:       "0075ca"
      description: "Documentation improvements"
    },
  ]
  features: {
    issues:                       true
    wiki:                         false
    projects:                     false
    discussions:                  false
    autoMerge:                    true
    mergeCommit:                  false
    rebaseMerge:                  false
    squashMerge:                  true
    squashMergeCommitMessage:     "pr-title"
    deleteBranchOnMerge:          true
    allowForking:                 false
    allowUpdateBranch:            false
    advancedSecurity:             true
    secretScanning:               true
    secretScanningPushProtection: true
  }
}

build: {
  language:             "go"
  mode:                 "binary"
  outputDir:            "build"
  checksumFile:         "build/checksums.txt"
  artifactTargetSuffix: true
  targets: [
    "linux-amd64",
    "linux-arm64",
    "darwin-arm64",
  ]
}

release: {
  versionSource:    "main.project.yaml"
  tagPrefix:        "v"
  notesMode:        "generate-notes"
  includeArtifacts: true
  artifactDir:      "build"
  includeChecksums: true
}
