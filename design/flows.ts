import { calls } from './calls';
import { type ComponentCall, type FlowContext, incrContext } from './common';
import { useCases } from './use_cases.ts';

export const cliRoot = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'cli.root',
    title: 'flyer CLI root command',
    directory: 'cmd/maat',
    note: '',
    level: context.level,
    useCases: [useCases.singleFileAnalysis.name],
  };
  calls.push(call);
  // Register commands under the root.
};

