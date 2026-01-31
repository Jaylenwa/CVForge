import React from 'react';
import { Briefcase, GraduationCap, Wrench, Layers, BookOpen, Award, Heart, Image as ImageIcon, User } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBluePhotoColumns: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#2c80b9');
  const { lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);

  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const infoPairs = React.useMemo(() => {
    const base: Array<{ label: string; value: string }> = [
      { label: t('editor.fields.age'), value: normalizedAge },
      { label: t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ].filter(p => p.value);

    const extra: Array<{ label: string; value: string }> = [
      { label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      { label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() },
      { label: t('editor.fields.expectedSalary'), value: String(personal?.Money || '').trim() },
      { label: t('editor.fields.joinTime'), value: String(personal?.JoinTime || '').trim() },
    ].filter(p => p.value);

    return { base, extra };
  }, [normalizedAge, personal?.City, personal?.Phone, personal?.Email, personal?.Gender, personal?.Degree, personal?.Money, personal?.JoinTime, t]);

  const sectionIconSize = 20;

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Education:
        return <GraduationCap size={sectionIconSize} />;
      case ResumeSectionType.Experience:
      case ResumeSectionType.Internships:
        return <Briefcase size={sectionIconSize} />;
      case ResumeSectionType.Projects:
        return <Layers size={sectionIconSize} />;
      case ResumeSectionType.Skills:
        return <Wrench size={sectionIconSize} />;
      case ResumeSectionType.Portfolio:
        return <ImageIcon size={sectionIconSize} />;
      case ResumeSectionType.Awards:
        return <Award size={sectionIconSize} />;
      case ResumeSectionType.Interests:
        return <Heart size={sectionIconSize} />;
      case ResumeSectionType.Exam:
        return <BookOpen size={sectionIconSize} />;
      case ResumeSectionType.SelfEvaluation:
      default:
        return <User size={sectionIconSize} />;
    }
  };

  const SectionHeader: React.FC<{ section: any }> = ({ section }) => {
    const title = String(getSectionTitle(section) || '').trim();
    if (!title) return null;
    return (
      <div className="flex items-center gap-3">
        <div className="w-8 h-8 rounded-full flex items-center justify-center text-white" style={{ backgroundColor: color }}>
          <SectionIcon type={section.type} />
        </div>
        <div className="text-xl font-bold" style={{ color }}>
          {title}
        </div>
        <div className="h-[2px] flex-1" style={{ backgroundColor: `${color}66` }} />
      </div>
    );
  };

  const renderThreeColList = (items: any[], rightResolver: (item: any) => string) => (
    <div className={listMediumClass}>
      {items.map((item) => (
        <div key={item.id}>
          <div className="grid grid-cols-12 gap-4 items-baseline">
            <div className="col-span-3 text-sm text-slate-600 whitespace-nowrap" style={{ fontSize: styles.fontSize }}>
              {formatRange(item)}
            </div>
            <div className="col-span-6 text-sm font-semibold text-slate-900 text-center" style={{ fontSize: styles.fontSize }}>
              {item.title}
            </div>
            <div className="col-span-3 text-sm text-slate-700 text-right" style={{ fontSize: styles.fontSize }}>
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
        if (major && degree) return `${major} ${degree}`;
        return major || degree || String(item.subtitle || '').trim();
      });
    }

    if ([ResumeSectionType.Experience, ResumeSectionType.Projects, ResumeSectionType.Internships].includes(section.type)) {
      return renderThreeColList(items, (item) => String(item.subtitle || '').trim());
    }

    if ([ResumeSectionType.Skills, ResumeSectionType.SelfEvaluation, ResumeSectionType.Portfolio, ResumeSectionType.Awards, ResumeSectionType.Interests].includes(section.type)) {
      return renderTightList(items);
    }

    return renderTightList(items);
  };

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="relative flex overflow-visible h-[90px]" style={{ backgroundColor: color }}>
        <div className="w-[150px] flex-shrink-0 relative">
          <div className="h-full flex items-start justify-center pt-8">
            {personal?.AvatarURL ? (
              <img
                src={personal.AvatarURL}
                alt={t('a11y.avatarAlt')}
                className={getAvatarPhotoClassName('border-4 border-white shadow-md')}
                style={{ backgroundColor: '#ffffff' }}
              />
            ) : (
              <div className={getAvatarPlaceholderClassName('bg-white/25 shadow-md')} />
            )}
          </div>
        </div>

        <div className="flex-1 min-w-0">
          <div className="px-10 py-8">
            <div className="flex flex-wrap items-baseline gap-x-8 gap-y-2">
              <h1 className="text-4xl font-bold tracking-wide text-white">{personal?.FullName}</h1>
              {personal?.Job ? <div className="text-base font-semibold text-white/90">{t('editor.fields.jobApplication')}：{personal.Job}</div> : null}
            </div>
          </div>
        </div>
      </div>

      <div className="flex bg-white">
        <div className="w-[150px] flex-shrink-0" />
        <div className="flex-1 min-w-0 px-10 pt-2.5 pb-6">
          <div className="flex flex-wrap gap-x-10 gap-y-1 text-sm text-slate-800 leading-7">
            {[...infoPairs.base, ...infoPairs.extra].map((p, idx) => (
              <div key={`${p.label}-${idx}`} className="whitespace-nowrap">
                <span className="text-slate-600">{p.label}：</span>
                <span className="text-slate-900">{p.value}</span>
              </div>
            ))}
            {customPairs.map((ci, idx) => {
              const label = String(ci.label || '').trim();
              const value = String(ci.value || '').trim();
              if (!label && !value) return null;
              if (!label) return (
                <div key={`ci-${idx}`} className="min-w-0 break-words">
                  {value}
                </div>
              );
              return (
                <div key={`ci-${idx}`} className="whitespace-nowrap">
                  <span className="text-slate-600">{label}：</span>
                  <span className="text-slate-900">{value}</span>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      <div className={`px-10 pb-10 ${contentGapClass}`}>
        {orderedSections.map((section) => {
          if (!section?.isVisible) return null;
          if (section.type === ResumeSectionType.Summary) return null;
          const body = renderSectionBody(section);
          if (!body) return null;
          return (
            <section key={section.id}>
              <SectionHeader section={section} />
              <div className="mt-5">{body}</div>
            </section>
          );
        })}
      </div>
    </div>
  );
};

export default TemplateBluePhotoColumns;
