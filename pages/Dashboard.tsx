import React from 'react';
import { useNavigate } from 'react-router-dom';
import { FileText, MoreVertical, Plus, Clock, Eye } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { INITIAL_RESUME } from '../services/mockData';
import { AppRoute } from '../types';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  // Mock List
  const resumes = [
      { ...INITIAL_RESUME, id: '1', title: 'Software Engineer Resume', lastModified: Date.now() - 1000000 },
      { ...INITIAL_RESUME, id: '2', title: 'Management CV', lastModified: Date.now() - 86400000 * 5, templateId: 't2' }
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">My Resumes</h1>
        <Button onClick={() => navigate(AppRoute.Templates)}>
            <Plus size={18} className="mr-2"/> Create New
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {resumes.map(resume => (
            <div key={resume.id} className="bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
                <div className="h-40 bg-gray-100 rounded-t-lg flex items-center justify-center border-b border-gray-100">
                     <FileText size={48} className="text-gray-300" />
                </div>
                <div className="p-5">
                    <div className="flex justify-between items-start">
                        <div>
                            <h3 className="font-semibold text-lg text-gray-900 truncate pr-4" title={resume.title}>{resume.title}</h3>
                            <p className="text-sm text-gray-500 mt-1 flex items-center">
                                <Clock size={12} className="mr-1"/> Updated {new Date(resume.lastModified).toLocaleDateString()}
                            </p>
                        </div>
                        <button className="text-gray-400 hover:text-gray-600">
                            <MoreVertical size={20}/>
                        </button>
                    </div>
                    <div className="mt-6 flex space-x-3">
                         <Button size="sm" variant="outline" className="flex-1" onClick={() => navigate(`${AppRoute.Editor}?id=${resume.id}`)}>
                             Edit
                         </Button>
                         <Button size="sm" variant="ghost" className="flex-1 text-gray-600" title="View Public Link">
                             <Eye size={16} />
                         </Button>
                    </div>
                </div>
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
            <span className="font-medium text-gray-600 group-hover:text-blue-600">Create New Resume</span>
        </div>
      </div>
    </div>
  );
};
