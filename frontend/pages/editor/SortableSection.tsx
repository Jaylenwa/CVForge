import React, { useState, useEffect, useRef } from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { GripVertical, ChevronDown, ChevronUp, Trash2, Edit2 } from 'lucide-react';
import { ResumeSection } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

interface SortableSectionProps {
  section: ResumeSection;
  isActive: boolean;
  onToggle: () => void;
  onUpdate: (updates: Partial<ResumeSection>) => void;
  onRemove: () => void;
  children: React.ReactNode;
  icon?: React.ReactNode;
  onAddItem?: () => void;
}

export const SortableSection: React.FC<SortableSectionProps> = ({
  section,
  isActive,
  onToggle,
  onUpdate,
  onRemove,
  children,
  icon,
  onAddItem
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({ id: section.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    zIndex: isDragging ? 10 : 1,
    opacity: isDragging ? 0.5 : 1,
  };

  const { t } = useLanguage();
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [editedTitle, setEditedTitle] = useState(section.title);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditingTitle && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isEditingTitle]);

  const getDisplayTitle = () => {
    const defaultTitles = ['Professional Summary', 'Work Experience', 'Education', 'Skills', 'Projects'];
    if (defaultTitles.includes(section.title) || !section.title) {
        return t(`section.${section.type}`);
    }
    return section.title;
  };

  const handleSaveTitle = (e?: React.FormEvent) => {
    e?.preventDefault();
    if (editedTitle.trim()) {
        onUpdate({ title: editedTitle });
    }
    setIsEditingTitle(false);
  };

  const handleCancelEdit = () => {
    setEditedTitle(section.title);
    setIsEditingTitle(false);
  };

  return (
    <div ref={setNodeRef} style={style} className="group rounded-2xl overflow-hidden border border-gray-200 bg-white shadow-sm mb-4">
      <div
        className={`flex items-center justify-between px-5 py-4 ${isActive ? 'bg-gray-50' : 'bg-white'}`}
        onClick={onToggle}
      >
        <div className="flex items-center flex-grow mr-4">
          <button
            className="mr-3 text-gray-400 hover:text-gray-600 cursor-grab active:cursor-grabbing p-1 rounded hover:bg-gray-200 transition-colors opacity-0 group-hover:opacity-100"
            {...attributes}
            {...listeners}
            onMouseDown={(e) => e.stopPropagation()}
            onClick={(e) => e.stopPropagation()}
            aria-label={t('a11y.drag')}
          >
            <GripVertical size={16} />
          </button>
          <div className="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center mr-3">
            {icon}
          </div>
          <div className="flex-grow flex items-center">
            {isEditingTitle ? (
              <form onSubmit={handleSaveTitle} className="flex items-center w-full max-w-md">
                <input
                  ref={inputRef}
                  value={editedTitle}
                  onChange={(e) => setEditedTitle(e.target.value)}
                  className="flex-grow px-2 py-1 text-sm border border-indigo-500 rounded focus:outline-none focus:ring-2 focus:ring-indigo-200"
                  onBlur={() => handleSaveTitle()}
                  onKeyDown={(e) => {
                    if (e.key === 'Escape') handleCancelEdit();
                  }}
                />
              </form>
            ) : (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onToggle();
                }}
                className="font-semibold text-gray-800 hover:text-indigo-600 transition-colors text-left flex items-center w-full"
              >
                {getDisplayTitle()}
                <span
                  onClick={(e) => {
                    e.stopPropagation();
                    setIsEditingTitle(true);
                    setEditedTitle(section.title || getDisplayTitle());
                  }}
                  className="ml-2 text-gray-400 hover:text-indigo-500 p-1 rounded transition-all"
                  title={t('editor.rename')}
                >
                  <Edit2 size={12} />
                </span>
              </button>
            )}
          </div>
        </div>
        <div className="flex items-center space-x-2">
          {onAddItem && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onAddItem();
              }}
              className="text-xs bg-indigo-50 text-indigo-600 px-3 py-1 rounded-lg hover:bg-indigo-100 transition-colors"
            >
              {t('editor.addItem')}
            </button>
          )}
          <button
            onClick={(e) => {
              e.stopPropagation();
              onRemove();
            }}
            className="p-2 text-gray-400 hover:text-red-500 rounded-full hover:bg-red-50 transition-colors"
            title={t('editor.removeSection')}
          >
            <Trash2 size={16} />
          </button>
          <input
            type="checkbox"
            checked={section.isVisible}
            onChange={(e) => onUpdate({ isVisible: e.target.checked })}
            className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded cursor-pointer"
            title={t('editor.toggleVisibility')}
            onClick={(e) => e.stopPropagation()}
          />
          <button
            onClick={(e) => {
              e.stopPropagation();
              onToggle();
            }}
            className="p-1 ml-1 text-gray-500 hover:text-gray-700 transition-colors"
          >
            {isActive ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </button>
        </div>
      </div>
      {isActive && (
        <div className="px-6 pt-4 pb-6 space-y-6 bg-gray-50/50 border-t border-gray-100">
          {children}
        </div>
      )}
    </div>
  );
};
