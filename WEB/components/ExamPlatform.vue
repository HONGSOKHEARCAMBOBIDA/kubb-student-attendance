<script setup>
import { ref, reactive, computed, watch } from "vue";
import { useNotification } from "../composables/useNotification.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppDatePicker from "./AppDatePicker.vue";
import {
  getUsers,
  getExamPlatforms,
  createExamPlatform,
  getClassAvailableSubjects,
} from "../src/api/services.js";
import AppTable from "./AppTable.vue";

const props = defineProps({
  modelValue: Boolean,
  classRow: Object,
});

const emit = defineEmits(["update:modelValue"]);
const notify = useNotification();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value),
});

const loading = ref(false);
const submitting = ref(false);
const formRef = ref();
const exams = ref([]);
const subjectOptions = ref([]);
const teacherOptions = ref([]);

const errorMessage = (e, fallback) => e.response?.data?.error || fallback;

// ---------- Formatters (backend may return "2026-10-16T00:00:00Z") ----------
const typeOptions = [
  { label: "ពាក់កណ្តាលឆមាស", value: "MIDTERM" },
  { label: "បញ្ចប់ឆមាស", value: "FINAL" },
];
const typeLabels = Object.fromEntries(
  typeOptions.map((t) => [t.value, t.label]),
);

const fmtDate = (v) => (v ? String(v).slice(0, 10) : "-");
const fmtTime = (v) => {
  if (!v) return "";
  const s = String(v);
  // handles "08:00:00" and "0000-01-01T08:00:00Z"
  return (s.includes("T") ? s.split("T")[1] : s).slice(0, 5);
};

// ---------- Table columns (names come straight from the backend) ----------
const examColumns = [
  {
    label: "ប្រភេទ",
    minWidth: 120,
    formatter: (row) => typeLabels[row.type] || row.type,
  },
  {
    prop: "start_use",
    label: "ចាប់ផ្តើមប្រើ",
    minWidth: 120,
    formatter: (row) => fmtDate(row.start_use),
  },
  {
    prop: "expire_date",
    label: "ផុតកំណត់",
    minWidth: 120,
    formatter: (row) => fmtDate(row.expire_date),
  },
  {
    prop: "show_result",
    label: "បង្ហាញលទ្ធផល",
    minWidth: 130,
    formatter: (row) => fmtDate(row.show_result),
  },
  {
    prop: "create_by_name",
    label: "អ្នកបង្កើត",
    minWidth: 140,
    formatter: (row) => row.create_by_name || "-",
  },
  {
    prop: "check_by_name",
    label: "អ្នកត្រួតពិនិត្យ",
    minWidth: 140,
    formatter: (row) => row.check_by_name || "-",
  },
  {
    label: "មុខវិជ្ជា",
    minWidth: 100,
    formatter: (row) => row.details?.length || 0,
  },
];

const detailColumns = [
  {
    label: "មុខវិជ្ជា",
    minWidth: 160,
    formatter: (row) =>
      [row.subject_code, row.subject_name].filter(Boolean).join(" — ") || "-",
  },
  {
    prop: "exam_date",
    label: "ថ្ងៃប្រឡង",
    minWidth: 120,
    formatter: (row) => fmtDate(row.exam_date),
  },
  {
    label: "ម៉ោង",
    minWidth: 140,
    formatter: (row) => `${fmtTime(row.start_time)} – ${fmtTime(row.end_time)}`,
  },
  {
    prop: "teacher_name",
    label: "គ្រូ",
    minWidth: 180,
    formatter: (row) => row.teacher_name || "-",
  },
];

// ---------- Load data ----------
async function fetchExams() {
  if (!props.classRow) return;

  loading.value = true;
  try {
    const res = await getExamPlatforms(props.classRow.id);
    exams.value = res.data.data || [];
  } catch (e) {
    notify.error(errorMessage(e, "Failed to load exams"));
  } finally {
    loading.value = false;
  }
}

// Load unique active subjects (for the form dropdown)
async function fetchSubjects() {
  if (!props.classRow) return;
  try {
    const res = await getClassAvailableSubjects(props.classRow.id);
    subjectOptions.value = (res.data.data || []).map((s) => ({
      label: `${s.code} — ${s.name_kh}`,
      value: s.id,
    }));
  } catch (e) {
    notify.error(errorMessage(e, "Failed to load subjects"));
  }
}

