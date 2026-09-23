<script setup>
import { ref, reactive, onMounted } from "vue";
import { useNotification } from "../../composables/useNotification.js";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppInput from "../../components/AppInput.vue";
import {
  getShift,
  addShift,
  editShift,
} from "../api/services.js";

const notify = useNotification();

const shifts = ref([]);
const loading = ref(false);
const submitting = ref(false);
const formRef = ref();

const dialogVisible = ref(false);
const isEditMode = ref(false);
const editingId = ref(null);

const defaultForm = () => ({
  name: "",
  session1: "",
  session2: "",
  session3: "",
  session4: "",
  session5: "",
});
const form = reactive(defaultForm());

// Field names/required-ness mirror request.ShiftRequestCreate / ShiftRequestUpdate
const rules = {
  name: [{ required: true, message: "សូមបញ្ចូលឈ្មោះវេន", trigger: "blur" }],
  session1: [{ required: true, message: "សូមបញ្ចូលវេនទី១", trigger: "blur" }],
  session2: [{ required: true, message: "សូមបញ្ចូលវេនទី២", trigger: "blur" }],
  session3: [{ required: true, message: "សូមបញ្ចូលវេនទី៣", trigger: "blur" }],
  session4: [{ required: true, message: "សូមបញ្ចូលវេនទី៤", trigger: "blur" }],
  session5: [{ required: true, message: "សូមបញ្ចូលវេនទី៥", trigger: "blur" }],
};

// Backend (GetShift) has no pagination/filter params - it always
// returns every row ordered by id ASC.
async function fetchShift() {
  loading.value = true;
  try {
    const res = await getShift();
    shifts.value = res.data.data || res.data || [];
  } catch {
    notify.error("Failed to load shifts");
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
    name: row.name,
    session1: row.session1,
    session2: row.session2,
    session3: row.session3,
    session4: row.session4,
    session5: row.session5,
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
      await editShift(editingId.value, { ...form });
      notify.success("កែប្រែជោគជ័យ");
    } else {
      await addShift({ ...form });
      notify.success("បន្ថែមជោគជ័យ");
    }
    dialogVisible.value = false;
    await fetchShift();
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  fetchShift();
});
</script>

<template>
  <div>
    <div class="toolbar">
      <AppButton type="primary" @click="openCreateDialog">
        បន្ថែមវេន
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        :show-pagination="false"
        show-index
        :data="shifts"
        :loading="loading"
        :actions-width="100"
        :columns="[
          { label: 'ឈ្មោះវេន', prop: 'name', minWidth: 140 },
          { label: 'វេនទី១', prop: 'session1', minWidth: 120 },
          { label: 'វេនទី២', prop: 'session2', minWidth: 120 },
          { label: 'វេនទី៣', prop: 'session3', minWidth: 120 },
          { label: 'វេនទី៤', prop: 'session4', minWidth: 120 },
          { label: 'វេនទី៥', prop: 'session5', minWidth: 120 },
        ]"
      >
        <template #actions="{ row }">
          <el-tooltip content="កែប្រែវេន" placement="top">
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
      :title="isEditMode ? 'កែប្រែវេន' : 'បង្កើតវេន'"
      width="640px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="form-row">
          <AppInput
            label="ឈ្មោះវេន"
            prop="name"
            v-model="form.name"
            size="large"
          />
        </div>

        <div class="form-row">
          <AppInput
            label="វេនទី១"
            prop="session1"
            v-model="form.session1"
            size="large"
          />
          <AppInput
            label="វេនទី២"
            prop="session2"
            v-model="form.session2"
            size="large"
          />
        </div>

        <div class="form-row">
          <AppInput
            label="វេនទី៣"
            prop="session3"
            v-model="form.session3"
            size="large"
          />
          <AppInput
            label="វេនទី៤"
            prop="session4"
            v-model="form.session4"
            size="large"
          />
        </div>

        <div class="form-row">
          <AppInput
            label="វេនទី៥"
            prop="session5"
            v-model="form.session5"
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