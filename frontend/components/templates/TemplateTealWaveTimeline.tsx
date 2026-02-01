import React from 'react';
import { Award, BookOpen, Briefcase, GraduationCap, Heart, Layers, Mail, MapPin, Phone, User, Wrench } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getHeaderInfoTextClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateTealWaveTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#3b7f8a');
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

  const headerIntentParts = React.useMemo(() => {
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

  const leftHeaderPairs = React.useMemo(() => {
    const pairs: Array<{ icon: React.ReactNode; label: string; value: string }> = [
      { icon: <User size={16} />, label: t('editor.fields.age'), value: normalizedAge },
      { icon: <MapPin size={16} />, label: t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { icon: <Phone size={16} />, label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
    ];
    return pairs.filter(p => p.value);
  }, [normalizedAge, personal?.City, personal?.Phone, t]);

  const rightHeaderPairs = React.useMemo(() => {
    const workYearsValue = findWorkYearsPair ? String(findWorkYearsPair.value || '').trim() : '';
    const workYearsLabel = findWorkYearsPair ? String(findWorkYearsPair.label || '').trim() : '';
    const workYears =
      workYearsLabel && workYearsValue
        ? { icon: <Briefcase size={16} />, label: workYearsLabel, value: workYearsValue }
        : { icon: <GraduationCap size={16} />, label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() };

    const pairs: Array<{ icon: React.ReactNode; label: string; value: string }> = [
      { icon: <User size={16} />, label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      workYears,
      { icon: <Mail size={16} />, label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ];
    return pairs.filter(p => p.value);
  }, [findWorkYearsPair, personal?.Degree, personal?.Email, personal?.Gender, t]);

  const extraHeaderPairs = React.useMemo(() => {
    const excludedLabels = new Set([String(findWorkYearsPair?.label || '').trim()]);
    return customPairs
      .map(p => ({ icon: <User size={16} />, label: String(p.label || '').trim(), value: String(p.value || '').trim() }))
      .filter(p => p.label && p.value && !excludedLabels.has(p.label));
  }, [customPairs, findWorkYearsPair]);

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Education:
        return <GraduationCap size={18} />;
      case ResumeSectionType.Experience:
      case ResumeSectionType.Internships:
        return <Briefcase size={18} />;
      case ResumeSectionType.Projects:
        return <Layers size={18} />;
      case ResumeSectionType.Skills:
        return <Wrench size={18} />;
      case ResumeSectionType.Awards:
        return <Award size={18} />;
      case ResumeSectionType.Interests:
        return <Heart size={18} />;
      case ResumeSectionType.Exam:
        return <BookOpen size={18} />;
      case ResumeSectionType.SelfEvaluation:
      default:
        return <User size={18} />;
    }
  };

  const DiagonalEdge: React.FC = () => (
    <svg viewBox="0 0 1200 120" preserveAspectRatio="none" className="w-full h-[30px] block">
      <polygon points="0,80 1200,0 1200,120 0,120" fill="#ffffff" />
    </svg>
  );

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

  const headerName = String(personal?.FullName || '').trim();

  const renderable = React.useMemo(() => {
    return orderedSections.filter((section) => {
      if (!section?.isVisible) return false;
      if (section.type === ResumeSectionType.Summary) return false;
      if (section.type === ResumeSectionType.Exam) return true;
      const items = getOrderedItems(section.items || []);
      return items.length > 0;
    });
  }, [orderedSections]);

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="w-full" style={{ backgroundColor: color }}>
        <div className="px-10 pt-7 pb-0">
          <div className="flex items-start gap-10">
            <div className="flex-1 min-w-0">
              {headerName ? <h1 className="text-3xl font-bold text-white">{headerName}</h1> : null}

              {headerIntentParts.length > 0 ? (
                <div className={`mt-3 ${getHeaderInfoTextClassName('text-white/90 flex flex-wrap items-center')}`}>
                  <span className="text-white/80 whitespace-nowrap">{(data.language || 'zh') === 'zh' ? '求职意向' : t('editor.fields.jobApplication')}：</span>
                  {headerIntentParts.map((p, idx) => (
                    <React.Fragment key={`${p}-${idx}`}>
                      <span className="ml-2 font-semibold whitespace-nowrap">{p}</span>
                      {idx < headerIntentParts.length - 1 ? <span className="mx-3 text-white/55">|</span> : null}
                    </React.Fragment>
                  ))}
                </div>
              ) : null}

              {leftHeaderPairs.length > 0 || rightHeaderPairs.length > 0 ? (
                <div className={`mt-5 grid grid-cols-2 gap-x-14 ${getHeaderInfoTextClassName('text-white/90')}`}>
                  <div className="space-y-2">
                    {leftHeaderPairs.map((p, idx) => (
                      <div key={`${p.label}-${idx}`} className="flex items-center gap-3 min-w-0">
                        <div className="w-5 h-5 flex items-center justify-center text-white/90 flex-shrink-0">{p.icon}</div>
                        <div className="text-white/75 whitespace-nowrap">{p.label}：</div>
                        <div className="min-w-0 break-words font-semibold">{p.value}</div>
                      </div>
                    ))}
                  </div>
                  <div className="space-y-2">
                    {rightHeaderPairs.map((p, idx) => (
                      <div key={`${p.label}-${idx}`} className="flex items-center gap-3 min-w-0">
                        <div className="w-5 h-5 flex items-center justify-center text-white/90 flex-shrink-0">{p.icon}</div>
                        <div className="text-white/75 whitespace-nowrap">{p.label}：</div>
                        <div className="min-w-0 break-words font-semibold">{p.value}</div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : null}

              {extraHeaderPairs.length > 0 ? (
                <div className={`mt-3 grid grid-cols-2 gap-x-14 gap-y-2 ${getHeaderInfoTextClassName('text-white/90')}`}>
                  {extraHeaderPairs.map((p, idx) => (
                    <div key={`${p.label}-${idx}`} className="flex items-center gap-3 min-w-0">
                      <div className="w-5 h-5 flex items-center justify-center text-white/90 flex-shrink-0">{p.icon}</div>
                      <div className="text-white/75 whitespace-nowrap">{p.label}：</div>
                      <div className="min-w-0 break-words font-semibold">{p.value}</div>
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

        <DiagonalEdge />
      </div>

      <div className="px-8 pt-3 pb-10">
        <div className="relative">
          <div className={contentGapClass}>
            {renderable.map((section) => {
              const title = String(getSectionTitle(section) || '').trim();
              const body = renderSectionBody(section);
              if (!title || !body) return null;
              return (
                <div key={section.id} className="grid grid-cols-[40px_1fr] gap-x-6 items-start">
                  <div className="flex flex-col items-center" style={{ width: 56 }}>
                    <div className="relative z-10 w-9 h-9 rounded-full flex items-center justify-center bg-white border-2" style={{ borderColor: color, color }}>
                      <SectionIcon type={section.type} />
                    </div>
                  </div>

                  <div className="min-w-0">
                    <div className="flex items-center gap-4">
                      <div className="text-lg font-bold text-slate-900 whitespace-nowrap">{title}</div>
                      <div className="flex-1 h-px bg-slate-200" />
                    </div>
                    <div className="mt-3">{body}</div>
                  </div>
                </div>
              );
            })}
          </div>

          <div
            className="absolute top-[20px] bottom-0 w-[2px] pointer-events-none"
            style={{ left: 28, backgroundColor: `${color}AA` }}
            aria-hidden="true"
          />
        </div>
      </div>
    </div>
  );
};

export default TemplateTealWaveTimeline;
