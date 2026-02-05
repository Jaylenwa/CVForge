package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var nlpPresetJSON = map[string]any{
	"title":    "自然语言处理工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "李娜",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "自然语言处理工程师",
		"City":       "北京",
		"Money":      "35k-50k",
		"JoinTime":   "随时到岗",
		"Gender":     "女",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>4 年 NLP 经验，熟悉分词/标注/文本分类/序列标注</li><li>掌握 Transformer/BERT/LLM 微调与部署</li><li>熟悉数据标注与清洗、评测指标与实验设计</li><li>具备工程化与推理优化：量化/裁剪/蒸馏与服务治理</li></ul>",
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
					"subtitle":    "NLP 工程师（模型/平台）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责文本分类与问答任务，优化模型精度与性能</li><li>搭建微调与评测流水线，完善数据治理</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某技术公司",
					"subtitle":    "NLP 工程师（LLM/推理）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>研究与落地 LLM 微调与推理优化</li><li>建设服务与治理，提升稳定性与效率</li></ul>",
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
					"title":       "文本分类与问答系统",
					"subtitle":    "自然语言处理工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：Python / PyTorch / Transformers</li><li>搭建训练与评测流水线，优化模型与指标</li><li>落地服务与治理</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "LLM 微调与推理优化",
					"subtitle":    "自然语言处理工程师",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：Transformers / Quantization / Distillation</li><li>优化推理性能与成本，提升稳定性与响应</li></ul>",
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
					"title":       "模型与算法",
					"subtitle":    "NLP",
					"description": "Transformer / BERT / LLM / CRF",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与平台",
					"subtitle":    "AI 工程",
					"description": "PyTorch / Transformers / Serving",
				},
				map[string]any{
					"id":          "s3",
					"title":       "数据与评测",
					"subtitle":    "质量",
					"description": "Label / Clean / Metrics / AB Test",
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
					"title":       "北京大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为自然语言处理与机器学习；参与竞赛与项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "北京大学",
					"major":       "信息工程",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修计算机与信息处理；参与校内项目与竞赛。",
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
					"description": `<ul><li>NLP 示例：<a href="https://nlp-examples.zhwei.invalid">nlp-examples.zhwei.invalid</a></li><li>LLM 微调演示：<a href="https://llm-finetune.zhwei.invalid">llm-finetune.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateNlpPreset() []byte {
	nlpByte, _ := json.Marshal(nlpPresetJSON)
	return nlpByte
}
