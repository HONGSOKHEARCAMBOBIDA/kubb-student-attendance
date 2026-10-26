<template>
  <div>
    <h1 style="margin-bottom: 16px">របាយការណ៍វត្តមានសាស្រ្តាចារ្យ</h1>
    <AppFilterBar
      :fields="[
        { slot: 'programm', span: 5 },
        { slot: 'generation', span: 5 },
        { slot: 'startDate', span: 5 },
        { slot: 'endDate', span: 5 },
      ]"
      :action-span="4"
    >
      <template #programm>
        <AppSelect
          v-model="filters.programme_id"
          :options="programmes"
          placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
          size="large"
          filterable
          clearable
          style="width: 100%"
          @change="onFilterChange"
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
          style="width: 100%"
          @change="onFilterChange"
        />
      </template>

      <template #startDate>
        <AppDatePicker
          v-model="filters.start_date"
          type="date"
          placeholder="ថ្ងៃចាប់ផ្ដើម"
          value-format="YYYY-MM-DD"
          clearable
          style="width: 100%"
          size="large"
          @change="onFilterChange"
        />
      </template>

      <template #endDate>
        <AppDatePicker
          v-model="filters.end_date"
          type="date"
          placeholder="ថ្ងៃបញ្ចប់"
          value-format="YYYY-MM-DD"
          clearable
          style="width: 100%"
          size="large"
          @change="onFilterChange"
        />
      </template>
    </AppFilterBar>

    <el-card>
      <AppTable
        :data="schedules"
        :loading="loading"
        show-index
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchSchedule"
        :columns="[
          { prop: 'name_kh', label: 'ឈ្មោះសាស្ត្រាចារ្យ', minWidth: 140 },
          { slot: 'gender', label: 'ភេទ', width: 60,align:'center' },
          { prop: 'subject_name_kh', label: 'មុខវិជ្ជា', minWidth: 160 },
          { prop: 'major_name', label: 'ជំនាញ', minWidth: 140 },
                 { prop: 'programme_name', label: 'កម្រិត', minWidth: 100 },
          {
            prop: 'generation_name',
            label: 'ជំនាន់',
            minWidth: 90,
            align: 'center',
          },
          {
            slot: 'day_of_week',
            label: 'ថ្ងៃបង្រៀន',
            minWidth: 90,
            align: 'center',
          },
          {
            prop: 'schedule_date',
            label: 'កាលបរិច្ឆេទ',
            minWidth: 110,
            align: 'center',
          },
   

          {
            prop: 'shift_name',
            label: 'វេនសិក្សា',
            minWidth: 100,
            align: 'center',
          },

          {
            prop: 'phone',
            label: 'លេខទូរសព្ទ',
            minWidth: 100,
            align: 'center',
          },
          { prop: 'room', label: 'បន្ទប់', minWidth: 80, align: 'center' },
          { prop: '', label: 'ព្រឹក', minWidth: 80, align: 'center' },
          { prop: '', label: 'រសៀល', minWidth: 80, align: 'center' },
        ]"
      >
      <template #day_of_week="{row}">
        <el-text>{{ dayLabel(row.day_of_week) }}</el-text>
      </template>
        <template #gender="{ row }">
          <el-text>{{ row.gender === 1 ? "ប" : "ស" }}</el-text>
        </template>
      </AppTable>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getWeekSchedule, getProgramme, getGeneration } from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppTable from "../../components/AppTable.vue";
import AppSelect from "../../components/AppSelect.vue";
import AppDatePicker from "../../components/AppDatePicker.vue";

const schedules = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

const programmes = ref([]);
const generations = ref([]);

// format using local time (toISOString() can shift the date back one day in UTC+7)
function formatLocalDate(d) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

const DAYS = [
  { value: 1, label: "ច័ន្ទ" },
  { value: 2, label: "អង្គារ" },
  { value: 3, label: "ពុធ" },
  { value: 4, label: "ព្រហស្បតិ៍" },
  { value: 5, label: "សុក្រ" },
  { value: 6, label: "សៅរ៍" },
  { value: 7, label: "អាទិត្យ" },
];
function dayLabel(v) {
  return DAYS.find((d) => d.value === v)?.label || v;
}

// Monday -> Sunday of the current week
function getCurrentWeekRange() {
  const today = new Date();
  const dayOfWeek = today.getDay(); // 0 = Sunday
  const diffToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;

  const monday = new Date(today);
  monday.setDate(today.getDate() + diffToMonday);

  const sunday = new Date(monday);
  sunday.setDate(monday.getDate() + 6);

  return [formatLocalDate(monday), formatLocalDate(sunday)];
}

const [weekStart, weekEnd] = getCurrentWeekRange();

const filters = reactive({
  programme_id: null,
  generation_id: null,
  start_date: weekStart,
  end_date: weekEnd,
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
    ElMessage.error("Failed to load generation");
  }
}

function onFilterChange() {
  page.value = 1;
  fetchSchedule();
}

function onDateChange(val) {
  // val is null when the picker is cleared
  filters.start_date = val ? val[0] : "";
  filters.end_date = val ? val[1] : "";
  page.value = 1;
  fetchSchedule();
}

async function fetchSchedule() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.programme_id) params.programme_id = filters.programme_id;
    if (filters.generation_id) params.generation_id = filters.generation_id;
    if (filters.start_date) params.start_date = filters.start_date;
    if (filters.end_date) params.end_date = filters.end_date;

    const res = await getWeekSchedule(params);
    schedules.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load schedule");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchSchedule();
  fetchProgramme();
  fetchGeneration();
});
</script>

<style scoped></style>
