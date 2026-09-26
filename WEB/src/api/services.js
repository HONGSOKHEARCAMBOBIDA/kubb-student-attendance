import api from './index'

// Auth
export const login = (data) => api.post('/login', data)
export const loginByQr = (data) => api.post('/loginbyqr', data)
export const refreshToken = (data) => api.post('/refresh',{}, { withCredentials: true })

// Company
export const getMajor = () => api.get('/view.Major')
export const getShift = () => api.get('/view.Shift')
export const addShift = (data) => api.post('/add.shift',data)
export const editShift = (id,data) => api.put(`/edit.shift/${id}`,data)
export const getGeneration = () => api.get('/view.Generation')
export const addGeneration = (data) => api.post('/add.Generation',data)
export const editGeneration = (id,data) => api.put(`/update.Generation/${id}`,data)
export const toggleGeneration = (id) => api.put(`/toggle.Generation/${id}`)
export const getProgramme = () => api.get('/view.Programme')
export const getClass = (params) => api.get('/view.company', { params })
export const createClass = (data) => api.post('/add.company', data)
export const updateClass = (id, data) => api.put(`/edit.company/${id}`, data)
export const changeStatusClass = (id) => api.put(`/toggle.status.class/${id}`)
export const updateClassTelegram = (id,data) => api.put(`/edit.telegram/${id}`,data)
export const viewmanagecompany = () => api.get(`/view.manage.company`)
export const viewcompanycolor = () => api.get(`/view.company.color`)
export const viewcompanyscan = () => api.get(`/view.company.scan`)

// User
export const getUsers = (params) => api.get('/view.user', { params })
export const registerUsers = (data) => api.post('/add.user', data)
export const registerUsersMain = (data) => api.post('/add.user.main', data)
export const registerUsersExcel = (formData) =>
  api.post("add.user.excell", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
export const updateUser = (id, data) => api.put(`/edit.user/${id}`, data)
export const toggleUserStatus = (id) => api.put(`/toggle.status.user/${id}`)
export const deleteuser = (id) => api.delete(`/delete.user/${id}`)
export const changePassword = (data) => api.put('/change.password', data)
export const countuser = () => api.get('/count.user')
export const getrole = () => api.get('/view.role')
export const getuserdata = () => api.get(`/view.user.data`)
export const getuserapprove = () => api.get(`/view.user.approve`)
export const verifyuser = (id) => api.put(`verify.user/${id}`)
export const logoutUser = () => api.post('/logout')
export const adduserclass = (data) => api.post('/add.user.class',data)
export const editUserClass = (id,data) => api.put(`/update.user.class/${id}`,data)


// Shift
export const createShift = (data) => api.post('/add.shift', data)
export const updateShift = (data) => api.put('/edit.shift', data)

// Attendance
export const createAttendance = (data) => api.post('/add.attendance', data)
export const getAttendance = (params) => api.get('/view.attendance', { params })
export const getAttendanceReport  = (params) => api.get('/view.attendance.report', { params })
export const getAttendanceDraft = (params) => api.get('/view.attendance.draft')
export const exportAttendancePDF = (params) => api.get('/view.attendance', {
  params,
})
export const deleteattendance = (id) => api.delete(`/delete.attendance/${id}`)
export const getAttendanceforedit = (params) => api.get(`/view.attendance.for.edit`,{params}) 
export const updateAttendanceRecordStatus = (id,status) => api.put(`/update.attendance/${id}`,status)
// Payroll
export const getPayrollDraft = (params) => api.get('/view.payroll.draft', { params })
export const createPayroll = (data) => api.post('/add.payroll', data)
export const getPayroll = (params) => api.get('/view.payroll',{params})
export const deletePayroll = (data) =>
  api.post('/delete.payroll',data)

// Backup
export const triggerBackup = () => api.post('add.backup')
export const listBackups = () => api.get('view.backup')
export const downloadBackup = (filename) => api.get('view.download.backup', { 
  params: { file: filename },
  responseType: 'blob' 
})
export const deleteBackup = (filename) => api.delete('delete.backup',{
  params: {file: filename},
})

// LeaveDeductType
export const getLeaveDeductType = () => api.get('view.leave.deduct.type')

// LeaveType
export const getleavetype = () => api.get('view.leave.type')
export const createleavetype = (data) => api.post('/add.leave.type',data)
export const updateleavetype = (id, data) => api.put(`/edit.leave.type/${id}`, data)

// LeaveRequest
export const exportNotPermissionLeave = (params, config) =>
  request.get("/leave-request/not-permission/export", { params, ...config });
export const getNotPermissionLeave = (params) => api.get('/view.not.permission.leave',{params})
export const getLeaveRequest = (params) => api.get('/view.leave.request',{params})
export const addLeaveRequest = (data) => api.post('/add.leave.request',data)
export const editLeaveRequest = (id,data) => api.put(`/edit.leave.request/${id}`,data)
export const approveLeaveRequest = (id,data) => api.put(`/approve.leave/${id}`)
export const deleteLeaveRequest = (id) => api.delete(`/delete.leave.request/${id}`)
export const addNotPermissionLeave = (data) => api.post('/add.leave.not.permission',data)

// RoleHasPermission
export const getrolehaspermission = (id) => api.get(`/view.role.has.permission/${id}`)
export const addrolehaspermission = (data) => api.post('/add.role.has.permission',data)
export const deleterolehaspermission = (data) => api.delete('/delete.role.has.permission',{data})
export const editrole = (id,data) => api.put(`/edit.role/${id}`,data)

// Subject
export const getSubject = (params) => api.get('/view.subject',{params})
export const addSubject = (data) => api.post('/add.subject',data)
export const editSubject = (id,data) => api.put(`/edit.subject/${id}`,data)
export const toggleSubject = (id) => api.put(`/toggle.subject/${id}`,)

// Major
export const addMajor = (data) => api.post('add.Major',data)
export const editMajor = (id,data) => api.put(`edit.Major/${id}`,data)
export const viewMajorWithPagination = (params) => api.get('view.major.pagination',{params})
export const toggleMajor = (id) => api.put(`Toggle.Major/${id}`,)
export const getMajorSubjects = (id, params) =>
  api.post(`view.major.subject/${id}`, {}, { params });
export const addMajorSubject = (id,data) => api.post(`add.major.subject/${id}`,data)
export const toggleMajorSubject = (id) => api.put(`Toggle.major.subject/${id}`)
export const removeMajorSubject = (id) => api.delete(`Remove.major.subject/${id}`)

// class schedule

export const getClassSchedule = (classId) => api.get(`/class/${classId}/schedule`);
export const getClassAvailableSubjects = (classId) =>api.get(`/class/${classId}/schedule/subjects`);
export const createClassSchedule = (classId, data) =>api.post(`/class/${classId}/schedule`, data);
export const toggleClassSchedule = (id) => api.patch(`/class-schedule/${id}/toggle`);

// score
export const getGradecomponent = () =>  api.get('/view.grade.component')
export const addScore = (data) =>  api.post('/add.score',data)
export const getScore = (params) =>  api.get('/view.score',{params})
export const importScoreExcel = (formData) =>
  api.post("/add.score.from.excell", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });