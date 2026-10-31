<script setup>
import {
  ref,
  reactive,
  computed,
  watch,
  nextTick,
  onMounted,
  onBeforeUnmount,
} from "vue";
import { ElMessageBox,ElMessage } from "element-plus";
import { useNotification } from "../composables/useNotification.js";
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
  deleteStudentExam,
  getStaft,
} from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppDatePicker from "./AppDatePicker.vue";
import AppTable from "./AppTable.vue";
import PrintableTable from "./PrintableTable.vue";
import PrintExamPlatform from "./PrintExamPlatform.vue";

const props = defineProps({
  modelValue: Boolean,
  classRow: Object,
});

const emit = defineEmits(["update:modelValue"]);

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value),
});



const errorMessage = (error, fallback) =>
  error.response?.data?.error || fallback;

const fmtDate = (value) => (value ? String(value).slice(0, 10) : "-");

const fmtTime = (value) => {
  if (!value) return "";
  const text = String(value);
  const time = text.includes("T") ? text.split("T")[1] : text;
  return time.slice(0, 5);
};

const toTime = (value) => (value ? fmtTime(value) + ":00" : "");

const toKh = (value) =>
  String(value ?? "").replace(/\d/g, (digit) => "០១២៣៤៥៦៧៨៩"[digit]);

const khMonths = [
  "មករា",
  "កុម្ភៈ",
  "មីនា",
  "មេសា",
  "ឧសភា",
  "មិថុនា",
  "កក្កដា",
  "សីហា",
  "កញ្ញា",
  "តុលា",
  "វិច្ឆិកា",
  "ធ្នូ",
];

const khDate = (value) => {
  if (!value) return "";
  const [year, month, day] = String(value).slice(0, 10).split("-").map(Number);
  return `ថ្ងៃទី${toKh(day)} ខែ${khMonths[month - 1]} ឆ្នាំ${toKh(year)}`;
};

async function confirmDelete(message) {
  try {
    await ElMessageBox.confirm(message, "បញ្ជាក់", { type: "warning" });
    return true;
  } catch {
    return false; // the user clicked cancel
  }
}

const loading = ref(false);
const exams = ref([]);
const subjectOptions = ref([]);

const staftOptions = ref([]);

const typeOptions = [
  { label: "ពាក់កណ្តាលឆមាស", value: "MIDTERM" },
  { label: "បញ្ចប់ឆមាស", value: "FINAL" },
  { label: "ប្រឡងចូល", value: "ENTRANCE" },
  { label: "ប្រឡងចេញ", value: "EXIT" },
  { label: "ប្រឡងសង", value: "MAKEUP" },
];

const typeLabels = Object.fromEntries(
  typeOptions.map((type) => [type.value, type.label]),
);

async function fetchExams() {
  if (!props.classRow) return;

  loading.value = true;
  try {
    const res = await getExamPlatforms(props.classRow.id);
    exams.value = res.data.data || [];
  } catch (e) {
    ElMessage.error(errorMessage(e, "Failed to load exams"));
  } finally {
    loading.value = false;
  }
}

async function fetchSubjects() {
  if (!props.classRow) return;

  try {
    const res = await getClassAvailableSubjects(props.classRow.id);
    subjectOptions.value = (res.data.data || []).map((s) => ({
      label: `${s.code} — ${s.name_kh}`,
      value: s.id,
    }));
  } catch (e) {
    ElMessage.error(errorMessage(e, "Failed to load subjects"));
  }
}

async function fetchStaft() {
  try {
    const res = await getStaft({ page: 1, page_size: 100, role_id: 3 });
    staftOptions.value = (res.data.data || []).map((t) => ({
      label: `${t.name_kh} — ${t.name_en}`,
      value: t.id,
    }));
  } catch (e) {
    ElMessage.error(errorMessage(e, "Failed to load teachers"));
  }
}

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

const formRef = ref();
const submitting = ref(false);

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

