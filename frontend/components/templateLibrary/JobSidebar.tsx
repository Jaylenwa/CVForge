import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
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
  language?: string;
  categories: JobSidebarCategory[];
  roles: JobSidebarRole[];
  selectedCategoryId: string;
  onSelectCategory: (categoryId: string) => void;
  onSelectRole: (roleId: string) => void;
}> = ({ title, language, categories, roles, selectedCategoryId, onSelectCategory, onSelectRole }) => {
  const sortLocale = language === 'en' ? 'en' : 'zh';
  const [hoveredCategory, setHoveredCategory] = useState<string | null>(null);
  const [flyoutStyle, setFlyoutStyle] = useState<React.CSSProperties | undefined>(undefined);
  const buttonRefs = useRef<Record<string, HTMLButtonElement | null>>({});
  const hoverCloseTimerRef = useRef<number | null>(null);

  const computeFlyoutStyleForCategory = useCallback((categoryId: string) => {
    const btn = buttonRefs.current[categoryId];
    const rect = btn?.getBoundingClientRect();
    if (!rect) return undefined;

    const viewportPadding = 16;
    const menuWidth = 640;
    const minMenuHeight = 240;
    const leftGap = 16;
    const desiredLeft = rect.right + leftGap;
    const left = Math.min(
      Math.max(viewportPadding, desiredLeft),
      Math.max(viewportPadding, window.innerWidth - viewportPadding - menuWidth)
    );

    const topOffset = 40;
    const desiredTop = rect.top - topOffset;
    const top = Math.max(
      viewportPadding,
      Math.min(desiredTop, window.innerHeight - viewportPadding - minMenuHeight)
    );
    const maxHeight = Math.max(minMenuHeight, window.innerHeight - viewportPadding - top);

    return {
      position: 'fixed',
      left,
      top,
      width: menuWidth,
      maxHeight,
      overflowY: 'auto',
      overscrollBehavior: 'contain',
    } as React.CSSProperties;
  }, []);

  const cancelCloseTimer = useCallback(() => {
    if (hoverCloseTimerRef.current) {
      window.clearTimeout(hoverCloseTimerRef.current);
      hoverCloseTimerRef.current = null;
    }
  }, []);

  const scheduleClose = useCallback(() => {
    cancelCloseTimer();
    hoverCloseTimerRef.current = window.setTimeout(() => {
      setHoveredCategory(null);
      setFlyoutStyle(undefined);
      hoverCloseTimerRef.current = null;
    }, 120);
  }, [cancelCloseTimer]);

  const handleSelectRole = useCallback(
    (roleId: string) => {
      cancelCloseTimer();
      onSelectRole(roleId);
      setHoveredCategory(null);
      setFlyoutStyle(undefined);
    },
    [cancelCloseTimer, onSelectRole]
  );

  const updateFlyoutStyle = useCallback(
    (categoryId?: string) => {
      const id = categoryId ?? hoveredCategory;
      if (!id) {
        setFlyoutStyle(undefined);
        return;
      }
      setFlyoutStyle(computeFlyoutStyleForCategory(id));
    },
    [computeFlyoutStyleForCategory, hoveredCategory]
  );

  useEffect(() => {
    updateFlyoutStyle();
    if (!hoveredCategory) return;

    const onScroll = () => updateFlyoutStyle();
    const onResize = () => updateFlyoutStyle();
    document.addEventListener('scroll', onScroll, true);
    window.addEventListener('resize', onResize);
    return () => {
      document.removeEventListener('scroll', onScroll, true);
      window.removeEventListener('resize', onResize);
    };
  }, [hoveredCategory, updateFlyoutStyle]);

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
        .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || a.name.localeCompare(b.name, sortLocale))
        .map((r) => ({ id: r.id, name: r.name } as JobMegaMenuRole));
      if (roleList.length) groups.push({ id: child.id, title: child.name, roles: roleList });
    }
    const rootRoles = (rolesByCategory[rootCategoryId] || [])
      .slice()
      .sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0) || a.name.localeCompare(b.name, sortLocale))
      .map((r) => ({ id: r.id, name: r.name } as JobMegaMenuRole));
    if (rootRoles.length) groups.push({ id: `${rootCategoryId}__other`, title: '其它', roles: rootRoles });
    return groups;
  };

  return (
    <aside className="w-60 flex-shrink-0 relative hidden md:block z-40">
      <div className="bg-white rounded-3xl border border-slate-200 p-4 shadow-sm relative z-40 isolate">
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
                onMouseEnter={() => {
                  cancelCloseTimer();
                  updateFlyoutStyle(category.id);
                  setHoveredCategory(category.id);
                }}
                onMouseLeave={scheduleClose}
              >
                <button
                  ref={(el) => {
                    buttonRefs.current[category.id] = el;
                  }}
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
                  onSelectRole={handleSelectRole}
                  style={showFlyout ? flyoutStyle : undefined}
                  onMouseEnter={cancelCloseTimer}
                  onMouseLeave={scheduleClose}
                />
              </div>
            );
          })}
        </nav>
      </div>
    </aside>
  );
};
