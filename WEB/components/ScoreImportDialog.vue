<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`បញ្ចូលពិន្ទុតាម Excel — ${classRow?.name || ''}`"
    width="70%"
    @closed="resetForm"
  >
    <el-form label-position="top">
      <AppSelect
        v-model="form.subject_id"
        :options="subjectOptions"
        label="មុខវិជ្ជា"
        placeholder="ជ្រើសរើសមុខវិជ្ជា"
        size="large"
      />
      <div class="form-row">
        <AppInput label="ឆ្នាំ" v-model.number="form.year" type="number" size="large" disabled></AppInput>
         <AppInput label="ឆមាស" v-model.number="form.semester" type="number" size="large" disabled></AppInput>
      </div>
    </el-form>

    <el-upload
      drag
      :auto-upload="false"
      :show-file-list="true"
      :limit="1"
      accept=".xlsx,.xls"
      :on-change="handleFilePicked"
      :on-remove="() => (file = null)"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        អូសឯកសារមកទីនេះ ឬ <em>ចុចដើម្បីជ្រើសរើសឯកសារ</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          Column ដែលត្រូវការ: <b>អត្តលេខ</b> + column មួយៗសម្រាប់ផ្នែកពិន្ទុ
        </div>
      </template>
    </el-upload>

      <AppTable v-if="previewRows.length" :data="previewRows" show-index :show-pagination="false" :columns="[
      { prop: 'name', label: 'គោត្តនាម-នាម', minwidth: 80 },
      { prop: 'gender', label: 'ភេទ', minwidth: 80 },
      { prop: 'code', label: 'អត្តលេខ', minwidth: 80 },
      { prop: 'attendance', label: 'វត្តមាននិស្សិត', minwidth: 80 },
      { prop: 'research', label: 'កិច្ចការស្រាវជ្រាវ', minwidth: 80 },
      { prop: 'midterm', label: 'ប្រឡងពាក់កណ្តាលឆមាស', minwidth: 80 },
      { prop: 'final', label: 'ប្រឡងបញ្ចប់ឆមាស', minwidth: 80 },
    ]">
    </AppTable>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton
        @click="handleSubmit"
        type="primary"
        :loading="importing"
        size="large"
        :block="false"
        :disabled="!file || !form.subject_id"
      >
        នាំចូល
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { getClassAvailableSubjects, importScoreExcel } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import { useNotification } from "../composables/useNotification.js";
import AppInput from "./AppInput.vue";
import AppTable from "./AppTable.vue";
import * as XLSX from "xlsx";
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classRow: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const importing = ref(false);
const file = ref(null);
const subjectOptions = ref([]);
const previewRows = ref([])

const form = reactive({
  subject_id: null,
  year: null,
  semester: null,
});

function resetForm() {
  file.value = null;
  previewRows.value = [];
}

function handleFilePicked(uploadFile) {
  const rawFile = uploadFile.raw;
  file.value = rawFile;
  const reader = new FileReader();
   reader.onload = (e) => {
    const wb = XLSX.read(e.target.result, { type: "array" });
    const sheet = wb.Sheets[wb.SheetNames[0]];
    const rows = XLSX.utils.sheet_to_json(sheet, { defval: "" });

previewRows.value = rows.map((r) => {
  return {
    name: r["គោត្តនាម-នាម"] ?? r.name ?? "",
    gender: r["ភេទ"] ?? r.gender ?? "",
    code: r["អត្តលេខ"] ?? r.code ?? "",
    attendance: r["វត្តមាននិស្សិត"] ?? r.attendance ?? "",
    research: r["កិច្ចការស្រាវជ្រាវ"] ?? r.research ?? "",
    midterm: r["ប្រឡងពាក់កណ្តាលឆមាស"] ?? r.midterm ?? "",
    final: r["ប្រឡងបញ្ចប់ឆមាស"] ?? r.final ?? "",
  };
});

    if (!previewRows.value.length) {
      notify.error("រកមិនឃើញទិន្នន័យក្នុងឯកសារនេះទេ");
    }
  }; 
  reader.readAsArrayBuffer(rawFile);
}

watch(
  () => props.modelValue,
  async (open) => {
    if (!open || !props.classRow) return;
    const row = props.classRow;
    file.value = null;
    form.subject_id = null;
    form.year = row.year;
    form.semester = row.semester;

    try {
      const res = await getClassAvailableSubjects(row.id);
      subjectOptions.value = (res.data.data || []).map((s) => ({
        label: `${s.code} — ${s.name_kh}`,
        value: s.id,
      }));
    } catch (e) {
      notify.error(e.response?.data?.error || "Failed to load subjects");
    }
  },
);

async function handleSubmit() {
  if (!file.value || !form.subject_id) return;
  const row = props.classRow;

  const formData = new FormData();
  formData.append("file", file.value);
  formData.append("class_id", row.id);
  formData.append("subject_id", form.subject_id);
  formData.append("major_id", row.major_id);
  formData.append("generation_id", row.generation_id);
  formData.append("programme_id", row.programme_id);
  formData.append("year", form.year);
  formData.append("semester", form.semester);

  importing.value = true;
  try {
    const res = await importScoreExcel(formData);
    const data = res.data?.data;
    notify.success(`នាំចូល៖ ${data?.imported ?? 0} ជោគជ័យ, ${data?.skipped ?? 0} រំលង`);
    if (data?.errors?.length) console.warn("Import errors:", data.errors);
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "នាំចូលពិន្ទុមិនបានជោគជ័យ");
  } finally {
    importing.value = false;
  }
}
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
