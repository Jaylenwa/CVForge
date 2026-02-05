package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var bigdataPresetJSON = map[string]any{
	"title":    "大数据工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "李娜",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "大数据工程师",
		"City":       "北京",
		"Money":      "35k-50k",
		"JoinTime":   "随时到岗",
		"Gender":     "女",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#3b82f6",
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
					"description": "<ul><li>4 年大数据工程经验，熟悉批流一体与数据湖/仓</li><li>掌握 Hadoop/Spark/Flink 与 Kafka</li><li>数据开发与治理：ETL/调度/质量/血缘与权限</li><li>指标体系与数据服务，支撑分析与业务场景</li></ul>",
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
					"subtitle":    "大数据工程师（开发/平台）",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>建设数据湖与仓，完善层次与质量治理</li><li>维护批流任务与调度，支撑分析与服务</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某电商公司",
					"subtitle":    "大数据工程师（实时/指标）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>搭建实时计算与指标体系</li><li>完善数据服务与权限治理</li></ul>",
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
					"title":       "数据湖与仓治理",
					"subtitle":    "大数据工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：Hadoop / Spark / Hive / Iceberg</li><li>完善层次与质量治理、血缘与权限</li><li>支撑分析与数据服务</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "实时计算与指标体系",
					"subtitle":    "大数据工程师",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：Flink / Kafka / Druid</li><li>搭建实时指标与查询服务</li></ul>",
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
					"title":       "计算与存储",
					"subtitle":    "平台",
					"description": "Hadoop / Spark / Flink / Kafka / Hive",
				},
				map[string]any{
					"id":          "s2",
					"title":       "数据开发与治理",
					"subtitle":    "工程",
					"description": "ETL / Scheduler / Quality / Lineage / Access",
				},
				map[string]any{
					"id":          "s3",
					"title":       "指标与服务",
					"subtitle":    "应用",
					"description": "Metrics / Data Service / Query",
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
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式计算与数据工程；参与平台与治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "中国科学技术大学",
					"major":       "信息与计算科学",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修计算与数据处理；参与校内项目。",
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
					"description": `<ul><li>数据湖示例：<a href="https://datalake.zhwei.invalid">datalake.zhwei.invalid</a></li><li>实时指标演示：<a href="https://realtime-metrics.zhwei.invalid">realtime-metrics.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateBigdataPreset() []byte {
	bdByte, _ := json.Marshal(bigdataPresetJSON)
	return bdByte
}
