import { useLanguage } from '../contexts/LanguageContext';

export const useSectionTitle = () => {
    const { t } = useLanguage();
    
    return (section: any) => {
        const defaultTitles = ['Professional Summary', 'Work Experience', 'Education', 'Skills', 'Projects'];
        if (defaultTitles.includes(section.title) || !section.title) {
            return t(`section.${section.type}`);
        }
        return section.title;
    };
};
