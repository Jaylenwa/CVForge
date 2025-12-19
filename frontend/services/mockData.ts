import { ResumeData, ResumeSectionType, Template } from '../types';

export const MOCK_TEMPLATES: Template[] = [
  // Base
  { id: 't1', name: 'Classic Professional', tags: ['Professional', 'Simple', 'ATS Friendly'], popularity: 98, isPremium: false, category: 'General' },
  { id: 't2', name: 'Modern Dark', tags: ['Creative', 'Design', 'Startup'], popularity: 85, isPremium: true, category: 'Creative' },
  { id: 't3', name: 'Tech Minimalist', tags: ['Minimalist', 'Tech', 'Clean'], popularity: 92, isPremium: false, category: 'IT' },
  { id: 't4', name: 'Executive Serif', tags: ['Professional', 'Management', 'Senior'], popularity: 70, isPremium: true, category: 'Finance' },
  { id: 't5', name: 'Creative Bold', tags: ['Creative', 'Marketing', 'Colorful'], popularity: 65, isPremium: true, category: 'Creative' },
  { id: 't6', name: 'Elegant Teal', tags: ['Modern', 'Fresh', 'Entry Level'], popularity: 88, isPremium: false, category: 'General' },

  // CN
  { id: 't7', name: 'Chinese Blue', tags: ['General', 'Chinese', 'ATS Friendly'], popularity: 75, isPremium: false, category: 'General' },
  { id: 't8', name: 'Chinese Ziyuan', tags: ['General', 'Chinese', 'Classic'], popularity: 80, isPremium: false, category: 'General' },
  { id: 't9', name: 'Chinese Modern', tags: ['Modern', 'Chinese', 'Clean'], popularity: 82, isPremium: false, category: 'General' },
  { id: 't10', name: 'Chinese Business', tags: ['Business', 'Chinese', 'Professional'], popularity: 85, isPremium: false, category: 'Finance' },
  { id: 't11', name: 'Chinese Creative', tags: ['Creative', 'Chinese', 'Design'], popularity: 78, isPremium: true, category: 'Creative' },
  { id: 't12', name: 'Chinese Tech', tags: ['Tech', 'Chinese', 'Geek'], popularity: 88, isPremium: false, category: 'IT' },
  { id: 't13', name: 'Chinese Academic', tags: ['Academic', 'Chinese', 'Research'], popularity: 60, isPremium: false, category: 'General' },
  { id: 't14', name: 'Chinese Minimal', tags: ['Minimal', 'Chinese', 'Simple'], popularity: 90, isPremium: false, category: 'General' },
  { id: 't15', name: 'Chinese Marketing', tags: ['Marketing', 'Chinese', 'Bold'], popularity: 72, isPremium: true, category: 'Creative' },
  { id: 't16', name: 'Chinese Artistic', tags: ['Artistic', 'Chinese', 'Ink'], popularity: 68, isPremium: true, category: 'Creative' },
  { id: 't17', name: 'Chinese Retro', tags: ['Retro', 'Chinese', 'Newspaper'], popularity: 55, isPremium: true, category: 'Creative' },
  { id: 't18', name: 'Chinese Fresh', tags: ['Fresh', 'Chinese', 'Green'], popularity: 76, isPremium: false, category: 'General' },
  { id: 't19', name: 'Chinese Dark', tags: ['Dark', 'Chinese', 'Luxury'], popularity: 65, isPremium: true, category: 'Creative' },
  { id: 't20', name: 'Chinese Magazine', tags: ['Magazine', 'Chinese', 'Fashion'], popularity: 62, isPremium: true, category: 'Creative' },
  { id: 't21', name: 'Chinese Timeline', tags: ['Timeline', 'Chinese', 'Infographic'], popularity: 70, isPremium: false, category: 'General' },

  // New CN (t22-t41)
  { id: 't22', name: 'Chinese Architect', tags: ['Architecture', 'Chinese', 'Structure'], popularity: 60, isPremium: true, category: 'Creative' },
  { id: 't23', name: 'Chinese Block', tags: ['Block', 'Chinese', 'Modern'], popularity: 75, isPremium: false, category: 'General' },
  { id: 't24', name: 'Chinese Brush', tags: ['Brush', 'Chinese', 'Art'], popularity: 58, isPremium: true, category: 'Creative' },
  { id: 't25', name: 'Chinese Cloud', tags: ['Cloud', 'Chinese', 'Soft'], popularity: 66, isPremium: false, category: 'IT' },
  { id: 't26', name: 'Chinese Code', tags: ['Code', 'Chinese', 'Dev'], popularity: 85, isPremium: false, category: 'IT' },
  { id: 't27', name: 'Chinese Fashion', tags: ['Fashion', 'Chinese', 'Trend'], popularity: 70, isPremium: true, category: 'Creative' },
  { id: 't28', name: 'Chinese Finance', tags: ['Finance', 'Chinese', 'Data'], popularity: 82, isPremium: false, category: 'Finance' },
  { id: 't29', name: 'Chinese Game', tags: ['Game', 'Chinese', 'Playful'], popularity: 64, isPremium: false, category: 'IT' },
  { id: 't30', name: 'Chinese Geometric', tags: ['Geometric', 'Chinese', 'Shape'], popularity: 72, isPremium: false, category: 'General' },
  { id: 't31', name: 'Chinese Glitch', tags: ['Glitch', 'Chinese', 'Cyberpunk'], popularity: 50, isPremium: true, category: 'IT' },
  { id: 't32', name: 'Chinese Hexagon', tags: ['Hexagon', 'Chinese', 'Tech'], popularity: 78, isPremium: false, category: 'IT' },
  { id: 't33', name: 'Chinese Law', tags: ['Law', 'Chinese', 'Serious'], popularity: 68, isPremium: false, category: 'General' },
  { id: 't34', name: 'Chinese Medical', tags: ['Medical', 'Chinese', 'Clean'], popularity: 74, isPremium: false, category: 'General' },
  { id: 't35', name: 'Chinese Newspaper', tags: ['Newspaper', 'Chinese', 'Classic'], popularity: 52, isPremium: true, category: 'Creative' },
  { id: 't36', name: 'Chinese Origami', tags: ['Origami', 'Chinese', 'Paper'], popularity: 60, isPremium: true, category: 'Creative' },
  { id: 't37', name: 'Chinese Pixel', tags: ['Pixel', 'Chinese', 'Retro'], popularity: 55, isPremium: true, category: 'IT' },
  { id: 't38', name: 'Chinese Product', tags: ['Product', 'Chinese', 'Showcase'], popularity: 80, isPremium: false, category: 'General' },
  { id: 't39', name: 'Chinese Social', tags: ['Social', 'Chinese', 'Media'], popularity: 76, isPremium: false, category: 'Creative' },
  { id: 't40', name: 'Chinese Vogue', tags: ['Vogue', 'Chinese', 'Style'], popularity: 65, isPremium: true, category: 'Creative' },
  { id: 't41', name: 'Chinese Wave', tags: ['Wave', 'Chinese', 'Fluid'], popularity: 70, isPremium: false, category: 'Creative' },
];

