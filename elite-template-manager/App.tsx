
import React, { useState, useMemo } from 'react';
import { Header } from './components/Header';
import { TemplateTable } from './components/TemplateTable';
import { StatsSection } from './components/StatsSection';
import { FilterBar } from './components/FilterBar';
import { Template, Industry, FilterState } from './types';

const INITIAL_DATA: Template[] = [
  { id: 't1', name: '经典专业版', industry: 'General', popularity: 98, isPremium: false, tags: ['专业', '简洁', 'ATS 友好'], lastUpdated: '2024-03-20' },
  { id: 't3', name: '技术极简', industry: 'IT', popularity: 92, isPremium: false, tags: ['极简', '技术', '清爽'], lastUpdated: '2024-03-18' },
  { id: 't14', name: '中文极简', industry: 'General', popularity: 90, isPremium: false, tags: ['极简', '中文', '简洁'], lastUpdated: '2024-03-15' },
  { id: 't6', name: '优雅青绿', industry: 'General', popularity: 88, isPremium: false, tags: ['现代', '清新', '入门'], lastUpdated: '2024-03-10' },
  { id: 't12', name: '科技极客', industry: 'IT', popularity: 88, isPremium: false, tags: ['技术', '中文', '极客'], lastUpdated: '2024-03-05' },
  { id: 't2', name: '现代暗色', industry: 'Creative', popularity: 85, isPremium: true, tags: ['创意', '设计', '初创'], lastUpdated: '2024-03-01' },
  { id: 't10', name: '商务专业', industry: 'Finance', popularity: 85, isPremium: false, tags: ['商务', '中文', '专业'], lastUpdated: '2024-02-28' },
  { id: 't26', name: '代码开发', industry: 'IT', popularity: 85, isPremium: false, tags: ['代码', '中文', '开发'], lastUpdated: '2024-02-25' },
  { id: 't9', name: '现代中文', industry: 'General', popularity: 82, isPremium: false, tags: ['现代', '中文', '清爽'], lastUpdated: '2024-02-20' },
  { id: 't28', name: '金融数据', industry: 'Finance', popularity: 82, isPremium: false, tags: ['金融', '中文', '数据'], lastUpdated: '2024-02-15' },
];

const App: React.FC = () => {
  const [templates, setTemplates] = useState<Template[]>(INITIAL_DATA);
  const [filters, setFilters] = useState<FilterState>({ keyword: '', industry: 'All' });

  const filteredTemplates = useMemo(() => {
    return templates.filter(t => {
      const matchKeyword = t.name.toLowerCase().includes(filters.keyword.toLowerCase()) || 
                          t.tags.some(tag => tag.toLowerCase().includes(filters.keyword.toLowerCase()));
      const matchIndustry = filters.industry === 'All' || t.industry === filters.industry;
      return matchKeyword && matchIndustry;
    });
  }, [templates, filters]);

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this template?')) {
      setTemplates(prev => prev.filter(t => t.id !== id));
    }
  };

  const handleSync = () => {
    alert('Syncing mock data to backend...');
  };

  return (
    <div className="min-h-screen bg-slate-50/50">
      <Header onSync={handleSync} />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        <StatsSection templates={templates} />
        
        <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
          <FilterBar 
            filters={filters} 
            setFilters={setFilters} 
            onAdd={() => alert('Opening Create Modal')}
          />
          
          <TemplateTable 
            templates={filteredTemplates} 
            onDelete={handleDelete}
            onUpdate={(id) => alert(`Update template ${id}`)}
            onPreview={(id) => alert(`Preview template ${id}`)}
          />
        </div>
      </main>
      
      <footer className="py-8 text-center text-slate-400 text-sm">
        &copy; 2024 Elite Template Management System. All rights reserved.
      </footer>
    </div>
  );
};

export default App;
