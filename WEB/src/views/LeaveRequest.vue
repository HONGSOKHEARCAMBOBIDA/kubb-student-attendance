<script setup>
import { ref, reactive, onMounted, computed, watch, onUnmounted } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import {
  getleavetype,
  getleavedeductype,
  getuserapprove,
  getleaverequest,
  addleaverequest,
  editleaverequest,
  editstatusleaverequest,
  getCompany,
  deleteleaverequest,
} from "../api/services.js";
import { useUserDataStore } from "../stores/user_data.js";
import AppFilterBar from "../../components/AppFilterBar.vue";
let searchTimer = null;
const notify = useNotification();
const userDataStore = useUserDataStore();

const leaverequest = ref([]);
const leavetype = ref([]);
const leavededucttype = ref([]);
const userapprove = ref([]);
const companys = ref([]);
const loading = ref(false);
const submitting = ref(false);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0,
});

const filter = reactive({
  name: "",
  company_id: "",
  role_id: "",
  status: "",
});

const statusOptions = [
  { value: 1, label: "PENDING" },
  { value: 2, label: "APPROVED" },
  { value: 3, label: "REJECTED" },
  { value: 4, label: "CANCELLED" },
];

const defaultForm = () => ({
  leave_type_id: null,
  start_date: "",
  end_date: "",
  back_to_work_date: "",
  total_day: null,
  deduct_type_id: null,
  reason: "",
  approve_by: null,
});
const form = reactive(defaultForm());

const rules = {
  leave_type_id: [
    { required: true, message: "សូមជ្រើសរើសប្រភេទច្បាប់", trigger: "change" },
  ],
  start_date: [
    { required: true, message: "សូមជ្រើសរើសថ្ងៃចាប់ផ្តើម", trigger: "change" },
  ],
  end_date: [
    { required: true, message: "សូមជ្រើសរើសថ្ងៃបញ្ចប់", trigger: "change" },
  ],
  total_day: [
    { required: true, message: "សូមបញ្ចូលចំនួនថ្ងៃ", trigger: "blur" },
  ],
  reason: [{ required: true, message: "សូមបញ្ចូលមូលហេតុ", trigger: "blur" }],
};

const canEditLeave = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.leave.request"),
);
const canEditStatusLeave = computed(() =>
  userDataStore.permissions?.some(
    (p) => p.name === "edit.status.leave.request",
  ),
);

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

async function fetchLeaveType() {
  try {
    const res = await getleavetype();
    leavetype.value = res.data.data || [];
  } catch {
    notify.error("Failed to load leave types");
  }
}

async function fetchLeaveDeductType() {
  try {
    const res = await getleavedeductype();
    leavededucttype.value = res.data.data || [];
  } catch {
    notify.error("Failed to load deduct types");
  }
}

async function fetchUserApprove() {
  try {
    const res = await getuserapprove();
    userapprove.value = res.data.data || [];
  } catch {
    notify.error("Failed to load approvers");
  }
}

