
export type Industry = 'All' | 'General' | 'IT' | 'Finance' | 'Creative' | 'Medical' | 'Education';

export interface Template {
  id: string;
  name: string;
  industry: Industry;
  popularity: number;
  isPremium: boolean;
  tags: string[];
  lastUpdated: string;
}

export interface FilterState {
  keyword: string;
  industry: Industry;
}
