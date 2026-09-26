<template>
  <AppDialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)"
    :title="`បញ្ចូលសិស្សតាម Excel — ${className} ${programmeName} ${generationName} ឆ្នាំ${yearName} ឆមាស${semesterName} ជំនាញ${majorName}`"
    width="60%" @closed="resetImport">
    <el-upload drag :auto-upload="false" :show-file-list="false" accept=".xlsx,.xls" :on-change="handleFilePicked">
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        អូសឯកសារមកទីនេះ ឬ <em>ចុចដើម្បីជ្រើសរើសឯកសារ</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          Column ដែលត្រូវការ (ជួរទី ១): <b>name_kh, name_en, gender, code</b> — gender អាចជា 1/2 ឬ ប្រុស/ស្រី
        </div>
      </template>
    </el-upload>
    <AppTable v-if="previewRows.length" :data="previewRows" show-index :show-pagination="false" :columns="[
      { prop: 'name_kh', label: 'ឈ្មោះខ្មែរ', minwidth: 80 },
      { prop: 'name_en', label: 'ឈ្មោះឡាតាំង', minwidth: 80 },
      { prop: 'code', label: 'អត្តលេខ', minwidth: 80 },
      { slot: 'gender', label: 'ភេទ', minwidth: 80 },
      { slot: '_valid', label: 'ស្ថានភាព', minwidth: 80 },
    ]">
      <template #gender="{ row }">
        <el-text>{{ row.gender === 1 ? "ប្រុស" : row.gender === 2 ? "ស្រី" : "?" }}</el-text>
      </template>
      <template #_valid="{ row }">
        <el-tag :type="row._valid ? 'success' : 'danger'" size="small">
          {{ row._valid ? "OK" : "ខ្វះទិន្នន័យ" }}
        </el-tag>
      </template>
    </AppTable>



    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton @click="handleImportSubmit" type="primary" :loading="importing" size="large" :block="false"
        :disabled="!previewRows.length">
        បញ្ចូលសិស្ស ({{ validRowCount }}) នាក់
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import * as XLSX from "xlsx";
import { ref, computed } from "vue";
import { registerUsersExcel } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import { useNotification } from "../composables/useNotification.js";
import AppTable from "./AppTable.vue";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classId: { type: [Number, String], default: null },
  className: { type: String, default: "" },
  generationName: { type: String, default: "" },
  programmeName: { type: String, default: "" },
  yearName: { type: String, default: "" },
  semesterName: { type: String, default: "" },
  majorName: { type: String, default: "" }
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const importing = ref(false);
const previewRows = ref([]);
const pickedFile = ref(null);

const validRowCount = computed(() => previewRows.value.filter((r) => r._valid).length);

function resetImport() {
  previewRows.value = [];
  pickedFile.value = null;
}



function normalizeGender(val) {
  if (val === 1 || val === 2) return val;
  const s = String(val || "").trim();
  if (s === "1" || s === "ប្រុស" || /^m(ale)?$/i.test(s)) return 1;
  if (s === "2" || s === "ស្រី" || /^f(emale)?$/i.test(s)) return 2;
  return null;
}

// Client-side parse is only for an instant preview so the user can catch
// mistakes before uploading. The authoritative parse happens server-side.
function handleFilePicked(uploadFile) {
  const file = uploadFile.raw;
  pickedFile.value = file;

  const reader = new FileReader();
  reader.onload = (e) => {
    const wb = XLSX.read(e.target.result, { type: "array" });
    const sheet = wb.Sheets[wb.SheetNames[0]];
    const rows = XLSX.utils.sheet_to_json(sheet, { defval: "" });

    previewRows.value = rows.map((r) => {
      const name_kh = String(r.name_kh ?? r["ឈ្មោះខ្មែរ"] ?? "").trim();
      const name_en = String(r.name_en ?? r["ឈ្មោះឡាតាំង"] ?? "").trim();
      const code = String(r.code ?? r["អត្តលេខ"] ?? "").trim();
      const gender = normalizeGender(r.gender ?? r["ភេទ"]);
      return {
        name_kh,
        name_en,
        code,
        gender,
        _valid: !!(name_kh && name_en && code && gender),
      };
    });

    if (!previewRows.value.length) {
      notify.error("រកមិនឃើញទិន្នន័យក្នុងឯកសារនេះទេ");
    }
  };
  reader.readAsArrayBuffer(file);
}

async function handleImportSubmit() {
  if (!pickedFile.value) return;
  if (validRowCount.value === 0) {
    notify.error("គ្មានជួរដេលត្រឹមត្រូវសម្រាប់នាំចូលទេ");
    return;
  }

  importing.value = true;
  try {
    const formData = new FormData();
    formData.append("file", pickedFile.value);
    formData.append("class_id", props.classId);
    const res = await registerUsersExcel(formData);
    const created = res.data?.created ?? validRowCount.value;
    notify.success(`នាំចូលសិស្សបានជោគជ័យ (${created})`);
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to import file");
  } finally {
    importing.value = false;
  }
}

defineExpose({ resetImport });
</script>
