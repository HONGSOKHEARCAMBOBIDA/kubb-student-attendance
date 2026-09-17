package routes

import (
	"mysql/constant/permission"
	"mysql/constant/route"
	"mysql/controller"
	"mysql/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	// shiftcontroller := controller.NewShiftController()
	backupcontroller := controller.NewBackupController()
	//leavededucttypecontroller := controller.NewLeaveDeductTypeController()
	rolehaspermissioncontroller := controller.NewRoleHasPermissionController()
	companycontroller := controller.NewCompanyController()
	attendancecontroller := controller.NewAttendanceController()
	leavecontroller := controller.NewLeaveController()
	deductcontroller := controller.NewLeaveDeductTypeController()
	r.Static("/clientimage", "./public/clientimage")
	public := r.Group("/")
	public.Use(middleware.APIKeyAuth())
	{
		public.POST(route.Login, authcontroller.Login)
		public.POST(route.Refresh, authcontroller.Refresh)
	}
	auth := r.Group("/")
	auth.Use(middleware.APIKeyAuth())
	auth.Use(middleware.AuthMiddleware())
	{
		// Company
		auth.GET(route.ViewMajor, middleware.PermissionMiddleware(permission.ViewMajor), companycontroller.GetMajor)
		auth.GET(route.ViewShift, middleware.PermissionMiddleware(permission.ViewShift), companycontroller.GetShift)
		auth.GET(route.ViewGeneration, middleware.PermissionMiddleware(permission.ViewGeneration), companycontroller.GetGeneration)
		auth.GET(route.ViewProgramme, middleware.PermissionMiddleware(permission.ViewProgramme), companycontroller.GetProgramme)
		auth.POST(route.AddCompany, middleware.PermissionMiddleware(permission.AddCompany), companycontroller.CreateClass)
		auth.GET(route.ViewCompany, middleware.PermissionMiddleware(permission.ViewCompany), companycontroller.GetClass)
		auth.PUT(route.EditCompany, middleware.PermissionMiddleware(permission.EditCompany), companycontroller.UpdateClass)
		auth.PUT(route.ToggleClassStatus, middleware.PermissionMiddleware(permission.EditCompany), companycontroller.ChangeStatusClass)
		auth.GET(route.ViewCompanyScan, middleware.PermissionMiddleware(permission.ViewCompany), companycontroller.GetClassScan)
		// User
		auth.POST(route.AddUserClass, middleware.PermissionMiddleware(permission.AddUser), authcontroller.CreateUserClass)
		auth.POST(route.AddUser, middleware.PermissionMiddleware(permission.AddUser), authcontroller.Register)
		auth.POST(route.AddUserExcell, middleware.PermissionMiddleware(permission.AddUser), authcontroller.RegisterFromExcel)
		auth.PUT(route.ToggleUserStatus, middleware.PermissionMiddleware(permission.EditUser), authcontroller.ToggleUserStatus)
		auth.PUT(route.EditUser, middleware.PermissionMiddleware(permission.EditUser), authcontroller.UpdateUser)
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetRole)
		auth.GET(route.ViewUserData, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetUserData)
		//auth.DELETE(route.DeleteUser, middleware.PermissionMiddleware(permission.EditUser), authcontroller.DeleteUser)

		// Shift
		// auth.PUT(route.EditShift, middleware.PermissionMiddleware(permission.EditUser), shiftcontroller.UpdateShift)
		// auth.POST(route.AddShift, middleware.PermissionMiddleware(permission.EditUser), shiftcontroller.CreateShift)

		// Backup
		auth.POST(route.CreateBackup, middleware.PermissionMiddleware(permission.CreateBackup), backupcontroller.TriggerBackup)
		auth.GET(route.ViewBackup, middleware.PermissionMiddleware(permission.ViewBackup), backupcontroller.ListBackups)
		auth.GET(route.DownloadBackup, middleware.PermissionMiddleware(permission.DownloadBackup), backupcontroller.DownloadBackup)
		auth.DELETE(route.DeleteBackup, middleware.PermissionMiddleware(permission.DeleteBackup), backupcontroller.DeleteBackup)

		// LeaveDeductType
		//auth.GET(route.ViewLeaveDeductType, middleware.PermissionMiddleware(permission.ViewLeaveDeductType), leavededucttypecontroller.GetLeaveDeductType)

		// RoleHasPermission
		auth.GET(route.ViewRoleHasPermission, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.GetRolePermission)
		auth.POST(route.AddRoleHasPermission, middleware.PermissionMiddleware(permission.AddRoleHasPermission), rolehaspermissioncontroller.CreateRoleHasPermission)
		auth.DELETE(route.DeleteRoleHasPermission, middleware.PermissionMiddleware(permission.DeleteRoleHasPermission), rolehaspermissioncontroller.DeleteRoleHasPermission)
		auth.PUT(route.EditRole, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.UpdateRole)

		auth.GET(route.ViewAttendanceDraft, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendanceDraft)
		auth.POST(route.AddAttendance, middleware.PermissionMiddleware(permission.AddAttendance), attendancecontroller.CreateAttendance)
		auth.GET(route.ViewAttendance, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendancePDF)
		auth.GET(route.ViewAttendanceReport, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendanceReport)

		// Leave
		auth.POST(route.AddLeaveRequest, middleware.PermissionMiddleware(permission.AddLeaveRequest), leavecontroller.CreateLeaveRequest)
		auth.PUT(route.EditLeaveRequest, middleware.PermissionMiddleware(permission.EditLeaveRequest), leavecontroller.UpdateLeaveRequest)
		auth.GET(route.ViewLeaveRequest, middleware.PermissionMiddleware(permission.ViewLeaveRequest), leavecontroller.GetLeaveRequest)
		auth.PUT(route.ApproveLeave, middleware.PermissionMiddleware(permission.EditStatusLeaveRequest), leavecontroller.VerifyLeaveRequest)
		auth.GET(route.ViewLeaveDeductType, middleware.PermissionMiddleware(permission.ViewLeaveDeductType), deductcontroller.GetLeaveDeductType)
	}
}
