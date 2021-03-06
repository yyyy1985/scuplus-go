package model

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// Course 课程
type Course struct {
	Model
	College         string           `json:"college"`                                                 // 学院
	CourseID        string           `json:"course_id"`                                               // 课程号
	Name            string           `json:"name"`                                                    // 课程名
	LessonID        string           `json:"lesson_id"`                                               // 课序号
	Credit          float64          `json:"credit"`                                                  // 学分
	ExamType        string           `json:"exam_type"`                                               // 考试类型
	AllWeek         string           `json:"all_week"`                                                // 周次: 1,2,3,4
	Day             int              `json:"day"`                                                     // 星期
	Session         string           `json:"session"`                                                 // 节次 1,2
	Campus          string           `json:"campus"`                                                  // 校区
	Building        string           `json:"building"`                                                // 教学楼
	Classroom       string           `json:"classroom"`                                               // 教室
	Max             int              `json:"max"`                                                     // 课容量
	StudentNo       int              `json:"student_no"`                                              // 学生数
	CourseLimit     string           `json:"course_limit"`                                            // 选课限制说明
	CourseCount     CourseCount      `json:"course_count"`                                            // 课程统计信息
	CourseEvaluates []CourseEvaluate `gorm:"many2many:course_and_evaluates;" json:"course_evaluates"` // 评价
	Teachers        []Teacher        `gorm:"many2many:course_teachers;" json:"teachers"`              // 教师
}

// AfterSave 回调
func (c *Course) AfterSave(scope *gorm.Scope) error {
	cc := CourseCount{}
	DB().FirstOrCreate(&cc, CourseCount{
		CourseID: c.CourseID,
		LessonID: c.LessonID,
	})
	cc.Name = c.Name
	cc.Day = c.Day
	cc.Credit = c.Credit
	cc.Campus = c.Campus
	// 教师统计
	teachers := []Teacher{}
	DB().Model(c).Related(&teachers, "Teachers")
	cc.Teacher = ""
	for _, teacher := range teachers {
		cc.Teacher = cc.Teacher + "," + teacher.Name
	}
	cc.Teacher = strings.Trim(cc.Teacher, ",")
	DB().Save(&cc)
	return nil
}
