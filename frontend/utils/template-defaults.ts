import { ResumeData } from '../types';

export type ResumeTheme = NonNullable<ResumeData['Theme']>;

const TEMPLATE_DEFAULT_THEMES: Record<string, Partial<ResumeTheme>> = {
  TemplateClassic: { Color: '#050505ff' },
  TemplateMintTimeline: { Color: '#14b8a6' },
  TemplateSlate: { Color: '#2563eb' },
  TemplateMonoBar: { Color: '#050505ff' },
  TemplateSidebarLabel: { Color: '#111827' },
  TemplateBlueStripe: { Color: '#1b9d4a' },
  TemplateDarkHeaderIcons: { Color: '#1f4a5b' },
  TemplateWaveTimeline: { Color: '#2f7a75' },
  TemplateBluePhotoColumns: { Color: '#2c80b9' },
  TemplateBlueResumeHeader: { Color: '#2c80b9' },
  TemplateTealWaveTimeline: { Color: '#3b7f8a' },
};

export const getTemplateDefaultTheme = (templateId: string): Partial<ResumeTheme> | null => {
  const key = String(templateId || '');
  const theme = TEMPLATE_DEFAULT_THEMES[key];
  return theme ? { ...theme } : null;
};

export const applyTemplateThemeDefaults = (templateId: string, theme?: ResumeData['Theme']): ResumeData['Theme'] => {
  const defaults = getTemplateDefaultTheme(templateId);
  if (!defaults) return theme;
  return { ...(theme || {}), ...defaults };
};

export const applyTemplateDefaultsToResumeData = (data: ResumeData): ResumeData => {
  const nextTheme = applyTemplateThemeDefaults(data.templateId, data.Theme);
  if (nextTheme === data.Theme) return data;
  return { ...data, Theme: nextTheme };
};