const form = reactive(defaultForm());

// Validation rules (the user must fill these fields)
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

// Hide subjects already chosen in other rows (so no duplicates)
function availableSubjects(currentKey) {
  const usedIds = form.details
    .filter((d) => d._key !== currentKey && d.subject_id)
    .map((d) => d.subject_id);

  return subjectOptions.value.filter((s) => !usedIds.includes(s.value));
}

function addDetail() {
  form.details.push(newDetail());
}

function removeDetail(index) {
  // always keep at least one row
  if (form.details.length > 1) {
    form.details.splice(index, 1);
  }
}

async function handleSubmit() {
  if (!props.classRow) return;

  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    await createExamPlatform({
      class_id: props.classRow.id,
      type: form.type,
      start_use: form.start_use,
      expire_date: form.expire_date,
      show_result: form.show_result || null,
      check_by: form.check_by,
      // remove "_key" because the server does not need it
      details: form.details.map(({ _key, ...detail }) => detail),
    });

    ElMessage.success("បន្ថែមការប្រឡងជោគជ័យ");
    Object.assign(form, defaultForm()); // reset the form
    formRef.value.clearValidate();
    await fetchExams();
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចបន្ថែមការប្រឡងបានទេ"));
  } finally {
    submitting.value = false;
  }
}

/* ---------- 7. Edit exam ---------- */
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
    ElMessage.success("កែប្រែជោគជ័យ");
    editVisible.value = false;
    await fetchExams();
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចកែប្រែបានទេ"));
  } finally {
    editSubmitting.value = false;
  }
}

async function handleDelete(row) {
  const confirmed = await confirmDelete("តើអ្នកពិតជាចង់លុបការប្រឡងនេះមែនទេ?");
  if (!confirmed) return;

  try {
    await deleteExamPlatform(row.id);
    ElMessage.success("លុបការប្រឡងជោគជ័យ");
    await fetchExams();
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចលុបការប្រឡងបានទេ"));
  }
}

/* ---------- 8. Edit subject (exam detail) ---------- */
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

// Subjects the user can pick: hide the ones used by the other rows of the same exam
const detailEditSubjects = computed(() => {
  const exam = exams.value.find(
    (e) => e.id === detailEditForm.exam_platform_id,
  );
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
    ElMessage.error("ម៉ោងបញ្ចប់ត្រូវធំជាងម៉ោងចាប់ផ្តើម");
    return;
  }

  detailEditSubmitting.value = true;
  try {
    // "id" goes in the URL and "exam_platform_id" is not needed in the body
    const { id, exam_platform_id, ...payload } = detailEditForm;
    await updateExamPlatformDetail(id, payload);
    ElMessage.success("កែប្រែជោគជ័យ");
    detailEditVisible.value = false;
    await fetchExams();
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចកែប្រែបានទេ"));
  } finally {
    detailEditSubmitting.value = false;
  }
}

/* ---------- 9. Table numbers ---------- */
const tableVisible = ref(false);
const tableSubmitting = ref(false);
const tableLoading = ref(false);
const tableDetail = ref(null); // the subject we are giving table numbers to
const tableRoom = ref("");
const tableStudents = ref([]);

async function openTableNumber(detail) {
  tableDetail.value = detail;
  tableRoom.value = "";

  // Start with every student having no table number
  tableStudents.value = (props.classRow?.students || []).map((s) => ({
    student_id: s.student_id ?? s.id,
    name: s.name_kh || s.name_en || s.name || `#${s.id}`,
    table_number: null,
    gender: s.gender,
    code: s.code,
    id: null,
  }));
  tableVisible.value = true;

  // Then fill in the table numbers that were saved before
  tableLoading.value = true;
  try {
    const res = await getStudentExams(detail.id);
    const saved = res.data.data || [];
    if (saved.length) tableRoom.value = saved[0].room || "";

    const savedByStudentId = new Map(saved.map((r) => [r.student_id, r]));
    tableStudents.value.forEach((student) => {
      const record = savedByStudentId.get(student.student_id);
      if (record) {
        student.table_number = record.table_number;
        student.id = record.id;
      }
    });
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចទាញយកលេខតុបានទេ"));
  } finally {
    tableLoading.value = false;
  }
}

