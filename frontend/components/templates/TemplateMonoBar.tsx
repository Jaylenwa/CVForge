import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateMonoBar: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#050505ff');
  const { lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const customPairs = React.useMemo(() => {
    return normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo));
  }, [personal?.CustomInfo]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const headerParts = React.useMemo(() => {
    const parts: string[] = [
      normalizedAge,
      ...customPairs.map(ci => String(ci.value || ci.label || '').trim()).filter(Boolean),
      String(personal?.Degree || '').trim(),
      String(personal?.City || '').trim(),
      String(personal?.Money || '').trim(),
      String(personal?.JoinTime || '').trim(),
      String(personal?.Phone || '').trim(),
      String(personal?.Email || '').trim(),
      String(personal?.Gender || '').trim(),
    ].filter(Boolean);
    return parts;
  }, [normalizedAge, personal?.Degree, personal?.City, personal?.Money, personal?.JoinTime, personal?.Phone, personal?.Email, personal?.Gender, customPairs]);

  const SectionTitle: React.FC<{ section: any }> = ({ section }) => (
    <div className="flex items-center gap-3 mb-4">
      <div className="w-1 h-5" style={{ backgroundColor: color }} />
      <h3 className="text-lg font-bold text-gray-900">{getSectionTitle(section)}</h3>
    </div>
  );

  const renderTightList = (items: any[]) => {
    return (
      <div className={listTightClass}>
        {items.map((item: any) => (
          <div key={item.id}>
            {item.title ? <div className="text-sm font-semibold text-gray-900">{item.title}</div> : null}
            {item.subtitle ? <div className="text-sm text-gray-700 mt-0.5">{item.subtitle}</div> : null}
            {item.description ? <RichText html={item.description} className="text-gray-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} /> : null}
          </div>
        ))}
      </div>
    );
  };

  const renderEducation = (items: any[]) => {
    return (
      <div className={listMediumClass}>
        {items.map((item: any) => {
          const right = (item.major || item.degree)
            ? `${String(item.major || '').trim()}${item.major && item.degree ? `（${String(item.degree).trim()}）` : String(item.degree || '').trim()}`
            : String(item.subtitle || '').trim();
          return (
            <div key={item.id}>
              <div className="grid grid-cols-12 gap-3 items-baseline">
                <div className="col-span-3 text-sm text-gray-600 whitespace-nowrap">{formatRange(item)}</div>
                <div className="col-span-6 text-sm font-semibold text-gray-900 text-center">{item.title}</div>
                <div className="col-span-3 text-sm text-gray-700 text-right">{right}</div>
              </div>
              {item.description ? (
                <div className="mt-2 w-full">
                  <RichText html={item.description} className="text-gray-700 w-full" fontSize={styles.fontSize} lineHeight={lineHeight} />
                </div>
              ) : null}
            </div>
          );
        })}
      </div>
    );
  };

  const renderThreeColumnEntries = (items: any[]) => {
    return (
      <div className={listMediumClass}>
        {items.map((item: any) => (
          <div key={item.id}>
            <div className="grid grid-cols-12 gap-3 items-baseline">
              <div className="col-span-3 text-sm text-gray-600 whitespace-nowrap">{formatRange(item)}</div>
              <div className="col-span-6 text-sm font-semibold text-gray-900 text-center">{item.title}</div>
              <div className="col-span-3 text-sm text-gray-700 text-right">{item.subtitle}</div>
            </div>
            {item.description ? (
              <div className="mt-2 w-full">
                <RichText html={item.description} className="text-gray-700 w-full" fontSize={styles.fontSize} lineHeight={lineHeight} />
              </div>
            ) : null}
          </div>
        ))}
      </div>
    );
  };

  const renderSection = (section: any) => {
    const items = getOrderedItems(section.items || []);
    if (section.type === ResumeSectionType.Exam) {
      return <ExamSection section={section} color={color} t={t} />;
    }
    if (section.type === ResumeSectionType.Education) {
      return renderEducation(items);
    }
    if (section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Projects || section.type === ResumeSectionType.Internships) {
      return renderThreeColumnEntries(items);
    }
    return renderTightList(items);
  };

  const sectionsOrdered = React.useMemo(() => {
    return getOrderedVisibleSections(data.sections || []);
  }, [data.sections]);

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <header className={`flex items-start gap-8 ${headerSpaceClass}`}>
        <div className="flex-1 min-w-0">
          <h1 className="text-3xl font-bold text-gray-900 truncate">{personal?.FullName}</h1>
          {personal?.Job ? (
            <div className="mt-2 text-base text-gray-700">
              {t('editor.fields.jobApplication')}：{personal.Job}
            </div>
          ) : null}
          {headerParts.length > 0 ? (
            <div className="mt-3 text-sm text-gray-700 flex flex-wrap items-center leading-10">
              {headerParts.map((v, idx) => (
                <React.Fragment key={`${v}-${idx}`}>
                  {idx > 0 ? <span className="mx-2 text-gray-400">|</span> : null}
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
      </header>

      <div className={contentGapClass}>
        {sectionsOrdered.map((section: any) => (
          <section key={section.id}>
            <SectionTitle section={section} />
            {renderSection(section)}
          </section>
        ))}
      </div>
    </div>
  );
};

export default TemplateMonoBar;
