export type Info = {
  name: string;
  title: string;
  description: string;
  labels: string[];
  filepath?: string;
  url?: string;
};

const useCase: Info = {
  name: 'usecase.io.calls.count',
  title: 'Count I/O calls per function and method',
  description: 'Returns counts keyed by function and method.',
  labels: ['usecase'],
};

// Helpful when documenting the exact JSON output shape in markdown/docs.
export const exampleUseCase = JSON.stringify(useCase, null, 2);
