<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :title="
      showForm
        ? isEdit
          ? 'កែប្រែឯកសារ'
          : 'បន្ថែមឯកសារ'
        : 'ព័ត៍មានឯកសារនិស្សិតប្រគល់ជួន'
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
        <el-table-column
          prop="document_type_name"
          label="ប្រភេទឯកសារ"
          min-width="180"
        />
        <el-table-column
          prop="required_qty"
          label="ចំនួនត្រូវការ"
          width="130"
          align="center"
        />
        <el-table-column
          prop="received_qty"
          label="ចំនួនទទួលបាន"
          width="130"
          align="center"
        />
        <el-table-column label="ស្ថានភាព" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="statusOf(row).type">{{ statusOf(row).label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="ចំណាំ" min-width="160" />
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
          <el-empty description="មិនទាន់មានឯកសារ" :image-size="80" />
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
        <el-col :span="24">
          <el-form-item label="ប្រភេទឯកសារ" prop="document_type_id">
            <el-select
              v-model="form.document_type_id"
              placeholder="ជ្រើសរើសប្រភេទឯកសារ"
              filterable
              size="large"
              style="width: 100%"
            >
              <el-option
                v-for="t in availableTypes"
                :key="t.id"
                :label="t.name_kh"
                :value="t.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ចំនួនត្រូវការ" prop="required_qty">
            <el-input-number
              v-model="form.required_qty"
              :min="0"
              :max="255"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ចំនួនទទួលបាន" prop="received_qty">
            <el-input-number
              v-model="form.received_qty"
              :min="0"
              :max="255"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item label="ចំណាំ" prop="remark">
            <el-input
              v-model="form.remark"
              type="textarea"
              :rows="3"
              clearable
            />
          </el-form-item>
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
import { ref, reactive, computed } from "vue";
import { ElMessage } from "element-plus";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
// TODO: adjust names/path to match your api/services file
import {
  getDocumentType,
  getStudentDocuments,
  addStudentDocument,
  updateStudentDocument,
} from "../src/api/services";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  studentRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const list = ref([]);
const documentTypes = ref([]);
const loading = ref(false);
const saving = ref(false);
const showForm = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const formRef = ref();

const emptyForm = () => ({
  document_type_id: null,
  required_qty: 1,
  received_qty: 0,
  remark: "",
});
const form = reactive(emptyForm());

const rules = {
  document_type_id: [
    { required: true, message: "សូមជ្រើសរើសប្រភេទឯកសារ", trigger: "change" },
  ],
};

// (student_id, document_type_id) is unique in the DB, so while adding,
// hide types this student already has. While editing, keep the current one.
const availableTypes = computed(() => {
  const used = new Set(list.value.map((d) => d.document_type_id));
  return documentTypes.value.filter(
    (t) => !used.has(t.id) || (isEdit.value && t.id === form.document_type_id),
  );
});

function statusOf(row) {
  if (row.received_qty >= row.required_qty)
    return { label: "គ្រប់", type: "success" };
  if (row.received_qty > 0) return { label: "មិនទាន់គ្រប់", type: "error" };
  return { label: "មិនទាន់ទទួល", type: "danger" };
}

const nullIfEmpty = (v) => (v === "" || v === undefined ? null : v);

function onOpen() {
  showForm.value = false;
  fetchList();
  fetchTypes();
}

async function fetchTypes() {
  if (documentTypes.value.length) return;
  try {
    const res = await getDocumentType();
    documentTypes.value = res.data.data || [];
  } catch (e) {
    ElMessage.error(`Failed to load document types ${e}`);
  }
}

async function fetchList() {
  if (!props.studentRaw?.id) return;
  loading.value = true;
  try {
    const res = await getStudentDocuments(props.studentRaw.id);
    list.value = res.data.data || [];
  } catch (e) {
    ElMessage.error(`Failed to load documents ${e}`);
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
    document_type_id: row.document_type_id,
    required_qty: row.required_qty ?? 1,
    received_qty: row.received_qty ?? 0,
    remark: row.remark ?? "",
  });
  showForm.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  try {
    const payload = {
      document_type_id: form.document_type_id,
      required_qty: form.required_qty,
      received_qty: form.received_qty,
      remark: nullIfEmpty(form.remark),
    };
    if (isEdit.value) {
      await updateStudentDocument(editId.value, payload);
      ElMessage.success("កែប្រែបានជោគជ័យ");
    } else {
      await addStudentDocument({ ...payload, student_id: props.studentRaw.id });
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