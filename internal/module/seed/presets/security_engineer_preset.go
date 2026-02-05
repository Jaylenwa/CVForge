package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var securityEngineerPresetJSON = map[string]any{
	"title":    "安全工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "安全工程师",
		"City":       "上海",
		"Money":      "30k-45k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "29",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#ef4444",
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
					"description": "<ul><li>5 年安全工程经验，具备漏洞扫描与渗透测试能力</li><li>熟悉应用与平台安全：WAF/IDS/IPS、权限与审计、风控策略</li><li>掌握安全开发流程与合规：SDL/DevSecOps、等级保护与合规</li><li>熟悉攻防演练与应急响应，保障业务与平台安全</li></ul>",
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
					"subtitle":    "安全工程师（应用/平台）",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>建设安全防护与审计体系，完善权限与访问控制</li><li>推动 SDL 与 DevSecOps 落地，提升研发安全性</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网公司",
					"subtitle":    "安全工程师（渗透/演练）",
					"timeStart":   "2018-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>组织攻防演练与应急响应</li><li>完善风控与合规流程</li></ul>",
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
					"title":       "应用与平台安全治理",
					"subtitle":    "安全工程师",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：WAF / IDS / IPS / SIEM</li><li>建设防护与审计，完善风控与合规</li><li>推动 SDL 与 DevSecOps 流程</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "攻防演练与应急响应",
					"subtitle":    "安全工程师",
					"timeStart":   "2019-08",
					"timeEnd":     "2020-04",
					"description": "<ul><li>技术栈：Nmap / Burp / Metasploit</li><li>组织攻防演练与响应流程</li><li>沉淀风险与处置手册</li></ul>",
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
					"title":       "应用与平台安全",
					"subtitle":    "防护",
					"description": "WAF / IDS / IPS / SIEM / IAM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "渗透与测试",
					"subtitle":    "攻防",
					"description": "Nmap / Burp / Metasploit / OWASP",
				},
				map[string]any{
					"id":          "s3",
					"title":       "流程与合规",
					"subtitle":    "治理",
					"description": "SDL / DevSecOps / 等保 / 合规",
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
					"title":       "上海交通大学",
					"major":       "信息安全",
					"degree":      "硕士",
					"timeStart":   "2015-09",
					"timeEnd":     "2018-06",
					"description": "研究方向为应用与平台安全；参与攻防演练与治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "上海交通大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "主修计算机基础与网络安全；参与校内项目与竞赛。",
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
					"description": `<ul><li>安全工具示例：<a href="https://sec-tools.zhwei.invalid">sec-tools.zhwei.invalid</a></li><li>安全手册：<a href="https://sec-handbook.zhwei.invalid">sec-handbook.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateSecurityEngineerPreset() []byte {
	seByte, _ := json.Marshal(securityEngineerPresetJSON)
	return seByte
}
