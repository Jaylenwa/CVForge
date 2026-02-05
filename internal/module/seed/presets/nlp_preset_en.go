package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var nlpPresetJSONEn = map[string]any{
	"title":    "NLP Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Na Li",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "NLP Engineer",
		"City":       "Beijing",
		"Money":      "35k-50k",
		"JoinTime":   "Available immediately",
		"Gender":     "Female",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#4f46e5",
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
					"description": "<ul><li>4 years in NLP; tokenization/annotation/text classification/sequence labeling</li><li>Transformer/BERT/LLM fine-tuning and deployment</li><li>Data labeling/cleaning; metrics and experiment design</li><li>Engineering and inference optimization: quantization/pruning/distillation; service governance</li></ul>",
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
					"subtitle":    "NLP Engineer (Models/Platform)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Text classification and QA; model accuracy and performance</li><li>Fine-tuning and evaluation pipelines; data governance</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Technology Company",
					"subtitle":    "NLP Engineer (LLM/Inference)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>LLM fine-tuning and inference optimization</li><li>Service building and governance for stability and efficiency</li></ul>",
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
					"title":       "Text Classification & QA System",
					"subtitle":    "NLP Engineer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: Python / PyTorch / Transformers</li><li>Training and evaluation pipelines; model/metric optimization</li><li>Service and governance</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "LLM Fine-tuning & Inference Optimization",
					"subtitle":    "NLP Engineer",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Stack: Transformers / Quantization / Distillation</li><li>Inference performance and cost improvements</li></ul>",
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
					"title":       "Models & Algorithms",
					"subtitle":    "NLP",
					"description": "Transformer / BERT / LLM / CRF",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Engineering & Platform",
					"subtitle":    "AI Engineering",
					"description": "PyTorch / Transformers / Serving",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Data & Evaluation",
					"subtitle":    "Quality",
					"description": "Label / Clean / Metrics / AB Test",
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
					"title":       "Peking University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "NLP and machine learning; competitions and projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Peking University",
					"major":       "Information Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Computing and information processing; campus projects and competitions.",
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
					"description": `<ul><li>NLP examples: <a href="https://nlp-examples.zhwei.invalid">nlp-examples.zhwei.invalid</a></li><li>LLM fine-tuning demos: <a href="https://llm-finetune.zhwei.invalid">llm-finetune.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateNlpPresetEn() []byte {
	nlpByte, _ := json.Marshal(nlpPresetJSONEn)
	return nlpByte
}
