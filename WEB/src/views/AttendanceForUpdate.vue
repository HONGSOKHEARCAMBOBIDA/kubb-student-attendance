<template>
  <div>
    <!-- Filters -->
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 5 },
        { slot: 'date', span: 5 },
        { slot: 'class', span: 5 },
        { slot: 'subject', span: 5 },
      ]"
      :action-span="3"
    >
      <template #name>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរក"
          prefix-icon="Search"
          clearable
          @change="fetchAttendance"
          size="large"
        />
      </template>
      <template #date>
        <AppDatePicker
          v-model="filters.check_date"
          type="date"
          placeholder="ជ្រេីសរេីសថ្ងៃទី"
          value-format="YYYY-MM-DD"
          clearable
          @change="fetchAttendance"
          style="width: 100%"
          size="large"
        />
      </template>
      <template #class>
        <AppSelect
          :options="classes"
          v-model="filters.class_id"
          placeholder="ថ្នាក់"
          clearable
          style="width: 100%"
          size="large"
          @change="onclasschange"
        />
      </template>
      <template #subject>
        <AppSelect
          v-model="filters.subject_id"
          :options="subjectOptions"
          :loading="subjectsLoading"
          placeholder="ជ្រើសរើសមុខវិជ្ជា"
          size="large"
          filterable
          clearable
          @change="fetchAttendance"
        />
      </template>
      <template #actions>
        <AppButton type="primary" @click="fetchAttendance"> ស្វែងរក </AppButton>
      </template>
    </AppFilterBar>

    <el-card>
      <AppTable
        :data="attendance"
        :loading="loading"
        show-index
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchAttendance"
        :columns="[
          { prop: 'name_kh', label: 'ឈ្មោះ', minWidth: 130 },
          { prop: 'name_en', label: 'Name', minWidth: 130 },
          { prop: 'code', label: 'កូដ', minWidth: 90 },
          { label: 'ភេទ', slot: 'gender', minWidth: 80 },
          { prop: 'class_name', label: 'ថ្នាក់', minWidth: 110 },
          { prop: 'subject_name', label: 'មុខវិជ្ជា', minWidth: 110 },
          { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 110 },
          { prop: 'status', label: 'ស្ថានភាពរួម', minWidth: 110 },
          { label: 'កំណត់ត្រា', slot: 'records', minWidth: 260 },
        ]"
      >
        <template #gender="{ row }">
          {{ genderLabel(row.gender) }}
        </template>

        <!-- render row.record (RecordResponse[]) exactly as it comes from the backend -->
        <template #records="{ row }">
          <div style="display: flex; flex-wrap: wrap; gap: 6px">
            <el-tag
              size="large"
              v-for="record in row.record"
              :key="record.id"
              :type="tagType(record.status)"
              style="cursor: pointer"
              @click="openStatusEditor(row, record)"
            >
              {{ record.type }} — {{ statusLabel[record.status] || record.status || "—" }}
            </el-tag>
            <span v-if="!row.record || row.record.length === 0" style="color: #c0c4cc">
              គ្មានកំណត់ត្រា
            </span>
          </div>
        </template>
      </AppTable>
    </el-card>

    <!-- Edit-status dialog, shared by every record tag -->
    <el-dialog v-model="editDialogVisible" title="កែប្រែស្ថានភាពវត្តមាន" width="360px">
      <div v-if="editingRecord">
        <p style="margin-bottom: 12px; color: #606266">
          {{ editingRow?.name_kh }} — {{ editingRecord.type }}
        </p>
        <el-select v-model="editingStatus" placeholder="ជ្រើសរើសស្ថានភាព" style="width: 100%"  size="large">
          <el-option
           
            v-for="opt in statusOptions"
            :key="opt.value"
            :label="opt.label"
            :value="opt.value"
          />
        </el-select>
      </div>

      <template #footer>
        <AppButton @click="editDialogVisible = false">បោះបង់</AppButton>
        <AppButton
          type="primary"
          :loading="updating"
          :disabled="!editingRecord"
          @click="confirmUpdateStatus"
        >
          រក្សាទុក
        </AppButton>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage } from "element-plus";
