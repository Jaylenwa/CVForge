package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var cCppPresetJSON = map[string]any{
	"title":    "C/C++ 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "C/C++ 后端开发工程师",
		"City":       "上海 / 杭州",
		"Money":      "30k-45k·14薪",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 年 C/C++ 服务端开发经验，熟悉 Linux 网络编程、并发与内存管理，具备低延迟与高吞吐系统优化经验</li><li>熟悉 TCP/UDP、epoll、线程池、无锁/细粒度锁设计，掌握性能分析与定位（perf/gprof/valgrind）</li><li>具备工程化能力：CMake/Make、单元测试、CI/CD，熟悉容器化部署与可观测（日志/指标/追踪）</li></ul>",
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
					"subtitle":    "C/C++ 后端开发工程师（低延迟交易）",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>负责撮合与行情分发核心模块（C++17），优化内存分配与对象复用，P99 延迟下降 35%</li><li>重构异步 I/O 与线程模型（epoll + 线程池），吞吐提升 1.8 倍且降低抖动</li><li>建设端到端压测与容量评估体系，完善熔断/限流与故障演练流程</li><li>完善核心链路可观测能力（日志/指标/追踪），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网基础架构团队",
					"subtitle":    "C++ 工程师（网关/中间件）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>参与自研 RPC 框架与服务网关（protobuf + gRPC-like），支撑数百个服务稳定运行</li><li>优化连接池、超时与重试策略，降低下游抖动对上游的影响</li><li>推动代码规范、单测与 CI 流程落地，降低回归成本</li></ul>",
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
					"title":       "低延迟消息总线与订阅系统",
					"subtitle":    "C/C++ 后端开发工程师",
					"timeStart":   "2022-08",
					"timeEnd":     "2023-05",
					"description": "<ul><li>技术栈：C++17 / epoll / protobuf / shared memory</li><li>设计发布-订阅协议与批量发送策略，降低系统调用与网络开销</li><li>引入零拷贝与对象池机制，显著降低 GC/alloc 压力与尾延迟</li><li>实现回放与重连机制，提升链路可恢复与可运维性</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "高性能日志与追踪采集器",
					"subtitle":    "C/C++ 后端开发工程师",
					"timeStart":   "2020-05",
					"timeEnd":     "2020-12",
					"description": "<ul><li>技术栈：C / C++ / Linux / mmap</li><li>实现批量写入与异步刷盘机制，降低 I/O 抖动</li><li>完善采样与背压策略，保障高峰期数据链路稳定</li></ul>",
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
					"title":       "语言与基础",
					"subtitle":    "C / C++",
					"description": "C++17 / STL / RAII / 多线程 / 内存管理",
				},
				map[string]any{
					"id":          "s2",
					"title":       "系统与网络",
					"subtitle":    "Linux",
					"description": "TCP/UDP / epoll / 线程池 / perf / valgrind / gdb",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程化",
					"subtitle":    "构建与交付",
					"description": "CMake / Make / 单元测试 / CI/CD / Docker",
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
					"title":       "北京邮电大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-07",
					"description": "研究方向为高性能网络与并发编程；参与实验室服务端性能优化项目。",
				},
			},
		},
	},
}

func GenerateCCppPreset() []byte {
	cCppByte, _ := json.Marshal(cCppPresetJSON)
	return cCppByte
}
