import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ArrowLeft, Download, Printer, Share2, Layout } from 'lucide-react';
import { EditorForm } from './EditorForm';
import { ResumePreview } from './ResumePreview';
import { API_BASE } from '../../config';
import { Button } from '../../components/ui/Button';
import { INITIAL_RESUME } from '../../services/mockData';
import { ResumeData } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

export const Editor: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [resumeData, setResumeData] = useState<ResumeData>(INITIAL_RESUME);
  const [scale, setScale] = useState(0.8);
  const [isMobilePreview, setIsMobilePreview] = useState(false);
  const { t } = useLanguage();
  const [templates, setTemplates] = useState<Array<{ id: string }>>([]);
  const [exportOpen, setExportOpen] = useState(false);

  // Initialize data based on URL params
  useEffect(() => {
    (async () => {
      try {
        const res = await fetch(`${API_BASE}/templates`);
        const data = await res.json();
        const items = (data.items || []).map((t: any) => ({ id: t.ExternalID || t.id }));
        setTemplates(items);
      } catch {}
    })();
    const resumeId = searchParams.get('id');
    const templateId = searchParams.get('template');

    if (resumeId) {
        const token = localStorage.getItem('token');
        fetch(`${API_BASE}/resumes/${resumeId}`, { headers: { Authorization: `Bearer ${token}` } })
          .then(r => r.json())
          .then((res: any) => {
              const mapped: ResumeData = {
                id: res.ExternalID || resumeId,
                title: res.Title,
                templateId: res.TemplateID,
                themeConfig: { color: res.ThemeColor, fontFamily: res.ThemeFont, spacing: res.ThemeSpacing },
                lastModified: res.LastModified,
                personalInfo: {
                  fullName: res.FullName,
                  jobTitle: res.JobTitle || '',
                  email: res.Email,
                  phone: res.Phone,
                  address: res.Address,
                  website: res.Website,
                  avatarUrl: res.AvatarURL
                },
                sections: (res.Sections || []).map((s: any) => ({
                  id: s.ExternalID || s.ID,
                  type: s.Type,
                  title: s.Title,
                  isVisible: s.IsVisible,
                  items: (s.Items || []).map((i: any) => ({
                    id: i.ExternalID || i.ID,
                    title: i.Title,
                    subtitle: i.Subtitle,
                    dateRange: i.DateRange,
                    location: i.Location,
                    description: i.Description
                  }))
                }))
              }
              setResumeData(mapped);
          });
        return;
    }

    if (templateId) {
        const token = localStorage.getItem('token');
        const payload = { ...INITIAL_RESUME, templateId };
        fetch(`${API_BASE}/resumes`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) })
          .then(r => r.json())
          .then(({ id }) => {
            setResumeData(prev => ({ ...prev, templateId, id }));
            window.history.replaceState(null, '', `#${window.location.pathname}?id=${id}`);
          });
    }
  }, [searchParams]);

  const handlePrint = () => {
    window.print();
  };

  const handleChangeTemplate = () => {
      if (!templates.length) return;
      const currentIdx = templates.findIndex(t => t.id === resumeData.templateId);
      const nextIdx = currentIdx >= 0 ? (currentIdx + 1) % templates.length : 0;
      setResumeData({ ...resumeData, templateId: templates[nextIdx].id });
  };

  const handleSave = () => {
    const token = localStorage.getItem('token');
    const payload = {
      title: resumeData.title,
      templateId: resumeData.templateId,
      themeConfig: resumeData.themeConfig,
      personalInfo: resumeData.personalInfo,
      sections: resumeData.sections.map(s => ({
        id: s.id,
        type: s.type,
        title: s.title,
        isVisible: s.isVisible,
        items: s.items.map(i => ({ id: i.id, title: i.title, subtitle: i.subtitle, dateRange: i.dateRange, location: i.location, description: i.description }))
      }))
    };
    fetch(`${API_BASE}/resumes/${resumeData.id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) })
      .then(() => {
        // no-op; could show a toast
      });
  };

  const handleExportPDF = () => {
    const token = localStorage.getItem('token');
    if (!resumeData.id) return;
    fetch(`${API_BASE}/resumes/${resumeData.id}/pdf`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } })
      .then(async r => {
        if (!r.ok) throw new Error('export failed');
        const blob = await r.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${resumeData.title || 'resume'}.pdf`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        URL.revokeObjectURL(url);
      })
      .catch(() => {
        window.open(`#${window.location.pathname}?id=${resumeData.id}`.replace('editor', 'print'), '_blank');
      });
  };

  const handleExportImage = () => {
    const token = localStorage.getItem('token');
    if (!resumeData.id) return;
    fetch(`${API_BASE}/resumes/${resumeData.id}/image`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } })
      .then(async r => {
        if (!r.ok) throw new Error('export failed');
        const blob = await r.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${resumeData.title || 'resume'}.png`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        URL.revokeObjectURL(url);
      })
      .catch(() => {
        window.open(`#${window.location.pathname}?id=${resumeData.id}`.replace('editor', 'print'), '_blank');
      });
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
            <Button variant="primary" size="sm" onClick={handleSave}>
                Save
            </Button>
            <div className="relative">
              <Button variant="secondary" size="sm" icon={<Printer size={16}/>} onClick={() => setExportOpen(!exportOpen)}>
                  {t('common.download')}
              </Button>
              {exportOpen && (
                <div className="absolute right-0 mt-2 w-40 bg-white rounded-md shadow-lg border border-gray-200 z-30">
                  <button className="w-full text-left px-3 py-2 text-sm hover:bg-gray-50" onClick={() => { setExportOpen(false); handleExportPDF(); }}>
                    PDF
                  </button>
                  <button className="w-full text-left px-3 py-2 text-sm hover:bg-gray-50" onClick={() => { setExportOpen(false); handleExportImage(); }}>
                    PNG
                  </button>
                </div>
              )}
            </div>
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
