package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var iosPresetJSON = map[string]any{
	"title":    "iOS 开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "iOS 开发工程师",
		"City":       "上海",
		"Money":      "25k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#1e90ff",
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
					"description": "<ul><li>5 年 iOS 开发经验，熟悉 Swift/Objective-C、UIKit/SwiftUI、Combine</li><li>掌握网络与并发（URLSession/Alamofire/GCD/Async-Await）、存储与持久化（Core Data）</li><li>具备性能与稳定性：启动优化、内存泄漏治理、崩溃分析与监控</li><li>熟悉工程与交付：SPM/CocoaPods、单元测试与 UI 测试、TestFlight/企业发布</li></ul>",
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
					"title":       "某电商 App",
					"subtitle":    "iOS 开发工程师（商品/订单）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责商品与订单模块，采用 MVVM + Combine 架构</li><li>优化渲染与交互性能，关键页面掉帧率降低</li><li>完善崩溃与异常上报体系（Crashlytics/Sentry），定位效率提升</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区 App",
					"subtitle":    "iOS 开发工程师（账号/消息）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>重构账号与消息模块，引入 SPM 与动态图形渲染优化</li><li>搭建 CI/CD 与自动化测试流程，提升交付质量与效率</li></ul>",
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
					"title":       "SwiftUI 组件与主题系统",
					"subtitle":    "iOS 开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Swift / SwiftUI / Combine</li><li>抽象 UI 组件与主题系统，提升复用与一致性</li><li>优化渲染路径与状态管理，体验与性能改善显著</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "离线能力与数据同步",
					"subtitle":    "iOS 开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Core Data / Background Tasks / URLSession</li><li>实现离线缓存与后台同步、失败重试与冲突处理</li><li>通过指标评估与迭代优化用户体验</li></ul>",
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
					"subtitle":    "iOS",
					"description": "Swift / Objective-C / UIKit / SwiftUI / Combine",
				},
				map[string]any{
					"id":          "s2",
					"title":       "网络与并发",
					"subtitle":    "数据",
					"description": "URLSession / Alamofire / GCD / Async-Await / Core Data",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程与交付",
					"subtitle":    "发布",
					"description": "SPM / CocoaPods / 单元测试 / UI 测试 / TestFlight",
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
					"title":       "清华大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为 iOS 架构与性能优化；参与平台建设与质量工程。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "清华大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修移动开发、编译与链接；参与校内项目与竞赛。",
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
					"description": `<ul><li>组件库：<a href="https://ioskit.zhwei.invalid">ioskit.zhwei.invalid</a></li><li>技术博客：<a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateIOSPreset() []byte {
	iosByte, _ := json.Marshal(iosPresetJSON)
	return iosByte
}
