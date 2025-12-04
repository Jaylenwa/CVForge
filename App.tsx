import React from 'react';
import { HashRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import { MainLayout } from './components/Layout';
import { Home } from './pages/Home';
import { Templates } from './pages/Templates';
import { Dashboard } from './pages/Dashboard';
import { Editor } from './pages/editor/Editor';
import { AppRoute } from './types';

// Auth Login Placeholder
const Login = () => (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8 text-center">
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">Sign in to your account</h2>
            <p className="text-gray-500">Demo Mode: No login required for this preview.</p>
            <a href="/" className="font-medium text-blue-600 hover:text-blue-500">Return Home</a>
        </div>
    </div>
);

// Pricing Placeholder
const Pricing = () => (
     <div className="py-12 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center">
                <h2 className="text-3xl font-extrabold text-gray-900 sm:text-4xl">Pricing Plans</h2>
                <p className="mt-4 text-lg text-gray-500">Simple, transparent pricing for everyone.</p>
            </div>
             <div className="mt-12 space-y-4 sm:mt-16 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-6 lg:max-w-4xl lg:mx-auto xl:max-w-none xl:mx-0 xl:grid-cols-2">
                 {['Basic', 'Pro'].map((plan) => (
                     <div key={plan} className="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
                         <div className="p-6">
                             <h2 className="text-lg leading-6 font-medium text-gray-900">{plan}</h2>
                             <p className="mt-4 text-sm text-gray-500">Best for individuals.</p>
                             <p className="mt-8">
                                 <span className="text-4xl font-extrabold text-gray-900">{plan === 'Basic' ? '$0' : '$12'}</span>
                                 <span className="text-base font-medium text-gray-500">/mo</span>
                             </p>
                             <button className="mt-8 block w-full bg-blue-600 border border-transparent rounded-md py-2 text-sm font-semibold text-white text-center hover:bg-blue-700">Buy {plan}</button>
                         </div>
                     </div>
                 ))}
             </div>
        </div>
     </div>
)

const LayoutWrapper = () => (
    <MainLayout>
        <Outlet />
    </MainLayout>
);

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        {/* Routes with Main Navbar/Footer */}
        <Route element={<LayoutWrapper />}>
            <Route path={AppRoute.Home} element={<Home />} />
            <Route path={AppRoute.Templates} element={<Templates />} />
            <Route path={AppRoute.Dashboard} element={<Dashboard />} />
            <Route path={AppRoute.Pricing} element={<Pricing />} />
        </Route>
        
        {/* Standalone Routes */}
        <Route path={AppRoute.Editor} element={<Editor />} />
        <Route path={AppRoute.Login} element={<Login />} />
      </Routes>
    </Router>
  );
};

export default App;
