package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var html5PresetJSONEn = map[string]any{
	"title":    "HTML5 Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "HTML5 Frontend Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f59e0b",
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
					"description": "<ul><li>5 years of H5/mobile web experience; familiar with hybrid development and WebView, JSBridge communication, and native capability integration</li><li>Proficient in HTML5/CSS3/JavaScript/TypeScript; experienced with mobile adaptation, animations, Canvas, and WebSocket</li><li>Engineering and performance optimization: asset compression, image optimization, PWA, preloading, and caching strategies</li><li>Mobile reliability: compatibility/adaptation, touch/scroll performance, long list optimization/skeleton screens, weak network/offline handling</li></ul>",
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
					"title":       "Internet Finance",
					"subtitle":    "HTML5 Frontend Engineer (Campaign/Marketing)",
					"timeStart":   "2021-03",
					"today":       true,
					"description": "<ul><li>Built in-app H5 campaign pages and marketing components; delivered animations/lotteries/sharing modules</li><li>Optimized mobile performance: lazy loading for images/scripts, route/data preloading</li><li>Developed JSBridge and tracking SDK (reporting/AB testing) to enhance data collection/analysis</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Mobility Platform",
					"subtitle":    "HTML5 Frontend Engineer (Miniapp/H5)",
					"timeStart":   "2019-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Implemented driver/user H5 pages (Vue + TypeScript), standardized form and list rendering</li><li>Introduced PWA and caching strategies for better UX in weak network conditions</li><li>Improved observability with tracking and error reporting for faster incident resolution</li></ul>",
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
					"title":       "Campaign Page Engine & Visual Config",
					"subtitle":    "HTML5 Frontend Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: Vue 3 / TypeScript / Vite / Canvas</li><li>Componentized and configurable campaign pages supporting animations/countdown/lottery</li><li>Visual builder and data-driven rendering to accelerate delivery</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "PWA & Offline Capabilities",
					"subtitle":    "HTML5 Frontend Engineer",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: Workbox / Service Worker / IndexedDB</li><li>Implemented resource caching and offline storage for robust UX in weak networks</li><li>Optimized first paint and interaction responsiveness for mobile devices</li></ul>",
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
					"title":       "Languages & Tech",
					"subtitle":    "Frontend",
					"description": "HTML5 / CSS3 / JavaScript / TypeScript / Canvas / WebSocket",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Tooling & Build",
					"subtitle":    "Engineering",
					"description": "Vite / Webpack / ESLint / Prettier / Babel",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Mobile & Adaptation",
					"subtitle":    "Mobile",
					"description": "Responsive Layout / Animation Optimization / Touch & Scroll / PWA",
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
					"title":       "Fudan University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Research on mobile web and performance optimization; contributed to mobile platform projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Fudan University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Focused on web technologies, UX, and interaction design; participated in campus projects.",
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
					"description": `<ul><li>H5 Component Kit: <a href="https://h5kit.zhwei.invalid">h5kit.zhwei.invalid</a></li><li>PWA Demo: <a href="https://pwa.zhwei.invalid">pwa.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateHTML5PresetEn() []byte {
	htmlByte, _ := json.Marshal(html5PresetJSONEn)
	return htmlByte
}
