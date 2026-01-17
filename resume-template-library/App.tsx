
import React, { useState } from 'react';
import { Search, Filter, ChevronRight, LayoutGrid, ListFilter, Sparkles, TrendingUp, History } from 'lucide-react';
import { JOB_CATEGORIES, TEMPLATES } from './data/mockData';
import ResumeCard from './components/ResumeCard';
import JobMegaMenu from './components/JobMegaMenu';

const App: React.FC = () => {
  const [activeCategory, setActiveCategory] = useState('all');
  const [hoveredCategory, setHoveredCategory] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  return (
    <div className="min-h-screen bg-[#f8fafc] text-slate-900 font-sans selection:bg-blue-100 selection:text-blue-700">
      {/* Global Header */}
      <header className="bg-white/80 backdrop-blur-md border-b border-slate-200 sticky top-0 z-40 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-gradient-to-br from-blue-500 to-indigo-600 text-white p-2 rounded-xl shadow-lg shadow-blue-200">
            <LayoutGrid size={24} />
          </div>
          <div>
            <h1 className="text-xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-slate-900 to-slate-600">简历模板库</h1>
            <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Premium Resume Builder</p>
          </div>
        </div>

        <div className="hidden lg:flex items-center gap-8">
          <nav className="flex items-center gap-6 text-sm font-semibold">
            <a href="#" className="text-blue-600 relative after:content-[''] after:absolute after:bottom-[-4px] after:left-0 after:w-full after:h-0.5 after:bg-blue-600 after:rounded-full">模板市场</a>
            <a href="#" className="text-slate-500 hover:text-slate-900 transition-colors">简历攻略</a>
            <a href="#" className="text-slate-500 hover:text-slate-900 transition-colors flex items-center gap-1.5">
              <Sparkles size={14} className="text-amber-500" />
              AI 优化
            </a>
          </nav>
          <div className="h-6 w-px bg-slate-200 mx-2"></div>
          <div className="flex items-center gap-3">
            <button className="text-slate-600 text-sm font-semibold hover:text-slate-900 px-4 py-2">登录</button>
            <button className="bg-slate-900 text-white px-6 py-2.5 rounded-full text-sm font-bold hover:bg-slate-800 transition-all shadow-lg shadow-slate-200 active:scale-95">
              免费注册
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-[1440px] mx-auto px-6 py-8 flex gap-8">
        
        {/* Left Sidebar */}
        <aside className="w-64 flex-shrink-0 relative hidden md:block">
          <div className="bg-white rounded-3xl border border-slate-200 p-4 shadow-sm sticky top-28">
            <div className="flex items-center gap-2 mb-6 px-3 text-slate-400">
              <ListFilter size={16} strokeWidth={2.5} />
              <span className="text-[11px] font-extrabold uppercase tracking-widest">岗位分类导航</span>
            </div>
            
            <nav className="space-y-1.5">
              {JOB_CATEGORIES.map((category) => (
                <div 
                  key={category.id}
                  className="relative group"
                  onMouseEnter={() => setHoveredCategory(category.id)}
                  onMouseLeave={() => setHoveredCategory(null)}
                >
                  <button
                    onClick={() => setActiveCategory(category.id)}
                    className={`w-full flex items-center justify-between px-4 py-3.5 rounded-2xl text-sm font-bold transition-all duration-200 border-2 ${
                      activeCategory === category.id 
                        ? 'bg-blue-50/50 text-blue-600 border-blue-100 shadow-sm' 
                        : 'text-slate-600 hover:bg-slate-50 border-transparent hover:border-slate-100'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <span className={`${activeCategory === category.id ? 'text-blue-600' : 'text-slate-400 group-hover:text-blue-500'} transition-colors`}>
                        {category.icon}
                      </span>
                      {category.name}
                    </div>
                    {category.subCategories && (
                      <ChevronRight size={16} className={`transition-all ${hoveredCategory === category.id ? 'translate-x-1 text-blue-400' : 'text-slate-300 opacity-0 group-hover:opacity-100'}`} />
                    )}
                  </button>

                  {/* Mega Menu Flyout */}
                  {category.subCategories && (
                    <JobMegaMenu 
                      isVisible={hoveredCategory === category.id} 
                      subCategories={category.subCategories} 
                    />
                  )}
                </div>
              ))}
            </nav>

            <div className="mt-8 pt-6 border-t border-slate-100">
              <div className="flex flex-col gap-2">
                <button className="flex items-center gap-3 px-4 py-3 text-slate-500 hover:text-slate-900 transition-colors">
                  <TrendingUp size={18} />
                  <span className="text-xs font-bold">热门推荐</span>
                </button>
                <button className="flex items-center gap-3 px-4 py-3 text-slate-500 hover:text-slate-900 transition-colors">
                  <History size={18} />
                  <span className="text-xs font-bold">最近使用</span>
                </button>
              </div>
            </div>
          </div>
        </aside>

        {/* Content Area */}
        <section className="flex-1">
          {/* Search & Filter Header */}
          <div className="mb-10 space-y-6">
            <div className="flex flex-col md:flex-row md:items-end justify-between gap-6">
              <div className="max-w-xl flex-1">
                <h2 className="text-3xl font-black text-slate-900 mb-2">发现您的完美简历</h2>
                <p className="text-slate-500 font-medium">从 1000+ 个专业设计的模板中挑选，助您脱颖而出。</p>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex bg-white p-1 rounded-xl border border-slate-200 shadow-sm">
                  <button className="px-4 py-1.5 rounded-lg bg-slate-100 text-slate-900 text-xs font-bold">最新</button>
                  <button className="px-4 py-1.5 rounded-lg text-slate-400 hover:text-slate-600 text-xs font-bold">最热</button>
                </div>
                <button className="p-2.5 bg-white border border-slate-200 rounded-xl text-slate-400 hover:text-blue-600 hover:border-blue-100 hover:bg-blue-50 transition-all shadow-sm">
                  <Filter size={20} />
                </button>
              </div>
            </div>

            <div className="relative group">
              <div className="absolute inset-y-0 left-5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-500 transition-colors">
                <Search size={22} strokeWidth={2.5} />
              </div>
              <input
                type="text"
                placeholder="搜索模板名称、岗位、行业关键词..."
                className="w-full bg-white border-2 border-slate-200 rounded-[2rem] py-5 pl-14 pr-6 text-lg font-medium outline-none focus:border-blue-500 focus:shadow-2xl focus:shadow-blue-100/50 transition-all placeholder:text-slate-300"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button className="absolute right-3 top-2 bottom-2 bg-blue-600 text-white px-8 rounded-full font-bold text-sm hover:bg-blue-700 transition-all active:scale-95 shadow-lg shadow-blue-100">
                搜索
              </button>
            </div>
            
            {/* Quick Filter Tags */}
            <div className="flex flex-wrap gap-2 pt-2">
              <span className="text-xs font-bold text-slate-400 uppercase tracking-widest mr-2 py-1">热门搜索:</span>
              {['互联网', '应届生', '简约', '管理岗', '金融', '设计'].map(tag => (
                <button key={tag} className="px-4 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-full text-xs font-bold transition-colors">
                  {tag}
                </button>
              ))}
            </div>
          </div>

          {/* Templates Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-8">
            {TEMPLATES.map((template) => (
              <ResumeCard key={template.id} template={template} />
            ))}
          </div>

          {/* Load More Section */}
          <div className="mt-16 text-center">
            <button className="inline-flex items-center gap-2 text-slate-400 hover:text-slate-600 font-bold text-sm transition-colors group">
              加载更多模板
              <div className="w-1.5 h-1.5 bg-slate-300 rounded-full group-hover:bg-slate-400 transition-colors"></div>
              <div className="w-1.5 h-1.5 bg-slate-300 rounded-full group-hover:bg-slate-400 transition-colors"></div>
              <div className="w-1.5 h-1.5 bg-slate-300 rounded-full group-hover:bg-slate-400 transition-colors"></div>
            </button>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="mt-20 border-t border-slate-200 bg-white py-12 px-6">
        <div className="max-w-[1440px] mx-auto flex flex-col md:flex-row justify-between items-center gap-8">
          <div className="flex items-center gap-3">
            <div className="bg-slate-100 p-2 rounded-lg grayscale opacity-50">
              <LayoutGrid size={20} />
            </div>
            <span className="text-slate-400 font-bold">© 2024 ResumeTemplateLib. All rights reserved.</span>
          </div>
          <div className="flex gap-10 text-slate-400 text-xs font-bold uppercase tracking-widest">
            <a href="#" className="hover:text-slate-900 transition-colors">隐私政策</a>
            <a href="#" className="hover:text-slate-900 transition-colors">服务条款</a>
            <a href="#" className="hover:text-slate-900 transition-colors">联系我们</a>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;
