<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`ពិន្ទុថ្នាក់ — ${classRow?.name || ''}`"
    width="75%"
  >
    <div class="form-row">
      <AppSelect
        v-model="filters.subject_id"
        :options="subjectOptions"
        label="មុខវិជ្ជា"
        placeholder="ជ្រើសរើសមុខវិជ្ជា (ទាំងអស់)"
        clearable
        size="large"
        @change="() => { page = 1; fetchScoreView(); }"
      />
      <AppInput
        v-model.trim="filters.name"
        label="ស្វែងរកឈ្មោះ/អត្តលេខ"
        placeholder="ឈ្មោះ ឬ អត្តលេខ"
        clearable
        size="large"
        @input="debouncedFetch"
      />
    </div>

    <AppTable :data="scoreData" v-loading="loading" :columns="columns" show-index :show-pagination="false">
      <template #gender="{ row }">
        <el-text>{{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}</el-text>
      </template>
    </AppTable>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
// import { getClassAvailableSubjects, getScore } from "../api/services";
import AppDialog from "./AppDialog.vue";
import AppTable from "./AppTable.vue";
import AppSelect from "./AppSelect.vue";
import AppInput from "./AppInput.vue";
import { getClassAvailableSubjects,getScore } from "../src/api/services.js";
import { useNotification } from "../composables/useNotification.js";
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classRow: { type: Object, default: null },
});
defineEmits(["update:modelValue"]);

const notify = useNotification();
const loading = ref(false);
const scoreData = ref([]);
const subjectOptions = ref([]);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const filters = reactive({ subject_id: null, name: "" });

const columns = [
  { prop: "name_kh", label: "ឈ្មោះខ្មែរ", minwidth: 100 },
  { slot: "gender", label: "ភេទ", width: 70 },
  { prop: "code", label: "អត្តលេខ", width: 100 },
  { prop: "ProgrammeName", label: "កម្រិត", width: 100 },
  { prop: "GenerationName", label: "ជំនាន់", width: 100 },
  { prop: "MajorName", label: "ជំនាញ", minwidth: 100 },
  { prop: "SubjectName", label: "មុខវិជ្ជា", minwidth: 100 },
  { prop: "attendance", label: "វត្តមាននិស្សិត", minwidth: 100 },
  { prop: "research", label: "កិច្ចការស្រាវជ្រាវ", minwidth: 100 },
  { prop: "midterm", label: "ប្រឡងពាក់កណ្តាលឆមាស", minwidth: 150 },
  { prop: "final", label: "ប្រឡងបញ្ចប់ឆមាស", minwidth: 150 },
];

let debounceTimer = null;
function debouncedFetch() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    page.value = 1;
    fetchScoreView();
  }, 400);
}

async function fetchScoreView() {
  const row = props.classRow;
  if (!row) return;

  loading.value = true;
  try {
    const res = await getScore({
      page: page.value,
      page_size: pageSize.value,
      class_id: row.id, // class's own primary key = score.class_id
      generation_id: row.generation_id,
      major_id: row.major_id,
      programme_id: row.programme_id,
      year: row.year,
      semester: row.semester,
      subject_id: filters.subject_id || undefined,
      name: filters.name || undefined,
    });
    scoreData.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចទាញយកពិន្ទុបានទេ");
  } finally {
    loading.value = false;
  }
}

watch(
  () => props.modelValue,
  async (open) => {
    if (!open || !props.classRow) return;
    filters.subject_id = null;
    filters.name = "";
    page.value = 1;

    try {
      const res = await getClassAvailableSubjects(props.classRow.id);
      subjectOptions.value = (res.data.data || []).map((s) => ({
        label: `${s.code} — ${s.name_kh}`,
        value: s.id,
      }));
    } catch (e) {
      notify.error(e.response?.data?.error || "Failed to load subjects");
    }

    fetchScoreView();
  },
);
</script>

<style scoped>
.form-row {
  padding: 10px;
  display: flex;
  gap: 16px;
}
.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
</style>
