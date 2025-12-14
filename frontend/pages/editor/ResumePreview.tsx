import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeData, ResumeSectionType, ThemeConfig } from '../../types';

interface PreviewProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
}

// Map Font IDs to CSS Font Stacks
const getFontStack = (fontId: string): string => {
    switch (fontId) {
        case 'inter': return '"Inter", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
        case 'roboto': return '"Roboto", "Helvetica Neue", Arial, sans-serif';
        case 'merriweather': return '"Merriweather", "Georgia", serif';
        case 'playfair': return '"Playfair Display", "Times New Roman", serif';
        case 'mono': return '"Roboto Mono", "Courier New", monospace';
        
        // Chinese Fonts
        case 'yahei': return '"Microsoft YaHei", "Heiti SC", "PingFang SC", sans-serif';
        case 'notosans': return '"Noto Sans SC", "Microsoft YaHei", sans-serif';
        case 'simsun': return '"SimSun", "Songti SC", "Noto Serif SC", serif';
        case 'kaiti': return '"KaiTi", "STKaiti", "Kai", serif';
        
        default: return '"Inter", sans-serif';
    }
};

// Helper to get CSS styles based on ThemeConfig
const getThemeStyles = (config?: ThemeConfig) => {
    const fontFamily = getFontStack(config?.fontFamily || 'inter');

    const spacingMultiplier = {
        'compact': '0.85',
        'normal': '1',
        'spacious': '1.25'
    }[config?.spacing || 'normal'];

  return { fontFamily, spacingMultiplier };
};

const hasExtraPersonalInfo = (data: ResumeData) => {
    const p = data.personalInfo || ({} as ResumeData['personalInfo']);
    return !!(p.gender || p.age || p.maritalStatus || p.politicalStatus || p.birthplace || p.ethnicity || p.height || p.weight || (p.customInfo && p.customInfo.length > 0));
};

const sanitizeHtml = (html: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html || '', 'text/html');
  const allowedTags = new Set(['b','strong','i','em','u','br','p','div','ul','ol','li','span','a']);
  const escapeText = (s: string) => s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
  const sanitizeNode = (node: Node): string => {
    if (node.nodeType === Node.TEXT_NODE) {
      return escapeText(node.textContent || '');
    }
    if (node.nodeType !== Node.ELEMENT_NODE) return '';
    const el = node as HTMLElement;
    const tag = el.tagName.toLowerCase();
    if (!allowedTags.has(tag)) {
      let s = '';
      el.childNodes.forEach(child => { s += sanitizeNode(child); });
      return s;
    }
    let attrs = '';
    if (tag === 'a') {
      const raw = el.getAttribute('href') || '';
      try {
        const u = new URL(raw, window.location.origin);
        const proto = u.protocol.replace(':','');
        if (['http','https','mailto'].includes(proto)) {
          attrs = ` href="${escapeText(raw)}" rel="noopener noreferrer nofollow"`;
        }
      } catch {}
    }
    if (tag === 'br') return '<br/>';
    let content = '';
    el.childNodes.forEach(child => { content += sanitizeNode(child); });
    return `<${tag}${attrs}>${content}</${tag}>`;
  };
  let out = '';
  doc.body.childNodes.forEach(n => { out += sanitizeNode(n); });
  return out;
};

// Helper hook for section title translation
const useSectionTitle = () => {
    const { t } = useLanguage();
    
    return (section: any) => {
        const defaultTitles = ['Professional Summary', 'Work Experience', 'Education', 'Skills', 'Projects'];
        if (defaultTitles.includes(section.title) || !section.title) {
            return t(`section.${section.type}`);
        }
        return section.title;
    };
};

