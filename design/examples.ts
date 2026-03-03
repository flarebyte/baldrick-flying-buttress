export type Note = {
  name: string;
  title: string;
  description: string;
  labels: string[];
  filepath?: string;
  url?: string;
  arguments?: string[];
};

export type Relationship = {
  from: string;
  to: string;
  labels: string[];
};

export type H3Section = {
  title: string;
  description: string;
  notes?: Note[];
  relationships?: Relationship[];
  arguments?: string[];
};

export type H2Section = {
  title: string;
  description: string;
  sections: H3Section;
};

export type Report = {
  title: string;
  filepath: string;
  sections: H2Section[];
};
const useCase: Note = {
  name: "usecase.io.calls.count",
  title: "Count I/O calls per function and method",
  description: "Returns counts keyed by function and method.",
  labels: ["usecase"],
};

const rootCall: Note = {
  name: "cli.root",
  title: "flyb CLI root command",
  description: "baldrick-flying-buttress will be shorten to flyb",
  labels: ["call"],
};

const relRootCallToUseCase: Relationship = {
  from: rootCall.name,
  to: useCase.name,
  labels: ["satisfies-usecase"],
};

// Helpful when documenting the exact JSON output shape in markdown/docs.
export const exampleUseCase = JSON.stringify(useCase, null, 2);
export const exampleCall = JSON.stringify(rootCall, null, 2);
export const exampleRelationship = JSON.stringify(
  relRootCallToUseCase,
  null,
  2,
);
