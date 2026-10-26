import api from "./index";

// Auth
export const login = (data) => api.post("/login", data);
export const loginByQr = (data) => api.post("/loginbyqr", data);
export const refreshToken = (data) =>
  api.post("/refresh", {}, { withCredentials: true });

// Company
export const getclassnopagination = () => api.get("/view.class.no.pagination");
export const getMajor = () => api.get("/view.Major");
export const getShift = () => api.get("/view.Shift");
export const addShift = (data) => api.post("/add.shift", data);
export const editShift = (id, data) => api.put(`/edit.shift/${id}`, data);
export const getGeneration = () => api.get("/view.Generation");
export const addGeneration = (data) => api.post("/add.Generation", data);
export const editGeneration = (id, data) =>
  api.put(`/update.Generation/${id}`, data);
export const toggleGeneration = (id) => api.put(`/toggle.Generation/${id}`);
export const getProgramme = () => api.get("/view.Programme");
export const getFaculty = () => api.get("/view.Faculty");
// export const saveFaculty = (data) => api.post('/save.faculty',data)

export const addFaculty = (data) => api.post("/save.faculty", data);
export const editFaculty = (id, data) => api.put(`/save.faculty/${id}`, data);


export const getClass = (params) => api.get("/view.company", { params });
export const createClass = (data) => api.post("/add.company", data);
export const updateClass = (id, data) => api.put(`/edit.company/${id}`, data);
export const changeStatusClass = (id) => api.put(`/toggle.status.class/${id}`);
export const updateClassTelegram = (id, data) =>
  api.put(`/edit.telegram/${id}`, data);
export const viewmanagecompany = () => api.get(`/view.manage.company`);
export const viewcompanycolor = () => api.get(`/view.company.color`);
export const viewcompanyscan = () => api.get(`/view.company.scan`);
export const updateexamfinal = (id) => api.put(`/update.exam.final/${id}`);
export const createExamPlatform = (data) =>
  api.post("/create.exam.platform", data);
export const getExamPlatforms = (id) => api.get(`/get.exam.platform/${id}`);
export const uploadExamDocument = (file) => {
  const fd = new FormData();
  fd.append("file", file);
  return api.post("/exam-documents/upload", fd, {
    // use your axios instance name
    headers: { "Content-Type": "multipart/form-data" },
  });
};
export const deleteExamPlatform = (id) =>
  api.delete(`/delete.exam.platform/${id}`);
export const updateExamPlatform = (id, data) =>
  api.put(`/update.exam.platform/${id}`, data);
export const updateExamPlatformDetail = (id, data) =>
  api.put(`/update.exam.platform.detail/${id}`, data);
