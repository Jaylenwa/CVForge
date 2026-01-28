import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { sanitizeHtml } from '../../utils/resume-helpers';
import { ExamScoreTable } from './ExamScoreTable';

export const TemplateSlateSidebar: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const color = data.Theme?.Color || '#0f172a';
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;

  const spacingValue = Number.parseFloat(String(styles?.spacingMultiplier ?? '1'));
  const spacingMode = spacingValue <= 0.9 ? 'compact' : spacingValue >= 1.15 ? 'spacious' : 'normal';
  const lineHeight = (Number.isFinite(spacingValue) ? spacingValue : 1) * 1.5;
  const contentGapClass = spacingMode === 'compact' ? 'space-y-6' : spacingMode === 'spacious' ? 'space-y-10' : 'space-y-8';
  const headerSpaceClass = 'pb-6 mb-6';
  const listTightClass = spacingMode === 'compact' ? 'space-y-2' : spacingMode === 'spacious' ? 'space-y-4' : 'space-y-3';
  const listMediumClass = spacingMode === 'compact' ? 'space-y-4' : spacingMode === 'spacious' ? 'space-y-6' : 'space-y-5';

  const customPairs = React.useMemo(() => {
    try {
      const raw = personal?.CustomInfo;
      if (!raw) return [] as Array<{ label?: string; value?: string }>;
      const parsed = JSON.parse(raw);
      if (!Array.isArray(parsed)) return [];
      return parsed as Array<{ label?: string; value?: string }>;
    } catch {
      return [] as Array<{ label?: string; value?: string }>;
    }
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
    const visible = (data.sections || []).filter(s => s.isVisible);
    return visible.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [data.sections]);

  const formatRange = (item: any) => {
    if (!item?.timeStart && !item?.timeEnd && !item?.today) return '';
    const start = item?.timeStart || item?.timeEnd || '';
    const end = item?.today ? t('common.toPresent') : (item?.timeEnd || '');
    return `${start}${start || end ? ' ~ ' : ''}${end}`;
  };

  const SectionTitle: React.FC<{ section: any }> = ({ section }) => (
    <div className="flex items-center gap-3 mb-4">
      <h3 className="text-lg font-bold tracking-wide" style={{ color }}>
        {getSectionTitle(section)}
      </h3>
      <div className="h-[2px] flex-1 rounded" style={{ backgroundColor: `${color}22` }} />
    </div>
  );

  const renderExam = (section: any) => {
    const items = (section.items || []).slice().sort((a: any, b: any) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
    const meta = items[0];
    const scores = items.slice(1);
    return (
      <ExamScoreTable
        color={color}
        schoolLabel={t('exam.school')}
        majorLabel={t('exam.major')}
        scoreLabel={(meta?.description && String(meta.description).trim()) ? String(meta.description).trim() : t('exam.scoreLabel')}
        school={meta?.title || ''}
        major={meta?.subtitle || ''}
        items={scores.map((s: any) => ({ subject: s.title || '', score: s.subtitle || '' }))}
      />
    );
  };

  const renderSection = (section: any) => {
    const items = (section.items || []).slice().sort((a: any, b: any) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
    if (section.type === ResumeSectionType.Exam) {
      return renderExam(section);
    }
    if (section.type === ResumeSectionType.Skills || section.type === ResumeSectionType.SelfEvaluation || section.type === ResumeSectionType.Portfolio || section.type === ResumeSectionType.Awards || section.type === ResumeSectionType.Interests) {
      return (
        <div className={listTightClass}>
          {items.map((item: any) => (
            <div key={item.id}>
              {item.title && <div className="text-sm font-semibold text-slate-900">{item.title}</div>}
              {item.subtitle && <div className="text-xs text-slate-600">{item.subtitle}</div>}
              {item.description && (
                <div className="resume-rich-content text-slate-700 text-sm mt-1" style={{ fontSize: styles.fontSize, lineHeight }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
              )}
            </div>
          ))}
        </div>
      );
    }
    const timelineLike = section.type === ResumeSectionType.Experience || section.type === ResumeSectionType.Projects || section.type === ResumeSectionType.Internships;
    if (timelineLike) {
      return (
        <div className={listMediumClass}>
          {items.map((item: any) => (
            <div key={item.id} className="relative pl-5">
              <div className="absolute left-0 top-2 w-2 h-2 rotate-45" style={{ backgroundColor: color }} />
              <div className="absolute left-[3px] top-5 bottom-0 w-[1px]" style={{ backgroundColor: `${color}22` }} />
              <div className="flex items-baseline justify-between gap-4">
                <div className="min-w-0">
                  <div className="text-sm font-bold text-slate-900 truncate">{item.title}</div>
                  {item.subtitle && <div className="text-xs text-slate-600">{item.subtitle}</div>}
                </div>
                {(item.timeStart || item.timeEnd || item.today) && (
                  <div className="text-xs text-slate-600 whitespace-nowrap">{formatRange(item)}</div>
                )}
              </div>
              {item.description && (
                <div className="resume-rich-content text-slate-700 text-sm mt-2" style={{ fontSize: styles.fontSize, lineHeight }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
              <div className="resume-rich-content text-slate-700 text-sm mt-1" style={{ fontSize: styles.fontSize, lineHeight }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
                className="w-[28mm] h-[38mm] rounded-lg object-cover ring-2 ring-white shadow-sm border border-slate-200"
                style={{ backgroundColor: '#ffffff' }}
              />
            ) : (
              <div className="w-[28mm] h-[38mm] rounded-lg bg-slate-200 border border-slate-200" />
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
                  const label = String(ci?.label || '').trim();
                  const value = String(ci?.value || '').trim();
                  if (!label && !value) return null;
                  if (!label) {
                    return (
                      <div key={`ci-${idx}`} className="col-span-2 break-words">
                        {value}
                      </div>
                    );
                  }
                  return (
                    <div key={`ci-${idx}`} className="flex gap-2 min-w-0">
                      <div className="text-slate-500 whitespace-nowrap">{label}:</div>
                      <div className="min-w-0 break-words">{value}</div>
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

export default TemplateSlateSidebar;
