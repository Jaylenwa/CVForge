import React, { createContext, useContext, useState, useRef, useCallback } from 'react';
import { Modal } from './Modal';
import { Button } from './Button';
import { useLanguage } from '../../contexts/LanguageContext';
import { HelpCircle } from 'lucide-react';

interface ConfirmOptions {
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'danger' | 'primary';
}

interface DialogContextType {
  confirm: (options: ConfirmOptions) => Promise<boolean>;
}

const DialogContext = createContext<DialogContextType | undefined>(undefined);

export const useConfirm = () => {
  const context = useContext(DialogContext);
  if (!context) {
    throw new Error('useConfirm must be used within a ConfirmDialogProvider');
  }
  return context.confirm;
};

export const ConfirmDialogProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [options, setOptions] = useState<ConfirmOptions>({ message: '' });
  const resolveRef = useRef<(value: boolean) => void>(() => {});
  const { t } = useLanguage();

  const confirm = useCallback((opts: ConfirmOptions) => {
    setOptions(opts);
    setIsOpen(true);
    return new Promise<boolean>((resolve) => {
      resolveRef.current = resolve;
    });
  }, []);

  const handleClose = () => {
    setIsOpen(false);
    resolveRef.current(false);
  };

  const handleConfirm = () => {
    setIsOpen(false);
    resolveRef.current(true);
  };

  return (
    <DialogContext.Provider value={{ confirm }}>
      {children}
      <Modal 
        isOpen={isOpen} 
        onClose={handleClose} 
        title={options.title || t('common.confirmAction')}
      >
        <div className="flex flex-col space-y-4">
          <div className="flex items-start">
             <div className={`flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 mr-4`}>
                <HelpCircle className="h-6 w-6 text-blue-600" aria-hidden="true" />
             </div>
             <div className="mt-1">
                <p className="text-base font-medium text-gray-700">{options.message}</p>
             </div>
          </div>
          
          <div className="flex justify-end space-x-3 pt-4">
            <Button variant="outline" onClick={handleClose}>
              {options.cancelText || t('common.cancel')}
            </Button>
            <Button 
              variant="primary" 
              onClick={handleConfirm}
            >
              {options.confirmText || t('common.confirm')}
            </Button>
          </div>
        </div>
      </Modal>
    </DialogContext.Provider>
  );
};
