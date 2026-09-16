<template>
  <div>
    <AppFilterBar :fields="[
      { slot: 'name', span: 10 },
      { slot: 'role', span: 6 },
      { slot: 'add', span: 4 },
    ]" >
      <template #name>
        <AppInput v-model="filters.name" placeholder="ស្វែងរក" prefix-icon="Search" clearable @input="fetchUsers" />
      </template>

      <template #role>
        <el-select v-model="filters.role_id" placeholder="តួនាទី" clearable style="width: 100%" @change="fetchUsers"
          size="large">
          <el-option v-for="role in roles" :key="role.id" :label="role.display_name" :value="role.id" />
        </el-select>
      </template>

      <template #add>
        <AppButton v-if="candadd" type="primary" @click="openCreate">
          បន្ថែមសិស្ស
        </AppButton>
      </template>

      <template #actions>
        <AppButton v-if="candadd" type="success" icon="Upload" @click="openImport">
          នាំចូលពី Excel
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card>
      <AppTable :data="users" :loading="loading" show-index v-model:current-page="page" v-model:page-size="pageSize"
        :total="total" @page-change="fetchUsers" :columns="[
          { prop: 'name_kh', label: 'ឈ្មោះខ្មែរ', minWidth: 130 },
          { prop: 'name_en', label: 'ឈ្មោះឡាតាំង', minWidth: 130 },
          { prop: 'code', label: 'កូដ', minWidth: 110 },
          { prop: 'gender_string', label: 'ភេទ', minWidth: 90 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 110 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]" actionsWidth="140">
        <template #status="{ row }">
          <el-text :style="{ color: row.is_active ? 'black' : 'red' }">
            {{ row.is_active ? "កំពុងធ្វើការ" : "ឈប់ធ្វើការ" }}
          </el-text>
        </template>

        <template #actions="{ row }" v-if="canedit">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton size="small" icon="Edit" type="warning" circle @click="openEdit(row)" />
          </el-tooltip>
          <el-tooltip content="បិទ" placement="top">
            <AppButton size="small" :icon="row.is_active ? 'CircleClose' : 'CircleCheck'"
              :type="row.is_active ? 'danger' : 'success'" circle @click="toggleStatus(row)" />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <!-- Create / Edit dialog -->
    <AppDialog v-model="createDialog" :title="isEdit ? 'កែប្រែសិស្ស' : 'បន្ថែមសិស្ស'" width="600px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-position="top">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <AppInput label="ឈ្មោះខ្មែរ" prop="name_kh" clearable v-model="createForm.name_kh" />
          </el-col>
          <el-col :xs="24" :sm="12">
            <AppInput label="ឈ្មោះឡាតាំង" prop="name_en" clearable v-model="createForm.name_en" />
          </el-col>
          <el-col :xs="24" :sm="12">
            <AppInput label="កូដ" prop="code" clearable v-model="createForm.code" />
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="ភេទ" prop="gender">
              <el-select v-model="createForm.gender" size="large" style="width: 100%">
                <el-option label="ប្រុស" :value="1" />
                <el-option label="ស្រី" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <!--
            Role is intentionally NOT editable on create: the backend
            Register service hardcodes RoleID = 5 for every bulk-created user.
            It only shows up here on edit, where a dedicated update endpoint
            can change it.
          -->
          <el-col :xs="24" :sm="12" v-if="isEdit">
            <el-form-item label="តួនាទី" prop="role_id">
              <el-select v-model="createForm.role_id" placeholder="ជ្រើសតួនាទី" clearable size="large"
                style="width: 100%">
                <el-option v-for="role in roles" :key="role.id" :label="role.display_name" :value="role.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <AppButton @click="createDialog = false" size="large" :block="false" type="warning">
          បោះបង់
        </AppButton>
        <AppButton @click="handleSave" type="primary" :loading="saving" size="large" :block="false">
          {{ isEdit ? "កែប្រែ" : "បង្កេីត" }}
        </AppButton>
      </template>
    </AppDialog>

    <!-- Excel import dialog -->
    <AppDialog v-model="importDialog" title="នាំចូលសិស្សពី Excel" width="900px" @closed="resetImport">
      <el-upload drag :auto-upload="false" :show-file-list="false" accept=".xlsx,.xls" :on-change="handleFilePicked">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          អូសឯកសារមកទីនេះ ឬ <em>ចុចដើម្បីជ្រើសរើសឯកសារ</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            Column ដែលត្រូវការ (ជួរទី ១): <b>name_kh, name_en, gender, code</b> —
            gender អាចជា 1/2 ឬ ប្រុស/ស្រី
          </div>
        </template>
      </el-upload>

      <el-table v-if="previewRows.length" :data="previewRows" size="small" stripe border style="margin-top: 16px"
        max-height="360">
        <el-table-column type="index" width="50" label="#" />
        <el-table-column prop="name_kh" label="ឈ្មោះខ្មែរ" />
        <el-table-column prop="name_en" label="ឈ្មោះឡាតាំង" />
        <el-table-column prop="code" label="កូដ" />
        <el-table-column label="ភេទ" width="90">
          <template #default="{ row }">{{ row.gender === 1 ? "ប្រុស" : row.gender === 2 ? "ស្រី" : "?" }}</template>
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
        <AppButton @click="importDialog = false" size="large" :block="false" type="warning">
          បោះបង់
        </AppButton>
        <AppButton @click="handleImportSubmit" type="primary" :loading="importing" size="large" :block="false"
          :disabled="!previewRows.length">
          នាំចូល ({{ validRowCount }})
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import { UploadFilled } from "@element-plus/icons-vue";
import { debounce } from "lodash-es";
import * as XLSX from "xlsx";
import {
  getUsers,
  updateUser,
  toggleUserStatus,
  getrole,
  registerUsers,
  registerUsersExcel,
} from "../api/services";
import { watch } from "vue";
import { useUserDataStore } from "../stores/user_data";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppInput from "../../components/AppInput.vue";

