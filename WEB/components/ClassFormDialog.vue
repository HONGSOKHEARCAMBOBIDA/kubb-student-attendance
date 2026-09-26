<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="isEdit ? 'កែប្រែថ្នាក់' : 'បន្ថែមថ្នាក់'"
    width="640px"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
      <AppInput label="ឈ្មោះថ្នាក់" prop="name" v-model.trim="form.name" size="large" placeholder="បញ្ចូលឈ្មោះថ្នាក់"></AppInput>

      <div class="form-row">
        <AppSelect label="ជំនាញ" prop="major_id" v-model="form.major_id" size="large" placeholder="ជ្រើសរើសជំនាញ" :options="majors"></AppSelect>
        <AppSelect label="វេន" prop="shift_id" v-model="form.shift_id" size="large" placeholder="ជ្រើសរើសវេន" :options="shifts"></AppSelect>
      </div>

      <div class="form-row">
        <AppSelect label="ជំនាន់" prop="generation_id" v-model="form.generation_id" size="large" placeholder="ជ្រើសរើសជំនាន់" :options="generations"></AppSelect>
        <AppSelect label="កម្មវិធីសិក្សា" prop="programme_id" v-model="form.programme_id" size="large" placeholder="ជ្រើសរើសកម្មវិធីសិក្សា" :options="programmes"></AppSelect>
      </div>

      <div class="form-row">
        <AppInput type="number" label="ឆ្នាំ" prop="year" v-model.number="form.year" size="large" placeholder="1"></AppInput>
        <AppInput type="number" label="ឆមាស" prop="semester" v-model.number="form.semester" size="large" placeholder="1"></AppInput>
      </div>

      <div class="form-row">
        <AppInput type="number" label="ក្រុម" prop="group" v-model.number="form.group" size="large" placeholder="1"></AppInput>
         <AppInput type="number" label="Term" prop="term" v-model.number="form.term" size="large" placeholder="1"></AppInput>
        <AppSelect
          size="large"
          v-model="form.type"
          :options="classType"
          label="ប្រភេទថ្នាក់"
          placeholder="ប្រភេទថ្នាក់"
        />
      </div>
       <AppInput clearable label="Map Link" v-model.trim="form.map_link" size="large" placeholder="https://maps.google.com/..."></AppInput>
      <div class="form-row">
        <AppInput type="number" clearable label="ចម្ងាយអាចស្កែនបាន (ម៉េត្រ)" v-model.number="form.radius" size="large" placeholder="e.g. 100"></AppInput>
        <el-form-item label="អាចស្កែនក្រៅតំបន់" prop="can_scan_outsize">
          <el-radio-group v-model="form.can_scan_outsize" size="large">
            <el-radio-button :value="true">អាចស្កែនបាន</el-radio-button>
            <el-radio-button :value="false">មិនអាចស្កែនបាន</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton @click="handleSave" type="primary" :loading="saving" size="large" :block="false">
        {{ isEdit ? "កែប្រែ" : "បង្កើត" }}
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { createClass, updateClass } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import { useNotification } from "../composables/useNotification.js";
import AppInput from "./AppInput.vue";
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  isEdit: { type: Boolean, default: false },
  editRow: { type: Object, default: null },
  majors: { type: Array, default: () => [] },
  shifts: { type: Array, default: () => [] },
  generations: { type: Array, default: () => [] },
  programmes: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const formRef = ref();
const saving = ref(false);

const classType = [
  { value: "onclass", label: "ផ្ទាល់" },
  { value: "online", label: "អនឡាញ" },
];

const form = reactive({
  name: "",
  type: "",
  major_id: null,
  shift_id: null,
  generation_id: null,
  programme_id: null,
  year: null,
  semester: null,
  group: null,
  term: null,
  map_link: "",
  radius: "",
  can_scan_outsize: false,
});

const rules = {
  name: [{ required: true, message: "Class name is required" }],
  major_id: [{ required: true, message: "Major is required" }],
  shift_id: [{ required: true, message: "Shift is required" }],
  generation_id: [{ required: true, message: "Generation is required" }],
  programme_id: [{ required: true, message: "Programme is required" }],
  year: [{ required: true, message: "Year is required" }],
  semester: [{ required: true, message: "Semester is required" }],
  map_link: [{ required: true, message: "Map link is required" }],
  radius: [{ required: true, message: "Radius is required" }],
  can_scan_outsize: [{ required: true, message: "Employee can scan outside or not" }],
};

function resetForm() {
  form.name = "";
  form.type = "";
  form.major_id = null;
  form.shift_id = null;
  form.generation_id = null;
  form.programme_id = null;
  form.year = "";
  form.semester = "";
  form.group = "";
  form.term = "";
  form.map_link = "";
  form.radius = "";
  form.can_scan_outsize = null;
}

// Repopulate whenever the dialog opens
watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    if (props.isEdit && props.editRow) {
      const row = props.editRow;
      Object.assign(form, {
        name: row.name || "",
        type: row.type || "",
        major_id: row.major_id ?? null,
        shift_id: row.shift_id ?? null,
        generation_id: row.generation_id ?? null,
        programme_id: row.programme_id ?? null,
        year: row.year || "",
        semester: row.semester || "",
        group: row.group || "",
        term: row.term || "",
        map_link: "",
        radius: row.radius || "",
        can_scan_outsize: row.can_scan_outsize ?? null,
      });
    } else {
      resetForm();
    }
  },
);

async function handleSave() {
  await formRef.value.validate();
 
  try {
    if (props.isEdit) {
      const payload = {
        name: form.name,
        type: form.type,
        major_id: form.major_id,
        shift_id: form.shift_id,
        generation_id: form.generation_id,
        programme_id: form.programme_id,
        year: form.year,
        semester: form.semester,
        group: form.group,
        term: form.term,
        can_scan_outsize: form.can_scan_outsize,
      };
      await updateClass(props.editRow.id, payload);
      notify.success("កែប្រែបានជោគជ័យ");
    } else {
      await createClass(form);
      notify.success("បង្កើតថ្នាក់បានជោគជ័យ");
    }
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
   
  }
}
</script>

<style scoped>
.form-row {
  padding: 10px;
  display: flex;
  gap: 16px;
}
.form-row .el-form-item {
  flex: 1;
  min-width: 0;
}
</style>
