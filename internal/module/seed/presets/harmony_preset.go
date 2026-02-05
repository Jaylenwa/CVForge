package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var harmonyPresetJSON = map[string]any{
	"title":    "鸿蒙开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "鸿蒙开发工程师",
		"City":       "上海",
		"Money":      "25k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#0ea5e9",
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
					"description": "<ul><li>3+ 年鸿蒙/HarmonyOS 开发经验，熟悉 ArkTS/ETS，Ability 模型与 UI 框架</li><li>掌握 Stage 模型、分布式能力与设备协同、HAP 打包与签名</li><li>具备性能与稳定性：启动与绘制优化、内存与功耗分析、崩溃治理</li><li>熟悉工程化与对接：多模块工程、能力接入与权限、与后端 API/消息系统对接</li></ul>",
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
					"title":       "某 IoT 平台",
					"subtitle":    "鸿蒙开发工程师（设备/控制）",
					"timeStart":   "2021-09",
					"today":       true,
					"description": "<ul><li>负责设备管理与控制端应用，ArkTS + Ability 架构</li><li>实现分布式设备协同与能力共享，提升跨设备体验</li><li>优化性能与稳定性，减少崩溃与卡顿</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某生活服务平台",
					"subtitle":    "鸿蒙开发工程师（服务/运营）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-08",
					"description": "<ul><li>构建服务预约与运营模块，引入分布式数据与消息能力</li><li>完善权限与能力接入、埋点与异常上报</li></ul>",
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
					"title":       "分布式设备协同应用",
					"subtitle":    "鸿蒙开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：ArkTS / Ability / Distributed Data</li><li>实现跨设备能力协同与状态同步</li><li>优化性能与链路稳定性</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "UI 框架与组件库沉淀",
					"subtitle":    "鸿蒙开发工程师",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>技术栈：ArkTS / ETS / UI</li><li>沉淀通用组件与主题系统、状态管理与路由</li><li>提升一致性与开发效率</li></ul>",
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
					"subtitle":    "鸿蒙",
					"description": "ArkTS / ETS / Ability / UI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与性能",
					"subtitle":    "优化",
					"description": "多模块 / HAP / 启动与绘制 / 内存与功耗 / 崩溃治理",
				},
				map[string]any{
					"id":          "s3",
					"title":       "分布式与对接",
					"subtitle":    "协同",
					"description": "分布式能力 / 数据同步 / 权限与能力接入 / 后端 API",
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
					"title":       "华中科技大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式与移动端；参与鸿蒙生态相关项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "华中科技大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修系统开发与工程化；参与校内项目与竞赛。",
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
					"description": `<ul><li>组件库：<a href="https://arktskit.zhwei.invalid">arktskit.zhwei.invalid</a></li><li>分布式示例：<a href="https://distributed.zhwei.invalid">distributed.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateHarmonyPreset() []byte {
	harmonyByte, _ := json.Marshal(harmonyPresetJSON)
	return harmonyByte
}
