<script setup>
import { ref, reactive, watch, computed } from "vue";
import { useNotification } from "../composables/useNotification.js";
import AppTable from "./AppTable.vue";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppDatePicker from "./AppDatePicker.vue";
import {
  getClassSchedule,
  getClassAvailableSubjects,
  createClassSchedule,
  toggleClassSchedule,
  getUsers,
  updateTeacherRate,
  addClassDetail,
  viewClassDetail,
  updateClassDetail,
} from "../src/api/services.js";
import AppInput from "./AppInput.vue";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
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

const Teachers = ref([]);
const selectedDetail = ref(null);
const formRef = ref();
const defaultForm = () => ({
  subject_id: null,
  day_of_week: null,
  schedule_date: "",
  total_session: 0,
  teacher_id: null,
  hourly_rate: 0,
  effective_date: "",
  end_date: "",
  class_detail_id: null,
  total_student: 0,
  room: "",
  is_mission: false,
});

const examstatus = [
  { label: "មិនទាន់ប្រឡង", value: false },
  { label: "បានប្រឡង", value: true },
];

async function openScheduleForm(detail) {
  selectedDetail.value = detail;
  Object.assign(form, defaultForm(), {
    class_detail_id: detail.id,
    teacher_id: detail.teacher_id ?? null, // prefill, still editable
    subject_id: detail.subject_id ?? null,
  });
  await nextTick(); // formRef exists only after v-if renders
  formRef.value?.clearValidate();
}
function closeScheduleForm() {
  selectedDetail.value = null;
  Object.assign(form, defaultForm());
}

// class detail
const detailFormRef = ref();
const detailSubmitting = ref(false);
const newDetailRow = () => ({ subject_id: null, credit: 0, teacher_id: null });
const detailForm = reactive({ rows: [newDetailRow()] });
function addDetailRow() {
  detailForm.rows.push(newDetailRow());
}
function removeDetailRow(index) {
  if (detailForm.rows.length === 1) return;
  detailForm.rows.splice(index, 1);
}

function subjectOptionsFor(index) {
  const picked = detailForm.rows
    .filter((_, i) => i !== index)
    .map((r) => r.subject_id)
    .filter(Boolean);
  return subjectOptions.value.filter((o) => !picked.includes(o.value));
}

const detailRules = {
  subject_id: [
    { required: true, message: "សូមជ្រើសរើសមុខវិជ្ជា", trigger: "change" },
  ],
  credit: [
    {
      required: true,
      type: "number",
      min: 0,
      message: "សូមបញ្ចូលក្រេឌីត",
      trigger: "change",
    },
  ],
  teacher_id: [
    { required: true, message: "សូមជ្រើសរើសគ្រូបង្រៀន", trigger: "change" },
  ],
};

async function handleAddClassDetail() {
  if (!detailFormRef.value || !props.classRow) return;
  const valid = await detailFormRef.value.validate().catch(() => false);
  if (!valid) return;

  const payload = detailForm.rows.map((r) => ({
    class_id: props.classRow.id,
    subject_id: r.subject_id,
    credit: Number(r.credit),
    teacher_id: r.teacher_id,
  }));

  detailSubmitting.value = true;
  try {
    await addClassDetail(payload);
    notify.success("បន្ថែមមុខវិជ្ជាក្នុងថ្នាក់ជោគជ័យ");
    detailForm.rows = [newDetailRow()];
    detailFormRef.value.clearValidate();
    await fetchSubjectOptions(); // added subjects are no longer "available"
    fetchClassDetail();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចបន្ថែមបានទេ");
  } finally {
    detailSubmitting.value = false;
  }
}

const classDetails = ref([]);

async function fetchClassDetail() {
  if (!props.classRow) return;
  loading.value = true;
  try {
    const res = await viewClassDetail(props.classRow.id);
    classDetails.value = res.data.data || [];
    console.log(classDetails.value);
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to load schedule");
  } finally {
    loading.value = false;
  }
}

const editDetailVisible = ref(false);
const editSubjects = ref([]);
const editDetail = reactive({
  id: null,
  subject_id: null,
  credit: 0,
  teacher_id: null,
  midterm: null,
  final: null,
});

