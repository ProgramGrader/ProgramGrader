package configProcessors

import (
	"SubmissionProducer/internal/common"
	"strings"
)

func GetTeachers() map[string]common.FullTeacher {
	TeachersRaw := GetTeachersFromConfig(common.ConfigPath())

	var teacher = make(map[string]common.FullTeacher)

	for _, t := range TeachersRaw.Teachers {
		for index, value := range t {
			teacher[strings.ToLower(index)] =
				common.FullTeacher{
					TeacherID:    strings.ToLower(index),
					TeacherName:  value.TeacherName,
					TeacherEmail: value.TeacherEmail,
					Salutations:  value.TeacherEmail,
				}
		}
	}

	return teacher
}