// Give numbers 1, 2, 3 ... to all students
function autoFillTables() {
  tableStudents.value.forEach((student, index) => {
    student.table_number = index + 1;
  });
}

function clearTables() {
  tableStudents.value.forEach((student) => {
    student.table_number = null;
  });
}

async function handleSaveTables() {
  // Check for duplicate table numbers
  const numbers = tableStudents.value
    .map((s) => s.table_number)
    .filter((n) => n != null);

  if (new Set(numbers).size !== numbers.length) {
    ElMessage.error("លេខតុមានស្ទួន");
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
    ElMessage.success("ដាក់លេខតុជោគជ័យ");
    tableVisible.value = false;
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចដាក់លេខតុបានទេ"));
  } finally {
    tableSubmitting.value = false;
  }
}

async function handleDeleteStudentExam(row) {
  // Not saved yet: just clear the number on screen
  if (!row.id) {
    row.table_number = null;
    return;
  }

  const confirmed = await confirmDelete("តើអ្នកពិតជាចង់លុបមែនទេ?");
  if (!confirmed) return;

  try {
    await deleteStudentExam(row.id);
    ElMessage.success("ជោគជ័យ");
    row.id = null;
    row.table_number = null;
  } catch (e) {
    ElMessage.error(errorMessage(e, "មិនអាចលុបបានទេ"));
  }
}

/* ---------- 10. Printing ---------- */
// printMode decides which printable page is shown: null | "students" | "exam"
const printMode = ref(null);

/* 10a. Print the student list (with table numbers) */
const printableRef = ref(null);
const printTitle = ref("");
const secondPrintTitle = ref("");
const thirdPrintTitle = ref("");
const printRows = ref([]);

const studentPrintColumns = [
  { key: "name", label: "គោត្តនាម-នាម" },
  {
    key: "gender",
    label: "ភេទ",
    format: (v) => (v === 1 ? "ប្រុស" : v === 2 ? "ស្រី" : "-"),
  },
  { key: "code", label: "អត្តលេខ" },
  { key: "table_number", label: "លេខតុ", strikeIfEmpty: true },
  { key: "score50", label: "ពិន្ទុ50%" },
  { key: "sign", label: "ហត្ថាលេខា" },
  { key: "other", label: "ផ្សេងៗ" },
];

// Values shown at the top of the printed student list
const printSubject = computed(() => tableDetail.value?.subject_name);
const printExamDate = computed(() => tableDetail.value?.exam_date);
const printExamTime = computed(() =>
  [tableDetail.value?.start_time, tableDetail.value?.end_time]
    .filter(Boolean)
    .join(" — "),
);
const totalStudents = computed(() => tableStudents.value.length);
const totalGirls = computed(
  () => tableStudents.value.filter((s) => s.gender === 2).length,
);

function printStudentList() {
  const c = props.classRow || {};

  printTitle.value = `បញ្ជីរាយនាមបេក្ខជនប្រឡងឆ្នាំទី${c.year} ឆមាសទី${c.semester}`;
  secondPrintTitle.value = `ថ្នាក់${c.programme_name} មហាវិទ្យាល័យ ${c.faculty_name} ជំនាញ ${c.major_name}`;
  thirdPrintTitle.value = `${toKh(c.generation_name)} វគ្គទី${toKh(c.term)} ក្រុមទី${toKh(c.group)} ពេលសិក្សា${c.shift_name} ឆ្នាំសិក្សា${toKh(c.generation_start)}-${toKh(c.generation_end)}`;

  // Sort by table number (students without a number go last)
  printRows.value = [...tableStudents.value].sort(
    (a, b) => (a.table_number ?? 9999) - (b.table_number ?? 9999),
  );

  printMode.value = "students";
  // nextTick = wait until Vue has drawn the printable page, then print
  nextTick(() => printableRef.value?.print());
}

