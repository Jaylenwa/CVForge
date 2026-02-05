package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dbaPresetJSONEn = map[string]any{
	"title":    "Database Administrator (DBA) Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Database Administrator (DBA)",
		"City":       "Shanghai",
		"Money":      "30k-45k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "29",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#9333ea",
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
					"description": "<ul><li>6 years of DBA experience; MySQL/PostgreSQL operations, replication and backup</li><li>Performance & capacity: indexes/plans, connection pools/concurrency, sharding</li><li>High availability & DR: MHA/Orchestrator, hot standby and drills</li><li>Security & compliance: permissions/audit/encryption; data lifecycle management</li></ul>",
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
					"title":       "FinTech Company",
					"subtitle":    "DBA (HA/Performance)",
					"timeStart":   "2021-03",
					"today":       true,
					"description": "<ul><li>Governed MySQL replication and HA; improved monitoring and alerting</li><li>Optimized slow queries and indexing; improved core transaction performance</li><li>Designed backup and DR drill processes</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Company",
					"subtitle":    "DBA (Capacity/Architecture)",
					"timeStart":   "2017-07",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Sharding and data governance design</li><li>Standards and change processes for data quality and security</li></ul>",
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
					"title":       "Database High Availability Platform",
					"subtitle":    "Database Administrator (DBA)",
					"timeStart":   "2022-04",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: MySQL / Orchestrator / MHA</li><li>Built HA and DR drill systems</li><li>Integrated monitoring and alerting for incident response</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Performance & Capacity Optimization",
					"subtitle":    "Database Administrator (DBA)",
					"timeStart":   "2019-06",
					"timeEnd":     "2020-03",
					"description": "<ul><li>Stack: MySQL / PostgreSQL / Proxy</li><li>Slow query/index optimization for QPS/latency improvements</li><li>Sharding and governance design</li></ul>",
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
					"title":       "Database Operations",
					"subtitle":    "Core",
					"description": "MySQL / PostgreSQL / Replication / Backup",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Performance & Capacity",
					"subtitle":    "Optimization",
					"description": "Index / Execution Plan / Connection Pool / Sharding",
				},
				map[string]any{
					"id":          "s3",
					"title":       "High Availability & Security",
					"subtitle":    "Governance",
					"description": "MHA / Orchestrator / Audit / Encryption",
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
					"title":       "Nanjing University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2015-09",
					"timeEnd":     "2017-06",
					"description": "Database systems and reliability engineering.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Nanjing University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "Data structures and databases; campus projects.",
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
					"description": `<ul><li>SQL tuning examples: <a href="https://sql-tuning.zhwei.invalid">sql-tuning.zhwei.invalid</a></li><li>Backup scripts: <a href="https://dba-backup.zhwei.invalid">dba-backup.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDbaPresetEn() []byte {
	dbaByte, _ := json.Marshal(dbaPresetJSONEn)
	return dbaByte
}
