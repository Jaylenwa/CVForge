import React from 'react';
import { Award, BookOpen, Briefcase, GraduationCap, Heart, Layers, User, Wrench } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBlueResumeHeader: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#2c80b9');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);

  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);

  const selfEvaluationHtml = React.useMemo(() => {
    const selfEvaluation = orderedSections.find(s => s.type === ResumeSectionType.SelfEvaluation);
    const items = getOrderedItems(selfEvaluation?.items || []);
    const html = items.map(it => String(it.description || '').trim()).filter(Boolean).join('<br/>');
    return html;
  }, [orderedSections]);

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && (data.language || 'zh') === 'zh') return `${raw}岁`;
    return raw;
  }, [personal?.Age, data.language]);

  const basicInfoPairs = React.useMemo(() => {
    const pairs: Array<{ label: string; value: string }> = [
      { label: t('editor.fields.age'), value: normalizedAge },
      { label: t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
      { label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      { label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() },
      { label: t('editor.fields.expectedSalary'), value: String(personal?.Money || '').trim() },
    ].filter(p => p.value);

    const custom = customPairs.map(p => ({ label: String(p.label || '').trim(), value: String(p.value || '').trim() })).filter(p => p.label || p.value);
    return [...pairs, ...custom];
  }, [customPairs, normalizedAge, personal?.City, personal?.Degree, personal?.Email, personal?.Gender, personal?.Money, personal?.Phone, t]);

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

  const BandHeader: React.FC<{ title: string; iconType?: ResumeSectionType }> = ({ title, iconType }) => {
    const icon = iconType ? <SectionIcon type={iconType} /> : <User size={16} />;
    return (
      <div className="w-full bg-slate-100 flex items-stretch">
        <div className="relative flex items-center gap-2 px-4 py-1 text-white" style={{ backgroundColor: color }}>
          <div className="w-5 h-5 flex items-center justify-center">{icon}</div>
          <div className="text-base font-bold tracking-wide whitespace-nowrap">{title}</div>
          <div className="absolute right-[-14px] top-0 h-full w-4" style={{ backgroundColor: color, transform: 'skewX(-25deg)' }} />
          <div className="absolute right-[-26px] top-0 h-full w-3 opacity-55" style={{ backgroundColor: color, transform: 'skewX(-25deg)' }} />
        </div>
        <div className="flex-1" />
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

    return renderTightList(items);
  };

  const title = String(personal?.FullName || '').trim() || t('common.resume');
  const headerPrimaryGridClassName =
    spacingMode === 'compact'
      ? 'mt-3 grid grid-cols-2 gap-x-8 gap-y-1 text-slate-800'
      : spacingMode === 'spacious'
        ? 'mt-5 grid grid-cols-2 gap-x-12 gap-y-3 text-slate-800'
        : 'mt-4 grid grid-cols-2 gap-x-10 gap-y-2 text-slate-800';
  const headerMetaGridClassName = spacingMode === 'compact' ? 'mt-2 grid grid-cols-2 gap-x-8 gap-y-1 text-slate-800' : spacingMode === 'spacious' ? 'mt-4 grid grid-cols-2 gap-x-12 gap-y-3 text-slate-800' : 'mt-3 grid grid-cols-2 gap-x-10 gap-y-2 text-slate-800';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="relative w-full h-[20px]" style={{ backgroundColor: color }}>
        <div className="h-full flex items-center justify-center text-white text-sm font-bold tracking-[0.24em]">
          — RESUME —
        </div>
        <div
          className="absolute right-0 top-0 h-full w-[76px] opacity-90"
          style={{
            backgroundImage: 'repeating-linear-gradient(135deg, rgba(255,255,255,0.75) 0 6px, rgba(255,255,255,0) 6px 12px)',
          }}
        />
      </div>

      <div className="px-10 pt-6">
        <div className="flex items-start gap-10">
          <div className="flex-1 min-w-0">
            <h1 className="text-4xl font-bold text-slate-900">{title}</h1>

            {selfEvaluationHtml ? (
              <div className="mt-4">
                <RichText html={selfEvaluationHtml} className="text-slate-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
              </div>
            ) : null}

            {(personal?.Job || personal?.JoinTime) ? (
              <div className={headerPrimaryGridClassName}>
                {personal?.Job ? (
                  <div className="flex gap-2 min-w-0">
                    <span className="text-slate-600 whitespace-nowrap">{t('editor.fields.jobApplication')}：</span>
                    <span className="text-slate-900 min-w-0 break-words">{String(personal.Job).trim()}</span>
                  </div>
                ) : null}
                {personal?.JoinTime ? (
                  <div className="flex gap-2 min-w-0">
                    <span className="text-slate-600 whitespace-nowrap">{t('editor.fields.joinTime')}：</span>
                    <span className="text-slate-900 min-w-0 break-words">{String(personal.JoinTime).trim()}</span>
                  </div>
                ) : null}
              </div>
            ) : null}

            {basicInfoPairs.length > 0 ? (
              <div className={headerMetaGridClassName}>
                {basicInfoPairs.map((p, idx) => (
                  <div key={`${p.label}-${idx}`} className="flex gap-2 min-w-0">
                    {p.label ? <span className="text-slate-600 whitespace-nowrap">{p.label}：</span> : null}
                    <span className="text-slate-900 min-w-0 break-words">{p.value}</span>
                  </div>
                ))}
              </div>
            ) : null}
          </div>

          <div className="flex-shrink-0">
            <div className="bg-white p-1">
              {personal?.AvatarURL ? (
                <img
                  src={personal.AvatarURL}
                  alt={t('a11y.avatarAlt')}
                  className={getAvatarPhotoClassName('border-0 shadow-none')}
                  style={{ backgroundColor: '#ffffff' }}
                />
              ) : (
                <div className={getAvatarPlaceholderClassName('bg-slate-200 border-0')} />
              )}
            </div>
          </div>
        </div>
      </div>

      <div className={`px-10 pt-8 pb-10 ${contentGapClass}`}>
        {orderedSections.map((section) => {
          if (!section?.isVisible) return null;
          if (section.type === ResumeSectionType.SelfEvaluation) return null;
          const title = String(getSectionTitle(section) || '').trim();
          const body = renderSectionBody(section);
          if (!title || !body) return null;
          return (
            <section key={section.id}>
              <BandHeader title={title} iconType={section.type} />
              <div className="pt-4">{body}</div>
            </section>
          );
        })}
      </div>
    </div>
  );
};

export default TemplateBlueResumeHeader;