/* 10b. Print the exam schedule */
const printexamRef = ref(null);
const printExamData = reactive({
  title: "",
  secondtitle: "",
  thirdtitle: "",
  rows: [],
  notes: [],
  issuedDate: "",
  approver: "",
  headOfStudy: "",
});

function printExamList(exam) {
  const c = props.classRow || {};

  printExamData.title = `កម្មវិធីប្រឡងប្រចាំឆ្នាំទី${toKh(c.year)} ឆមាសទី${toKh(c.semester)}`;
  printExamData.secondtitle = `ថ្នាក់${c.programme_name} មហាវិទ្យាល័យ ${c.faculty_name} ជំនាញ ${c.major_name}`;
  printExamData.thirdtitle = `${toKh(c.generation_name)} វគ្គ${toKh(c.term)} ក្រុម${toKh(c.group)} (${c.shift_name}) ឆ្នាំសិក្សា ${toKh(c.generation_start)}-${toKh(c.generation_end)}`;

  // Sort subjects by date, then by start time
  printExamData.rows = [...(exam.details || [])]
    .sort((a, b) =>
      `${fmtDate(a.exam_date)} ${fmtTime(a.start_time)}`.localeCompare(
        `${fmtDate(b.exam_date)} ${fmtTime(b.start_time)}`,
      ),
    )
    .map((d) => ({
      subject: d.subject_name,
      date: khDate(d.exam_date),
      time: toKh(`${fmtTime(d.start_time)}-${fmtTime(d.end_time)}`),
      teacher: d.teacher_name || "",
    }));

  // The "show_result" note only appears when that date exists
  printExamData.notes = [
    `កម្មវិធីនេះចាប់ផ្តើមអនុវត្តនៅ${khDate(exam.start_use)}`,
    `កម្មវិធីសិក្សាសម្រាប់ឆ្នាំទី${toKh(c.year)} ឆមាសទី${toKh(c.semester)} នឹងបញ្ចប់នៅ${khDate(exam.expire_date)}`,
    exam.show_result &&
      `ប្រកាសលទ្ធផលប្រឡងឆ្នាំទី${toKh(c.year)} ឆមាសទី${toKh(c.semester)} នៅ${khDate(exam.show_result)}`,
    "ព័ត៌មានផ្សេងៗ សូមទាក់ទងការិយាល័យសិក្សា តាមទូរសព្ទលេខ ០១២ ៨២ ៥២ ៥៦",
  ].filter(Boolean);

  printExamData.issuedDate = khDate(new Date().toISOString());
  printExamData.headOfStudy = exam.check_by_name || "";
  printExamData.approver = "";

  printMode.value = "exam";
  nextTick(() => printexamRef.value?.print());
}

// After the print window closes, hide the printable page again
function resetPrint() {
  printMode.value = null;
}

/* ---------- 11. Lifecycle (things that run automatically) ---------- */
onMounted(() => window.addEventListener("afterprint", resetPrint));
onBeforeUnmount(() => window.removeEventListener("afterprint", resetPrint));

// When the dialog opens: reset the form and load fresh data
watch(
  () => [props.modelValue, props.classRow?.id],
  ([isOpen]) => {
    if (!isOpen || !props.classRow) return;

    Object.assign(form, defaultForm());
    fetchExams();
    fetchSubjects();
    fetchStaft();
  },
);
</script>

