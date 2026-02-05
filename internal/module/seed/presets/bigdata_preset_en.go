package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var bigdataPresetJSONEn = map[string]any{
	"title":    "Big Data Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Na Li",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Big Data Engineer",
		"City":       "Beijing",
		"Money":      "35k-50k",
		"JoinTime":   "Available immediately",
		"Gender":     "Female",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>4 years in big data engineering; batch/stream and data lake/warehouse</li><li>Hadoop/Spark/Flink with Kafka</li><li>Data development and governance: ETL/scheduling/quality/lineage/access</li><li>Metrics systems and data services for analytics and business</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "Work Experience",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "Internet Company",
					"subtitle":    "Big Data Engineer (Development/Platform)",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Data lake and warehouse; layered architecture and quality governance</li><li>Batch/stream jobs and scheduling; analytics and services</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "E-commerce Company",
					"subtitle":    "Big Data Engineer (Real-time/Metrics)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Real-time computation and metrics systems</li><li>Data services and access governance</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "Projects",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "Data Lake & Warehouse Governance",
					"subtitle":    "Big Data Engineer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: Hadoop / Spark / Hive / Iceberg</li><li>Layering and quality; lineage and access</li><li>Analytics and data services</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Real-time Compute & Metrics System",
					"subtitle":    "Big Data Engineer",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Stack: Flink / Kafka / Druid</li><li>Real-time metrics and query services</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "Skills",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "Compute & Storage",
					"subtitle":    "Platforms",
					"description": "Hadoop / Spark / Flink / Kafka / Hive",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Data Dev & Governance",
					"subtitle":    "Engineering",
					"description": "ETL / Scheduler / Quality / Lineage / Access",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Metrics & Services",
					"subtitle":    "Applications",
					"description": "Metrics / Data Service / Query",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "Education",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "University of Science and Technology of China",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "Distributed computing and data engineering; platform and governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "University of Science and Technology of China",
					"major":       "Information and Computing Science",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Computing and data processing; campus projects.",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "Portfolio",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>Data lake examples: <a href="https://datalake.zhwei.invalid">datalake.zhwei.invalid</a></li><li>Real-time metrics demos: <a href="https://realtime-metrics.zhwei.invalid">realtime-metrics.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateBigdataPresetEn() []byte {
	bdByte, _ := json.Marshal(bigdataPresetJSONEn)
	return bdByte
}
