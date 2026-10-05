<script setup>
import { computed, reactive, ref, watch } from "vue";
import { addMajorPrice } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppInput from "./AppInput.vue";
import { useNotification } from "../composables/useNotification.js";

const props = defineProps({
  raw: {type: Object,default: () => []},
  isEdit: { type: Boolean, default: false },
  modelValue: { type: Boolean, default: false },
  majorID: { type: Number, default: null },
  generations: { type: Array, default: () => [] },
  programmes: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const formRef = ref();
const saving = ref(false);

const title = computed(() => {
  const base = props.isEdit ? "កែប្រែថ្លៃសិក្សា" : "បន្ថែមថ្លៃសិក្សា";
  return props.raw?.name_kh ? `${base} - ជំនាញ${props.raw.name_kh} | មហាវិទ្យាល័យ${props.raw.faculty_name}`  : base;
});

const defaultForm = () => ({
  major_id: null,
  generation_id: null,
  programme_id: null,
  year: null,
  monthly_fee: 0,
  quarter_fee: 0,
  semester_fee: 0,
  year_fee: 0,
});
const form = reactive(defaultForm());

const yearOptions = [1, 2, 3, 4].map((y) => ({ label: `ឆ្នាំទី ${y}`, value: y }));

const feeRule = (message) => [
  { required: true, message, trigger: "blur" },
  {
    validator: (_, v, cb) =>
      v === "" || v === null || v === undefined || Number(v) < 0
        ? cb(new Error("តម្លៃមិនត្រឹមត្រូវ"))
        : cb(),
    trigger: "blur",
  },
];

const rules = {
  generation_id: [{ required: true, message: "សូមជ្រើសរើសជំនាន់", trigger: "change" }],
  programme_id: [{ required: true, message: "សូមជ្រើសរើសកម្មវិធីសិក្សា", trigger: "change" }],
  year: [{ required: true, message: "សូមជ្រើសរើសឆ្នាំ", trigger: "change" }],
  monthly_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំខែ"),
  quarter_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំត្រីមាស"),
  semester_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំឆមាស"),
  year_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំឆ្នាំ"),
};

// Reset the form and attach the major every time the dialog opens
watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    Object.assign(form, defaultForm(), { major_id: props.majorID });
    formRef.value?.clearValidate();
  },
);

function close() {
  emit("update:modelValue", false);
}

async function handleSave() {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    await addMajorPrice({
      major_id: props.majorID,
      generation_id: form.generation_id,
      programme_id: form.programme_id,
      year: form.year,
      monthly_fee: Number(form.monthly_fee),
      quarter_fee: Number(form.quarter_fee),
      semester_fee: Number(form.semester_fee),
      year_fee: Number(form.year_fee),
    });
    notify.success("បង្កើតបានជោគជ័យ");
    close();
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title= title
    width="640px"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
      <div class="form-row">
        <AppSelect
          label="ជំនាន់"
          prop="generation_id"
          v-model="form.generation_id"
          :options="generations"
          placeholder="ជ្រើសរើសជំនាន់"
          size="large"
          filterable
          clearable
        />
        <AppSelect
          label="កម្មវិធីសិក្សា"
          prop="programme_id"
          v-model="form.programme_id"
          :options="programmes"
          placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
          size="large"
          filterable
          clearable
        />
      </div>

      <div class="form-row">
        <AppSelect
          label="ឆ្នាំ"
          prop="year"
          v-model="form.year"
          :options="yearOptions"
          placeholder="ជ្រើសរើសឆ្នាំ"
          size="large"
        />
        <AppInput label="ថ្លៃប្រចាំខែ" prop="monthly_fee" v-model.number="form.monthly_fee" type="number" size="large" />
      </div>

      <div class="form-row">
        <AppInput label="ថ្លៃប្រចាំត្រីមាស" prop="quarter_fee" v-model.number="form.quarter_fee" type="number" size="large" />
        <AppInput label="ថ្លៃប្រចាំឆមាស" prop="semester_fee" v-model.number="form.semester_fee" type="number" size="large" />
      </div>

      <div >
        <AppInput label="ថ្លៃប្រចាំឆ្នាំ" prop="year_fee" v-model.number="form.year_fee" type="number" size="large" />
        <div class="form-spacer" />
      </div>
    </el-form>

    <template #footer>
      <AppButton @click="close">បោះបង់</AppButton>
      <AppButton type="primary" :loading="saving" @click="handleSave">រក្សាទុក</AppButton>
    </template>
  </AppDialog>
</template>

<style scoped>
.form-row {
  display: flex;
  gap: 20px;
}
.form-row .el-form-item,
.form-spacer {
  flex: 1;
  min-width: 0;
}
</style>