import React from 'react';
import { Mail, Phone, MapPin, GraduationCap, User, Briefcase, Wrench, Layers, BookOpen, Award, Heart, Image as ImageIcon } from 'lucide-react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

const ScallopEdge: React.FC = () => {
  const radius = 12;
  const gap = 10;
  const diameter = radius * 2;
  const count = 64;
  return (
    <div className="absolute left-0 right-0 bottom-0 pointer-events-none" style={{ height: `${radius}px` }} aria-hidden="true">
      <div className="w-full h-full overflow-hidden">
        <div className="flex" style={{ gap: `${gap}px`, paddingLeft: `${radius}px` }}>
          {Array.from({ length: count }).map((_, idx) => (
            <div key={idx} className="shrink-0" style={{ width: `${diameter}px`, height: `${diameter}px`, backgroundColor: '#ffffff', borderRadius: 9999 }} />
          ))}
        </div>
      </div>
    </div>
  );
};

export const TemplateWaveTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const color = getAccentColor(data, '#2f7a75');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const sectionsOrdered = React.useMemo(() => {
    return getOrderedVisibleSections(data.sections || []);
  }, [data.sections]);

  const customPairs = React.useMemo(() => {
    return normalizeCustomPairs(parseCustomPairs(data.Personal?.CustomInfo));
  }, [data.Personal?.CustomInfo]);

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Experience:
      case ResumeSectionType.Internships:
        return <Briefcase size={18} />;
      case ResumeSectionType.Education:
        return <GraduationCap size={18} />;
      case ResumeSectionType.Skills:
        return <Wrench size={18} />;
      case ResumeSectionType.Projects:
        return <Layers size={18} />;
      case ResumeSectionType.Portfolio:
        return <ImageIcon size={18} />;
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

  const formatRange = (item: any) => formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });

  const isTimelineLike = (sectionType: ResumeSectionType) => {
    return [ResumeSectionType.Experience, ResumeSectionType.Education, ResumeSectionType.Projects, ResumeSectionType.Internships].includes(sectionType);
  };

  const hasMeaningfulContent = (item: any, type: ResumeSectionType) => {
    if (!item) return false;
    if (type === ResumeSectionType.Skills) {
      return !!(item.description && String(item.description).trim());
    }
    return !!(
      (item.title && String(item.title).trim()) ||
      (item.subtitle && String(item.subtitle).trim()) ||
      (item.major && String(item.major).trim()) ||
      (item.degree && String(item.degree).trim()) ||
      (item.description && String(item.description).trim()) ||
      (item.timeStart && String(item.timeStart).trim()) ||
      (item.timeEnd && String(item.timeEnd).trim()) ||
      item.today
    );
  };

  const renderItemTime = (item: any) => {
    const range = formatRange(item);
    if (!range) return null;
    return (
      <span className="text-sm text-gray-600 font-bold" style={{ fontSize: styles.fontSize }}>
        {range}
      </span>
    );
  };

  const showMarkerForSection = (sectionType: ResumeSectionType) => {
    return isTimelineLike(sectionType);
  };

  const personal = data.Personal || {};

  const infoItems = React.useMemo(() => {
    const items: Array<{ key: string; icon?: React.ReactNode; text: string }> = [];
    const age = String(personal.Age || '').trim();
    if (age) items.push({ key: 'age', icon: <User size={16} />, text: age });
    const phone = String(personal.Phone || '').trim();
    if (phone) items.push({ key: 'phone', icon: <Phone size={16} />, text: phone });
    const email = String(personal.Email || '').trim();
    if (email) items.push({ key: 'email', icon: <Mail size={16} />, text: email });
    const city = String(personal.City || '').trim();
    if (city) items.push({ key: 'city', icon: <MapPin size={16} />, text: city });
    const degree = String(personal.Degree || '').trim();
    if (degree) items.push({ key: 'degree', icon: <GraduationCap size={16} />, text: degree });
    const money = String(personal.Money || '').trim();
    if (money) items.push({ key: 'money', text: `${t('editor.fields.expectedSalary')}：${money}` });
    const joinTime = String(personal.JoinTime || '').trim();
    if (joinTime) items.push({ key: 'joinTime', text: `${t('editor.fields.joinTime')}：${joinTime}` });
    customPairs.forEach((p, idx) => {
      const text = p.label ? `${p.label} ${p.value}` : p.value;
      const safe = String(text || '').trim();
      if (safe) items.push({ key: `custom-${idx}`, text: safe });
    });
    return items;
  }, [customPairs, personal.Age, personal.Phone, personal.Email, personal.City, personal.Degree, personal.Money, personal.JoinTime, t]);
  const headerPaddingClassName = spacingMode === 'compact' ? 'px-10 pt-6 pb-8' : spacingMode === 'spacious' ? 'px-10 pt-10 pb-12' : 'px-10 pt-8 pb-10';
  const headerRowGapClassName = spacingMode === 'compact' ? 'gap-6' : spacingMode === 'spacious' ? 'gap-10' : 'gap-8';
  const headerJobMarginTopClassName = spacingMode === 'compact' ? 'mt-1' : spacingMode === 'spacious' ? 'mt-3' : 'mt-2';
  const headerInfoMarginTopClassName = spacingMode === 'compact' ? 'mt-3' : spacingMode === 'spacious' ? 'mt-5' : 'mt-4';
  const headerInfoGapClassName =
    spacingMode === 'compact'
      ? 'gap-x-4 gap-y-2'
      : spacingMode === 'spacious'
        ? 'gap-x-6 gap-y-4'
        : 'gap-x-5 gap-y-3';

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="relative">
        <div className="relative" style={{ backgroundColor: color }}>
          <div className={`${headerPaddingClassName} flex items-start ${headerRowGapClassName}`}>
            <div className="flex-1 min-w-0">
              <h1 className="text-3xl font-bold tracking-wide text-white">{personal.FullName}</h1>
              {personal.Job ? <p className={`${headerJobMarginTopClassName} text-white/90`}>{personal.Job}</p> : null}

              {infoItems.length ? (
                <div className={`${headerInfoMarginTopClassName} flex flex-wrap ${headerInfoGapClassName} text-white/90`}>
                  {infoItems.map((it) => (
                    <span key={it.key} className="inline-flex items-center gap-1">
                      {it.icon ? it.icon : null}
                      {it.text}
                    </span>
                  ))}
                </div>
              ) : null}
            </div>

            <div className="flex-shrink-0">
              {personal.AvatarURL ? (
                <img
                  src={personal.AvatarURL}
                  alt={t('a11y.avatarAlt')}
                  className={getAvatarPhotoClassName()}
                  style={{ backgroundColor: '#ffffff' }}
                />
              ) : (
                <div className={getAvatarPlaceholderClassName('bg-white/30')} />
              )}
            </div>
          </div>

          <ScallopEdge />
        </div>
      </div>

      <div className={`px-10 pt-8 pb-10 ${contentGapClass}`}>
        {(() => {
          const renderable = sectionsOrdered.filter((section) => {
            if (section.type === ResumeSectionType.Exam) return true;
            const items = getOrderedItems(section.items || []).filter((it: any) => hasMeaningfulContent(it, section.type));
            return items.length > 0;
          });

          return renderable.map((section, idx) => {
          const items = getOrderedItems(section.items || []).filter((it: any) => hasMeaningfulContent(it, section.type));
          if (!items.length && section.type !== ResumeSectionType.Exam) return null;
          const isLast = idx === renderable.length - 1;
          const dashTop = 24;
          const dashBottom = isLast ? 2 : -55;
          const title = String(getSectionTitle(section) || '').trim();
          const timelineLike = isTimelineLike(section.type);

          return (
            <div key={section.id} className="flex gap-6 items-start">
              <div className="relative mt-4 -translate-y-1/2" style={{ width: 140 }}>
                <div className="text-base font-bold tracking-wide pl-3" style={{ color }}>
                  {title}
                </div>
              </div>

              <div className="flex-1">
                <div className="mt-4 border-t-2 pt-4 pb-4 relative" style={{ borderColor: color }}>
                  <div className="absolute border-l border-dashed" style={{ borderColor: color, left: -32, top: dashTop, bottom: dashBottom }} />
                  <div className="absolute left-[-48px] w-8 h-8 rounded-full flex items-center justify-center text-white ring-2 ring-white/60 print:ring-0 z-10" style={{ backgroundColor: color, top: -17 }}>
                    <SectionIcon type={section.type} />
                  </div>

                  {section.type === ResumeSectionType.Exam ? (
                    <ExamSection section={section} color={color} t={t} />
                  ) : section.type === ResumeSectionType.SelfEvaluation ? (
                    <div className={listTightClass}>
                      {items.map((item: any) => (
                        <div key={item.id} className="relative">
                          {item.title ? <div className="text-sm font-semibold text-gray-900">{item.title}</div> : null}
                          {item.subtitle ? <div className="text-sm text-gray-700 mt-0.5">{item.subtitle}</div> : null}
                          {item.description ? (
                            <RichText html={item.description} className="text-gray-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
                          ) : null}
                        </div>
                      ))}
                    </div>
                  ) : section.type === ResumeSectionType.Skills ? (
                    <div className={listTightClass}>
                      {items.map((item: any) => (
                        <div key={item.id} className="relative">
                          {showMarkerForSection(section.type) ? (
                            <div className="absolute -translate-x-1/2 top-2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                          ) : null}
                          {item.description ? (
                            <RichText html={item.description} className="text-gray-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
                          ) : null}
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className={items.length > 1 ? listMediumClass : ''}>
                      {items.map((item: any) => {
                        const right = String(item.subtitle || '').trim();
                        return (
                          <div key={item.id} className="relative">
                            <div className="relative flex justify-between items-baseline mb-1">
                              {showMarkerForSection(section.type) ? (
                                <div className="absolute -translate-x-1/2 top-1/2 -translate-y-1/2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                              ) : null}

                              {timelineLike && renderItemTime(item) ? (
                                <div className="absolute -translate-y-1/2 text-gray-600 text-sm whitespace-nowrap" style={{ left: -164, top: '50%' }}>
                                  {renderItemTime(item)}
                                </div>
                              ) : null}

                              <div className="font-semibold text-gray-900 text-sm">
                                {item.title}
                              </div>

                              {right ? (
                                <span className="text-gray-900 font-medium">{right}</span>
                              ) : null}
                            </div>

                            {item.description ? (
                              <RichText html={item.description} className="text-gray-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
                            ) : null}
                          </div>
                        );
                      })}
                    </div>
                  )}
                </div>
              </div>
            </div>
          );
        });
        })()}
      </div>
    </div>
  );
};

export default TemplateWaveTimeline;
