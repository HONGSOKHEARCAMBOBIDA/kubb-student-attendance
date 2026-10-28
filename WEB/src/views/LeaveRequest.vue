<script setup>
import { ref, reactive, onMounted, computed, watch, onUnmounted } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import {
  getLeaveRequest,
  addLeaveRequest,
  editLeaveRequest,
  deleteLeaveRequest,
  approveLeaveRequest,
  getLeaveDeductType,
  getClass,
  viewcompanyscan,
  getClassAvailableSubjects,
  getProgramme,
  getGeneration,
  getMajor,
  getShift,
  rollbackleave
} from "../api/services.js";
import { useUserDataStore } from "../stores/user_data.js";
import AppSelect from "../../components/AppSelect.vue";
import AppDatePicker from "../../components/AppDatePicker.vue";
import { ElMessageBox } from "element-plus";

let searchTimer = null;
const notify = useNotification();
const userDataStore = useUserDataStore();

const leaveRequests = ref([]);
const deductTypes = ref([]);
const classes = ref([]);
const loading = ref(false);
const submitting = ref(false);
const approvingId = ref(null);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const subjectOptions = ref([]);
const subjectsLoading = ref(false)

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
});

// state filter
const programmes = ref([]);
const generations = ref([]);
const shifts = ref([]);
const majors = ref([]);

async function confirmDelete(message) {
  try {
    await ElMessageBox.confirm(message, "បញ្ជាក់", { type: "warning" });
    return true;
  } catch {
    return false; // the user clicked cancel
  }
}

async function handleRollback(row) {
  const confirmed = await confirmDelete("តើអ្នកពិតជាចង់ត្រឡប់ការអនុម័តច្បាប់នេះមែនទេ?");
  if (!confirmed) return;

  try {
    await rollbackleave(row.id);
    notify.success("ជោគជ័យ");
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចបានទេ"));
  }
}


const years = [
  { label: "ឆ្នាំទី 1", value: 1 },
  { label: "ឆ្នាំទី 2", value: 2 },
  { label: "ឆ្នាំទី 3", value: 3 },
  { label: "ឆ្នាំទី 4", value: 4 },
];

const semesters = [
  { label: "ឆមាសទី 1", value: 1 },
  { label: "ឆមាសទី 2", value: 2 },
];



// Matches the only filter the backend actually supports: name, class_id, status
const filter = reactive({
  name: "",
  class_id: "",
  status: "",
  subject_id: null,
  programme_id: null,
  generation_id: null,
  major_id: null,
  shift_id: null,
  year: null,
  semester: null,
  start_date: "",
  end_date: ""
});

const filterClass = reactive({
  major_id: "",
  shift_id: "",
  generation_id: "",
  programme_id: "",
  year: "",
  semester: "",
});

