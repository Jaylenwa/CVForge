
export interface Option {
  value: string;
  label: string;
}

export interface FormData {
  projectName: string;
  description: string;
  technologies: string[];
  priority: 'low' | 'medium' | 'high';
}
