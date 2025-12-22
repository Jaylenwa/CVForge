
import React from 'react';
import { Template } from '../types';

interface StatsProps {
  templates: Template[];
}

export const StatsSection: React.FC<StatsProps> = ({ templates }) => {
  const totalTemplates = templates.length;
  const premiumTemplates = templates.filter(t => t.isPremium).length;
  const avgPopularity = Math.round(templates.reduce((acc, t) => acc + t.popularity, 0) / (totalTemplates || 1));
  const topCategories = Array.from(new Set(templates.map(t => t.industry))).length;

  const stats = [
    { label: 'Total Templates', value: totalTemplates, icon: '📄', color: 'indigo' },
    { label: 'Premium Assets', value: premiumTemplates, icon: '💎', color: 'amber' },
    { label: 'Avg. Popularity', value: `${avgPopularity}%`, icon: '🔥', color: 'rose' },
    { label: 'Industries', value: topCategories, icon: '🏢', color: 'emerald' },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat, i) => (
        <div key={i} className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
          <div className="flex items-center gap-4">
            <div className={`w-12 h-12 flex items-center justify-center text-2xl rounded-xl bg-${stat.color}-50 text-${stat.color}-600`}>
              {stat.icon}
            </div>
            <div>
              <p className="text-sm font-medium text-slate-500">{stat.label}</p>
              <p className="text-2xl font-bold text-slate-900">{stat.value}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
