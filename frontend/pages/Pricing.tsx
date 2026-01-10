import React, { useState } from 'react';
import { useLanguage } from '../contexts/LanguageContext';
import { Modal } from '../components/ui/Modal';
import { Button } from '../components/ui/Button';

export const Pricing: React.FC = () => {
    const { t } = useLanguage();
    const [isModalOpen, setIsModalOpen] = useState(false);
    React.useEffect(() => {
        document.body.classList.add('no-scrollbar');
        document.documentElement.classList.add('no-scrollbar');
        return () => {
            document.body.classList.remove('no-scrollbar');
            document.documentElement.classList.remove('no-scrollbar');
        };
    }, []);

    const handleBuy = () => {
        setIsModalOpen(true);
    };

    return (
        <div className="bg-white min-h-[calc(100vh-4rem-1px)] flex items-center">
            <div className="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 -translate-y-6 md:-translate-y-12">
                <div className="text-center">
                    <h2 className="text-3xl font-extrabold text-gray-900 sm:text-4xl">{t('pricing.title')}</h2>
                    <p className="mt-4 text-lg text-gray-500">{t('pricing.subtitle')}</p>
                </div>
                <div className="mt-12 space-y-4 sm:mt-16 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-6 max-w-6xl mx-auto">
                    {/* Basic Plan */}
                    <div className="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
                        <div className="p-6">
                            <h2 className="text-lg leading-6 font-medium text-gray-900">{t('pricing.basic')}</h2>
                            <p className="mt-4 text-sm text-gray-500">{t('pricing.basicDesc')}</p>
                            <p className="mt-8">
                                <span className="text-4xl font-extrabold text-gray-900">{t('pricing.price.free')}</span>
                                <span className="text-base font-medium text-gray-500">{t('pricing.month')}</span>
                            </p>
                            <Button 
                                variant="primary" 
                                className="mt-8 w-full"
                                onClick={handleBuy}
                            >
                                {t('pricing.buy')} {t('pricing.basic')}
                            </Button>
                        </div>
                    </div>

                    {/* Pro Plan */}
                    <div className="border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200">
                        <div className="p-6">
                            <h2 className="text-lg leading-6 font-medium text-gray-900">{t('pricing.pro')}</h2>
                            <p className="mt-4 text-sm text-gray-500">{t('pricing.proDesc')}</p>
                            <p className="mt-8">
                                <span className="text-4xl font-extrabold text-gray-900">{t('pricing.price.pro')}</span>
                                <span className="text-base font-medium text-gray-500">{t('pricing.month')}</span>
                            </p>
                            <Button 
                                variant="primary" 
                                className="mt-8 w-full"
                                onClick={handleBuy}
                            >
                                {t('pricing.buy')} {t('pricing.pro')}
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={t('pricing.modal.title')}
            >
                <div className="text-center">
                    <p className="text-gray-600 mb-6">{t('pricing.modal.planning')}</p>
                    <Button onClick={() => setIsModalOpen(false)}>
                        {t('pricing.modal.close')}
                    </Button>
                </div>
            </Modal>
        </div>
    );
};
