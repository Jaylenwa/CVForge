package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var performanceTestPresetJSON = map[string]any{
	"title":    "性能测试工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "性能测试工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f43f5e",
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
					"description": "<ul><li>5 年性能测试经验，熟悉压测/容量评估/稳定性测试，能定位与分析性能瓶颈</li><li>掌握工具与框架：JMeter/Gatling/Locust，脚本设计与数据准备</li><li>熟悉指标与可观测：TPS/RT/错误率/资源利用率，Prometheus/Grafana/APM</li><li>具备性能治理与优化实践：缓存/并发/连接池/异步化与降级</li></ul>",
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
					"subtitle":    "性能测试工程师（交易/订单）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>构建交易与订单链路压测方案，评估容量与瓶颈</li><li>沉淀场景与数据准备、脚本复用与结果分析</li><li>推动优化与治理，显著提升稳定性与吞吐</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区",
					"subtitle":    "性能测试工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>搭建用户增长与消息等场景压测方案，完善指标与告警</li><li>与研发协同定位瓶颈并优化资源与并发策略</li></ul>",
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
					"title":       "统一压测平台与报告系统",
					"subtitle":    "性能测试工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：JMeter / Grafana / Prometheus</li><li>沉淀场景管理与数据准备、报告与告警</li><li>接入 APM 与链路追踪，定位瓶颈与优化路径</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "容量评估与稳定性治理",
					"subtitle":    "性能测试工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：JMeter / Locust / APM</li><li>开展容量评估与稳定性测试，提出治理方案并落地</li><li>实现压测数据管理与复用</li></ul>",
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
					"title":       "工具与框架",
					"subtitle":    "压测",
					"description": "JMeter / Gatling / Locust",
				},
				map[string]any{
					"id":          "s2",
					"title":       "指标与可观测",
					"subtitle":    "监控",
					"description": "TPS / RT / 错误率 / 资源利用率 / Prometheus / Grafana / APM",
				},
				map[string]any{
					"id":          "s3",
					"title":       "治理与优化",
					"subtitle":    "稳定性",
					"description": "缓存 / 并发 / 连接池 / 异步化 / 降级",
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
					"title":       "浙江大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为性能与稳定性工程；参与压测平台项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "浙江大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修操作系统与网络；参与校内项目与竞赛。",
				},
			},
		},
	},
}

func GeneratePerformanceTestPreset() []byte {
	perfByte, _ := json.Marshal(performanceTestPresetJSON)
	return perfByte
}