// 1. Classic Professional
const TemplateClassic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  return (
  <div className={`w-full p-8 md:p-12 bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:p-0 print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
    <div className="border-b-2 pb-6 mb-6 flex flex-col md:flex-row items-center md:items-start gap-6" style={{ borderColor: data.themeConfig?.color || '#333' }}>
      <div className="flex-1 text-center md:text-left order-2 md:order-1">
          <h1 className="text-4xl font-bold uppercase tracking-wider" style={{ color: data.themeConfig?.color }}>{data.personalInfo.fullName}</h1>
          <p className="text-xl mt-2 text-gray-600">{data.personalInfo.jobTitle}</p>
          <div className="mt-4 flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-600">
            {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
            {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
          </div>
          {hasExtraPersonalInfo(data) && (
            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-600">
              {data.personalInfo.gender && <span>{t('editor.fields.gender')}: {data.personalInfo.gender}</span>}
              {data.personalInfo.age && <span>{t('editor.fields.age')}: {data.personalInfo.age}</span>}
              {data.personalInfo.maritalStatus && <span>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</span>}
              {data.personalInfo.politicalStatus && <span>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</span>}
              {data.personalInfo.birthplace && <span>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</span>}
              {data.personalInfo.ethnicity && <span>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</span>}
              {data.personalInfo.height && <span>{t('editor.fields.height')}: {data.personalInfo.height}</span>}
              {data.personalInfo.weight && <span>{t('editor.fields.weight')}: {data.personalInfo.weight}</span>}
              {(data.personalInfo.customInfo || []).map((ci, idx) => (
                <span key={idx}>{ci.label}: {ci.value}</span>
              ))}
            </div>
          )}
      </div>
      {data.personalInfo.avatarUrl && (
          <div className="order-1 md:order-2 flex-shrink-0">
             <img 
               src={data.personalInfo.avatarUrl} 
               alt="Profile" 
               className="w-32 h-32 rounded-lg object-cover border-2 shadow-sm"
               style={{ borderColor: data.themeConfig?.color || '#e5e7eb' }}
             />
          </div>
      )}
    </div>

    <div className="space-y-6" style={{ gap: `${parseFloat(styles.spacingMultiplier) * 1.5}rem` }}>
      {data.sections.filter(s => s.isVisible).map(section => (
        <div key={section.id}>
          <h3 className="text-lg font-bold uppercase border-b mb-4 pb-1 tracking-wide" style={{ borderColor: '#e5e7eb', color: '#1f2937' }}>
            {getSectionTitle(section)}
          </h3>
          <div className="space-y-4">
            {section.items.map(item => (
              <div key={item.id} className="relative">
                {section.type !== ResumeSectionType.Skills && (
                  <div className="flex justify-between items-baseline mb-1">
                    <h4 className="font-bold text-gray-900">{item.title}</h4>
                    {item.dateRange && <span className="text-sm italic text-gray-600">{item.dateRange}</span>}
                  </div>
                )}
                {item.subtitle && (
                   <div className="mb-1">
                    <span className="font-semibold text-gray-700">{item.subtitle}</span>
                   </div>
                )}
                <div className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  </div>
);
};

const TemplateCNBlue: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const SectionHeader: React.FC<{ title: string }> = ({ title }) => (
    <div className="bg-blue-100 text-blue-900 font-bold px-4 py-2 rounded">{title}</div>
  );
  const SkillBar: React.FC<{ text: string }> = ({ text }) => {
    const raw = text || '';
    const name = raw.includes(':') ? raw.split(':')[0] : (raw.includes(' ') ? raw.split(' ')[0] : raw);
    const lower = raw.toLowerCase();
    let level = 3;
    let label = '良好';
    if (raw.includes('精通') || lower.includes('expert') || lower.includes('advanced')) { level = 5; label = '精通'; }
    else if (raw.includes('熟练') || lower.includes('proficient')) { level = 4; label = '熟练'; }
    else if (raw.includes('一般') || lower.includes('beginner')) { level = 2; label = '一般'; }
    const pct = `${Math.min(100, Math.max(0, level * 20))}%`;
    return (
      <div className="mt-2">
        <div className="flex justify-between items-center text-sm">
          <span className="font-medium text-gray-800">{name}</span>
          <span className="text-gray-500">{label}</span>
        </div>
        <div className="h-1.5 bg-gray-200 rounded mt-1 overflow-hidden">
          <div className="h-full bg-blue-500" style={{ width: pct }} />
        </div>
      </div>
    );
  };
  return (
    <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[12mm] -z-10" style={{ backgroundColor: '#4b5563' }}></div>
      <div className="px-10 pt-10">
        <div className="text-center mb-6">
          <h1 className="text-3xl font-bold text-blue-700 tracking-wide">个人简历</h1>
          <p className="text-gray-500 mt-1 text-sm">{data.personalInfo.jobTitle || '努力超越自己，每天进步一点点'}</p>
        </div>
        <div className="mb-8">
          <SectionHeader title="基本信息" />
          <div className="grid grid-cols-1 md:grid-cols-12 gap-6 mt-4 items-start">
            <div className="md:col-span-9 grid grid-cols-2 gap-y-2 text-sm">
              <div><span className="text-gray-500 mr-2">姓名：</span><span className="text-gray-900 font-medium">{data.personalInfo.fullName}</span></div>
              <div><span className="text-gray-500 mr-2">电话：</span><span className="text-gray-900">{data.personalInfo.phone}</span></div>
              <div className="col-span-2"><span className="text-gray-500 mr-2">邮箱：</span><span className="text-gray-900 break-all">{data.personalInfo.email}</span></div>
              {/* website removed */}
              {data.personalInfo.gender && <div><span className="text-gray-500 mr-2">性别：</span><span className="text-gray-900">{data.personalInfo.gender}</span></div>}
              {data.personalInfo.age && <div><span className="text-gray-500 mr-2">年龄：</span><span className="text-gray-900">{data.personalInfo.age}</span></div>}
              {data.personalInfo.maritalStatus && <div><span className="text-gray-500 mr-2">婚姻：</span><span className="text-gray-900">{data.personalInfo.maritalStatus}</span></div>}
              {data.personalInfo.politicalStatus && <div><span className="text-gray-500 mr-2">政治面貌：</span><span className="text-gray-900">{data.personalInfo.politicalStatus}</span></div>}
              {data.personalInfo.birthplace && <div><span className="text-gray-500 mr-2">籍贯：</span><span className="text-gray-900">{data.personalInfo.birthplace}</span></div>}
              {data.personalInfo.ethnicity && <div><span className="text-gray-500 mr-2">民族：</span><span className="text-gray-900">{data.personalInfo.ethnicity}</span></div>}
              {data.personalInfo.height && <div><span className="text-gray-500 mr-2">身高：</span><span className="text-gray-900">{data.personalInfo.height}</span></div>}
              {data.personalInfo.weight && <div><span className="text-gray-500 mr-2">体重：</span><span className="text-gray-900">{data.personalInfo.weight}</span></div>}
              {(data.personalInfo.customInfo || []).map((ci, idx) => (
                <div key={idx} className="col-span-2">
                  <span className="text-gray-500 mr-2">{ci.label}：</span><span className="text-gray-900 break-all">{ci.value}</span>
                </div>
              ))}
            </div>
            {data.personalInfo.avatarUrl && (
              <div className="md:col-span-3 flex md:justify-end">
                <img
                  src={data.personalInfo.avatarUrl}
                  alt="Profile"
                  className="w-28 h-28 rounded-md object-cover border-2 shadow-sm"
                  style={{ borderColor: data.themeConfig?.color || '#93c5fd' }}
                />
              </div>
            )}
          </div>
        </div>
        {data.sections.filter(s => s.isVisible).map(section => (
          <div key={section.id} className="mb-6">
            <SectionHeader title={getSectionTitle(section)} />
            <div className="pt-4 space-y-4">
              {section.type === ResumeSectionType.Skills
                ? (
                  <div className="space-y-2">
                    {section.items.map(it => (
                      <div key={it.id} className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap">
                        {it.description}
                      </div>
                    ))}
                  </div>
                )
                : section.items.map(item => (
                  <div key={item.id} className="border-b border-gray-200 pb-3">
                    <div className="flex justify-between items-baseline">
                      <h4 className="font-bold text-gray-900">{item.title}</h4>
                      {item.dateRange && <span className="text-sm text-gray-600">{item.dateRange}</span>}
                    </div>
                    {item.subtitle && (
                      <div className="mt-0.5 text-sm font-medium text-gray-700">
                        {item.subtitle}
                      </div>
                    )}
                    <div className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap mt-1" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                  </div>
                ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// 2. Modern Dark
const TemplateModern: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full grid grid-cols-12 h-full min-h-[1123px] bg-white ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto print:bg-transparent`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
        <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[70mm] -z-10" style={{ backgroundColor: '#0f172a' }}></div>
        <div className="col-span-4 text-white p-8" style={{ backgroundColor: '#0f172a' }}>
            <div className="mb-8 flex flex-col items-center md:items-start">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-slate-700 shadow-lg" style={{ borderRadius: '50%' }}/>
                 )}
                 <h1 className="text-2xl font-bold leading-tight break-words" style={{ color: data.themeConfig?.color || 'white' }}>{data.personalInfo.fullName}</h1>
                 <p className="text-slate-300 mt-1 font-medium">{data.personalInfo.jobTitle}</p>
            </div>

                 <div className="space-y-4 text-sm text-slate-300 mb-8 break-all">
                     <div className="block">{data.personalInfo.email}</div>
                     <div className="block">{data.personalInfo.phone}</div>
                     {/* website removed */}
                 {data.personalInfo.gender && <div className="block">{t('editor.fields.gender')}: {data.personalInfo.gender}</div>}
                 {data.personalInfo.age && <div className="block">{t('editor.fields.age')}: {data.personalInfo.age}</div>}
                 {data.personalInfo.maritalStatus && <div className="block">{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</div>}
                 {data.personalInfo.politicalStatus && <div className="block">{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</div>}
                 {data.personalInfo.birthplace && <div className="block">{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</div>}
                 {data.personalInfo.ethnicity && <div className="block">{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</div>}
                 {data.personalInfo.height && <div className="block">{t('editor.fields.height')}: {data.personalInfo.height}</div>}
                 {data.personalInfo.weight && <div className="block">{t('editor.fields.weight')}: {data.personalInfo.weight}</div>}
                 {(data.personalInfo.customInfo || []).map((ci, idx) => (
                   <div key={idx} className="block">{ci.label}: {ci.value}</div>
                 ))}
            </div>

            {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                <div key={section.id} className="mb-6">
                    <h3 className="text-white font-bold uppercase tracking-wider mb-4 text-sm border-b border-slate-700 pb-2">{getSectionTitle(section)}</h3>
                    <div className="flex flex-wrap gap-2">
                        {section.items.map(item => (
                             <span key={item.id} className="text-xs bg-slate-800 px-2 py-1 rounded text-slate-200 border border-slate-700">{item.description}</span>
                        ))}
                    </div>
                </div>
            ))}
        </div>
        
        <div className="col-span-8 p-8">
             {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id} className="mb-8">
                      <h3 className="text-slate-900 font-bold uppercase tracking-wider mb-4 text-sm border-b-2 inline-block pb-1" style={{ borderColor: data.themeConfig?.color || '#3b82f6' }}>{getSectionTitle(section)}</h3>
                      <div className="space-y-5">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-start">
                                      <div>
                                          {item.title && <h4 className="font-bold text-gray-800">{item.title}</h4>}
                                          {item.subtitle && <p className="font-medium text-sm" style={{ color: data.themeConfig?.color || '#2563eb' }}>{item.subtitle}</p>}
                                      </div>
                                      <div className="text-right shrink-0 ml-2">
                                           {item.dateRange && <p className="text-xs text-gray-500 font-medium bg-gray-100 px-2 py-1 rounded inline-block">{item.dateRange}</p>}
                                      </div>
                                  </div>
                                  <div className="mt-2 text-sm text-gray-600 leading-relaxed whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
)
};

