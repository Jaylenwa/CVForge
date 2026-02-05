package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var mlAiPresetJSON = map[string]any{
	"title":    "机器学习/AI工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "李娜",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "机器学习/AI工程师",
		"City":       "北京",
		"Money":      "35k-50k",
		"JoinTime":   "随时到岗",
		"Gender":     "女",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#a855f7",
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
					"description": "<ul><li>4 年 ML/AI 经验，熟悉监督/非监督、深度学习与部署</li><li>掌握 PyTorch/TensorFlow，训练与评估流程、MLOps</li><li>具备特征工程与数据治理，指标体系与 A/B 实验</li><li>有推理与服务治理经验：性能优化与成本控制</li></ul>",
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
					"title":       "某互联网公司",
					"subtitle":    "机器学习/AI工程师（建模/平台）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责推荐/分类/排序等任务，建设训练与评估体系</li><li>推进 MLOps 与模型治理，完善服务与监控</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某技术公司",
					"subtitle":    "机器学习/AI工程师（推理/交付）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>优化推理性能与成本，提升稳定性与响应</li><li>完善交付与治理流程</li></ul>",
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
					"title":       "推荐与排序系统",
					"subtitle":    "机器学习/AI工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：Python / PyTorch / TensorFlow / MLFlow</li><li>建设训练与评估体系，落地 A/B 实验</li><li>完善服务与监控治理</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "推理性能与成本优化",
					"subtitle":    "机器学习/AI工程师",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：TensorRT / ONNX / Serving</li><li>优化吞吐与延迟，控制成本</li></ul>",
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
					"title":       "建模与训练",
					"subtitle":    "AI",
					"description": "PyTorch / TensorFlow / Sklearn",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与治理",
					"subtitle":    "MLOps",
					"description": "MLFlow / Serving / Monitoring",
				},
				map[string]any{
					"id":          "s3",
					"title":       "推理与优化",
					"subtitle":    "性能",
					"description": "TensorRT / ONNX / Quantization",
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
					"title":       "清华大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为机器学习与工程化；参与平台与治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "清华大学",
					"major":       "自动化",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修控制与计算；参与校内项目与竞赛。",
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
					"description": `<ul><li>推荐系统示例：<a href="https://recsys.zhwei.invalid">recsys.zhwei.invalid</a></li><li>模型服务演示：<a href="https://ai-serving.zhwei.invalid">ai-serving.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateMlAiPreset() []byte {
	mlByte, _ := json.Marshal(mlAiPresetJSON)
	return mlByte
}
