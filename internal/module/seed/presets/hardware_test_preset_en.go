package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var hardwareTestPresetJSONEn = map[string]any{
	"title":    "Hardware Testing Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Hardware Testing Engineer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
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
					"description": "<ul><li>5 years of hardware testing; reliability and environmental tests; failure analysis and remediation</li><li>Electrical/signal/interface tests; EMC/ESD/safety tests and certification processes</li><li>Test jig automation, data acquisition and reporting</li><li>Collaborate with firmware/hardware to improve stability and yield</li></ul>",
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
					"title":       "Smart Hardware Company",
					"subtitle":    "Hardware Testing Engineer (Reliability/Certification)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned reliability and environmental tests (temperature/humidity/vibration)</li><li>Drove EMC/ESD/safety certifications; design optimizations and remediation</li><li>Built automated test jigs and data capture for efficiency and consistency</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Consumer Electronics Enterprise",
					"subtitle":    "Hardware Testing Engineer (Interfaces/Firmware)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Interface and signal integrity testing; issue localization and design changes</li><li>Firmware/hardware collaboration to optimize protocols and timing</li></ul>",
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
					"title":       "Automated Test Jig & Data Platform",
					"subtitle":    "Hardware Testing Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: LabVIEW / Python / NI instruments</li><li>Automated jigs and data platform; unified test process and reporting</li><li>Improved efficiency and consistency across product lines</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "EMC/ESD & Safety Certifications",
					"subtitle":    "Hardware Testing Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: EMC/ESD instruments / standard processes</li><li>Certification tests and remediation; shielding/grounding improvements</li><li>Better stability and yield</li></ul>",
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
					"title":       "Testing & Certification",
					"subtitle":    "Hardware",
					"description": "EMC / ESD / Safety / Environmental & Reliability / Signal Integrity",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Automation & Jigs",
					"subtitle":    "Engineering",
					"description": "LabVIEW / Python / NI Instruments / Data Acquisition & Reporting",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Collaboration & Optimization",
					"subtitle":    "Quality",
					"description": "Firmware/Hardware Collaboration / Design Changes / Yield Improvements",
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
					"title":       "Xidian University",
					"major":       "Electronic Information Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Hardware reliability and certification; test platforms and jig building.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Xidian University",
					"major":       "Electronic Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Circuits and signals; campus projects and competitions.",
				},
			},
		},
	},
}

func GenerateHardwareTestPresetEn() []byte {
	hwByte, _ := json.Marshal(hardwareTestPresetJSONEn)
	return hwByte
}
