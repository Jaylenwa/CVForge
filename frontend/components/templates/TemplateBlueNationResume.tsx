import React from 'react';
import { Briefcase, CalendarDays, Mail, MapPin, Phone, User, Wallet } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getHeaderInfoTextClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateBlueNationResume: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const color = getAccentColor(data, '#3aa6d8');
  const { lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);
  const isZh = (data.language || 'zh') === 'zh';

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);
  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);
  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'tilde', normalizeMonthSeparator: '-' });

  const normalizedAge = React.useMemo(() => {
    const raw = String(personal?.Age || '').trim();
    if (!raw) return '';
    const isPureNumber = /^[0-9]+$/.test(raw);
    if (isPureNumber && isZh) return `${raw}岁`;
    return raw;
  }, [isZh, personal?.Age]);

  const findWorkYearsPair = React.useMemo(() => {
    const targetLabels = ['工作年限', '工作经验', '经验', '年限'];
    return (
      customPairs.find((p) => {
        const label = String(p.label || '').trim();
        return label && targetLabels.some((k) => label.includes(k));
      }) || null
    );
  }, [customPairs]);

  const intentCells = React.useMemo(() => {
    const pairs: Array<{ icon: React.ReactNode; label: string; value: string }> = [
      { icon: <Briefcase size={16} />, label: t('editor.fields.jobApplication'), value: String(personal?.Job || '').trim() },
      { icon: <MapPin size={16} />, label: t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { icon: <Wallet size={16} />, label: t('editor.fields.expectedSalary'), value: String(personal?.Money || '').trim() },
      { icon: <CalendarDays size={16} />, label: t('editor.fields.joinTime'), value: String(personal?.JoinTime || '').trim() },
    ].filter((x) => x.value);
    return pairs;
  }, [personal?.City, personal?.Job, personal?.JoinTime, personal?.Money, t]);

  const basicCells = React.useMemo(() => {
    const workYearsValue = findWorkYearsPair ? String(findWorkYearsPair.value || '').trim() : '';
    const workYearsLabel = findWorkYearsPair ? String(findWorkYearsPair.label || '').trim() : '';
    const workYears = workYearsLabel && workYearsValue ? { label: workYearsLabel, value: workYearsValue } : null;

    const pairs: Array<{ icon: React.ReactNode; label: string; value: string }> = [
      { icon: <CalendarDays size={16} />, label: t('editor.fields.age'), value: normalizedAge },
      { icon: <Phone size={16} />, label: t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
      { icon: <User size={16} />, label: t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      workYears
        ? { icon: <Briefcase size={16} />, label: workYears.label, value: workYears.value }
        : { icon: <Briefcase size={16} />, label: t('editor.fields.degree'), value: String(personal?.Degree || '').trim() },
      { icon: <Mail size={16} />, label: t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ].filter((x) => x.value);
    return pairs;
  }, [findWorkYearsPair, normalizedAge, personal?.Degree, personal?.Email, personal?.Gender, personal?.Phone, t]);

  const SectionHeader: React.FC<{ section: any }> = ({ section }) => {
    const title = String(getSectionTitle(section) || '').trim();
    if (!title) return null;
    return (
      <div className="w-full bg-sky-50 border-y border-slate-200 px-4 py-1 border-l-[3px]" style={{ borderLeftColor: color }}>
        <div className="text-base font-bold" style={{ color }}>
          {title}
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
          {item.description ? <RichText html={item.description} className="text-slate-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} /> : null}
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

  const renderable = React.useMemo(() => {
    return orderedSections.filter((section) => {
      if (!section?.isVisible) return false;
      if (section.type === ResumeSectionType.Summary) return false;
      if (section.type === ResumeSectionType.Exam) return true;
      const items = getOrderedItems(section.items || []);
      return items.length > 0;
    });
  }, [orderedSections]);

  const headerTitle = String(personal?.FullName || '').trim();

  const InfoRow: React.FC<{ icon: React.ReactNode; label: string; value: string }> = ({ icon, label, value }) => (
    <div className="flex items-center gap-3 min-w-0">
      <div className="w-4 h-4 flex items-center justify-center flex-shrink-0" style={{ color }}>
        {icon}
      </div>
      <div className="text-slate-600 whitespace-nowrap flex-shrink-0">{label}：</div>
      <div className={`min-w-0 font-semibold text-slate-900 ${String(value || '').includes('@') ? 'break-all' : 'break-words'}`}>{value}</div>
    </div>
  );

  const extraHeaderCells = React.useMemo(() => {
    const excludedLabels = new Set<string>([
      String(findWorkYearsPair?.label || '').trim(),
      ...intentCells.map((x) => x.label),
      ...basicCells.map((x) => x.label),
    ]);

    const cityValue = String(personal?.City || '').trim();

    return customPairs
      .map((p) => ({ icon: <User size={16} />, label: String(p.label || '').trim(), value: String(p.value || '').trim() }))
      .filter((p) => {
        if (!p.label || !p.value) return false;
        if (excludedLabels.has(p.label)) return false;
        if (cityValue && p.value === cityValue && (p.label.includes('城市') || p.label === t('editor.fields.city'))) return false;
        return true;
      });
  }, [basicCells, customPairs, findWorkYearsPair?.label, intentCells, personal?.City, t]);

  const headerInfo = React.useMemo(() => {
    const cells = [...intentCells, ...basicCells, ...extraHeaderCells];
    return {
      left: cells.filter((_, idx) => idx % 2 === 0),
      right: cells.filter((_, idx) => idx % 2 === 1),
    };
  }, [basicCells, extraHeaderCells, intentCells]);

  const hasHeaderInfo = headerInfo.left.length > 0 || headerInfo.right.length > 0;

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <header className={`flex items-start gap-8 ${headerSpaceClass}`}>
        <div className="flex-shrink-0 mt-3">
          {personal?.AvatarURL ? (
            <img src={personal.AvatarURL} alt={t('a11y.avatarAlt')} className={getAvatarPhotoClassName('shadow-none')} style={{ backgroundColor: '#ffffff' }} />
          ) : (
            <div className={getAvatarPlaceholderClassName()} />
          )}
        </div>

        <div className="flex-1 min-w-0">
          <div className="text-4xl font-bold tracking-wide" style={{ color }}>
            {headerTitle}
          </div>

          {hasHeaderInfo ? (
            <div className={`mt-3 grid grid-cols-2 gap-x-14 gap-y-2 ${getHeaderInfoTextClassName('text-slate-800')}`}>
              <div className="space-y-2">
                {headerInfo.left.map((x, idx) => (
                  <InfoRow key={`h-l-${x.label}-${idx}`} icon={x.icon} label={x.label} value={x.value} />
                ))}
              </div>
              <div className="space-y-2">
                {headerInfo.right.map((x, idx) => (
                  <InfoRow key={`h-r-${x.label}-${idx}`} icon={x.icon} label={x.label} value={x.value} />
                ))}
              </div>
            </div>
          ) : null}
        </div>
      </header>

      <div className={contentGapClass}>
        {renderable.map((section) => {
          const body = renderSectionBody(section);
          if (!body) return null;
          return (
            <section key={section.id} className="pt-1">
              <SectionHeader section={section} />
              <div className="mt-4">{body}</div>
            </section>
          );
        })}
      </div>
    </div>
  );
};

export default TemplateBlueNationResume;
