<script setup>
import { ref, reactive, onMounted, computed, watch, onUnmounted } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import {
  viewMajorWithPagination,
  addMajor,
  editMajor,
  toggleMajor,
  getMajorSubjects,
  addMajorSubject,
  toggleMajorSubject,
  removeMajorSubject,
  getSubject,
} from "../api/services.js";
import { useUserDataStore } from "../stores/user_data.js";
import AppSelect from "../../components/AppSelect.vue";

let searchTimer = null;
const notify = useNotification();
const userDataStore = useUserDataStore();

const majors = ref([]);
const loading = ref(false);
const submitting = ref(false);
const togglingId = ref(null);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
});

// Backend only supports filtering by name today
const filter = reactive({
  name: "",
});

const defaultForm = () => ({
  code: "",
  name_kh: "",
  name_en: "",
});
const form = reactive(defaultForm());

const rules = {
  code: [{ required: true, message: "សូមបញ្ចូលលេខកូដ", trigger: "blur" }],
  name_kh: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាខ្មែរ", trigger: "blur" }],
  name_en: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាអង់គ្លេស", trigger: "blur" }],
};

const canCreateMajor = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.Major"),
);
const canEditMajor = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.Major"),
);

async function fetchMajor() {
  loading.value = true;
  try {
    const res = await viewMajorWithPagination({
      page: pagination.page,
      page_size: pagination.page_size,
      ...filter,
    });
    majors.value = res.data.data || [];
    pagination.total = res.data.pagination?.totalCount || 0;
    console.log(majors.value)
  } catch {
    notify.error("Failed to load majors");
  } finally {
    loading.value = false;
  }
}