async function fetchProgramme() {
  try {
    const res = await getProgramme();
    programmes.value = (res.data.data || []).map((s) => ({
      label: `${s.name}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load programme");
  }
}

async function fetchGeneration() {
  try {
    const res = await getGeneration();
    generations.value = (res.data.data || []).map((s) => ({
      label: `${s.name_kh}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load programme");
  }
}

async function fetchMajor() {
  try {
    const res = await getMajor();
    majors.value = (res.data.data || []).map((s) => ({
      label: `${s.name_kh}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load programme");
  }
}

async function fetchShift() {
  try {
    const res = await getShift();
    shifts.value = (res.data.data || []).map((s) => ({
      label: `${s.name}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load programme");
  }
}

async function onProgrammchange(programmID) {
  filterClass.programme_id = programmID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}

async function onGenerationchange(generationID) {
  filterClass.generation_id = generationID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}

async function onMajorchange(majorID) {
  filterClass.major_id = majorID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}

async function onShiftchange(shiftID) {
  filterClass.shift_id = shiftID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}

async function onYearchange(yearID) {
  filterClass.year = yearID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}

async function onSemesterchange(semesterID) {
  filterClass.semester = semesterID;
  filter.class_id = null;
  filter.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchLeaveRequest();
}



// Model enum is only PENDING / APPROVE — no REJECTED/CANCELLED on the backend today
const statusOptions = [
  { value: "PENDING", label: "កំពុងរង់ចាំ" },
  { value: "APPROVE", label: "អនុម័តរួច" },
];

async function fetchSubjectOptions(classID) {
  subjectsLoading.value = true;
  try {
    const res = await getClassAvailableSubjects(classID);
    subjectOptions.value = (res.data.data || []).map((s) => ({
      label: `${s.code} — ${s.name_kh}`,
      value: s.id,
    }));
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to load subjects");
  } finally {
    subjectsLoading.value = false;
  }
}

async function onclasschange(classID) {
  filter.subject_id = null
  subjectOptions.value = []
  if (classID) {
    await fetchSubjectOptions(classID)
  }
  fetchLeaveRequest();
}

const defaultForm = () => ({
  class_id: null,
  start_date: "",
  end_date: "",
  back_to_work_date: "",
  total_day: null,
  deduct_type_id: null,
  reason: "",
});
const form = reactive(defaultForm());

// Field names/required-ness mirror request.LeaveRequestCreate / LeaveRequestUpdate
const rules = {
  class_id: [{ required: true, message: "សូមជ្រើសរើសថ្នាក់", trigger: "change" }],
  start_date: [{ required: true, message: "សូមជ្រើសរើសថ្ងៃចាប់ផ្តើម", trigger: "change" }],
  end_date: [{ required: true, message: "សូមជ្រើសរើសថ្ងៃបញ្ចប់", trigger: "change" }],
  back_to_work_date: [{ required: true, message: "សូមជ្រើសរើសថ្ងៃចូលធ្វើការវិញ", trigger: "change" }],
  total_day: [{ required: true, message: "សូមបញ្ចូលចំនួនថ្ងៃ", trigger: "blur" }],
  deduct_type_id: [{ required: true, message: "សូមជ្រើសរើសឯកតាកាត់ថ្ងៃ", trigger: "change" }],
  reason: [{ required: true, message: "សូមបញ្ចូលមូលហេតុ", trigger: "blur" }],
};

const canCreateLeave = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.leave.request"),
);
const canEditLeave = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.leave.request"),
);
const canApproveLeave = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.status.leave.request"),
);
const canDeleteLeave = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "delete.leave.request"),
);

async function fetchClasses() {
  try {
    const params = {};
    if (filterClass.programme_id)
      params.programme_id = filterClass.programme_id;
    if (filterClass.major_id) params.major_id = filterClass.major_id;
    if (filterClass.shift_id) params.shift_id = filterClass.shift_id;
    if (filterClass.generation_id)
      params.generation_id = filterClass.generation_id;
    if (filterClass.year) params.year = filterClass.year;
    if (filterClass.semester) params.semester = filterClass.semester;
    const res = await getClass(params);
    classes.value = (res.data.data || []).map((s) => ({
      label: `${s.programme_name} - ${s.major_name} - ${s.shift_name} - ឆ្នាំទី${s.year} - ឆមាសទី${s.semester}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load classes");
  }
}

async function fetchDeductTypes() {
  try {
    const res = await getLeaveDeductType();
    deductTypes.value = res.data.data || [];
  } catch {
    notify.error("Failed to load deduct types");
  }
}

async function fetchLeaveRequest() {
  loading.value = true;
  try {
    const params = { page: pagination.page, page_size: pagination.page_size };
    if (filter.name) params.name = filter.name;
    if (filter.status) params.status = filter.status;
    if (filter.class_id) params.class_id = filter.class_id;
    if (filter.subject_id) params.subject_id = filter.subject_id;
    if (filter.programme_id) params.programme_id = filter.programme_id;
    if (filter.generation_id) params.generation_id = filter.generation_id;
    if (filter.major_id) params.major_id = filter.major_id;
    if (filter.shift_id) params.shift_id = filter.shift_id;
    if (filter.year) params.year = filter.year;
    if (filter.semester) params.semester = filter.semester;
    if (filter.start_date) params.start_date = filter.start_date;
    if (filter.end_date) params.end_date = filter.end_date;
    const res = await getLeaveRequest(params);
    leaveRequests.value = res.data.data || [];
    pagination.total = res.data.pagination?.totalCount || 0;
  } catch {
    notify.error("Failed to load leave requests");
  } finally {
    loading.value = false;
  }
}

function handlePageChange() {
  fetchLeaveRequest();
}

function openCreateDialog() {
  isEditMode.value = false;
  editingId.value = null;
  Object.assign(form, defaultForm());
  dialogVisible.value = true;
}

function openEditDialog(row) {
  isEditMode.value = true;
  editingId.value = row.id;
  Object.assign(form, {
    class_id: row.class_id,
    start_date: row.start_date,
    end_date: row.end_date,
    back_to_work_date: row.back_to_work_date,
    total_day: row.total_day,
    deduct_type_id: row.deduct_type_id,
    reason: row.reason,
  });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    if (isEditMode.value) {
      await editLeaveRequest(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addLeaveRequest({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

async function handleApprove(row) {
  const confirmed = await confirmDelete("តេីអ្នកពិតជាចង់អនុម័តច្បាប់នេះមែនទេ?");
  if (!confirmed) return;
  try {
    await approveLeaveRequest(row.id);
    notify.success("អនុម័តច្បាប់ជោគជ័យ");
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចអនុម័តច្បាប់នេះបានទេ");
  } finally {
    
  }
}

async function handleDelete(row) {
  const confirmed = await confirmDelete("តេីអ្នកពិតជាចង់លុបច្បាប់នេះមែនទេ?");
  if (!confirmed) return;
  try {
    await deleteLeaveRequest(row.id);
    notify.success("លុបច្បាប់ជោគជ័យ");
    await fetchLeaveRequest();
  } catch (e) {
    // backend refuses when the row is already linked to payroll
    notify.error(e.response?.data?.error || "មិនអាចលុបច្បាប់នេះបានទេ");
  }
}

function statusTagType(status) {
  return status === "APPROVE" ? "success" : "warning";
}

function statusLabel(status) {
  return statusOptions.find((s) => s.value === status)?.label || status;
}

watch(
  filter,
  () => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      pagination.page = 1;
      fetchLeaveRequest();
    }, 200); // debounce so typing doesn't fire a request per keystroke
  },
  { deep: true },
);

onMounted(() => {
  fetchClasses()
  fetchDeductTypes();
  fetchLeaveRequest();
  fetchProgramme();
  fetchGeneration();
  fetchMajor();
  fetchShift();
});
onUnmounted(() => clearTimeout(searchTimer));
</script>

<template>
  <div>
    <AppFilterBar :fields="[
      { slot: 'name', span: 2 },
      { slot: 'programm', span: 2 },
      { slot: 'generation', span: 2 },
      { slot: 'major', span: 2 },
      { slot: 'shift', span: 2 },
      { slot: 'year', span: 2 },
      { slot: 'semester', span: 2 },
      { slot: 'class', span: 2 },
      { slot: 'subject', span: 2 },
      { slot: 'status', span: 2 },
      { slot: 'start_date', span: 2 },
      { slot: 'end_date', span: 2 },
    ]" :action-span="4">

      <template #name>
        <AppInput v-model="filter.name" placeholder="រកតាមឈ្មោះ" prefix-icon="Search" clearable />
      </template>
      <template #programm>
        <AppSelect v-model="filter.programme_id" :options="programmes" placeholder="កម្មវិធីសិក្សា" size="large"
          filterable clearable @change="onProgrammchange" />
      </template>
      <template #generation>
        <AppSelect v-model="filter.generation_id" :options="generations" placeholder="ជំនាន់" size="large" filterable
          clearable @change="onGenerationchange" />
      </template>
      <template #major>
        <AppSelect v-model="filter.major_id" :options="majors" placeholder="ជំនាញ" size="large" filterable clearable
          @change="onMajorchange" />
      </template>
      <template #shift>
        <AppSelect v-model="filter.shift_id" :options="shifts" placeholder="វេនសិក្សា" size="large" filterable clearable
          @change="onShiftchange" />
      </template>
      <template #year>
        <AppSelect v-model="filter.year" :options="years" placeholder="ឆ្នាំសិក្សា" size="large" clearable
          @change="onYearchange" />
      </template>
      <template #semester>
        <AppSelect v-model="filter.semester" :options="semesters" placeholder="ឆមាស" size="large" clearable
          @change="onSemesterchange" />
      </template>
      <template #class>
        <AppSelect :options="classes" filterable v-model="filter.class_id" placeholder="ថ្នាក់" clearable
          style="width: 100%" size="large" @change="onclasschange">

        </AppSelect>
      </template>

      <template #subject>
        <AppSelect v-model="filter.subject_id" :options="subjectOptions" :loading="subjectsLoading"
          placeholder="មុខវិជ្ជា" size="large" filterable clearable @change="fetchLeaveRequest" />
      </template>

      <template #status>
        <AppSelect :options="statusOptions" v-model="filter.status" placeholder="ស្ថានភាព" clearable style="width: 100%"
          size="large">

        </AppSelect>
      </template>
      <template #start_date>
        <AppDatePicker v-model="filter.start_date" type="date" placeholder="ជ្រេីសរេីសថ្ងៃទី" value-format="YYYY-MM-DD"
          clearable @change="fetchLeaveRequest" style="width: 100%" size="large">
        </AppDatePicker>
      </template>
      <template #end_date>
        <AppDatePicker v-model="filter.end_date" type="date" placeholder="ជ្រេីសរេីសថ្ងៃទី" value-format="YYYY-MM-DD"
          clearable @change="fetchLeaveRequest" style="width: 100%" size="large">
        </AppDatePicker>
      </template>
      <template #actions>
        <AppButton v-if="canCreateLeave" type="primary" @click="openCreateDialog">
          បន្ថែមច្បាប់
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable show-index :data="leaveRequests" :loading="loading" :actions-width="200"
        v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :total="pagination.total"
        @page-change="handlePageChange" :columns="[
          { label: 'ឈ្មោះ', slot: 'user', minWidth: 170 },
          { label: 'ថ្នាក់', prop: 'class_name', minWidth: 120 },
          { label: 'ថ្ងៃចាប់ផ្តើម', prop: 'start_date', minWidth: 110 },
          { label: 'ថ្ងៃបញ្ចប់', prop: 'end_date', minWidth: 110 },
          { label: 'ថ្ងៃចូលរៀនវិញ', prop: 'back_to_work_date', minWidth: 130 },
          { label: 'ចំនួនថ្ងៃ', slot: 'total_day', width: 150 },
          { label: 'មូលហេតុ', prop: 'reason', minWidth: 150 },
          { label: 'ស្ថានភាព', slot: 'status', width: 110 },
          { label: 'អ្នកអនុម័ត', slot: 'approver', minWidth: 150 },
        ]">
        <template #user="{ row }">
          <el-space direction="vertical" alignment="start" :size="0">
            <el-text tag="b" size="large" type="primary">
              {{ row.user_name_kh }} | <el-text>{{ row.user_name_en }}</el-text>
            </el-text>

            <el-text size="small">
              {{ row.user_code }} | {{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}
            </el-text>
          </el-space>
        </template>

        <template #total_day="{ row }">
          <el-text>{{ row.total_day }}{{ row.deduct_type_name }}</el-text>


        </template>

        <template #status="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>

        <template #approver="{ row }">
          <template v-if="row.approve_by">
            <div>{{ row.approve_by_name }}</div>
            <div style="color: #909399; font-size: 12px">{{ row.approved_at }}</div>
          </template>
          <el-text v-else type="info">—</el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែច្បាប់" placement="top">
            <AppButton v-if="canEditLeave && row.status === 'PENDING'" size="small" icon="Edit" type="warning" circle
              @click="openEditDialog(row)" />
          </el-tooltip>
          <el-tooltip content="អនុម័ត" placement="top">
            <AppButton v-if="canApproveLeave && row.status === 'PENDING'" size="small" icon="Check" type="success"
              circle :loading="approvingId === row.id" @click="handleApprove(row)" />
          </el-tooltip>
          <el-tooltip content="លុប" placement="top">
            <AppButton v-if="canDeleteLeave && row.status === 'PENDING'" size="small" icon="Delete" type="danger" circle
              @click="handleDelete(row)" />
          </el-tooltip>
          <el-tooltip content="ត្រឡប់" placement="top">
            <AppButton v-if="canDeleteLeave && row.status === 'APPROVE'" size="small" icon="Refresh" type="danger" circle
              @click="handleRollback(row)" />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <AppDialog v-model="dialogVisible" :title="isEditMode ? 'កែប្រែការសុំច្បាប់' : 'បង្កើតការសុំច្បាប់'" width="740px">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
          <el-form-item label="ថ្នាក់" prop="class_id">
            <AppSelect :options="classes" v-model="form.class_id" placeholder="ជ្រើសរើសថ្នាក់" style="width: 100%"
              size="large">

            </AppSelect>
          </el-form-item>
          <el-form-item label="ឯកតាកាត់ថ្ងៃ" prop="deduct_type_id">
            <el-select v-model="form.deduct_type_id" placeholder="ជ្រើសរើសឯកតា" style="width: 100%" size="large">
              <el-option v-for="dt in deductTypes" :key="dt.id" :label="`${dt.code} - ${dt.name}`" :value="dt.id" />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ថ្ងៃចាប់ផ្តើម" prop="start_date">
            <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%"
              size="large" />
          </el-form-item>
          <el-form-item label="ថ្ងៃបញ្ចប់" prop="end_date">
            <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%"
              size="large" />
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ថ្ងៃចូលរៀនវិញ" prop="back_to_work_date">
            <el-date-picker v-model="form.back_to_work_date" type="date" value-format="YYYY-MM-DD" style="width: 100%"
              size="large" />
          </el-form-item>
          <AppInput label="ចំនួនថ្ងៃសរុប" prop="total_day" v-model.number="form.total_day" type="number" size="large" />
        </div>

        <AppInput prop="reason" label="មូលហេតុ" clearable v-model="form.reason" type="textarea" :block="true" />
      </el-form>

      <template #footer>
        <AppButton @click="dialogVisible = false">បោះបង់</AppButton>
        <AppButton type="primary" :loading="submitting" @click="handleSubmit">
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
</style>