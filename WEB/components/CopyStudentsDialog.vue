<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`Copy សិស្សពី ${sourceClass?.name || ''}`"
    width="60%"
  >
    <el-form label-position="top">
      <el-form-item label="ថ្នាក់គោលដៅ" required>
        <el-select
          v-model="copyTargetClassId"
          placeholder="ជ្រើសរើសថ្នាក់"
          filterable
          size="large"
          style="width: 100%"
        >
          <el-option v-for="c in targetOptions" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
      </el-form-item>

      <el-form-item label="ជ្រើសរើសសិស្ស">
        <el-checkbox v-model="copySelectAll" @change="handleCopySelectAll" style="margin-bottom: 8px">
          ជ្រើសរើសទាំងអស់
        </el-checkbox>
        <el-table :data="copyStudents" size="small" border max-height="320">
          <el-table-column width="50">
            <template #default="{ row }">
              <el-checkbox v-model="row._checked" />
            </template>
          </el-table-column>
          <el-table-column prop="name_kh" label="ឈ្មោះខ្មែរ" />
          <el-table-column prop="name_en" label="ឈ្មោះឡាតាំង" />
          <el-table-column prop="code" label="អត្តលេខ" width="100" />
        </el-table>
      </el-form-item>
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton
        @click="handleCopySubmit"
        type="primary"
        :loading="copying"
        size="large"
        :block="false"
        :disabled="!copyTargetClassId || !copySelectedCount"
      >
        Copy ({{ copySelectedCount }})
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { adduserclass } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import { useNotification } from "../composables/useNotification.js";
import AppSelect from "./AppSelect.vue";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  sourceClass: { type: Object, default: null },
  targetOptions: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const copying = ref(false);
const copyTargetClassId = ref(null);
const copyStudents = ref([]);
const copySelectAll = ref(true);

const copySelectedCount = computed(() => copyStudents.value.filter((s) => s._checked).length);

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    copyStudents.value = (props.sourceClass?.students || []).map((s) => ({ ...s, _checked: true }));
    copyTargetClassId.value = null;
    copySelectAll.value = true;
  },
);

function handleCopySelectAll(val) {
  copyStudents.value.forEach((s) => (s._checked = val));
}

async function handleCopySubmit() {
  const selected = copyStudents.value.filter((s) => s._checked);
  if (!selected.length || !copyTargetClassId.value) return;

  copying.value = true;
  try {
    await adduserclass({
      class_id: copyTargetClassId.value,
      user_id: selected.map((s) => ({ user_id: s.id })),
    });
    notify.success("Copy ទិន្ន័យសិស្សបានជោគជ័យ");
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "Copy បរាជ័យ");
  } finally {
    copying.value = false;
  }
}
</script>
