import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { Mail, Phone, MapPin, GraduationCap, User, Briefcase, Wrench, Layers, BookOpen, Award, Heart, Image as ImageIcon } from 'lucide-react';
import { ExamSection } from './shared/ExamSection';
import { RichText } from './shared/RichText';
import { formatDateRange, getAccentColor, getAvatarPhotoClassName, getAvatarPlaceholderClassName, getOrderedItems, getOrderedVisibleSections, getSpacingTokens, normalizeCustomPairs, parseCustomPairs } from './shared/templateTokens';

export const TemplateMintTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const color = getAccentColor(data, '#14b8a6');
  const { spacingMode, lineHeight, contentGapClass, listTightClass, listMediumClass } = getSpacingTokens(styles);

  const findJobTarget = () => {
    return data.Personal?.Job || '';
  };
  const findCity = () => {
    return data.Personal?.City || '';
  };
  const findMoney = () => {
    return data.Personal?.Money || '';
  };
  const findJoinTime = () => {
    return data.Personal?.JoinTime || '';
  };
  const findDegree = () => {
    return data.Personal?.Degree || '';
  };

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Experience: return <Briefcase size={18} />;
      case ResumeSectionType.Education: return <GraduationCap size={18} />;
      case ResumeSectionType.Skills: return <Wrench size={18} />;
      case ResumeSectionType.Projects: return <Layers size={18} />;
      case ResumeSectionType.Internships: return <Briefcase size={18} />;
      case ResumeSectionType.Portfolio: return <ImageIcon size={18} />;
      case ResumeSectionType.Awards: return <Award size={18} />;
      case ResumeSectionType.Interests: return <Heart size={18} />;
      case ResumeSectionType.Exam: return <BookOpen size={18} />;
      case ResumeSectionType.SelfEvaluation: return <User size={18} />;
      default: return <User size={18} />;
    }
  };

  const renderItemTime = (item: any) => {
    const range = formatDateRange(item, t, { separatorVariant: 'dash', normalizeMonthSeparator: '.' });
    if (!range) return null;
    return (
      <span className="text-sm text-gray-600 font-bold" style={{ fontSize: styles.fontSize }}>
        {range}
      </span>
    );
  };

  const getRightTag = (sectionType: ResumeSectionType, item: any) => {
    if (sectionType === ResumeSectionType.Experience) {
      return item?.subtitle || null;
    }
    if (sectionType === ResumeSectionType.Internships) {
      return item?.subtitle || null;
    }
    if (sectionType === ResumeSectionType.Education) {
      if (item?.major || item?.degree) {
        const major = item?.major || '';
        const degree = item?.degree || '';
        if (major && degree) return `${major}（${degree}）`;
        return major || degree || null;
      }
      return null;
    }
    if (sectionType === ResumeSectionType.Projects) {
      return item?.subtitle || null;
    }
    return null;
  };
  
  const showMarkerForSection = (sectionType: ResumeSectionType) => {
    return [
      ResumeSectionType.Experience,
      ResumeSectionType.Education,
      ResumeSectionType.Projects,
      ResumeSectionType.Internships
    ].includes(sectionType);
  };

  const sectionsOrdered = React.useMemo(() => {
    return getOrderedVisibleSections(data.sections || []);
  }, [data.sections]);
  const headerPaddingClassName = spacingMode === 'compact' ? 'px-10 pt-6 pb-5' : spacingMode === 'spacious' ? 'px-10 pt-10 pb-7' : 'px-10 pt-8 pb-6';
  const headerRowGapClassName = spacingMode === 'compact' ? 'gap-5' : spacingMode === 'spacious' ? 'gap-8' : 'gap-6';
  const headerJobMarginTopClassName = spacingMode === 'compact' ? 'mt-1' : spacingMode === 'spacious' ? 'mt-3' : 'mt-2';
  const headerInfoMarginTopClassName = spacingMode === 'compact' ? 'mt-3' : spacingMode === 'spacious' ? 'mt-5' : 'mt-4';
  const headerInfoGapClassName = spacingMode === 'compact' ? 'gap-3' : spacingMode === 'spacious' ? 'gap-5' : 'gap-4';

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight, fontSize: styles.fontSize }}>
      <div className="relative">
        <div className="bg-teal-500" style={{ backgroundColor: color }}>
          <div className={`${headerPaddingClassName} flex items-start ${headerRowGapClassName}`}>
            <div className="flex-1">
              <h1 className="text-3xl font-bold tracking-wide text-white">{data.Personal?.FullName}</h1>
              {findJobTarget() && <p className={`${headerJobMarginTopClassName} text-white/90`}>{t('editor.fields.jobApplication')}：{findJobTarget()}</p>}
              <div className={`${headerInfoMarginTopClassName} flex flex-wrap ${headerInfoGapClassName} text-white/90`}>
                {data.Personal?.Age && (
                  <span className="inline-flex items-center gap-1">
                    <User size={16} /> {data.Personal.Age}
                  </span>
                )}
                {data.Personal?.Phone && (
                  <span className="inline-flex items-center gap-1">
                    <Phone size={16} /> {data.Personal.Phone}
                  </span>
                )}
                {data.Personal?.Gender && (
                  <span className="inline-flex items-center gap-1">
                    <User size={16} /> {data.Personal.Gender}
                  </span>
                )}
                {data.Personal?.Email && (
                  <span className="inline-flex items-center gap-1">
                    <Mail size={16} /> {data.Personal.Email}
                  </span>
                )}
                {findDegree() && (
                  <span className="inline-flex items-center gap-1">
                    <GraduationCap size={16} /> {findDegree()}
                  </span>
                )}
                {findCity() && (
                  <span className="inline-flex items-center gap-1">
                    <MapPin size={16} /> {findCity()}
                  </span>
                )}
                {findMoney() && (
                  <span className="inline-flex items-center gap-1">
                    {t('editor.fields.expectedSalary')}：{findMoney()}
                  </span>
                )}
                {findJoinTime() && (
                  <span className="inline-flex items-center gap-1">
                    {t('editor.fields.joinTime')}：{findJoinTime()}
                  </span>
                )}
                {normalizeCustomPairs(parseCustomPairs(data.Personal?.CustomInfo)).map((ci, idx) => (
                  <span key={idx} className="inline-flex items-center gap-1">
                    {ci.label ? `${ci.label}: ${ci.value}` : ci.value}
                  </span>
                ))}
              </div>
            </div>
            <div className="flex-shrink-0">
              {data.Personal?.AvatarURL ? (
                <img
                  src={data.Personal.AvatarURL}
                  alt={t('a11y.avatarAlt')}
                  className={getAvatarPhotoClassName()}
                  style={{ backgroundColor: '#ffffff' }}
                />
              ) : (
                <div className={getAvatarPlaceholderClassName('bg-white/30')} />
              )}
            </div>
          </div>
        </div>
        {/* <div className="w-full" style={{ backgroundColor: color }}>
        {/* <div className="w-full" style={{ backgroundColor: color }}>
          <svg viewBox="0 0 1200 60" preserveAspectRatio="none" className="w-full h-6 block">
            <path d="M0 30 Q 300 60 600 30 T 1200 30 L 1200 60 L 0 60 Z" fill="white"></path>
          </svg>
        </div> */}
      </div>

      <div className={`px-10 pt-6 pb-10 ${contentGapClass}`}>
                {sectionsOrdered.map((section, idx) => (
                  (() => {
                    const items = getOrderedItems(section.items || []);
                    const isLast = idx === sectionsOrdered.length - 1;
                    const dashTop = 24;
                    const dashBottom = isLast ? 2 : -55;
                    return (
                      <div key={section.id} className="flex gap-6 items-start">
                        <div className="relative mt-4 -translate-y-1/2" style={{ width: 140 }}>
                          <div className="text-base font-bold tracking-wide pl-3" style={{ color }}>
                            {(() => {
                              const title = String(getSectionTitle(section) || '').trim();
                              const parts = title.match(/^([A-Za-z]+)\s+([A-Za-z]+)$/);
                              if (parts) {
                                return (
                                  <span className="inline-block leading-tight">
                                    <span className="block">{parts[1]}</span>
                                    <span className="block">{parts[2]}</span>
                                  </span>
                                );
                              }
                              return title;
                            })()}
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
                    ) : section.type === ResumeSectionType.Skills ? (
                      <div className={listTightClass}>
                        {items.map(item => (
                          <div key={item.id} className="relative">
                            {showMarkerForSection(section.type) && (
                              <div className="absolute -translate-x-1/2 top-2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                            )}
                            {item.description && (
                              <RichText html={item.description} className="text-gray-700" fontSize={styles.fontSize} lineHeight={lineHeight} />
                            )}
                          </div>
                        ))}
                      </div>
                ) : (
                      <div className={items.length > 1 ? listMediumClass : ''}>
                        {items.map(item => (
                          <div key={item.id} className="relative">
                            <div className="relative flex justify-between items-baseline mb-1">
                              {showMarkerForSection(section.type) && (
                                <div className="absolute -translate-x-1/2 top-1/2 -translate-y-1/2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                              )}
                              {renderItemTime(item) && (
                                <div className="absolute -translate-y-1/2 text-gray-600 text-sm whitespace-nowrap" style={{ left: -164, top: '50%' }}>
                                  {renderItemTime(item)}
                                </div>
                              )}
                              <div className="font-semibold text-gray-900 text-sm">
                                {item.title}
                              </div>
                              {getRightTag(section.type, item) && (
                                <span className="text-gray-900 font-medium">{getRightTag(section.type, item)}</span>
                              )}
                            </div>
                            {item.description && (
                              <RichText html={item.description} className="text-gray-700 mt-1" fontSize={styles.fontSize} lineHeight={lineHeight} />
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            );
          })()
        ))}
      </div>
    </div>
  );
};

export default TemplateMintTimeline;
