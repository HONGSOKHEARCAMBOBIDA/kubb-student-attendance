<script setup>
import { ref, reactive, computed, onMounted } from "vue";

import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import {
  getSubject,
  addSubject,
  editSubject,
  toggleSubject,
  importsubjectexcell,
} from "../api/services.js";
import { ElMessageBox } from "element-plus";
import { useUserDataStore } from "../stores/user_data.js";

const notify = useNotification();
const userDataStore = useUserDataStore();

const subjects = ref([]);
const loading = ref(false);
const submitting = ref(false);
const importing = ref(false);

const togglingId = ref(null);
const editingId = ref(null);

const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);

const pagination = reactive({
  page: 1,
  page_size: 15,
  total: 0,
});

const searchText = ref("");

const filter = reactive({
  name: "",
});

const defaultForm = () => ({
  code: "",
  name_kh: "",
  name_en: "",
  credit_hour: null,
});

const form = reactive(defaultForm());

const rules = {
  code: [
    {
      required: true,
      message: "សូមបញ្ចូលលេខកូដ",
      trigger: "blur",
    },
  ],

  name_kh: [
    {
      required: true,
      message: "សូមបញ្ចូលឈ្មោះជាភាសាខ្មែរ",
      trigger: "blur",
    },
  ],

  name_en: [
    {
      required: true,
      message: "សូមបញ្ចូលឈ្មោះជាភាសាអង់គ្លេស",
      trigger: "blur",
    },
  ],

  credit_hour: [
    {
      required: true,
      message: "សូមបញ្ចូលចំនួនម៉ោងក្រេឌីត",
      trigger: "blur",
    },
  ],
};

// ==================== Permissions ====================

const canCreateSubject = computed(() =>
  userDataStore.permissions?.some(
    (permission) => permission.name === "add.subject",
  ),
);

const canEditSubject = computed(() =>
  userDataStore.permissions?.some(
    (permission) => permission.name === "edit.subject",
  ),
);

// ==================== Search ====================

function handleSearch() {
  filter.name = searchText.value.trim();
  pagination.page = 1;

  fetchSubject();
}

function handleClearSearch() {
  searchText.value = "";
  filter.name = "";
  pagination.page = 1;

  fetchSubject();
}

// ==================== Fetch Data ====================

async function fetchSubject() {
  loading.value = true;

  try {
    const res = await getSubject({
      page: pagination.page,
      page_size: pagination.page_size,
      ...filter,
    });

    subjects.value = res.data.data || [];
    pagination.total = res.data.pagination?.totalCount || 0;
  } catch {
    notify.error("Failed to load subjects");
  } finally {
    loading.value = false;
  }
}

// ==================== Pagination ====================

function handlePageChange() {
  fetchSubject();
}

// ==================== Import ====================

async function handleImport(file) {
  try {
    await ElMessageBox.confirm(
      `តើអ្នកចង់នាំចូលឯកសារ "${file.name}" មែនទេ?`,
      "បញ្ជាក់ការនាំចូល",
      {
        confirmButtonText: "នាំចូល",
        cancelButtonText: "បោះបង់",
        type: "warning",
      },
    );
  } catch {
    return; // user clicked cancel
  }

  importing.value = true;

  try {
    await importsubjectexcell(file);

    notify.success("នាំចូលជោគជ័យ");

    await fetchSubject();
  } catch (error) {
    notify.error(error.response?.data?.error || "នាំចូលមិនបានជោគជ័យ");
  } finally {
    importing.value = false;
  }
}

function handleUploadRequest(options) {
  return handleImport(options.file);
}

// ==================== Create ====================

function openCreateDialog() {
  isEditMode.value = false;
  editingId.value = null;

  Object.assign(form, defaultForm());

  dialogVisible.value = true;
}

// ==================== Edit ====================

function openEditDialog(row) {
  isEditMode.value = true;
  editingId.value = row.id;

  Object.assign(form, {
    code: row.code,
    name_kh: row.name_kh,
    name_en: row.name_en,
    credit_hour: row.credit_hour,
  });

  dialogVisible.value = true;
}

// ==================== Submit ====================

