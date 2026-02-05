package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dbaPresetJSON = map[string]any{
	"title":    "数据库管理员（DBA）简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "数据库管理员（DBA）",
		"City":       "上海",
		"Money":      "30k-45k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "29",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>6 年 DBA 经验，熟悉 MySQL/PostgreSQL 运维、复制与备份</li><li>性能与容量优化：索引/执行计划、连接池与并发、分库分表</li><li>高可用与容灾：MHA/Orchestrator、热备与演练</li><li>安全与合规：权限/审计/加密，数据生命周期管理</li></ul>",
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
					"title":       "某金融科技公司",
					"subtitle":    "DBA（高可用/性能）",
					"timeStart":   "2021-03",
					"today":       true,
					"description": "<ul><li>治理 MySQL 主从与高可用，完善监控与告警</li><li>优化慢查询与索引策略，提升核心交易性能</li><li>设计备份与容灾演练流程</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网公司",
					"subtitle":    "DBA（容量/架构）",
					"timeStart":   "2017-07",
					"timeEnd":     "2021-02",
					"description": "<ul><li>设计分库分表与数据治理方案</li><li>沉淀规范与变更流程，保障数据质量与安全</li></ul>",
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
					"title":       "数据库高可用平台",
					"subtitle":    "数据库管理员（DBA）",
					"timeStart":   "2022-04",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：MySQL / Orchestrator / MHA</li><li>建设高可用与容灾演练体系</li><li>打通监控与告警，完善故障响应</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "性能与容量优化",
					"subtitle":    "数据库管理员（DBA）",
					"timeStart":   "2019-06",
					"timeEnd":     "2020-03",
					"description": "<ul><li>技术栈：MySQL / PostgreSQL / Proxy</li><li>优化慢查询与索引策略，提升 QPS 与响应</li><li>设计分库分表与治理方案</li></ul>",
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
					"title":       "数据库运维",
					"subtitle":    "核心",
					"description": "MySQL / PostgreSQL / Replication / Backup",
				},
				map[string]any{
					"id":          "s2",
					"title":       "性能与容量",
					"subtitle":    "优化",
					"description": "Index / Execution Plan / Connection Pool / Sharding",
				},
				map[string]any{
					"id":          "s3",
					"title":       "高可用与安全",
					"subtitle":    "治理",
					"description": "MHA / Orchestrator / Audit / Encryption",
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
					"timeStart":   "2015-09",
					"timeEnd":     "2017-06",
					"description": "研究方向为数据库系统与可靠性工程。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "南京大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "主修数据结构与数据库；参与校内项目。",
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
					"description": `<ul><li>SQL 优化示例：<a href="https://sql-tuning.zhwei.invalid">sql-tuning.zhwei.invalid</a></li><li>备份脚本集：<a href="https://dba-backup.zhwei.invalid">dba-backup.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDbaPreset() []byte {
	dbaByte, _ := json.Marshal(dbaPresetJSON)
	return dbaByte
}
