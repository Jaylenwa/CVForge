 import React, { useState, useEffect, useRef } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Download, Printer, Share2, Layout, Globe, Eye } from 'lucide-react';
import { EditorForm } from './EditorForm';
import { ResumePreview } from './ResumePreview';
import { API_BASE } from '../../config';
import { Button } from '../../components/ui/Button';
import { Modal } from '../../components/ui/Modal';
import { DownloadModal } from '../../components/ui/DownloadModal';
import { INITIAL_RESUME } from '../../services/mockData';
import { ResumeData, AppRoute } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import { useToast } from '../../components/ui/Toast';

export const Editor: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const resumeIdParam = searchParams.get('id');
  const savedThemeStr = resumeIdParam ? localStorage.getItem(`resume:${resumeIdParam}:theme`) : null;
  const initialTheme = (() => {
    if (savedThemeStr) {
      try {
        const obj = JSON.parse(savedThemeStr);
        if (obj && typeof obj === 'object') {
          return { ...(INITIAL_RESUME.Theme || {}), ...obj };
        }
      } catch {}
    }
    return INITIAL_RESUME.Theme;
  })();
  const initialResumeData: ResumeData = { ...INITIAL_RESUME, Theme: initialTheme };
  const [resumeData, setResumeData] = useState<ResumeData>(initialResumeData);
  const [scale, setScale] = useState(1);
  const [isMobilePreview, setIsMobilePreview] = useState(false);
  const { t, language, setLanguage } = useLanguage();
  const { showToast } = useToast();
  const [templates, setTemplates] = useState<Array<{ id: string }>>([]);
  const [downloadOpen, setDownloadOpen] = useState(false);
  const [exportError, setExportError] = useState<{ open: boolean; title: string; details: string }>({ open: false, title: '', details: '' });
  const hasCreatedFromTemplate = useRef(false);
  const [loading, setLoading] = useState<boolean>(!!resumeIdParam);

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
              const incoming = res || {};
              const sectionsRaw = (incoming.Sections || []);
              const mapped: ResumeData = {
                id: incoming.ExternalID || resumeId,
                title: incoming.Title,
                templateId: incoming.TemplateID,
                language: (incoming.Language || '') === 'en' ? 'en' : 'zh',
                lastModified: incoming.LastModified,
                Personal: { ...(INITIAL_RESUME.Personal || {}), ...(resumeData.Personal || {}), ...(incoming.Personal || {}) },
                Theme: incoming.Theme,
                sections: sectionsRaw.map((s: any) => ({
                  id: s.ExternalID || s.ID,
                  type: s.Type,
                  title: s.Title,
                  isVisible: s.IsVisible,
                  orderNum: s.OrderNum,
                  items: (s.Items || []).map((i: any) => ({
                    id: i.ExternalID || i.ID,
                    title: i.Title,
                    subtitle: i.Subtitle,
                    major: i.Major,
                    degree: i.Degree,
                    timeStart: i.TimeStart,
                    timeEnd: i.TimeEnd,
                    today: !!i.Today,
                    description: i.Description,
                    orderNum: i.OrderNum
                  })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
                })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
              }
              setResumeData(mapped);
              try {
                if (mapped.language) setLanguage(mapped.language);
              } catch {}
              setLoading(false);
          });
        return;
    }

    if (templateId) {
        if (hasCreatedFromTemplate.current) {
          return;
        }
        hasCreatedFromTemplate.current = true;
        const token = localStorage.getItem('token');
        const payload = {
          Title: INITIAL_RESUME.title,
          TemplateID: templateId,
          Language: language,
          Personal: INITIAL_RESUME.Personal || {},
          Theme: INITIAL_RESUME.Theme || {},
          Sections: INITIAL_RESUME.sections.map((s, si) => ({
            ExternalID: s.id,
            Type: s.type,
            Title: s.title,
            IsVisible: s.isVisible,
            OrderNum: si,
            Items: s.items.map((i, ii) => ({
              ExternalID: i.id,
              Title: i.title || '',
              Subtitle: i.subtitle || '',
              Major: i.major || '',
              Degree: i.degree || '',
              TimeStart: i.timeStart || '',
              TimeEnd: i.today ? '' : (i.timeEnd || ''),
              Today: !!i.today,
              Description: i.description || '',
              OrderNum: ii
            }))
          }))
        };
        fetch(`${API_BASE}/resumes`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) })
          .then(r => r.json())
          .then(({ id }) => {
            setResumeData(prev => ({ ...prev, templateId, id, language }));
            window.history.replaceState(null, '', `#${AppRoute.Editor}?id=${id}`);
            setLoading(false);
          });
    }
  }, [searchParams]);

  useEffect(() => {
    const id = resumeData.id;
    const theme = resumeData.Theme;
    if (id && theme) {
      try {
        localStorage.setItem(`resume:${id}:theme`, JSON.stringify(theme));
      } catch {}
    }
  }, [
    resumeData.id,
    resumeData.Theme?.Color,
    resumeData.Theme?.Font,
    resumeData.Theme?.Spacing,
    resumeData.Theme?.FontSize
  ]);

  const handlePrint = () => {
    window.print();
  };

  const handlePreview = () => {
    if (!resumeData.id) return;
    window.open(`#${AppRoute.Print}?id=${resumeData.id}`, '_blank');
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
      Title: resumeData.title,
      TemplateID: resumeData.templateId,
      Language: resumeData.language || language,
      Personal: resumeData.Personal || {},
      Theme: resumeData.Theme || {},
      Sections: resumeData.sections.map((s, si) => ({
        ExternalID: s.id,
        Type: s.type,
        Title: s.title,
        IsVisible: s.isVisible,
        OrderNum: si,
        Items: s.items.map((i, ii) => ({
          ExternalID: i.id,
          Title: i.title || '',
          Subtitle: i.subtitle || '',
          Major: i.major || '',
          Degree: i.degree || '',
          TimeStart: i.timeStart || '',
          TimeEnd: i.today ? '' : (i.timeEnd || ''),
          Today: !!i.today,
          Description: i.description || '',
          OrderNum: ii
        }))
      }))
    };
    fetch(`${API_BASE}/resumes/${resumeData.id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) })
      .then(() => {
        showToast(t('editor.success.saved'), 'success');
      });
  };

  const handleExportPDF = async (): Promise<void> => {
    const token = localStorage.getItem('token');
    if (!resumeData.id) return;
    const submit = await fetch(`${API_BASE}/pdf/exports?resumeId=${encodeURIComponent(resumeData.id)}`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } });
    if (!submit.ok) {
      let txt = '';
      try { txt = await submit.text(); } catch {}
      throw new Error(`HTTP ${submit.status} ${submit.statusText}${txt ? ' - ' + txt : ''}`);
    }
    const { job_id } = await submit.json();
    const start = Date.now();
    while (true) {
      const st = await fetch(`${API_BASE}/pdf/exports/${job_id}`, { headers: { Authorization: `Bearer ${token}` } });
      const data = await st.json();
      if (data.status === 'done') {
        const a = document.createElement('a');
        const base = API_BASE.replace(/\/+$/, '');
        let href = `${base}/pdf/exports/${encodeURIComponent(job_id)}/download`;
        if (data.token) href += `?token=${encodeURIComponent(data.token)}`;
        a.href = href;
        a.referrerPolicy = 'no-referrer';
        a.rel = 'noreferrer';
        a.download = `${resumeData.title || 'resume'}.pdf`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        return;
      }
      if (data.status === 'failed') {
        throw new Error(data.error || 'Export failed');
      }
      if (Date.now() - start > 120000) {
        throw new Error('Export timeout');
      }
      await new Promise(res => setTimeout(res, 2000));
    }
  };

  const handleExportImage = async (): Promise<void> => {
    const token = localStorage.getItem('token');
    if (!resumeData.id) return;
    const r = await fetch(`${API_BASE}/resumes/${resumeData.id}/image`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } });
    if (!r.ok) {
      let txt = '';
      try { txt = await r.text(); } catch {}
      throw new Error(`HTTP ${r.status} ${r.statusText}${txt ? ' - ' + txt : ''}`);
    }
    const blob = await r.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${resumeData.title || 'resume'}.png`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="h-screen flex flex-col bg-gray-100 overflow-hidden">
      {/* Editor Toolbar */}
      <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-4 sm:px-6 z-20 print:hidden">
        <div className="flex items-center">
            <Button variant="ghost" size="sm" className="mr-4" onClick={() => navigate(AppRoute.Templates)}>
                <ArrowLeft size={18} className="mr-2"/> {t('editor.back')}
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
                {t('editor.template')}
            </Button>

            <Button 
                variant="ghost" 
                size="sm" 
                onClick={() => {
                  const next = language === 'en' ? 'zh' : 'en';
                  setLanguage(next);
                  setResumeData(prev => ({ ...prev, language: next }));
                }}
                title={t('lang.switchTitle')}
            >
                <Globe size={16} className="mr-2"/>
                {language === 'en' ? t('lang.en_short') : t('lang.zh_short')}
            </Button>

            <Button variant="outline" size="sm" icon={<Eye size={16}/>} onClick={handlePreview}>
                {t('common.preview')}
            </Button>
            <Button variant="primary" size="sm" onClick={handleSave} className="ml-2">
                {t('editor.save')}
            </Button>
            <Button variant="secondary" size="sm" icon={<Printer size={16}/>} onClick={() => setDownloadOpen(true)}>
                {t('common.download')}
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

      <div className="flex-grow flex overflow-hidden relative">
        {loading ? (
          <div className="flex-1 flex items-center justify-center bg-white">
            <div className="text-sm text-gray-500">{t('common.loading')}</div>
          </div>
        ) : (
          <>
            <div className={`w-full md:w-2/5 lg:w-2/5 xl:w-2/5 bg-white h-full z-10 transition-transform duration-300 absolute md:relative ${isMobilePreview ? '-translate-x-full md:translate-x-0' : 'translate-x-0'}`}>
              <EditorForm data={resumeData} onChange={setResumeData} />
            </div>
            <div className={`w-full md:w-3/5 lg:w-3/5 xl:w-3/5 bg-white h-full overflow-auto no-scrollbar min-h-0 absolute md:relative transition-transform duration-300 md:border-l md:border-gray-200 md:px-6 lg:px-8 xl:px-10 ${isMobilePreview ? 'translate-x-0' : 'translate-x-full md:translate-x-0'}`}>
              <ResumePreview data={resumeData} scale={scale} disableShadow scrollInside={false} />
            </div>
          </>
        )}
      </div>
      
      <Modal isOpen={exportError.open} onClose={() => setExportError(prev => ({ ...prev, open: false }))} title={exportError.title || t('editor.export.failed')}>
        <div className="space-y-4">
          <div className="text-sm text-gray-700 break-words whitespace-pre-wrap">{exportError.details}</div>
          <div className="flex justify-end">
            <Button variant="primary" onClick={() => setExportError(prev => ({ ...prev, open: false }))}>
              {t('common.close')}
            </Button>
          </div>
        </div>
      </Modal>
      <DownloadModal
        isOpen={downloadOpen}
        onClose={() => setDownloadOpen(false)}
        onExportPDF={handleExportPDF}
        onExportImage={handleExportImage}
        onError={(err) => setExportError({ open: true, title: t('editor.export.failed'), details: err?.message ? String(err.message) : String(err) })}
      />
    </div>
  );
};
