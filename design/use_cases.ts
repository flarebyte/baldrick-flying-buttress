import type { UseCase } from './common.ts';

// Use cases for parsing a single source file (Go, Dart, TypeScript).
export const useCases: Record<string, UseCase> = {
  sameAsName: {
    name: '',
    title: 'Create notes with title and description',
    note: 'free text description usually markdown'
  
  },
  sameAsName: {
    name: '',
    title: 'Note can referenced a filepath to a file',
    note: 'The file could be embeded in final markdown report by using triple quote'
  },
  sameAsName: {
    name: '',
    title: 'Note to a filepath of type csv',
    note: 'When a file is of type csv, it can be embedded as mardown table or csv'
  },
  sameAsName: {
    name: '',
    title: 'Filter csv by column',
    note: 'A csv table can be reduced to a subset by filtering by a column'
  },
  sameAsName: {
    name: '',
    title: 'file can be an image',
    note: 'A file can be an image that will be shown as preview'
  },
  sameAsName: {
    name: '',
    title: 'file can be mermaid',
    note: 'A file can be an a mermaid content that will be shown with using triple quote'
  },
  sameAsName: {
    name: '',
    title: 'A note can have a url link',
    note: 'The link will be converted to a markdown link'
  },
  sameAsName: {
    name: '',
    title: 'the config should have notes with labelled relationships',
    note: 'The config could be written in CUE to offer greater flexibility'
  },
  sameAsName: {
    name: '',
    title: 'The config will declare multiple markdown reports',
    note: ''
  },
  sameAsName: {
    name: '',
    title: 'the cli produce a list or all markdown reports',
    note: ''
  },
  sameAsName: {
    name: '',
    title: 'Notes and relationships can be exported to JSON',
    note: ''
  },
  sameAsName: {
    name: '',
    title: 'A report can use a sub-graph by filter relationship per label',
    note: ''
  },
};

export const getByName = (expectedName: string) =>
  Object.values(useCases).find(({ name }) => name === expectedName);

export const mustUseCases = new Set([
  ...Object.values(useCases).map(({ name }) => name),
]);

export const useCaseCatalogByName: Record<
  string,
  { name: string; title: string; note?: string }
> = Object.fromEntries(Object.values(useCases).map((u) => [u.name, u]));