// Load teachers (for the form dropdowns)
async function fetchTeachers() {
  try {
    const res = await getUsers({
      page: 1,
      page_size: 100,
      role_id: 3,
    });

    teacherOptions.value = (res.data.data || []).map((t) => ({
      label: `${t.name_kh} — ${t.name_en}`,
      value: t.id,
    }));
  } catch (e) {
    notify.error(errorMessage(e, "Failed to load teachers"));
  }
}

// ---------- Form ----------
let keySeq = 0;
const newDetail = () => ({
  _key: ++keySeq,
  subject_id: null,
  exam_date: "",
  start_time: "",
  end_time: "",
  teacher_id: null,
  document_url: "",
});

const defaultForm = () => ({
  type: "FINAL",
  start_use: "",
  expire_date: "",
  show_result: "",
  check_by: null,
  details: [newDetail()],
});

const availableSubjects = (currentKey) => {
  const usedIds = form.details
    .filter((d) => d._key !== currentKey && d.subject_id)
    .map((d) => d.subject_id);

  return subjectOptions.value.filter((s) => !usedIds.includes(s.value));
};

const availableSubject = (key) =>
  subjectOptions.value.filter(
    (s) => !form.details.some((d) => d._key !== key && d.subject_id === s.value),
  );

const form = reactive(defaultForm());

const required = (message) => [{ required: true, message, trigger: "change" }];

const rules = {
  type: required("សូមជ្រើសរើសប្រភេទ"),
  start_use: required("សូមជ្រើសរើសថ្ងៃចាប់ផ្តើមប្រើ"),
  expire_date: required("សូមជ្រើសរើសថ្ងៃផុតកំណត់"),
  check_by: required("សូមជ្រើសរើសអ្នកត្រួតពិនិត្យ"),
};

const detailRules = {
  subject_id: required("សូមជ្រើសរើសមុខវិជ្ជា"),
  exam_date: required("សូមជ្រើសរើសថ្ងៃប្រឡង"),
  start_time: required("សូមជ្រើសរើសម៉ោងចាប់ផ្តើម"),
  end_time: required("សូមជ្រើសរើសម៉ោងបញ្ចប់"),
  teacher_id: required("សូមជ្រើសរើសគ្រូ"),
};

const addDetail = () => form.details.push(newDetail());

function removeDetail(index) {
  if (form.details.length > 1) {
    form.details.splice(index, 1);
  }
}

// ---------- Submit ----------
async function handleSubmit() {
  if (!props.classRow) return;

  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;

  try {
    const data = {
      class_id: props.classRow.id,
      type: form.type,
      start_use: form.start_use,
      expire_date: form.expire_date,
      show_result: form.show_result || null,
      check_by: form.check_by,
      details: form.details.map(({ _key, ...d }) => d),
    };

    await createExamPlatform(data);

    notify.success("បន្ថែមការប្រឡងជោគជ័យ");
    Object.assign(form, defaultForm());
    formRef.value.clearValidate();
    await fetchExams();
  } catch (e) {
    console.error("Create exam error:", e);
    notify.error(errorMessage(e, "មិនអាចបន្ថែមការប្រឡងបានទេ"));
  } finally {
    submitting.value = false;
  }
}

// ---------- Open dialog ----------
watch(
  () => [props.modelValue, props.classRow?.id],
  ([open]) => {
    if (!open || !props.classRow) return;

    Object.assign(form, defaultForm());
    fetchExams();
    fetchSubjects();
    fetchTeachers();
  },
);
</script>

