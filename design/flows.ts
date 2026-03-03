import { calls } from "./calls";
import { incrContext, type ComponentCall, type FlowContext } from "./common";
import { useCases } from "./use_cases.ts";

const primaryUseCase = useCases["cli.report.generate"].name;

export const cliRoot = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "cli.root",
    title: "flyb CLI root command",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
  generateMarkdownAction(incrContext(context));
  generateJsonAction(incrContext(context));
  validateAction(incrContext(context));
};

export const generateMarkdownAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "action.generate.markdown",
    title: "Generate the markdown reports",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  generateMarkdownSections(incrContext(context));
};

export const generateMarkdownSections = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "action.generate.markdown.sections",
    title: "Generate the markdown sections",
    note: "",
    level: context.level,
    useCases: [],
  };
  calls.push(call);
  generateSingleh3Section(incrContext(context));
};

export const generateSingleh3Section = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "action.generate.markdown.section.h3",
    title: "Generate the H3 section",
    note: "",
    level: context.level,
    useCases: [],
  };
  calls.push(call);
};

export const generateJsonAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "action.generate.json",
    title: "Generate as json",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const validateAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "action.validate",
    title: "Validate the CUE file",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const loadAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "load.app.data",
    title: "Load CLUE application data",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
};

export const validateAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: "validate.app.data",
    title: "Validate CLUE application data",
    note: "",
    level: context.level,
    useCases: [primaryUseCase],
  };
  calls.push(call);
};
