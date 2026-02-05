package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dataMiningPresetJSONEn = map[string]any{
	"title":    "Data Mining Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Na Li",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Data Mining Engineer",
		"City":       "Beijing",
		"Money":      "30k-45k",
		"JoinTime":   "Available immediately",
		"Gender":     "Female",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>4 years in data mining; feature engineering and model evaluation</li><li>SQL/Python/Pandas; ETL and data quality governance</li><li>Familiar with Spark/Flink; batch/stream processing and metrics systems</li><li>Business understanding and metric adoption to drive decisions</li></ul>",
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
					"subtitle":    "Data Mining Engineer (Modeling/Analytics)",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>User profiles and feature store; training data design</li><li>Retention/conversion/churn modeling and evaluation</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "E-commerce Company",
					"subtitle":    "Data Mining Engineer (ETL/Reporting)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Data warehouse and ETL jobs; data quality improvement</li><li>Metric and reporting systems to support analytics</li></ul>",
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
					"title":       "User Churn Prediction & Retention",
					"subtitle":    "Data Mining Engineer",
					"timeStart":   "2022-01",
					"timeEnd":     "2022-10",
					"description": "<ul><li>Stack: Python / Pandas / Sklearn</li><li>Feature engineering and modeling to improve retention</li><li>A/B experiments and metric evaluation</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Data Warehouse & Metrics System",
					"subtitle":    "Data Mining Engineer",
					"timeStart":   "2020-05",
					"timeEnd":     "2020-12",
					"description": "<ul><li>Stack: SQL / Spark / Airflow</li><li>ETL and data quality; analytics and reporting support</li></ul>",
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
					"title":       "Data & Modeling",
					"subtitle":    "Core",
					"description": "SQL / Python / Pandas / Sklearn",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Compute & Engineering",
					"subtitle":    "Platforms",
					"description": "Spark / Flink / Airflow",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Metrics & Experimentation",
					"subtitle":    "Analytics",
					"description": "Metrics / AB Test / Evaluation",
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
					"title":       "Renmin University of China",
					"major":       "Statistics",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "Data analysis and modeling; competitions and projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Renmin University of China",
					"major":       "Mathematics and Applied Mathematics",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Probability and numerical analysis; campus projects.",
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
					"description": `<ul><li>Data mining examples: <a href="https://datamining.zhwei.invalid">datamining.zhwei.invalid</a></li><li>AB testing examples: <a href="https://abtest.zhwei.invalid">abtest.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDataMiningPresetEn() []byte {
	dmByte, _ := json.Marshal(dataMiningPresetJSONEn)
	return dmByte
}
