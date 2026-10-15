<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 4 },
        { slot: 'check_date', span: 4 },
        { slot: 'class', span: 4 },
        { slot: 'major', span: 4 },
        { slot: 'shift', span: 4 },
        { slot: 'generation', span: 4 },
      ]"
    >
      <template #name>
        <el-input
          v-model.trim="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះសិស្ស"
          clearable
          size="large"
        />
      </template>

      <template #check_date>
        <el-date-picker
          v-model="filters.check_date"
          type="date"
          placeholder="ជ្រើសរើសកាលបរិច្ឆេទ"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          size="large"
          style="width: 100%"
          clearable="false"
        />
      </template>

      <template #class>
        <AppSelect
          v-model="filters.class_id"
          placeholder="ថ្នាក់"
          clearable
          filterable
          size="large"
          :options="classOptions"
        >
        </AppSelect>
      </template>

      <template #major>
        <el-select
          v-model="filters.major_id"
          placeholder="ជំនាញ"
          clearable
          filterable
          size="large"
        >
          <el-option
            v-for="item in majors"
            :key="item.id"
            :label="item.name_kh"
            :value="item.id"
          />
        </el-select>
      </template>

      <template #shift>
        <el-select
          v-model="filters.shift_id"
          placeholder="វេន"
          clearable
          size="large"
        >
          <el-option
            v-for="item in shifts"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
      </template>

      <template #generation>
        <el-select
          v-model="filters.generation_id"
          placeholder="ជំនាន់"
          clearable
          filterable
          size="large"
        >
          <el-option
            v-for="item in generations"
            :key="item.id"
            :label="item.name_kh"
            :value="item.id"
          />
        </el-select>
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
        <template #name="{row}">
          <el-text size="large" tag="b" type="primary">{{ row.user_name_kh }} | <el-text>{{ row.user_name_en }}</el-text></el-text>
        </template>
        <template #gender="{ row }">
          <el-text :type="row.gender === 1 ? 'primary' : 'danger'" size="small">
            {{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}
          </el-text>
        </template>
        <template #class="{row}">
          <el-text>{{ row.class_name }}-{{ row.programme_name }}-ឆ្នាំទី{{ row.year }}-ឆមាសទី{{ row.semester }}</el-text>
        </template>
        <template #major="{row}">
          <el-text size="large" type="primary" tag="b">{{ row.major_name }}</el-text>
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
} from "../api/services";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import { useNotification } from "../../composables/useNotification.js";
import AppSelect from "../../components/AppSelect.vue";
const notify = useNotification();
const selectedRows = ref([]);
const submitting = ref(false);
const list = ref([]);
const classOptions = ref([]);
const majors = ref([]);
const shifts = ref([]);
const generations = ref([]);

const loading = ref(false);
const exporting = ref(false);
const page = ref(1);
const pageSize = ref(50);
const total = ref(0);

function todayStr() {
  return new Date().toISOString().slice(0, 10);
}

const filters = reactive({
  name: "",
  check_date: todayStr(),
  class_id: "",
  major_id: "",
  shift_id: "",
  generation_id: "",
});

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
    const res = await getNotPermissionLeave({
      page: page.value,
      page_size: pageSize.value,
      name: filters.name || undefined,
      check_date: filters.check_date || undefined,
      class_id: filters.class_id || undefined,
      major_id: filters.major_id || undefined,
      shift_id: filters.shift_id || undefined,
      generation_id: filters.generation_id || undefined,
    });
    list.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
    console.log(list.value)
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

async function fetchLookups() {
  try {
    const [classRes, majorRes, shiftRes, generationRes] = await Promise.all([
      viewcompanyscan(),
      getMajor(),
      getShift(),
      getGeneration(),
    ]);
    classOptions.value = (classRes.data.data || []).map((s) => ({
      label: `${s.name}-${s.programme_name}-ឆ្នាំទី${s.year}-ឆមាសទី${s.semester}`,
      value: s.id,
    }))
    majors.value = majorRes.data.data || [];
    shifts.value = shiftRes.data.data || [];
    generations.value = generationRes.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

async function handleExport() {
  exporting.value = true;
  try {
    const res = await exportNotPermissionLeave(
      {
        name: filters.name || undefined,
        check_date: filters.check_date || undefined,
        class_id: filters.class_id || undefined,
        major_id: filters.major_id || undefined,
        shift_id: filters.shift_id || undefined,
        generation_id: filters.generation_id || undefined,
      },
      { responseType: "blob" },
    );

    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute(
      "download",
      `not-permission-leave-${filters.check_date}.xlsx`,
    );
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (e) {
    notify.error(e.response?.data?.error || "នាំចេញបរាជ័យ");
  } finally {
    exporting.value = false;
  }
}

watch(
  () => filters.name,
  () => {
    page.value = 1;
    debounce(() => fetchData());
  },
);

watch(
  [
    () => filters.check_date,
    () => filters.class_id,
    () => filters.major_id,
    () => filters.shift_id,
    () => filters.generation_id,
  ],
  () => {
    page.value = 1;
    fetchData();
  },
);

onMounted(() => {
  fetchData();
  fetchLookups();
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
