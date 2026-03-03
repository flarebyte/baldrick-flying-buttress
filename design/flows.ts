import { calls } from './calls';
import { type ComponentCall, type FlowContext, incrContext } from './common';
import { useCases } from './use_cases.ts';

const uc = {
  reportGenerate: useCases['cli.report.generate'].name,
  noteBasicMarkdown: useCases['cli.note.basic-markdown'].name,
  noteFilepathReference: useCases['cli.note.filepath.reference'].name,
  noteCsvEmbed: useCases['cli.note.csv.embed'].name,
  noteCsvFilterColumn: useCases['cli.note.csv.filter-column'].name,
  noteImagePreview: useCases['cli.note.image.preview'].name,
  noteMermaidEmbed: useCases['cli.note.mermaid.embed'].name,
  noteLinkMarkdown: useCases['cli.note.link.markdown'].name,
  configRelationshipsLabeled: useCases['cli.config.relationships.labeled'].name,
  configReportsMultiple: useCases['cli.config.reports.multiple'].name,
  reportList: useCases['cli.report.list'].name,
  exportJsonGraph: useCases['cli.export.json.graph'].name,
  reportSubgraphByLabel: useCases['cli.report.subgraph.by-label'].name,
};

export const cliRoot = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'cli.root',
    title: 'flyb CLI root command',
    note: 'Entry point for report generation, listing, JSON export, and config validation.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
  listReportsAction(incrContext(context));
  generateMarkdownAction(incrContext(context));
  generateJsonAction(incrContext(context));
  validateAction(incrContext(context));
};

export const listReportsAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.list.reports',
    title: 'List configured markdown reports',
    note: 'Enumerate report targets without generating files.',
    level: context.level,
    useCases: [uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const generateMarkdownAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown',
    title: 'Generate markdown reports',
    note: 'Renders one or more markdown outputs from the validated config.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  generateMarkdownSections(incrContext(context));
};

export const generateMarkdownSections = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown.sections',
    title: 'Generate markdown sections',
    note: 'Build H3 sections from note subsets and renderers.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.noteBasicMarkdown],
  };
  calls.push(call);
  generateSingleH3Section(incrContext(context));
};

export const renderPlainSection = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.plain',
    title: 'Render plain section',
    note: 'Render title and markdown body, including markdown links.',
    level: context.level,
    useCases: [uc.noteBasicMarkdown, uc.noteLinkMarkdown],
  };
  calls.push(call);
};

export const renderGraphSection = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.graph',
    title: 'Render section as a graph',
    note: 'Renders graph-focused content from filtered notes and relationships.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.reportSubgraphByLabel],
  };
  calls.push(call);
};
export const renderSectionWithFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file',
    title: 'Render section with referenced file content',
    note: 'Dispatches file rendering by type (CSV, media, code/diagram).',
    level: context.level,
    useCases: [uc.noteFilepathReference],
  };
  calls.push(call);
  renderSectionWithCsvFile(incrContext(context));
  renderSectionWithMedia(incrContext(context));
  renderSectionWithCodeSnippet(incrContext(context));
};

export const renderSectionWithCsvFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.csv',
    title: 'Render section with CSV file',
    note: 'Render as a markdown table or raw CSV code block.',
    level: context.level,
    useCases: [uc.noteFilepathReference, uc.noteCsvEmbed],
  };
  calls.push(call);
  filterCsvFile(incrContext(context));
};

export const renderSectionWithCodeSnippet = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.code',
    title: 'Render section with code or Mermaid snippet',
    note: 'Preserve fenced-block formatting for code and Mermaid content.',
    level: context.level,
    useCases: [uc.noteFilepathReference, uc.noteMermaidEmbed],
  };
  calls.push(call);
};
export const filterCsvFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'file.csv.filter',
    title: 'Filter CSV rows by column',
    note: 'Apply include/exclude filters before rendering CSV output.',
    level: context.level,
    useCases: [uc.noteCsvFilterColumn],
  };
  calls.push(call);
};

export const selectSubGraph = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.select',
    title: 'Extract subgraph using labels',
    note: 'Filter notes and relationships by labels and optional starting node.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.reportSubgraphByLabel],
  };
  calls.push(call);
};

export const renderSectionWithMedia = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.media',
    title: 'Render section with media file',
    note: 'Embed image previews for supported media types.',
    level: context.level,
    useCases: [uc.noteFilepathReference, uc.noteImagePreview],
  };
  calls.push(call);
};

export const generateSingleH3Section = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown.section.h3',
    title: 'Generate a single H3 section',
    note: 'Compose subgraph, plain content, and file-backed content for one section.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.reportSubgraphByLabel],
  };
  calls.push(call);
  selectSubGraph(incrContext(context));
  renderGraphSection(incrContext(context));
  renderPlainSection(incrContext(context));
  renderSectionWithFile(incrContext(context));
};

export const generateJsonAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.json',
    title: 'Generate JSON graph export',
    note: 'Export notes and relationships in machine-readable JSON format.',
    level: context.level,
    useCases: [uc.exportJsonGraph],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const validateAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.validate',
    title: 'Validate the CUE file',
    note: 'Validate configuration structure and constraints without rendering output.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const loadAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'load.app.data',
    title: 'Load CUE application data',
    note: 'Read notes, relationships, and report definitions from config.',
    level: context.level,
    useCases: [
      uc.reportGenerate,
      uc.exportJsonGraph,
      uc.configRelationshipsLabeled,
      uc.configReportsMultiple,
      uc.reportList,
    ],
  };
  calls.push(call);
};

export const validateAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'validate.app.data',
    title: 'Validate CUE application data',
    note: 'Ensure required fields and cross-reference integrity are valid.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.configReportsMultiple],
  };
  calls.push(call);
};