async function fetchLeaveRequest() {
  loading.value = true;
  try {
    const res = await getleaverequest({
      page: pagination.page,
      page_size: pagination.page_size,
      ...filter,
    });
    leaverequest.value = res.data.data || [];
    pagination.total = res.data.pagination?.totalCount || 0;
  } catch {
    notify.error("Failed to load leave requests");
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  pagination.page = 1;
  fetchLeaveRequest();
}

function handlePageChange(page) {
  pagination.page = page;
  fetchLeaveRequest();
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
    leave_type_id: row.leave_type_id,
    start_date: row.start_date,
    end_date: row.end_date,
    back_to_work_date: row.back_to_work_date,
    total_day: row.total_day,
    deduct_type_id: row.deduct_type_id,
    reason: row.reason,
    approve_by: row.approve_by || null,
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
      await editleaverequest(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addleaverequest({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

async function updateStatus(row, status) {
  try {
    await editstatusleaverequest(row.id, { status });
    notify.success("ធ្វើបច្ចុប្បន្នភាពស្ថានភាពជោគជ័យ");
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(
      e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ",
    );
  }
}

async function deleteleave(row) {
  try {
    await deleteleaverequest(row.id);
    notify.success("លុបច្បាប់បានជោគជ័យ");
    await fetchLeaveRequest();
  } catch (e) {
    notify.error(
      e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ",
    );
  }
}

function statusTagType(status) {
  if (status === 2) return "success";
  if (status === 3) return "danger";

  return "warning";
}
watch(
  filter,
  () => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      pagination.page = 1;
      fetchLeaveRequest();
    }, 300); // debounce so typing doesn't fire a request per keystroke
  },
  { deep: true },
);
onMounted(() => {
  fetchCompany();
  fetchLeaveType();
  fetchLeaveDeductType();
  fetchUserApprove();
  fetchLeaveRequest();
});
onUnmounted(() => clearTimeout(searchTimer));
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 7 },
        { slot: 'status', span: 6 },
        { slot: 'company', span: 6 },
        { slot: 'add', span: 4 },
      ]"
      :action-span="4"
    >
      <template #name>
        <AppInput
          v-model="filter.name"
          placeholder="ស្វែង"
          prefix-icon="Search"
          clearable
        >
        </AppInput>
      </template>

      <template #status>
        <el-select
          v-model="filter.status"
          placeholder="ស្ថានភាព"
          clearable
          style="width: 100%"
          size="large"
        >
          <el-option
            v-for="s in statusOptions"
            :key="s.value"
            :label="s.label"
            :value="s.value"
          />
        </el-select>
      </template>
      <template #company>
        <el-select
          v-model="filter.company_id"
          placeholder="ក្រុមហ៑ុន"
          clearable
          style="width: 100%"
          size="large"
          @change="fetchUsers"
        >
          <el-option
            v-for="company in companys"
            :key="company.id"
            :label="company.name"
            :value="company.id"
          />
        </el-select>
      </template>

      <template #add>
        <AppButton type="primary" @click="openCreateDialog">
          បន្ថែមច្បាប់
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        :data="leaverequest"
        :loading="loading"
        :actions-width="240"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        @page-change="fetchLeaveRequest"
        :columns="[
          { label: 'ឈ្មោះបុគ្គលិក', slot: 'user_name', minWidth: 140 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 100 },
          { prop: 'company_name', label: 'ក្រុមហ៊ុន', minWidth: 140 },
          { label: 'ប្រភេទច្បាប់', slot: 'leave_type_name', minWidth: 170 },
          { prop: 'start_date', label: 'ថ្ងៃចាប់ផ្តើម', minWidth: 110 },
          { prop: 'end_date', label: 'ថ្ងៃបញ្ចប់', minWidth: 110 },
          { prop: 'back_to_work_date', label: 'ថ្ងៃមកធ្វេីការ', minWidth: 110 },
          { label: 'ចំនួនថ្ងៃ', slot: 'total_day', width: 150 },
          { label: 'មូលហេតុ', prop: 'reason', minWidth: 150 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
          { prop: 'approve_by_name', label: 'អ្នកអនុម័ត្ត', minWidth: 120 },
          { prop: 'approved_at', label: 'ថ្ងៃអនុម័ត្ត', minWidth: 150 },
        ]"
      >
        <template #user_name="{ row }">
          <el-text tag="b">
            {{ row.user_name }}
          </el-text>
          <el-tag :type="row.gender === 1 ? 'success' : 'warning'" size="small">
            {{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}
          </el-tag>
        </template>
        <template #leave_type_name="{ row }">
          <el-text type="primary" tag="b">
            {{ row.leave_type_code }}
          </el-text>
          <el-text>
            {{ row.leave_type_name }}
          </el-text>
          <el-tag
            :type="row.leave_type_is_deduct === true ? 'danger' : 'success'"
            size="small"
          >
            {{ row.leave_type_is_deduct === true ? "កាត់លុយ" : "មិនកាត់លុយ" }}
          </el-tag>
        </template>
        <template #total_day="{ row }">
          <el-text>
            {{ row.total_day }}
            {{ row.deduct_type_name }}
          </el-text>
          <el-text tag="b" type="danger">
            {{ row.deduct_type_code }}
          </el-text>
        </template>
        <template #status="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ row.status_string }}
          </el-tag>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែច្បាប់"c placement="top">
            <AppButton
              v-if="canEditLeave && row.status === 1"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEditDialog(row)"
            />
          </el-tooltip>
          <template v-if="canEditStatusLeave">
            <el-tooltip content="អនុញ្ញាត" placement="top">
              <AppButton
                size="small"
                icon="Check"
                type="success"
                circle
                @click="updateStatus(row, 2)"
              />
            </el-tooltip>
            <el-tooltip content="បដិសេធ" placement="top">
              <AppButton
                size="small"
                icon="Close"
                type="danger"
                circle
                @click="updateStatus(row, 3)"
              />
            </el-tooltip>
            <el-tooltip content="ត្រឡប់" placement="top">
              <AppButton
                size="small"
                icon="Refresh"
                type="warning"
                circle
                @click="updateStatus(row, 1)"
              />
            </el-tooltip>
            <el-tooltip content="លុប" placement="top">
              <AppButton
                v-if="row.status === 1"
                size="small"
                icon="Delete"
                type="danger"
                circle
                @click="deleteleave(row)"
              />
            </el-tooltip>
          </template>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEditMode ? 'កែប្រែការសុំច្បាប់' : 'បង្កើតការសុំច្បាប់'"
      width="740px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
          <el-form-item label="ប្រភេទច្បាប់" prop="leave_type_id">
            <el-select
              v-model="form.leave_type_id"
              placeholder="ជ្រើសរើសប្រភេទច្បាប់"
              style="width: 100%"
              size="large"
            >
              <el-option
                v-for="lt in leavetype"
                :key="lt.id"
                :label="`${lt.name} (${lt.company_name}) ${lt.is_deduct ? '- កាត់លុយ' : '- មិនកាត់លុយ'}`"
                :value="lt.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="ប្រភេទឯកតា" prop="deduct_type_id">
            <el-select
              v-model="form.deduct_type_id"
              placeholder="ជ្រើសរើសឯកតា"
              style="width: 100%"
              clearable
              size="large"
            >
              <el-option
                v-for="dt in leavededucttype"
                :key="dt.id"
                :label="dt.name"
                :value="dt.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ថ្ងៃចាប់ផ្តើម" prop="start_date">
            <el-date-picker
              v-model="form.start_date"
              type="date"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              size="large"
            />
          </el-form-item>
          <el-form-item label="ថ្ងៃបញ្ចប់" prop="end_date">
            <el-date-picker
              v-model="form.end_date"
              type="date"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              size="large"
            />
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ថ្ងៃចូលធ្វើការវិញ" prop="back_to_work_date">
            <el-date-picker
              v-model="form.back_to_work_date"
              type="date"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              size="large"
            />
          </el-form-item>
          <AppInput
            label="ចំនួនថ្ងៃសរុប"
            prop="total_day"
            v-model.number="form.total_day"
            type="number"
            size="large"
          />
        </div>

        <el-form-item label="អ្នកអនុម័ត" prop="approve_by">
          <el-select
            v-model="form.approve_by"
            placeholder="ជ្រើសរើសអ្នកអនុម័ត"
            style="width: 100%"
            clearable
            size="large"
          >
            <el-option
              v-for="u in userapprove"
              :key="u.id"
              :label="u.user_name"
              :value="u.id"
            />
          </el-select>
        </el-form-item>

        <AppInput
          prop="reason"
          label="មូលហេតុ"
          clearable="true"
          v-model="form.reason"
          type="textarrea"
          :block="true"
        >
        </AppInput>
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
.filter-card {
  margin-bottom: 16px;
  border-radius: 6px;
}
.filter-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}
.table-card {
  border-radius: 6px;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
.form-row {
  display: flex;
  gap: 16px;
}
.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
</style>
