package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"SubmissionProducer/internal/configProcessors"
	"strconv"
	"strings"
	"time"
)

type classReturn struct {
	teacher  common.FullTeacher
	dueDates []time.Time
}

func GetAssignmentsAndDueDates() map[string]map[string]classReturn {

	// Set up the config repository
	configProcessors.GetOrUpdateConfigs(common.ConfigPath())

	Teachers := configProcessors.GetTeachers()
	//fmt.Println(Teachers)
	// Get Classes for this semester
	semesterClassesRaw := configProcessors.GetClassesFromConfig(common.ConfigPath()+"\\Courses\\", common.GetCurrentSemesterIdentifier())

	// class, course, teacher, and due dates
	var semesterClasses = make(map[string]map[string]classReturn)

	for CourseNumber, Course := range semesterClassesRaw.Courses {
		//fmt.Println(CourseNumber)
		for _, CourseValue := range Course {
			//fmt.Println(CourseValue)
			for _, class := range CourseValue {

				//fmt.Println()
				if _, ok := semesterClasses[CourseNumber]; !ok {
					semesterClasses[CourseNumber] = make(map[string]classReturn)
				}

				//Convert strings to times
				var DueDates []time.Time

				for _, stringTime := range class.ClassDueDates {
					tempTime, _ := time.Parse("01-02", stringTime)

					DueDates = append(DueDates, time.Date(time.Now().Year(), tempTime.Month(), tempTime.Day(), 11, 59, 59, 0, tempTime.Location()))
				}

				if teacher, ok := Teachers[strings.ToLower(class.TeacherID)]; ok {
					temp := classReturn{teacher, DueDates}

					semesterClasses[CourseNumber][strconv.Itoa(class.CourseID)] = temp
				} else {
					common.Error("Teacher was not found.")
				}
			}
		}
	}

	return semesterClasses

}
