package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var opsManagerPresetJSON = map[string]any{
	"title":    "运维经理/主管简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "运维经理/主管",
		"City":       "上海",
		"Money":      "40k-60k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "32",
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
					"description": "<ul><li>8 年运维与平台经验，带队建设可观测、交付与 SRE 体系</li><li>负责稳定性与成本治理：SLA/SLO/SLI、容量与成本优化</li><li>完善流程与协作：变更/发布/应急响应与演练，跨团队协作</li><li>人才培养与组织能力：梯队建设、制度与评审、指标驱动改进</li></ul>",
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
					"title":       "某零售平台",
					"subtitle":    "运维经理（平台/SRE）",
					"timeStart":   "2020-03",
					"today":       true,
					"description": "<ul><li>带队建设 GitOps 与可观测平台</li><li>完善稳定性指标与演练机制</li><li>推动跨团队协作与流程改进</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网公司",
					"subtitle":    "运维主管（交付/治理）",
					"timeStart":   "2016-07",
					"timeEnd":     "2020-02",
					"description": "<ul><li>建设交付平台与流程规范</li><li>组织应急响应与复盘改进</li></ul>",
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
					"title":       "平台与稳定性体系建设",
					"subtitle":    "运维经理/主管",
					"timeStart":   "2021-06",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：Argo CD / Helm / Prometheus / Grafana</li><li>统一交付与可观测治理，完善稳定性指标</li><li>组织演练与持续改进</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "流程与组织能力建设",
					"subtitle":    "运维经理/主管",
					"timeStart":   "2018-08",
					"timeEnd":     "2019-06",
					"description": "<ul><li>建立变更/发布/应急响应流程</li><li>梯队与制度建设，指标驱动改进</li></ul>",
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
					"title":       "平台与交付",
					"subtitle":    "工程",
					"description": "GitOps / Helm / CI/CD",
				},
				map[string]any{
					"id":          "s2",
					"title":       "稳定性与可靠性",
					"subtitle":    "SRE",
					"description": "Observability / SLO / SLI / Incident",
				},
				map[string]any{
					"id":          "s3",
					"title":       "流程与组织",
					"subtitle":    "治理",
					"description": "Change / Release / Drill / Training",
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
					"timeStart":   "2013-09",
					"timeEnd":     "2016-06",
					"description": "研究方向为平台工程与组织治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "复旦大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2009-09",
					"timeEnd":     "2013-06",
					"description": "主修分布式系统与工程化；参与校内项目与竞赛。",
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
					"description": `<ul><li>平台治理手册：<a href="https://platform-governance.zhwei.invalid">platform-governance.zhwei.invalid</a></li><li>稳定性指标方案：<a href="https://sre-metrics.zhwei.invalid">sre-metrics.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateOpsManagerPreset() []byte {
	omByte, _ := json.Marshal(opsManagerPresetJSON)
	return omByte
}
