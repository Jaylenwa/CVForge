import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FileText, MoreVertical, Plus, Clock, Copy, Trash2, Edit2 } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { INITIAL_RESUME, MOCK_USER_RESUMES } from '../services/mockData';
import { AppRoute, ResumeData } from '../types';
import { useLanguage } from '../contexts/LanguageContext';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();
  
  // Local state to simulate database, initializing with mock user resumes
  const [resumes, setResumes] = useState<ResumeData[]>(MOCK_USER_RESUMES);
  const [activeMenu, setActiveMenu] = useState<string | null>(null);
  const [renamingId, setRenamingId] = useState<string | null>(null);
  const [tempTitle, setTempTitle] = useState('');

  const handleEdit = (id: string) => {
    navigate(`${AppRoute.Editor}?id=${id}`);
  };

  const handleDelete = (id: string, e: React.MouseEvent) => {
      e.stopPropagation(); // Prevent navigation when clicking delete
      if (window.confirm('Are you sure you want to delete this resume?')) {
          setResumes(prev => prev.filter(r => r.id !== id));
      }
      setActiveMenu(null);
  };

  const handleDuplicate = (resume: ResumeData, e: React.MouseEvent) => {
      e.stopPropagation();
      const newResume = {
          ...resume,
          id: Math.random().toString(36).substr(2, 9),
          title: `${resume.title} (Copy)`,
          lastModified: Date.now()
      };
      setResumes(prev => [newResume, ...prev]);
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

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 min-h-[calc(100vh-4rem)]">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('dashboard.title')}</h1>
        <Button onClick={() => navigate(AppRoute.Templates)}>
            <Plus size={18} className="mr-2"/> {t('dashboard.createNew')}
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {resumes.map(resume => (
            <div 
                key={resume.id} 
                className="bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-300 transition-all relative cursor-pointer group"
                onClick={() => handleEdit(resume.id)}
            >
                <div className="h-40 bg-gray-100 rounded-t-lg flex items-center justify-center border-b border-gray-100 relative overflow-hidden">
                     {/* Preview Mockup */}
                     <div className="absolute inset-0 bg-white opacity-50"></div>
                     <FileText size={48} className="text-gray-300 relative z-10" />
                     
                     <div className="absolute inset-0 bg-black/50 hidden group-hover:flex items-center justify-center space-x-3 transition-all">
                        <Button size="sm">
                            {t('common.edit')}
                        </Button>
                     </div>
                </div>
                <div className="p-5">
                    <div className="flex justify-between items-start mb-4">
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
                            <p className="text-sm text-gray-500 mt-1 flex items-center">
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
        
        {/* Create New Placeholder Card */}
        <div 
            onClick={() => navigate(AppRoute.Templates)}
            className="border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center h-full min-h-[250px] cursor-pointer hover:border-blue-500 hover:bg-blue-50 transition-colors group"
        >
            <div className="h-12 w-12 rounded-full bg-gray-100 flex items-center justify-center group-hover:bg-blue-100 mb-4 transition-colors">
                <Plus size={24} className="text-gray-400 group-hover:text-blue-600"/>
            </div>
            <span className="font-medium text-gray-600 group-hover:text-blue-600">{t('dashboard.createNew')}</span>
        </div>
      </div>
    </div>
  );
};