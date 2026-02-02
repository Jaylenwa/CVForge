package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var golangPresetJSONEn = map[string]any{
	"title":    "Go Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Go Backend Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#14b8a6",
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
					"description": "<ul><li>4+ years of Go backend experience; familiar with microservices, gRPC, distributed consistency, and performance tuning for high-concurrency systems</li><li>Hands-on with MySQL/Redis/Kafka; experienced in performance optimization, incident response, and reliability improvements</li><li>Observability and cloud-native practices: Prometheus/Grafana, tracing, Docker/Kubernetes</li></ul>",
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
					"title":       "Large Internet Company",
					"subtitle":    "Go Backend Developer (Orders/Fulfillment)",
					"timeStart":   "2021-08",
					"today":       true,
					"description": "<ul><li>Owned core fulfillment services (Go + gRPC), supporting tens of millions of requests daily and reducing P95 latency by 35%</li><li>Built idempotency, distributed locks, and compensation mechanisms to ensure consistency and recoverability under high concurrency</li><li>Drove Redis hot-key mitigation and degradation strategies, significantly improving critical-path stability</li><li>Improved observability (metrics/logs/traces), reducing major incident recovery time from hours to under 20 minutes</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "FinTech Company",
					"subtitle":    "Backend Engineer (Gateway/Platform)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-07",
					"description": "<ul><li>Worked on a unified API gateway (auth, rate limiting, canary release, circuit breaking), improving onboarding efficiency</li><li>Implemented Kafka async flows and delay queues to decouple core services and smooth traffic spikes</li><li>Drove database governance: slow SQL optimization, indexing, sharding strategies; reduced critical query time by 60%</li></ul>",
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
					"title":       "Unified Authentication and Authorization Platform",
					"subtitle":    "Go Backend Developer",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-10",
					"description": "<ul><li>Tech stack: Go / gRPC / Redis / MySQL</li><li>Designed a hybrid RBAC/ABAC model and permission cache strategy to balance consistency and performance</li><li>Implemented token rotation and multi-device session control, adding audit and alerting capabilities</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Order Fulfillment Pipeline Refactoring",
					"subtitle":    "Go Backend Developer",
					"timeStart":   "2023-01",
					"timeEnd":     "2023-12",
					"description": "<ul><li>Tech stack: Go / Kafka / Redis</li><li>Split fulfillment sub-domains and introduced event-driven workflows to reduce coupling and improve evolvability</li><li>Added Saga-style compensation and retry strategies, significantly lowering abnormal order ratio</li></ul>",
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
					"title":       "Backend",
					"subtitle":    "Go",
					"description": "Gin / gRPC / Concurrency & Performance Tuning / Microservices Architecture",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Databases & Middleware",
					"subtitle":    "MySQL / Redis / Kafka",
					"description": "Transactions & Locks, Index Optimization, Cache Consistency, Reliable Messaging",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Reliability & Cloud Native",
					"subtitle":    "Observability",
					"description": "Prometheus / Grafana / Logging & Tracing / Docker / Kubernetes",
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
					"title":       "Sun Yat-sen University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on distributed systems and high-concurrency services; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Sun Yat-sen University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied software engineering and distributed systems; participated in campus project practice.",
				},
			},
		},
	},
}

func GenerateGolangPresetEn() []byte {
	golangByte, _ := json.Marshal(golangPresetJSONEn)
	return golangByte
}
