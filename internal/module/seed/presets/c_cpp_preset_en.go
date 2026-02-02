package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var cCppPresetJSONEn = map[string]any{
	"title":    "C/C++ Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "C/C++ Backend Developer",
		"City":       "Shanghai / Hangzhou",
		"Money":      "30k-45k · 14-month salary",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#0ea5e9",
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
					"description": "<ul><li>5 years of C/C++ backend experience; proficient in Linux network programming, concurrency, and memory management, with hands-on optimization for low-latency and high-throughput systems</li><li>Solid understanding of TCP/UDP, epoll, thread pools, lock-free/fine-grained locking; experienced in performance profiling and debugging (perf/gprof/valgrind)</li><li>Strong engineering skills: CMake/Make, unit testing, CI/CD; familiar with containerized deployment and observability (logs/metrics/tracing)</li></ul>",
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
					"subtitle":    "C/C++ Backend Developer (Low-Latency Trading)",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>Owned matching and market-data distribution modules (C++17); optimized memory allocation and object reuse, reducing P99 latency by 35%</li><li>Refactored async I/O and threading model (epoll + thread pool), increasing throughput by 1.8x while reducing jitter</li><li>Built end-to-end load testing and capacity planning; improved circuit breaking/rate limiting and incident drills</li><li>Enhanced observability on critical paths (logs/metrics/tracing), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Infrastructure Team",
					"subtitle":    "C++ Engineer (Gateway/Middleware)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Worked on an in-house RPC framework and service gateway (protobuf + gRPC-like), supporting stable operations for hundreds of services</li><li>Optimized connection pooling, timeouts, and retry strategies to reduce downstream jitter impact</li><li>Drove coding standards, unit tests, and CI adoption, lowering regression cost</li></ul>",
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
					"title":       "Low-Latency Message Bus and Subscription System",
					"subtitle":    "C/C++ Backend Developer",
					"timeStart":   "2022-08",
					"timeEnd":     "2023-05",
					"description": "<ul><li>Tech stack: C++17 / epoll / protobuf / shared memory</li><li>Designed pub/sub protocol and batch sending to reduce syscalls and network overhead</li><li>Introduced zero-copy and object pooling to reduce alloc pressure and tail latency</li><li>Implemented replay and reconnect mechanisms to improve resilience and operability</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "High-Performance Logging and Tracing Collector",
					"subtitle":    "C/C++ Backend Developer",
					"timeStart":   "2020-05",
					"timeEnd":     "2020-12",
					"description": "<ul><li>Tech stack: C / C++ / Linux / mmap</li><li>Implemented batch writes and async flush to reduce I/O jitter</li><li>Improved sampling and backpressure strategies to keep data pipelines stable during peaks</li></ul>",
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
					"title":       "Languages & Fundamentals",
					"subtitle":    "C / C++",
					"description": "C++17 / STL / RAII / Multithreading / Memory Management",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Systems & Networking",
					"subtitle":    "Linux",
					"description": "TCP/UDP / epoll / Thread Pool / perf / valgrind / gdb",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering",
					"subtitle":    "Build & Delivery",
					"description": "CMake / Make / Unit Testing / CI/CD / Docker",
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
					"title":       "Beijing University of Posts and Telecommunications",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on high-performance networking and concurrent programming; participated in lab performance optimization projects for server systems.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Beijing University of Posts and Telecommunications",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied data structures, operating systems, and computer networks; participated in campus systems programming projects.",
				},
			},
		},
	},
}

func GenerateCCppPresetEn() []byte {
	cCppByte, _ := json.Marshal(cCppPresetJSONEn)
	return cCppByte
}
