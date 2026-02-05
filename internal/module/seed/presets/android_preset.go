package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var androidPresetJSON = map[string]any{
	"title":    "Android 开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Android 开发工程师",
		"City":       "上海",
		"Money":      "25k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#3ddc84",
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
					"description": "<ul><li>5 年 Android 开发经验，熟悉 Kotlin/Java、Jetpack 架构组件、Compose/MVVM</li><li>掌握并发与数据流（Coroutines/Flow）、网络与存储（Retrofit/Room/OKHttp）</li><li>具备性能与稳定性建设：冷启动优化、内存与耗电分析、崩溃治理</li><li>熟悉依赖注入与工程化：Hilt/Dagger、Gradle 构建、模块化与多渠道发布</li></ul>",
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
					"title":       "某在线教育 App",
					"subtitle":    "Android 开发工程师（课程/支付）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责课程与支付模块，采用 MVVM + Jetpack 架构，提升可维护性</li><li>优化冷启动与页面渲染，首页启动时长降低 30%</li><li>建设崩溃分析与埋点（Crashlytics + 自研 SDK），提升稳定性与定位效率</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区 App",
					"subtitle":    "Android 开发工程师（账号/消息）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>重构账号与消息模块，引入 Hilt 与协程，降低复杂度</li><li>优化网络与缓存策略，关键链路 P95 延迟下降 40%</li><li>构建多模块与多渠道发布流程，提升交付效率</li></ul>",
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
					"title":       "Compose 组件化改造",
					"subtitle":    "Android 开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Kotlin / Jetpack Compose / Hilt / Room</li><li>沉淀通用 UI 组件与状态管理方案，提升复用性与一致性</li><li>引入快照与差分渲染，改善交互体验与性能</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "离线缓存与弱网优化",
					"subtitle":    "Android 开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：OKHttp / Retrofit / Room / WorkManager</li><li>实现数据与资源缓存、断点续传与任务重试，增强弱网可用性</li><li>通过指标与埋点评估体验改善，迭代优化效果显著</li></ul>",
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
					"subtitle":    "Android",
					"description": "Kotlin / Java / Jetpack / Compose / MVVM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "网络与存储",
					"subtitle":    "数据",
					"description": "Retrofit / OKHttp / Room / SQLite",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程化与发布",
					"subtitle":    "交付",
					"description": "Gradle / 多模块 / 多渠道 / CI/CD / Crashlytics",
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
					"title":       "北京大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为移动端架构与性能优化；参与平台建设与工具开发。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "北京大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修移动开发、计算机网络与操作系统；参与校内项目与竞赛。",
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
					"description": `<ul><li>组件库与工具：<a href="https://androidkit.zhwei.invalid">androidkit.zhwei.invalid</a></li><li>技术博客：<a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateAndroidPreset() []byte {
	androidByte, _ := json.Marshal(androidPresetJSON)
	return androidByte
}
