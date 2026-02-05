package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var securityEngineerPresetJSONEn = map[string]any{
	"title":    "Security Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Security Engineer",
		"City":       "Shanghai",
		"Money":      "30k-45k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "29",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 years in security engineering; vulnerability scanning and penetration testing</li><li>Application/platform security: WAF/IDS/IPS, permissions and audits, risk control</li><li>Secure development and compliance: SDL/DevSecOps; classified protection and compliance</li><li>Red/blue exercises and incident response to protect platforms and businesses</li></ul>",
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
					"subtitle":    "Security Engineer (Application/Platform)",
					"timeStart":   "2021-04",
					"today":       true,
					"description": "<ul><li>Security defense and audit; permissions and access control</li><li>SDL and DevSecOps practices to improve development security</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Company",
					"subtitle":    "Security Engineer (Pentest/Exercises)",
					"timeStart":   "2018-07",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Organized red/blue team exercises and incident response</li><li>Risk control and compliance processes</li></ul>",
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
					"title":       "Application & Platform Security Governance",
					"subtitle":    "Security Engineer",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: WAF / IDS / IPS / SIEM</li><li>Defense and audit; risk control and compliance</li><li>SDL and DevSecOps workflows</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Red/Blue Exercises & Incident Response",
					"subtitle":    "Security Engineer",
					"timeStart":   "2019-08",
					"timeEnd":     "2020-04",
					"description": "<ul><li>Stack: Nmap / Burp / Metasploit</li><li>Exercise design and incident procedures</li><li>Risk and playbook documentation</li></ul>",
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
					"title":       "Application & Platform Security",
					"subtitle":    "Defense",
					"description": "WAF / IDS / IPS / SIEM / IAM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Pentesting & Testing",
					"subtitle":    "Offense/Defense",
					"description": "Nmap / Burp / Metasploit / OWASP",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Process & Compliance",
					"subtitle":    "Governance",
					"description": "SDL / DevSecOps / Classified Protection / Compliance",
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
					"major":       "Information Security",
					"degree":      "Master's",
					"timeStart":   "2015-09",
					"timeEnd":     "2018-06",
					"description": "Application and platform security; exercises and governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Shanghai Jiao Tong University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "Computer fundamentals and network security; campus projects.",
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
					"description": `<ul><li>Security tools: <a href="https://sec-tools.zhwei.invalid">sec-tools.zhwei.invalid</a></li><li>Security handbook: <a href="https://sec-handbook.zhwei.invalid">sec-handbook.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateSecurityEngineerPresetEn() []byte {
	seByte, _ := json.Marshal(securityEngineerPresetJSONEn)
	return seByte
}
