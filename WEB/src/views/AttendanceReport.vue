<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 5 },
        { slot: 'dateFrom', span: 5 },
        { slot: 'class', span: 5 },
        { slot: 'subject', span: 5 },
      ]"
      :action-span="3"
    >
      <template #name>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរក"
          clearable
          size="large"
        >
        </AppInput>
      </template>
      <template #dateFrom>
        <AppDatePicker
          v-model="filters.check_date"
          type="date"
          placeholder="ពីថ្ងៃទី"
          value-format="YYYY-MM-DD"
          clearable
          size="large"
        >
        </AppDatePicker>
      </template>
      <template #class>
        <AppSelect
          v-model="filters.class_id"
          :options="classes"
          placeholder="ជ្រើសរើសថ្នាក់"
          size="large"
          filterable
          clearable
          @change="onclasschange"
        >
        </AppSelect>
      </template>
      <template #subject>
        <AppSelect
          v-model="filters.subject_id"
          :options="subjectOptions"
          placeholder="ជ្រើសរើសមុខវិជ្ជា"
          size="large"
          filterable
          clearable
          @change="fetchReport"
        >
        </AppSelect>
      </template>
      <template #actions>
        <AppButton type="primary" @click="fetchReport">ស្វែងរក</AppButton>
      </template>
    </AppFilterBar>

    <el-card v-loading="loading">
      <div align="center" class="p-5 moul-regular">សរុបអវត្តមាននិស្សិត</div>
      <div align="center" v-if="report.rows.length" class="report-class-name">
        ថ្នាក់៖ {{ report.rows[0].class_name }} ជំនាញ៖
        {{ report.rows[0].major_name }}
        {{ report.rows[0].generation }} ក្រុមទី{{
          report.rows[0].group_name
        }}
        វគ្គទី{{ report.rows[0].term }} ឆ្នាំទី{{
          report.rows[0].year
        }}
        ឆមាសទី{{ report.rows[0].semester }} (ពេលសិក្សា វគ្គ
        {{ report.rows[0].shift_name }})
      </div>
      <div align="center" v-if="report.rows.length" class="report-class-name">
        មុខវិជ្ជា {{ report.rows[0].subject_name }}
      </div>
      <div class="report-scroll">
        <table class="report-table" v-if="report.rows.length">
          <thead>
            <tr>
              <th rowspan="3">ល.រ</th>
              <th rowspan="3">គោត្តនាម-នាម</th>
              <th rowspan="3">ភេទ</th>
              <th rowspan="3">អត្តលេខ</th>
              <th
                v-for="g in dateGroups"
                :key="g.date"
                :colspan="g.columns.length"
              >
                {{ formatDate(g.date) }}
                <!-- {{ g.date }} -->
              </th>
              <th colspan="3" rowspan="2">សរុ​ប</th>
            </tr>
            <tr>
              <th
                v-for="g in dateGroups"
                :key="g.date + '-l'"
                :colspan="g.columns.length"
              >
                ពេលសិក្សា
              </th>
            </tr>
            <tr>
              <th
                v-for="col in report.columns"
                :key="col.column_no"
                class="diagonal-text"
              >
                Session {{ col.column_no }}
              </th>
              <th>អត់ច្បាប់</th>
              <th>មានច្បាប់</th>
              <th>ពិន្ទុ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in report.rows" :key="row.user_id">
              <td>{{ row.index }}</td>
              <td class="text-left">{{ row.name_kh }}</td>
              <td>{{ genderLabel(row.gender) }}</td>
              <td>{{ row.code }}</td>
              <td
                v-for="(cell, i) in row.cells"
                :key="i"
                :class="cellClass(cell.status)"
              >
                <span>
                  {{
                    cell.status === "PR"
                      ? "✓"
                      : cell.status === "P"
                        ? "P"
                        : cell.status === "A"
                          ? "A"
                          : ""
                  }}
                </span>
              </td>
              <td>{{ row.absent_count }}</td>
              <td>{{ row.permission_count }}</td>
              <td>
                {{ 10 - row.absent_count * 1 - row.permission_count * 0.5 }}
              </td>
            </tr>
          </tbody>
        </table>
        <el-empty v-else description="គ្មានទិន្នន័យ" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";
