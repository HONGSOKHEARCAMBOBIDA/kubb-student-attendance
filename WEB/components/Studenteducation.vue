<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :title="
      showForm
        ? isEdit
          ? 'កែប្រែប្រវត្តិសិក្សា'
          : 'បន្ថែមប្រវត្តិសិក្សា'
        : 'ព័ត៍មានអំពីប្រវត្តិសិក្សា'
    "
    width="60%"
    @open="onOpen"
  >
    <!-- ===== LIST VIEW ===== -->
    <div v-if="!showForm">
      <div class="toolbar">
        <el-text>{{ studentRaw?.name_kh }} ({{ studentRaw?.code }})</el-text>
        <AppButton
          type="primary"
          size="large"
          :block="false"
          @click="openCreate"
        >
          បន្ថែម
        </AppButton>
      </div>

      <el-table :data="list" v-loading="loading" border style="width: 100%">
        <el-table-column type="index" label="#" width="55" align="center" />
        <el-table-column prop="level" label="កម្រិតសិក្សា" min-width="120" />
        <el-table-column prop="school_name" label="សាលា" min-width="160" />
        <el-table-column prop="start_date" label="ចាប់ផ្តើម" width="120" />
        <el-table-column prop="end_date" label="បញ្ចប់" width="120" />
        <el-table-column
          prop="cerificate_date"
          label="ថ្ងៃទទួលសញ្ញាបត្រ"
          width="150"
        />
        <el-table-column prop="score" label="ពិន្ទុ" width="90" />
        <el-table-column prop="gpa" label="GPA" width="80" />
        <el-table-column prop="grade" label="និទ្ទេស" width="90" />
        <el-table-column label="សកម្មភាព" width="90" align="center">
          <template #default="{ row }">
            <el-tooltip content="កែប្រែ" placement="top">
              <AppButton
                size="small"
                icon="Edit"
                type="warning"
                circle
                @click="openEdit(row)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="មិនទាន់មានប្រវត្តិសិក្សា" :image-size="80" />
        </template>
      </el-table>
    </div>

    <!-- ===== FORM VIEW ===== -->
    <el-form
      v-else
      :model="form"
      :rules="rules"
      ref="formRef"
      label-position="top"
    >
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <AppInput
            label="កម្រិតសិក្សា"
            prop="level"
            clearable
            v-model="form.level"
          />
        </el-col>
        <el-col :xs="24" :sm="12">
          <AppInput
            label="ឈ្មោះសាលា"
            prop="school_name"
            clearable
            v-model="form.school_name"
          />
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-form-item label="ថ្ងៃចាប់ផ្តើម" prop="start_date">
            <el-date-picker
              v-model="form.start_date"
              type="date"
              value-format="YYYY-MM-DD"
              format="DD/MM/YYYY"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-form-item label="ថ្ងៃបញ្ចប់" prop="end_date">
            <el-date-picker
              v-model="form.end_date"
              type="date"
              value-format="YYYY-MM-DD"
              format="DD/MM/YYYY"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-form-item label="ថ្ងៃទទួលសញ្ញាបត្រ" prop="cerificate_date">
            <el-date-picker
              v-model="form.cerificate_date"
              type="date"
              value-format="YYYY-MM-DD"
              format="DD/MM/YYYY"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8">
          <AppInput
            label="ពិន្ទុ"
            prop="score"
            clearable
            v-model="form.score"
          />
        </el-col>
        <el-col :xs="24" :sm="8">
          <AppInput label="GPA" prop="gpa" clearable v-model="form.gpa" />
        </el-col>
        <el-col :xs="24" :sm="8">
          <AppInput
            label="និទ្ទេស"
            prop="grade"
            clearable
            v-model="form.grade"
          />
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <template v-if="!showForm">
        <AppButton
          @click="emit('update:modelValue', false)"
          size="large"
          :block="false"
          type="warning"
        >
          បិទ
        </AppButton>
      </template>
      <template v-else>
        <AppButton
          @click="showForm = false"
          size="large"
          :block="false"
          type="warning"
        >
          បោះបង់
        </AppButton>
        <AppButton
          @click="handleSave"
          type="primary"
          :loading="saving"
          size="large"
          :block="false"
        >
          {{ isEdit ? "កែប្រែ" : "បង្កេីត" }}
        </AppButton>
      </template>
    </template>
  </AppDialog>
</template>

<script setup>
import { ref, reactive } from "vue";
import { ElMessage } from "element-plus";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";
// TODO: adjust names/path to match your api/services file
import {
  getStudentEducation,
  addStudentEducation,
  updateStudentEducation,
} from "../src/api/services";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  studentRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const list = ref([]);
const loading = ref(false);
const saving = ref(false);
const showForm = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const formRef = ref();

const emptyForm = () => ({
  level: "",
  school_name: "",
  start_date: null,
  end_date: null,
  cerificate_date: null, // spelling matches backend json tag
  score: "",
  gpa: "",
  grade: "",
});
const form = reactive(emptyForm());

const rules = {
  level: [{ required: true, message: "សូមបញ្ចូលកម្រិតសិក្សា", trigger: "blur" }],
  school_name: [
    { required: true, message: "សូមបញ្ចូលឈ្មោះសាលា", trigger: "blur" },
  ],
};

// Backend returns dates through helper.FormatDate. The date picker needs
// YYYY-MM-DD, so normalise DD/MM/YYYY or DD-MM-YYYY (if that's what
// FormatDate returns) back to ISO before editing.


// Send null instead of "" so the *string columns / date columns stay NULL
const nullIfEmpty = (v) => (v === "" || v === undefined ? null : v);

function onOpen() {
  showForm.value = false;
  fetchList();
}

async function fetchList() {
  if (!props.studentRaw?.id) return;
  loading.value = true;
  try {
    const res = await getStudentEducation(props.studentRaw.id);
    list.value = res.data.data || [];
  } catch(e) {
    ElMessage.error(`Failed to load education ${e}`);
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  isEdit.value = false;
  editId.value = null;
  Object.assign(form, emptyForm());
  showForm.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    level: row.level ?? "",
    school_name: row.school_name ?? "",
    start_date: row.start_date,
    end_date: row.end_date,
    cerificate_date: row.cerificate_date,
    score: row.score ?? "",
    gpa: row.gpa ?? "",
    grade: row.grade ?? "",
  });
  showForm.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  try {
    const payload = {
      student_id: props.studentRaw.id,
      level: nullIfEmpty(form.level),
      school_name: nullIfEmpty(form.school_name),
      start_date: nullIfEmpty(form.start_date),
      end_date: nullIfEmpty(form.end_date),
      cerificate_date: nullIfEmpty(form.cerificate_date),
      score: nullIfEmpty(form.score),
      gpa: nullIfEmpty(form.gpa),
      grade: nullIfEmpty(form.grade),
    };
    if (isEdit.value) {
      await updateStudentEducation(editId.value, payload);
      ElMessage.success("កែប្រែបានជោគជ័យ");
    } else {
      await addStudentEducation(payload);
      ElMessage.success("បង្កេីតបានជោគជ័យ");
    }
    showForm.value = false;
    await fetchList();
    emit("saved");
  } catch (e) {
    ElMessage.error(
      e.response?.data?.message || e.response?.data?.error || "Failed to save",
    );
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>