// 3. Tech Minimalist
const TemplateMinimalist: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full p-8 md:p-14 bg-white h-full min-h-[1123px] text-gray-800 ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.4 }}>
        <header className="border-b-4 pb-6 mb-10 flex flex-col md:flex-row justify-between items-start md:items-end gap-6" style={{ borderColor: data.themeConfig?.color || 'black' }}>
            <div className="flex-1">
                <h1 className="text-5xl font-black tracking-tight uppercase mb-4 leading-none" style={{ color: data.themeConfig?.color || 'black' }}>{data.personalInfo.fullName}</h1>
                <div className="flex flex-wrap text-sm font-semibold gap-x-6 gap-y-2 text-gray-500 uppercase tracking-wide">
                    <span>{data.personalInfo.jobTitle}</span>
                    {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
                    {data.personalInfo.phone && <span>{data.personalInfo.phone}</span>}
                </div>
                {hasExtraPersonalInfo(data) && (
                  <div className="mt-2 flex flex-wrap text-xs gap-x-6 gap-y-2 text-gray-600">
                      {data.personalInfo.gender && <span>{t('editor.fields.gender')}: {data.personalInfo.gender}</span>}
                      {data.personalInfo.age && <span>{t('editor.fields.age')}: {data.personalInfo.age}</span>}
                      {data.personalInfo.maritalStatus && <span>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</span>}
                      {data.personalInfo.politicalStatus && <span>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</span>}
                      {data.personalInfo.birthplace && <span>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</span>}
                      {data.personalInfo.ethnicity && <span>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</span>}
                      {data.personalInfo.height && <span>{t('editor.fields.height')}: {data.personalInfo.height}</span>}
                      {data.personalInfo.weight && <span>{t('editor.fields.weight')}: {data.personalInfo.weight}</span>}
                      {(data.personalInfo.customInfo || []).map((ci, idx) => (
                        <span key={idx}>{ci.label}: {ci.value}</span>
                      ))}
                  </div>
                )}
            </div>
             {data.personalInfo.avatarUrl && (
                <img 
                    src={data.personalInfo.avatarUrl} 
                    alt="Profile" 
                    className="w-24 h-24 object-cover border border-gray-200 grayscale shadow-sm flex-shrink-0 self-center md:self-end"
                />
            )}
        </header>

        <div className="grid grid-cols-1 gap-10">
             {data.sections.filter(s => s.isVisible).map(section => (
                 <div key={section.id}>
                     <h3 className="text-xs font-bold uppercase tracking-[0.2em] text-gray-400 mb-6">{getSectionTitle(section)}</h3>
                     <div className="space-y-8">
                         {section.items.map(item => (
                             <div key={item.id} className="grid grid-cols-1 md:grid-cols-12 gap-4">
                                 <div className="md:col-span-3 text-xs font-bold text-gray-400 pt-1 uppercase tracking-wide">
                                    {item.dateRange}
                                 </div>
                                 <div className="md:col-span-9">
                                     <h4 className="font-bold text-gray-900 text-lg leading-none mb-1">{item.title}</h4>
                                     {item.subtitle && <p className="text-sm font-semibold text-gray-600 mb-3">{item.subtitle}</p>}
                                     <div className="text-sm leading-relaxed text-gray-700 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                 </div>
                             </div>
                         ))}
                     </div>
                 </div>
             ))}
        </div>
    </div>
);
};

