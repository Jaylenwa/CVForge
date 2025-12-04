import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';

interface PreviewProps {
  data: ResumeData;
  scale?: number;
}

// 1. Classic Professional
const TemplateClassic: React.FC<{ data: ResumeData }> = ({ data }) => (
  <div className="p-8 md:p-12 bg-white text-gray-900 h-full min-h-[1123px] font-serif shadow-lg print:shadow-none print:p-0">
    <div className="border-b-2 border-gray-800 pb-6 mb-6 flex flex-col md:flex-row items-center md:items-start gap-6">
      <div className="flex-1 text-center md:text-left order-2 md:order-1">
          <h1 className="text-4xl font-bold uppercase tracking-wider">{data.personalInfo.fullName}</h1>
          <p className="text-xl mt-2 text-gray-600">{data.personalInfo.jobTitle}</p>
          <div className="mt-4 flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-600 font-sans">
            {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
            {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
            {data.personalInfo.address && <span>• {data.personalInfo.address}</span>}
            {data.personalInfo.website && <span>• {data.personalInfo.website}</span>}
          </div>
      </div>
      {data.personalInfo.avatarUrl && (
          <div className="order-1 md:order-2 flex-shrink-0">
             <img 
               src={data.personalInfo.avatarUrl} 
               alt="Profile" 
               className="w-32 h-32 rounded-lg object-cover border-2 border-gray-200 shadow-sm"
             />
          </div>
      )}
    </div>

    <div className="space-y-6">
      {data.sections.filter(s => s.isVisible).map(section => (
        <div key={section.id}>
          <h3 className="text-lg font-bold uppercase border-b border-gray-300 mb-4 pb-1 tracking-wide text-gray-800">
            {section.title}
          </h3>
          <div className="space-y-4">
            {section.items.map(item => (
              <div key={item.id} className="relative">
                {section.type !== ResumeSectionType.Skills && (
                  <div className="flex justify-between items-baseline mb-1">
                    <h4 className="font-bold text-gray-900">{item.title}</h4>
                    {item.dateRange && <span className="text-sm italic text-gray-600 font-sans">{item.dateRange}</span>}
                  </div>
                )}
                {item.subtitle && (
                   <div className="flex justify-between items-baseline mb-1">
                    <span className="font-semibold text-gray-700">{item.subtitle}</span>
                     {item.location && <span className="text-sm text-gray-500 font-sans">{item.location}</span>}
                   </div>
                )}
                <div className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap">
                    {item.description}
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  </div>
);

// 2. Modern Dark
const TemplateModern: React.FC<{ data: ResumeData }> = ({ data }) => (
    <div className="grid grid-cols-12 h-full min-h-[1123px] bg-white shadow-lg print:shadow-none font-sans">
        <div className="col-span-4 bg-slate-900 text-white p-8">
            <div className="mb-8 flex flex-col items-center md:items-start">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-slate-700 shadow-lg"/>
                 )}
                 <h1 className="text-2xl font-bold leading-tight break-words">{data.personalInfo.fullName}</h1>
                 <p className="text-blue-300 mt-1 font-medium">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-4 text-sm text-slate-300 mb-8 break-all">
                 <div className="block">{data.personalInfo.email}</div>
                 <div className="block">{data.personalInfo.phone}</div>
                 <div className="block">{data.personalInfo.address}</div>
                 <div className="block">{data.personalInfo.website}</div>
            </div>

            {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                <div key={section.id} className="mb-6">
                    <h3 className="text-white font-bold uppercase tracking-wider mb-4 text-sm border-b border-slate-700 pb-2">{section.title}</h3>
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
                      <h3 className="text-slate-900 font-bold uppercase tracking-wider mb-4 text-sm border-b-2 border-blue-500 inline-block pb-1">{section.title}</h3>
                      <div className="space-y-5">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-start">
                                      <div>
                                          {item.title && <h4 className="font-bold text-gray-800">{item.title}</h4>}
                                          {item.subtitle && <p className="text-blue-600 font-medium text-sm">{item.subtitle}</p>}
                                      </div>
                                      <div className="text-right shrink-0 ml-2">
                                           {item.dateRange && <p className="text-xs text-gray-500 font-medium bg-gray-100 px-2 py-1 rounded inline-block">{item.dateRange}</p>}
                                      </div>
                                  </div>
                                  <div className="mt-2 text-sm text-gray-600 leading-relaxed whitespace-pre-wrap">
                                      {item.description}
                                  </div>
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
)

// 3. Tech Minimalist
const TemplateMinimalist: React.FC<{ data: ResumeData }> = ({ data }) => (
    <div className="p-8 md:p-14 bg-white h-full min-h-[1123px] font-sans text-gray-800 shadow-lg print:shadow-none">
        <header className="border-b-4 border-black pb-6 mb-10 flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
            <div className="flex-1">
                <h1 className="text-5xl font-black tracking-tight uppercase mb-4 leading-none">{data.personalInfo.fullName}</h1>
                <div className="flex flex-wrap text-sm font-semibold gap-x-6 gap-y-2 text-gray-500 uppercase tracking-wide">
                    <span>{data.personalInfo.jobTitle}</span>
                    {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
                    {data.personalInfo.phone && <span>{data.personalInfo.phone}</span>}
                    {data.personalInfo.website && <span>{data.personalInfo.website}</span>}
                </div>
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
                     <h3 className="text-xs font-bold uppercase tracking-[0.2em] text-gray-400 mb-6">{section.title}</h3>
                     <div className="space-y-8">
                         {section.items.map(item => (
                             <div key={item.id} className="grid grid-cols-1 md:grid-cols-12 gap-4">
                                 <div className="md:col-span-3 text-xs font-bold text-gray-400 pt-1 uppercase tracking-wide">
                                    {item.dateRange}
                                 </div>
                                 <div className="md:col-span-9">
                                     <h4 className="font-bold text-gray-900 text-lg leading-none mb-1">{item.title}</h4>
                                     {item.subtitle && <p className="text-sm font-semibold text-gray-600 mb-3">{item.subtitle} {item.location && <span className="font-normal text-gray-400">• {item.location}</span>}</p>}
                                     <div className="text-sm leading-relaxed text-gray-700 whitespace-pre-wrap">{item.description}</div>
                                 </div>
                             </div>
                         ))}
                     </div>
                 </div>
             ))}
        </div>
    </div>
);

// 4. Executive Serif
const TemplateExecutive: React.FC<{ data: ResumeData }> = ({ data }) => (
    <div className="p-10 md:p-12 bg-white h-full min-h-[1123px] font-serif text-gray-900 shadow-lg print:shadow-none">
        <div className="text-center border-b border-gray-900 pb-6 mb-8 flex flex-col items-center">
            {data.personalInfo.avatarUrl && (
                 <img 
                    src={data.personalInfo.avatarUrl} 
                    alt="Profile" 
                    className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-double border-gray-200 shadow-sm"
                 />
            )}
            <h1 className="text-3xl font-bold uppercase mb-2 tracking-widest">{data.personalInfo.fullName}</h1>
            <p className="italic text-lg text-gray-700 mb-3">{data.personalInfo.jobTitle}</p>
            <div className="text-sm text-gray-600 space-x-3 font-sans flex flex-wrap justify-center">
                 <span>{data.personalInfo.phone}</span> <span className="text-gray-300">•</span> <span>{data.personalInfo.email}</span> <span className="text-gray-300">•</span> <span>{data.personalInfo.address}</span>
            </div>
        </div>

        {data.sections.filter(s => s.isVisible).map(section => (
            <div key={section.id} className="mb-6">
                <h3 className="text-md font-bold uppercase border-b border-gray-400 mb-4 pb-1 flex justify-between items-end">
                    <span>{section.title}</span>
                </h3>
                <div className="space-y-5">
                    {section.items.map(item => (
                        <div key={item.id}>
                            <div className="flex justify-between items-baseline mb-1">
                                <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                <span className="text-sm font-bold font-sans text-gray-700">{item.dateRange}</span>
                            </div>
                            {item.subtitle && (
                                <div className="flex justify-between items-baseline mb-2">
                                    <span className="italic text-gray-800">{item.subtitle}</span>
                                    <span className="text-sm italic text-gray-500">{item.location}</span>
                                </div>
                            )}
                             <div className="text-sm leading-normal text-gray-800 whitespace-pre-wrap">{item.description}</div>
                        </div>
                    ))}
                </div>
            </div>
        ))}
    </div>
);

// 5. Creative Bold
const TemplateBold: React.FC<{ data: ResumeData }> = ({ data }) => (
     <div className="bg-white h-full min-h-[1123px] shadow-lg print:shadow-none font-sans flex flex-col">
        <div className="bg-blue-700 text-white p-10 print:bg-blue-700 print:text-white flex flex-col md:flex-row justify-between items-center gap-6">
             <div className="order-2 md:order-1 flex-1">
                <h1 className="text-5xl font-extrabold mb-2 tracking-tight">{data.personalInfo.fullName}</h1>
                <p className="text-blue-200 text-2xl font-light">{data.personalInfo.jobTitle}</p>
                <div className="mt-6 flex flex-wrap gap-x-6 gap-y-2 text-sm text-blue-100 opacity-90 font-medium">
                    {data.personalInfo.email && <div className="flex items-center gap-2">{data.personalInfo.email}</div>}
                    {data.personalInfo.phone && <div className="flex items-center gap-2">• {data.personalInfo.phone}</div>}
                    {data.personalInfo.website && <div className="flex items-center gap-2">• {data.personalInfo.website}</div>}
                    {data.personalInfo.address && <div className="flex items-center gap-2">• {data.personalInfo.address}</div>}
                </div>
             </div>
             {data.personalInfo.avatarUrl && (
                 <img src={data.personalInfo.avatarUrl} alt="Avatar" className="order-1 md:order-2 w-32 h-32 rounded-full border-4 border-white object-cover shadow-xl flex-shrink-0" />
             )}
        </div>

        <div className="p-10 grid grid-cols-1 md:grid-cols-12 gap-8 flex-grow">
            <div className="md:col-span-8 pr-4">
                 {data.sections.filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-10">
                         <h3 className="text-blue-700 font-bold text-lg uppercase mb-6 flex items-center tracking-wider">
                             <span className="w-1.5 h-6 bg-blue-700 mr-3 rounded-sm"></span>
                             {section.title}
                         </h3>
                          <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id} className="relative">
                                    <div className="flex justify-between items-center mb-1">
                                        <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                        {item.dateRange && <span className="text-xs font-bold bg-blue-50 text-blue-700 px-3 py-1 rounded-full whitespace-nowrap ml-4">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <p className="text-sm text-gray-600 font-medium mb-3">{item.subtitle} {item.location && <span className="text-gray-400 font-normal">| {item.location}</span>}</p>}
                                    <div className="text-sm text-gray-600 whitespace-pre-wrap leading-relaxed">{item.description}</div>
                                </div>
                            ))}
                        </div>
                     </div>
                 ))}
            </div>
            
            <div className="md:col-span-4 bg-gray-50 p-6 rounded-2xl h-fit border border-gray-100">
                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                     <div key={section.id} className="mb-8">
                         <h3 className="text-gray-900 font-bold uppercase mb-4 text-sm tracking-wider">{section.title}</h3>
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

