package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var webFrontendPresetJSON = map[string]any{
	"title":    "Web 前端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Web 前端开发工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 年 Web 前端开发经验，熟悉 TypeScript，精通 React/Vue 生态，具备组件化与复杂业务前端架构设计能力</li><li>掌握性能优化与体验提升：Core Web Vitals、代码分割、懒加载、SSR/SSG、预渲染与缓存策略</li><li>熟悉工程化体系：Vite/Webpack、ESLint/Prettier、Jest/Vitest、Playwright/Cypress、CI/CD 与灰度发布</li><li>具备稳定性与可观测建设：前端异常上报、用户行为埋点、性能监控与告警</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "工作经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "某电商平台",
					"subtitle":    "Web 前端工程师（商品/交易）",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>负责商品与交易前端模块（React + TypeScript），沉淀表单/表格等高复用组件</li><li>引入路由与数据层解耦（React Query + Zustand），提升可维护性与可测试性</li><li>优化首屏与交互性能：代码分割、图片优化与懒加载，首页 LCP 降低 35%</li><li>建设埋点与异常上报体系（Sentry + 自研 SDK），定位线上问题效率提升 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区",
					"subtitle":    "Web 前端工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>参与管理后台重构（Vue 3 + Vite + TypeScript），抽象权限/菜单/表单引擎</li><li>推动设计系统与组件库建设（Tailwind CSS + Headless UI），统一视觉与交互规范</li><li>完善单元与端到端测试（Vitest/Playwright），发布稳定性提升</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "项目经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "SSR 电商站点与动态路由",
					"subtitle":    "Web 前端工程师",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：Next.js / React / TypeScript / Tailwind CSS</li><li>实现 SSR + 增量静态生成，提升 SEO 与首屏速度</li><li>沉淀多租户主题与组件方案，支持按需扩展</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "微前端与多应用整合",
					"subtitle":    "Web 前端工程师",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>技术栈：Qiankun / React / Vue / TypeScript</li><li>整合多业务前端应用与共享组件库，提升复用性与交付效率</li><li>构建公共通信与路由策略，降低耦合度</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "技能清单",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "语言与框架",
					"subtitle":    "前端",
					"description": "TypeScript / React / Vue / Next.js / Nuxt",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与构建",
					"subtitle":    "工程化",
					"description": "Vite / Webpack / Rollup / ESLint / Prettier",
				},
				map[string]any{
					"id":          "s3",
					"title":       "质量与测试",
					"subtitle":    "测试",
					"description": "Jest / Vitest / React Testing Library / Playwright",
				},
				map[string]any{
					"id":          "s4",
					"title":       "性能与体验",
					"subtitle":    "优化",
					"description": "Core Web Vitals / 懒加载 / 代码分割 / 缓存策略",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "教育背景",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "上海交通大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为前端工程化与性能优化；参与前端平台搭建。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "上海交通大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修计算机网络、Web 技术与人机交互；参与校内项目与竞赛。",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "个人作品",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>组件库示例：<a href="https://ui.zhwei.invalid">ui.zhwei.invalid</a></li><li>性能优化示例站：<a href="https://perf.zhwei.invalid">perf.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateWebFrontendPreset() []byte {
	webByte, _ := json.Marshal(webFrontendPresetJSON)
	return webByte
}
