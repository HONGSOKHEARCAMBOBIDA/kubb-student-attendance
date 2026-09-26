<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="កែប្រែសិស្ស"
    width="500px"
  >
    <el-form :model="studentForm" ref="studentFormRef" label-position="top">
      <el-form-item label="ឈ្មោះខ្មែរ" prop="name_kh">
        <el-input v-model.trim="studentForm.name_kh" size="large" />
      </el-form-item>
      <el-form-item label="ឈ្មោះឡាតាំង" prop="name_en">
        <el-input v-model.trim="studentForm.name_en" size="large" />
      </el-form-item>
      <el-form-item label="ភេទ" prop="gender">
        <el-radio-group v-model="studentForm.gender" size="large">
          <el-radio-button :value="1">ប្រុស</el-radio-button>
          <el-radio-button :value="2">ស្រី</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="អត្តលេខ" prop="code">
        <el-input v-model.trim="studentForm.code" size="large" />
      </el-form-item>
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton @click="handleSave" type="primary" :loading="saving" size="large" :block="false">
        កែប្រែ
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { updateUser } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import { useNotification } from "../composables/useNotification.js";
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  student: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const studentFormRef = ref();
const saving = ref(false);
const editStudentId = ref(null);

const studentForm = reactive({
  name_kh: "",
  name_en: "",
  gender: null,
  code: "",
});

watch(
  () => props.modelValue,
  (open) => {
    if (!open || !props.student) return;
    const row = props.student;
    editStudentId.value = row.id;
    studentForm.name_kh = row.name_kh || "";
    studentForm.name_en = row.name_en || "";
    studentForm.gender = row.gender !== undefined && row.gender !== null ? Number(row.gender) : null;
    studentForm.code = row.code || "";
  },
);

async function handleSave() {
  await studentFormRef.value.validate();
  saving.value = true;
  try {
    // Only send fields that make sense as a partial update — matches
    // the *string / *int pointer semantics on the Go side (nil = untouched).
    const payload = {
      name_kh: studentForm.name_kh,
      name_en: studentForm.name_en,
      gender: studentForm.gender,
      code: studentForm.code,
    };
    await updateUser(editStudentId.value, payload);
    notify.success("កែប្រែសិស្សបានជោគជ័យ");
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "កែប្រែបរាជ័យ");
  } finally {
    saving.value = false;
  }
}
</script>
