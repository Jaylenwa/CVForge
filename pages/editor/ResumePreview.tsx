import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';

interface PreviewProps {
  data: ResumeData;
  scale?: number;
}

// Basic Template Renderer
const TemplateClassic: React.FC<{ data: ResumeData }> = ({ data }) => (
  <div className="p-8 md:p-12 bg-white text-gray-900 h-full min-h-[1123px] font-serif shadow-lg print:shadow-none print:p-0">
    {/* Header */}
    <div className="border-b-2 border-gray-800 pb-6 mb-6">
      <h1 className="text-4xl font-bold uppercase tracking-wider">{data.personalInfo.fullName}</h1>
      <p className="text-xl mt-2 text-gray-600">{data.personalInfo.jobTitle}</p>
      <div className="mt-4 flex flex-wrap gap-4 text-sm text-gray-600">
        {data.personalInfo.email && <span>{data.personalInfo.email}</span>}
        {data.personalInfo.phone && <span>• {data.personalInfo.phone}</span>}
        {data.personalInfo.address && <span>• {data.personalInfo.address}</span>}
        {data.personalInfo.website && <span>• {data.personalInfo.website}</span>}
      </div>
    </div>

    {/* Sections */}
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
                
                {/* Description Rendering - simplistic HTML rendering or text */}
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

// Modern Template Renderer
const TemplateModern: React.FC<{ data: ResumeData }> = ({ data }) => (
    <div className="grid grid-cols-12 h-full min-h-[1123px] bg-white shadow-lg print:shadow-none font-sans">
        {/* Sidebar */}
        <div className="col-span-4 bg-slate-900 text-white p-8">
            <div className="mb-8">
                 {data.personalInfo.avatarUrl && (
                     <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-24 h-24 rounded-full mb-4 object-cover border-2 border-white"/>
                 )}
                 <h1 className="text-2xl font-bold leading-tight">{data.personalInfo.fullName}</h1>
                 <p className="text-blue-300 mt-1">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-4 text-sm text-slate-300 mb-8">
                 <div className="block">{data.personalInfo.email}</div>
                 <div className="block">{data.personalInfo.phone}</div>
                 <div className="block">{data.personalInfo.address}</div>
                 <div className="block truncate">{data.personalInfo.website}</div>
            </div>

            {/* Skills in Sidebar for Modern */}
            {data.sections.filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                <div key={section.id} className="mb-6">
                    <h3 className="text-white font-bold uppercase tracking-wider mb-4 text-sm border-b border-slate-700 pb-2">{section.title}</h3>
                    <div className="flex flex-wrap gap-2">
                        {section.items.map(item => (
                             <span key={item.id} className="text-xs bg-slate-800 px-2 py-1 rounded text-slate-200">{item.description}</span>
                        ))}
                    </div>
                </div>
            ))}
        </div>
        
        {/* Main Content */}
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
                                      <div className="text-right">
                                           {item.dateRange && <p className="text-xs text-gray-500 font-medium bg-gray-100 px-2 py-1 rounded inline-block">{item.dateRange}</p>}
                                      </div>
                                  </div>
                                  <div className="mt-2 text-sm text-gray-600 leading-relaxed">
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

export const ResumePreview: React.FC<PreviewProps> = ({ data, scale = 1 }) => {
  const style = {
    transform: `scale(${scale})`,
    transformOrigin: 'top center',
  };

  return (
    <div className="w-full flex justify-center bg-gray-100 p-8 overflow-auto print:p-0 print:bg-white">
      <div 
        className="w-[210mm] min-h-[297mm] print:w-full print:transform-none bg-white mx-auto transition-transform duration-200"
        style={style}
      >
        {data.templateId === 't2' ? <TemplateModern data={data} /> : <TemplateClassic data={data} />}
      </div>
    </div>
  );
};
