import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { FileText, Grid, User, Menu, X, Star } from 'lucide-react';
import { Button } from './ui/Button';
import { AppRoute } from '../types';

export const Navbar: React.FC = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const location = useLocation();

  const navItems = [
    { label: 'Templates', href: AppRoute.Templates, icon: <Grid size={18} /> },
    { label: 'My Resumes', href: AppRoute.Dashboard, icon: <FileText size={18} /> },
    { label: 'Pricing', href: AppRoute.Pricing, icon: <Star size={18} /> },
  ];

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <Link to="/" className="flex-shrink-0 flex items-center">
              <span className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 text-transparent bg-clip-text">CVForge</span>
            </Link>
            <div className="hidden sm:ml-8 sm:flex sm:space-x-8">
              {navItems.map((item) => (
                <Link
                  key={item.href}
                  to={item.href}
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
                    location.pathname === item.href
                      ? 'border-blue-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </div>
          </div>
          <div className="hidden sm:ml-6 sm:flex sm:items-center sm:space-x-4">
            <Link to={AppRoute.Login}>
                <Button variant="ghost">Sign In</Button>
            </Link>
            <Link to={AppRoute.Templates}>
                <Button>Get Started</Button>
            </Link>
          </div>
          <div className="-mr-2 flex items-center sm:hidden">
            <button
              onClick={() => setIsOpen(!isOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
            >
              {isOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {isOpen && (
        <div className="sm:hidden bg-white border-b border-gray-200">
          <div className="pt-2 pb-3 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                to={item.href}
                onClick={() => setIsOpen(false)}
                className={`block pl-3 pr-4 py-2 border-l-4 text-base font-medium ${
                   location.pathname === item.href
                    ? 'bg-blue-50 border-blue-500 text-blue-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                }`}
              >
                <div className="flex items-center">
                    <span className="mr-3">{item.icon}</span>
                    {item.label}
                </div>
              </Link>
            ))}
             <div className="pt-4 pb-3 border-t border-gray-200">
                <div className="mt-3 space-y-1 px-2">
                     <Link to={AppRoute.Login} onClick={() => setIsOpen(false)}>
                        <Button className="w-full justify-center" variant="outline">Sign In</Button>
                     </Link>
                     <div className="h-2"></div>
                     <Link to={AppRoute.Templates} onClick={() => setIsOpen(false)}>
                        <Button className="w-full justify-center">Get Started</Button>
                     </Link>
                </div>
            </div>
          </div>
        </div>
      )}
    </nav>
  );
};

export const Footer: React.FC = () => {
    return (
        <footer className="bg-slate-900 text-white py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 grid grid-cols-1 md:grid-cols-4 gap-8">
                <div>
                    <h3 className="text-xl font-bold mb-4">CVForge</h3>
                    <p className="text-slate-400 text-sm">Empowering careers with professional, AI-enhanced resumes.</p>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">Product</h4>
                    <ul className="space-y-2 text-slate-400 text-sm">
                        <li><Link to="/templates" className="hover:text-white">Templates</Link></li>
                        <li><Link to="/pricing" className="hover:text-white">Pricing</Link></li>
                        <li><Link to="/features" className="hover:text-white">Features</Link></li>
                    </ul>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">Legal</h4>
                    <ul className="space-y-2 text-slate-400 text-sm">
                        <li>Privacy Policy</li>
                        <li>Terms of Service</li>
                        <li>Cookie Policy</li>
                    </ul>
                </div>
                <div>
                    <h4 className="font-semibold mb-4">Language</h4>
                    <div className="flex space-x-2">
                        <span className="text-sm text-slate-400 hover:text-white cursor-pointer">English</span>
                        <span className="text-sm text-slate-400">|</span>
                        <span className="text-sm text-slate-400 hover:text-white cursor-pointer">中文</span>
                    </div>
                </div>
            </div>
            <div className="max-w-7xl mx-auto px-4 mt-8 pt-8 border-t border-slate-800 text-center text-slate-500 text-sm">
                &copy; {new Date().getFullYear()} CVForge. All rights reserved.
            </div>
        </footer>
    )
}

export const MainLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar />
      <main className="flex-grow">
        {children}
      </main>
      <Footer />
    </div>
  );
};
