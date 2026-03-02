import { writeFile } from 'node:fs/promises';
import {
  exampleUseCase,
} from './examples.ts';

const EXAMPLES_PATH = 'doc/EXAMPLES.md';

export const generateExamplesReport = async () => {
  const lines: string[] = [];
  lines.push('# Examples (Generated)');
  lines.push('');
  lines.push('## Use case');
  lines.push('');
  lines.push('```json');
  lines.push(exampleUseCase);
  lines.push('```');
  lines.push('');
 

  await writeFile(EXAMPLES_PATH, lines.join('\n'), 'utf8');
};
