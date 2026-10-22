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
  deleteExamPlatform,
  updateExamPlatform,
  updateExamPlatformDetail,
  addStudentExam,
  getStudentExams,
  deleteStudentExam
} from "../src/api/services.js";
import AppTable from "./AppTable.vue";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import PrintableTable from "./PrintableTable.vue";
import { nextTick } from "vue";
// print
const printableRef = ref(null)
const printTitle = ref('')
const SecondprintTitle = ref('')
const ThirdprintTitle = ref('')
const printRows = ref([])
const studentPrintColumns = [
  { key: 'name', label: 'គោត្តនាម-នាម' },
  {
    key: 'gender',
    label: 'ភេទ',
    format: (v) => (v === 1 ? 'ប្រុស' : v === 2 ? 'ស្រី' : '-'),
  },
  { key: 'code', label: 'អត្តលេខ' },
  { key: 'table_number', label: 'លេខតុ', strikeIfEmpty: true },
  { key: 'score50', label: 'ពិន្ទុ50%' },
  { key: 'sign', label: 'ហត្ថាលេខា' },
  { key: 'other', label: 'ផ្សេងៗ' },
]

function printStudentList() {
  const d = tableDetail.value
  const subject = [d?.subject_code, d?.subject_name].filter(Boolean).join(' — ')
  printTitle.value =
    `បញ្ជីរាយនាមបេក្ខជនប្រឡងឆ្នាំទី${props.classRow?.year} ឆមាសទី${props.classRow?.semester}`
  SecondprintTitle.value = `
    ថ្នាក់${props.classRow?.programme_name} មហាវិទ្យាល័យ ${props.classRow?.faculty_name} ជំនាញ ${props.classRow?.major_name}
  `
  ThirdprintTitle.value = `
    ${props.classRow?.generation_name} វគ្គទី${props.classRow?.term} ក្រុមទី${props.classRow?.group} ពេលសិក្សា${props.classRow?.shift_name} ឆ្នាំសិក្សា${props.classRow?.generation_start}-${props.classRow?.generation_end}
  `
  // sort by table number (empty ones last)
  printRows.value = [...tableStudents.value].sort(
    (a, b) => (a.table_number ?? 9999) - (b.table_number ?? 9999),
  )

  nextTick(() => printableRef.value?.print())
}

// print
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

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      "តើអ្នកពិតជាចង់លុបការប្រឡងនេះមែនទេ?",
      "បញ្ជាក់",
      {
        type: "warning",
      },
    );
  } catch {
    return; // cancelled
  }
  try {
    await deleteExamPlatform(row.id);
    notify.success("លុបការប្រឡងជោគជ័យ");
    await fetchExams();
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចលុបការប្រឡងបានទេ"));
  }
}


async function handleDeleteStudentExam(row) {
  if (!row.id) {
    row.table_number = null;
    return;
  }

  try {
    await ElMessageBox.confirm("តើអ្នកពិតជាចង់លុបមែនទេ?", "បញ្ជាក់", {
      type: "warning",
    });
  } catch {
    return; // cancelled
  }

  try {
    await deleteStudentExam(row.id);
    notify.success("ជោគជ័យ");
    row.id = null;
    row.table_number = null;
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចលុបបានទេ"));
  }
}

const editVisible = ref(false);
const editSubmitting = ref(false);
const editFormRef = ref();
const editForm = reactive({
  id: null,
  start_use: "",
  expire_date: "",
  show_result: "",
  check_by: null,
});

function openEdit(row) {
  Object.assign(editForm, {
    id: row.id,
    start_use: fmtDate(row.start_use),
    expire_date: fmtDate(row.expire_date),
    show_result: row.show_result ? fmtDate(row.show_result) : "",
    check_by: row.check_by,
  });
  editVisible.value = true;
}

