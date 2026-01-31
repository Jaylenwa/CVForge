import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateSlate: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const color = getAccentColor(data, '#0f172a');
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;

  const { lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const customPairs = React.useMemo(() => {
    return normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo));
  }, [personal?.CustomInfo]);

  const basePairs: Array<{ label: string; value: string }> = [
    { label: t('editor.fields.jobApplication'), value: personal?.Job || '' },
    { label: t('editor.fields.phone'), value: personal?.Phone || '' },
    { label: t('editor.fields.email'), value: personal?.Email || '' },
    { label: t('editor.fields.city'), value: personal?.City || '' },
    { label: t('editor.fields.degree'), value: personal?.Degree || '' },
    { label: t('editor.fields.gender'), value: personal?.Gender || '' },
    { label: t('editor.fields.age'), value: personal?.Age || '' },
    { label: t('editor.fields.expectedSalary'), value: personal?.Money || '' },
    { label: t('editor.fields.joinTime'), value: personal?.JoinTime || '' },
  ].filter(p => p.value && String(p.value).trim());

  const sectionsOrdered = React.useMemo(() => {
    return getOrderedVisibleSections(data.sections || []);
  }, [data.sections]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'tilde' });

  const SectionTitle: React.FC<{ section: any }> = ({ section }) => (
    <div className="flex items-center gap-3 mb-4">
      <h3 className="text-lg font-bold tracking-wide" style={{ color }}>
        {getSectionTitle(section)}
      </h3>
      <div className="h-[2px] flex-1 rounded" style={{ backgroundColor: `${color}88` }} />
    </div>
  );

  const renderSection = (section: any) => {
    const items = getOrderedItems(section.items || []);
    if (section.type === ResumeSectionType.Exam) {
      return <ExamSection section={section} color={color} t={t} />;
    }
    if (section.type === ResumeSectionType.Skills || section.type === ResumeSectionType.SelfEvaluation || section.type === ResumeSectionType.Portfolio || section.type === ResumeSectionType.Awards || section.type === ResumeSectionType.Interests) {
      return (
        <div className={listTightClass}>
          {items.map((item: any) => (
            <div key={item.id}>
              {item.title && <div className="text-sm font-semibold text-slate-900">{item.title}</div>}
              {item.subtitle && <div className="text-xs text-slate-600">{item.subtitle}</div>}
              {item.description && (
                <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
              )}
            </div>
          ))}
        </div>
      );
    }
    const timelineLike = section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Projects || section.type === ResumeSectionType.Internships || section.type === ResumeSectionType.Education;
    if (timelineLike) {
      return (
        <div className={listMediumClass}>
          {items.map((item: any) => (
            <div key={item.id} className="relative pl-5">
              <div className="absolute left-0 top-2 w-2 h-2 rotate-45" style={{ backgroundColor: color }} />
              <div className="absolute left-[3px] top-5 bottom-0 w-[1px]" style={{ backgroundColor: `${color}44` }} />
              <div className="flex items-baseline justify-between gap-4">
                <div className="min-w-0">
                  <div className="text-sm font-bold text-slate-900 truncate">{item.title}</div>
                  {section.type === ResumeSectionType.Education ? (
                    (item.major || item.degree) ? (
                      <div className="text-xs text-slate-600">
                        {item.major}{item.major && item.degree ? ' • ' : ''}{item.degree}
                      </div>
                    ) : (item.subtitle ? <div className="text-xs text-slate-600">{item.subtitle}</div> : null)
                  ) : (
                    item.subtitle ? <div className="text-xs text-slate-600">{item.subtitle}</div> : null
                  )}
                </div>
                {(item.timeStart || item.timeEnd || item.today) && (
                  <div className="text-xs text-slate-600 whitespace-nowrap">{formatRange(item)}</div>
                )}
              </div>
              {item.description && (
                <RichText html={item.description} className="text-slate-700 mt-2" fontSize={styles.fontSize} lineHeight={lineHeight} />
              )}
            </div>
          ))}
        </div>
      );
    }
    return (
      <div className={listTightClass}>
        {items.map((item: any) => (
          <div key={item.id}>
            {(item.title || item.subtitle) && (
              <div className="flex items-baseline justify-between gap-3">
                <div className="text-sm font-semibold text-slate-900">{item.title}</div>
                {(item.timeStart || item.timeEnd || item.today) && (
                  <div className="text-xs text-slate-600 whitespace-nowrap">{formatRange(item)}</div>
                )}
              </div>
            )}
            {(item.major || item.degree) && (
              <div className="text-xs text-slate-600 mt-0.5">
                {item.major}{item.major && item.degree ? ' • ' : ''}{item.degree}
              </div>
            )}
            {item.subtitle && section.type !== ResumeSectionType.Education && (
              <div className="text-xs text-slate-600 mt-0.5">{item.subtitle}</div>
            )}
            {item.description && (
              <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
            )}
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="px-10 py-8">
        <header className={headerSpaceClass}>
          <div className="flex items-start gap-6">
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
            <div className="flex-1 min-w-0 pl-8">
              <div className="min-w-0">
                <h1 className="text-3xl font-extrabold tracking-wide text-slate-900 truncate">{personal?.FullName}</h1>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-x-8 gap-y-2 text-sm text-slate-700">
                {basePairs.map((p, idx) => (
                  <div key={`${p.label}-${idx}`} className="flex gap-2 min-w-0">
                    <div className="text-slate-500 whitespace-nowrap">{p.label}:</div>
                    <div className="min-w-0 break-words">{p.value}</div>
                  </div>
                ))}
                {customPairs.map((ci, idx) => {
                  if (!ci.label) {
                    return (
                      <div key={`ci-${idx}`} className="col-span-2 break-words">
                        {ci.value}
                      </div>
                    );
                  }
                  return (
                    <div key={`ci-${idx}`} className="flex gap-2 min-w-0">
                      <div className="text-slate-500 whitespace-nowrap">{ci.label}:</div>
                      <div className="min-w-0 break-words">{ci.value}</div>
                    </div>
                  );
                })}
              </div>
            </div>
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
    </div>
  );
};

export default TemplateSlate;
