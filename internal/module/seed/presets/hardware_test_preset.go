package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var hardwareTestPresetJSON = map[string]any{
	"title":    "硬件测试工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "硬件测试工程师",
		"City":       "上海",
		"Money":      "25k-35k",
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
					"description": "<ul><li>5 年硬件测试经验，熟悉硬件可靠性与环境测试，具备失效分析与整改能力</li><li>掌握电气/信号/接口测试、EMC/ESD/安规测试与认证流程</li><li>熟悉测试治具与自动化搭建、数据采集与报告输出</li><li>协同研发优化硬件与固件设计，提升稳定性与良率</li></ul>",
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
					"title":       "某智能硬件公司",
					"subtitle":    "硬件测试工程师（可靠性/认证）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责产品可靠性与环境测试（高低温/湿热/振动）</li><li>推进 EMC/ESD/安规认证，优化设计与整改方案</li><li>搭建自动化测试治具与数据采集，提升效率与一致性</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某消费电子企业",
					"subtitle":    "硬件测试工程师（接口/固件）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>开展接口与信号完整性测试，定位问题并提出整改</li><li>协同固件与硬件团队优化协议与时序，提升稳定性</li></ul>",
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
					"title":       "自动化测试治具与数据平台",
					"subtitle":    "硬件测试工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：LabVIEW / Python / NI 测试仪器</li><li>搭建自动化治具与数据平台，统一测试流程与报告输出</li><li>提升测试效率与一致性，支撑多产品线交付</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "EMC/ESD 与安规测试认证",
					"subtitle":    "硬件测试工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：EMC/ESD 仪器 / 标准流程</li><li>完成认证测试与整改，优化硬件设计与屏蔽/接地方案</li><li>提升产品稳定性与良率</li></ul>",
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
					"title":       "测试与认证",
					"subtitle":    "硬件",
					"description": "EMC / ESD / 安规 / 环境与可靠性 / 信号完整性",
				},
				map[string]any{
					"id":          "s2",
					"title":       "自动化与治具",
					"subtitle":    "工程",
					"description": "LabVIEW / Python / NI 仪器 / 数据采集与报告",
				},
				map[string]any{
					"id":          "s3",
					"title":       "协作与优化",
					"subtitle":    "质量",
					"description": "固件/硬件协作 / 设计整改 / 良率提升",
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
					"title":       "西安电子科技大学",
					"major":       "电子信息工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为硬件可靠性与认证；参与测试平台与治具搭建。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "西安电子科技大学",
					"major":       "电子科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修电路与信号；参与校内项目与竞赛。",
				},
			},
		},
	},
}

func GenerateHardwareTestPreset() []byte {
	hwByte, _ := json.Marshal(hardwareTestPresetJSON)
	return hwByte
}
