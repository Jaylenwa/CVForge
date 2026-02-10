import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { MoreVertical, Plus, Clock, Copy, Trash2, Edit2, Share2 } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { AppRoute, ResumeData } from '../types';
import { useLanguage } from '../contexts/LanguageContext';
import { API_BASE } from '../config';
import { ResumeArtboard } from './editor/ResumePreview';
import { LanguageProvider } from '../contexts/LanguageContext';
import { Modal } from '../components/ui/Modal';
import { useToast } from '../components/ui/Toast';
 

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();
  const { showToast } = useToast();
  
  
  const [resumes, setResumes] = useState<ResumeData[]>([]);
  const [activeMenu, setActiveMenu] = useState<number | string | null>(null);
  const [renamingId, setRenamingId] = useState<number | string | null>(null);
  const [tempTitle, setTempTitle] = useState('');
  type ShareModalState = {
    open: boolean;
    url: string;
    slug: string;
    resumeId: number | string;
    isPublic: boolean;
    hasPassword: boolean;
    passwordEnabled: boolean;
    passwordEditable: boolean;
    password: string;
    expiryEnabled: boolean;
    expiresAt: string;
    views: number;
    lastAccessAt: string;
    saving: boolean;
  };
  const [shareModal, setShareModal] = useState<ShareModalState>({
    open: false,
    url: '',
    slug: '',
    resumeId: '',
    isPublic: true,
    hasPassword: false,
    passwordEnabled: false,
    passwordEditable: true,
    password: '',
    expiryEnabled: false,
    expiresAt: '',
    views: 0,
    lastAccessAt: '',
    saving: false
  });
  const [copied, setCopied] = useState(false);
  const [shareSettingsOpen, setShareSettingsOpen] = useState(false);

  useEffect(() => {
    (async () => {
      const token = localStorage.getItem('token');
      if (!token) return;
      const res = await fetch(`${API_BASE}/resumes`, { headers: { Authorization: `Bearer ${token}` } });
      const data = await res.json();
      const items = (data.items || []).map((r: any) => ({
        id: r.ID ?? r.id,
        title: r.Title ?? r.title,
        templateId: r.TemplateID ?? r.templateId,
        language: (r.language || r.Language) === 'en' ? 'en' : 'zh',
        Theme: r.Theme || { Color: r.Theme?.Color, Font: r.Theme?.Font, Spacing: r.Theme?.Spacing, FontSize: r.Theme?.FontSize },
        lastModified: r.lastModified || r.LastModified || Date.now(),
        Personal: r.Personal || {
          FullName: r.Personal?.FullName,
          Email: r.Personal?.Email,
          Phone: r.Personal?.Phone,
          AvatarURL: r.Personal?.AvatarURL,
        },
        sections: (r.sections || r.Sections || []).map((s: any) => ({
          id: s.ID ?? s.id,
          type: s.type || s.Type,
          title: s.title || s.Title,
          isVisible: (s.isVisible ?? s.IsVisible) ?? true,
          items: (s.items || s.Items || []).map((it: any) => ({
            id: it.ID ?? it.id,
            title: it.title || it.Title,
            subtitle: it.subtitle || it.Subtitle,
            major: it.major || it.Major,
            degree: it.degree || it.Degree,
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

  const handleEdit = (id: number | string) => {
    window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Editor}?id=${id}&returnTo=${encodeURIComponent(AppRoute.Dashboard)}`, '_blank');
  };

  const handleDelete = async (id: number | string, e: React.MouseEvent) => {
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
        Language: resume.language || 'zh',
        Personal: resume.Personal || {},
        Theme: resume.Theme || {},
        Sections: (resume.sections || []).map(s => ({
          Type: s.type,
          Title: s.title,
          IsVisible: s.isVisible,
          Items: s.items.map(i => ({
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

  const saveRename = (id: number | string) => {
      setResumes(prev => prev.map(r => r.id === id ? { ...r, title: tempTitle } : r));
      setRenamingId(null);
  };

  const handleShare = async (resumeId: number | string, e: React.MouseEvent) => {
      e.stopPropagation();
      const token = localStorage.getItem('token');
      if (!token) return;
      try {
        const mapToUiUrl = (slug: string, pageUrl?: string) => {
          const rel = String(pageUrl || '').trim();
          if (rel.includes('#')) {
            const hash = rel.slice(rel.indexOf('#'));
            return `${window.location.origin}${window.location.pathname}${hash}`;
          }
          return `${window.location.origin}${window.location.pathname}#${AppRoute.Public.replace(':slug', slug)}`;
        };

        const openWithShareData = (slug: string, sdata: any) => {
          const exp = sdata?.expiresAt ? String(sdata.expiresAt) : '';
          const initialIsPublic = sdata?.isPublic === undefined ? true : !!sdata?.isPublic;
          const initialHasPassword = sdata?.hasPassword === undefined ? false : !!sdata?.hasPassword;
          const initialExpiresAt = exp ? new Date(exp).toISOString().slice(0, 16) : '';
          const initialPassword = typeof sdata?.password === 'string' ? sdata.password : '';
          const uiUrl = mapToUiUrl(slug, sdata?.pageUrl);
          setShareModal(prev => ({
            ...prev,
            open: true,
            url: uiUrl,
            slug,
            resumeId,
            isPublic: initialIsPublic,
            hasPassword: initialHasPassword,
            password: initialHasPassword ? initialPassword : '',
            passwordEnabled: !!initialIsPublic && !!initialHasPassword,
            passwordEditable: !initialHasPassword,
            expiryEnabled: !!initialIsPublic && !!initialExpiresAt,
            expiresAt: initialExpiresAt,
            views: Number.isFinite(sdata?.views) ? Number(sdata.views) : 0,
            lastAccessAt: sdata?.lastAccessAt ? new Date(String(sdata.lastAccessAt)).toLocaleString() : ''
          }));
          const openSettings = !initialIsPublic || initialHasPassword || !!initialExpiresAt;
          setShareSettingsOpen(!!openSettings);
        };

        const rs = await fetch(`${API_BASE}/resumes/${resumeId}/share`, { headers: { Authorization: `Bearer ${token}` } });
        if (rs.ok) {
          const sdata = await rs.json();
          const slug = String(sdata?.slug || '').trim();
          if (slug) {
            openWithShareData(slug, sdata);
            return;
          }
        }

        const r = await fetch(`${API_BASE}/resumes/${resumeId}/publish`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } });
        if (!r.ok) {
          let txt = '';
          try { txt = await r.text(); } catch {}
          throw new Error(`HTTP ${r.status} ${r.statusText}${txt ? ' - ' + txt : ''}`);
        }
        const data = await r.json();
        openWithShareData(String(data?.slug || ''), data);
      } catch (err) {
        showToast(t('editor.share.failed'), 'error');
      } finally {
        setActiveMenu(null);
      }
  };

  useEffect(() => {
    (async () => {
      if (!shareModal.open || !shareModal.resumeId) return;
      const token = localStorage.getItem('token');
      if (!token) return;
      try {
        const r = await fetch(`${API_BASE}/resumes/${shareModal.resumeId}/share`, { headers: { Authorization: `Bearer ${token}` } });
        if (!r.ok) return;
        const data = await r.json();
        const expiresAt = data?.expiresAt ? String(data.expiresAt) : '';
        const password = typeof data?.password === 'string' ? data.password : '';
        setShareModal(prev => ({
          ...prev,
          isPublic: !!data?.isPublic,
          hasPassword: !!data?.hasPassword,
          passwordEnabled: !!data?.isPublic && !!data?.hasPassword,
          passwordEditable: !data?.hasPassword,
          password: data?.hasPassword ? password : '',
          expiryEnabled: !!data?.isPublic && !!expiresAt,
          expiresAt: expiresAt ? new Date(expiresAt).toISOString().slice(0, 16) : '',
          views: Number.isFinite(data?.views) ? Number(data.views) : 0,
          lastAccessAt: data?.lastAccessAt ? new Date(String(data.lastAccessAt)).toLocaleString() : ''
        }));
        const shouldOpenSettings = !data?.isPublic || !!data?.hasPassword || !!expiresAt;
        setShareSettingsOpen(!!shouldOpenSettings);
      } catch {}
    })();
  }, [shareModal.open, shareModal.resumeId]);

  useEffect(() => {
    if (!shareModal.open) return;
    setCopied(false);
  }, [shareModal.open]);

  const saveShareSettings = async (override?: Partial<ShareModalState>, options?: { silent?: boolean }) => {
    const token = localStorage.getItem('token');
    if (!token || !shareModal.resumeId) return;

    setShareModal(prev => ({ ...prev, ...(override || {}), saving: true }));
    const nextState: ShareModalState = { ...shareModal, ...(override || {}), saving: true };

    try {
      const body: any = { isPublic: nextState.isPublic };
      if (nextState.isPublic) {
        if (nextState.expiryEnabled && nextState.expiresAt) {
          body.expiresAt = new Date(nextState.expiresAt).toISOString();
        } else if (!nextState.expiryEnabled && !!shareModal.expiresAt) {
          body.expiresAt = '';
        }

        if (nextState.passwordEnabled) {
          if ((!nextState.hasPassword || nextState.passwordEditable) && nextState.password) {
            body.password = nextState.password;
          }
        } else if (nextState.hasPassword) {
          body.password = '';
        }
      }

      const r = await fetch(`${API_BASE}/resumes/${nextState.resumeId}/share`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(body)
      });
      if (!r.ok) throw new Error('save failed');
      if (!options?.silent) showToast(t('share.settings.saved'), 'success');

      const rr = await fetch(`${API_BASE}/resumes/${nextState.resumeId}/share`, { headers: { Authorization: `Bearer ${token}` } });
      if (rr.ok) {
        const data = await rr.json();
        const expiresAt = data?.expiresAt ? String(data.expiresAt) : '';
        const password = typeof data?.password === 'string' ? data.password : '';
        setShareModal(prev => ({
          ...prev,
          isPublic: !!data?.isPublic,
          hasPassword: !!data?.hasPassword,
          passwordEnabled: !!data?.isPublic && !!data?.hasPassword,
          passwordEditable: !data?.hasPassword,
          password: data?.hasPassword ? password : '',
          expiryEnabled: !!data?.isPublic && !!expiresAt,
          expiresAt: expiresAt ? new Date(expiresAt).toISOString().slice(0, 16) : '',
          views: Number.isFinite(data?.views) ? Number(data.views) : prev.views,
          lastAccessAt: data?.lastAccessAt ? new Date(String(data.lastAccessAt)).toLocaleString() : prev.lastAccessAt
        }));
      }
    } catch {
      showToast(t('share.settings.saveFailed'), 'error');
    } finally {
      setShareModal(prev => ({ ...prev, saving: false }));
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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
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
                        <LanguageProvider languageOverride={resume.language || 'zh'}>
                          <ResumeArtboard 
                              data={resume} 
                              scale={0.25} 
                              disableShadow={true} 
                              showPageHint={false}
                              style={{ margin: 0 }}
                          />
                        </LanguageProvider>
                     </div>
                     
                     <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors z-10" />

                     <div className="absolute inset-0 hidden group-hover:flex items-center justify-center z-20">
                        <div className="flex flex-col gap-2">
                          <Button size="sm" onClick={(e) => { e.stopPropagation(); handleEdit(resume.id); }}>
                              {t('common.edit')}
                          </Button>
                          <Button size="sm" variant="outline" onClick={(e) => { e.stopPropagation(); window.open(`${window.location.origin}${window.location.pathname}#${AppRoute.Print}?id=${resume.id}`, '_blank'); }}>
                              {t('common.preview')}
                          </Button>
                        </div>
                     </div>
                </div>
                <div className="p-4">
                    <div className="flex justify-between items-start mb-2 min-w-0">
                        <div className="flex-1 mr-2 min-w-0">
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
                                <div className="relative group/title min-w-0">
                                    <h3 className="block font-semibold text-lg text-gray-900 truncate" aria-label={resume.title}>{resume.title}</h3>
                                    <div className="pointer-events-none absolute left-0 bottom-full mb-1 z-50 w-max max-w-[320px] rounded-md bg-gray-900 px-2 py-1 text-xs text-white opacity-0 group-hover/title:opacity-100">
                                        <div className="whitespace-normal break-words">{resume.title}</div>
                                    </div>
                                </div>
                            )}
                            <p className="text-sm text-gray-500 mt-0.5 flex items-center">
                                <Clock size={12} className="mr-1"/> {new Date(resume.lastModified).toLocaleDateString()}
                            </p>
                        </div>
                        
                        <div className="relative flex-shrink-0">
                            <button 
                                onClick={(e) => { e.stopPropagation(); setActiveMenu(activeMenu === resume.id ? null : resume.id); }}
                                className="p-1 rounded-full hover:bg-gray-100 text-gray-500 relative z-20"
                            >
                                <MoreVertical size={20}/>
                            </button>
                            
                            {activeMenu === resume.id && (
                                <div className="absolute right-0 mt-2 w-36 bg-white rounded-md shadow-lg py-1 z-30 border border-gray-100" onClick={e => e.stopPropagation()}>
                                    <button onClick={(e) => startRename(resume, e)} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Edit2 size={14} className="mr-2"/> {t('common.rename')}
                                    </button>
                                    <button onClick={(e) => handleShare(resume.id, e)} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Share2 size={14} className="mr-2"/> {t('common.share')}
                                    </button>
                                    <button onClick={(e) => handleDuplicate(resume, e)} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center">
                                        <Copy size={14} className="mr-2"/> {t('common.duplicate')}
                                    </button>
                                    <button onClick={(e) => handleDelete(resume.id, e)} className="w-full text-left px-3 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center">
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
      <Modal
        isOpen={shareModal.open}
        onClose={() => setShareModal(prev => ({ ...prev, open: false }))}
        hideHeader
        closeOnBackdrop
        overlayClassName="bg-slate-900/60 backdrop-blur-sm"
        panelClassName="max-w-lg rounded-[2rem] overflow-hidden border border-slate-200 shadow-2xl"
        bodyClassName="p-0"
      >
        <div className="px-8 pt-8 pb-4 flex items-start justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-800">{t('editor.share')}</h2>
            <p className="text-slate-500 text-sm mt-1">
              {t('share.modal.desc')}
            </p>
          </div>
          <button
            onClick={() => setShareModal(prev => ({ ...prev, open: false }))}
            className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-slate-100 text-slate-400 transition-colors"
            aria-label={t('common.close')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
          </button>
        </div>

        <div className="px-8 pb-6">
          <div className="space-y-4">
            <div className="relative">
              <label className="text-xs font-bold text-slate-400 uppercase tracking-widest mb-1.5 block px-1">
                {t('editor.share.link')}
              </label>
              <div className="flex items-center gap-2">
                <div className="flex-[1.2] relative flex items-center">
                  <input
                    readOnly
                    value={shareModal.url}
                    className="w-full h-12 bg-slate-50 border-2 border-slate-100 rounded-2xl px-4 pr-4 text-slate-700 text-sm font-medium focus:outline-none focus:border-indigo-500 transition-all cursor-default"
                  />
                </div>
                <button
                  onClick={async () => {
                    try {
                      await navigator.clipboard.writeText(shareModal.url);
                      setCopied(true);
                      showToast(t('editor.share.copied'), 'success');
                      setTimeout(() => setCopied(false), 2000);
                    } catch {
                      showToast(t('common.copyFailed'), 'error');
                    }
                  }}
                  className={`relative flex items-center justify-center gap-2 h-12 px-3 min-w-[96px] rounded-2xl text-sm font-bold text-white transition-all transform active:scale-95 shadow-lg ${
                    copied
                      ? 'bg-emerald-500 shadow-emerald-200'
                      : 'bg-indigo-600 hover:bg-indigo-700 shadow-indigo-200'
                  }`}
                >
                  {copied ? (
                    <>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                      <span>{t('common.copied')}</span>
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

            <div className="flex items-center justify-between px-1">
              <div className="flex items-center gap-2 text-sm text-slate-600 font-medium">
                <div className={`w-2 h-2 rounded-full ${shareModal.isPublic ? 'bg-emerald-500' : 'bg-amber-500'}`}></div>
                <span>
                  {shareModal.isPublic
                    ? t('share.status.public')
                    : t('share.status.restricted')}
                </span>
              </div>
              <button
                onClick={() => setShareSettingsOpen(p => {
                  return !p;
                })}
                className={`text-sm font-semibold flex items-center gap-1.5 transition-colors ${shareSettingsOpen ? 'text-indigo-600' : 'text-slate-400 hover:text-slate-600'}`}
              >
                <svg className={`transition-transform ${shareSettingsOpen ? 'rotate-90' : ''}`} xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="4" x2="4" y1="21" y2="14"/><line x1="4" x2="4" y1="10" y2="3"/><line x1="12" x2="12" y1="21" y2="12"/><line x1="12" x2="12" y1="8" y2="3"/><line x1="20" x2="20" y1="21" y2="16"/><line x1="20" x2="20" y1="12" y2="3"/><line x1="2" x2="6" y1="14" y2="14"/><line x1="10" x2="14" y1="8" y2="8"/><line x1="18" x2="22" y1="16" y2="16"/></svg>
                {t('share.settings.entry')}
              </button>
            </div>
          </div>
        </div>

        <div className={`overflow-hidden transition-all duration-300 ease-in-out ${shareSettingsOpen ? 'max-h-[520px] opacity-100 mb-6' : 'max-h-0 opacity-0'}`}>
          <div className="px-8 border-t border-slate-100 pt-6 space-y-6">
            <div className="grid grid-cols-2 gap-3 p-1 bg-slate-50 rounded-2xl border border-slate-200">
              <button
                onClick={() => setShareModal(prev => ({ ...prev, isPublic: true }))}
                className={`flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all ${shareModal.isPublic ? 'bg-white shadow-sm text-slate-800' : 'text-slate-400 hover:text-slate-600'}`}
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>
                {t('share.visibility.public')}
              </button>
              <button
                onClick={async () => {
                  if (!shareModal.isPublic) return;
                  setShareModal(prev => ({ ...prev, isPublic: false, passwordEnabled: false, passwordEditable: true, expiryEnabled: false, password: '' }));
                  await saveShareSettings({ isPublic: false }, { silent: true });
                }}
                className={`flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all ${!shareModal.isPublic ? 'bg-white shadow-sm text-slate-800' : 'text-slate-400 hover:text-slate-600'}`}
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                {t('share.visibility.restricted')}
              </button>
            </div>

            {shareModal.isPublic ? (
              <>
                <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex flex-col">
                  <span className="text-sm font-bold text-slate-700">{t('share.settings.password')}</span>
                  <span className="text-xs text-slate-500">{t('share.settings.passwordDesc')}</span>
                </div>
                <div className="flex items-center gap-3">
                  {shareModal.passwordEnabled && shareModal.hasPassword && !shareModal.passwordEditable ? (
                    <button
                      type="button"
                      onClick={() => setShareModal(prev => ({ ...prev, passwordEditable: true, password: '' }))}
                      className="text-xs font-bold text-indigo-600 hover:text-indigo-800 transition-colors"
                    >
                      {t('share.settings.resetPassword')}
                    </button>
                  ) : null}
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      checked={shareModal.passwordEnabled}
                      onChange={(e) => {
                        const enabled = e.target.checked;
                        setShareModal(prev => ({
                          ...prev,
                          passwordEnabled: enabled,
                          passwordEditable: enabled ? !prev.hasPassword : true,
                          password: enabled ? prev.password : ''
                        }));
                      }}
                    />
                    <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                  </label>
                </div>
              </div>

              {shareModal.passwordEnabled && (
                <div className="animate-in slide-in-from-top-2 duration-200">
                  <input
                    type="text"
                    value={shareModal.password}
                    onChange={(e) => setShareModal(prev => ({ ...prev, password: e.target.value }))}
                    disabled={shareModal.hasPassword && !shareModal.passwordEditable}
                    placeholder={shareModal.passwordEditable ? (shareModal.hasPassword ? t('share.settings.passwordHintHas') : t('share.settings.passwordHint')) : ''}
                    className={`w-full border border-slate-200 rounded-xl py-3 px-4 text-sm transition-all ${shareModal.hasPassword && !shareModal.passwordEditable ? 'bg-slate-50 text-slate-400 cursor-not-allowed' : 'bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500'}`}
                  />
                </div>
              )}

              <div className="flex items-center justify-between">
                <div className="flex flex-col">
                  <span className="text-sm font-bold text-slate-700">{t('share.settings.expiresAt')}</span>
                  <span className="text-xs text-slate-500">{t('share.settings.expiresDesc')}</span>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    className="sr-only peer"
                    checked={shareModal.expiryEnabled}
                    onChange={(e) => {
                      const enabled = e.target.checked;
                      setShareModal(prev => ({
                        ...prev,
                        expiryEnabled: enabled,
                        expiresAt: enabled
                          ? (prev.expiresAt || new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().slice(0, 16))
                          : '',
                      }));
                    }}
                  />
                  <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                </label>
              </div>

              {shareModal.expiryEnabled && (
                <div className="animate-in slide-in-from-top-2 duration-200">
                  <input
                    type="datetime-local"
                    value={shareModal.expiresAt}
                    onChange={(e) => setShareModal(prev => ({ ...prev, expiresAt: e.target.value }))}
                    className="w-full bg-white border border-slate-200 rounded-xl py-3 px-4 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                  />
                </div>
              )}
            </div>

            <div className="pt-2">
              <button
                onClick={() => saveShareSettings()}
                disabled={shareModal.saving}
                className="w-full py-3 rounded-xl bg-slate-800 hover:bg-slate-900 disabled:bg-slate-300 text-white font-bold text-sm transition-all shadow-lg shadow-slate-200"
              >
                {t('share.settings.save')}
              </button>
            </div>
                </>
              ) : null}
          </div>
        </div>

        <div className="bg-slate-50/50 px-8 py-5 border-t border-slate-100 flex items-center justify-between">
          <div className="flex items-center gap-6">
            <div className="flex flex-col">
              <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{t('share.settings.views')}</span>
              <span className="text-sm font-bold text-slate-700">{shareModal.views}</span>
            </div>
            <div className="flex flex-col">
              <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{t('share.settings.lastAccessAt')}</span>
              <span className="text-sm font-bold text-slate-700">{shareModal.lastAccessAt || '-'}</span>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={() => saveShareSettings({ isPublic: false })}
              disabled={shareModal.saving || !shareModal.isPublic}
              className="text-sm font-semibold text-slate-500 hover:text-slate-800 disabled:text-slate-300 transition-colors"
            >
              {t('share.actions.disableLink')}
            </button>
            <div className="w-[1px] h-4 bg-slate-200"></div>
            <button
              onClick={() => { if (shareModal.url) window.open(shareModal.url, '_blank'); }}
              className="flex items-center gap-1.5 text-sm font-bold text-indigo-600 hover:text-indigo-800 transition-colors"
            >
              {t('common.preview')}
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M7 7h10v10"/><path d="M7 17 17 7"/></svg>
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
