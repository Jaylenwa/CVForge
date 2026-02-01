import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBluePillRail: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#607691');
  const { spacingMode, lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);
  const isZh = (data.language || 'zh') === 'zh';

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);
  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);
  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'tilde', normalizeMonthSeparator: '-' });

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const findWorkYearsPair = React.useMemo(() => {
    const targetLabels = ['工作年限', '工作经验', '经验', '年限'];
    return (
      customPairs.find((p) => {
        const label = String(p.label || '').trim();
        return label && targetLabels.some(k => label.includes(k));
      }) || null
    );
  }, [customPairs]);

  const basicInfoPairs = React.useMemo(() => {
    const workYearsValue = findWorkYearsPair ? String(findWorkYearsPair.value || '').trim() : '';
    const workYearsLabel = findWorkYearsPair ? String(findWorkYearsPair.label || '').trim() : '';
    const workYears = workYearsLabel && workYearsValue ? { label: workYearsLabel, value: workYearsValue } : null;

    const base: Array<{ label: string; value: string }> = [
      { label: isZh ? '姓名' : t('editor.fields.fullName'), value: String(personal?.FullName || '').trim() },
      { label: t('editor.fields.age'), value: normalizedAge },
      { label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      { label: isZh ? '籍贯' : t('editor.fields.city'), value: String(personal?.City || '').trim() },
      workYears ? { label: workYears.label, value: workYears.value } : { label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() },
      { label: isZh ? '求职岗位' : t('editor.fields.jobApplication'), value: String(personal?.Job || '').trim() },
      { label: isZh ? '联系电话' : t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { label: isZh ? '邮箱' : t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ].filter(p => p.value);

    return base;
  }, [findWorkYearsPair, isZh, normalizedAge, personal?.Age, personal?.City, personal?.Degree, personal?.Email, personal?.FullName, personal?.Gender, personal?.Job, personal?.Phone, t]);

  const PillHeader: React.FC<{ title: string }> = ({ title }) => {
    return (
      <div className="flex items-center -ml-10">
        <div className="inline-flex items-center gap-2 px-6 py-1 rounded-full text-white relative" style={{ backgroundColor: color }}>
          <div className="text-base font-bold tracking-wide whitespace-nowrap">{title}</div>
        </div>
        <div className="flex-1 h-0.5 bg-slate-200 ml-0" />
      </div>
    );
  };

  const renderEducation = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map((item) => {
        const major = String(item.major || '').trim();
        const degree = String(item.degree || '').trim();
        const right = major && degree ? `${major}（${degree}）` : (major || degree || String(item.subtitle || '').trim());
        return (
          <div key={item.id}>
            <div className="grid grid-cols-12 gap-4 items-baseline">
              <div className="col-span-3 text-sm font-semibold text-slate-700 whitespace-nowrap" style={{ fontSize: styles.fontSize }}>
                {formatRange(item)}
              </div>
              <div className="col-span-6 text-sm font-semibold text-slate-900 text-center" style={{ fontSize: styles.fontSize }}>
                {item.title}
              </div>
              <div className="col-span-3 text-sm font-semibold text-slate-700 text-right" style={{ fontSize: styles.fontSize }}>
                {right}
              </div>
            </div>
            {item.description ? (
              <div className="mt-2 text-sm text-slate-700 leading-7" style={{ fontSize: styles.fontSize }}>
                <RichText html={item.description} className="text-slate-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
              </div>
            ) : null}
          </div>
        );
      })}
    </div>
  );

  const renderThreeColList = (items: any[]) => (
    <div className={listMediumClass}>
      {items.map((item) => (
        <div key={item.id}>
          <div className="grid grid-cols-12 gap-4 items-baseline">
            <div className="col-span-3 text-sm font-semibold text-slate-700 whitespace-nowrap" style={{ fontSize: styles.fontSize }}>
              {formatRange(item)}
            </div>
            <div className="col-span-6 text-sm font-semibold text-slate-900 text-center" style={{ fontSize: styles.fontSize }}>
              {item.title}
            </div>
            <div className="col-span-3 text-sm font-semibold text-slate-700 text-right" style={{ fontSize: styles.fontSize }}>
              {String(item.subtitle || '').trim()}
            </div>
          </div>
          {item.description ? (
            <div className="mt-2 text-sm text-slate-700 leading-7" style={{ fontSize: styles.fontSize }}>
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
          {item.description ? <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} /> : null}
        </div>
      ))}
    </div>
  );

  const renderSectionBody = (section: any) => {
    const items = getOrderedItems(section.items || []);
    if (!items.length && section.type !== ResumeSectionType.Exam) return null;
    if (section.type === ResumeSectionType.Exam) return <ExamSection section={section} color={color} t={t} />;
    if (section.type === ResumeSectionType.Education) return renderEducation(items);
    if ([ResumeSectionType.Experience, ResumeSectionType.Projects, ResumeSectionType.Internships].includes(section.type)) return renderThreeColList(items);
    return renderTightList(items);
  };

  const renderable = React.useMemo(() => {
    return orderedSections.filter((section) => {
      if (!section?.isVisible) return false;
      if (section.type === ResumeSectionType.Summary) return false;
      if (section.type === ResumeSectionType.Exam) return true;
      const items = getOrderedItems(section.items || []);
      return items.length > 0;
    });
  }, [orderedSections]);

  const basicInfoTitle = (data.language || 'zh') === 'zh' ? '基本信息' : 'Basic Information';
  const headerMain = (data.language || 'zh') === 'zh' ? '个人简历' : 'Resume';
  const headerEn = (data.language || 'zh') === 'zh' ? 'PERSONAL RESUME' : 'PERSONAL RESUME';
  const basicInfoGridClassName =
    spacingMode === 'compact'
      ? 'grid grid-cols-2 gap-x-8 gap-y-2 text-slate-800'
      : spacingMode === 'spacious'
        ? 'grid grid-cols-2 gap-x-12 gap-y-4 text-slate-800'
        : 'grid grid-cols-2 gap-x-10 gap-y-3 text-slate-800';
  const basicInfoPaddingTopClassName = spacingMode === 'compact' ? 'pt-4' : spacingMode === 'spacious' ? 'pt-6' : 'pt-5';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="px-13 pt-4">
        <div className="w-1/2 rounded-r-[48px] px-3 py-3 text-white" style={{ backgroundColor: color }}>
          <div className="flex items-baseline gap-3">
            <div className="pl-5 text-3xl font-bold tracking-wide whitespace-nowrap">{headerMain}</div>
            <div className="text-white/75 font-bold">|</div>
            <div className="text-lg font-semibold tracking-wide whitespace-nowrap">{headerEn}</div>
          </div>
        </div>
      </div>

      <div className={`px-10 pt-6 pb-10 ${headerSpaceClass}`}>
        <div className="relative pl-7">
          <div className="absolute left-3 top-0 bottom-0 w-px bg-slate-200" aria-hidden="true" />

          <section>
            <PillHeader title={basicInfoTitle} />
            <div className={basicInfoPaddingTopClassName}>
              <div className="flex items-start gap-8">
                <div className="flex-1 min-w-0">
                  <div className={basicInfoGridClassName}>
                    {basicInfoPairs.map((p, idx) => (
                      <div key={`${p.label}-${idx}`} className="flex items-baseline gap-3 min-w-0">
                        <div className="w-20 text-right text-slate-500 whitespace-nowrap">{p.label}：</div>
                        <div className={`min-w-0 ${String(p.value || '').includes('@') ? 'break-all' : 'break-words'} font-semibold text-slate-900`}>{p.value}</div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="flex-shrink-0 mt-2">
                  {personal?.AvatarURL ? (
                    <img
                      src={personal.AvatarURL}
                      alt={t('a11y.avatarAlt')}
                      className={getAvatarPhotoClassName('shadow-none')}
                      style={{ backgroundColor: '#ffffff' }}
                    />
                  ) : (
                    <div className={getAvatarPlaceholderClassName()} />
                  )}
                </div>
              </div>
            </div>
          </section>

          <div className={`mt-7 ${contentGapClass}`}>
            {renderable.map((section) => {
              const title = String(getSectionTitle(section) || '').trim();
              const body = renderSectionBody(section);
              if (!title || !body) return null;
              return (
                <section key={section.id}>
                  <PillHeader title={title} />
                  <div className="pt-4">{body}</div>
                </section>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export default TemplateBluePillRail;
