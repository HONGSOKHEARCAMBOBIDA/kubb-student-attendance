<script setup>
import { ref, reactive, watch } from "vue";
import { ElMessage } from "element-plus";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";
import {
  getStudentFamily, 
  addStudentFamily, 
  updateStudentFamily,
} from "../src/api/services";

const props = defineProps({
  modelValue: Boolean,
  studentRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const formRef = ref();
const saving = ref(false);
const loading = ref(false);
const familyId = ref(null); 

// Father and mother share the same fields, so build them from one list.
const parents = [
  { key: "father", title: "ព័ត៌មានឪពុក" },
  { key: "mother", title: "ព័ត៌មានម្ដាយ" },
];

const textFields = [
  { suffix: "name", label: "ឈ្មោះខ្មែរ" },
  { suffix: "english_name", label: "ឈ្មោះឡាតាំង" },
  { suffix: "age", label: "អាយុ" },
  { suffix: "phone_number", label: "លេខទូរស័ព្ទ" },
  { suffix: "occupation", label: "មុខរបរ" },
  { suffix: "workplace", label: "កន្លែងធ្វើការ" },
];

const emptyForm = () => {
  const f = {};
  for (const { key } of parents) {
    f[`${key}_is_alive`] = true;
    for (const { suffix } of textFields) f[`${key}_${suffix}`] = "";
  }
  return f;
};

const form = reactive(emptyForm());

function resetForm() {
  familyId.value = null;
  Object.assign(form, emptyForm());
}

async function loadFamily() {
  resetForm();
  if (!props.studentRaw) return;
  loading.value = true;
  try {
    let f = null;
    try {
      const res = await getStudentFamily(props.studentRaw.id);
      f = res.data.data;
    } catch (e) {
      if (e.response?.status !== 404) throw e; // 404 => no record yet, create mode
    }
    if (!f) return;

    familyId.value = f.id;
    for (const { key } of parents) {
      form[`${key}_is_alive`] = f[`${key}_is_alive`] ?? true;
      for (const { suffix } of textFields) {
        form[`${key}_${suffix}`] = f[`${key}_${suffix}`] ?? "";
      }
    }
  } catch {
    ElMessage.error("Failed to load student family");
  } finally {
    loading.value = false;
  }
}

watch(
  () => props.modelValue,
  (open) => open && loadFamily(),
);

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  try {
    // Nullable string columns: send null instead of "" (matches *string in Go)
    const payload = { student_id: props.studentRaw.id };
    for (const { key } of parents) {
      payload[`${key}_is_alive`] = form[`${key}_is_alive`];
      for (const { suffix } of textFields) {
        payload[`${key}_${suffix}`] = form[`${key}_${suffix}`] || null;
      }
    }

    if (familyId.value) {
      await updateStudentFamily(familyId.value, payload);
    } else {
      await addStudentFamily(payload);
    }
    ElMessage.success("រក្សាទុកបានជោគជ័យ");
    emit("saved");
    emit("update:modelValue", false);
  } catch (e) {
    ElMessage.error(
      e.response?.data?.message || e.response?.data?.error || "Failed to save",
    );
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :title="`ព័ត៌មានអំពីគ្រួសារ${studentRaw ? ' - ' + studentRaw.name_kh : ''}`"
    width="40%"
  >
    <el-form
      ref="formRef"
      :model="form"
      label-position="top"
      v-loading="loading"
    >
      <template v-for="(p, i) in parents" :key="p.key">
        <el-divider v-if="i > 0" />
        <h4 style="margin: 0 0 12px">{{ p.title }}</h4>

        <el-row :gutter="16">
          <el-col :xs="24" :sm="12" v-for="f in textFields" :key="f.suffix">
            <AppInput
              :label="f.label"
              :prop="`${p.key}_${f.suffix}`"
              clearable
              v-model="form[`${p.key}_${f.suffix}`]"
            />
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="ស្ថានភាព" :prop="`${p.key}_is_alive`">
              <el-switch
                v-model="form[`${p.key}_is_alive`]"
                active-text="នៅរស់"
                inactive-text="ស្លាប់"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </template>
    </el-form>

    <template #footer>
      <AppButton
        @click="emit('update:modelValue', false)"
        size="large"
        :block="false"
        type="warning"
      >
        បោះបង់
      </AppButton>
      <AppButton
        @click="handleSave"
        type="primary"
        :loading="saving"
        size="large"
        :block="false"
      >
        {{ familyId ? "កែប្រែ" : "បង្កេីត" }}
      </AppButton>
    </template>
  </AppDialog>
</template>