export const INITIAL_RESUME: ResumeData = {
  id: 'new',
  title: 'My Resume',
  templateId: 't1',
  themeConfig: {
    color: '#000000',
    fontFamily: 'inter',
    spacing: 'normal'
  },
  lastModified: Date.now(),
  personalInfo: {
    fullName: 'YOUR NAME',
    email: 'email@example.com',
    phone: '123-456-7890',
    avatarUrl: 'https://i.pravatar.cc/256?img=12',
    jobTitle: 'Software Engineer',
    gender: 'Male',
    age: '25',
    maritalStatus: 'Single',
    politicalStatus: 'Member',
    birthplace: 'Beijing',
    ethnicity: 'Han',
    height: '180cm',
    weight: '70kg',
    customInfo: []
  },
  sections: [
    {
      id: 'summary',
      type: ResumeSectionType.Summary,
      title: 'Professional Summary',
      isVisible: true,
      items: [
        {
          id: 's1',
          description: 'Experienced software engineer with a passion for building scalable web applications.'
        }
      ]
    },
    {
      id: 'exp',
      type: ResumeSectionType.Experience,
      title: 'Work Experience',
      isVisible: true,
      items: [
        {
          id: 'e1',
          title: 'Senior Developer',
          subtitle: 'Tech Corp',
          dateRange: '2020 - Present',
          description: 'Led a team of 5 developers to build a new SaaS product.'
        }
      ]
    },
    {
      id: 'edu',
      type: ResumeSectionType.Education,
      title: 'Education',
      isVisible: true,
      items: [
        {
          id: 'ed1',
          title: 'Computer Science',
          subtitle: 'University of Tech',
          dateRange: '2016 - 2020',
          description: 'Bachelor of Science'
        }
      ]
    },
    {
      id: 'skills',
      type: ResumeSectionType.Skills,
      title: 'Skills',
      isVisible: true,
      items: [
        { id: 'sk1', description: 'React' },
        { id: 'sk2', description: 'TypeScript' },
        { id: 'sk3', description: 'Go' },
        { id: 'sk4', description: 'Node.js' }
      ]
    },
    {
      id: 'projects',
      type: ResumeSectionType.Projects,
      title: 'Projects',
      isVisible: true,
      items: [
        {
          id: 'p1',
          title: 'OpenResume',
          subtitle: 'Open Source',
          dateRange: '2023',
          description: 'A free and open source resume builder.'
        }
      ]
    }
  ]
};
