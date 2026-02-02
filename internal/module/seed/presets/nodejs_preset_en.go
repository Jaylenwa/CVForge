package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var nodejsPresetJSONEn = map[string]any{
	"title":    "Node.js Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Node.js Backend Developer",
		"City":       "Shanghai / Hangzhou",
		"Money":      "25k-40k · 14-month salary",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#16a34a",
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
					"description": "<ul><li>5 years of Node.js backend experience; proficient with TypeScript and NestJS/Express, with hands-on optimization for high-concurrency APIs and reliability improvements</li><li>Hands-on with MySQL/PostgreSQL, Redis, and Kafka/RabbitMQ; experienced in cache consistency, reliable messaging, idempotency, and retry design</li><li>Engineering and cloud-native practices: Docker/Kubernetes, CI/CD, metrics/logging/tracing</li></ul>",
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
					"title":       "Online Education Platform",
					"subtitle":    "Node.js Backend Developer (Courses/Payments)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned core course and payment services (NestJS + TypeScript); standardized APIs and unified authentication to improve iteration efficiency</li><li>Implemented Redis caching and hot-key mitigation, reducing P95 latency from 230ms to 140ms</li><li>Introduced Kafka-based async notifications and reconciliation flows, reducing critical-path latency and improving reliability</li><li>Built observability (Prometheus/Grafana + tracing), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community Company",
					"subtitle":    "Node.js Engineer (Platform/Middle Platform)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Built account and authorization platform (Express + TypeScript + PostgreSQL) to support multiple business lines</li><li>Optimized slow queries and indexing strategies, reducing critical query time by 60%</li><li>Established CI/CD and canary release processes to lower release risk</li></ul>",
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
					"title":       "Unified Notification and Task Scheduling Platform",
					"subtitle":    "Node.js Backend Developer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Tech stack: Node.js / TypeScript / NestJS / Redis / Kafka</li><li>Abstracted task templates and retry/idempotency mechanisms to reduce duplicated work and improve reliability</li><li>Added task observability: status tracking, failure aggregation, and alert notifications</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "High-Concurrency Campaign and Coupon System",
					"subtitle":    "Node.js Backend Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Tech stack: Node.js / Redis / RabbitMQ</li><li>Designed idempotent inventory deduction and issuance flows to prevent duplicate claims and overselling</li><li>Used Redis Lua atomic scripts and hot-key governance, reducing P95 latency by 40% during peaks</li></ul>",
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
					"title":       "Languages & Frameworks",
					"subtitle":    "Node.js",
					"description": "TypeScript / NestJS / Express / REST API / OpenAPI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Storage & Middleware",
					"subtitle":    "Database / Cache / Messaging",
					"description": "MySQL / PostgreSQL / Redis / Kafka / RabbitMQ",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering & Reliability",
					"subtitle":    "Delivery & Observability",
					"description": "Docker / Kubernetes / CI/CD / Prometheus / Grafana",
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
					"title":       "Shanghai Jiao Tong University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on distributed systems and engineering; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Shanghai Jiao Tong University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied data structures, databases, and computer networks; participated in campus projects and competitions.",
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
					"description": `<ul><li>Open-source tools: <a href="https://npm.zhwei.invalid">npm.zhwei.invalid</a></li><li>API demo site: <a href="https://api-demo.zhwei.invalid">api-demo.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateNodejsPresetEn() []byte {
	nodejsByte, _ := json.Marshal(nodejsPresetJSONEn)
	return nodejsByte
}
