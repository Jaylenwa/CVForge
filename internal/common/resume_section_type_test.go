package common

import "testing"

func TestNormalizeResumeSectionType(t *testing.T) {
	cases := []struct {
		in   string
		want ResumeSectionType
		ok   bool
	}{
		{in: "", ok: false},
		{in: "experience", want: ResumeSectionTypeExperience, ok: true},
		{in: "Experience", want: ResumeSectionTypeExperience, ok: true},
		{in: " self_evaluation ", want: ResumeSectionTypeSelfEvaluation, ok: true},
		{in: "self-evaluation", want: ResumeSectionTypeSelfEvaluation, ok: true},
		{in: "SELF_EVALUATION", want: ResumeSectionTypeSelfEvaluation, ok: true},
		{in: "unknown", ok: false},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, ok := NormalizeResumeSectionType(tc.in)
			if ok != tc.ok {
				t.Fatalf("ok=%v want=%v got=%q", ok, tc.ok, got)
			}
			if tc.ok && got != tc.want {
				t.Fatalf("want=%q got=%q", tc.want, got)
			}
		})
	}
}
