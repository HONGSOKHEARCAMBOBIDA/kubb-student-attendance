<script setup>
import { ref, reactive, onMounted,computed } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import { getleavetype, createleavetype, updateleavetype,getCompany } from "../api/services.js";
import { useUserDataStore } from '../stores/user_data'
const leavetype = ref([]);
const notify = useNotification();
const loading = ref(false);
const submitting = ref(false);
const formRef = ref();
const userDataStore = useUserDataStore()
const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);
const companys = ref([]);
const defaultForm = () => ({
  company_id: null,
  code: "",
  name: "",
  is_active: true,
  is_deduct: false,
});

async function fetchCompany() {
  loading.value = true;
  try {
    const res = await getCompany();
    companys.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

const canAddLeaveType = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.leave.type"),
);
const canEditLeaveType = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.leave.type"),
);
const form = reactive(defaultForm());

const rules = {
  code: [{ required: true, message: "សូមបញ្ចូលលេខកូដ", trigger: "blur" }],
  name: [{ required: true, message: "សូមបញ្ចូលឈ្មោះ", trigger: "blur" }],
  company_id: [{ required: true, message: "សូមជ្រើសរើសក្រុមហ៊ុន", trigger: "change" }],
};

async function fetchLeaveType() {
  loading.value = true;
  try {
    const res = await getleavetype();
    leavetype.value = res.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការទាញយកទិន្នន័យ");
  } finally {
    loading.value = false;
  }
}

function openCreateDialog() {
  isEditMode.value = false;
  editingId.value = null;
  Object.assign(form, defaultForm());
  dialogVisible.value = true;
}

function openEditDialog(row) {
  isEditMode.value = true;
  editingId.value = row.id;
  Object.assign(form, {
    company_id: row.company_id,
    code: row.code,
    name: row.name,
    is_active: row.is_active,
    is_deduct: row.is_deduct,
  });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    if (isEditMode.value) {
      await updateleavetype(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await createleavetype({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchLeaveType();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  fetchCompany();
  fetchLeaveType();
});
</script>

<template>
  <div>
    <div class="page-header">
      <AppButton v-if="canAddLeaveType" type="primary" :block="false" @click="openCreateDialog">
        បន្ថែមប្រភេទច្បាប់
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        :data="leavetype"
        :loading="loading"
        show-index
        :columns="[
          { prop: 'code', label: 'លេខកូដ', minWidth: 110 },
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', width: 200 },
          { label: 'កាត់លុយ', slot: 'is_deduct', width: 150 },
          { label: 'ស្ថានភាព', slot: 'is_active', width: 100 },
        //   { label: '', slot: 'actions', width: 90 },
        ]"
      >
        <template #is_deduct="{ row }">
          <el-tag :type="row.is_deduct ? 'danger' : 'success'">
            {{ row.is_deduct ? "កាត់លុយ" : "មិនកាត់លុយ" }}
          </el-tag>
        </template>

        <template #is_active="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
            {{ row.is_active ? "Active" : "Inactive" }}
          </el-tag>
        </template>

        <template #actions="{ row }" v-if="canEditLeaveType">
          <AppButton
         
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="openEditDialog(row)"
          />
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEditMode ? 'កែប្រែប្រភេទច្បាប់' : 'បន្ថែមប្រភេទច្បាប់'"
      width="640px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
            
            <AppInput
            label="លេខកូដ" prop="code" v-model="form.code"
            ></AppInput>
          <AppInput label="ឈ្មោះ" v-model="form.name" prop="name">

          </AppInput>
        </div>
    <el-select v-model="form.company_id" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%" size="large"
          @change="fetchUsers">
          <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
        </el-select>
        <div class="form-row">
          <el-form-item label="ស្ថានភាព">
            <el-switch v-model="form.is_active" size="large" />
          </el-form-item>
          <el-form-item label="កាត់លុយ">
            <el-switch v-model="form.is_deduct" size="large"/>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <AppButton @click="dialogVisible = false">បោះបង់</AppButton>
        <AppButton type="primary" :loading="submitting" @click="handleSubmit">
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}

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