async function handleUpdate() {
  const valid = await editFormRef.value.validate().catch(() => false);
  if (!valid) return;

  editSubmitting.value = true;
  try {
    await updateExamPlatform(editForm.id, {
      start_use: editForm.start_use,
      expire_date: editForm.expire_date,
      show_result: editForm.show_result || null,
      check_by: editForm.check_by,
    });
    notify.success("កែប្រែជោគជ័យ");
    editVisible.value = false;
    await fetchExams();
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចកែប្រែបានទេ"));
  } finally {
    editSubmitting.value = false;
  }
}


const toTime = (v) => (v ? fmtTime(v) + ":00" : "");
const detailEditVisible = ref(false);
const detailEditSubmitting = ref(false);
const detailEditFormRef = ref();
const detailEditForm = reactive({
  id: null,
  exam_platform_id: null,
  subject_id: null,
  exam_date: "",
  start_time: "",
  end_time: "",
  teacher_id: null,
  document_url: "",
});
const detailEditSubjects = computed(() => {
  const exam = exams.value.find((e) => e.id === detailEditForm.exam_platform_id);
  const usedIds = (exam?.details || [])
    .filter((d) => d.id !== detailEditForm.id)
    .map((d) => d.subject_id);
  return subjectOptions.value.filter((s) => !usedIds.includes(s.value));
});
function openEditDetail(exam, row) {
  Object.assign(detailEditForm, {
    id: row.id,
    exam_platform_id: exam.id,
    subject_id: row.subject_id,
    exam_date: fmtDate(row.exam_date),
    start_time: toTime(row.start_time),
    end_time: toTime(row.end_time),
    teacher_id: row.teacher_id,
    document_url: row.document_url || "",
  });
  detailEditVisible.value = true;
}
async function handleUpdateDetail() {
  const valid = await detailEditFormRef.value.validate().catch(() => false);
  if (!valid) return;

  if (detailEditForm.end_time <= detailEditForm.start_time) {
    notify.error("ម៉ោងបញ្ចប់ត្រូវធំជាងម៉ោងចាប់ផ្តើម");
    return;
  }

  detailEditSubmitting.value = true;
  try {
    const { id, exam_platform_id, ...payload } = detailEditForm;
    await updateExamPlatformDetail(id, payload);
    notify.success("កែប្រែជោគជ័យ");
    detailEditVisible.value = false;
    await fetchExams();
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចកែប្រែបានទេ"));
  } finally {
    detailEditSubmitting.value = false;
  }
}

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
    console.log(exams.value)
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


const tableVisible = ref(false);
const tableSubmitting = ref(false);
const tableLoading = ref(false);
const tableDetail = ref(null);
const tableRoom = ref("");
const tableStudents = ref([])

const examdate = computed(() =>
  tableDetail.value?.exam_date 
);
const subjects = computed(() =>
  tableDetail.value?.subject_name 
);
const roomprop = computed(() =>
  tableDetail.value?.room 
);

const totalstudent = computed(() => tableStudents.value?.length)
const totalstudentgirl = computed(() =>
  tableStudents.value?.filter(s => s.gender === 2).length || 0
)
const examtime = computed(() =>
  [tableDetail.value?.start_time, tableDetail.value?.end_time]
    .filter(Boolean)
    .join(' — ')
);

async function openTableNumber(detail) {
  tableDetail.value = detail;
  tableRoom.value = "";
  tableStudents.value = (props.classRow?.students || []).map((s) => ({
    student_id: s.student_id ?? s.id,
    name: s.name_kh || s.name_en || s.name || `#${s.id}`,
    table_number: null,
    gender: s.gender,
    code: s.code,
    id: null
  }));
  tableVisible.value = true;

  tableLoading.value = true;
  try {
    const res = await getStudentExams(detail.id);
    const saved = res.data.data || [];
    if (saved.length) tableRoom.value = saved[0].room || "";
    const map = new Map(saved.map((r) => [r.student_id, r]));
tableStudents.value.forEach((s) => {
  const rec = map.get(s.student_id);
  if (rec) {
    s.table_number = rec.table_number;
    s.id = rec.id;
  }
});
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចទាញយកលេខតុបានទេ"));
  } finally {
    tableLoading.value = false;
  }
}

