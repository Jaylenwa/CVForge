import React, { useState, useRef } from 'react';
import { Trash2, Plus, GripVertical, Sparkles, ChevronDown, ChevronUp, Upload, X, Image as ImageIcon } from 'lucide-react';
import { ResumeData, ResumeSection, ResumeItem, ResumeSectionType } from '../../types';
import { Button } from '../../components/ui/Button';
import { polishText, generateSummary } from '../../services/geminiService';

interface EditorFormProps {
  data: ResumeData;
  onChange: (data: ResumeData) => void;
}

export const EditorForm: React.FC<EditorFormProps> = ({ data, onChange }) => {
  const [activeSection, setActiveSection] = useState<string | null>('personal');
  const [isAiLoading, setIsAiLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const updatePersonalInfo = (field: string, value: string) => {
    onChange({
      ...data,
      personalInfo: { ...data.personalInfo, [field]: value }
    });
  };

  const handleAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
          alert("Image size too large. Please upload an image smaller than 2MB.");
          return;
      }
      
      const reader = new FileReader();
      reader.onloadend = () => {
        updatePersonalInfo('avatarUrl', reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const updateSection = (sectionId: string, updates: Partial<ResumeSection>) => {
    onChange({
      ...data,
      sections: data.sections.map(s => s.id === sectionId ? { ...s, ...updates } : s)
    });
  };

  const updateItem = (sectionId: string, itemId: string, field: keyof ResumeItem, value: string) => {
    onChange({
      ...data,
      sections: data.sections.map(s => {
        if (s.id !== sectionId) return s;
        return {
          ...s,
          items: s.items.map(item => item.id === itemId ? { ...item, [field]: value } : item)
        };
      })
    });
  };

  const addItem = (sectionId: string) => {
    const newItem: ResumeItem = {
      id: Math.random().toString(36).substr(2, 9),
      title: '',
      description: ''
    };
    onChange({
      ...data,
      sections: data.sections.map(s => 
        s.id === sectionId ? { ...s, items: [...s.items, newItem] } : s
      )
    });
  };

  const removeItem = (sectionId: string, itemId: string) => {
    onChange({
      ...data,
      sections: data.sections.map(s => 
        s.id === sectionId ? { ...s, items: s.items.filter(i => i.id !== itemId) } : s
      )
    });
  };

  const handleAiPolish = async (sectionId: string, itemId: string, text: string) => {
    if (!text) return;
    setIsAiLoading(true);
    const improved = await polishText(text);
    updateItem(sectionId, itemId, 'description', improved);
    setIsAiLoading(false);
  };

  const handleAiSummary = async () => {
     setIsAiLoading(true);
     // Find relevant skills for context
     const skillsSection = data.sections.find(s => s.type === ResumeSectionType.Skills);
     const skills = skillsSection?.items.map(i => i.description).join(', ') || 'General';
     
     const summary = await generateSummary(data.personalInfo.jobTitle, skills);
     
     // Update summary section
     const summarySection = data.sections.find(s => s.type === ResumeSectionType.Summary);
     if (summarySection && summarySection.items.length > 0) {
         updateItem(summarySection.id, summarySection.items[0].id, 'description', summary);
     }
     setIsAiLoading(false);
  }

  return (
    <div className="h-full overflow-y-auto bg-white border-r border-gray-200">
      
      {/* Personal Info */}
      <div className="border-b border-gray-200">
        <button 
          className="w-full px-6 py-4 flex items-center justify-between font-semibold hover:bg-gray-50"
          onClick={() => setActiveSection(activeSection === 'personal' ? null : 'personal')}
        >
            <span>Personal Information</span>
            {activeSection === 'personal' ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
        </button>
        
        {activeSection === 'personal' && (
          <div className="px-6 pb-6 space-y-5 animate-fadeIn">
             
             {/* Avatar Upload */}
             <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                 <label className="block text-sm font-medium text-gray-700 mb-3 flex items-center">
                    <ImageIcon size={16} className="mr-2 text-gray-500" /> 
                    Profile Photo
                 </label>
                 <div className="flex items-center space-x-4">
                    {data.personalInfo.avatarUrl ? (
                        <div className="relative group shrink-0">
                            <img 
                                src={data.personalInfo.avatarUrl} 
                                alt="Avatar" 
                                className="w-16 h-16 rounded-full object-cover border-2 border-white shadow-sm"
                            />
                            <button
                                onClick={() => updatePersonalInfo('avatarUrl', '')}
                                className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full p-1 shadow-sm hover:bg-red-600 transition-colors"
                                title="Remove photo"
                            >
                                <X size={10} />
                            </button>
                        </div>
                    ) : (
                        <div className="w-16 h-16 shrink-0 rounded-full bg-gray-100 flex items-center justify-center text-gray-400 border-2 border-gray-200 border-dashed">
                            <Upload size={20} />
                        </div>
                    )}
                    
                    <div className="flex-1">
                        <input 
                            ref={fileInputRef}
                            type="file" 
                            className="hidden" 
                            accept="image/png, image/jpeg, image/jpg" 
                            onChange={handleAvatarUpload} 
                        />
                        <Button 
                            variant="outline" 
                            size="sm"
                            onClick={() => fileInputRef.current?.click()}
                            className="w-full"
                        >
                            {data.personalInfo.avatarUrl ? 'Change Photo' : 'Upload Photo'}
                        </Button>
                        <p className="mt-2 text-xs text-gray-500">
                            Square image recommended (1:1).<br/>
                            Max size 2MB (JPG/PNG).
                        </p>
                    </div>
                 </div>
             </div>

            <div>
              <label className="block text-sm font-medium text-gray-700">Full Name</label>
              <input 
                type="text" 
                value={data.personalInfo.fullName}
                onChange={e => updatePersonalInfo('fullName', e.target.value)}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Job Title</label>
              <input 
                type="text" 
                value={data.personalInfo.jobTitle}
                onChange={e => updatePersonalInfo('jobTitle', e.target.value)}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
             <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700">Email</label>
                    <input type="email" value={data.personalInfo.email} onChange={e => updatePersonalInfo('email', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700">Phone</label>
                    <input type="tel" value={data.personalInfo.phone} onChange={e => updatePersonalInfo('phone', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
                </div>
             </div>
             <div>
                <label className="block text-sm font-medium text-gray-700">Address / Location</label>
                <input type="text" value={data.personalInfo.address} onChange={e => updatePersonalInfo('address', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
             </div>
             <div>
                <label className="block text-sm font-medium text-gray-700">Website / LinkedIn</label>
                <input type="text" value={data.personalInfo.website} onChange={e => updatePersonalInfo('website', e.target.value)} className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 sm:text-sm"/>
             </div>
          </div>
        )}
      </div>

      {/* Dynamic Sections */}
      {data.sections.map((section) => (
        <div key={section.id} className="border-b border-gray-200">
           <div className="flex items-center justify-between px-6 py-4 hover:bg-gray-50">
             <button 
                className="flex-grow text-left font-semibold flex items-center"
                onClick={() => setActiveSection(activeSection === section.id ? null : section.id)}
             >
                 <span className="mr-2 cursor-grab active:cursor-grabbing"><GripVertical size={16} className="text-gray-400"/></span>
                 {section.title}
             </button>
             <div className="flex items-center space-x-2">
                 {/* Toggle Visibility */}
                 <input 
                    type="checkbox" 
                    checked={section.isVisible} 
                    onChange={(e) => updateSection(section.id, { isVisible: e.target.checked })}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                 />
                  {activeSection === section.id ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
             </div>
           </div>

           {activeSection === section.id && (
             <div className="px-6 pb-6 space-y-6 bg-gray-50/50">
                {/* AI Summary Generator specific to Summary Section */}
                {section.type === ResumeSectionType.Summary && (
                    <div className="flex justify-end">
                         <Button 
                            size="sm" 
                            variant="secondary" 
                            icon={<Sparkles size={14}/>} 
                            onClick={handleAiSummary}
                            isLoading={isAiLoading}
                            className="mb-2"
                        >
                            Auto-Generate Summary
                        </Button>
                    </div>
                )}

                {section.items.map((item, index) => (
                   <div key={item.id} className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 relative group">
                      <button 
                        onClick={() => removeItem(section.id, item.id)}
                        className="absolute top-2 right-2 text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                          <Trash2 size={16} />
                      </button>
                      
                      {section.type !== ResumeSectionType.Skills && section.type !== ResumeSectionType.Summary && (
                          <div className="grid grid-cols-1 gap-4 mb-3">
                              <input 
                                placeholder="Title (e.g. Senior Manager)" 
                                className="w-full font-medium border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                value={item.title}
                                onChange={e => updateItem(section.id, item.id, 'title', e.target.value)}
                              />
                              <div className="grid grid-cols-2 gap-4">
                                  <input 
                                    placeholder="Subtitle (e.g. Company)" 
                                    className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                    value={item.subtitle}
                                    onChange={e => updateItem(section.id, item.id, 'subtitle', e.target.value)}
                                  />
                                   <input 
                                    placeholder="Date Range" 
                                    className="w-full text-sm border-b border-gray-200 focus:border-blue-500 outline-none pb-1"
                                    value={item.dateRange}
                                    onChange={e => updateItem(section.id, item.id, 'dateRange', e.target.value)}
                                  />
                              </div>
                          </div>
                      )}

                      <div className="relative">
                          <textarea 
                            rows={section.type === ResumeSectionType.Skills ? 2 : 4}
                            className="w-full text-sm border border-gray-300 rounded p-2 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                            placeholder={section.type === ResumeSectionType.Skills ? "List skills here..." : "Describe your achievements..."}
                            value={item.description}
                            onChange={e => updateItem(section.id, item.id, 'description', e.target.value)}
                          />
                          {section.type === ResumeSectionType.Experience && (
                              <button 
                                onClick={() => handleAiPolish(section.id, item.id, item.description)}
                                className="absolute bottom-2 right-2 text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded flex items-center hover:bg-blue-200 transition-colors"
                                disabled={isAiLoading}
                              >
                                  <Sparkles size={12} className="mr-1"/> 
                                  {isAiLoading ? 'Polishing...' : 'Polish'}
                              </button>
                          )}
                      </div>
                   </div>
                ))}

                <Button variant="outline" size="sm" className="w-full border-dashed" onClick={() => addItem(section.id)}>
                    <Plus size={16} className="mr-1"/> Add Item
                </Button>
             </div>
           )}
        </div>
      ))}
    </div>
  );
};