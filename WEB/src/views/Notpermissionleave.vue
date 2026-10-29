<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 2 },
        { slot: 'programm', span: 2 },
        { slot: 'generation', span: 2 },
        { slot: 'major', span: 2 },
        { slot: 'shift', span: 2 },
        { slot: 'year', span: 2 },
        { slot: 'semester', span: 2 },

        { slot: 'class', span: 6 },
        { slot: 'subject', span: 2 },
        { slot: 'check_date', span: 2 },
      ]"
    >
      <template #name>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរក"
          prefix-icon="Search"
          clearable
          @change="fetchData"
          size="large"
        >
        </AppInput>
      </template>

      <template #check_date>
        <AppDatePicker
          v-model="filters.check_date"
          type="date"
          placeholder="ជ្រេីសរេីសថ្ងៃទី"
          value-format="YYYY-MM-DD"
          clearable
          @change="fetchData"
          style="width: 100%"
          size="large"
        >
        </AppDatePicker>
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
          @change="fetchData"
        />
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
    </AppFilterBar>

    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <el-text tag="b">
            សិស្សមិនបានស្កែនវត្តមាន និងមិនបានសុំច្បាប់ —
            {{ filters.check_date }} សិស្សសរុប {{ total }}នាក់
          </el-text>
          <AppButton
            size="small"
            icon="Upload"
            type="primary"
            :loading="submitting"
            :disabled="!selectedRows.length"
            @click="handleAddNotPermission"
          >
            បញ្ចូនទិន្ន័យ
          </AppButton>
        </div>
      </template>

      <AppTable
        :data="list"
        :loading="loading"
        selectable
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchData"
        :columns="[
          { prop: 'user_code', label: 'អត្តលេខ', width: 150 },
          { slot: 'name', label: 'ឈ្មោះ', minWidth: 140 },
          { slot: 'gender', label: 'ភេទ', width: 80 },
          { slot: 'class', label: 'ថ្នាក់', minWidth: 120 },
          { slot: 'major', label: 'ជំនាញ', minWidth: 110 },
          { prop: 'shift_name', label: 'វេន', width: 100 },
          { prop: 'generation_name', label: 'ជំនាន់', width: 110 },
          { prop: 'programme_name', label: 'កម្មវិធីសិក្សា', minWidth: 130 },
          { prop: 'year', label: 'ឆ្នាំ', width: 80 },
          { prop: 'semester', label: 'ឆមាស', width: 80 },
          { prop: 'group', label: 'ក្រុម', width: 80 },
          { prop: 'term', label: 'វគ្គ', width: 80 },
        ]"
        @selection-change="(rows) => (selectedRows = rows)"
      >
        <template #name="{ row }">
          <el-text size="large" tag="b" type="primary"
            >{{ row.user_name_kh }} |
            <el-text>{{ row.user_name_en }}</el-text></el-text
          >
        </template>
        <template #gender="{ row }">
          <el-text :type="row.gender === 1 ? 'primary' : 'danger'" size="small">
            {{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}
          </el-text>
        </template>
        <template #class="{ row }">
          <el-text
            >{{ row.class_name }}-{{ row.programme_name }}-ឆ្នាំទី{{
              row.year
            }}-ឆមាសទី{{ row.semester }}</el-text
          >
        </template>
        <template #major="{ row }">
          <el-text size="large" type="primary" tag="b">{{
            row.major_name
          }}</el-text>
        </template>
      </AppTable>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from "vue";
import {
  getNotPermissionLeave,
  exportNotPermissionLeave,
  viewcompanyscan,
  getMajor,
  getShift,
  getGeneration,
  addNotPermissionLeave,
  getProgramme,
  getClassAvailableSubjects,
  getClass,
} from "../api/services";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import { useNotification } from "../../composables/useNotification.js";
import AppSelect from "../../components/AppSelect.vue";
import AppDatePicker from "../../components/AppDatePicker.vue";
import AppInput from "../../components/AppInput.vue";
const notify = useNotification();
const selectedRows = ref([]);
const submitting = ref(false);
const list = ref([]);



const loading = ref(false);
const page = ref(1);
const pageSize = ref(50);
const total = ref(0);


// state filter
const classes = ref([]);
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
  fetchData();
}

async function onGenerationchange(generationID) {
  filterClass.generation_id = generationID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchData();
}

async function onMajorchange(majorID) {
  filterClass.major_id = majorID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchData();
}

async function onShiftchange(shiftID) {
  filterClass.shift_id = shiftID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchData();
}

async function onYearchange(yearID) {
  filterClass.year = yearID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchData();
}

async function onSemesterchange(semesterID) {
  filterClass.semester = semesterID;
  filters.class_id = null;
  filters.subject_id = null;
  classes.value = [];
  subjectOptions.value = [];
  await fetchClasses(); //
  page.value = 1;
  fetchData();
}

async function onclasschange(classID) {
  filters.subject_id = null;
  subjectOptions.value = [];
  if (classID) {
    await fetchSubjectOptions(classID);
  }
  fetchData();
}

const subjectOptions = ref([]);
const subjectsLoading = ref(false);
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

let debounceTimer = null;
function debounce(fn, delay = 400) {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fn, delay);
}

async function handleAddNotPermission() {
  if (!selectedRows.value.length) {
    notify.error("សូមជ្រើសរើសសិស្សយ៉ាងហោចណាស់ម្នាក់");
    return;
  }
  submitting.value = true;
  try {
    await addNotPermissionLeave({
      check_date: filters.check_date,
      data: selectedRows.value.map((r) => ({
        user_id: r.user_id,
        class_id: r.class_id,
        shift_id: r.shift_id,
      })),
    });
    notify.success("បានកត់ត្រាជោគជ័យ");
    selectedRows.value = [];
    fetchData();
  } catch (e) {
    notify.error(e.response?.data?.error || "បរាជ័យ");
  } finally {
    submitting.value = false;
  }
}

async function fetchData() {
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
    const res = await getNotPermissionLeave(params);
    list.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
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

onMounted(() => {
  fetchClasses();
  fetchProgramme();
  fetchGeneration();
  fetchMajor();
  fetchShift();
  fetchData();
});
</script>

<style scoped>
.table-card {
  border-radius: 6px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
