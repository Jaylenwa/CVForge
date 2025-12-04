import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ArrowLeft, Download, Printer, Share2, Layout } from 'lucide-react';
import { EditorForm } from './EditorForm';
import { ResumePreview } from './ResumePreview';
import { Button } from '../../components/ui/Button';
import { INITIAL_RESUME, MOCK_TEMPLATES, MOCK_USER_RESUMES } from '../../services/mockData';
import { ResumeData } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

export const Editor: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [resumeData, setResumeData] = useState<ResumeData>(INITIAL_RESUME);
  const [scale, setScale] = useState(0.8);
  const [isMobilePreview, setIsMobilePreview] = useState(false);
  const { t } = useLanguage();

  // Initialize data based on URL params
  useEffect(() => {
    const resumeId = searchParams.get('id');
    const templateId = searchParams.get('template');

    if (resumeId) {
        // Mock fetching existing resume
        const found = MOCK_USER_RESUMES.find(r => r.id === resumeId);
        if (found) {
            setResumeData(found);
            return;
        }
    }

    if (templateId) {
        // Initialize new resume with selected template
        setResumeData(prev => ({ 
            ...prev, 
            templateId,
            id: Math.random().toString(36).substr(2, 9) // Assign new random ID
        }));
    }
  }, [searchParams]);

  const handlePrint = () => {
    window.print();
  };

  const handleChangeTemplate = () => {
      // Cycle template for demo
      const currentIdx = MOCK_TEMPLATES.findIndex(t => t.id === resumeData.templateId);
      const nextIdx = (currentIdx + 1) % MOCK_TEMPLATES.length;
      setResumeData({ ...resumeData, templateId: MOCK_TEMPLATES[nextIdx].id });
  };

  return (
    <div className="h-screen flex flex-col bg-gray-100 overflow-hidden">
      {/* Editor Toolbar */}
      <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-4 sm:px-6 z-20 print:hidden">
        <div className="flex items-center">
            <Button variant="ghost" size="sm" className="mr-4" onClick={() => window.history.back()}>
                <ArrowLeft size={18} className="mr-2"/> Back
            </Button>
            <div className="hidden md:block h-6 w-px bg-gray-300 mx-2"></div>
            <input 
                type="text" 
                value={resumeData.title} 
                onChange={(e) => setResumeData({...resumeData, title: e.target.value})}
                className="text-lg font-semibold text-gray-900 border-none focus:ring-0 px-2 hover:bg-gray-50 rounded"
            />
        </div>
        
        <div className="flex items-center space-x-2">
            <div className="hidden md:flex items-center mr-4">
                 <button onClick={() => setScale(Math.max(0.5, scale - 0.1))} className="p-1 hover:bg-gray-100 rounded">-</button>
                 <span className="text-xs text-gray-500 w-12 text-center">{Math.round(scale * 100)}%</span>
                 <button onClick={() => setScale(Math.min(1.5, scale + 0.1))} className="p-1 hover:bg-gray-100 rounded">+</button>
            </div>
            
            <Button variant="outline" size="sm" icon={<Layout size={16}/>} onClick={handleChangeTemplate} className="hidden sm:flex">
                Template
            </Button>
            <Button variant="secondary" size="sm" icon={<Printer size={16}/>} onClick={handlePrint}>
                {t('common.download')} / PDF
            </Button>
            <Button 
                className="md:hidden" 
                size="sm" 
                onClick={() => setIsMobilePreview(!isMobilePreview)}
            >
                {isMobilePreview ? t('common.edit') : t('common.preview')}
            </Button>
        </div>
      </header>

      {/* Main Workspace */}
      <div className="flex-grow flex overflow-hidden relative">
        {/* Left: Form Editor */}
        <div className={`w-full md:w-1/2 lg:w-5/12 xl:w-1/3 bg-white h-full z-10 transition-transform duration-300 absolute md:relative ${isMobilePreview ? '-translate-x-full md:translate-x-0' : 'translate-x-0'}`}>
            <EditorForm data={resumeData} onChange={setResumeData} />
        </div>

        {/* Right: Preview */}
        <div className={`w-full md:w-1/2 lg:w-7/12 xl:w-2/3 bg-gray-100 h-full overflow-hidden absolute md:relative transition-transform duration-300 ${isMobilePreview ? 'translate-x-0' : 'translate-x-full md:translate-x-0'}`}>
             <ResumePreview data={resumeData} scale={scale} />
        </div>
      </div>
    </div>
  );
};