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
          បន្ថែម
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card>
      <AppTable :data="users" :loading="loading" show-index v-model:current-page="page" v-model:page-size="pageSize"
        :total="total" @page-change="fetchUsers" :columns="[
          { prop: 'name_kh', label: 'ឈ្មោះខ្មែរ', minWidth: 130 },
          { prop: 'name_en', label: 'ឈ្មោះឡាតាំង', minWidth: 130 },
          { prop: 'code', label: 'កូដ', minWidth: 110 },
          { slot: 'gender', label: 'ភេទ', minWidth: 90 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 110 },
        
        ]" actionsWidth="140">
        <template #gender="{row}">
          <el-text>{{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}</el-text>
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
          <el-col >
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
  registerUsersMain,
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
    await registerUsersMain({
      name_kh: createForm.name_kh,
      name_en: createForm.name_en,
      code: createForm.code,
      gender: createForm.gender,
      role_id: createForm.role_id,
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