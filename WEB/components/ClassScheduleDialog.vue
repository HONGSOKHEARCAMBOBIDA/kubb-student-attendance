<script setup>
import { ref, reactive, watch, computed } from "vue";
import { useNotification } from "../composables/useNotification.js";
import AppTable from "./AppTable.vue";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";

import {
  getClassSchedule,
  getClassAvailableSubjects,
  createClassSchedule,
  toggleClassSchedule,
} from "../src/api/services.js";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classRow: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue"]);

const notify = useNotification();

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});

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

const schedules = ref([]);
const loading = ref(false);
const togglingId = ref(null);

const subjectOptions = ref([]);
const subjectsLoading = ref(false);

const formRef = ref();
const defaultForm = () => ({ subject_id: null, day_of_week: null });
const form = reactive(defaultForm());
const rules = {
  subject_id: [{ required: true, message: "សូមជ្រើសរើសមុខវិជ្ជា", trigger: "change" }],
  day_of_week: [{ required: true, message: "សូមជ្រើសរើសថ្ងៃ", trigger: "change" }],
};
const submitting = ref(false);

async function fetchSchedules() {
  if (!props.classRow) return;
  loading.value = true;
  try {
    const res = await getClassSchedule(props.classRow.id);
    schedules.value = res.data.data || [];
    console.log(schedules.value)
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to load schedule");
  } finally {
    loading.value = false;
  }
}

async function fetchSubjectOptions() {
  if (!props.classRow) return;
  subjectsLoading.value = true;
  try {
    const res = await getClassAvailableSubjects(props.classRow.id);
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

async function handleAdd() {
  if (!formRef.value || !props.classRow) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    await createClassSchedule(props.classRow.id, { ...form });
    notify.success("បន្ថែមកាលវិភាគជោគជ័យ");
    Object.assign(form, defaultForm());
    formRef.value.clearValidate();
    await fetchSchedules();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចបន្ថែមកាលវិភាគបានទេ");
  } finally {
    submitting.value = false;
  }
}

async function handleToggle(row) {
  togglingId.value = row.id;
  try {
    await toggleClassSchedule(row.id);
    notify.success("ប្តូរស្ថានភាពជោគជ័យ");
    await fetchSchedules();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចប្តូរស្ថានភាពបានទេ");
  } finally {
    togglingId.value = null;
  }
}

// Load fresh data every time the dialog is opened for a (possibly new) class
watch(
  () => [props.modelValue, props.classRow?.id],
  ([open]) => {
    if (open && props.classRow) {
      Object.assign(form, defaultForm());
      fetchSchedules();
      fetchSubjectOptions();
    }
  },
);
</script>

<template>
  <AppDialog
    v-model="visible"
    :title="classRow ? `កាលវិភាគសម្រាប់ ${classRow.name}` : 'កាលវិភាគ'"
    width="70%"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="add-schedule-form"
    >
      <div class="form-row">
        <el-form-item label="មុខវិជ្ជា" prop="subject_id" class="grow">
          <AppSelect
            v-model="form.subject_id"
            :options="subjectOptions"
            :loading="subjectsLoading"
            placeholder="ជ្រើសរើសមុខវិជ្ជា"
            size="large"
            filterable
            clearable
          />
        </el-form-item>
        <el-form-item label="ថ្ងៃ" prop="day_of_week">
          <el-select v-model="form.day_of_week" placeholder="ជ្រើសរើសថ្ងៃ" size="large">
            <el-option v-for="d in DAYS" :key="d.value" :label="d.label" :value="d.value" />
          </el-select>
        </el-form-item>
      </div>

      <AppButton type="primary" :loading="submitting" @click="handleAdd">
        បន្ថែមកាលវិភាគ
      </AppButton>
    </el-form>

    <AppTable
      class="schedule-table"
      show-index
      :data="schedules"
      :loading="loading"
      :show-pagination="false"
      :columns="[
        { label: 'លេខកូដ', prop: 'subject_code', minWidth: 100 },
        { label: 'មុខវិជ្ជា', prop: 'subject_name_kh', minWidth: 160 },
        { label: 'ថ្ងៃ', slot: 'day', width: 120 },
        { label: 'ស្ថានភាព', slot: 'status', width: 110 },
      ]"
    >
      <template #day="{ row }">
        <el-text>{{ dayLabel(row.day_of_week) }}</el-text>
      </template>

      <template #status="{ row }">
        <el-switch
          :model-value="row.is_active"
          :loading="togglingId === row.id"
          @change="handleToggle(row)"
        />
      </template>
    </AppTable>

    <template #footer>
      <AppButton @click="visible = false">បិទ</AppButton>
    </template>
  </AppDialog>
</template>

<style scoped>
.add-schedule-form {
  margin-bottom: 16px;
}
.form-row {
  display: flex;
  gap: 16px;
}
.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
.schedule-table {
  margin-top: 8px;
}
</style>