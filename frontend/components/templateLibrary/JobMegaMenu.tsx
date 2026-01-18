import React from 'react';
import { AnimatePresence, motion } from 'framer-motion';

export type JobMegaMenuRole = { id: string; name: string };

export type JobMegaMenuGroup = { id: string; title: string; roles: JobMegaMenuRole[] };

export const JobMegaMenu: React.FC<{
  isVisible: boolean;
  title?: string;
  roles?: JobMegaMenuRole[];
  groups?: JobMegaMenuGroup[];
  onSelectRole: (roleId: string) => void;
}> = ({ isVisible, title, roles, groups, onSelectRole }) => {
  const normalizedGroups: JobMegaMenuGroup[] =
    groups && groups.length
      ? groups
      : [
          {
            id: 'all',
            title: title || '',
            roles: roles || [],
          },
        ];
  const hasRoles = normalizedGroups.some(g => g.roles.length > 0);
  return (
    <AnimatePresence>
      {isVisible && hasRoles ? (
        <motion.div
          initial={{ opacity: 0, x: -8, y: -2 }}
          animate={{ opacity: 1, x: 0, y: 0 }}
          exit={{ opacity: 0, x: -8, y: -2 }}
          transition={{ duration: 0.16, ease: 'easeOut' }}
          className="absolute left-full top-0 ml-4 w-[640px] bg-white rounded-2xl shadow-2xl border border-slate-100 z-40 p-8"
        >
          <div className="space-y-6">
            {title ? (
              <h4 className="text-base font-bold text-slate-900 flex items-center">
                {title}
                <div className="h-px flex-1 bg-slate-100 ml-4" />
              </h4>
            ) : null}
            {normalizedGroups.map((g) =>
              g.roles.length ? (
                <div key={g.id} className="space-y-3">
                  {g.title ? <div className="text-sm font-bold text-slate-900">{g.title}</div> : null}
                  <div className="flex flex-wrap gap-x-4 gap-y-2">
                    {g.roles.map((r) => (
                      <button
                        key={r.id}
                        onClick={() => onSelectRole(r.id)}
                        className="text-sm text-left text-slate-600 hover:text-blue-600 transition-colors"
                      >
                        {r.name}
                      </button>
                    ))}
                  </div>
                </div>
              ) : null
            )}
          </div>
        </motion.div>
      ) : null}
    </AnimatePresence>
  );
};
