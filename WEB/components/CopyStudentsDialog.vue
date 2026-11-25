<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`Copy សិស្សពី ${sourceClass?.name || ''}`"
    width="60%"
  >
    <el-form label-position="top">
      <AppSelect
        v-model="copyTargetClassId"
        label="ថ្នាក់គោលដៅ"
        placeholder="ជ្រើសរើសថ្នាក់"
        filterable
        :options="targetOptions"
        size="large"
      />

      <el-form-item label="ជ្រើសរើសសិស្ស">
        <AppTable
          ref="tableRef"
          :data="copyStudents"
          selectable
          :columns="studentColumns"
          :show-pagination="false"
          @selection-change="onSelectionChange"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <AppButton
        @click="$emit('update:modelValue', false)"
        size="large"
        :block="false"
        type="warning"
      >
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
import { ref, computed, watch, nextTick } from "vue";
import { adduserclass } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import { useNotification } from "../composables/useNotification.js";
import AppSelect from "./AppSelect.vue";
import AppTable from "./AppTable.vue";

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
const selectedStudents = ref([]);
const tableRef = ref(null);

const studentColumns = [
  { prop: "name_kh", label: "ឈ្មោះខ្មែរ" },
  { prop: "name_en", label: "ឈ្មោះឡាតាំង" },
  { prop: "code", label: "អត្តលេខ" },
];

const copySelectedCount = computed(() => selectedStudents.value.length);

function onSelectionChange(rows) {
  selectedStudents.value = rows;
}

watch(
  () => props.modelValue,
  async (open) => {
    if (!open) return;
    copyStudents.value = [...(props.sourceClass?.students || [])];
    selectedStudents.value = [];
    copyTargetClassId.value = null;
    await nextTick();
    await nextTick();
    tableRef.value?.toggleAllSelection?.();
  },
);

async function handleCopySubmit() {
  const selected = selectedStudents.value;
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