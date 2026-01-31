import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBlueStripe: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#2563eb');
  const { lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const hexToRgb = React.useCallback((hex: string): { r: number; g: number; b: number } | null => {
    const raw = String(hex || '').trim().replace(/^#/, '');
    const normalized = raw.length === 3 ? raw.split('').map(c => c + c).join('') : raw;
    if (!/^[0-9a-fA-F]{6}$/.test(normalized)) return null;
    const num = parseInt(normalized, 16);
    return { r: (num >> 16) & 255, g: (num >> 8) & 255, b: num & 255 };
  }, []);

  const gradientColor = React.useMemo(() => {
    const rgb = hexToRgb(color) || { r: 37, g: 99, b: 235 };
    return `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, 0.18)`;
  }, [color, hexToRgb]);

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const customPairs = React.useMemo(() => {
    return normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo));
  }, [personal?.CustomInfo]);

  const headerParts = React.useMemo(() => {
    const customValues = customPairs
      .map(p => String(p.value || p.label || '').trim())
      .filter(Boolean);
    const parts: string[] = [
      normalizedAge,
      String(personal?.City || '').trim(),
      String(personal?.Phone || '').trim(),
      String(personal?.Email || '').trim(),
      String(personal?.Degree || '').trim(),
      String(personal?.Gender || '').trim(),
      String(personal?.Money || '').trim(),
      String(personal?.JoinTime || '').trim(),
      ...customValues,
    ].filter(Boolean);
    const seen = new Set<string>();
    return parts.filter((p) => {
      if (seen.has(p)) return false;
      seen.add(p);
      return true;
    });
  }, [customPairs, normalizedAge, personal?.City, personal?.Phone, personal?.Email, personal?.Degree, personal?.Gender, personal?.Money, personal?.JoinTime]);

  const orderedSections = React.useMemo(() => {
    return getOrderedVisibleSections(data.sections || []);
  }, [data.sections]);

  const summaryHtml = React.useMemo(() => {
    const summary = orderedSections.find(s => s.type === ResumeSectionType.Summary);
    const items = getOrderedItems(summary?.items || []);
    const html = items.map(it => String(it.description || '').trim()).filter(Boolean).join('<br/>');
    return html;
  }, [orderedSections]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const HeaderBar: React.FC<{ label: string }> = ({ label }) => (
    <div className="flex items-stretch w-full">
      <div className="w-2" style={{ backgroundColor: color }} />
      <div className="flex-1 bg-slate-100 py-2 px-4 text-base font-bold text-slate-900">{label}</div>
    </div>
  );

  const SectionBar: React.FC<{ section: any }> = ({ section }) => (
    <div className="flex items-stretch w-full">
      <div className="w-2" style={{ backgroundColor: color }} />
      <div className="flex-1 bg-slate-100 py-2 px-4 text-base font-bold text-slate-900">{getSectionTitle(section)}</div>
    </div>
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

  const renderThreeColumnEntries = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map(item => (
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
      {items.map(item => (
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

  const renderSection = (section: any) => {
    if (section.type === ResumeSectionType.Summary) return null;
    if (section.type === ResumeSectionType.Exam) return <ExamSection section={section} color={color} t={t} />;

    const items = getOrderedItems(section.items || []);

    if (section.type === ResumeSectionType.Education) return renderEducation(items);
    if (section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Internships || section.type === ResumeSectionType.Projects) {
      return renderThreeColumnEntries(items);
    }
    return renderTightList(items);
  };

  const basicInfoTitle = (data.language || 'zh') === 'zh' ? '基本信息' : 'Basic Information';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <HeaderBar label={basicInfoTitle} />

      <header className={`relative pt-6 ${headerSpaceClass}`}>
        <div className="flex items-start gap-8">
          <div className="flex-1 min-w-0">
            <h1 className="text-3xl font-bold text-slate-900 truncate">{personal?.FullName}</h1>

            {summaryHtml ? (
              <>
                {personal?.Job ? <div className="mt-2 text-base text-slate-700">{personal.Job}</div> : null}
                <div className="mt-2">
                  <RichText html={summaryHtml} className="text-slate-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
                </div>
              </>
            ) : personal?.Job ? (
              <div className="mt-2 text-base text-slate-700">{personal.Job}</div>
            ) : null}

            {headerParts.length > 0 ? (
              <div className="mt-3 text-sm text-slate-800 flex flex-wrap items-center leading-8">
                {headerParts.map((v, idx) => (
                  <React.Fragment key={`${v}-${idx}`}>
                    {idx > 0 ? <span className="mx-2 text-slate-400">|</span> : null}
                    <span className="whitespace-nowrap">{v}</span>
                  </React.Fragment>
                ))}
              </div>
            ) : null}
          </div>

          <div className="flex-shrink-0">
            {personal?.AvatarURL ? (
              <img
                src={personal.AvatarURL}
                alt={t('a11y.avatarAlt')}
                className={getAvatarPhotoClassName()}
                style={{ backgroundColor: '#ffffff' }}
              />
            ) : (
              <div className={getAvatarPlaceholderClassName()} />
            )}
          </div>
        </div>
      </header>

      <div className={contentGapClass}>
        {orderedSections.map(section => {
          if (section.type === ResumeSectionType.Summary) return null;
          return (
            <section key={section.id}>
              <SectionBar section={section} />
              <div className="pt-4">{renderSection(section)}</div>
            </section>
          );
        })}
      </div>
    </div>
  );
};

export default TemplateBlueStripe;
