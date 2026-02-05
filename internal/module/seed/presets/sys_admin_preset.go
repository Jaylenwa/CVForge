package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var sysAdminPresetJSON = map[string]any{
	"title":    "系统/网络管理员简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "系统/网络管理员",
		"City":       "上海",
		"Money":      "20k-30k",
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
					"description": "<ul><li>5 年系统/网络管理经验，熟悉 Linux/Windows、TCP/IP、路由与交换</li><li>掌握安全与访问控制：防火墙/WAF/VPN/ACL，权限与审计</li><li>熟悉虚拟化与存储：VMware/KVM、NAS/SAN、备份与容灾</li><li>具备自动化与脚本能力：Shell/Python，提升效率与一致性</li></ul>",
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
					"title":       "某互联网企业",
					"subtitle":    "系统/网络管理员（网络/安全）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>维护企业网络与安全设备，优化路由与访问控制</li><li>统一日志与审计，完善资产与权限管理</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某制造企业",
					"subtitle":    "系统/网络管理员（主机/存储）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>维护主机与虚拟化平台、存储与备份，保障生产稳定</li><li>推行自动化脚本与流程，提升交付效率</li></ul>",
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
					"title":       "网络与安全治理",
					"subtitle":    "系统/网络管理员",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Firewall / WAF / VPN</li><li>优化路由与访问控制、统一审计与日志</li><li>完善权限与资产管理</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "虚拟化与存储改造",
					"subtitle":    "系统/网络管理员",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：VMware / KVM / NAS / SAN</li><li>建设备份与容灾方案，提升稳定性与恢复能力</li><li>沉淀自动化脚本与流程</li></ul>",
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
					"title":       "系统与网络",
					"subtitle":    "基础",
					"description": "Linux / Windows / TCP/IP / Routing / Switching",
				},
				map[string]any{
					"id":          "s2",
					"title":       "安全与访问",
					"subtitle":    "治理",
					"description": "Firewall / WAF / VPN / ACL / Audit",
				},
				map[string]any{
					"id":          "s3",
					"title":       "虚拟化与存储",
					"subtitle":    "平台",
					"description": "VMware / KVM / NAS / SAN / Backup / DR",
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
					"major":       "网络工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为企业网络与安全；参与网络平台与治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "北京邮电大学",
					"major":       "通信工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修网络与通信；参与校内项目与竞赛。",
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
					"description": `<ul><li>自动化脚本集：<a href="https://net-scripts.zhwei.invalid">net-scripts.zhwei.invalid</a></li><li>运维手册：<a href="https://ops-handbook.zhwei.invalid">ops-handbook.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateSysAdminPreset() []byte {
	saByte, _ := json.Marshal(sysAdminPresetJSON)
	return saByte
}
