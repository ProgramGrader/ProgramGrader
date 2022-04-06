package common

import "time"

type AssignmentConfig struct {
	Language                   string `mapstructure:"language" validate:"required"`
	GradeDocs                  bool
	NonCodeSubmissions         bool
	StudentTestsEnabled        bool
	NumberStudentTestsRequired int
	TeacherUnitTests           bool
}

// -----------------

type TeacherConfig struct {
	Teachers []map[string]Teacher `mapstructure:"Teachers"`
}

type Teacher struct {
	TeacherName  string `mapstructure:"Teacher_Name"`
	TeacherEmail string `mapstructure:"Teacher_Email"`
	Salutations  string `mapstructure:"Salutations"`
}

type FullTeacher struct {
	TeacherName  string
	TeacherEmail string
	Salutations  string
	TeacherID    string
}

// -----------------

type Classes struct {
	Courses map[string][]map[string]Course `mapstructure:"Classes"`
}

type Course struct {
	TeacherID     string   `mapstructure:"Teacher_Id"`
	CourseID      int      `mapstructure:"Course_Id"`
	ClassDueDates []string `mapstructure:"Class_Due_Dates"`
}

// -----------------

type AssignmentSubmissions struct {
	AssignmentName string
	DueDate        time.Time
	Teacher        FullTeacher
	CourseNumber   string
	CourseName     string
	SubmissionURL  string
	GradeDocs      bool
	GradeTests     bool
	Data           []UserAssignmentData
}

type UserAssignmentData struct {
	UserName    string
	FullName    string
	GitURL      string
	CommitTimes []string
	FailedTests TestDocs
	MissingDocs TestDocs
}

type TestDocs struct {
	Total     int
	TestNames []string
	AllPass   bool
	Completed bool
}

// -------

type course struct {
	Name        string
	id          string
	Assignments []Assignment
}

type Assignment struct {
	Name    string
	DueDate time.Time
}

// --------

type DueAssignment struct {
	CourseType        string
	CourseIDSemID     string
	CourseTeacher     FullTeacher
	AssignmentDueDate time.Time
	AssignmentName    string
	AssignmentNumber  int
	AssignmentConfig
}
