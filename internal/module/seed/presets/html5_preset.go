package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var html5PresetJSON = map[string]any{
	"title":    "HTML5 前端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "HTML5 前端开发工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 年 H5/移动 Web 开发经验，熟悉混合开发与 WebView 相关能力，掌握 JSBridge 通信与端能力接入</li><li>精通 HTML5/CSS3/JavaScript/TypeScript，熟悉移动端适配、动画与交互、Canvas/WebSocket</li><li>具备工程化与性能优化实践：资源压缩、图片优化、PWA、预加载与缓存策略</li><li>熟悉移动端问题治理：兼容性与适配、触摸与滚动性能、长列表与骨架屏、弱网与离线处理</li></ul>",
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
					"title":       "某互联网金融",
					"subtitle":    "HTML5 前端工程师（活动/营销）",
					"timeStart":   "2021-03",
					"today":       true,
					"description": "<ul><li>负责 App 内嵌 H5 活动页与营销组件，沉淀动画/抽奖/分享等模块</li><li>优化移动端性能与体验：图片与脚本懒加载、预加载路由与数据</li><li>建设 JSBridge 与埋点 SDK（上报/AB 实验），提升数据采集与分析能力</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某出行平台",
					"subtitle":    "HTML5 前端工程师（小程序/H5）",
					"timeStart":   "2019-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>负责车主与用户 H5 页面建设（Vue + TypeScript），统一表单与列表渲染方案</li><li>引入 PWA 与资源缓存策略，弱网场景体验提升显著</li><li>完善埋点与异常上报，缩短问题定位时间</li></ul>",
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
					"title":       "活动页引擎与可视化配置",
					"subtitle":    "HTML5 前端工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：Vue 3 / TypeScript / Vite / Canvas</li><li>实现活动页组件化与可配置化，支持动画/倒计时/抽奖等能力</li><li>沉淀可视化搭建与数据驱动渲染，提升活动交付效率</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "PWA 与离线能力建设",
					"subtitle":    "HTML5 前端工程师",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Workbox / Service Worker / IndexedDB</li><li>实现资源缓存与离线数据存储，增强弱网与断网场景可用性</li><li>优化首屏加载与互动响应，显著提升移动端体验</li></ul>",
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
					"title":       "语言与技术",
					"subtitle":    "前端",
					"description": "HTML5 / CSS3 / JavaScript / TypeScript / Canvas / WebSocket",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与构建",
					"subtitle":    "工程化",
					"description": "Vite / Webpack / ESLint / Prettier / Babel",
				},
				map[string]any{
					"id":          "s3",
					"title":       "移动端与适配",
					"subtitle":    "移动端",
					"description": "响应式布局 / 动画优化 / 触摸与滚动性能 / PWA",
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
					"title":       "复旦大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为移动 Web 与性能优化；参与移动端平台项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "复旦大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修 Web 技术、用户体验与交互设计；参与校园项目。",
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
					"description": `<ul><li>活动组件库：<a href="https://h5kit.zhwei.invalid">h5kit.zhwei.invalid</a></li><li>PWA 示例站：<a href="https://pwa.zhwei.invalid">pwa.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateHTML5Preset() []byte {
	htmlByte, _ := json.Marshal(html5PresetJSON)
	return htmlByte
}