import {
  getAttendanceforedit,
  updateAttendanceRecordStatus,
  viewcompanyscan,
  getClassAvailableSubjects,
} from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";
import AppTable from "../../components/AppTable.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useUserDataStore } from "../stores/user_data.js";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppDatePicker from "../../components/AppDatePicker.vue";

const userDataStore = useUserDataStore();
const notify = useNotification();

const attendance = ref([]); // AttendanceResponse[]
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const classes = ref([]);
const subjectOptions = ref([]);
const subjectsLoading = ref(false);
const filters = reactive({
  name: "",
  check_date: new Date().toISOString().split("T")[0],
  class_id: null,
  subject_id: null,
});

const candeleteattendance = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "delete.backup")
);

// gender comes back as an int from the backend — adjust mapping to match your DB convention
function genderLabel(gender) {
  if (gender === 1) return "ប្រុស";
  if (gender === 2) return "ស្រី";
  return "—";
}

// --- status display config ---------------------------------------------
// RecordResponse.Status is a plain string from the backend — adjust these
// values/labels to match whatever your `attendance_record.status` column
// actually stores.
const statusOptions = [
  { label: "មកដល់ទាន់ម៉ោង", value: "PR" },
  { label: "អវត្តមាន", value: "A" },
  { label: "មានច្បាប់", value: "P" },
];
const statusLabel = Object.fromEntries(statusOptions.map((o) => [o.value, o.label]));
const statusTagType = { PR: "success", A: "danger", P: "primary" };
function tagType(status) {
  return statusTagType[status] || "info";
}

// --- edit-status dialog state --------------------------------------------
const editDialogVisible = ref(false);
const editingRow = ref(null); // the AttendanceResponse row
const editingRecord = ref(null); // the specific RecordResponse being edited
const editingStatus = ref("");
const updating = ref(false);

function openStatusEditor(row, record) {
  editingRow.value = row;
  editingRecord.value = record;
  editingStatus.value = record.status;
  editDialogVisible.value = true;
}

async function confirmUpdateStatus() {
  if (!editingRecord.value) return;
  updating.value = true;
  try {
    await updateAttendanceRecordStatus(editingRecord.value.id, {
      status: editingStatus.value,
    });
    // reflect the change locally instead of refetching the whole page
    const row = editingRow.value;
    const record = (row.record || []).find((r) => r.id === editingRecord.value.id);
    if (record) record.status = editingStatus.value;
    notify.success("ធ្វើបច្ចុប្បន្នភាពស្ថានភាពជោគជ័យ");
    editDialogVisible.value = false;
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ");
  } finally {
    updating.value = false;
  }
}

// --- data loading -----------------------------------------------------------
async function onclasschange(classID) {
  filters.subject_id = null;
  subjectOptions.value = [];
  if (classID) {
    await fetchSubjectOptions(classID);
  }
  fetchAttendance();
}

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

async function fetchClasses() {
  try {
    const res = await viewcompanyscan();
    classes.value = (res.data.data || []).map((s) => ({
      label: `${s.name}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load classes");
  }
}

async function fetchAttendance() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.class_id) params.class_id = filters.class_id;
    if (filters.subject_id) params.subject_id = filters.subject_id;
    const res = await getAttendanceforedit(params);
    attendance.value = res.data.data || []; // AttendanceResponse[]
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load attendance");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchAttendance();
  fetchClasses();
});
</script>

<style scoped>
.checkin-card {
  border-radius: 6px;
}
.filter-card {
  border-radius: 6px;
}
.card-title {
  font-weight: 600;
  font-size: 15px;
}
.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.detail-action {
  padding: 10px 0;
}
</style>