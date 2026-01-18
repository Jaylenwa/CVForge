import React, { useMemo, useState } from 'react';
import { Briefcase, Code, Palette, Settings, User, ChevronRight, ListFilter } from 'lucide-react';
import { JobMegaMenu, JobMegaMenuGroup, JobMegaMenuRole } from './JobMegaMenu';

export type JobSidebarCategory = { id: string; name: string; parentId?: string; orderNum?: number };
export type JobSidebarRole = { id: string; name: string; categoryId: string; orderNum?: number };

const categoryIcon = (id: string) => {
  switch (id) {
    case 'tech':
      return <Code size={18} />;
    case 'product':
      return <Palette size={18} />;
    case 'ops':
      return <Settings size={18} />;
    case 'business':
      return <Briefcase size={18} />;
    case 'all':
    default:
      return <User size={18} />;
  }
};

export const JobSidebar: React.FC<{
  title: string;
  categories: JobSidebarCategory[];
  roles: JobSidebarRole[];
  selectedCategoryId: string;
  onSelectCategory: (categoryId: string) => void;
  onSelectRole: (roleId: string) => void;
}> = ({ title, categories, roles, selectedCategoryId, onSelectCategory, onSelectRole }) => {
  const [hoveredCategory, setHoveredCategory] = useState<string | null>(null);

  const rolesByCategory = useMemo(() => {
    const m: Record<string, JobSidebarRole[]> = {};
    for (const r of roles) {
      if (!m[r.categoryId]) m[r.categoryId] = [];
      m[r.categoryId].push(r);
    }
    return m;
  }, [roles]);

  const childrenByParent = useMemo(() => {
    const m: Record<string, JobSidebarCategory[]> = {};
    for (const c of categories) {
      const parent = c.parentId || '';
      if (!m[parent]) m[parent] = [];
      m[parent].push(c);
    }
    for (const k of Object.keys(m)) {
      m[k] = m[k].slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
    }
    return m;
  }, [categories]);

  const sortedRootCategories = useMemo(() => {
    return (childrenByParent[''] || []).slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
  }, [childrenByParent]);

  const buildFlyoutGroups = (rootCategoryId: string): JobMegaMenuGroup[] => {
    const groups: JobMegaMenuGroup[] = [];
    const childCategories = childrenByParent[rootCategoryId] || [];
    for (const child of childCategories) {
      const roleList = (rolesByCategory[child.id] || [])
        .slice()
        .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || a.name.localeCompare(b.name, 'zh'))
        .map((r) => ({ id: r.id, name: r.name } as JobMegaMenuRole));
      if (roleList.length) groups.push({ id: child.id, title: child.name, roles: roleList });
    }
    const rootRoles = (rolesByCategory[rootCategoryId] || [])
      .slice()
      .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || a.name.localeCompare(b.name, 'zh'))
      .map((r) => ({ id: r.id, name: r.name } as JobMegaMenuRole));
    if (rootRoles.length) groups.push({ id: `${rootCategoryId}__other`, title: '其它', roles: rootRoles });
    return groups;
  };

  return (
    <aside className="w-64 flex-shrink-0 relative hidden md:block z-50">
      <div className="bg-white rounded-3xl border border-slate-200 p-4 shadow-sm sticky top-28 relative z-50 isolate">
        <div className="flex items-center gap-2 mb-6 px-3 text-slate-400">
          <ListFilter size={16} strokeWidth={2.5} />
          <span className="text-[11px] font-extrabold uppercase tracking-widest">{title}</span>
        </div>

        <nav className="space-y-1.5">
          {sortedRootCategories.map((category) => {
            const isActive = selectedCategoryId === category.id;
            const flyoutGroups = buildFlyoutGroups(category.id);
            const showFlyout = hoveredCategory === category.id && flyoutGroups.some(g => g.roles.length > 0);
            return (
              <div
                key={category.id}
                className="relative group"
                onMouseEnter={() => setHoveredCategory(category.id)}
                onMouseLeave={() => setHoveredCategory(null)}
              >
                <button
                  onClick={() => onSelectCategory(category.id)}
                  className={`w-full flex items-center justify-between px-4 py-3.5 rounded-2xl text-sm font-bold transition-all duration-200 border-2 ${
                    isActive
                      ? 'bg-blue-50/50 text-blue-600 border-blue-100 shadow-sm'
                      : 'text-slate-600 hover:bg-slate-50 border-transparent hover:border-slate-100'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <span className={`${isActive ? 'text-blue-600' : 'text-slate-400 group-hover:text-blue-500'} transition-colors`}>
                      {categoryIcon(category.id)}
                    </span>
                    {category.name}
                  </div>
                  {flyoutGroups.length > 0 ? (
                    <ChevronRight
                      size={16}
                      className={`transition-all ${showFlyout ? 'translate-x-1 text-blue-400' : 'text-slate-300 opacity-0 group-hover:opacity-100'}`}
                    />
                  ) : null}
                </button>

                <JobMegaMenu
                  isVisible={showFlyout}
                  title={category.name}
                  groups={flyoutGroups}
                  onSelectRole={onSelectRole}
                />
              </div>
            );
          })}
        </nav>
      </div>
    </aside>
  );
};
