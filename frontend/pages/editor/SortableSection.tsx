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
}

export const SortableSection: React.FC<SortableSectionProps> = ({
  section,
  isActive,
  onToggle,
  onUpdate,
  onRemove,
  children
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
    <div ref={setNodeRef} style={style} className="border-b border-gray-200 bg-white">
       <div className={`flex items-center justify-between px-6 py-4 hover:bg-gray-50 ${isActive ? 'bg-gray-50' : ''}`}>
         <div className="flex items-center flex-grow mr-4">
             {/* Drag Handle */}
             <button 
                className="mr-3 text-gray-400 hover:text-gray-600 cursor-grab active:cursor-grabbing p-1 rounded hover:bg-gray-200 transition-colors"
                {...attributes} 
                {...listeners}
             >
                 <GripVertical size={16} />
             </button>

             {/* Title */}
             <div className="flex-grow flex items-center">
                 {isEditingTitle ? (
                     <form onSubmit={handleSaveTitle} className="flex items-center w-full max-w-md">
                         <input
                            ref={inputRef}
                            value={editedTitle}
                            onChange={(e) => setEditedTitle(e.target.value)}
                            className="flex-grow px-2 py-1 text-sm border border-blue-500 rounded focus:outline-none focus:ring-2 focus:ring-blue-200"
                            onBlur={() => handleSaveTitle()} 
                            onKeyDown={(e) => {
                                if (e.key === 'Escape') handleCancelEdit();
                            }}
                         />
                     </form>
                 ) : (
                     <button 
                        onClick={onToggle}
                        className="font-semibold text-gray-700 hover:text-blue-600 transition-colors text-left flex items-center group"
                     >
                        {getDisplayTitle()}
                        <span 
                            onClick={(e) => {
                                e.stopPropagation();
                                setIsEditingTitle(true);
                                setEditedTitle(section.title || getDisplayTitle());
                            }}
                            className="ml-2 opacity-0 group-hover:opacity-100 text-gray-400 hover:text-blue-500 p-1 rounded transition-all"
                            title={t('editor.rename')}
                        >
                            <Edit2 size={12} />
                        </span>
                     </button>
                 )}
             </div>
         </div>

         <div className="flex items-center space-x-1">
             <button
                onClick={(e) => { e.stopPropagation(); onRemove(); }}
                className="p-2 text-gray-400 hover:text-red-500 rounded-full hover:bg-red-50 transition-colors"
                title={t('editor.removeSection')}
             >
                <Trash2 size={16} />
             </button>

             <div className="h-4 w-px bg-gray-300 mx-2"></div>

             <input 
                type="checkbox" 
                checked={section.isVisible} 
                onChange={(e) => onUpdate({ isVisible: e.target.checked })}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded cursor-pointer"
                title={t('editor.toggleVisibility')}
             />
             
             <button 
                onClick={onToggle}
                className="p-1 ml-2 text-gray-500 hover:text-gray-700 transition-colors"
             >
                {isActive ? <ChevronUp size={16}/> : <ChevronDown size={16}/>}
             </button>
         </div>
       </div>

       {isActive && (
         <div className="px-6 pb-6 space-y-6 bg-gray-50/50 border-t border-gray-100 animate-fadeIn">
            {children}
         </div>
       )}
    </div>
  );
};
