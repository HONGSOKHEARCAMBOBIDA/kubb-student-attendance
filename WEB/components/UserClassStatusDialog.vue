<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="ផ្លាស់ប្តូរស្ថានភាពសិស្ស"
    width="420px"
  >
    <el-form label-position="top">
      <AppSelect v-model="statusForm.status" size="large" :options="userClassStatus" />
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton @click="handleUpdate" type="primary" :loading="saving" size="large" :block="false">
        កែប្រែ
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { editUserClass } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import { useNotification } from "../composables/useNotification.js";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  // row comes from GetClass -> UserResponse, built from
  // response.StudentWithClass { ..., user_class_id }, so use UserClassID,
  // not the user's own id.
  student: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const saving = ref(false);

const userClassStatus = [
  { value: "STUDY", label: "កំពុងសិក្សា" },
  { value: "SUSPEND", label: "ព្យួរការសិក្សា" },
  { value: "TRANSFER", label: "ដូរ/ផ្ទេរជំនាញ" },
  { value: "DROPPED", label: "បោះបង់ការសិក្សា" },
];

const statusForm = reactive({
  id: null,
  status: null,
});

watch(
  () => props.modelValue,
  (open) => {
    if (!open || !props.student) return;
    statusForm.id = props.student.UserClassID ?? props.student.id;
    statusForm.status = props.student.status;
  },
);

async function handleUpdate() {
  if (!statusForm.id) return;
  saving.value = true;
  try {
    await editUserClass(statusForm.id, { status: statusForm.status });
    notify.success("ផ្លាស់ប្តូរស្ថានភាពបានជោគជ័យ");
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "ផ្លាស់ប្តូរបរាជ័យ");
  } finally {
    saving.value = false;
  }
}
</script>
