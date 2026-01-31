import { ResumeItem, ResumeSection, ResumeData } from '../../../types';

export type TemplateStyles = {
  fontFamily?: string;
  fontSize?: number;
  spacingMultiplier?: string | number;
};

export type SpacingMode = 'compact' | 'normal' | 'spacious';

export type SpacingTokens = {
  spacingMode: SpacingMode;
  spacingValue: number;
  lineHeight: number;
  headerSpaceClass: string;
  contentGapClass: string;
  listTightClass: string;
  listMediumClass: string;
};

export const ID_PHOTO_WIDTH_MM = 30;
export const ID_PHOTO_HEIGHT_MM = 42;
export const ID_PHOTO_CLASSNAME = `w-[${ID_PHOTO_WIDTH_MM}mm] h-[${ID_PHOTO_HEIGHT_MM}mm]`;

export const getIdPhotoClassName = (extraClassName?: string) => {
  return extraClassName ? `${ID_PHOTO_CLASSNAME} ${extraClassName}` : ID_PHOTO_CLASSNAME;
};

export const AVATAR_PHOTO_BASE_CLASSNAME = 'rounded-sm object-cover border border-slate-200 shadow-sm';
export const AVATAR_PLACEHOLDER_BASE_CLASSNAME = 'rounded-sm bg-slate-200 border border-slate-200';

export const getAvatarPhotoClassName = (extraClassName?: string) => {
  return getIdPhotoClassName(extraClassName ? `${AVATAR_PHOTO_BASE_CLASSNAME} ${extraClassName}` : AVATAR_PHOTO_BASE_CLASSNAME);
};

export const getAvatarPlaceholderClassName = (extraClassName?: string) => {
  return getIdPhotoClassName(extraClassName ? `${AVATAR_PLACEHOLDER_BASE_CLASSNAME} ${extraClassName}` : AVATAR_PLACEHOLDER_BASE_CLASSNAME);
};

export const getAccentColor = (data: ResumeData, fallback: string) => {
  return data?.Theme?.Color || fallback;
};

export const getSpacingValue = (styles: TemplateStyles): number => {
  const spacingValue = Number.parseFloat(String(styles?.spacingMultiplier ?? '1'));
  return Number.isFinite(spacingValue) ? spacingValue : 1;
};

export const getLineHeight = (styles: TemplateStyles, base: number = 1.5): number => {
  return getSpacingValue(styles) * base;
};

export const getSpacingTokens = (styles: TemplateStyles): SpacingTokens => {
  const spacingValue = getSpacingValue(styles);
  const spacingMode: SpacingMode = spacingValue <= 0.9 ? 'compact' : spacingValue >= 1.15 ? 'spacious' : 'normal';
  const lineHeight = spacingValue * 1.5;

  return {
    spacingMode,
    spacingValue,
    lineHeight,
    headerSpaceClass: 'pb-6 mb-6',
    contentGapClass: spacingMode === 'compact' ? 'space-y-6' : spacingMode === 'spacious' ? 'space-y-10' : 'space-y-8',
    listTightClass: spacingMode === 'compact' ? 'space-y-2' : spacingMode === 'spacious' ? 'space-y-4' : 'space-y-3',
    listMediumClass: spacingMode === 'compact' ? 'space-y-4' : spacingMode === 'spacious' ? 'space-y-6' : 'space-y-5',
  };
};

export const getOrderedVisibleSections = (sections: ResumeSection[] = []): ResumeSection[] => {
  return (sections || []).filter(s => s.isVisible).slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
};

export const getOrderedItems = <T extends { orderNum?: number }>(items: T[] = []): T[] => {
  return (items || []).slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
};

export type DateRangeSeparatorVariant = 'tilde' | 'dash';

export type FormatDateRangeOptions = {
  separatorVariant?: DateRangeSeparatorVariant;
  normalizeMonthSeparator?: '.' | '-';
};

const normalizeMonth = (raw: string, normalizeMonthSeparator?: '.' | '-') => {
  const s = String(raw || '').trim();
  if (!s) return '';
  if (!normalizeMonthSeparator) return s;
  return normalizeMonthSeparator === '.' ? s.replace('-', '.') : s.replace('.', '-');
};

export const formatDateRange = (
  item: Pick<ResumeItem, 'timeStart' | 'timeEnd' | 'today'> | undefined,
  t: (key: string) => string,
  options: FormatDateRangeOptions = {},
) => {
  if (!item?.timeStart && !item?.timeEnd && !item?.today) return '';

  const separator = options.separatorVariant === 'dash' ? ' - ' : ' ~ ';
  const start = normalizeMonth(item?.timeStart || item?.timeEnd || '', options.normalizeMonthSeparator);
  const end = item?.today ? t('common.toPresent') : normalizeMonth(item?.timeEnd || '', options.normalizeMonthSeparator);

  if (!start && !end) return '';
  return `${start}${start || end ? separator : ''}${end}`;
};

export type CustomPair = { label?: string; value?: string };
export type NormalizedPair = { label: string; value: string };

export const parseCustomPairs = (raw?: string): CustomPair[] => {
  try {
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed as CustomPair[];
  } catch {
    return [];
  }
};

export const normalizeCustomPairs = (pairs: CustomPair[] = []): NormalizedPair[] => {
  return (pairs || [])
    .map(p => ({ label: String(p?.label || '').trim(), value: String(p?.value || '').trim() }))
    .filter(p => p.label || p.value);
};
