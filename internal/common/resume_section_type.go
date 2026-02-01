package common

import "strings"

type ResumeSectionType string

const (
	ResumeSectionTypePersonal       ResumeSectionType = "personal"
	ResumeSectionTypeSummary        ResumeSectionType = "summary"
	ResumeSectionTypeExperience     ResumeSectionType = "experience"
	ResumeSectionTypeEducation      ResumeSectionType = "education"
	ResumeSectionTypeSkills         ResumeSectionType = "skills"
	ResumeSectionTypeProjects       ResumeSectionType = "projects"
	ResumeSectionTypeInternships    ResumeSectionType = "internships"
	ResumeSectionTypePortfolio      ResumeSectionType = "portfolio"
	ResumeSectionTypeAwards         ResumeSectionType = "awards"
	ResumeSectionTypeInterests      ResumeSectionType = "interests"
	ResumeSectionTypeExam           ResumeSectionType = "exam"
	ResumeSectionTypeSelfEvaluation ResumeSectionType = "selfEvaluation"
	ResumeSectionTypeCustom         ResumeSectionType = "custom"
)

var resumeSectionTypeByKey = map[string]ResumeSectionType{
	"personal":       ResumeSectionTypePersonal,
	"summary":        ResumeSectionTypeSummary,
	"experience":     ResumeSectionTypeExperience,
	"education":      ResumeSectionTypeEducation,
	"skills":         ResumeSectionTypeSkills,
	"projects":       ResumeSectionTypeProjects,
	"internships":    ResumeSectionTypeInternships,
	"portfolio":      ResumeSectionTypePortfolio,
	"awards":         ResumeSectionTypeAwards,
	"interests":      ResumeSectionTypeInterests,
	"exam":           ResumeSectionTypeExam,
	"selfevaluation": ResumeSectionTypeSelfEvaluation,
	"custom":         ResumeSectionTypeCustom,
}

func IsValidResumeSectionType(v string) bool {
	_, ok := NormalizeResumeSectionType(v)
	return ok
}

func NormalizeResumeSectionType(v string) (ResumeSectionType, bool) {
	key := canonicalSectionTypeKey(v)
	if key == "" {
		return "", false
	}
	out, ok := resumeSectionTypeByKey[key]
	return out, ok
}

func canonicalSectionTypeKey(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.ToLower(v)
	var b strings.Builder
	b.Grow(len(v))
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
