package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var pythonPresetJSONEn = map[string]any{
	"title":    "Python Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Python Backend Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#2563eb",
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
					"description": "<ul><li>5 years of Python backend experience; familiar with FastAPI/Django, capable of delivering products from 0 to 1 and optimizing performance</li><li>Hands-on with PostgreSQL/MySQL, Redis, Kafka, and Celery; able to build high-concurrency and highly available services</li><li>Engineering and cloud-native practices: Docker, Kubernetes, CI/CD, logging/monitoring/alerting</li></ul>",
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
					"title":       "Leading E-commerce Platform",
					"subtitle":    "Python Backend Developer (Trading/Risk Control)",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Owned core trading-domain APIs and async flows (FastAPI + Kafka + Celery), raising order creation success rate to 99.95%</li><li>Implemented multi-layer Redis caching and hot-key mitigation, reducing P95 latency from 240ms to 130ms</li><li>Led risk-rule engine refactoring with configuration-driven rules and canary release, tripling iteration efficiency</li><li>Improved monitoring and alerting (Prometheus/Grafana), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "SaaS Startup",
					"subtitle":    "Python Backend Engineer (Data & Platform)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Built multi-tenant IAM, audit logs, and configuration center with Django to support 100+ enterprise customers</li><li>Implemented ETL processing and workflow orchestration (Celery + PostgreSQL), reducing report generation time by 40%</li><li>Contributed to API standards and automated testing (OpenAPI + pytest), significantly improving regression efficiency</li></ul>",
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
					"title":       "Real-Time Risk Control Rule Engine",
					"subtitle":    "Python Backend Developer",
					"timeStart":   "2022-06",
					"timeEnd":     "2023-03",
					"description": "<ul><li>Tech stack: FastAPI / Kafka / Redis / PostgreSQL</li><li>Designed a rules DSL and executor to support dynamic rollout by merchant/audience and fast rollback</li><li>Integrated Kafka streaming events and Redis profiling cache to enable millisecond-level decisions</li><li>Improved load testing and capacity planning, doubling peak QPS while keeping the system stable</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Task Orchestration and Async Processing Platform",
					"subtitle":    "Python Backend Developer",
					"timeStart":   "2020-03",
					"timeEnd":     "2020-11",
					"description": "<ul><li>Tech stack: Django / Celery / Redis</li><li>Abstracted task templates and retry/idempotency mechanisms to reduce duplicated work</li><li>Added task observability: status tracking, failure aggregation, and alert notifications</li></ul>",
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
					"subtitle":    "Python",
					"description": "FastAPI / Django / REST API / OpenAPI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Storage & Middleware",
					"subtitle":    "Database / Cache / Messaging",
					"description": "PostgreSQL / MySQL / Redis / Kafka / Celery",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering",
					"subtitle":    "Observability & Delivery",
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
					"title":       "Huazhong University of Science and Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on distributed systems and highly available architectures; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Huazhong University of Science and Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied data structures, databases, and computer networks; participated in campus projects and competitions.",
				},
			},
		},
	},
}

func GeneratePythonPresetEn() []byte {
	pythonByte, _ := json.Marshal(pythonPresetJSONEn)
	return pythonByte
}
