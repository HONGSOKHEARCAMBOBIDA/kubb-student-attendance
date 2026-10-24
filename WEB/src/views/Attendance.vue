<template>
  <div>
    <!-- Filters -->
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 2 },
        { slot: 'programm', span: 2 },
        { slot: 'generation', span: 2 },
        { slot: 'major', span: 2 },
        { slot: 'shift', span: 2 },
        { slot: 'year', span: 2 },
        { slot: 'semester', span: 2 },
        { slot: 'class', span: 5 },

        { slot: 'subject', span: 2 },
        { slot: 'date', span: 3 },
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
        >
        </AppInput>
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
        >
        </AppDatePicker>
      </template>
      <template #class>
        <AppSelect
          :options="classes"
          v-model="filters.class_id"
          placeholder="ថ្នាក់"
          filterable
          clearable
          style="width: 100%"
          size="large"
          @change="onclasschange"
        >
        </AppSelect>
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
      <template #programm>
        <AppSelect
          v-model="filters.programme_id"
          :options="programmes"
          placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
          size="large"
          filterable
          clearable
          @change="onProgrammchange"
        />
      </template>
      <template #generation>
        <AppSelect
          v-model="filters.generation_id"
          :options="generations"
          placeholder="ជ្រើសរើសជំនាន់"
          size="large"
          filterable
          clearable
          @change="onGenerationchange"
        />
      </template>
      <template #major>
        <AppSelect
          v-model="filters.major_id"
          :options="majors"
          placeholder="ជ្រើសរើសជំនាញ"
          size="large"
          filterable
          clearable
          @change="onMajorchange"
        />
      </template>
      <template #shift>
        <AppSelect
          v-model="filters.shift_id"
          :options="shifts"
          placeholder="ជ្រើសរើសវេនសិក្សា"
          size="large"
          filterable
          clearable
          @change="onShiftchange"
        />
      </template>
      <template #year>
        <AppSelect
          v-model="filters.year"
          :options="years"
          placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
          size="large"
          clearable
          @change="onYearchange"
        />
      </template>
      <template #semester>
        <AppSelect
          v-model="filters.semester"
          :options="semesters"
          placeholder="ជ្រើសរើសឆមាស"
          size="large"
          clearable
          @change="onSemesterchange"
        />
      </template>
      <!-- <template #actions>
        <AppButton type="primary" @click="fetchAttendance"> ស្វែងរក </AppButton>
      </template> -->
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
          { slot: 'name', label: 'ឈ្មោះ', minWidth: 200 },
          { prop: 'programme_name', label: 'កម្រិត', minWidth: 90 },
          { prop: 'class_name', label: 'ថ្នាក់', minWidth: 110 },
          { prop: 'major_name', label: 'ជំនាញ', minWidth: 110 },
          { prop: 'faculty_name', label: 'មហាវិទ្យាល័យ', minWidth: 110 },
          { prop: 'shift_name', label: 'វេនសិក្សា', minWidth: 100,align:'center'  },
          { prop: 'generation_name', label: 'ជំនាន់', minWidth: 90,align:'center'  },
          { prop: 'year', label: 'ឆ្នាំ', minWidth: 40,align:'center' },
          { prop: 'semester', label: 'ឆមាស', minWidth: 50,align:'center' },
          { prop: 'group', label: 'ក្រុម', minWidth: 50,align:'center' },
          { prop: 'term', label: 'វគ្គ', minWidth:50,align:'center' },
          { prop: 'subject_name', label: 'មុខវិជ្ជា', minWidth: 120 },
          { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 80,align:'center' },
          { label: 'ម៉ោងទី ១', slot: 'session1', minWidth: 60,align:'center'  },
          { label: 'ម៉ោងទី ២', slot: 'session2', minWidth: 60,align:'center'  },
          { label: 'ម៉ោងទី ៣', slot: 'session3', minWidth: 60,align:'center'  },
          { label: 'ម៉ោងទី ៤', slot: 'session4', minWidth: 60,align:'center'  },
          { label: 'ម៉ោងទី ៥', slot: 'session5', minWidth: 60,align:'center'  },
        ]"
      >