const users = ref([]);
const roles = ref([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const filters = reactive({ name: "", role_id: null });

const createDialog = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const createFormRef = ref();

const importDialog = ref(false);
const importing = ref(false);
const previewRows = ref([]);
const pickedFile = ref(null);

const userDataStore = useUserDataStore();
const canedit = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.user"),
);
const candadd = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.user"),
);

const createForm = reactive({
  name_kh: "",
  name_en: "",
  code: "",
  gender: null,
  role_id: null,
});

const createRules = {
  name_kh: [{ required: true, message: "សូមបញ្ចូលឈ្មោះខ្មែរ", trigger: "blur" }],
  name_en: [{ required: true, message: "សូមបញ្ចូលឈ្មោះឡាតាំង", trigger: "blur" }],
  code: [{ required: true, message: "សូមបញ្ចូលកូដ", trigger: "blur" }],
  gender: [{ required: true, message: "សូមជ្រើសភេទ", trigger: "change" }],
};

const validRowCount = computed(
  () => previewRows.value.filter((r) => r._valid).length,
);

const debouncedFetch = debounce(() => {
  page.value = 1;
  fetchUsers();
}, 500);
watch(() => filters.name, debouncedFetch);

async function fetchUsers() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.role_id) params.role_id = filters.role_id;
    const res = await getUsers(params);
    users.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

async function fetchRole() {
  try {
    const res = await getrole();
    roles.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load roles");
  }
}

function openCreate() {
  isEdit.value = false;
  editId.value = null;
  Object.assign(createForm, { name_kh: "", name_en: "", code: "", gender: null, role_id: null });
  createDialog.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(createForm, {
    name_kh: row.name_kh,
    name_en: row.name_en,
    code: row.code,
    gender: row.gender,
    role_id: row.role_id,
  });
  createDialog.value = true;
}

async function handleSave() {
  await createFormRef.value.validate();
  saving.value = true;
  try {
    if (isEdit.value) {
      // Editing an existing user keeps the dedicated update endpoint,
      // since role changes are allowed there.
      await updateUser(editId.value, {
        name_kh: createForm.name_kh,
        name_en: createForm.name_en,
        code: createForm.code,
        gender: createForm.gender,
        role_id: createForm.role_id,
      });
      ElMessage.success("កែប្រែសិស្សបានជោគជ័យ");
    } else {
      // Creating goes through the same Register endpoint the bulk-import
      // uses, just with a single-element array — this matches
      // RegisterRequest{ users: []UserInput } on the backend exactly.
      await registerUsers({
        users: [
          {
            name_kh: createForm.name_kh,
            name_en: createForm.name_en,
            code: createForm.code,
            gender: createForm.gender,
          },
        ],
      });
      ElMessage.success("បង្កេីតសិស្សបានជោគជ័យ");
    }
    createDialog.value = false;
    fetchUsers();
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.response?.data?.error || "Failed to save");
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row) {
  await ElMessageBox.confirm(
    `${row.is_active ? "បិទ" : "បេីក"} ${row.name_kh}?`,
    "សូមបញ្ជាក់",
    { type: "warning" },
  );
  try {
    await toggleUserStatus(row.id);
    ElMessage.success("កែប្រែស្ថានភាពបានជោគជ័យ");
    fetchUsers();
  } catch (e) {
    ElNotification.error({ title: "Error", message: e.response?.data?.error, offset: 100 });
  }
}

/* ---------------- Excel import ---------------- */

function openImport() {
  resetImport();
  importDialog.value = true;
}

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
// mistakes before uploading. The authoritative parse happens on the
// server (see registerUsersExcel / the Go excelize handler) so there is
// a single source of truth for validation.
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
      ElMessage.warning("រកមិនឃើញទិន្នន័យក្នុងឯកសារនេះទេ");
    }
  };
  reader.readAsArrayBuffer(file);
}

async function handleImportSubmit() {
  if (!pickedFile.value) return;
  if (validRowCount.value === 0) {
    ElMessage.warning("គ្មានជួរដេលត្រឹមត្រូវសម្រាប់នាំចូលទេ");
    return;
  }

  importing.value = true;
  try {
    const formData = new FormData();
    formData.append("file", pickedFile.value);
    const res = await registerUsersExcel(formData);
    const created = res.data?.created ?? validRowCount.value;
    ElMessage.success(`នាំចូលសិស្សបានជោគជ័យ (${created})`);
    importDialog.value = false;
    fetchUsers();
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.response?.data?.error || "Failed to import file");
  } finally {
    importing.value = false;
  }
}

onMounted(() => {
  fetchUsers();
  fetchRole();
});
</script>

<style scoped>
.name-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>