// 6. Elegant Teal
const TemplateElegant: React.FC<{ data: ResumeData }> = ({ data }) => (
    <div className="grid grid-cols-12 h-full min-h-[1123px] bg-white shadow-lg print:shadow-none font-sans">
        <div className="col-span-4 bg-teal-800 text-white p-8 flex flex-col">
            <div className="mb-10 text-center">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-32 h-32 rounded-full mb-6 object-cover border-4 border-teal-600 mx-auto shadow-md"/>
                 )}
                 <h1 className="text-2xl font-serif font-bold leading-tight mb-2">{data.personalInfo.fullName}</h1>
                 <p className="text-teal-200 uppercase tracking-widest text-xs font-semibold">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-6 text-sm text-teal-50 flex-grow">
                 <div>
                    <span className="block text-teal-400 text-xs font-bold uppercase mb-1">Contact</span>
                    <div className="space-y-1">
                        <div className="break-all">{data.personalInfo.email}</div>
                        <div>{data.personalInfo.phone}</div>
                        <div>{data.personalInfo.address}</div>
                        <div className="break-all">{data.personalInfo.website}</div>
                    </div>
                 </div>

                 {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <span className="block text-teal-400 text-xs font-bold uppercase mb-2">{section.title}</span>
                        <ul className="list-disc list-inside space-y-1 text-teal-100">
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
                      <h3 className="text-teal-900 font-serif font-bold uppercase tracking-widest text-lg border-b border-teal-100 pb-2 mb-6">{section.title}</h3>
                      <div className="space-y-6">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-baseline mb-1">
                                      <h4 className="font-bold text-gray-900 text-lg">{item.title}</h4>
                                      <span className="text-sm text-teal-700 font-medium">{item.dateRange}</span>
                                  </div>
                                  {item.subtitle && <p className="text-gray-600 italic mb-2">{item.subtitle} {item.location && <span className="not-italic text-sm text-gray-400">| {item.location}</span>}</p>}
                                  <div className="text-sm text-gray-600 leading-relaxed whitespace-pre-wrap font-light">
                                      {item.description}
                                  </div>
                              </div>
                          ))}
                      </div>
                 </div>
             ))}
        </div>
    </div>
);

export const ResumePreview: React.FC<PreviewProps> = ({ data, scale = 1 }) => {
  const style = {
    transform: `scale(${scale})`,
    transformOrigin: 'top center',
  };

  const renderTemplate = () => {
      switch (data.templateId) {
          case 't1': return <TemplateClassic data={data} />;
          case 't2': return <TemplateModern data={data} />;
          case 't3': return <TemplateMinimalist data={data} />;
          case 't4': return <TemplateExecutive data={data} />;
          case 't5': return <TemplateBold data={data} />;
          case 't6': return <TemplateElegant data={data} />;
          default: return <TemplateClassic data={data} />;
      }
  };

  return (
    <div className="w-full flex justify-center bg-gray-100 p-8 overflow-auto print:p-0 print:bg-white h-full scrollbar-thin scrollbar-thumb-gray-300">
      <div 
        className="w-[210mm] min-h-[297mm] print:w-full print:min-h-0 print:transform-none bg-white mx-auto transition-transform duration-200 shadow-2xl print:shadow-none"
        style={style}
      >
        {renderTemplate()}
      </div>
    </div>
  );
};