<template>
  <AppDialog
    v-model="visible"
    :title="classRow ? `ប្រឡងបញ្ចប់ឆមាស សម្រាប់ ${classRow.name}` : 'ប្រឡង'"
    width="80%"
  >
    <!-- ===== Create exam form ===== -->
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
            :options="staftOptions"
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

      <!-- One card for each subject -->
      <div v-for="(d, i) in form.details" :key="d._key" class="detail-card">
        <AppSelect
          v-model="d.subject_id"
          label="មុខវិជ្ជា"
          placeholder="មុខវិជ្ជា"
          :prop="`details.${i}.subject_id`"
          :rules="detailRules.subject_id"
          :options="availableSubjects(d._key)"
          size="large"
          filterable
        />
        <AppDatePicker
          v-model="d.exam_date"
          label="ថ្ងៃប្រឡង"
          :prop="`details.${i}.exam_date`"
          :rules="detailRules.exam_date"
          size="large"
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
          v-model="d.teacher_id"
          label="គ្រូបង្រៀន"
          placeholder="គ្រូបង្រៀន"
          :prop="`details.${i}.teacher_id`"
          :rules="detailRules.teacher_id"
          :options="staftOptions"
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

    <!-- ===== List of exams (click a row to see its subjects) ===== -->
    <AppTable
      class="exam-table"
      expandable
      row-key="id"
      :data="exams"
      :loading="loading"
      :columns="examColumns"
      :show-pagination="false"
      actions-width="150px"
    >
      <!-- Subjects inside one exam -->
      <template #expand="{ row }">
        <AppTable
          show-index
          :show-pagination="false"
          :data="row.details || []"
          :columns="detailColumns"
          actions-width="120px"
        >
          <template #actions="{ row: detail }">
            <el-tooltip content="កែប្រែ" placement="top">
              <AppButton
                type="warning"
                icon="Edit"
                size="small"
                circle
                @click="openEditDetail(row, detail)"
              />
            </el-tooltip>
            <el-tooltip content="ដាក់លេខតុ" placement="top">
              <AppButton
                type="primary"
                icon="Pointer"
                size="small"
                circle
                @click="openTableNumber(detail)"
              />
            </el-tooltip>
          </template>
        </AppTable>
      </template>

      <!-- Buttons for one exam -->
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton
            type="warning"
            icon="Edit"
            size="small"
            circle
            @click="openEdit(row)"
          />
        </el-tooltip>
        <el-tooltip content="លុប" placement="top">
          <AppButton
            type="danger"
            icon="Delete"
            size="small"
            circle
            @click="handleDelete(row)"
          />
        </el-tooltip>
        <el-tooltip content="បោះពុម្ព" placement="top">
          <AppButton
            type="primary"
            icon="Printer"
            size="small"
            circle
            @click="printExamList(row)"
          />
        </el-tooltip>
      </template>
    </AppTable>

    <!-- ===== Dialog: edit exam ===== -->
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
          v-model="editForm.check_by"
          label="អ្នកត្រួតពិនិត្យ"
          :options="staftOptions"
          size="large"
          filterable
        />
      </el-form>
      <template #footer>
        <AppButton
          type="warning"
          :loading="editSubmitting"
          @click="handleUpdate"
        >
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>

    <!-- ===== Dialog: edit subject ===== -->
    <AppDialog
      v-model="detailEditVisible"
      title="កែប្រែមុខវិជ្ជា"
      width="500px"
    >
      <el-form
        ref="detailEditFormRef"
        :model="detailEditForm"
        label-position="top"
      >
        <AppSelect
          v-model="detailEditForm.subject_id"
          label="មុខវិជ្ជា"
          prop="subject_id"
          :rules="detailRules.subject_id"
          :options="detailEditSubjects"
          size="large"
          filterable
        />
        <AppDatePicker
          v-model="detailEditForm.exam_date"
          label="ថ្ងៃប្រឡង"
          prop="exam_date"
          :rules="detailRules.exam_date"
          size="large"
        />
        <el-form-item
          label="ម៉ោងចាប់ផ្តើម"
          prop="start_time"
          :rules="detailRules.start_time"
        >
          <el-time-picker
            v-model="detailEditForm.start_time"
            value-format="HH:mm:ss"
            format="HH:mm"
            size="large"
          />
        </el-form-item>
        <el-form-item
          label="ម៉ោងបញ្ចប់"
          prop="end_time"
          :rules="detailRules.end_time"
        >
          <el-time-picker
            v-model="detailEditForm.end_time"
            value-format="HH:mm:ss"
            format="HH:mm"
            size="large"
          />
        </el-form-item>
        <AppSelect
          v-model="detailEditForm.teacher_id"
          label="គ្រូបង្រៀន"
          prop="teacher_id"
          :rules="detailRules.teacher_id"
          :options="staftOptions"
          size="large"
          filterable
        />
      </el-form>
      <template #footer>
        <AppButton
          type="warning"
          :loading="detailEditSubmitting"
          @click="handleUpdateDetail"
        >
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>

    <!-- ===== Dialog: table numbers ===== -->
    <AppDialog v-model="tableVisible" title="រៀបលេខតុ" width="40%">
      <div class="table-toolbar">
        <el-input
          v-model.trim="tableRoom"
          placeholder="បន្ទប់ (Room)"
          size="large"
          style="max-width: 220px"
        />
        <AppButton size="small" plain type="primary" @click="autoFillTables">
          លេខរៀងស្វ័យប្រវត្តិ
        </AppButton>
        <AppButton size="small" plain type="danger" @click="clearTables">
          សម្អាត
        </AppButton>
        <AppButton size="small" plain type="success" @click="printStudentList">
          បោះពុម្ភ
        </AppButton>
      </div>

      <el-table
        v-loading="tableLoading"
        :data="tableStudents"
        max-height="450"
        border
      >
        <el-table-column type="index" label="ល.រ" width="60" />
        <el-table-column prop="name" label="គោត្តនាម-នាម" min-width="200" />
        <el-table-column
          prop="gender"
          label="ភេទ"
          width="60"
          :formatter="
            (row) =>
              row.gender === 1 ? 'ប្រុស' : row.gender === 2 ? 'ស្រី' : '-'
          "
        />
        <el-table-column prop="code" label="អត្តលេខ" width="200" />
        <el-table-column label="លេខតុ" width="160" align="center">
          <template #default="{ row }">
            <el-input-number
              v-model="row.table_number"
              :min="1"
              :controls="false"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column label="សកម្មភាព" width="160" align="center">
          <template #default="{ row }">
            <el-tooltip content="លុប" placement="top">
              <AppButton
                size="small"
                circle
                icon="Delete"
                type="danger"
                @click="handleDeleteStudentExam(row)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <AppButton
          type="warning"
          :loading="tableSubmitting"
          @click="handleSaveTables"
        >
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>

    <!-- ===== Printable pages (hidden on screen, shown only when printing) ===== -->
    <!-- Only ONE of them exists at a time, chosen by printMode -->
    <Teleport to="body">
      <div class="print-only-wrapper">
        <PrintableTable
          v-if="printMode === 'students'"
          ref="printableRef"
          :title="printTitle"
          :secondtitle="secondPrintTitle"
          :thirdtitle="thirdPrintTitle"
          :columns="studentPrintColumns"
          :rows="printRows"
          :subject="printSubject"
          :examdate="printExamDate"
          :examtime="printExamTime"
          :room="tableRoom"
          :totalstudent="totalStudents"
          :totalstudentgirld="totalGirls"
        />
        <PrintExamPlatform
          v-if="printMode === 'exam'"
          ref="printexamRef"
          v-bind="printExamData"
        />
      </div>
    </Teleport>
  </AppDialog>
</template>

<style scoped>
/* Printable pages are hidden on screen and visible only when printing */
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

.table-toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 12px;
}

.exam-table {
  margin-top: 8px;
}

/* On small screens, show two columns instead of six */
@media (max-width: 768px) {
  .detail-card {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
