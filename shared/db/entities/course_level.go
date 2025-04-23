package entities

import "slices"

type CourseLevel string

const (
	CourseLevel_Beginner        = "beginner"
	CourseLevel_PreIntermediate = "pre-intermediate"
	CourseLevel_Intermediate    = "intermediate"
	CourseLevel_Advance         = "advance"
)

func (courseLevel CourseLevel) IsValid(canBeEmpty bool) bool {
	if canBeEmpty && courseLevel == "" {
		return true
	}
	courseLevels := []CourseLevel{
		CourseLevel_Beginner,
		CourseLevel_PreIntermediate,
		CourseLevel_Intermediate,
		CourseLevel_Advance,
	}
	return slices.Contains(courseLevels, courseLevel)
}