// 4. Executive Serif
const TemplateExecutive: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full p-10 md:p-12 bg-white h-full min-h-[1123px] text-gray-900 ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.6 }}>
        <div className="text-center border-b pb-6 mb-8 flex flex-col items-center" style={{ borderColor: data.themeConfig?.color || '#111827' }}>
            {data.personalInfo.avatarUrl && (
                 <img 
                    src={data.personalInfo.avatarUrl} 
                    alt="Profile" 
                    className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-double border-gray-200 shadow-sm"
                    style={{ borderRadius: '50%' }}
                 />
            )}
            <h1 className="text-3xl font-bold uppercase mb-2 tracking-widest" style={{ color: data.themeConfig?.color }}>{data.personalInfo.fullName}</h1>
            <p className="italic text-lg text-gray-700 mb-3">{data.personalInfo.jobTitle}</p>
                <div className="text-sm text-gray-600 space-x-3 font-sans flex flex-wrap justify-center">
                 <span>{data.personalInfo.phone}</span> <span className="text-gray-300">•</span> <span>{data.personalInfo.email}</span>
            </div>
            {hasExtraPersonalInfo(data) && (
              <div className="mt-2 text-xs text-gray-600 space-x-3 font-sans flex flex-wrap justify-center">
                  {data.personalInfo.gender && <span>{t('editor.fields.gender')}: {data.personalInfo.gender}</span>}
                  {data.personalInfo.age && <span>{t('editor.fields.age')}: {data.personalInfo.age}</span>}
                  {data.personalInfo.maritalStatus && <span>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</span>}
                  {data.personalInfo.politicalStatus && <span>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</span>}
                  {data.personalInfo.birthplace && <span>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</span>}
                  {data.personalInfo.ethnicity && <span>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</span>}
                  {data.personalInfo.height && <span>{t('editor.fields.height')}: {data.personalInfo.height}</span>}
                  {data.personalInfo.weight && <span>{t('editor.fields.weight')}: {data.personalInfo.weight}</span>}
                  {(data.personalInfo.customInfo || []).map((ci, idx) => (
                    <span key={idx}>{ci.label}: {ci.value}</span>
                  ))}
              </div>
            )}
        </div>

        {data.sections.filter(s => s.isVisible).map(section => (
            <div key={section.id} className="mb-6">
                <h3 className="text-md font-bold uppercase border-b border-gray-400 mb-4 pb-1 flex justify-between items-end">
                    <span>{getSectionTitle(section)}</span>
                </h3>
                <div className="space-y-5">
                    {section.items.map(item => (
                        <div key={item.id}>
                            <div className="flex justify-between items-baseline mb-1">
                                <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                <span className="text-sm font-bold font-sans text-gray-700">{item.dateRange}</span>
                            </div>
                            {item.subtitle && (
                                <div className="mb-2">
                                    <span className="italic text-gray-800">{item.subtitle}</span>
                                </div>
                            )}
                             <div className="text-sm leading-normal text-gray-800 whitespace-pre-wrap" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                        </div>
                    ))}
                </div>
            </div>
        ))}
    </div>
);
};