function handlePageChange() {
  fetchMajor();
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
    code: row.code,
    name_kh: row.name_kh,
    name_en: row.name_en,
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
      await editMajor(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addMajor({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchMajor();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

async function handleToggle(row) {
  togglingId.value = row.id;
  try {
    await toggleMajor(row.id);
    notify.success("ប្តូរស្ថានភាពជោគជ័យ");
    await fetchMajor();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចប្តូរស្ថានភាពបានទេ");
  } finally {
    togglingId.value = null;
  }
}

watch(
  filter,
  () => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      pagination.page = 1;
      fetchMajor();
    }, 300); // debounce so typing doesn't fire a request per keystroke
  },
  { deep: true },
);

/* ---------------- Major subjects dialog (major_subject) ---------------- */

const subjectsDialogVisible = ref(false);
const activeMajor = ref(null);
const majorSubjects = ref([]);
const majorSubjectsLoading = ref(false);
const addingSubject = ref(false);
const removingId = ref(null);
const subjectTogglingId = ref(null);

const subjectOptions = ref([]);
const subjectOptionsLoading = ref(false);

const subjectPagination = reactive({
  page: 1,
  page_size: 40,
  total: 0,
});

const subjectFormRef = ref();
const defaultSubjectForm = () => ({
  subject_id: null,
  year: null,
  semester: null,
});
const subjectForm = reactive(defaultSubjectForm());

const subjectRules = {
  subject_id: [{ required: true, message: "សូមជ្រើសរើសមុខវិជ្ជា", trigger: "change" }],
  year: [{ required: true, message: "សូមបញ្ចូលឆ្នាំ", trigger: "blur" }],
  semester: [{ required: true, message: "សូមបញ្ចូលឆមាស", trigger: "blur" }],
};

async function openSubjectsDialog(row) {
  activeMajor.value = row;
  Object.assign(subjectForm, defaultSubjectForm());
  subjectPagination.page = 1;
  subjectsDialogVisible.value = true;
  await Promise.all([fetchMajorSubjects(), searchSubjects()]);
}

async function fetchMajorSubjects() {
  if (!activeMajor.value) return;
  majorSubjectsLoading.value = true;
  try {
    const res = await getMajorSubjects(activeMajor.value.id, {
      page: subjectPagination.page,
      page_size: subjectPagination.page_size,
    });
    majorSubjects.value = res.data.data || [];
    subjectPagination.total = res.data.pagination?.totalCount || 0;
  } catch {
    notify.error("Failed to load subjects for this major");
  } finally {
    majorSubjectsLoading.value = false;
  }
}

async function searchSubjects(query) {
  try {
    const params = {
      page: 1,
      page_size: 10,
    };

    if (query.trim()) {
      params.name = query.trim();
    }

    const res = await getSubject(params);

    subjectOptions.value = (res.data.data || []).map((c) => ({
      label: c.name_kh,
      value: c.id,
      raw: c,
    }));
  } catch (e) {
    
  }
}



function handleSubjectsPageChange() {
  fetchMajorSubjects();
}

async function handleAddSubject() {
  if (!subjectFormRef.value || !activeMajor.value) return;
  const valid = await subjectFormRef.value.validate().catch(() => false);
  if (!valid) return;

  addingSubject.value = true;
  try {
    await addMajorSubject(activeMajor.value.id, { ...subjectForm });
    notify.success("បន្ថែមមុខវិជ្ជាជោគជ័យ");
    Object.assign(subjectForm, defaultSubjectForm());
    subjectFormRef.value.clearValidate();
    subjectPagination.page = 1;
    await fetchMajorSubjects();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចបន្ថែមមុខវិជ្ជាបានទេ");
  } finally {
    addingSubject.value = false;
  }
}

async function handleToggleMajorSubject(row) {
  subjectTogglingId.value = row.id;
  try {
    await toggleMajorSubject(row.id);
    notify.success("ប្តូរស្ថានភាពជោគជ័យ");
    await fetchMajorSubjects();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចប្តូរស្ថានភាពបានទេ");
  } finally {
    subjectTogglingId.value = null;
  }
}

const getYearType = (year) => {
  const types = {
    1: "primary",
    2: "success",
    3: "warning",
    4: "danger",
  }

  return types[year] || "info"
}

async function handleRemoveSubject(row) {
  removingId.value = row.id;
  try {
    await removeMajorSubject(row.id);
    notify.success("លុបមុខវិជ្ជាជោគជ័យ");
    await fetchMajorSubjects();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចលុបមុខវិជ្ជាបានទេ");
  } finally {
    removingId.value = null;
  }
}

onMounted(() => {
  fetchMajor();
});
onUnmounted(() => clearTimeout(searchTimer));
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 16 },
        { slot: 'add', span: 8 },
      ]"
      :action-span="4"
    >
      <template #name>
        <AppInput
          v-model="filter.name"
          placeholder="ស្វែងរកតាមឈ្មោះជំនាញ"
          prefix-icon="Search"
          clearable
        />
      </template>

      <template #actions>
        <AppButton v-if="canCreateMajor" type="primary" @click="openCreateDialog">
          បន្ថែមជំនាញ
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        show-index
        :data="majors"
        :loading="loading"
        :actions-width="180"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        @page-change="handlePageChange"
        :columns="[
          { label: 'លេខកូដ', prop: 'code', minWidth: 100 },
          { label: 'ឈ្មោះខ្មែរ', prop: 'name_kh', minWidth: 160 },
          { label: 'ឈ្មោះអង់គ្លេស', prop: 'name_en', minWidth: 160 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
        ]"
      >
        <template #status="{ row }">
          <el-switch
            :model-value="row.is_active"
            :loading="togglingId === row.id"
            :disabled="!canEditMajor"
            @change="handleToggle(row)"
          />
        </template>

        <template #actions="{ row }">
          <el-tooltip content="គ្រប់គ្រងមុខវិជ្ជា" placement="top">
            <AppButton
              :disabled="row.is_active === false"
              size="small"
              icon="Notebook"
              type="primary"
              circle
              @click="openSubjectsDialog(row)"
            />
          </el-tooltip>
          <el-tooltip content="កែប្រែជំនាញ" placement="top">
            <AppButton
              :disabled="row.is_active === false"
              v-if="canEditMajor"
              size="small"
              icon="Edit"
              type="warning"
              circle
              @click="openEditDialog(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>

    <!-- Create / edit major -->
    <AppDialog
      v-model="dialogVisible"
      :title="isEditMode ? 'កែប្រែជំនាញ' : 'បង្កើតជំនាញ'"
      width="640px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
          <AppInput label="លេខកូដ" prop="code" v-model="form.code" size="large" />
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
        <AppButton @click="dialogVisible = false">បោះបង់</AppButton>
        <AppButton type="primary" :loading="submitting" @click="handleSubmit">
          រក្សាទុក
        </AppButton>
      </template>
    </AppDialog>

    <!-- Subjects belonging to the selected major (major_subject) -->
    <AppDialog
      v-model="subjectsDialogVisible"
      :title="activeMajor ? `មុខវិជ្ជាសម្រាប់ ${activeMajor.name_kh}` : 'មុខវិជ្ជា'"
      width="70%"
    >
      <el-form
        ref="subjectFormRef"
        :model="subjectForm"
        :rules="subjectRules"
        label-position="top"
        class="add-subject-form"
      >
        <div class="form-row">
          <el-form-item label="មុខវិជ្ជា" prop="subject_id" class="grow">
            <AppSelect
          v-model="subjectForm.subject_id"
          :options="subjectOptions"
          placeholder="ជ្រើសរើសមុខវិជ្ជា"
          size="large"
          filterable
          remote
          :remote-method="searchSubjects"
          clearable
        />
          </el-form-item>
          <AppInput
            label="ឆ្នាំ"
            prop="year"
            v-model.number="subjectForm.year"
            type="number"
            size="large"
          />
          <AppInput
            label="ឆមាស"
            prop="semester"
            v-model.number="subjectForm.semester"
            type="number"
            size="large"
          />
        </div>

        <AppButton
          type="primary"
          :loading="addingSubject"
          @click="handleAddSubject"
        >
          បន្ថែមមុខវិជ្ជា
        </AppButton>
      </el-form>

      <AppTable
        
        class="subjects-table"
        show-index
        :data="majorSubjects"
        :loading="majorSubjectsLoading"
        :actions-width="140"
        v-model:current-page="subjectPagination.page"
        v-model:page-size="subjectPagination.page_size"
        :total="subjectPagination.total"
        @page-change="handleSubjectsPageChange"
        :columns="[
          { label: 'លេខកូដ', prop: 'subject_code', minWidth: 100 },
          { label: 'ឈ្មោះខ្មែរ', prop: 'subject_name_kh', minWidth: 160 },
          { label: 'ម៉ោងក្រេឌីត', prop: 'credit_hour', width: 100 },
          { label: 'ឆ្នាំ', slot: 'year', width: 80 },
          { label: 'ឆមាស', prop: 'semester', width: 80 },
          { label: 'ស្ថានភាព', slot: 'status', width: 110 },
        ]"
      >
      <template #year="{ row }">
  <el-text :type="getYearType(row.year)">
    ឆ្នាំទី {{ row.year }}
  </el-text>
</template>
        <template #status="{ row }">
          <el-switch
            :model-value="row.is_active"
            :loading="subjectTogglingId === row.id"
            @change="handleToggleMajorSubject(row)"
          />
        </template>

        <template #actions="{ row }">
          <el-tooltip content="លុបចេញពីជំនាញ" placement="top">
            <AppButton
              size="small"
              icon="Delete"
              type="danger"
              circle
              :loading="removingId === row.id"
              @click="handleRemoveSubject(row)"
            />
          </el-tooltip>
        </template>
      </AppTable>

      <template #footer>
        <AppButton @click="subjectsDialogVisible = false">បិទ</AppButton>
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
.add-subject-form {
  margin-bottom: 16px;
}
.subjects-table {
  margin-top: 8px;
}
</style>