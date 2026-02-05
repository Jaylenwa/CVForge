package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var sysAdminPresetJSONEn = map[string]any{
	"title":    "System/Network Administrator Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "System/Network Administrator",
		"City":       "Shanghai",
		"Money":      "20k-30k",
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
					"description": "<ul><li>5 years of system/network administration; Linux/Windows; TCP/IP, routing/switching</li><li>Security and access control: firewalls/WAF/VPN/ACL; permissions and auditing</li><li>Virtualization and storage: VMware/KVM, NAS/SAN, backup and DR</li><li>Automation and scripting: Shell/Python for efficiency and consistency</li></ul>",
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
					"subtitle":    "System/Network Administrator (Networking/Security)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Maintained enterprise networks and security devices; optimized routing and access control</li><li>Unified logs and audits; asset and permission management</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Manufacturing Company",
					"subtitle":    "System/Network Administrator (Compute/Storage)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Maintained servers/virtualization; storage/backup for production stability</li><li>Introduced automation scripts and processes for delivery efficiency</li></ul>",
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
					"title":       "Networking & Security Governance",
					"subtitle":    "System/Network Administrator",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Firewall / WAF / VPN</li><li>Routing and access control; unified audits and logs</li><li>Permissions and asset management</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Virtualization & Storage Modernization",
					"subtitle":    "System/Network Administrator",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: VMware / KVM / NAS / SAN</li><li>Backup and DR for reliability; automation scripts and processes</li></ul>",
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
					"title":       "Systems & Networking",
					"subtitle":    "Foundation",
					"description": "Linux / Windows / TCP/IP / Routing / Switching",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Security & Access",
					"subtitle":    "Governance",
					"description": "Firewall / WAF / VPN / ACL / Audit",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Virtualization & Storage",
					"subtitle":    "Platforms",
					"description": "VMware / KVM / NAS / SAN / Backup / DR",
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
					"major":       "Network Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Enterprise networking and security; platform building and governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Beijing University of Posts and Telecommunications",
					"major":       "Communication Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Networking and communications; campus projects and competitions.",
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
					"description": `<ul><li>Automation scripts: <a href="https://net-scripts.zhwei.invalid">net-scripts.zhwei.invalid</a></li><li>Ops handbook: <a href="https://ops-handbook.zhwei.invalid">ops-handbook.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateSysAdminPresetEn() []byte {
	saByte, _ := json.Marshal(sysAdminPresetJSONEn)
	return saByte
}
