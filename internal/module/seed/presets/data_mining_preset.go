package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dataMiningPresetJSON = map[string]any{
	"title":    "数据挖掘工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "李娜",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "数据挖掘工程师",
		"City":       "北京",
		"Money":      "30k-45k",
		"JoinTime":   "随时到岗",
		"Gender":     "女",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
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
					"description": "<ul><li>4 年数据挖掘经验，熟悉特征工程与建模评估</li><li>掌握 SQL/Python/Pandas，ETL 与数据质量治理</li><li>了解 Spark/Flink，批流一体与指标体系建设</li><li>具备业务理解与指标落地，推动数据驱动决策</li></ul>",
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
					"title":       "某互联网公司",
					"subtitle":    "数据挖掘工程师（建模/分析）",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>搭建用户画像与特征库，设计训练数据</li><li>完成留存/转化/流失建模与评估</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某电商公司",
					"subtitle":    "数据挖掘工程师（ETL/报表）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>维护数据仓库与 ETL 作业，完善数据质量</li><li>建设指标与报表，支持业务分析</li></ul>",
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
					"title":       "用户流失预测与挽回",
					"subtitle":    "数据挖掘工程师",
					"timeStart":   "2022-01",
					"timeEnd":     "2022-10",
					"description": "<ul><li>技术栈：Python / Pandas / Sklearn</li><li>构建特征与模型，提升召回效率</li><li>落地 A/B 实验与指标评估</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "数据仓库与指标体系建设",
					"subtitle":    "数据挖掘工程师",
					"timeStart":   "2020-05",
					"timeEnd":     "2020-12",
					"description": "<ul><li>技术栈：SQL / Spark / Airflow</li><li>完善 ETL 与数据质量，支持分析与报表</li></ul>",
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
					"title":       "数据与建模",
					"subtitle":    "核心",
					"description": "SQL / Python / Pandas / Sklearn",
				},
				map[string]any{
					"id":          "s2",
					"title":       "计算与工程",
					"subtitle":    "平台",
					"description": "Spark / Flink / Airflow",
				},
				map[string]any{
					"id":          "s3",
					"title":       "指标与实验",
					"subtitle":    "分析",
					"description": "Metrics / AB Test / Evaluation",
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
					"title":       "中国人民大学",
					"major":       "统计学",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为数据分析与建模；参与数据竞赛与项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "中国人民大学",
					"major":       "数学与应用数学",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修概率统计与数值分析；参与校内项目。",
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
					"description": `<ul><li>数据挖掘示例：<a href="https://datamining.zhwei.invalid">datamining.zhwei.invalid</a></li><li>AB 测试示例：<a href="https://abtest.zhwei.invalid">abtest.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDataMiningPreset() []byte {
	dmByte, _ := json.Marshal(dataMiningPresetJSON)
	return dmByte
}