<template>
  <AppDialog
    v-model="visible"
    :title="classRow ? `ប្រឡងបញ្ចប់ឆមាស សម្រាប់ ${classRow.name}` : 'ប្រឡង'"
    width="80%"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="exam-form"
    >
      <div class="form-row">
        <el-form-item label="ប្រភេទ" prop="type">
          <AppSelect v-model="form.type" :options="typeOptions" size="large" />
        </el-form-item>
        <el-form-item label="ថ្ងៃចាប់ផ្តើមអនុវត្ត" prop="start_use">
          <AppDatePicker v-model="form.start_use" size="large" clearable />
        </el-form-item>
        <el-form-item label="ថ្ងៃបញ្ចប់កម្មវិធី" prop="expire_date">
          <AppDatePicker v-model="form.expire_date" size="large" clearable />
        </el-form-item>
        <el-form-item label="ថ្ងៃបង្ហាញលទ្ធផល" prop="show_result">
          <AppDatePicker v-model="form.show_result" size="large" clearable />
        </el-form-item>
        <el-form-item label="អ្នកត្រួតពិនិត្យ" prop="check_by">
          <AppSelect
            v-model="form.check_by"
            :options="teacherOptions"
            placeholder="ជ្រើសរើសអ្នកត្រួតពិនិត្យ"
            size="large"
            filterable
            clearable
          />
        </el-form-item>
      </div>

      <div class="details-head">
        <strong>កាលវិភាគប្រឡងតាមមុខវិជ្ជា</strong>
        <AppButton
          size="small"
          type="primary"
          plain
          icon="Plus"
          @click="addDetail"
        >
          បន្ថែមមុខវិជ្ជា
        </AppButton>
      </div>

      <div v-for="(d, i) in form.details" :key="d._key" class="detail-card">
        <AppSelect
          :rules="detailRules.subject_id"
          :prop="`details.${i}.subject_id`"
          label="មុខវិជ្ជា"
          v-model="d.subject_id"
          :options="availableSubjects(d._key)"
          placeholder="មុខវិជ្ជា"
          size="large"
          filterable
        />
        <AppDatePicker
          label="ថ្ងៃប្រឡង"
          :prop="`details.${i}.exam_date`"
          v-model="d.exam_date"
          size="large"
          :rules="detailRules.exam_date"
        />

        <el-form-item
          label="ម៉ោងចាប់ផ្តើម"
          :prop="`details.${i}.start_time`"
          :rules="detailRules.start_time"
        >
          <el-time-picker
            v-model="d.start_time"
            value-format="HH:mm:ss"
            format="HH:mm"
            size="large"
            placeholder="ចាប់ផ្តើម"
          />
        </el-form-item>

        <el-form-item
          label="ម៉ោងបញ្ចប់"
          :prop="`details.${i}.end_time`"
          :rules="detailRules.end_time"
        >
          <el-time-picker
            v-model="d.end_time"
            value-format="HH:mm:ss"
            format="HH:mm"
            size="large"
            placeholder="បញ្ចប់"
          />
        </el-form-item>
        <AppSelect
          :rules="detailRules.teacher_id"
          :prop="`details.${i}.teacher_id`"
          label="គ្រូបង្រៀន"
          v-model="d.teacher_id"
          :options="teacherOptions"
          placeholder="គ្រូបង្រៀន"
          size="large"
          filterable
        />

        <AppButton
          size="small"
          circle
          type="danger"
          icon="Delete"
          :disabled="form.details.length === 1"
          @click="removeDetail(i)"
        />
      </div>

      <AppButton type="primary" :loading="submitting" @click="handleSubmit">
        បង្កើតការប្រឡង
      </AppButton>
    </el-form>

    <AppTable
      class="exam-table"
      expandable
      row-key="id"
      :data="exams"
      :loading="loading"
      :columns="examColumns"
      :show-pagination="false"
    >
      <template #expand="{ row }">
        <AppTable
          :show-pagination="false"
          show-index
          :data="row.details || []"
          :columns="detailColumns"
        />
      </template>
    </AppTable>
  </AppDialog>
</template>

<style scoped>
.exam-form {
  margin-bottom: 16px;
}
.form-row {
  display: flex;
  gap: 6px;
  align-items: flex-start;
}
.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
.details-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 8px 0 12px;
}
.detail-card {
  display: grid;
  grid-template-columns: 2fr 1.2fr 1fr 1fr 1.5fr auto;
  gap: 8px;
  align-items: start;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 8px;
}

.detail-card .el-form-item {
  min-width: 0;
  margin-bottom: 12px;
}

@media (max-width: 768px) {
  .detail-card {
    grid-template-columns: 1fr 1fr;
  }
}
.exam-table {
  margin-top: 8px;
}
</style>