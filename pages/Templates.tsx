import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, Filter, Star } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { MOCK_TEMPLATES } from '../services/mockData';
import { AppRoute } from '../types';

export const Templates: React.FC = () => {
  const navigate = useNavigate();
  const [filter, setFilter] = useState('');
  const [selectedTag, setSelectedTag] = useState<string | null>(null);

  const tags = Array.from(new Set(MOCK_TEMPLATES.flatMap(t => t.tags)));

  const filteredTemplates = MOCK_TEMPLATES.filter(t => {
    const matchesSearch = t.name.toLowerCase().includes(filter.toLowerCase());
    const matchesTag = selectedTag ? t.tags.includes(selectedTag) : true;
    return matchesSearch && matchesTag;
  });

  const handleUseTemplate = (templateId: string) => {
    // In a real app, this would create a new resume ID and redirect
    navigate(`${AppRoute.Editor}?template=${templateId}`);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Resume Templates</h1>
          <p className="mt-2 text-gray-500">Choose a professionally designed template to get started.</p>
        </div>
        <div className="mt-4 md:mt-0 flex items-center space-x-4">
           {/* Simple Search */}
           <div className="relative rounded-md shadow-sm">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              className="focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md p-2 border"
              placeholder="Search templates..."
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
        </div>
      </div>

      {/* Filter Chips */}
      <div className="flex flex-wrap gap-2 mb-8">
        <button
          onClick={() => setSelectedTag(null)}
          className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
            selectedTag === null ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          All
        </button>
        {tags.map(tag => (
          <button
            key={tag}
            onClick={() => setSelectedTag(tag)}
            className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
              selectedTag === tag ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            {tag}
          </button>
        ))}
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
                {/* Overlay on hover */}
                <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all duration-300 flex items-center justify-center opacity-0 group-hover:opacity-100">
                    <Button onClick={() => handleUseTemplate(template.id)}>Use Template</Button>
                </div>
                {template.isPremium && (
                    <div className="absolute top-2 right-2 bg-yellow-400 text-yellow-900 text-xs font-bold px-2 py-1 rounded flex items-center">
                        <Star size={12} className="mr-1 fill-current" /> Premium
                    </div>
                )}
             </div>
             <div className="p-4">
                <h3 className="text-lg font-medium text-gray-900">{template.name}</h3>
                <div className="mt-2 flex items-center justify-between text-sm text-gray-500">
                    <span>{template.popularity}% Popularity</span>
                </div>
                <div className="mt-3 flex flex-wrap gap-1">
                    {template.tags.slice(0, 2).map(t => (
                        <span key={t} className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">{t}</span>
                    ))}
                </div>
             </div>
          </div>
        ))}
      </div>
      
      {filteredTemplates.length === 0 && (
          <div className="text-center py-20">
              <p className="text-gray-500 text-lg">No templates found matching your criteria.</p>
              <Button variant="ghost" onClick={() => {setFilter(''); setSelectedTag(null)}} className="mt-4">Clear Filters</Button>
          </div>
      )}
    </div>
  );
};
