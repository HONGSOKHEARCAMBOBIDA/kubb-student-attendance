<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 5 },
        { slot: 'dateFrom', span: 5 },
        { slot: 'class', span: 5 },
      ]"
      :action-span="3"
    >
      <template #name>
        <el-input v-model="filters.name" placeholder="ស្វែងរក" clearable size="large" />
      </template>
      <template #dateFrom>
        <el-date-picker v-model="filters.check_date" type="date" placeholder="ពីថ្ងៃទី" value-format="YYYY-MM-DD" clearable style="width:100%" size="large" />
      </template>
      <template #class>
        <el-select v-model="filters.class_id" placeholder="ថ្នាក់" clearable style="width:100%" size="large">
          <el-option v-for="cls in classes" :key="cls.id" :label="cls.name" :value="cls.id" />
        </el-select>
      </template>
      <template #actions>
        <AppButton type="primary" @click="fetchReport">ស្វែងរក</AppButton>
      </template>
    </AppFilterBar>

    <el-card v-loading="loading">
      <div class="report-scroll">
        <table class="report-table" v-if="report.rows.length">
          <thead>
            <tr>
              <th rowspan="3">ល.រ</th>
              <th rowspan="3">គោត្តនាម-នាម</th>
              <th rowspan="3">ភេទ</th>
              <th rowspan="3">អត្តលេខ</th>
              <th v-for="g in dateGroups" :key="g.date" :colspan="g.columns.length">
                {{ formatDate(g.date) }}
                  <!-- {{ g.date }} -->
              </th>
              <th colspan="3" rowspan="2">សរុ​ប</th>
            </tr>
            <tr>
              <th v-for="g in dateGroups" :key="g.date + '-l'" :colspan="g.columns.length">ពេលសិក្សា</th>
            </tr>
            <tr>
              <th v-for="col in report.columns" :key="col.column_no">Session {{ col.column_no }}</th>
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
              <td v-for="(cell, i) in row.cells" :key="i" :class="cellClass(cell.status)">
                              <span>
                {{
                  cell.status === 'PR'
                    ? '✓'
                    : cell.status === 'P'
                      ? 'P'
                      : cell.status === 'A'
                        ? 'A'
                        : ''
                }}
              </span>
              </td>
              <td>{{ row.absent_count }}</td>
              <td>{{ row.permission_count }}</td>
              <td>
              {{ 10 - (row.absent_count * 1) - (row.permission_count * 0.5) }}
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
import { ref, reactive, computed, onMounted,watch } from "vue";
import { ElMessage } from "element-plus";
import { getAttendanceReport, getClass } from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";

const loading = ref(false);
const classes = ref([]);
const report = reactive({ columns: [], rows: [] });

const filters = reactive({ 
  name: "",
  check_date: new Date().toISOString().split("T")[0],
   class_id: "" });

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
    classes.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load classes");
  }
}

async function fetchReport() {
  loading.value = true;
  try {
    const params = {};
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.class_id) params.class_id = filters.class_id;

    const res = await getAttendanceReport(params);
    report.columns = res.data.data?.columns || [];
    report.rows = res.data.data?.rows || [];
  } catch {
    ElMessage.error("Failed to load report");
  } finally {
    loading.value = false;
  }
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
  }
);

onMounted(() => {
  fetchReport();
  fetchClasses();
});
</script>

<style scoped>
.report-scroll { overflow-x: auto; }
.report-table { border-collapse: collapse; width: 100%; font-size: 12px; }
.report-table th, .report-table td { border: 1px solid #dcdfe6; padding: 4px 6px; text-align: center; white-space: nowrap; }
.text-left { text-align: left; }
.cell-absent { color: #f56c6c; font-weight: 600; }
.cell-present { color: #67c23a; font-weight: 600; }
.cell-permission { color: #1e65ff; font-weight: 600; }
.cell-empty { color: #c0c4cc; }
</style>