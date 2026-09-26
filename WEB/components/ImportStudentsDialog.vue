<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`បញ្ចូលសិស្សតាម Excel — ${className} ${generationName}`"
    width="900px"
    @closed="resetImport"
  >
    <el-upload
      drag
      :auto-upload="false"
      :show-file-list="false"
      accept=".xlsx,.xls"
      :on-change="handleFilePicked"
    >
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

    <el-table
      v-if="previewRows.length"
      :data="previewRows"
      size="small"
      stripe
      border
      style="margin-top: 16px"
      max-height="360"
    >
      <el-table-column type="index" width="50" label="#" />
      <el-table-column prop="name_kh" label="ឈ្មោះខ្មែរ" />
      <el-table-column prop="name_en" label="ឈ្មោះឡាតាំង" />
      <el-table-column prop="code" label="កូដ" />
      <el-table-column label="ភេទ" width="90">
        <template #default="{ row }">{{
          row.gender === 1 ? "ប្រុស" : row.gender === 2 ? "ស្រី" : "?"
        }}</template>
      </el-table-column>
      <el-table-column label="ស្ថានភាព" width="110">
        <template #default="{ row }">
          <el-tag :type="row._valid ? 'success' : 'danger'" size="small">
            {{ row._valid ? "OK" : "ខ្វះទិន្នន័យ" }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton
        @click="handleImportSubmit"
        type="primary"
        :loading="importing"
        size="large"
        :block="false"
        :disabled="!previewRows.length"
      >
        នាំចូល ({{ validRowCount }})
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

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classId: { type: [Number, String], default: null },
  className: { type: String, default: "" },
  generationName: {type: String, default: ""}
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
      const code = String(r.code ?? r["កូដ"] ?? "").trim();
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
