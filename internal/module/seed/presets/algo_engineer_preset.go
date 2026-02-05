package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var algoEngineerPresetJSON = map[string]any{
	"title":    "算法工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "王强",
		"Email":      "wangqiang@example.com",
		"Phone":      "13700000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "算法工程师",
		"City":       "北京",
		"Money":      "35k-55k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#ef6c00",
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
					"description": "<ul><li>4 年算法工程经验，扎实数据结构与算法功底</li><li>熟悉搜索排序与推荐、路径规划与图算法</li><li>掌握 C++/Java 与工程化，性能优化与并发治理</li><li>具备建模与评估能力，指标体系与 A/B 实验</li></ul>",
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
					"title":       "某搜索公司",
					"subtitle":    "算法工程师（排序/推荐）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>优化排序与推荐策略，提升点击与转化</li><li>建设特征与评估体系，落地 A/B 实验</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某物流公司",
					"subtitle":    "算法工程师（路径/图）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>研究路径规划与图算法，优化路线与成本</li><li>完善工程与服务治理，提升性能与稳定</li></ul>",
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
					"title":       "搜索排序与推荐优化",
					"subtitle":    "算法工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：C++ / Java / Feature / AB</li><li>优化排序与推荐策略，提升 CTR 与 CVR</li><li>建设特征与评估体系</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "路径规划与图算法",
					"subtitle":    "算法工程师",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：Graph / Shortest Path / Heuristic</li><li>优化路线与成本，提升效率与稳定</li></ul>",
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
					"title":       "数据结构与算法",
					"subtitle":    "基础",
					"description": "Array / LinkedList / Tree / Graph / DP",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与优化",
					"subtitle":    "性能",
					"description": "C++ / Java / Concurrency / Profiling",
				},
				map[string]any{
					"id":          "s3",
					"title":       "建模与评估",
					"subtitle":    "指标",
					"description": "Feature / Metrics / AB Test",
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
					"title":       "北京航空航天大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为算法与工程；参与竞赛与项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "北京航空航天大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修数据结构与算法；参与校内项目。",
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
					"description": `<ul><li>算法题库笔记：<a href="https://algo-notes.wq.invalid">algo-notes.wq.invalid</a></li><li>排序推荐示例：<a href="https://ranking-demo.wq.invalid">ranking-demo.wq.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateAlgoEngineerPreset() []byte {
	aeByte, _ := json.Marshal(algoEngineerPresetJSON)
	return aeByte
}