// 5. Creative Bold
const TemplateBold: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
     const { t } = useLanguage();
     const getSectionTitle = useSectionTitle();
     return (
     <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto flex flex-col`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
        <div className="text-white p-10 print:text-white flex flex-col md:flex-row justify-between items-center gap-6" style={{ backgroundColor: data.themeConfig?.color || '#1d4ed8' }}>
             <div className="order-2 md:order-1 flex-1">
                <h1 className="text-5xl font-extrabold mb-2 tracking-tight">{data.personalInfo.fullName}</h1>
                <p className="text-white/80 text-2xl font-light">{data.personalInfo.jobTitle}</p>
                <div className="mt-6 flex flex-wrap gap-x-6 gap-y-2 text-sm text-white/90 font-medium">
                    {data.personalInfo.email && <div className="flex items-center gap-2">{data.personalInfo.email}</div>}
                    {data.personalInfo.phone && <div className="flex items-center gap-2">• {data.personalInfo.phone}</div>}
                </div>
                <div className="mt-2 flex flex-wrap gap-x-6 gap-y-1 text-xs text-white/80">
                    {hasExtraPersonalInfo(data) && (
                      <>
                        {data.personalInfo.gender && <div>{t('editor.fields.gender')}: {data.personalInfo.gender}</div>}
                        {data.personalInfo.age && <div>{t('editor.fields.age')}: {data.personalInfo.age}</div>}
                        {data.personalInfo.maritalStatus && <div>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</div>}
                        {data.personalInfo.politicalStatus && <div>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</div>}
                        {data.personalInfo.birthplace && <div>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</div>}
                        {data.personalInfo.ethnicity && <div>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</div>}
                        {data.personalInfo.height && <div>{t('editor.fields.height')}: {data.personalInfo.height}</div>}
                        {data.personalInfo.weight && <div>{t('editor.fields.weight')}: {data.personalInfo.weight}</div>}
                        {(data.personalInfo.customInfo || []).map((ci, idx) => (
                          <div key={idx}>{ci.label}: {ci.value}</div>
                        ))}
                      </>
                    )}
                </div>
             </div>
             {data.personalInfo.avatarUrl && (
                 <img src={data.personalInfo.avatarUrl} alt={t('a11y.avatarAlt')} className="order-1 md:order-2 w-32 h-32 rounded-full border-4 border-white object-cover shadow-xl flex-shrink-0" style={{ borderRadius: '50%' }} />
             )}
        </div>

        <div className="p-10 grid grid-cols-1 md:grid-cols-12 gap-8 flex-grow">
            <div className="md:col-span-8 pr-4">
                 {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-10">
                         <h3 className="font-bold text-lg uppercase mb-6 flex items-center tracking-wider" style={{ color: data.themeConfig?.color || '#1d4ed8' }}>
                             <span className="w-1.5 h-6 mr-3 rounded-sm" style={{ backgroundColor: data.themeConfig?.color || '#1d4ed8' }}></span>
                             {getSectionTitle(section)}
                         </h3>
                          <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id} className="relative">
                                    <div className="flex justify-between items-center mb-1">
                                        <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                        {item.dateRange && <span className="text-xs font-bold bg-gray-100 text-gray-700 px-3 py-1 rounded-full whitespace-nowrap ml-4">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <p className="text-sm text-gray-600 font-medium mb-3">{item.subtitle}</p>}
                                    <div className="text-sm text-gray-600 whitespace-pre-wrap leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                </div>
                            ))}
                        </div>
                     </div>
                 ))}
            </div>
            
            <div className="md:col-span-4 bg-gray-50 p-6 rounded-2xl h-fit border border-gray-100">
                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-8">
                         <h3 className="text-gray-900 font-bold uppercase mb-4 text-sm tracking-wider">{getSectionTitle(section)}</h3>
                         <div className="flex flex-wrap gap-2">
                             {section.items.map(item => (
                                 <span key={item.id} className="bg-white border border-gray-200 px-3 py-1.5 rounded-lg text-sm shadow-sm font-semibold text-gray-700">{item.description}</span>
                             ))}
                         </div>
                     </div>
                 ))}
            </div>
        </div>
     </div>
     );
};

// 6. Elegant Teal
const TemplateElegant: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
    const getSectionTitle = useSectionTitle();
    const { t } = useLanguage();
    return (
    <div className={`w-full grid grid-cols-12 h-full min-h-[1123px] bg-white ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto print:bg-transparent`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.6 }}>
        <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[70mm] -z-10" style={{ backgroundColor: data.themeConfig?.color || '#115e59' }}></div>
        <div className="col-span-4 text-white p-8 flex flex-col" style={{ backgroundColor: data.themeConfig?.color || '#115e59' }}>
            <div className="mb-10 text-center">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-white/20 mx-auto shadow-md" style={{ borderRadius: '50%' }}/>
                 )}
                <h1 className="text-2xl font-serif font-bold leading-tight mb-2" style={{ fontFamily: styles.fontFamily }}>{data.personalInfo.fullName}</h1>
                 <p className="text-white/70 uppercase tracking-widest text-xs font-semibold">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-6 text-sm text-white/90 flex-grow">
                 <div>
                    <span className="block text-white/60 text-xs font-bold uppercase mb-1">Contact</span>
                    <div className="space-y-1">
                        <div className="break-all">{data.personalInfo.email}</div>
                        <div>{data.personalInfo.phone}</div>
                        {/* website removed */}
                        {data.personalInfo.gender && <div>{t('editor.fields.gender')}: {data.personalInfo.gender}</div>}
                        {data.personalInfo.age && <div>{t('editor.fields.age')}: {data.personalInfo.age}</div>}
                        {data.personalInfo.maritalStatus && <div>{t('editor.fields.maritalStatus')}: {data.personalInfo.maritalStatus}</div>}
                        {data.personalInfo.politicalStatus && <div>{t('editor.fields.politicalStatus')}: {data.personalInfo.politicalStatus}</div>}
                        {data.personalInfo.birthplace && <div>{t('editor.fields.birthplace')}: {data.personalInfo.birthplace}</div>}
                        {data.personalInfo.ethnicity && <div>{t('editor.fields.ethnicity')}: {data.personalInfo.ethnicity}</div>}
                        {data.personalInfo.height && <div>{t('editor.fields.height')}: {data.personalInfo.height}</div>}
                        {data.personalInfo.weight && <div>{t('editor.fields.weight')}: {data.personalInfo.weight}</div>}
                        {(data.personalInfo.customInfo || []).map((ci, idx) => (
                          <div key={idx}>{ci.label}: {ci.value}</div>
                        ))}
                    </div>
                 </div>

                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <span className="block text-white/60 text-xs font-bold uppercase mb-2">{getSectionTitle(section)}</span>
                        <ul className="list-disc list-inside space-y-1 text-white/80">
                            {section.items.map(item => (
                                <li key={item.id}>{item.description}</li>
                            ))}
                        </ul>
                    </div>
                 ))}
            </div>
        </div>
        
        <div className="col-span-8 p-10 bg-white">
             {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id} className="mb-10">
                      <h3 className="font-serif font-bold uppercase tracking-widest text-lg border-b pb-2 mb-6" style={{ borderColor: '#f3f4f6', color: data.themeConfig?.color || '#134e4a' }}>{getSectionTitle(section)}</h3>
                      <div className="space-y-6">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-baseline mb-1">
                                      <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                      <span className="text-sm font-medium" style={{ color: data.themeConfig?.color || '#115e59' }}>{item.dateRange}</span>
                                  </div>
                                  {item.subtitle && <p className="text-gray-600 italic mb-2">{item.subtitle}</p>}
                                  <div className="text-sm text-gray-600 leading-relaxed whitespace-pre-wrap font-light" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
);
};

export interface ArtboardProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

export const ResumeArtboard: React.FC<ArtboardProps> = ({ data, scale = 1, disableShadow = false, className = '', style = {} }) => {
  const styles = getThemeStyles(data.themeConfig);

  const containerStyle = {
    transform: `scale(${scale})`,
    transformOrigin: 'top left',
    ...style,
  };

  const renderTemplate = () => {
      switch (data.templateId) {
          case 't1': return <TemplateClassic data={data} styles={styles} disableShadow={disableShadow} />;
          case 't2': return <TemplateModern data={data} styles={styles} disableShadow={disableShadow} />;
          case 't3': return <TemplateMinimalist data={data} styles={styles} disableShadow={disableShadow} />;
          case 't4': return <TemplateExecutive data={data} styles={styles} disableShadow={disableShadow} />;
          case 't5': return <TemplateBold data={data} styles={styles} disableShadow={disableShadow} />;
          case 't6': return <TemplateElegant data={data} styles={styles} disableShadow={disableShadow} />;
          case 't7': return <TemplateCNBlue data={data} styles={styles} disableShadow={disableShadow} />;
          default: return <TemplateClassic data={data} styles={styles} />;
      }
  };

  return (
      <div 
        id="resume-export-root"
        className={`w-[210mm] min-h-[297mm] print:w-full print:min-h-0 print:transform-none bg-white mx-auto ${disableShadow ? 'shadow-none' : 'shadow-md'} print:shadow-none ${className}`}
        style={containerStyle}
      >
        {renderTemplate()}
      </div>
  );
};

export const ResumePreview: React.FC<PreviewProps> = ({ data, scale = 1, disableShadow = false }) => {
  return (
    <div className={`w-full flex justify-center bg-white p-8 overflow-auto print:p-0 print:bg-white h-full scrollbar-thin scrollbar-thumb-gray-300 ${disableShadow ? 'shadow-none' : ''} print:shadow-none`}>
      <ResumeArtboard data={data} scale={scale} disableShadow={disableShadow} />
    </div>
  );
};