function autoFillTables() {
  tableStudents.value.forEach((s, i) => (s.table_number = i + 1));
}
function clearTables() {
  tableStudents.value.forEach((s) => (s.table_number = null));
}

async function handleSaveTables() {
  // client-side duplicate check
  const nums = tableStudents.value.map((s) => s.table_number).filter((n) => n != null);
  if (new Set(nums).size !== nums.length) {
    notify.error("លេខតុមានស្ទួន");
    return;
  }

  tableSubmitting.value = true;
  try {
    await addStudentExam({
      exam_platform_detail_id: tableDetail.value.id,
      room: tableRoom.value,
      student: tableStudents.value.map((s) => ({
        class_id: props.classRow.id,
        student_id: s.student_id,
        table_number: s.table_number ?? null,
      })),
    });
    notify.success("ដាក់លេខតុជោគជ័យ");
    tableVisible.value = false;
  } catch (e) {
    notify.error(errorMessage(e, "មិនអាចដាក់លេខតុបានទេ"));
  } finally {
    tableSubmitting.value = false;
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

      <AppButton type="warning" :loading="submitting" @click="handleSubmit">
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
      actions-width="120px"
    >
      <template #expand="{ row }">
        <AppTable
          :show-pagination="false"
          show-index
          :data="row.details || []"
          :columns="detailColumns"
          actions-width="120px"
        >
<template #actions="{ row: detail }">
  <el-tooltip content="កែប្រែ" placement="top">
    <AppButton type="warning" icon="Edit" size="small" circle
      @click="openEditDetail(row, detail)" />
  </el-tooltip>
  <el-tooltip content="ដាក់លេខតុ" placement="top">
    <AppButton type="primary" icon="Pointer" size="small" circle
      @click="openTableNumber(detail)" />
  </el-tooltip>
</template>
        </AppTable>
      </template>
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top"
          ><AppButton
            type="warning"
            icon="Edit"
            size="small"
            circle
            @click="openEdit(row)"
          ></AppButton
        ></el-tooltip>
        <el-tooltip content="លុប" placement="top"
          ><AppButton
            type="danger"
            icon="Delete"
            size="small"
            circle
            @click="handleDelete(row)"
          ></AppButton
        ></el-tooltip>
      </template>
    </AppTable>

    <AppDialog v-model="editVisible" title="កែប្រែការប្រឡង" width="500px">
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="ថ្ងៃចាប់ផ្តើមអនុវត្ត" prop="start_use">
          <AppDatePicker v-model="editForm.start_use" size="large" />
        </el-form-item>
        <el-form-item label="ថ្ងៃបញ្ចប់កម្មវិធី" prop="expire_date">
          <AppDatePicker v-model="editForm.expire_date" size="large" />
        </el-form-item>
        <el-form-item label="ថ្ងៃបង្ហាញលទ្ធផល" prop="show_result">
          <AppDatePicker
            v-model="editForm.show_result"
            size="large"
            clearable
          />
        </el-form-item>
        <AppSelect
          label="អ្នកត្រួតពិនិត្យ"
          v-model="editForm.check_by"
          :options="teacherOptions"
          size="large"
          filterable
        />
      </el-form>
      <template #footer>
        <AppButton
          type="warning"
          :loading="editSubmitting"
          @click="handleUpdate"
          >រក្សាទុក</AppButton
        >
      </template>
    </AppDialog>
    <AppDialog v-model="detailEditVisible" title="កែប្រែមុខវិជ្ជា" width="500px">
  <el-form ref="detailEditFormRef" :model="detailEditForm" label-position="top">
    <AppSelect
      label="មុខវិជ្ជា"
      prop="subject_id"
      :rules="detailRules.subject_id"
      v-model="detailEditForm.subject_id"
      :options="detailEditSubjects"
      size="large"
      filterable
    />
    <AppDatePicker
      label="ថ្ងៃប្រឡង"
      prop="exam_date"
      :rules="detailRules.exam_date"
      v-model="detailEditForm.exam_date"
      size="large"
    />
    <el-form-item label="ម៉ោងចាប់ផ្តើម" prop="start_time" :rules="detailRules.start_time">
      <el-time-picker v-model="detailEditForm.start_time"
        value-format="HH:mm:ss" format="HH:mm" size="large" />
    </el-form-item>
    <el-form-item label="ម៉ោងបញ្ចប់" prop="end_time" :rules="detailRules.end_time">
      <el-time-picker v-model="detailEditForm.end_time"
        value-format="HH:mm:ss" format="HH:mm" size="large" />
    </el-form-item>
    <AppSelect
      label="គ្រូបង្រៀន"
      prop="teacher_id"
      :rules="detailRules.teacher_id"
      v-model="detailEditForm.teacher_id"
      :options="teacherOptions"
      size="large"
      filterable
    />
  </el-form>
  <template #footer>
    <AppButton type="warning" :loading="detailEditSubmitting" @click="handleUpdateDetail">
      រក្សាទុក
    </AppButton>
  </template>
</AppDialog>

<AppDialog v-model="tableVisible" title="រៀបលេខតុ" width="40%">
  <div class="table-toolbar">
    <el-input v-model.trim="tableRoom" placeholder="បន្ទប់ (Room)" size="large" style="max-width: 220px" />
    <AppButton size="small" plain type="primary" @click="autoFillTables">លេខរៀងស្វ័យប្រវត្តិ</AppButton>
    <AppButton size="small" plain type="danger" @click="clearTables">សម្អាត</AppButton>
    <AppButton size="small" plain type="success" @click="printStudentList()" >បោះពុម្ភ</AppButton>
  </div>

  <el-table :data="tableStudents" v-loading="tableLoading" max-height="450" border>
    <el-table-column type="index" label="ល.រ" width="60" />
    <el-table-column prop="name" label="គោត្តនាម-នាម" min-width="200" />
    <el-table-column prop="gender" label="ភេទ" width="60"  :formatter="(row) => row.gender === 1 ? 'ប្រុស' : row.gender === 2 ? 'ស្រី' : '-'"/>
    <el-table-column prop="code" label="អត្តលេខ" width="200" />
    <el-table-column label="លេខតុ" width="160" align="center">
      <template #default="{ row }">
        <el-input-number v-model="row.table_number" :min="1" :controls="false" style="width: 100%" />
      </template>
    </el-table-column>
    <el-table-column label="សកម្មភាព" width="160" align="center">
        <template #default="{row}">
          <el-tooltip content="លុប" placement="top"><AppButton @click="handleDeleteStudentExam(row)" size="small" circle icon="Delete" type="danger"></AppButton></el-tooltip>
        </template>
    </el-table-column>
  </el-table>

  <template #footer>
    <AppButton type="warning" :loading="tableSubmitting" @click="handleSaveTables">រក្សាទុក</AppButton>
  </template>
</AppDialog>

<Teleport to="body">
  <div class="print-only-wrapper">
    <PrintableTable
      ref="printableRef"
      :title="printTitle"
      :columns="studentPrintColumns"
      :rows="printRows"
      :subject="subjects"
      :examdate="examdate"
      :examtime="examtime"
      :room="tableRoom"
      :secondtitle="SecondprintTitle"
      :thirdtitle="ThirdprintTitle"
      :totalstudent="totalstudent"
      :totalstudentgirld="totalstudentgirl"
    />
  </div>
</Teleport>
  </AppDialog>
</template>

<style scoped>
.print-only-wrapper {
  display: none;
}
@media print {
  .print-only-wrapper {
    display: block;
  }
}
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

.table-toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }

@media (max-width: 768px) {
  .detail-card {
    grid-template-columns: 1fr 1fr;
  }
}
.exam-table {
  margin-top: 8px;
}
</style>
