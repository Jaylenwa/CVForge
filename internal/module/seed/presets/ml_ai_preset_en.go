package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var mlAiPresetJSONEn = map[string]any{
	"title":    "Machine Learning/AI Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Na Li",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Machine Learning/AI Engineer",
		"City":       "Beijing",
		"Money":      "35k-50k",
		"JoinTime":   "Available immediately",
		"Gender":     "Female",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>4 years in ML/AI; supervised/unsupervised, deep learning and deployment</li><li>PyTorch/TensorFlow; training/evaluation pipelines and MLOps</li><li>Feature engineering and data governance; metrics and A/B experiments</li><li>Inference and service governance: performance optimization and cost control</li></ul>",
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
					"subtitle":    "ML/AI Engineer (Modeling/Platform)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Recommendation/classification/ranking; training and evaluation systems</li><li>MLOps and model governance; serving and monitoring</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Technology Company",
					"subtitle":    "ML/AI Engineer (Inference/Delivery)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Inference performance and cost optimization</li><li>Delivery and governance processes</li></ul>",
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
					"title":       "Recommendation & Ranking System",
					"subtitle":    "Machine Learning/AI Engineer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: Python / PyTorch / TensorFlow / MLFlow</li><li>Training and evaluation systems; A/B experiments</li><li>Serving and monitoring governance</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Inference Performance & Cost Optimization",
					"subtitle":    "Machine Learning/AI Engineer",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Stack: TensorRT / ONNX / Serving</li><li>Throughput/latency improvements and cost control</li></ul>",
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
					"title":       "Modeling & Training",
					"subtitle":    "AI",
					"description": "PyTorch / TensorFlow / Sklearn",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Engineering & Governance",
					"subtitle":    "MLOps",
					"description": "MLFlow / Serving / Monitoring",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Inference & Optimization",
					"subtitle":    "Performance",
					"description": "TensorRT / ONNX / Quantization",
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
					"title":       "Tsinghua University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "Machine learning and engineering; platform and governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Tsinghua University",
					"major":       "Automation",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Control and computing; campus projects and competitions.",
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
					"description": `<ul><li>Recommendation system demos: <a href="https://recsys.zhwei.invalid">recsys.zhwei.invalid</a></li><li>Model serving demos: <a href="https://ai-serving.zhwei.invalid">ai-serving.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateMlAiPresetEn() []byte {
	mlByte, _ := json.Marshal(mlAiPresetJSONEn)
	return mlByte
}