function openEditDetail(row) {
  Object.assign(editDetail, {
    id: row.id,
    subject_id: row.subject_id,
    credit: row.credit,
    teacher_id: row.teacher_id,
    midterm: row.midterm,
    final: row.final,
  });
  // available list excludes the current subject, so add it back
  editSubjects.value = [
    { label: row.subject_name, value: row.subject_id },
    ...subjectOptions.value,
  ];
  editDetailVisible.value = true;
}

async function handleUpdateClassDetail() {
  const { id, ...body } = editDetail;
  try {
    await updateClassDetail(id, {
      ...body,
      class_id: props.classRow.id,
      credit: Number(body.credit),
    });
    notify.success("កែប្រែជោគជ័យ");
    editDetailVisible.value = false;
    await Promise.all([fetchClassDetail(), fetchSubjectOptions()]);
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចកែប្រែបានទេ");
  }
}

// class detail

const form = reactive(defaultForm());
const rules = {
  subject_id: [
    { required: true, message: "សូមជ្រើសរើសមុខវិជ្ជា", trigger: "change" },
  ],
  day_of_week: [
    { required: true, message: "សូមជ្រើសរើសថ្ងៃ", trigger: "change" },
  ],
  schedule_date: [
    { required: true, message: "សូមជ្រើសរើសថ្ងៃទី", trigger: "change" },
  ],
  total_session: [
    {
      required: true,
      message: "សូម Session សរុបដែលសិស្សត្រូវរៀន",
      trigger: "change",
    },
  ],
  teacher_id: [
    { required: true, message: "សូមជ្រើសរើសគ្រូបង្រៀន", trigger: "change" },
  ],
  hourly_rate: [
    { required: true, message: "សូមបញ្ចូលតម្លៃម៉ោងគ្រូ", trigger: "change" },
  ],
  effective_date: [
    { required: true, message: "សូមរើសថ្ងៃមានប្រសិទ្ធភាព", trigger: "change" },
  ],
  end_date: [
    { required: true, message: "សូមរើសថ្ងៃអស់ប្រសិទ្ធភាព", trigger: "change" },
  ],
};
const submitting = ref(false);

async function fetchSchedules() {
  if (!props.classRow) return;
  loading.value = true;
  try {
    const res = await getClassSchedule(props.classRow.id);
    schedules.value = res.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to load schedule");
  } finally {
    loading.value = false;
  }
}

