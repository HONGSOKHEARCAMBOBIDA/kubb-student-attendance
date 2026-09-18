package route

const (

	// Authentication
	Login     = "login"
	LoginByQr = "loginbyqr"
	Refresh   = "refresh"
	Logout    = "logout"

	// Company
	AddCompany = "add.company"

	ViewCompany       = "view.company"
	ViewMajor         = "view.Major"
	ViewShift         = "view.Shift"
	ViewGeneration    = "view.Generation"
	ViewProgramme     = "view.Programme"
	EditCompany       = "edit.company/:id"
	EditTelegram      = "edit.telegram/:id"
	ViewManageCompany = "view.manage.company"
	ViewCompanyColor  = "view.company.color"
	ViewCompanyScan   = "view.company.scan"
	ToggleClassStatus = "toggle.status.class/:id"

	// User
	AddUser          = "add.user"
	AddUserClass     = "add.user.class"
	AddUserExcell    = "add.user.excell"
	ViewUser         = "view.user"
	EditUser         = "edit.user/:id"
	ToggleUserStatus = "toggle.status.user/:id"
	ChangePassword   = "change.password"
	CountUser        = "count.user"
	DeleteUser       = "delete.user/:id"
	ViewUserData     = "view.user.data"
	ViewUserApprove  = "view.user.approve"
	VerifyUser       = "verify.user/:id"

	// Shift
	EditShift = "edit.shift"
	AddShift  = "add.shift"

	// Attendance
	AddAttendance        = "add.attendance"
	ViewAttendance       = "view.attendance"
	ViewAttendanceReport = "view.attendance.report"
	ViewAttendanceDraft  = "view.attendance.draft"
	EditAttendance       = "edit.attendance"
	DeleteAttendance     = "delete.attendance/:id"

	// Payroll
	AddPayroll            = "add.payroll"
	ViewPayroll           = "view.payroll"
	ViewPayrollDraft      = "view.payroll.draft"
	EditPayroll           = "edit.payroll"
	DeletePayroll         = "delete.payroll"
	GenerateAttendancePDF = "generate.attendance.pdf"

	// Role
	ViewRole = "view.role"
	EditRole = "edit.role/:id"

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
	EditLeaveType = "edit.leave.type/:id"

	// LeaveRequest
	ViewLeaveRequest       = "view.leave.request"
	ViewNotPermisionLeave  = "view.not.permission.leave"
	AddLeaveRequest        = "add.leave.request"
	EditLeaveRequest       = "edit.leave.request/:id"
	EditStatusLeaveRequest = "edit.status.leave.request/:id"
	DeleteLeaveRequest     = "delete.leave.request/:id"
	ApproveLeave           = "approve.leave/:id"
	AddLeaveNotPermission  = "add.leave.not.permission"

	// RoleHasPermission
	ViewRoleHasPermission   = "view.role.has.permission/:id"
	AddRoleHasPermission    = "add.role.has.permission"
	DeleteRoleHasPermission = "delete.role.has.permission"
)
