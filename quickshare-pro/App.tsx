
import React, { useState } from 'react';
import ShareModal from './components/ShareModal';

const App: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="max-w-4xl w-full bg-white rounded-3xl shadow-xl overflow-hidden border border-slate-100">
        <div className="relative h-48 bg-gradient-to-r from-blue-600 to-indigo-700 p-8 flex flex-col justify-end">
          <div className="absolute top-4 right-4 bg-white/20 backdrop-blur-md px-3 py-1 rounded-full text-xs text-white font-medium uppercase tracking-wider">
            Draft Project
          </div>
          <h1 className="text-3xl font-bold text-white">Q1 Performance Report.pdf</h1>
          <p className="text-blue-100 mt-2">Last edited 2 hours ago • 4.2 MB</p>
        </div>
        
        <div className="p-8">
          <div className="flex items-center justify-between mb-8">
            <div className="flex -space-x-2">
              {[1, 2, 3].map(i => (
                <img 
                  key={i}
                  src={`https://picsum.photos/seed/${i + 40}/64`} 
                  className="w-10 h-10 rounded-full border-2 border-white"
                  alt="User"
                />
              ))}
              <div className="w-10 h-10 rounded-full border-2 border-white bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-500">
                +5
              </div>
            </div>
            <button 
              onClick={() => setIsModalOpen(true)}
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2.5 rounded-xl font-semibold transition-all flex items-center gap-2 shadow-lg shadow-indigo-200"
            >
              <i className="fa-solid fa-share-nodes"></i>
              Share with others
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="p-4 rounded-2xl bg-slate-50 border border-slate-100">
              <p className="text-sm text-slate-500 mb-1">Total Views</p>
              <p className="text-2xl font-bold text-slate-800">1,284</p>
            </div>
            <div className="p-4 rounded-2xl bg-slate-50 border border-slate-100">
              <p className="text-sm text-slate-500 mb-1">Active Links</p>
              <p className="text-2xl font-bold text-slate-800">3</p>
            </div>
            <div className="p-4 rounded-2xl bg-slate-50 border border-slate-100">
              <p className="text-sm text-slate-500 mb-1">Last Download</p>
              <p className="text-2xl font-bold text-slate-800">Just now</p>
            </div>
          </div>
        </div>
      </div>

      <ShareModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        fileTitle="Q1 Performance Report.pdf"
      />

      <div className="mt-8 text-slate-400 text-sm flex items-center gap-4">
        <span>© 2024 QuickShare Pro</span>
        <span className="w-1 h-1 bg-slate-300 rounded-full"></span>
        <a href="#" className="hover:text-slate-600 transition-colors">Privacy</a>
        <span className="w-1 h-1 bg-slate-300 rounded-full"></span>
        <a href="#" className="hover:text-slate-600 transition-colors">Terms</a>
      </div>
    </div>
  );
};

export default App;
