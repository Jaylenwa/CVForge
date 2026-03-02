import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateDarkHeaderIcons: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#1f4a5b');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const hexToRgb = React.useCallback((hex: string): { r: number; g: number; b: number } | null => {
    const raw = String(hex || '').trim().replace(/^#/, '');
    const normalized = raw.length === 3 ? raw.split('').map(c => c + c).join('') : raw;
    const rgbPart = normalized.length >= 6 ? normalized.slice(0, 6) : normalized;
    if (!/^[0-9a-fA-F]{6}$/.test(rgbPart)) return null;
    const num = parseInt(rgbPart, 16);
    return { r: (num >> 16) & 255, g: (num >> 8) & 255, b: num & 255 };
  }, []);

  const headerBg = React.useMemo(() => {
    const rgb = hexToRgb(color) || { r: 31, g: 74, b: 91 };
    const r = Math.round(rgb.r * 0.22);
    const g = Math.round(rgb.g * 0.34);
    const b = Math.round(rgb.b * 0.44);
    return `rgb(${r}, ${g}, ${b})`;
  }, [color, hexToRgb]);

  const headerDivider = React.useMemo(() => {
    const rgb = hexToRgb(color) || { r: 31, g: 74, b: 91 };
    return `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, 0.24)`;
  }, [color, hexToRgb]);

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);

  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const SectionTitle: React.FC<{ title: string }> = ({ title }) => (
    <div className="text-xl font-bold text-slate-900">{title}</div>
  );

  const renderEducation = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map(item => {
        const right = (item.major || item.degree)
          ? `${String(item.major || '').trim()}${item.major && item.degree ? `（${String(item.degree).trim()}）` : String(item.degree || '').trim()}`
          : String(item.subtitle || '').trim();
        return (
          <div key={item.id}>
            <div className="grid grid-cols-12 gap-3 items-baseline">
              <div className="col-span-3 text-sm text-slate-600 whitespace-nowrap">{formatRange(item)}</div>
              <div className="col-span-6 text-sm font-semibold text-slate-900 text-center">{item.title}</div>
              <div className="col-span-3 text-sm text-slate-700 text-right">{right}</div>
            </div>
            {item.description ? (
              <div className="mt-2">
                <RichText html={item.description} className="text-slate-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
              </div>
            ) : null}
          </div>
        );
      })}
    </div>
  );

  const renderExperienceLike = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map((item: any) => (
        <div key={item.id}>
          <div className="grid grid-cols-12 gap-3 items-baseline">
            <div className="col-span-3 text-sm text-slate-600 whitespace-nowrap">{formatRange(item)}</div>
            <div className="col-span-6 text-sm font-semibold text-slate-900 text-center">{item.title}</div>
            <div className="col-span-3 text-sm text-slate-700 text-right">{item.subtitle}</div>
          </div>
          {item.description ? (
            <div className="mt-2">
              <RichText html={item.description} className="text-slate-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
            </div>
          ) : null}
        </div>
      ))}
    </div>
  );

  const renderTightList = (items: any[]) => (
    <div className={listTightClass}>
      {items.map((item: any) => (
        <div key={item.id}>
          {item.title ? <div className="text-sm font-semibold text-slate-900">{item.title}</div> : null}
          {item.subtitle ? <div className="text-sm text-slate-700 mt-0.5">{item.subtitle}</div> : null}
          {item.description ? (
            <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
          ) : null}
        </div>
      ))}
    </div>
  );

  const renderSectionBody = (section: any) => {
    if (section.type === ResumeSectionType.SelfEvaluation) {
      const items = getOrderedItems(section.items || []);
      if (!items.length) return null;
      return renderTightList(items);
    }
    if (section.type === ResumeSectionType.Exam) return <ExamSection section={section} color={color} t={t} />;
    const items = getOrderedItems(section.items || []);
    if (section.type === ResumeSectionType.Education) return renderEducation(items);
    if (section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Projects || section.type === ResumeSectionType.Internships) return renderExperienceLike(items);
    return renderTightList(items);
  };

  const infoPairs = React.useMemo(() => {
    const base: Array<{ label: string; value: string }> = [
      { label: t('editor.fields.jobApplication'), value: String(personal?.Job || '').trim() },
      { label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
      { label: t('editor.fields.city'), value: String(personal?.City || '').trim() },
    ].filter(p => p.value);

    const extra: Array<{ label: string; value: string }> = [
      { label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      { label: t('editor.fields.age'), value: String(normalizedAge || '').trim() },
      { label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() },
      { label: t('editor.fields.expectedSalary'), value: String(personal?.Money || '').trim() },
      { label: t('editor.fields.joinTime'), value: String(personal?.JoinTime || '').trim() },
    ].filter(p => p.value);

    return { base, extra };
  }, [t, personal?.Job, personal?.Phone, personal?.Email, personal?.City, personal?.Gender, normalizedAge, personal?.Degree, personal?.Money, personal?.JoinTime]);

  const contentSections = orderedSections;
  const headerInfoMarginTopClassName = spacingMode === 'compact' ? 'mt-4' : spacingMode === 'spacious' ? 'mt-6' : 'mt-5';
  const headerInfoGridClassName =
    spacingMode === 'compact'
      ? 'grid grid-cols-2 gap-x-8 gap-y-1 text-white/85'
      : spacingMode === 'spacious'
        ? 'grid grid-cols-2 gap-x-12 gap-y-3 text-white/85'
        : 'grid grid-cols-2 gap-x-10 gap-y-2 text-white/85';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="w-full border-b" style={{ backgroundColor: headerBg, borderColor: headerDivider }}>
        <div className="px-10 pt-10 pb-8">
          <div className="flex items-start gap-8">
            <div className="flex-shrink-0 mt-1">
              {personal?.AvatarURL ? (
                <img
                  src={personal.AvatarURL}
                  alt={t('a11y.avatarAlt')}
                  className={getAvatarPhotoClassName()}
                  style={{ backgroundColor: '#ffffff' }}
                />
              ) : (
                <div className={getAvatarPlaceholderClassName('bg-white/20')} />
              )}
            </div>

            <div className="min-w-0 flex-1">
              <h1 className="text-3xl font-bold text-white break-words">{personal?.FullName}</h1>

              {(infoPairs.base.length > 0 || infoPairs.extra.length > 0 || customPairs.length > 0) ? (
                <div className={headerInfoMarginTopClassName}>
                  <div className={headerInfoGridClassName}>
                    {[...infoPairs.base, ...infoPairs.extra].map((p, idx) => (
                      <div key={`${p.label}-${idx}`} className="flex gap-2 min-w-0">
                        <div className="text-white/65 whitespace-nowrap">{p.label}:</div>
                        <div className="min-w-0 break-words">{p.value}</div>
                      </div>
                    ))}
                    {customPairs.map((ci, idx) => {
                      const label = String(ci.label || '').trim();
                      const value = String(ci.value || '').trim();
                      if (!label) {
                        return (
                          <div key={`ci-${idx}`} className="col-span-2 min-w-0 break-words">
                            {value}
                          </div>
                        );
                      }
                      return (
                        <div key={`ci-${idx}`} className="flex gap-2 min-w-0">
                          <div className="text-white/65 whitespace-nowrap">{label}:</div>
                          <div className="min-w-0 break-words">{value}</div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              ) : null}
            </div>
          </div>
        </div>
      </div>

      <div className="px-10 py-8">
        <div className={contentGapClass}>
          {contentSections.map(section => (
            <section key={section.id}>
              <SectionTitle title={getSectionTitle(section)} />
              <div className="mt-4">{renderSectionBody(section)}</div>
            </section>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TemplateDarkHeaderIcons;