async function handleSubmit() {
  if (!formRef.value) return;

  const valid = await formRef.value.validate().catch(() => false);

  if (!valid) return;

  submitting.value = true;

  try {
    if (isEditMode.value && editingId.value !== null) {
      await editSubject(editingId.value, { ...form });

      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addSubject({ ...form });

      notify.success("បន្ថែមជោគជ័យ");
    }

    dialogVisible.value = false;

    await fetchSubject();
  } catch (error) {
    notify.error(error.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

// ==================== Toggle Status ====================

async function handleToggle(row) {
  togglingId.value = row.id;

  try {
    await toggleSubject(row.id);

    notify.success("ប្តូរស្ថានភាពជោគជ័យ");

    await fetchSubject();
  } catch (error) {
    notify.error(error.response?.data?.error || "មិនអាចប្តូរស្ថានភាពបានទេ");
  } finally {
    togglingId.value = null;
  }
}

// ==================== Lifecycle ====================

onMounted(() => {
  fetchSubject();
});
</script>

<template>
  <div>
    <!-- Filter -->
    <AppFilterBar :fields="[
      { slot: 'name', span: 6 }
      
    ]"
    show-back
    >
      <template #name>
        <AppInput
          label="ស្វែងរក"
          v-model="searchText"
          placeholder="ស្វែងរកតាមឈ្មោះមុខវិជ្ជា"
          prefix-icon="Search"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleClearSearch"
        />
      </template>

      <template #actions>
        <AppButton v-if="canCreateSubject" @click="openCreateDialog">
          បន្ថែមមុខវិជ្ជា
        </AppButton>
      </template>
    </AppFilterBar>

    <!-- Table -->
    <el-card class="table-card">
      <AppTable
        show-index
        :data="subjects"
        :loading="loading"
        :actions-width="80"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        @page-change="handlePageChange"
        :columns="[
          { label: 'លេខកូដ', prop: 'code', minWidth: 100 },
          { label: 'ឈ្មោះខ្មែរ', prop: 'name_kh', minWidth: 160 },
          { label: 'ឈ្មោះអង់គ្លេស', prop: 'name_en', minWidth: 160 },
          { label: 'ម៉ោងក្រេឌីត', prop: 'credit_hour', width: 80 },
          {
            label: 'ស្ថានភាព',
            slot: 'status',
            width: 70,
            align: 'center',
          },
        ]"
      >
        <template #filter>
          <el-upload
            v-if="canCreateSubject"
            :show-file-list="false"
            accept=".xlsx"
            :http-request="handleUploadRequest"
          >
            <AppButton icon="Bottom" :loading="importing">នាំចូល</AppButton>
          </el-upload>
        </template>

        <!-- Status -->
        <template #status="{ row }">
          <el-switch
            size="small"
            :model-value="row.is_active"
            :loading="togglingId === row.id"
            :disabled="!canEditSubject"
            @change="handleToggle(row)"
          />
        </template>

        <!-- Actions -->
        <template #actions="{ row }">
          <AppButton
            v-if="canEditSubject"
            size="small"
            @click="openEditDialog(row)"
          >
            កែប្រែ
          </AppButton>
        </template>
      </AppTable>
    </el-card>

    <!-- Dialog -->
    <AppDialog
      v-model="dialogVisible"
      :title="isEditMode ? 'កែប្រែមុខវិជ្ជា' : 'បង្កើតមុខវិជ្ជា'"
      width="640px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
          <AppInput
            label="លេខកូដ"
            prop="code"
            v-model="form.code"
            size="large"
          />

          <AppInput
            label="ម៉ោងក្រេឌីត"
            prop="credit_hour"
            v-model.number="form.credit_hour"
            type="number"
            size="large"
          />
        </div>

        <div class="form-row">
          <AppInput
            label="ឈ្មោះជាភាសាខ្មែរ"
            prop="name_kh"
            v-model="form.name_kh"
            size="large"
          />

          <AppInput
            label="ឈ្មោះជាភាសាអង់គ្លេស"
            prop="name_en"
            v-model="form.name_en"
            size="large"
          />
        </div>
      </el-form>

      <template #footer>
        <AppButton @click="dialogVisible = false"> បោះបង់ </AppButton>

        <AppButton type="primary" :loading="submitting" @click="handleSubmit">
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<style scoped>
.table-card {
  border-radius: 6px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}

@media (max-width: 600px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style>