package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var automationTestPresetJSON = map[string]any{
	"title":    "自动化测试工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "自动化测试工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#22c55e",
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
					"description": "<ul><li>5 年自动化测试经验，熟悉 Web/App 端自动化框架与最佳实践</li><li>精通 UI 端到端与接口自动化：Playwright/Cypress/Selenium、pytest、REST/GraphQL</li><li>具备稳定性与可维护性建设：用例分层/数据驱动/持续集成与质量门禁</li><li>熟悉环境准备与 Mock/Stub、容器化与并行执行、测试报告与可视化</li></ul>",
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
					"subtitle":    "自动化测试工程师（平台/交易）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>构建 Web 端 E2E 与接口自动化体系（Playwright + pytest），覆盖核心链路</li><li>实现数据工厂与环境隔离，稳定性与可重复性提升</li><li>接入 CI/CD 与质量门禁，减少回归人力与风险</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区",
					"subtitle":    "自动化测试工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>搭建组件与页面层自动化框架，完善定位与稳定性策略</li><li>构建并行执行与报告可视化，提升反馈速度与可视性</li></ul>",
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
					"title":       "统一自动化框架与用例治理",
					"subtitle":    "自动化测试工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Playwright / pytest / Docker</li><li>沉淀分层架构与数据驱动、Mock/Stub 与环境隔离</li><li>接入 CI 并行执行与质量门禁，稳定回归与快速反馈</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "接口自动化与契约测试",
					"subtitle":    "自动化测试工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：pytest / Pact / REST / GraphQL</li><li>实现接口自动化与契约测试，降低回归风险与耦合</li><li>完善数据工厂与告警通知，提升定位效率</li></ul>",
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
					"title":       "自动化框架",
					"subtitle":    "Web/App",
					"description": "Playwright / Cypress / Selenium / Appium",
				},
				map[string]any{
					"id":          "s2",
					"title":       "接口与契约",
					"subtitle":    "API",
					"description": "pytest / REST / GraphQL / Pact",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程化与执行",
					"subtitle":    "CI",
					"description": "Docker / 并行执行 / 报告可视化 / 质量门禁",
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
					"title":       "中国科学技术大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为测试自动化与工程化；参与质量平台项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "中国科学技术大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修软件工程与测试技术；参与校内项目与竞赛。",
				},
			},
		},
	},
}

func GenerateAutomationTestPreset() []byte {
	autoByte, _ := json.Marshal(automationTestPresetJSON)
	return autoByte
}
