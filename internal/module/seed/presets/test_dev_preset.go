package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var testDevPresetJSON = map[string]any{
	"title":    "测试开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "测试开发工程师",
		"City":       "上海",
		"Money":      "30k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#8b5cf6",
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
					"description": "<ul><li>5 年测试开发经验，具备测试框架与平台建设能力</li><li>熟悉后端与前端自动化、接口契约、数据工厂与环境管理</li><li>具备质量平台、用例治理、覆盖率与质量度量体系建设能力</li><li>熟悉 CI/CD 与质量门禁、稳定性建设与可观测</li></ul>",
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
					"title":       "某平台中台",
					"subtitle":    "测试开发工程师（平台/质量）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>建设统一测试平台与框架，沉淀分层架构与用例治理</li><li>构建数据工厂、环境隔离与并行执行，提升效率与稳定性</li><li>接入 CI/CD 与质量门禁，完善度量与看板</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某电商中台",
					"subtitle":    "测试开发工程师（接口/契约）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>建设接口自动化与契约测试平台，降低耦合与回归风险</li><li>完善数据生成与模拟方案，提升定位与修复效率</li></ul>",
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
					"title":       "统一测试平台与质量看板",
					"subtitle":    "测试开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Go / Python / Playwright / Grafana</li><li>统一平台与用例治理、覆盖率与质量度量看板</li><li>接入告警与门禁，保障交付质量</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "契约测试平台与数据工厂",
					"subtitle":    "测试开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Pact / pytest / Postgres / Redis</li><li>实现契约测试、数据生成与隔离，降低耦合与回归风险</li><li>完善报告与定位工具，提高问题响应速度</li></ul>",
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
					"subtitle":    "测试开发",
					"description": "Go / Python / Playwright / pytest / Pact",
				},
				map[string]any{
					"id":          "s2",
					"title":       "平台与工程",
					"subtitle":    "质量",
					"description": "测试平台 / 数据工厂 / 并行执行 / 度量与看板",
				},
				map[string]any{
					"id":          "s3",
					"title":       "CI/CD 与治理",
					"subtitle":    "交付",
					"description": "CI/CD / 质量门禁 / 覆盖率 / 告警与可观测",
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
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为测试工程与平台建设；参与质量体系与度量。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "上海交通大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修软件工程与系统开发；参与校内项目与竞赛。",
				},
			},
		},
	},
}

func GenerateTestDevPreset() []byte {
	tdByte, _ := json.Marshal(testDevPresetJSON)
	return tdByte
}