<template #name="{ row }">
  <div style="display: flex; flex-direction: column; text-align: start; gap: 2px;">
    <div>
      <el-text type="primary" tag="b">
        {{ row.name_kh }} |
      </el-text>
      <el-text type="default">
        {{ row.name_en }} | {{ genderLabel(row.gender) }}
      </el-text>
    </div>

    <div>
      <el-text type="info" size="small">
      {{ row.code }}
    </el-text>
    </div>
  </div>
</template>

        <template #status="{ row }">
          <el-tag
            :type="row.status === 'COMPLETE' ? 'success' : 'warning'"
            size="small"
          >
            {{ row.status }}
          </el-tag>
        </template>

        <template #session1="{ row }"
          ><CheckCell :time="row.session1"
        /></template>
        <template #session2="{ row }"
          ><CheckCell :time="row.session2"
        /></template>
        <template #session3="{ row }"
          ><CheckCell :time="row.session3"
        /></template>
        <template #session4="{ row }"
          ><CheckCell :time="row.session4"
        /></template>
        <template #session5="{ row }"
          ><CheckCell :time="row.session5"
        /></template>

        <!-- <template #actions="{ row }">
          <AppButton
            v-if="candeleteattendance"
            size="small"
            icon="Delete"
            type="danger"
            circle
            @click="deleteattendancev1(row)"
          />
        </template> -->
      </AppTable>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h, computed } from "vue";
import { ElMessage } from "element-plus";
import {
  getAttendance,
  exportAttendancePDF,
  viewcompanyscan,
  deleteattendance,
  getClassAvailableSubjects,
  getProgramme,
  getClass,
  getGeneration,
  getMajor,
  getShift,
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

const attendance = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const classes = ref([]);
const subjectOptions = ref([]);
const subjectsLoading = ref(false);

// state filter
const programmes = ref([]);
const generations = ref([]);
const shifts = ref([]);
const majors = ref([]);

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
// state filter
const filters = reactive({
  name: "",
  check_date: new Date().toISOString().split("T")[0],
  class_id: null,
  subject_id: null,
  programme_id: null,
  generation_id: null,
  major_id: null,
  shift_id: null,
  year: null,
  semester: null,
});

const filterClass = reactive({
  major_id: "",
  shift_id: "",
  generation_id: "",
  programme_id: "",
  year: "",
  semester: "",
});

const candeleteattendance = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "delete.backup"),
);

// gender comes back as an int from the backend — adjust mapping to match your DB convention
function genderLabel(gender) {
  if (gender === 1) return "ប្រុស";
  if (gender === 2) return "ស្រី";
  return "—";
}

// backend has no per-session "diff" field, so this just renders the time or a dash
const CheckCell = (props) => {
  if (!props.time) {
    return h("span", { style: "color:#c0c4cc" }, "—");
  }
  return h("div", { style: "font-weight:600" }, props.time);
};
CheckCell.props = ["time"];

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
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

async function onGenerationchange(generationID) {
  filterClass.generation_id = generationID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

async function onMajorchange(majorID) {
  filterClass.major_id = majorID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

async function onShiftchange(shiftID) {
  filterClass.shift_id = shiftID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

async function onYearchange(yearID) {
  filterClass.year = yearID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

async function onSemesterchange(semesterID) {
  filterClass.semester = semesterID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchAttendance();
}

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

async function deleteattendancev1(row) {
  try {
    await deleteattendance(row.id);
    notify.success("Delete success");
    await fetchAttendance();
  } catch (e) {
    notify.error(
      e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ",
    );
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
    if (filters.programme_id) params.programme_id = filters.programme_id;
    if (filters.generation_id) params.generation_id = filters.generation_id;
    if (filters.major_id) params.major_id = filters.major_id;
    if (filters.shift_id) params.shift_id = filters.shift_id;
    if (filters.year) params.year = filters.year;
    if (filters.semester) params.semester = filters.semester;
    const res = await exportAttendancePDF(params);
    attendance.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load attendance");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchAttendance();
  // fetchClasses();
  fetchProgramme();
  fetchGeneration();
  fetchMajor();
  fetchShift();
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
