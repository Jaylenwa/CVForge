import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import {
  formatDateRange,
  getAccentColor,
  getAvatarPhotoClassName,
  getAvatarPlaceholderClassName,
  getHeaderInfoTextClassName,
  getOrderedItems,
  getOrderedVisibleSections,
  getSpacingTokens,
  normalizeCustomPairs,
  parseCustomPairs,
} from './shared/templateTokens';

export const TemplateBlueFormSheet: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const personal = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
  const isZh = (data.language || 'zh') === 'zh';
  const color = getAccentColor(data, '#2c80b9');
  const { lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const colon = isZh ? '：' : ':';

  const orderedSections = React.useMemo(() => getOrderedVisibleSections(data.sections || []), [data.sections]);
  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'tilde', normalizeMonthSeparator: '-' });

  const customPairs = React.useMemo(() => normalizeCustomPairs(parseCustomPairs(personal?.CustomInfo)), [personal?.CustomInfo]);

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

  const baseInfo = React.useMemo(() => {
    const workYearsValue = findWorkYearsPair ? String(findWorkYearsPair.value || '').trim() : '';
    const workYearsLabel = findWorkYearsPair ? String(findWorkYearsPair.label || '').trim() : '';

    const left: Array<{ label: string; value: string }> = [
      { label: isZh ? '姓名' : t('editor.fields.fullName'), value: String(personal?.FullName || '').trim() },
      { label: isZh ? '性别' : t('editor.fields.gender'), value: String(personal?.Gender || '').trim() },
      { label: workYearsLabel || (isZh ? '工作年限' : 'Work Years'), value: workYearsValue },
      { label: isZh ? '邮箱' : t('editor.fields.email'), value: String(personal?.Email || '').trim() },
    ].filter((x) => x.value);

    const right: Array<{ label: string; value: string }> = [
      { label: isZh ? '年龄' : t('editor.fields.age'), value: normalizedAge },
      { label: isZh ? '籍贯' : t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { label: isZh ? '电话' : t('editor.fields.phone'), value: String(personal?.Phone || '').trim() },
    ].filter((x) => x.value);

    const excludedLabels = new Set<string>([workYearsLabel, ...left.map((x) => x.label), ...right.map((x) => x.label)].filter(Boolean));
    const cityValue = String(personal?.City || '').trim();
    const extras = customPairs
      .map((p) => ({ label: String(p.label || '').trim(), value: String(p.value || '').trim() }))
      .filter((p) => {
        if (!p.label || !p.value) return false;
        if (excludedLabels.has(p.label)) return false;
        if (cityValue && p.value === cityValue && (p.label.includes('城市') || p.label.includes('籍贯') || p.label === t('editor.fields.city'))) return false;
        return true;
      });

    const balancedLeft = [...left];
    const balancedRight = [...right];
    extras.forEach((p) => {
      if (balancedLeft.length <= balancedRight.length) balancedLeft.push(p);
      else balancedRight.push(p);
    });

    return { left: balancedLeft, right: balancedRight };
  }, [customPairs, findWorkYearsPair, isZh, normalizedAge, personal?.City, personal?.Email, personal?.FullName, personal?.Gender, personal?.Phone, t]);

  const basicInfoWithIntent = React.useMemo(() => {
    const intentPairs: Array<{ label: string; value: string }> = [
      { label: isZh ? '职位' : t('editor.fields.jobApplication'), value: String(personal?.Job || '').trim() },
      { label: isZh ? '城市' : t('editor.fields.city'), value: String(personal?.City || '').trim() },
      { label: isZh ? '期望薪资' : t('editor.fields.expectedSalary'), value: String(personal?.Money || '').trim() },
      { label: isZh ? '到岗时间' : t('editor.fields.joinTime'), value: String(personal?.JoinTime || '').trim() },
    ].filter((x) => x.value);

    const left = [...baseInfo.left];
    const right = [...baseInfo.right];
    intentPairs.forEach((p) => {
      if (left.length <= right.length) left.push(p);
      else right.push(p);
    });
    return { left, right };
  }, [baseInfo.left, baseInfo.right, isZh, personal?.City, personal?.Job, personal?.JoinTime, personal?.Money, t]);

  const BlockHeader: React.FC<{ title: string }> = ({ title }) => (
    <div className="w-full bg-sky-100 h-8 flex items-center">
      <div className="h-full flex items-center px-6 text-white font-bold" style={{ backgroundColor: color }}>
        {title}
      </div>
      <div className="flex-1" />
    </div>
  );

  const FieldRow: React.FC<{ label: string; value: string }> = ({ label, value }) => (
    <div className="flex items-baseline gap-2 min-w-0">
      <div className="whitespace-nowrap text-slate-700 font-semibold">{label + colon}</div>
      <div className="min-w-0 text-slate-900 break-words">{value}</div>
    </div>
  );

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

  const headerTitle = isZh ? '个人简历' : 'Resume';

  return (
    <div className={`w-full bg-white text-slate-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none p-10`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <header className="text-center">
        <div className="text-3xl font-bold text-slate-900">{headerTitle}</div>
      </header>

      <div className="mt-6 space-y-6">
        <section>
          <BlockHeader title={isZh ? '基本信息' : 'Basic Information'} />
          <div className={`mt-5 grid grid-cols-12 gap-8 items-start ${getHeaderInfoTextClassName('text-slate-800')}`}>
            <div className="col-span-5 space-y-3">
              {basicInfoWithIntent.left.map((p, idx) => (
                <FieldRow key={`bi-l-${p.label}-${idx}`} label={p.label} value={p.value} />
              ))}
            </div>
            <div className="col-span-5 space-y-3">
              {basicInfoWithIntent.right.map((p, idx) => (
                <FieldRow key={`bi-r-${p.label}-${idx}`} label={p.label} value={p.value} />
              ))}
            </div>
            <div className="col-span-2 flex justify-end">
              {personal?.AvatarURL ? (
                <img src={personal.AvatarURL} alt={t('a11y.avatarAlt')} className={getAvatarPhotoClassName('shadow-none')} style={{ backgroundColor: '#ffffff' }} />
              ) : (
                <div className={getAvatarPlaceholderClassName()} />
              )}
            </div>
          </div>
        </section>

        <div className={contentGapClass}>
          {renderable.map((section) => {
            const body = renderSectionBody(section);
            if (!body) return null;
            const title = String(getSectionTitle(section) || '').trim();
            if (!title) return null;
            return (
              <section key={section.id}>
                <BlockHeader title={title} />
                <div className="mt-4">{body}</div>
              </section>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default TemplateBlueFormSheet;
