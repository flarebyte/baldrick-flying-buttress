import { writeFile } from 'node:fs/promises';
import {
  exampleH2Section,
  exampleH3Section,
  exampleNote,
  exampleRelationship,
  exampleReport,
  exampleUseCase,
} from './examples.ts';

const EXAMPLES_PATH = 'doc/EXAMPLES.md';

export const generateExamplesReport = async () => {
  const lines: string[] = [];
  lines.push('# Examples (Generated)');
  lines.push('');
  lines.push('## Note');
  lines.push('');
  lines.push('```json');
  lines.push(exampleNote);
  lines.push('```');
  lines.push('');

  lines.push('## Use case');
  lines.push('');
  lines.push('```json');
  lines.push(exampleUseCase);
  lines.push('```');
  lines.push('');

  lines.push('## Relationship');
  lines.push('');
  lines.push('```json');
  lines.push(exampleRelationship);
  lines.push('```');
  lines.push('');

  lines.push('## H3 section');
  lines.push('');
  lines.push('```json');
  lines.push(exampleH3Section);
  lines.push('```');
  lines.push('');

  lines.push('## H2 section');
  lines.push('');
  lines.push('```json');
  lines.push(exampleH2Section);
  lines.push('```');
  lines.push('');

  lines.push('## Report');
  lines.push('');
  lines.push('```json');
  lines.push(exampleReport);
  lines.push('```');
  lines.push('');

  await writeFile(EXAMPLES_PATH, lines.join('\n'), 'utf8');
};