import {
  getAttendanceReport,
  getClass,
  getClassAvailableSubjects,
} from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";
import AppSelect from "../../components/AppSelect.vue";
import AppInput from "../../components/AppInput.vue";
import AppDatePicker from "../../components/AppDatePicker.vue";
const loading = ref(false);
const classes = ref([]);
const report = reactive({ columns: [], rows: [] });
const subjectOptions = ref([]);
const subjectsLoading = ref(false);

const filters = reactive({
  name: "",
  check_date: new Date().toISOString().split("T")[0],
  class_id: null,
  subject_id: null,
});

const dateGroups = computed(() => {
  const groups = [];
  const byDate = new Map();
  for (const col of report.columns) {
    if (!byDate.has(col.check_date)) {
      byDate.set(col.check_date, { date: col.check_date, columns: [] });
      groups.push(byDate.get(col.check_date));
    }
    byDate.get(col.check_date).columns.push(col);
  }
  return groups;
});

function genderLabel(gender) {
  if (gender === 1) return "ប្រុស";
  if (gender === 2) return "ស្រី";
  return "—";
}
function cellClass(status) {
  if (status === "A") return "cell-absent";
  if (status === "PR") return "cell-present";
  if (status === "P") return "cell-permission";
  return "cell-empty";
}
function formatDate(d) {
  if (!d) return "";
  const dt = new Date(d);
  return `${dt.getMonth() + 1}/${dt.getDate()}/${dt.getFullYear()}`;
}

async function fetchClasses() {
  try {
    const res = await getClass();
    classes.value = (res.data.data || []).map((s) => ({
      label: `${s.name}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load classes");
  }
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

async function fetchReport() {
  loading.value = true;
  try {
    const params = {};
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.class_id) params.class_id = filters.class_id;
    if (filters.subject_id) params.subject_id = filters.subject_id;
    const res = await getAttendanceReport(params);
    report.columns = res.data.data?.columns || [];
    report.rows = res.data.data?.rows || [];
  } catch {
    ElMessage.error("Failed to load report");
  } finally {
    loading.value = false;
  }
}

async function onclasschange(classID) {
  filters.subject_id = null;
  subjectOptions.value = [];
  if (classID) {
    await fetchSubjectOptions(classID);
  }
  fetchReport();
}

function debounce(fn, delay = 300) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}
const debouncedFetch = debounce(fetchReport);

watch(
  () => filters.check_date,
  () => {
    debouncedFetch();
  },
);

onMounted(() => {
  fetchReport();
  fetchClasses();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Moul&display=swap");
.diagonal-text {
  height: 80px;
  vertical-align: middle;
  white-space: nowrap;
  transform: rotate(-75deg);
}
.moul-regular {
  font-family: "Moul", serif;
  font-weight: 400;
  font-style: normal;
  font-size: large;
}
.report-class-name {
  text-align: center;
  font-weight: 600;
  padding-bottom: 10px;
}
.report-scroll {
  padding-top: 10px;
  overflow-x: auto;
}
.report-table {
  border-collapse: collapse;
  width: 100%;
  font-size: 12px;
}
.report-table th,
.report-table td {
  border: 1px solid #dcdfe6;
  padding: 4px 6px;
  text-align: center;
  white-space: nowrap;
}
.text-left {
  text-align: left;
}
.cell-absent {
  color: #f56c6c;
  font-weight: 600;
}
.cell-present {
  color: #67c23a;
  font-weight: 600;
}
.cell-permission {
  color: #1e65ff;
  font-weight: 600;
}
.cell-empty {
  color: #c0c4cc;
}
@media (max-width: 767px) { .report-scroll { max-height: 70vh; } .report-table { font-size: 11px; } .report-table th, .report-table td { padding: 4px 6px; } }
</style>
