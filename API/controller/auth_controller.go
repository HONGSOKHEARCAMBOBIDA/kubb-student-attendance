package controller

import (
	"context"
	"errors"
	"log"
	"mysql/constant/apperror"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/pattern/res"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController() AuthController {
	return AuthController{
		service: service.NewAuthService(),
	}
}

func (cr *AuthController) Login(c *gin.Context) {
	var input request.AuthRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := cr.service.Login(input, c)
	if err != nil {
		log.Printf("Login Error: %v", err)
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) Refresh(c *gin.Context) {

	// var input request.RefreshTokenRequest
	///	log.Printf(input.RefreshToken)
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := cr.service.RefreshToken(cookie, c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) Register(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	var input request.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Register(c, input, c, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) RegisterMain(c *gin.Context) {
	var input request.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.RegisterMain(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) CreateUserClass(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	var input request.UserClass
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateUserClass(c, input, c, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (ctrl *AuthController) RegisterFromExcel(c *gin.Context) {
	classIDStr := c.PostForm("class_id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil || classID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id is required"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	xf, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid excel file"})
		return
	}
	defer xf.Close()

	sheet := xf.GetSheetName(0)
	rows, err := xf.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no data rows found"})
		return
	}

	// Map header -> column index so column order in the file doesn't matter.
	header := rows[0]
	col := map[string]int{}
	for i, h := range header {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	required := []string{"name_kh", "name_en", "gender", "code"}
	for _, r := range required {
		if _, ok := col[r]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing column: " + r})
			return
		}
	}

	get := func(row []string, key string) string {
		i := col[key]
		if i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}
	parseGender := func(v string) (int, bool) {
		switch strings.ToLower(v) {
		case "1", "ប្រុស", "m", "male":
			return 1, true
		case "2", "ស្រី", "f", "female":
			return 2, true
		}
		return 0, false
	}

	var input request.RegisterRequest
	input.ClassID = classID
	var skipped []int
	for idx, row := range rows[1:] {
		nameKH := get(row, "name_kh")
		nameEN := get(row, "name_en")
		code := get(row, "code")
		genderRaw := get(row, "gender")

		if nameKH == "" || nameEN == "" || code == "" {
			skipped = append(skipped, idx+2) // +2: header row + 1-index
			continue
		}
		gender, ok := parseGender(genderRaw)
		if !ok {
			skipped = append(skipped, idx+2)
			continue
		}

		input.UserInput = append(input.UserInput, request.UserInput{
			NameKH: nameKH,
			NameEN: nameEN,
			Gender: gender,
			Code:   code,
		})
	}

	if len(input.UserInput) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid rows to import"})
		return
	}

	userID, _ := strconv.Atoi(strconv.Itoa(c.GetInt("user_id"))) // adjust to however you extract the authed user id elsewhere
	if err := ctrl.service.Register(c.Request.Context(), input, c, userID); err != nil {
		if ae, ok := err.(*apperror.AppError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ae.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "users imported",
		"created": len(input.UserInput),
		"skipped": skipped,
	})
}

func (cr *AuthController) ToggleUserStatus(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.ToggleUserStatus(c, id, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "status changed")
}

func (cr *AuthController) UpdateUser(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.UserRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateUser(c, input, id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "updated user")
}

func (cr *AuthController) GetRole(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	data, err := cr.service.GetRole(c, userID)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) GetUserData(c *gin.Context) {
	userlog, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	data, err := cr.service.GetUserData(c, userlog)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) UpdateUserClass(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	var input request.UserClassUpdateStatus
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateUserClass(c, id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "company Updated")
}

func (cr *AuthController) GetUserNotStudent(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":    c.Query("name"),
		"role_id": c.Query("role_id"),
	}

	data, meta, err := cr.service.GetUserIncludeStudent(c.Request.Context(), userID, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, 200, data, meta)
}

func (cr *AuthController) GetStudentCategory(c *gin.Context) {
	data, err := cr.service.GetStudentCategory(c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) AddStudentProfile(c *gin.Context) {
	var input request.StudentProfileRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddStudentProfile(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetStudentProfile(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetStudentProfile(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
func (cr *AuthController) UpdateStudentProfile(c *gin.Context) {
	var input request.StudentProfileRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateStudentProfile(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) AddStudentFamily(c *gin.Context) {
	var input request.StudentFamily
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddStudentFamily(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) UpdateStudentFamily(c *gin.Context) {
	var input request.StudentFamily
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateStudentFamily(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetStudentFamily(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetStudentFamily(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) AddStudentEducation(c *gin.Context) {
	var input request.StudentEducation
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddStudentEducation(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetStudentEducation(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetStudentEducation(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) UpdateStudentEducation(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	var input request.StudentEducation
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateStudentEducation(c, id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetDocumentType(c *gin.Context) {
	data, err := cr.service.GetDocumentType(c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) AddStudentDocument(c *gin.Context) {
	var input request.StudentDocumentRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddStudentDocument(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) UpdateStudentDocument(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	var input request.StudentDocumentRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateStudentDocument(c, id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetStudentDocuments(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetStudentDocuments(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) AddTeacherProfile(c *gin.Context) {
	var input request.TeacherProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddTeacherProfile(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) GetTeacherProfile(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetTeacherProfile(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusNotFound, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) UpdateTeacherProfile(c *gin.Context) {
	var input request.TeacherProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateTeacherProfile(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "user create")
}

func (cr *AuthController) UpdateTeacherRate(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	var input request.TeacherRate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateTeacherRate(c, id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Udpated")
}

func (cr *AuthController) GetStaft(c *gin.Context) {
	data, err := cr.service.GetStaft(c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) AddStudentExam(c *gin.Context) {
	var input request.StudentExam
	if err := c.ShouldBindJSON(&input); err != nil {
		res.Fail(c, 400, "BAD_REQUEST", err.Error())
		return
	}
	if err := cr.service.AddStudentExam(c, input); err != nil {
		res.Fail(c, 500, "", err.Error())
		return
	}
	res.OK(c, 200)
}

func (cr *AuthController) GetStudentExam(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetStudentExam(c, id)
	if err != nil {
		res.Fail(c, http.StatusNoContent, "NOT_FOUND", err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) DeleteStudentExam(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	if err := cr.service.DeleteStudentExam(c, id); err != nil {
		res.Fail(c, http.StatusNoContent, "CAN NOT DELETE", err.Error())
		return
	}
	res.OK(c, 200)
}
