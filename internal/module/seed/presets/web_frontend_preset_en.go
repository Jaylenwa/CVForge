package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var webFrontendPresetJSONEn = map[string]any{
	"title":    "Web Frontend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Web Frontend Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#06b6d4",
		"Font":     "notosans",
		"Spacing":  "normal",
		"FontSize": "13",
	},
	"sections": []any{
		map[string]any{
			"id":        "summary",
			"type":      common.ResumeSectionTypeSummary,
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 years of web frontend experience; proficient in TypeScript and React/Vue ecosystems, with solid componentization and architecture capabilities</li><li>Performance-focused: Core Web Vitals, code splitting, lazy loading, SSR/SSG, pre-rendering, and caching strategies</li><li>Engineering practices: Vite/Webpack, ESLint/Prettier, Jest/Vitest, Playwright/Cypress, CI/CD and progressive delivery</li><li>Reliability and observability: error reporting, user behavior tracking, performance monitoring, and alerting</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "Work Experience",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "E-commerce Platform",
					"subtitle":    "Web Frontend Engineer (Product/Checkout)",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>Led product and checkout modules (React + TypeScript), building reusable form/table components</li><li>Decoupled routing and data layers (React Query + Zustand) to improve maintainability and testability</li><li>Optimized performance: code splitting, image optimization, lazy loading; reduced homepage LCP by 35%</li><li>Established observability (Sentry + in-house SDK) for errors and tracking, cutting incident time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community",
					"subtitle":    "Web Frontend Engineer (Platform/Admin)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Refactored admin console (Vue 3 + Vite + TypeScript), abstracting permissions/menu/form engines</li><li>Built a design system and component library (Tailwind CSS + Headless UI) to unify UX</li><li>Improved testing coverage (Vitest/Playwright) to increase release stability</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "Projects",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "SSR E-commerce Site with Dynamic Routing",
					"subtitle":    "Web Frontend Engineer",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: Next.js / React / TypeScript / Tailwind CSS</li><li>Implemented SSR with incremental static generation for better SEO and first paint</li><li>Built multi-tenant theming and component strategy for extensibility</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Micro-Frontend Integration",
					"subtitle":    "Web Frontend Engineer",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Stack: Qiankun / React / Vue / TypeScript</li><li>Integrated multiple frontend apps with a shared component library to improve reuse</li><li>Designed communication and routing strategies to reduce coupling</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "Skills",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "Languages & Frameworks",
					"subtitle":    "Frontend",
					"description": "TypeScript / React / Vue / Next.js / Nuxt",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Tooling & Build",
					"subtitle":    "Engineering",
					"description": "Vite / Webpack / Rollup / ESLint / Prettier",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Quality & Testing",
					"subtitle":    "Testing",
					"description": "Jest / Vitest / React Testing Library / Playwright",
				},
				map[string]any{
					"id":          "s4",
					"title":       "Performance & UX",
					"subtitle":    "Optimization",
					"description": "Core Web Vitals / Lazy Loading / Code Splitting / Caching",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "Education",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "Shanghai Jiao Tong University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Research on frontend engineering and performance; contributed to frontend platform building.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Shanghai Jiao Tong University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied computer networks, web technologies, and HCI; participated in projects and competitions.",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "Portfolio",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>UI library demo: <a href="https://ui.zhwei.invalid">ui.zhwei.invalid</a></li><li>Performance showcase: <a href="https://perf.zhwei.invalid">perf.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateWebFrontendPresetEn() []byte {
	webByte, _ := json.Marshal(webFrontendPresetJSONEn)
	return webByte
}
