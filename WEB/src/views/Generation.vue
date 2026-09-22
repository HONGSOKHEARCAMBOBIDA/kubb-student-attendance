<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import {
  getGeneration,
  addGeneration,
  editGeneration,
  toggleGeneration,
} from "../api/services.js";
import { useUserDataStore } from "../stores/user_data.js";

const notify = useNotification();
const userDataStore = useUserDataStore();

const generations = ref([]);
const loading = ref(false);
const submitting = ref(false);
const togglingId = ref(null);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const defaultForm = () => ({
  code: "",
  name_kh: "",
  name_en: "",
  start_year: null,
  end_year: null,
});
const form = reactive(defaultForm());

// Field names/required-ness mirror request.GenerationRequestCreate / GenerationRequestUpdate
const rules = {
  code: [{ required: true, message: "សូមបញ្ចូលលេខកូដ", trigger: "blur" }],
  name_kh: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាខ្មែរ", trigger: "blur" }],
  name_en: [{ required: true, message: "សូមបញ្ចូលឈ្មោះជាភាសាអង់គ្លេស", trigger: "blur" }],
  start_year: [{ required: true, message: "សូមបញ្ចូលឆ្នាំចាប់ផ្តើម", trigger: "blur" }],
  end_year: [{ required: true, message: "សូមបញ្ចូលឆ្នាំបញ្ចប់", trigger: "blur" }],
};



// Backend (GetGeneration) has no pagination/filter params - it always
// returns every active (is_active = 1) row ordered by id ASC.
async function fetchGeneration() {
  loading.value = true;
  try {
    const res = await getGeneration();
    generations.value = res.data.data || res.data || [];
  } catch {
    notify.error("Failed to load generations");
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
    code: row.code,
    name_kh: row.name_kh,
    name_en: row.name_en,
    start_year: row.start_year,
    end_year: row.end_year,
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
      await editGeneration(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addGeneration({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchGeneration();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

async function handleToggle(row) {
  togglingId.value = row.id;
  try {
    await toggleGeneration(row.id);
    notify.success("ប្តូរស្ថានភាពជោគជ័យ");
    await fetchGeneration();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចប្តូរស្ថានភាពបានទេ");
  } finally {
    togglingId.value = null;
  }
}

onMounted(() => {
  fetchGeneration();
});
</script>

<template>
  <div>
    <div class="toolbar">
      <AppButton  type="primary" @click="openCreateDialog">
        បន្ថែមជំនាន់
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        :show-pagination="false"
        show-index
        :data="generations"
        :loading="loading"
        :actions-width="140"
        :columns="[
          { label: 'លេខកូដ', prop: 'code', minWidth: 100 },
          { label: 'ឈ្មោះខ្មែរ', prop: 'name_kh', minWidth: 160 },
          { label: 'ឈ្មោះអង់គ្លេស', prop: 'name_en', minWidth: 160 },
          { label: 'ឆ្នាំចាប់ផ្តើម', prop: 'start_year', width: 120 },
          { label: 'ឆ្នាំបញ្ចប់', prop: 'end_year', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 120 },
        ]"
      >
        <template #status="{ row }">
          <el-switch
            :model-value="row.is_active"
            :loading="togglingId === row.id"
            
            @change="handleToggle(row)"
          />
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែជំនាន់" placement="top">
            <AppButton
             
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
      :title="isEditMode ? 'កែប្រែជំនាន់' : 'បង្កើតជំនាន់'"
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

        <div class="form-row">
          <AppInput
            label="ឆ្នាំចាប់ផ្តើម"
            prop="start_year"
            v-model.number="form.start_year"
            type="number"
            size="large"
          />
          <AppInput
            label="ឆ្នាំបញ្ចប់"
            prop="end_year"
            v-model.number="form.end_year"
            type="number"
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
.toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}
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