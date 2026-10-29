package routes

import (
	"mysql/config"
	"mysql/constant/permission"
	"mysql/constant/route"
	"mysql/constant/share"
	"mysql/controller"
	"mysql/middleware"
	"mysql/repository"
	"mysql/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	// shiftcontroller := controller.NewShiftController()
	backupcontroller := controller.NewBackupController()
	//leavededucttypecontroller := controller.NewLeaveDeductTypeController()
	rolehaspermissioncontroller := controller.NewRoleHasPermissionController()
	classcontroller := controller.NewCompanyController()
	attendancecontroller := controller.NewAttendanceController()
	leavecontroller := controller.NewLeaveController()
	deductcontroller := controller.NewLeaveDeductTypeController()
	subjectcontroller := controller.NewSubjectController()
	majorcontroller := controller.NewMajorController()
	classschedulecontroller := controller.NewClassScheduleController()
	scorecontroller := controller.NewScoreController()
	transcriptcontroller := controller.NewTranscriptController()
	reportcontroller := controller.NewReportController()
	//repo
	locationrepo := repository.NewLocationRepository(config.DB)
	feerepo := repository.NewFeeRepository(config.DB)
	incomerepo := repository.NewIncomeRepository(config.DB)
	//service
	locationservice := service.NewLocationService(locationrepo)
	feeservice := service.NewFeeService(feerepo)
	incomeservice := service.NewIncomeService(incomerepo)
	//controller
	locationcontroller := controller.NewLocationController(locationservice)
	feecontroller := controller.NewFeeController(feeservice)
	incomecontroller := controller.NewIncomeController(incomeservice)
	r.Static("/clientimage", "./public/clientimage")
	r.Static("/uploads", "./uploads")
	r.GET("/health", func(ctx *gin.Context) {
		if db, err := config.DB.DB(); err != nil || db.Ping() != nil {
			share.ResponseError(ctx, 500, "Down")
		}
		share.ResponseSuccess(ctx, 200, "UP")
	})
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
		// Class
		auth.GET(route.ViewClassNoPagination, middleware.PermissionMiddleware(permission.ViewClass), classcontroller.GetClassNoPagination)

		// Shift
		auth.GET(route.ViewShift, middleware.PermissionMiddleware(permission.ViewShift), classcontroller.GetShift)
		auth.POST(route.AddShift, middleware.PermissionMiddleware(permission.AddShift), classcontroller.CreateShift)
		auth.PUT(route.EditShift, middleware.PermissionMiddleware(permission.EditShift), classcontroller.UpdateShift)

		// Generation
		auth.GET(route.ViewGeneration, middleware.PermissionMiddleware(permission.ViewGeneration), classcontroller.GetGeneration)
		auth.POST(route.AddGeneration, middleware.PermissionMiddleware(permission.AddGeneration), classcontroller.CreateGeneration)
		auth.PUT(route.UpdateGeneration, middleware.PermissionMiddleware(permission.UpdateGeneration), classcontroller.UpdateGeneration)
		auth.PUT(route.ToggleGeneration, middleware.PermissionMiddleware(permission.UpdateGeneration), classcontroller.ToggleGeneration)

		// Programme
		auth.GET(route.ViewProgramme, middleware.PermissionMiddleware(permission.ViewProgramme), classcontroller.GetProgramme)

		// Faculty
		auth.GET(route.ViewFaculty, middleware.PermissionMiddleware(permission.ViewFaculty), classcontroller.GetFaculty)
		auth.POST(route.SavePOSTFaculty, middleware.PermissionMiddleware(permission.SaveFaculty), classcontroller.SaveFaculty)
		auth.PUT(route.SavePUTFaculty, middleware.PermissionMiddleware(permission.SaveFaculty), classcontroller.SaveFaculty)

		// Class
		auth.POST(route.AddCompany, middleware.PermissionMiddleware(permission.AddClass), classcontroller.CreateClass)
		auth.GET(route.ViewCompany, middleware.PermissionMiddleware(permission.ViewClass), classcontroller.GetClass)
		auth.PUT(route.EditCompany, middleware.PermissionMiddleware(permission.EditClass), classcontroller.UpdateClass)
		auth.GET(route.ViewCompanyScan, middleware.PermissionMiddleware(permission.ViewClass), classcontroller.GetClassScan)
		auth.PUT(route.ToggleClassStatus, middleware.PermissionMiddleware(permission.EditClass), classcontroller.ChangeStatusClass)
		auth.PUT(route.UpdateExamFinal, middleware.PermissionMiddleware(permission.EditClass), classcontroller.UpdateExamFinal)

		// Examplatform
		auth.POST(route.CreateExamPlatform, middleware.PermissionMiddleware(permission.CreateExamPlatform), classcontroller.CreateExamPlatform)
		auth.GET(route.GetExamPlatform, middleware.PermissionMiddleware(permission.GetExamPlatform), classcontroller.ListExamPlatforms)
		auth.PUT(route.UpdateExamPlatform, middleware.PermissionMiddleware(permission.UpdateExamPlatform), classcontroller.UpdateExamPlatform)
		auth.PUT(route.UpdateExamPlatformDetail, middleware.PermissionMiddleware(permission.UpdateExamPlatformDetail), classcontroller.UpdateExamPlatformDetail)
		auth.DELETE(route.DeleteExamPlatform, middleware.PermissionMiddleware(permission.DeleteExamPlatform), classcontroller.DeleteExamPlatform)

		// ClassDetail
		auth.POST(route.AddClassDetail, middleware.PermissionMiddleware(permission.AddClassDetail), classcontroller.CreateClassDetail)
		auth.GET(route.ViewClassDetail, middleware.PermissionMiddleware(permission.ViewClassDetail), classcontroller.GetClassDetail)
		auth.PUT(route.UpdateClassDetail, middleware.PermissionMiddleware(permission.UpdateClassDetail), classcontroller.UpdateClassDetail)

		// User Class
		auth.POST(route.AddUserClass, middleware.PermissionMiddleware(permission.AddUserClass), authcontroller.CreateUserClass)
		auth.PUT(route.UpdateUserClass, middleware.PermissionMiddleware(permission.UpdateUserClass), authcontroller.UpdateUserClass)
		auth.GET(route.ViewUserClass, middleware.PermissionMiddleware(permission.ViewUserClass), feecontroller.GetUserClass)

		// Role
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewRole), authcontroller.GetRole)
		auth.POST(route.AddRole, middleware.PermissionMiddleware(permission.AddRole), rolehaspermissioncontroller.CreateRole)

		// StudentExam
		auth.POST(route.AddStudentExam, middleware.PermissionMiddleware(permission.AddStudentExam), authcontroller.AddStudentExam)
		auth.GET(route.ViewStudentExam, middleware.PermissionMiddleware(permission.ViewStudentExam), authcontroller.GetStudentExam)
		auth.DELETE(route.DeleteStudentExam, middleware.PermissionMiddleware(permission.DeleteStudentExam), authcontroller.DeleteStudentExam)

		// User
		auth.POST(route.AddUser, middleware.PermissionMiddleware(permission.AddUser), authcontroller.Register)
		auth.POST(route.AddUserMain, middleware.PermissionMiddleware(permission.AddUser), authcontroller.RegisterMain)
		auth.POST(route.AddUserExcell, middleware.PermissionMiddleware(permission.AddUser), authcontroller.RegisterFromExcel)
		auth.PUT(route.ToggleUserStatus, middleware.PermissionMiddleware(permission.EditUser), authcontroller.ToggleUserStatus)
		auth.PUT(route.EditUser, middleware.PermissionMiddleware(permission.EditUser), authcontroller.UpdateUser)
		auth.GET(route.ViewUserData, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetUserData)
		auth.GET(route.ViewUser, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetUserNotStudent)
		//auth.DELETE(route.DeleteUser, middleware.PermissionMiddleware(permission.EditUser), authcontroller.DeleteUser)
		auth.GET(route.ViewStaft, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStaft)
		auth.GET(route.ViewStudentStudyHistory, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentStudyHistory)

		// Backup
		auth.POST(route.CreateBackup, middleware.PermissionMiddleware(permission.CreateBackup), backupcontroller.TriggerBackup)
		auth.GET(route.ViewBackup, middleware.PermissionMiddleware(permission.ViewBackup), backupcontroller.ListBackups)
		auth.GET(route.DownloadBackup, middleware.PermissionMiddleware(permission.DownloadBackup), backupcontroller.DownloadBackup)
		auth.DELETE(route.DeleteBackup, middleware.PermissionMiddleware(permission.DeleteBackup), backupcontroller.DeleteBackup)

		// RoleHasPermission
		auth.GET(route.ViewRoleHasPermission, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.GetRolePermission)
		auth.POST(route.AddRoleHasPermission, middleware.PermissionMiddleware(permission.AddRoleHasPermission), rolehaspermissioncontroller.CreateRoleHasPermission)
		auth.DELETE(route.DeleteRoleHasPermission, middleware.PermissionMiddleware(permission.DeleteRoleHasPermission), rolehaspermissioncontroller.DeleteRoleHasPermission)
		auth.PUT(route.EditRole, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.UpdateRole)

		// Attendance
		auth.GET(route.ViewAttendanceDraft, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendanceDraft)
		auth.POST(route.AddAttendance, middleware.PermissionMiddleware(permission.AddAttendance), attendancecontroller.CreateAttendance)
		auth.GET(route.ViewAttendance, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendancePDF)
		auth.GET(route.ViewAttendanceReport, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendanceReport)
		auth.GET(route.ViewAttendanceForEdit, middleware.PermissionMiddleware(permission.ViewAttendance), attendancecontroller.GetAttendance)
		auth.PUT(route.UpdateAttendance, middleware.PermissionMiddleware(permission.EditAttendance), attendancecontroller.UpdateAttendanceRecordStatus)

		// Leave
		auth.POST(route.AddLeaveRequest, middleware.PermissionMiddleware(permission.AddLeaveRequest), leavecontroller.CreateLeaveRequest)
		auth.PUT(route.EditLeaveRequest, middleware.PermissionMiddleware(permission.EditLeaveRequest), leavecontroller.UpdateLeaveRequest)
		auth.GET(route.ViewLeaveRequest, middleware.PermissionMiddleware(permission.ViewLeaveRequest), leavecontroller.GetLeaveRequest)
		auth.PUT(route.ApproveLeave, middleware.PermissionMiddleware(permission.EditStatusLeaveRequest), leavecontroller.VerifyLeaveRequest)
		auth.GET(route.ViewLeaveDeductType, middleware.PermissionMiddleware(permission.ViewLeaveDeductType), deductcontroller.GetLeaveDeductType)
		auth.DELETE(route.DeleteLeaveRequest, middleware.PermissionMiddleware(permission.DeleteLeaveRequest), leavecontroller.DeleteLeaveRequest)
		auth.PUT(route.RollbackLeave, middleware.PermissionMiddleware(permission.DeleteLeaveRequest), leavecontroller.RollBackLeaveRequest)
		auth.GET(route.ViewNotPermisionLeave, middleware.PermissionMiddleware(permission.ViewLeave), leavecontroller.GetNotPermissionLeave)
		auth.POST(route.AddLeaveNotPermission, middleware.PermissionMiddleware(permission.AddLeaveRequest), leavecontroller.AddNotPermission)
		auth.GET(route.CountLeave, middleware.PermissionMiddleware(permission.ViewLeave), leavecontroller.CountLeave)

		// Subject
		auth.POST(route.Addsubjec, middleware.PermissionMiddleware(permission.Addsubject), subjectcontroller.Create)
		auth.GET(route.Viewsubject, middleware.PermissionMiddleware(permission.Viewsubject), subjectcontroller.Get)
		auth.PUT(route.Editsubject, middleware.PermissionMiddleware(permission.Editsubject), subjectcontroller.Update)
		auth.PUT(route.ToggleSubject, middleware.PermissionMiddleware(permission.Editsubject), subjectcontroller.Toggle)
		auth.GET(route.GetClassAvailableSubjects, middleware.PermissionMiddleware(permission.Viewsubject), classschedulecontroller.GetAvailableSubjects)

		// Majoir
		auth.GET(route.ViewMajor, middleware.PermissionMiddleware(permission.ViewMajor), classcontroller.GetMajor)
		auth.POST(route.AddMajor, middleware.PermissionMiddleware(permission.AddMajor), majorcontroller.Create)
		auth.PUT(route.EditMajor, middleware.PermissionMiddleware(permission.EditMajor), majorcontroller.Update)
		auth.GET(route.ViewMajorWithPagination, middleware.PermissionMiddleware(permission.ViewMajor), majorcontroller.GetWithPagination)
		auth.PUT(route.ToggleMajor, middleware.PermissionMiddleware(permission.EditMajor), majorcontroller.Toggle)
		auth.POST(route.AddMajorSubject, middleware.PermissionMiddleware(permission.AddMajor), majorcontroller.AddSubject)
		auth.POST(route.GetMajorSubject, middleware.PermissionMiddleware(permission.ViewMajor), majorcontroller.GetSubjects)
		auth.PUT(route.ToggleMajorSubject, middleware.PermissionMiddleware(permission.Editsubject), majorcontroller.ToggleSubject)
		auth.DELETE(route.RemoveMajorSubject, middleware.PermissionMiddleware(permission.Editsubject), majorcontroller.RemoveSubject)

		// Class Schedule
		auth.GET(route.GetClassSchedule, middleware.PermissionMiddleware(permission.GetClassSchedule), classschedulecontroller.GetByClass)
		auth.POST(route.CreateClassSchedule, middleware.PermissionMiddleware(permission.CreateClassSchedule), classschedulecontroller.Create)
		auth.PATCH(route.ToggleClassSchedule, middleware.PermissionMiddleware(permission.ToggleClassSchedule), classschedulecontroller.Toggle)

		// grade component
		auth.GET(route.ViewGradeComponent, middleware.PermissionMiddleware(permission.ViewGradeComponent), scorecontroller.GetGradeComponent)

		// Transcript
		auth.GET(route.ViewTranscript, middleware.PermissionMiddleware(permission.ViewUser), transcriptcontroller.Transcript)

		// Score
		auth.POST(route.AddScore, middleware.PermissionMiddleware(permission.AddScore), scorecontroller.CreateScore)
		auth.POST(route.AddScoreFromExcel, middleware.PermissionMiddleware(permission.AddScore), scorecontroller.ImportScoreExcel)
		auth.GET(route.ViewScore, middleware.PermissionMiddleware(permission.ViewScore), scorecontroller.GetScore)
		auth.PUT(route.EditScore, middleware.PermissionMiddleware(permission.EditScore), scorecontroller.UpdateScore)
		auth.GET(route.ViewScoreReport, middleware.PermissionMiddleware(permission.ViewScore), scorecontroller.GetScoreReport)

		// Location
		auth.GET(route.ViewProvince, middleware.PermissionMiddleware(permission.ViewLocation), locationcontroller.GetProvince)
		auth.GET(route.ViewDistrict, middleware.PermissionMiddleware(permission.ViewLocation), locationcontroller.GetDistrict)
		auth.GET(route.ViewCommune, middleware.PermissionMiddleware(permission.ViewLocation), locationcontroller.GetCommune)
		auth.GET(route.ViewVillage, middleware.PermissionMiddleware(permission.ViewLocation), locationcontroller.GetVillage)

		// MajorPrice
		auth.POST(route.AddMajorPrice, middleware.PermissionMiddleware(permission.AddMajor), majorcontroller.AddMajorPrice)
		auth.GET(route.ViewMajorPrice, middleware.PermissionMiddleware(permission.ViewMajor), majorcontroller.GetMajorPrice)
		auth.PUT(route.UpdateMajorPrice, middleware.PermissionMiddleware(permission.AddMajor), majorcontroller.UpdateMajorPrice)

		// Fee
		auth.GET(route.ViewFeeSchedule, middleware.PermissionMiddleware(permission.ViewFeeSchedule), feecontroller.GetFeeSchedule)
		auth.POST(route.AddFee, middleware.PermissionMiddleware(permission.AddFee), feecontroller.AddFee)
		auth.GET(route.ViewSchoolarship, middleware.PermissionMiddleware(permission.ViewFeeSchedule), feecontroller.GetSchoolarship)
		auth.POST(route.AddFeeTransaction, middleware.PermissionMiddleware(permission.AddFee), feecontroller.AddFeeTransaction)
		auth.GET(route.PrintInvoice, middleware.PermissionMiddleware(permission.ViewFeeSchedule), feecontroller.PrintInvoice)
		auth.DELETE(route.DeleteFeeTransaction, middleware.PermissionMiddleware(permission.DeleteFeeTransaction), feecontroller.DeleteFeeTransaction)
		auth.DELETE(route.DeleteFee, middleware.PermissionMiddleware(permission.DeleteFee), feecontroller.DeleteFee)

		// Income
		auth.GET(route.ViewIncomeCategory, middleware.PermissionMiddleware(permission.ViewIncomeCategory), incomecontroller.GetIncomeCategory)
		auth.POST(route.AddIncome, middleware.PermissionMiddleware(permission.AddIncome), incomecontroller.AddIncome)
		auth.GET(route.ViewIncome, middleware.PermissionMiddleware(permission.ViewIncome), incomecontroller.GetIncome)
		auth.DELETE(route.DeleteIncome, middleware.PermissionMiddleware(permission.DeleteIncome), incomecontroller.DeleteIncome)

		// Student
		auth.GET(route.ViewStudentCategory, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentCategory)
		auth.POST(route.AddStudentProfile, middleware.PermissionMiddleware(permission.AddUser), authcontroller.AddStudentProfile)
		auth.GET(route.ViewStudentProfile, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentProfile)
		auth.PUT(route.UpdateStudentProfile, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateStudentProfile)
		auth.POST(route.AddStudentFamily, middleware.PermissionMiddleware(permission.AddUser), authcontroller.AddStudentFamily)
		auth.GET(route.ViewStudentFamily, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentFamily)
		auth.PUT(route.UpdateStudentFamily, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateStudentFamily)
		auth.POST(route.AddStudentEducation, middleware.PermissionMiddleware(permission.AddUser), authcontroller.AddStudentEducation)
		auth.GET(route.ViewStudentEducation, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentEducation)
		auth.PUT(route.UpdateStudentEducation, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateStudentEducation)
		auth.GET(route.ViewDocumnetType, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetDocumentType)
		auth.POST(route.AddStudentDocument, middleware.PermissionMiddleware(permission.AddUser), authcontroller.AddStudentDocument)
		auth.GET(route.ViewStudentDocument, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetStudentDocuments)
		auth.PUT(route.UpdateStudentDocument, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateStudentDocument)

		// Teacher
		auth.POST(route.AddTeacherProfile, middleware.PermissionMiddleware(permission.AddUser), authcontroller.AddTeacherProfile)
		auth.GET(route.ViewTeacherProfile, middleware.PermissionMiddleware(permission.ViewUser), authcontroller.GetTeacherProfile)
		auth.PUT(route.UpdateTeacherProfile, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateTeacherProfile)
		auth.PUT(route.UpdateTeacherRate, middleware.PermissionMiddleware(permission.AddUser), authcontroller.UpdateTeacherRate)

		// Report
		auth.GET(route.ViewReportWeekSchedule, middleware.PermissionMiddleware(permission.ViewReport), reportcontroller.GetWeekSchedule)
	}
}
