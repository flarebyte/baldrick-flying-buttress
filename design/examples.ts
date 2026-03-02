export type Note = {
  name: string;
  title: string;
  description: string;
  labels: string[];
  filepath?: string;
  url?: string;
};

export type Relationship = {
  from: string;
  to: string;
  otherTo: string[];
  labels: string[];
}

const useCase: Note = {
  name: 'usecase.io.calls.count',
  title: 'Count I/O calls per function and method',
  description: 'Returns counts keyed by function and method.',
  labels: ['usecase'],
};

 const call: Note = {
    name: 'cli.root',
    title: 'flyb CLI root command',
    description: 'baldrick-flying-buttress will be shorten to flyb',
    labels: ['call']
  };

// Helpful when documenting the exact JSON output shape in markdown/docs.
export const exampleUseCase = JSON.stringify(useCase, null, 2);
export const exampleCall = JSON.stringify(call, null, 2);
