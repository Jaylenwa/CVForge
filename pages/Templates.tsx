import React, { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Search, Star, Filter } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { MOCK_TEMPLATES } from '../services/mockData';
import { AppRoute } from '../types';
import { useLanguage } from '../contexts/LanguageContext';

export const Templates: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { t } = useLanguage();
  
  const [filter, setFilter] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [selectedLevel, setSelectedLevel] = useState<string>('All');

  useEffect(() => {
    const tag = searchParams.get('tag');
    if (tag) {
        // Simple mapping for demo purposes
        if (['IT', 'Finance', 'Creative', 'General'].includes(tag)) {
            setSelectedCategory(tag);
        } else if (['Intern', 'Junior', 'Senior', 'Executive'].includes(tag)) {
            setSelectedLevel(tag);
        }
    }
  }, [searchParams]);

  const categories = ['All', 'IT', 'Finance', 'Creative', 'General'];
  const levels = ['All', 'Intern', 'Junior', 'Senior', 'Executive'];

  const filteredTemplates = MOCK_TEMPLATES.filter(t => {
    const matchesSearch = t.name.toLowerCase().includes(filter.toLowerCase());
    const matchesCategory = selectedCategory === 'All' || t.category === selectedCategory;
    const matchesLevel = selectedLevel === 'All' || t.level === selectedLevel;
    return matchesSearch && matchesCategory && matchesLevel;
  });

  const handleUseTemplate = (templateId: string) => {
    navigate(`${AppRoute.Editor}?template=${templateId}`);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t('templates.title')}</h1>
          <p className="mt-2 text-gray-500">{t('templates.desc')}</p>
        </div>
        <div className="mt-4 md:mt-0 flex items-center space-x-4">
           <div className="relative rounded-md shadow-sm w-full md:w-64">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              className="focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md p-2 border"
              placeholder={t('templates.search')}
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
        </div>
      </div>

      {/* Advanced Filters */}
      <div className="flex flex-wrap items-center gap-4 mb-8 bg-gray-50 p-4 rounded-lg border border-gray-200">
        <div className="flex items-center space-x-2">
            <Filter size={18} className="text-gray-500" />
            <span className="font-medium text-gray-700 text-sm">{t('templates.filters.label')}</span>
        </div>
        
        {/* Category Dropdown */}
        <div className="flex items-center space-x-2">
            <label className="text-sm text-gray-500">{t('templates.filter.industry')}:</label>
            <select 
                value={selectedCategory} 
                onChange={(e) => setSelectedCategory(e.target.value)}
                className="block w-full pl-3 pr-10 py-1.5 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md border"
            >
                {categories.map(c => <option key={c} value={c}>{t(`templates.category.${c}`)}</option>)}
            </select>
        </div>

        {/* Level Dropdown */}
        <div className="flex items-center space-x-2">
            <label className="text-sm text-gray-500">{t('templates.filter.level')}:</label>
            <select 
                value={selectedLevel} 
                onChange={(e) => setSelectedLevel(e.target.value)}
                className="block w-full pl-3 pr-10 py-1.5 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md border"
            >
                {levels.map(l => <option key={l} value={l}>{t(`templates.level.${l}`)}</option>)}
            </select>
        </div>

         <div className="ml-auto">
            <button 
                onClick={() => { setSelectedCategory('All'); setSelectedLevel('All'); setFilter(''); }}
                className="text-sm text-blue-600 hover:text-blue-800"
            >
                {t('templates.actions.clearAll')}
            </button>
         </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
        {filteredTemplates.map(template => (
          <div key={template.id} className="group relative bg-white border border-gray-200 rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all duration-300">
             <div className="aspect-[3/4] w-full bg-gray-200 overflow-hidden relative">
                <img 
                  src={template.thumbnail} 
                  alt={template.name} 
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                />
                <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all duration-300 flex items-center justify-center opacity-0 group-hover:opacity-100">
                    <Button onClick={() => handleUseTemplate(template.id)}>{t('templates.actions.useTemplate')}</Button>
                </div>
                {template.isPremium && (
                    <div className="absolute top-2 right-2 bg-yellow-400 text-yellow-900 text-xs font-bold px-2 py-1 rounded flex items-center">
                        <Star size={12} className="mr-1 fill-current" /> {t('templates.badge.premium')}
                    </div>
                )}
             </div>
             <div className="p-4">
                <h3 className="text-lg font-medium text-gray-900">{template.name}</h3>
                <div className="mt-2 flex items-center justify-between text-sm text-gray-500">
                    <span>{template.popularity}% {t('templates.meta.popularity')}</span>
                </div>
                <div className="mt-3 flex flex-wrap gap-1">
                    <span className="px-2 py-0.5 bg-blue-50 text-blue-700 text-xs rounded border border-blue-100">{template.category}</span>
                    <span className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">{template.level}</span>
                </div>
             </div>
          </div>
        ))}
      </div>
      
      {filteredTemplates.length === 0 && (
          <div className="text-center py-20 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
              <p className="text-gray-500 text-lg">{t('templates.empty')}</p>
              <Button variant="ghost" onClick={() => {setFilter(''); setSelectedCategory('All'); setSelectedLevel('All')}} className="mt-4">{t('templates.actions.clearFilters')}</Button>
          </div>
      )}
    </div>
  );
};