async function fetchTeacher() {
  loading.value = true;
  try {
    const params = { page: 1, page_size: 100 };
    params.role_id = 3;
    const res = await getUsers(params);
    Teachers.value = (res.data.data || []).map((s) => ({
      label: `${s.name_kh} — ${s.name_en}`,
      value: s.id,
    }));
    console.log(Teachers.value);
  } catch (e) {
    ElMessage.error("Failed to load employees");
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
    fetchClassDetail();
    closeScheduleForm();
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

// ---- Edit teacher rate ----
const editVisible = ref(false);
const editRow = ref(null);
const editFormRef = ref();
const editSubmitting = ref(false);
const editForm = reactive({
  user_id: null,
  hourly_rate: 0,
  subject_id: null,
  day_of_week: null,
  schedule_date: "",
  total_session: 0,
  total_student: 0,
  room: "",
  is_mission: false,
});
const editRules = {
  user_id: [
    { required: true, message: "សូមជ្រើសរើសគ្រូបង្រៀន", trigger: "change" },
  ],
  hourly_rate: [
    { required: true, message: "សូមបញ្ចូលតម្លៃម៉ោងគ្រូ", trigger: "change" },
  ],
};

function openEdit(row) {
  editRow.value = row;
  // make sure your GET schedule API returns the teacher's user id (teacher_id / user_id)
  editForm.user_id = row.teacher_id ?? row.user_id ?? null;
  editForm.hourly_rate = Number(row.hourly_rate) || 0;
  editForm.subject_id = row.subject_id ?? null;
  editForm.day_of_week = row.day_of_week ?? null;
  editForm.schedule_date = row.schedule_date ?? null;
  editForm.total_session = row.total_session ?? null;
  editForm.total_student = row.total_student ?? null;
  editForm.is_mission = row.is_mission ?? false;
  editForm.room = row.room ?? "";
  editVisible.value = true;
  editFormRef.value?.clearValidate();
}

async function handleUpdateRate() {
  const valid = await editFormRef.value.validate().catch(() => false);
  if (!valid || !editRow.value) return;

  editSubmitting.value = true;
  try {
    await updateTeacherRate(editRow.value.id, {
      user_id: editForm.user_id,
      hourly_rate: editForm.hourly_rate,
      subject_id: editForm.subject_id,
      day_of_week: editForm.day_of_week,
      schedule_date: editForm.schedule_date,
      total_session: editForm.total_session,
      total_student: editForm.total_student,
      room: editForm.room,
      is_mission: editForm.is_mission,
    });
    notify.success("កែប្រែតម្លៃម៉ោងគ្រូជោគជ័យ");
    editVisible.value = false;
    await fetchSchedules();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចកែប្រែបានទេ");
  } finally {
    editSubmitting.value = false;
  }
}

// Load fresh data every time the dialog is opened for a (possibly new) class
watch(
  () => [props.modelValue, props.classRow?.id],
  ([open]) => {
    if (open && props.classRow) {
      Object.assign(form, defaultForm());
      detailForm.rows = [newDetailRow()];
      fetchClassDetail();
      fetchSchedules();
      fetchSubjectOptions();
      fetchTeacher();
    }
  },
);
</script>

<template>
  <AppDialog
    v-model="visible"
    :title="classRow ? `កាលវិភាគសម្រាប់ ${classRow.name}` : 'កាលវិភាគ'"
    width="78%"
  >
    <el-form
      ref="detailFormRef"
      :model="detailForm"
      label-position="top"
      class="add-schedule-form"
    >
      <el-card shadow="always">
        <template #header>
          <div
            style="
              display: flex;
              justify-content: space-between;
              align-items: center;
            "
          >
            <span>បន្ថែមមុខវិជ្ជាក្នុងថ្នាក់</span>
            <AppButton
              size="small"
              type="primary"
              icon="Plus"
              @click="addDetailRow"
            >
              បន្ថែមជួរ
            </AppButton>
          </div>
        </template>

        <div
          v-for="(row, i) in detailForm.rows"
          :key="i"
          class="form-row detail-row"
        >
          <el-form-item
            label="មុខវិជ្ជា"
            :prop="`rows.${i}.subject_id`"
            class="grow"
          >
            <AppSelect
              v-model="row.subject_id"
              :options="subjectOptionsFor(i)"
              :loading="subjectsLoading"
              placeholder="ជ្រើសរើសមុខវិជ្ជា"
              size="large"
              filterable
              clearable
            />
          </el-form-item>

          <el-form-item
            label="ក្រេឌីត"
            :prop="`rows.${i}.credit`"
            :rules="detailRules.credit"
          >
            <AppInput
              type="number"
              v-model.number="row.credit"
              placeholder="ក្រេឌីត"
              size="large"
            />
          </el-form-item>

          <el-form-item
            label="គ្រូបង្រៀន"
            :prop="`rows.${i}.teacher_id`"
            :rules="detailRules.teacher_id"
          >
            <AppSelect
              v-model="row.teacher_id"
              :options="Teachers"
              placeholder="ជ្រើសរើសគ្រូបង្រៀន"
              size="large"
              filterable
              clearable
            />
          </el-form-item>

          <el-form-item label=" " class="remove-col">
            <AppButton
              circle
              type="danger"
              icon="Delete"
              :disabled="detailForm.rows.length === 1"
              @click="removeDetailRow(i)"
            />
          </el-form-item>
        </div>

        <AppButton
          type="primary"
          :loading="detailSubmitting"
          @click="handleAddClassDetail"
        >
          រក្សាទុក
        </AppButton>

        <AppTable
          class="schedule-table"
          show-index
          :data="classDetails"
          :loading="loading"
          :show-pagination="false"
          actions-width="150px"
          :columns="[
            { label: 'មុខវិជ្ជា', prop: 'subject_name', minWidth: 100 },
            { label: 'សាស្រ្តាចារ្យ', prop: 'teacher_name', minWidth: 100 },
            { label: 'ក្រេឌីតត្រូវបង្រៀន', prop: 'credit', minWidth: 100 },
            { label: 'ម៉ោងសរុបនៅសល់', slot: 'remaining_hours', minWidth: 100 },
            { label: 'ប្រឡង Midterm', slot: 'midterm', minWidth: 100 },
            { label: 'ប្រឡង Final', slot: 'final', minWidth: 100 },
          ]"
        >
          <template #remaining_hours="{ row }">
            <el-text
              size="large"
              style="color: black"
              :tag="row.remaining_hours === 0 ? 'mark' : 'b'"
            >
              {{ row.remaining_hours }} ម៉ោង
            </el-text>
          </template>
          <template #midterm="{ row }">
            <el-text>{{
              row.midterm === false ? "មិនទាន់ប្រឡង" : "បានប្រឡង"
            }}</el-text>
          </template>
          <template #final="{ row }">
            <el-text>{{
              row.final === false ? "មិនទាន់ប្រឡង" : "បានប្រឡង"
            }}</el-text>
          </template>
          <template #actions="{ row }">
            <el-tooltip content="បញ្ចូលកាលវិភាគ" placement="top">
              <AppButton
                circle
                size="small"
                icon="Clock"
                @click="openScheduleForm(row)"
              />
            </el-tooltip>
            <el-tooltip content="កែប្រែ" placement="top">
              <AppButton
                circle
                size="small"
                icon="Edit"
                @click="openEditDetail(row)"
              />
            </el-tooltip>
            <el-tooltip content="លុប" placement="top">
              <AppButton
                circle
                size="small"
                icon="Delete"
                @click="openEditDetail(row)"
              />
            </el-tooltip>
          </template>
        </AppTable>
      </el-card>
    </el-form>

    <el-card v-if="selectedDetail" shadow="always" class="add-schedule-form">
      <template #header>
        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
          "
        >
          <span>
            បញ្ចូលកាលវិភាគ: {{ selectedDetail.subject_name }} —
            {{ selectedDetail.teacher_name }}
          </span>
          <AppButton
            size="small"
            icon="Close"
            circle
            @click="closeScheduleForm"
          />
        </div>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        class="add-schedule-form"
      >
        <div class="form-row">
          <!-- <el-form-item label="មុខវិជ្ជា" prop="subject_id" class="grow">
          <AppSelect
            v-model="form.subject_id"
            :options="subjectOptions"
            :loading="subjectsLoading"
            placeholder="ជ្រើសរើសមុខវិជ្ជា"
            size="large"
            filterable
            clearable
          />
        </el-form-item> -->

          <el-form-item label="ថ្ងៃ" prop="day_of_week">
            <el-select
              v-model="form.day_of_week"
              placeholder="ជ្រើសរើសថ្ងៃ"
              size="large"
            >
              <el-option
                v-for="d in DAYS"
                :key="d.value"
                :label="d.label"
                :value="d.value"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="ថ្ងៃទី" prop="schedule_date">
            <AppDatePicker
              v-model="form.schedule_date"
              placeholder="ជ្រើសរើសថ្ងៃទី"
              size="large"
              clearable
            />
          </el-form-item>

          <el-form-item label="Session សរុប" prop="total_session">
            <AppInput
              type="number"
              v-model.number="form.total_session"
              placeholder="Session សរុប"
              size="large"
              clearable
            />
          </el-form-item>
          <el-form-item label="និស្សិតសរុប">
            <AppInput
              type="number"
              v-model.number="form.total_student"
              placeholder="និស្សិតសរុប"
              size="large"
              clearable
            />
          </el-form-item>
        </div>

        <div class="form-row">
          <!-- <el-form-item label="គ្រូបង្រៀន" prop="teacher_id">
          <AppSelect
            v-model="form.teacher_id"
            :options="Teachers"
            :loading="loading"
            placeholder="ជ្រើសរើសគ្រូបង្រៀន"
            size="large"
            filterable
            clearable
          />
        </el-form-item> -->

          <el-form-item label="តម្លៃម៉ោងគ្រូ" prop="hourly_rate">
            <AppInput
              suffix="$"
              type="number"
              v-model.number="form.hourly_rate"
              placeholder="តម្លៃម៉ោងគ្រូ"
              size="large"
            >
            </AppInput>
          </el-form-item>

          <el-form-item label="ថ្ងៃមានប្រសិទ្ធភាព" prop="effective_date">
            <AppDatePicker
              v-model="form.effective_date"
              placeholder="ថ្ងៃមានប្រសិទ្ធភាព"
              size="large"
              clearable
            />
          </el-form-item>

          <el-form-item label="ថ្ងៃអស់ប្រសិទ្ធភាព" prop="end_date">
            <AppDatePicker
              v-model="form.end_date"
              placeholder="ថ្ងៃអស់ប្រសិទ្ធភាព"
              size="large"
              clearable
            />
          </el-form-item>
          <el-form-item label="បន្ទប់">
            <AppInput v-model="form.room" placeholder="បន្ទប់" size="large">
            </AppInput>
          </el-form-item>
          <el-form-item label="បេសកកម្ម" prop="is_mission">
            <el-switch
              v-model="form.is_mission"
              size="large"
              active-text="បេសកកម្ម"
              inactive-text="ធម្មតា"
            />
          </el-form-item>
        </div>

        <AppButton type="primary" :loading="submitting" @click="handleAdd">
          បន្ថែមកាលវិភាគ
        </AppButton>
      </el-form>
    </el-card>

    <AppTable
      class="schedule-table"
      show-index
      :data="schedules"
      :loading="loading"
      :show-pagination="false"
      actions-width="90px"
      :columns="[
        { label: 'លេខកូដ', prop: 'subject_code', minWidth: 100 },
        { label: 'មុខវិជ្ជា', prop: 'subject_name_kh', minWidth: 160 },
        { label: 'ថ្ងៃ', slot: 'day', width: 120 },
        { label: 'កាលបរិច្ឆេទ', prop: 'schedule_date', width: 120 },
        { label: 'Session សរុប', prop: 'total_session', width: 120 },
        { label: 'និស្សិតសរុប', slot: 'total_student', width: 120 },
        { label: 'បន្ទប់', prop: 'room', width: 120 },
        { label: 'សាស្រ្តាចារ្យ', slot: 'name_kh', minWidth: 120 },
        { label: 'តម្លៃម៉ោង', slot: 'hourly_rate', width: 120 },
        { label: 'ប្រាក់ម៉ោងគ្រូសរុប', slot: 'total_rate', width: 150 },
        { label: 'បង្កើតដោយ', prop: 'created_by', minWidth: 120 },
        { label: 'បាន់ប្រគល់ប្រាក់ម៉ោង', slot: 'status', width: 150 },
        { label: 'ផ្សេងៗ', slot: 'is_mission', width: 110 },
      ]"
    >
      <template #is_mission="{ row }">
        <el-tag v-if="row.is_mission" type="warning">បេសកកម្ម</el-tag>
      </template>
      <template #total_student="{ row }">
        <el-text tag="b" size="large" style="color: black"
          >{{ row.total_student }} នាក់</el-text
        >
      </template>
      <template #name_kh="{ row }">
        <el-text>{{ row.name_kh }} - {{ row.name_en }}</el-text>
      </template>
      <template #day="{ row }">
        <el-text>{{ dayLabel(row.day_of_week) }}</el-text>
      </template>
      <template #hourly_rate="{ row }">
        <el-statistic
          :value="row.hourly_rate"
          :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '20px', color: '#FF0000' }"
        />
      </template>
      <template #total_rate="{ row }">
        <el-statistic
          :value="row.total_session * row.hourly_rate"
          :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '20px', color: '#FF0000' }"
        />
      </template>
      <template #status="{ row }">
        <el-switch
          :model-value="row.is_active"
          :loading="togglingId === row.id"
          @change="handleToggle(row)"
        />
      </template>
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton
            :disabled="row.is_active"
            circle
            size="small"
            type="warning"
            icon="Edit"
            @click="openEdit(row)"
          />
        </el-tooltip>
      </template>
    </AppTable>
  </AppDialog>

  <AppDialog
    v-model="editVisible"
    title="កែប្រែម៉ោងបង្រៀន និងតម្លៃម៉ោងគ្រូ"
    width="420px"
    append-to-body
  >
    <el-form
      ref="editFormRef"
      :model="editForm"
      :rules="editRules"
      label-position="top"
    >
      <AppSelect
        label="គ្រូបង្រៀន"
        v-model="editForm.user_id"
        :options="Teachers"
        placeholder="ជ្រើសរើសគ្រូបង្រៀន"
        size="large"
        filterable
      />

      <AppSelect
        label="មុខវិជ្ជា"
        v-model="editForm.subject_id"
        :options="subjectOptions"
        placeholder="ជ្រើសរើសមុខវិជ្ជា"
        size="large"
        filterable
      />

      <AppSelect
        label="ថ្ងៃ"
        v-model="editForm.day_of_week"
        :options="DAYS"
        placeholder="ជ្រើសរើសថ្ងៃ"
        size="large"
        filterable
      />

      <AppDatePicker
        label="ថ្ងៃបង្រៀន"
        v-model="editForm.schedule_date"
        placeholder="ជ្រើសរើសថ្ងៃបង្រៀន"
        size="large"
        filterable
      >
      </AppDatePicker>
      <AppInput
        label="Session សរុប"
        suffix="$"
        type="number"
        v-model.number="editForm.total_session"
        placeholder="Session សរុប  "
        size="large"
      />
      <AppInput
        label="និស្សិតសរុប"
        suffix="$"
        type="number"
        v-model.number="editForm.total_student"
        placeholder="Session សរុប  "
        size="large"
      />
      <AppInput
        label="តម្លៃម៉ោងគ្រូ"
        suffix="$"
        type="number"
        v-model.number="editForm.hourly_rate"
        placeholder="តម្លៃម៉ោងគ្រូ"
        size="large"
      />

      <AppInput
        label="បន្ទប់"
        v-model="editForm.room"
        placeholder="បន្ទប់"
        size="large"
      />
      <el-form-item label="បេសកកម្ម">
        <el-switch
          v-model="editForm.is_mission"
          size="large"
          active-text="បេសកកម្ម"
          inactive-text="ធម្មតា"
        />
      </el-form-item>

      <div style="display: flex; justify-content: flex-end; gap: 8px">
        <AppButton @click="editVisible = false">បោះបង់</AppButton>
        <AppButton
          type="primary"
          :loading="editSubmitting"
          @click="handleUpdateRate"
        >
          រក្សាទុក
        </AppButton>
      </div>
    </el-form>
  </AppDialog>
  <AppDialog
    v-model="editDetailVisible"
    title="កែប្រែមុខវិជ្ជាក្នុងថ្នាក់"
    width="460px"
    append-to-body
  >
    <AppSelect
      label="មុខវិជ្ជា"
      v-model="editDetail.subject_id"
      :options="editSubjects"
      size="large"
      filterable
    />
    <AppInput
      label="ក្រេឌីត"
      type="number"
      v-model.number="editDetail.credit"
      size="large"
    />
    <AppSelect
      label="គ្រូបង្រៀន"
      v-model="editDetail.teacher_id"
      :options="Teachers"
      size="large"
      filterable
    />
    <AppSelect
      label="Midterm"
      v-model="editDetail.midterm"
      :options="examstatus"
      size="large"
      filterable
    />
    <AppSelect
      label="Final"
      v-model="editDetail.final"
      :options="examstatus"
      size="large"
      filterable
    />
    <div style="display: flex; justify-content: flex-end; gap: 8px">
      <AppButton @click="editDetailVisible = false">បោះបង់</AppButton>
      <AppButton type="primary" @click="handleUpdateClassDetail"
        >រក្សាទុក</AppButton
      >
    </div>
  </AppDialog>
</template>

<style scoped>
.form-row + .form-row {
  margin-top: 4px;
}
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
.detail-row .remove-col {
  flex: 0 0 48px;
}
</style>
