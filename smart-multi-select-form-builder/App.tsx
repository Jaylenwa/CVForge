
import React, { useState } from 'react';
import { Layout, Send, Info, Layers, Zap, ShieldCheck } from 'lucide-react';
import MultiSelect from './components/MultiSelect';
import { Option, FormData } from './types';

const TECH_OPTIONS: Option[] = [
  { value: 'react', label: 'React.js' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'tailwind', label: 'Tailwind CSS' },
  { value: 'node', label: 'Node.js' },
  { value: 'python', label: 'Python' },
  { value: 'rust', label: 'Rust' },
  { value: 'go', label: 'Go' },
  { value: 'docker', label: 'Docker' },
  { value: 'kubernetes', label: 'Kubernetes' },
  { value: 'aws', label: 'AWS' },
  { value: 'firebase', label: 'Firebase' },
  { value: 'supabase', label: 'Supabase' },
  { value: 'graphql', label: 'GraphQL' },
  { value: 'postgre', label: 'PostgreSQL' },
  { value: 'redis', label: 'Redis' },
];

const App: React.FC = () => {
  const [formData, setFormData] = useState<FormData>({
    projectName: '',
    description: '',
    technologies: [],
    priority: 'medium'
  });

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    // Simulate API call
    setTimeout(() => {
      setIsSubmitting(false);
      setSuccess(true);
      setTimeout(() => setSuccess(false), 5000);
    }, 1500);
  };

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col items-center py-12 px-4">
      <div className="w-full max-w-2xl bg-white rounded-2xl shadow-xl shadow-slate-200/60 border border-slate-100 overflow-hidden">
        
        {/* Header */}
        <div className="bg-indigo-600 px-8 py-10 text-white relative overflow-hidden">
          <div className="relative z-10">
            <h1 className="text-3xl font-bold flex items-center gap-2">
              <Layers className="text-indigo-200" />
              Project Blueprint
            </h1>
            <p className="text-indigo-100 mt-2 text-sm max-w-md">
              Configure your next big idea. Fill in the details below to initialize your project architecture.
            </p>
          </div>
          {/* Decorative background element */}
          <div className="absolute -right-12 -top-12 w-48 h-48 bg-white/10 rounded-full blur-3xl"></div>
          <div className="absolute -left-12 -bottom-12 w-48 h-48 bg-indigo-400/20 rounded-full blur-3xl"></div>
        </div>

        {/* Form Body */}
        <form onSubmit={handleSubmit} className="p-8 space-y-6">
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-slate-700 flex items-center gap-1.5">
                Project Name
                <Info size={14} className="text-slate-400" />
              </label>
              <input 
                type="text"
                required
                className="w-full px-4 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
                placeholder="Ex: Galactic Voyager"
                value={formData.projectName}
                onChange={(e) => setFormData({...formData, projectName: e.target.value})}
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-slate-700 flex items-center gap-1.5">
                Priority Level
              </label>
              <div className="flex p-1 bg-slate-100 rounded-lg border border-slate-200">
                {(['low', 'medium', 'high'] as const).map((p) => (
                  <button
                    key={p}
                    type="button"
                    onClick={() => setFormData({...formData, priority: p})}
                    className={`
                      flex-1 py-1 text-xs font-bold capitalize rounded-md transition-all
                      ${formData.priority === p 
                        ? 'bg-white text-indigo-600 shadow-sm' 
                        : 'text-slate-500 hover:text-slate-700'
                      }
                    `}
                  >
                    {p}
                  </button>
                ))}
              </div>
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-slate-700">Project Overview</label>
            <textarea 
              rows={3}
              className="w-full px-4 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all resize-none"
              placeholder="Briefly describe the mission objectives..."
              value={formData.description}
              onChange={(e) => setFormData({...formData, description: e.target.value})}
            />
          </div>

          {/* THE MULTI-SELECT COMPONENT */}
          <MultiSelect 
            label="Tech Stack (Choose as many as needed)"
            options={TECH_OPTIONS}
            value={formData.technologies}
            onChange={(vals) => setFormData({...formData, technologies: vals})}
            placeholder="Search and select technologies..."
          />

          <div className="pt-4 border-t border-slate-100">
            <button
              type="submit"
              disabled={isSubmitting}
              className={`
                w-full py-3.5 px-6 rounded-xl text-white font-bold flex items-center justify-center gap-2 transition-all
                ${isSubmitting 
                  ? 'bg-indigo-400 cursor-not-allowed' 
                  : 'bg-indigo-600 hover:bg-indigo-700 shadow-lg shadow-indigo-200 hover:-translate-y-0.5 active:translate-y-0'
                }
              `}
            >
              {isSubmitting ? (
                <>
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                  Generating Blueprint...
                </>
              ) : (
                <>
                  <Send size={18} />
                  Deploy Configuration
                </>
              )}
            </button>
          </div>

          {success && (
            <div className="p-4 bg-emerald-50 border border-emerald-100 rounded-lg flex items-start gap-3 animate-in slide-in-from-top-4 duration-300">
              <div className="bg-emerald-100 p-1.5 rounded-full text-emerald-600">
                <ShieldCheck size={18} />
              </div>
              <div>
                <h4 className="text-sm font-bold text-emerald-900">Project saved successfully!</h4>
                <p className="text-xs text-emerald-700 mt-0.5">Your configurations have been synchronized with the cloud.</p>
              </div>
            </div>
          )}
        </form>

        {/* Footer Quick Info */}
        <div className="px-8 py-6 bg-slate-50 border-t border-slate-100 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex items-center text-slate-500 gap-1.5">
              <Zap size={14} className="text-amber-400" />
              <span className="text-xs font-medium">Fast Build</span>
            </div>
            <div className="flex items-center text-slate-500 gap-1.5">
              <ShieldCheck size={14} className="text-emerald-500" />
              <span className="text-xs font-medium">Auto-Safe</span>
            </div>
          </div>
          <span className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">v2.4.0 Engine</span>
        </div>
      </div>

      <p className="mt-8 text-slate-400 text-sm flex items-center gap-2">
        <Layout size={14} />
        Built with React and Tailwind CSS
      </p>
    </div>
  );
};

export default App;
