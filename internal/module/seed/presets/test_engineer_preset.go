package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var testEngineerPresetJSON = map[string]any{
	"title":    "测试工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "测试工程师",
		"City":       "上海",
		"Money":      "20k-30k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#ef4444",
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
					"description": "<ul><li>5 年软件测试经验，熟悉测试方法与流程，能够独立完成需求评审与测试方案设计</li><li>掌握用例设计（边界值/等价类/判定表/因果图）、缺陷管理与质量度量</li><li>熟悉功能/回归/兼容性/安全测试，具备跨端与多浏览器/机型适配经验</li><li>具备基础自动化能力与脚本编写，能协同研发保障交付质量</li></ul>",
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
					"subtitle":    "测试工程师（商品/交易）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责商品与交易模块测试，完善用例库与冒烟/回归策略</li><li>搭建缺陷管理与质量看板，提升问题跟踪与决策效率</li><li>与研发协作优化上线流程，降低线上故障率</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区",
					"subtitle":    "测试工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>参与权限/消息等基础平台测试，推进冒烟与回归自动化覆盖</li><li>建立兼容性测试矩阵与适配方案，降低线上兼容问题</li></ul>",
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
					"title":       "质量度量与缺陷治理",
					"subtitle":    "测试工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Jira / Confluence / 自研质量看板</li><li>沉淀缺陷与用例度量指标，提升质量可视化与决策效率</li><li>推动缺陷分类与治理方案，缩短问题修复周期</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "兼容性与适配测试矩阵",
					"subtitle":    "测试工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：BrowserStack / 设备农场</li><li>构建多浏览器与机型适配矩阵，完善跨端测试策略</li><li>显著降低线上兼容问题与用户反馈</li></ul>",
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
					"title":       "测试方法与设计",
					"subtitle":    "理论",
					"description": "用例设计 / 边界值 / 等价类 / 判定表 / 因果图",
				},
				map[string]any{
					"id":          "s2",
					"title":       "测试类型与流程",
					"subtitle":    "实践",
					"description": "功能 / 回归 / 兼容性 / 安全 / 探索性测试",
				},
				map[string]any{
					"id":          "s3",
					"title":       "质量与协作",
					"subtitle":    "治理",
					"description": "缺陷管理 / 质量度量 / 看板 / 跨团队协作",
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
					"title":       "南京大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为软件质量与测试工程；参与质量平台建设。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "南京大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修软件工程与测试技术；参与校内项目与竞赛。",
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
					"description": `<ul><li>测试用例规范站：<a href="https://test-cases.zhwei.invalid">test-cases.zhwei.invalid</a></li><li>质量看板示例：<a href="https://qa-dashboard.zhwei.invalid">qa-dashboard.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateTestEngineerPreset() []byte {
	testByte, _ := json.Marshal(testEngineerPresetJSON)
	return testByte
}