// User
export const getUsers = (params) => api.get("/view.user", { params });
export const registerUsers = (data) => api.post("/add.user", data);
export const registerUsersMain = (data) => api.post("/add.user.main", data);
export const registerUsersExcel = (formData) =>
  api.post("add.user.excell", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
export const updateUser = (id, data) => api.put(`/edit.user/${id}`, data);
export const toggleUserStatus = (id) => api.put(`/toggle.status.user/${id}`);
export const deleteuser = (id) => api.delete(`/delete.user/${id}`);
export const changePassword = (data) => api.put("/change.password", data);
export const countuser = () => api.get("/count.user");
export const getrole = () => api.get("/view.role");
export const getuserdata = () => api.get(`/view.user.data`);
export const getuserapprove = () => api.get(`/view.user.approve`);
export const verifyuser = (id) => api.put(`verify.user/${id}`);
export const logoutUser = () => api.post("/logout");
export const adduserclass = (data) => api.post("/add.user.class", data);
export const editUserClass = (id, data) =>
  api.put(`/update.user.class/${id}`, data);
export const getStaft = () => api.get("/view.staft");
// Shift
export const createShift = (data) => api.post("/add.shift", data);
export const updateShift = (data) => api.put("/edit.shift", data);

// Attendance
export const createAttendance = (data) => api.post("/add.attendance", data);
export const getAttendance = (params) =>
  api.get("/view.attendance", { params });
export const getAttendanceReport = (params) =>
  api.get("/view.attendance.report", { params });
export const getAttendanceDraft = (params) => api.get("/view.attendance.draft");
export const exportAttendancePDF = (params) =>
  api.get("/view.attendance", {
    params,
  });
export const deleteattendance = (id) => api.delete(`/delete.attendance/${id}`);
export const getAttendanceforedit = (params) =>
  api.get(`/view.attendance.for.edit`, { params });
export const updateAttendanceRecordStatus = (id, status) =>
  api.put(`/update.attendance/${id}`, status);
// Payroll
export const getPayrollDraft = (params) =>
  api.get("/view.payroll.draft", { params });
export const createPayroll = (data) => api.post("/add.payroll", data);
export const getPayroll = (params) => api.get("/view.payroll", { params });
export const deletePayroll = (data) => api.post("/delete.payroll", data);

// Backup
export const triggerBackup = () => api.post("add.backup");
export const listBackups = () => api.get("view.backup");
export const downloadBackup = (filename) =>
  api.get("view.download.backup", {
    params: { file: filename },
    responseType: "blob",
  });
export const deleteBackup = (filename) =>
  api.delete("delete.backup", {
    params: { file: filename },
  });

// LeaveDeductType
export const getLeaveDeductType = () => api.get("view.leave.deduct.type");

// LeaveType
export const getleavetype = () => api.get("view.leave.type");
export const createleavetype = (data) => api.post("/add.leave.type", data);
export const updateleavetype = (id, data) =>
  api.put(`/edit.leave.type/${id}`, data);

// LeaveRequest
export const exportNotPermissionLeave = (params, config) =>
  request.get("/leave-request/not-permission/export", { params, ...config });
export const getNotPermissionLeave = (params) =>
  api.get("/view.not.permission.leave", { params });
export const getLeaveRequest = (params) =>
  api.get("/view.leave.request", { params });
export const addLeaveRequest = (data) => api.post("/add.leave.request", data);
export const editLeaveRequest = (id, data) =>
  api.put(`/edit.leave.request/${id}`, data);
export const approveLeaveRequest = (id, data) =>
  api.put(`/approve.leave/${id}`);
export const deleteLeaveRequest = (id) =>
  api.delete(`/delete.leave.request/${id}`);
export const addNotPermissionLeave = (data) =>
  api.post("/add.leave.not.permission", data);
export const getcountleave = () => api.get("/count.leave");

// RoleHasPermission
export const getrolehaspermission = (id) =>
  api.get(`/view.role.has.permission/${id}`);
export const addrolehaspermission = (data) =>
  api.post("/add.role.has.permission", data);
export const deleterolehaspermission = (data) =>
  api.delete("/delete.role.has.permission", { data });
export const editrole = (id, data) => api.put(`/edit.role/${id}`, data);

// Subject
export const getSubject = (params) => api.get("/view.subject", { params });
export const addSubject = (data) => api.post("/add.subject", data);
export const editSubject = (id, data) => api.put(`/edit.subject/${id}`, data);
export const toggleSubject = (id) => api.put(`/toggle.subject/${id}`);

// Major
export const addMajor = (data) => api.post("add.Major", data);
export const editMajor = (id, data) => api.put(`edit.Major/${id}`, data);
export const viewMajorWithPagination = (params) =>
  api.get("view.major.pagination", { params });
export const toggleMajor = (id) => api.put(`Toggle.Major/${id}`);
export const getMajorSubjects = (id, params) =>
  api.post(`view.major.subject/${id}`, {}, { params });
export const addMajorSubject = (id, data) =>
  api.post(`add.major.subject/${id}`, data);
export const toggleMajorSubject = (id) => api.put(`Toggle.major.subject/${id}`);
export const removeMajorSubject = (id) =>
  api.delete(`Remove.major.subject/${id}`);

// class schedule

export const getClassSchedule = (classId) =>
  api.get(`/class/${classId}/schedule`);
export const getClassAvailableSubjects = (classId) =>
  api.get(`/class/${classId}/schedule/subjects`);
export const createClassSchedule = (classId, data) =>
  api.post(`/class/${classId}/schedule`, data);
export const toggleClassSchedule = (id) =>
  api.patch(`/class-schedule/${id}/toggle`);

// score
export const getGradecomponent = () => api.get("/view.grade.component");
export const addScore = (data) => api.post("/add.score", data);
export const getScore = (params) => api.get("/view.score", { params });
export const importScoreExcel = (formData) =>
  api.post("/add.score.from.excell", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
export const updateScore = (data, id) => api.put(`/edit.score/${id}`, data);
export const getScoreReport = (params) =>
  api.get("/view.score.report", { params });

export const getTranscript = (id) => api.get(`/view.transcript/${id}`);

export const addMajorPrice = (data) => api.post("/add.Major.Price", data);
export const getMajorPrice = (params) =>
  api.get("/view.Major.Price", { params });
export const updateMajorPrice = (id, data) =>
  api.put(`/update.Major.Price/${id}`, data);

export const getFeeschedule = () => api.get("/view.fee.schedule");
export const addFee = (data) => api.post("/add.fee", data);
export const getUserClasss = (id) => api.get(`/view.user.class/${id}`);
export const getSchoolarship = () => api.get(`/view.schoolarship`);
export const addFeeTransaction = (data) =>
  api.post("/add.fee.transaction", data);
export const printInvoice = (id) => api.get(`/print.invoice/${id}`);
export const deletefeetransaction = (id) =>
  api.delete(`/delete.fee.transaction/${id}`);
export const deletefee = (id) => api.delete(`/delete.fee/${id}`);

export const getincomecategory = () => api.get("/view.income.category");
export const addincome = (data) => api.post("/add.income", data);
export const getincome = (params) => api.get("/view.income", { params });
export const deleteincome = (id) => api.delete(`/delete.income/${id}`);

export const getprovince = () => api.get("/location/province");
export const getdistrict = (id) => api.get(`/location/district/${id}`);
export const getcommunce = (id) => api.get(`/location/commune/${id}`);
export const getvillage = (id) => api.get(`/location/village/${id}`);

export const getstudentcategory = () => api.get("/view.student.category");
export const addStudentProfile = (data) =>
  api.post("/add.student.profile", data);
export const getStudentProfile = (id) => api.get(`/view.student.profile/${id}`);
export const updateStudentProfile = (id, data) =>
  api.put(`/update.student.profile/${id}`, data);

export const addStudentFamily = (data) => api.post("/add.student.Family", data);
export const getStudentFamily = (id) => api.get(`/view.student.Family/${id}`);
export const updateStudentFamily = (id, data) =>
  api.put(`/update.student.Family/${id}`, data);

export const addStudentEducation = (data) =>
  api.post("/add.student.Education", data);
export const getStudentEducation = (id) =>
  api.get(`/view.student.Education/${id}`);
export const updateStudentEducation = (id, data) =>
  api.put(`/update.student.Education/${id}`, data);

export const getDocumentType = () => api.get("/view.document.type");

export const addStudentDocument = (data) =>
  api.post("/add.student.Document", data);
export const getStudentDocuments = (id) =>
  api.get(`/view.student.Document/${id}`);
export const updateStudentDocument = (id, data) =>
  api.put(`/update.student.Document/${id}`, data);

export const addTeacherProfile = (data) =>
  api.post("/add.Teacher.Profile", data);
export const getTeacherProfile = (id) => api.get(`/view.Teacher.Profile/${id}`);
export const updateTeacherProfile = (id, data) =>
  api.put(`/update.Teacher.Profile`, data);
export const updateTeacherRate = (id, data) =>
  api.put(`/update.Teacher.Rate/${id}`, data);

export const addStudentExam = (data) => api.post("/add.student.exam", data);
export const getStudentExams = (id) => api.get(`/view.student.exam/${id}`);
export const deleteStudentExam = (id) =>
  api.delete(`/delete.student.exam/${id}`);

export const getStudentStudyHistory = (id) =>
  api.get(`/view.student.study.history/${id}`);

export const addClassDetail = (data) => api.post("/add.class.detail", data);
export const viewClassDetail = (id) => api.get(`/view.class.detail/${id}`)
export const updateClassDetail = (id,data) => api.put(`/updaste.class.detail/${id}`,data)


export const getWeekSchedule  = (params) => api.get(`/view.report.weed.schedule`,{params})