import React from 'react';
import { Award, BookOpen, Briefcase, GraduationCap, Heart, Layers, Mail, Phone, User, Wrench } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBlueTabLine: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#4b6076');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

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

  const intentParts = React.useMemo(() => {
    const parts = [
      String(personal?.Job || '').trim(),
      String(personal?.City || '').trim(),
      String(personal?.Money || '').trim(),
      String(personal?.JoinTime || '').trim(),
    ].filter(Boolean);
    return parts;
  }, [personal?.City, personal?.Job, personal?.JoinTime, personal?.Money]);

  const findWorkYearsPair = React.useMemo(() => {
    const targetLabels = ['工作年限', '工作经验', '经验', '年限'];
    return (
      customPairs.find((p) => {
        const label = String(p.label || '').trim();
        return label && targetLabels.some(k => label.includes(k));
      }) || null
    );
  }, [customPairs]);

  const infoCells = React.useMemo(() => {
    const workYearsValue = findWorkYearsPair ? String(findWorkYearsPair.value || '').trim() : '';
    const workYearsLabel = findWorkYearsPair ? String(findWorkYearsPair.label || '').trim() : '';
    const workYears =
      workYearsLabel && workYearsValue
        ? { icon: <Briefcase size={16} />, label: workYearsLabel, value: workYearsValue }
        : { icon: <GraduationCap size={16} />, label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() };

    const base: Array<{ icon: React.ReactNode; label: string; value: string }> = [
      { icon: <User size={16} />, label: t('editor.fields.age'), value: normalizedAge },
      { icon: <User size={16} />, label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      workYears,
      { icon: <Phone size={16} />, label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { icon: <Mail size={16} />, label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ].filter(x => x.value);

    const cityValue = String(personal?.City || '').trim();
    const excludedLabels = new Set<string>([String(findWorkYearsPair?.label || '').trim(), ...base.map(x => x.label)]);
    const custom = customPairs
      .map(p => ({ icon: <User size={16} />, label: String(p.label || '').trim(), value: String(p.value || '').trim() }))
      .filter(p => {
        if (!p.label || !p.value) return false;
        if (excludedLabels.has(p.label)) return false;
        if (cityValue && p.value === cityValue && (p.label.includes('城市') || p.label === t('editor.fields.city'))) return false;
        return true;
      });

    return [...base, ...custom];
  }, [customPairs, findWorkYearsPair, normalizedAge, personal?.City, personal?.Degree, personal?.Email, personal?.Gender, personal?.Phone, t]);

  const infoRows = React.useMemo(() => {
    const rows: Array<Array<{ icon: React.ReactNode; label: string; value: string }>> = [];
    for (let i = 0; i < infoCells.length; i += 2) {
      rows.push(infoCells.slice(i, i + 2));
    }
    return rows;
  }, [infoCells]);

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Education:
        return <GraduationCap size={16} />;
      case ResumeSectionType.Experience:
      case ResumeSectionType.Internships:
        return <Briefcase size={16} />;
      case ResumeSectionType.Projects:
        return <Layers size={16} />;
      case ResumeSectionType.Skills:
        return <Wrench size={16} />;
      case ResumeSectionType.Awards:
        return <Award size={16} />;
      case ResumeSectionType.Interests:
        return <Heart size={16} />;
      case ResumeSectionType.Exam:
        return <BookOpen size={16} />;
      case ResumeSectionType.SelfEvaluation:
      default:
        return <User size={16} />;
    }
  };

  const TabLineHeader: React.FC<{ title: string; iconType?: ResumeSectionType }> = ({ title, iconType }) => {
    const icon = iconType ? <SectionIcon type={iconType} /> : <User size={16} />;
    return (
      <div className="relative">
        <div className="absolute left-0 right-0 top-1/2 -translate-y-1/2 h-[5px]" style={{ backgroundColor: color }} />
        <div className="relative inline-flex items-center gap-2 pl-3 pr-5 py-1 text-white" style={{ backgroundColor: color }}>
          <div className="w-6 h-6 rounded-full bg-white/15 flex items-center justify-center">{icon}</div>
          <div className="text-base font-bold tracking-wide whitespace-nowrap">{title}</div>
          <div className="absolute left-8 -bottom-[6px] w-3 h-3" style={{ backgroundColor: color, transform: 'rotate(45deg)' }} />
        </div>
      </div>
    );
  };

  const renderThreeColList = (items: any[], rightResolver: (item: any) => string) => (
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
              {rightResolver(item)}
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
          {item.description ? (
            <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
          ) : null}
        </div>
      ))}
    </div>
  );

  const renderSectionBody = (section: any) => {
    const items = getOrderedItems(section.items || []);
    if (!items.length && section.type !== ResumeSectionType.Exam) return null;

    if (section.type === ResumeSectionType.Exam) return <ExamSection section={section} color={color} t={t} />;

    if (section.type === ResumeSectionType.Education) {
      return renderThreeColList(items, (item) => {
        const major = String(item.major || '').trim();
        const degree = String(item.degree || '').trim();
        if (major && degree) return `${major}（${degree}）`;
        return major || degree || String(item.subtitle || '').trim();
      });
    }

    if ([ResumeSectionType.Experience, ResumeSectionType.Projects, ResumeSectionType.Internships].includes(section.type)) {
      return renderThreeColList(items, (item) => String(item.subtitle || '').trim());
    }

    return renderTightList(items);
  };

  const headerName = String(personal?.FullName || '').trim() || ((data.language || 'zh') === 'zh' ? '简历' : 'Resume');

  const renderable = React.useMemo(() => {
    return orderedSections.filter((section) => {
      if (!section?.isVisible) return false;
      if (section.type === ResumeSectionType.Summary) return false;
      if (section.type === ResumeSectionType.Exam) return true;
      const items = getOrderedItems(section.items || []);
      return items.length > 0;
    });
  }, [orderedSections]);
  const headerIntentClassName =
    spacingMode === 'compact'
      ? 'mt-2 text-slate-700 flex flex-wrap items-center'
      : spacingMode === 'spacious'
        ? 'mt-4 text-slate-700 flex flex-wrap items-center'
        : 'mt-3 text-slate-700 flex flex-wrap items-center';
  const headerIntentSeparatorClassName = spacingMode === 'compact' ? 'mx-2' : spacingMode === 'spacious' ? 'mx-4' : 'mx-3';
  const headerInfoBlockClassName = spacingMode === 'compact' ? 'mt-2 border-y border-slate-200 text-slate-700' : spacingMode === 'spacious' ? 'mt-4 border-y border-slate-200 text-slate-700' : 'mt-3 border-y border-slate-200 text-slate-700';
  const headerInfoCellPaddingClassName = spacingMode === 'compact' ? 'py-1.5 px-2' : spacingMode === 'spacious' ? 'py-3 px-2' : 'py-2 px-2';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="relative h-[18px]">
        <div className="absolute left-0 right-0 top-0 h-[6px]" style={{ backgroundColor: color }} />
        <div className="absolute left-1/2 -translate-x-1/2 top-0 w-[180px] h-[6px]" style={{ backgroundColor: color }} />
        <div className="absolute left-1/2 -translate-x-1/2 top-[6px] w-0 h-0 border-l-[26px] border-r-[26px] border-t-[10px] border-l-transparent border-r-transparent" style={{ borderTopColor: color }} />
      </div>

      <div className="px-10 pt-6">
        <div className="flex items-start gap-10">
          <div className="flex-1 min-w-0">
            <div className="text-4xl font-bold tracking-wide" style={{ color }}>{headerName}</div>

            {intentParts.length > 0 ? (
              <div className={headerIntentClassName}>
                <span className="text-slate-600 whitespace-nowrap">{(data.language || 'zh') === 'zh' ? '求职意向' : t('editor.fields.jobApplication')}：</span>
                {intentParts.map((p, idx) => (
                  <React.Fragment key={`${p}-${idx}`}>
                    <span className="ml-2 font-semibold whitespace-nowrap">{p}</span>
                    {idx < intentParts.length - 1 ? <span className={`${headerIntentSeparatorClassName} text-slate-400`}>|</span> : null}
                  </React.Fragment>
                ))}
              </div>
            ) : null}

            {infoRows.length > 0 ? (
              <div className={headerInfoBlockClassName}>
                {infoRows.map((row, rowIdx) => (
                  <div key={`row-${rowIdx}`} className={`grid grid-cols-2 gap-x-6 ${rowIdx > 0 ? 'border-t border-slate-200' : ''}`}>
                    {row.map((cell, cellIdx) => (
                      <div key={`${cell.label}-${cellIdx}`} className={`flex items-center gap-3 ${headerInfoCellPaddingClassName} min-w-0`}>
                        <div className="w-7 h-7 rounded-full bg-slate-100 text-slate-600 flex items-center justify-center flex-shrink-0">
                          {cell.icon}
                        </div>
                        <div className="text-slate-500 whitespace-nowrap">{cell.label}：</div>
                        <div className="min-w-0 break-words font-semibold text-slate-800">{cell.value}</div>
                      </div>
                    ))}
                    {row.length === 1 ? <div /> : null}
                  </div>
                ))}
              </div>
            ) : null}
          </div>

          <div className="flex-shrink-0 mt-10">
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

      <div className={`px-10 pt-7 pb-10 ${contentGapClass}`}>
        {renderable.map((section) => {
          const title = String(getSectionTitle(section) || '').trim();
          const body = renderSectionBody(section);
          if (!title || !body) return null;
          return (
            <section key={section.id}>
              <TabLineHeader title={title} iconType={section.type} />
              <div className="pt-3">{body}</div>
            </section>
          );
        })}
      </div>
    </div>
  );
};

export default TemplateBlueTabLine;
