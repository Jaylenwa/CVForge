import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateSidebarLabel: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#111827');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const customInfo = normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo));
  const basePairs: Array<{ label: string; value: string }> = [
    { label: t('editor.fields.jobApplication'), value: personal?.Job || '' },
    { label: t('editor.fields.phone'), value: personal?.Phone || '' },
    { label: t('editor.fields.email'), value: personal?.Email || '' },
  ].filter(p => p.value && String(p.value).trim());
  const extraPairs: Array<{ label: string; value: string }> = [
    { label: t('editor.fields.gender'), value: personal?.Gender || '' },
    { label: t('editor.fields.age'), value: personal?.Age || '' },
    { label: t('editor.fields.degree'), value: personal?.Degree || '' },
    { label: t('editor.fields.city'), value: personal?.City || '' },
    { label: t('editor.fields.expectedSalary'), value: personal?.Money || '' },
    { label: t('editor.fields.joinTime'), value: personal?.JoinTime || '' },
  ].filter(p => p.value && String(p.value).trim());
  const hasCustomInfo = customInfo.length > 0;

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const SectionHeader: React.FC<{ title: string }> = ({ title }) => (
    <div>
      <div className="text-lg font-bold text-gray-900">{title}</div>
      <div className="mt-2 relative h-[10px] w-full">
        <div className="absolute left-0 top-0 h-[3px] w-28" style={{ backgroundColor: color }} />
        <div className="absolute left-0 bottom-0 h-px w-full bg-gray-300" />
      </div>
    </div>
  );

  const renderExperienceLike = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map((item: any) => (
        <div key={item.id} className="w-full">
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

  const renderTightSection = (items: any[]) => (
    <div className={listTightClass}>
      {items.map((item: any) => (
        <div key={item.id} className="w-full">
          {item.title ? <div className="text-sm font-semibold text-gray-900">{item.title}</div> : null}
          {item.subtitle ? <div className="text-sm text-gray-700 mt-0.5">{item.subtitle}</div> : null}
          {item.description ? (
            <div className="mt-1 w-full">
              <RichText html={item.description} className="text-gray-700 w-full" fontSize={styles.fontSize} lineHeight={lineHeight} />
            </div>
          ) : null}
        </div>
      ))}
    </div>
  );

  const renderEducation = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map((item: any) => {
        const right = (item.major || item.degree)
          ? `${String(item.major || '').trim()}${item.major && item.degree ? `（${String(item.degree).trim()}）` : String(item.degree || '').trim()}`
          : String(item.subtitle || '').trim();
        return (
          <div key={item.id} className="w-full">
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

  const sectionsOrdered = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);

  const getDisplayedTitle = (section: any) => {
    if ((data.language || 'zh') === 'zh') {
      if (section.type === ResumeSectionType.Summary) return '优势分析';
      if (section.type === ResumeSectionType.Personal) return '基本信息';
    }
    return getSectionTitle(section);
  };

  const renderSectionBody = (section: any) => {
    const items = getOrderedItems(section.items || []);
    if (section.type === ResumeSectionType.Exam) return <ExamSection section={section} color={color} t={t} />;
    if (section.type === ResumeSectionType.Education) return renderEducation(items);
    if (section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Projects || section.type === ResumeSectionType.Internships) return renderExperienceLike(items);
    return renderTightSection(items);
  };
  const remainingSections = sectionsOrdered;

  const headerInfoGridClassName =
    spacingMode === 'compact'
      ? 'mt-4 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-1 text-gray-600'
      : spacingMode === 'spacious'
        ? 'mt-8 grid grid-cols-1 sm:grid-cols-2 gap-x-10 gap-y-3 text-gray-600'
        : 'mt-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-2 text-gray-600';

  const headerLayoutGapClassName = spacingMode === 'compact' ? 'gap-6' : spacingMode === 'spacious' ? 'gap-10' : 'gap-8';
  const headerSpaceClassName = spacingMode === 'compact' ? 'pb-3 mb-3' : spacingMode === 'spacious' ? 'pb-5 mb-5' : 'pb-4 mb-4';

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="p-6">
        <div className="relative border-2 border-gray-400 p-4">
          <div className={headerSpaceClassName}>
            <div className={`flex items-start ${headerLayoutGapClassName}`}>
              <div className="flex-shrink-0 mt-1">
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
              <div className="min-w-0 flex-1">
                <h1 className="text-3xl font-bold text-gray-900">{personal?.FullName}</h1>
                {(basePairs.length > 0 || extraPairs.length > 0 || hasCustomInfo) ? (
                  <div className={headerInfoGridClassName}>
                    {[...basePairs, ...extraPairs].map((p, idx) => (
                      <div key={`${p.label}-${idx}`} className="flex gap-2 min-w-0">
                        <div className="text-gray-500 whitespace-nowrap">{p.label}:</div>
                        <div className="min-w-0 break-words">{p.value}</div>
                      </div>
                    ))}
                    {customInfo.map((ci, idx) => {
                      if (!ci.label) {
                        return (
                          <div key={`ci-${idx}`} className="col-span-1 sm:col-span-2 break-words">
                            {ci.value}
                          </div>
                        );
                      }
                      return (
                        <div key={`ci-${idx}`} className="flex gap-2 min-w-0">
                          <div className="text-gray-500 whitespace-nowrap">{ci.label}:</div>
                          <div className="min-w-0 break-words">{ci.value}</div>
                        </div>
                      );
                    })}
                  </div>
                ) : null}
              </div>
            </div>
          </div>

          <div className={contentGapClass}>
            {remainingSections.map(section => (
              <section key={section.id} className="pt-2">
                <SectionHeader title={getDisplayedTitle(section)} />
                <div className="mt-3">{renderSectionBody(section)}</div>
              </section>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default TemplateSidebarLabel;
