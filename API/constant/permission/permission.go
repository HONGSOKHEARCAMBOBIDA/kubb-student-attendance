package permission

const (
	// Company
	AddClass         = "add.class"
	ViewClass        = "view.class"
	EditClass        = "edit.class"
	ViewMajor        = "view.Major"
	ViewShift        = "view.Shift"
	ViewGeneration   = "view.Generation"
	AddGeneration    = "add.Generation"
	UpdateGeneration = "update.Generation"
	ViewProgramme    = "view.Programme"

	// Exam
	CreateExamPlatform       = "create.exam.platform"
	UpdateExamPlatform       = "update.exam.platform"
	UpdateExamPlatformDetail = "update.exam.platform.detail"
	DeleteExamPlatform       = "delete.exam.platform"
	GetExamPlatform          = "get.exam.platform"

	// Class Detail
	AddClassDetail    = "add.class.detail"
	ViewClassDetail   = "view.class.detail"
	UpdateClassDetail = "updaste.class.detail"

	// UserClass
	AddUserClass    = "add.user.class"
	UpdateUserClass = "update.user.class"
	ViewUserClass   = "view.user.class"

	// Student Exam
	AddStudentExam    = "add.student.exam"
	ViewStudentExam   = "view.student.exam"
	DeleteStudentExam = "delete.student.exam"

	// Class Schedule
	GetClassSchedule    = "get.class.schedule"
	CreateClassSchedule = "add.class.schedule"
	ToggleClassSchedule = "toggle.class.schedule"

	// Role
	ViewRole = "view.role"
	EditRole = "edit.role"
	AddRole  = "add.role"

	// Faculty
	ViewFaculty = "view.Faculty"
	SaveFaculty = "save.Faculty"
	// User
	AddUser        = "add.user"
	ViewUser       = "view.user"
	EditUser       = "edit.user"
	ChangePassword = "change.password"

	// Attendance
	AddAttendance  = "add.attendance"
	ViewAttendance = "view.attendance"
	EditAttendance = "edit.attendance"

	// Payroll
	AddPayroll  = "add.payroll"
	ViewPayroll = "view.payroll"
	EditPayroll = "edit.payroll"

	// BackUp
	CreateBackup   = "add.backup"
	ViewBackup     = "view.backup"
	DownloadBackup = "view.download.backup"
	DeleteBackup   = "delete.backup"

	// LeaveDeduction
	ViewLeaveDeductType = "view.leave.deduct.type"

	// LeaveType
	ViewLeave     = "view.leave.type"
	AddLeaveType  = "add.leave.type"
	EditLeaveType = "edit.leave.type"

	// LeaveRequest
	ViewLeaveRequest       = "view.leave.request"
	AddLeaveRequest        = "add.leave.request"
	EditLeaveRequest       = "edit.leave.request"
	EditStatusLeaveRequest = "edit.status.leave.request"
	DeleteLeaveRequest     = "delete.leave.request"

	// RoleHasPermission
	ViewRoleHasPermission   = "view.role.has.permission"
	AddRoleHasPermission    = "add.role.has.permission"
	DeleteRoleHasPermission = "delete.role.has.permission"

	// Subject
	Viewsubject = "view.subject"
	Addsubject  = "add.subject"
	Editsubject = "edit.subject"

	// Major

	AddMajor  = "add.Major"
	EditMajor = "edit.Major"

	EditShift = "edit.shift"
	AddShift  = "add.shift"

	ViewGradeComponent = "view.grade.component"
	AddScore           = "add.score"
	ViewScore          = "view.score"
	EditScore          = "edit.score"

	ViewLocation = "view.location"

	ViewFeeSchedule      = "view.fee.schedule"
	AddFee               = "add.fee"
	DeleteFeeTransaction = "delete.fee.transaction"
	DeleteFee            = "delete.fee"

	ViewIncomeCategory = "view.income.category"
	ViewIncome         = "view.income"
	AddIncome          = "add.income"
	DeleteIncome       = "delete.income"

	ViewReport = "view.report"
)
