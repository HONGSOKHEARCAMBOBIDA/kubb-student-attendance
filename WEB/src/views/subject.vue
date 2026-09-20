<script setup>
import { ref, reactive, onMounted, computed, watch, onUnmounted } from "vue";
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
} from "../api/services.js";
import { useUserDataStore } from "../stores/user_data.js";

let searchTimer = null;
const notify = useNotification();
const userDataStore = useUserDataStore();

const subjects = ref([]);
const loading = ref(false);
const submitting = ref(false);
const togglingId = ref(null);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const pagination = reactive({
  page: 1,
  page_size: 15,
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
  credit_hour: null,
});
const form = reactive(defaultForm());

// Field names/required-ness mirror request.SubjectRequestCreate / SubjectRequestUpdate
const rules = {
  code: [{ required: true, message: "សូមបញ្ចូលលេខកូដ", trigger: "blur" }],
  name_kh: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាខ្មែរ", trigger: "blur" }],
  name_en: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាអង់គ្លេស", trigger: "blur" }],
  credit_hour: [{ required: true, message: "សូមបញ្ចូលចំនួនម៉ោងក្រេឌីត", trigger: "blur" }],
};

const canCreateSubject = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.subject"),
);
const canEditSubject = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.subject"),
);


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

function handlePageChange() {
  fetchSubject();
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
    credit_hour: row.credit_hour,
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
      await editSubject(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addSubject({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchSubject();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

async function handleToggle(row) {
  togglingId.value = row.id;
  try {
    await toggleSubject(row.id);
    notify.success("ប្តូរស្ថានភាពជោគជ័យ");
    await fetchSubject();
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
      fetchSubject();
    }, 300); // debounce so typing doesn't fire a request per keystroke
  },
  { deep: true },
);

onMounted(() => {
  fetchSubject();
});
onUnmounted(() => clearTimeout(searchTimer));
</script>

<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 10 },
      ]"
      :action-span="4"
    >
      <template #name>
        <AppInput
          v-model="filter.name"
          placeholder="ស្វែងរកតាមឈ្មោះមុខវិជ្ជា"
          prefix-icon="Search"
          clearable
        />
      </template>

      <template #actions>
        <AppButton v-if="canCreateSubject" type="primary" @click="openCreateDialog">
          បន្ថែមមុខវិជ្ជា
        </AppButton>
      </template>
    </AppFilterBar>

    <el-card class="table-card">
      <AppTable
        show-index
        :data="subjects"
        :loading="loading"
        :actions-width="140"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        @page-change="handlePageChange"
        :columns="[
          { label: 'លេខកូដ', prop: 'code', minWidth: 100 },
          { label: 'ឈ្មោះខ្មែរ', prop: 'name_kh', minWidth: 160 },
          { label: 'ឈ្មោះអង់គ្លេស', prop: 'name_en', minWidth: 160 },
          { label: 'ម៉ោងក្រេឌីត', prop: 'credit_hour', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
        ]"
      >
        <template #status="{ row }">
          <el-switch
            :model-value="row.is_active"
            :loading="togglingId === row.id"
            :disabled="!canEditSubject"
            @change="handleToggle(row)"
          />
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែមុខវិជ្ជា" placement="top">
            <AppButton
              v-if="canEditSubject"
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
        <AppButton @click="dialogVisible = false">បោះបង់</AppButton>
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
</style>