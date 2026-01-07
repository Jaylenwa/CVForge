import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { FileText, MoreVertical, Plus, Clock, Copy, Trash2, Edit2, Share2 } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { INITIAL_RESUME } from '../services/mockData';
import { AppRoute, ResumeData } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { API_BASE } from '../config';
import { ResumeArtboard } from './editor/ResumePreview';
import { Modal } from '../components/ui/Modal';
import { useToast } from '../components/ui/Toast';
 

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();
  const { showToast } = useToast();
  
  
  const [resumes, setResumes] = useState<ResumeData[]>([]);
  const [activeMenu, setActiveMenu] = useState<string | null>(null);
  const [renamingId, setRenamingId] = useState<string | null>(null);
  const [tempTitle, setTempTitle] = useState('');
  const [shareModal, setShareModal] = useState<{ open: boolean; url: string; slug: string }>({ open: false, url: '', slug: '' });
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    (async () => {
      const token = localStorage.getItem('token');
      if (!token) return;
      const res = await fetch(`${API_BASE}/resumes`, { headers: { Authorization: `Bearer ${token}` } });
      const data = await res.json();
      const items = (data.items || []).map((r: any) => ({
        id: r.id || r.ExternalID,
        title: r.title || r.Title,
        templateId: r.templateId || r.TemplateID,
        Theme: r.Theme || { Color: r.Theme?.Color, Font: r.Theme?.Font, Spacing: r.Theme?.Spacing, FontSize: r.Theme?.FontSize },
        lastModified: r.lastModified || r.LastModified || Date.now(),
        Personal: r.Personal || {
          FullName: r.Personal?.FullName,
          JobTitle: r.Personal?.JobTitle || '',
          Email: r.Personal?.Email,
          Phone: r.Personal?.Phone,
          AvatarURL: r.Personal?.AvatarURL,
        },
        sections: (r.sections || []).map((s: any) => ({
          id: s.id || s.ExternalID,
          type: s.type || s.Type,
          title: s.title || s.Title,
          isVisible: s.isVisible ?? s.IsVisible,
          items: (s.items || []).map((it: any) => ({
            id: it.id || it.ExternalID,
            title: it.title || it.Title,
            subtitle: it.subtitle || it.Subtitle,
            timeStart: it.timeStart || it.TimeStart,
            timeEnd: it.timeEnd || it.TimeEnd,
            today: !!(it.today ?? it.Today),
            description: it.description || it.Description,
          })),
        })),
      }));
      setResumes(items);
    })();
  }, []);

  const handleEdit = (id: string) => {
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Editor}?id=${id}`, '_blank');
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
      e.stopPropagation();
      const token = localStorage.getItem('token');
      fetch(`${API_BASE}/resumes/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
        .then(() => setResumes(prev => prev.filter(r => r.id !== id)));
      setActiveMenu(null);
  };

  const handleDuplicate = (resume: ResumeData, e: React.MouseEvent) => {
      e.stopPropagation();
      const token = localStorage.getItem('token');
      const cloneTitle = `${resume.title}${t('dashboard.copySuffix')}`;
      const payload = {
        Title: cloneTitle,
        TemplateID: resume.templateId,
        Personal: resume.Personal || INITIAL_RESUME.Personal || {},
        Job: resume.Job || INITIAL_RESUME.Job || {},
        Theme: resume.Theme || INITIAL_RESUME.Theme || {},
        Sections: (resume.sections || INITIAL_RESUME.sections).map(s => ({
          ExternalID: s.id,
          Type: s.type,
          Title: s.title,
          IsVisible: s.isVisible,
          Items: s.items.map(i => ({
            ExternalID: i.id,
            Title: i.title || '',
            Subtitle: i.subtitle || '',
            Major: i.major || '',
            Degree: i.degree || '',
            TimeStart: i.timeStart || '',
            TimeEnd: i.today ? '' : (i.timeEnd || ''),
            Today: !!i.today,
            Description: i.description || '',
          })),
        })),
      };
      fetch(`${API_BASE}/resumes`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) })
        .then(r => r.json())
        .then(({ id }) => setResumes(prev => [{ ...resume, id, title: cloneTitle, lastModified: Date.now() }, ...prev]));
      setActiveMenu(null);
  };

  const startRename = (resume: ResumeData, e: React.MouseEvent) => {
      e.stopPropagation();
      setRenamingId(resume.id);
      setTempTitle(resume.title);
      setActiveMenu(null);
  };

  const saveRename = (id: string) => {
      setResumes(prev => prev.map(r => r.id === id ? { ...r, title: tempTitle } : r));
      setRenamingId(null);
  };

  const handleShare = async (resumeId: string, e: React.MouseEvent) => {
      e.stopPropagation();
      const token = localStorage.getItem('token');
      if (!token) return;
      try {
        const r = await fetch(`${API_BASE}/resumes/${resumeId}/publish`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } });
        if (!r.ok) {
          let txt = '';
          try { txt = await r.text(); } catch {}
          throw new Error(`HTTP ${r.status} ${r.statusText}${txt ? ' - ' + txt : ''}`);
        }
        const data = await r.json();
        const uiUrl = `${window.location.origin}${window.location.pathname}#${AppRoute.Public.replace(':slug', data.slug)}`;
        setShareModal({ open: true, url: uiUrl, slug: data.slug });
      } catch (err) {
        showToast(t('editor.share.failed'), 'error');
      } finally {
        setActiveMenu(null);
      }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 min-h-[calc(100vh-4rem)]">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('dashboard.title')}</h1>
        <Button onClick={() => navigate(AppRoute.Templates)}>
            <Plus size={18} className="mr-2"/> {t('dashboard.createNew')}
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Create New Placeholder Card - Move to the front */}
        <div 
            onClick={() => navigate(AppRoute.Templates)}
            className="border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center min-h-[180px] cursor-pointer hover:border-blue-500 hover:bg-blue-50 transition-colors group"
        >
            <div className="h-12 w-12 rounded-full bg-gray-100 flex items-center justify-center group-hover:bg-blue-100 mb-4 transition-colors">
                <Plus size={24} className="text-gray-400 group-hover:text-blue-600"/>
            </div>
            <span className="font-medium text-gray-600 group-hover:text-blue-600">{t('dashboard.createNew')}</span>
        </div>

        {resumes.map(resume => (
            <div 
                key={resume.id} 
                className="bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-300 transition-all relative cursor-pointer group"
                onClick={() => handleEdit(resume.id)}
            >
                <div className="h-36 bg-gray-50 rounded-t-lg flex justify-center items-start border-b border-gray-100 relative overflow-hidden">
                     {/* Resume Thumbnail */}
                     <div style={{ width: 'calc(210mm * 0.25)', height: 'calc(297mm * 0.25)' }} className="relative select-none pointer-events-none shadow-sm bg-white">
                        <ResumeArtboard 
                            data={resume} 
                            scale={0.25} 
                            disableShadow={true} 
                            style={{ margin: 0 }}
                        />
                     </div>
                     
                     <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors z-10" />

                     <div className="absolute inset-0 hidden group-hover:flex items-center justify-center z-20">
                        <Button size="sm">
                            {t('common.edit')}
                        </Button>
                     </div>
                </div>
                <div className="p-4">
                    <div className="flex justify-between items-start mb-2">
                        <div className="flex-1 mr-2">
                            {renamingId === resume.id ? (
                                <div className="flex items-center" onClick={e => e.stopPropagation()}>
                                    <input 
                                        className="w-full border border-blue-300 rounded px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-200"
                                        value={tempTitle}
                                        onChange={(e) => setTempTitle(e.target.value)}
                                        autoFocus
                                        onBlur={() => saveRename(resume.id)}
                                        onKeyDown={(e) => e.key === 'Enter' && saveRename(resume.id)}
                                    />
                                </div>
                            ) : (
                                <h3 className="font-semibold text-lg text-gray-900 truncate" title={resume.title}>{resume.title}</h3>
                            )}
                            <p className="text-sm text-gray-500 mt-0.5 flex items-center">
                                <Clock size={12} className="mr-1"/> {new Date(resume.lastModified).toLocaleDateString()}
                            </p>
                        </div>
                        
                        <div className="relative">
                            <button 
                                onClick={(e) => { e.stopPropagation(); setActiveMenu(activeMenu === resume.id ? null : resume.id); }}
                                className="p-1 rounded-full hover:bg-gray-100 text-gray-500 relative z-20"
                            >
                                <MoreVertical size={20}/>
                            </button>
                            
                            {activeMenu === resume.id && (
                                <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-30 border border-gray-100" onClick={e => e.stopPropagation()}>
                                    <button onClick={(e) => startRename(resume, e)} className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Edit2 size={14} className="mr-2"/> {t('common.rename')}
                                    </button>
                                    <button onClick={(e) => handleShare(resume.id, e)} className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Share2 size={14} className="mr-2"/> {t('common.share')}
                                    </button>
                                    <button onClick={(e) => handleDuplicate(resume, e)} className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Copy size={14} className="mr-2"/> {t('common.duplicate')}
                                    </button>
                                    <button onClick={(e) => handleDelete(resume.id, e)} className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center">
                                        <Trash2 size={14} className="mr-2"/> {t('common.delete')}
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
                {/* Click outside to close menu handler - simplistic approach */}
                {activeMenu === resume.id && (
                    <div className="fixed inset-0 z-10" onClick={(e) => { e.stopPropagation(); setActiveMenu(null); }}></div>
                )}
            </div>
        ))}
      </div>
      <Modal isOpen={shareModal.open} onClose={() => setShareModal(prev => ({ ...prev, open: false }))} title={t('editor.share')}>
        <div className="space-y-6">
          <div className="space-y-2">
            <label className="text-[13px] font-bold text-slate-700 ml-1">{t('editor.share.link')}</label>
            <div className="flex items-center gap-2">
              <div className="flex-1 bg-white border border-slate-200 rounded-xl px-4 py-2.5 flex items-center gap-2 overflow-hidden focus-within:ring-2 focus-within:ring-indigo-500/20 focus-within:border-indigo-500 transition-all">
                <svg className="text-slate-400 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
                <input
                  type="text"
                  readOnly
                  value={shareModal.url}
                  className="flex-1 text-sm text-slate-600 outline-none bg-transparent truncate"
                />
              </div>
              <button
                onClick={async () => {
                  try {
                    await navigator.clipboard.writeText(shareModal.url);
                    setCopied(true);
                    showToast(t('editor.share.copied'), 'success');
                    setTimeout(() => setCopied(false), 2000);
                  } catch {}
                }}
                className={`flex-shrink-0 h-10 px-4 rounded-xl font-semibold text-sm transition-all duration-200 flex items-center gap-2 ${
                  copied
                    ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-200'
                    : 'bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-100'
                }`}
              >
                {copied ? (
                  <>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                    <span>{t('editor.share.copied')}</span>
                  </>
                ) : (
                  <>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                    <span>{t('editor.share.copy')}</span>
                  </>
                )}
              </button>
            </div>
          </div>
          <div className="-mx-6 px-6 mt-2 pt-4 border-t bg-slate-50 flex items-center justify-between">
            <button
              onClick={() => setShareModal(prev => ({ ...prev, open: false }))}
              className="px-4 py-2 text-sm font-bold text-slate-600 hover:text-slate-800 transition-colors"
            >
              {t('common.cancel')}
            </button>
            <button
              onClick={() => { window.open(`#${AppRoute.Public.replace(':slug', shareModal.slug)}`, '_blank'); }}
              className="px-5 py-2 bg-white border border-slate-200 rounded-xl text-sm font-bold text-slate-700 hover:bg-slate-50 hover:border-slate-300 transition-all shadow-sm"
            >
              {t('common.preview')}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
