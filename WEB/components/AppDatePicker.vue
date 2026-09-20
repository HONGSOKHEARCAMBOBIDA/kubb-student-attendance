<script setup>
defineProps({
  modelValue: { type: [String, Date, Array], default: null },
  label: { type: String, default: '' },
  prop: { type: String, default: '' },
  type: { type: String, default: 'date' }, // date | daterange | datetime | month
  placeholder: { type: String, default: 'Select date' },
  format: { type: String, default: 'YYYY-MM-DD' },
  valueFormat: { type: String, default: 'YYYY-MM-DD' },
  clearable: { type: Boolean, default: true },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue', 'change'])
</script>

<template>
  <el-form-item :prop="prop" :label="label" class="app-date-picker mb-4">
    <el-date-picker
      :model-value="modelValue"
      :type="type"
      :placeholder="placeholder"
      :start-placeholder="type.includes('range') ? 'Start date' : undefined"
      :end-placeholder="type.includes('range') ? 'End date' : undefined"
      :format="format"
      :value-format="valueFormat"
      :clearable="clearable"
      :disabled="disabled"
      class="app-date-picker-input"
      @update:model-value="$emit('update:modelValue', $event)"
      @change="$emit('change', $event)"
    />
  </el-form-item>
</template>
<style >
.app-date-picker {
  width: 100%;
}

.app-date-picker :deep(.el-form-item__content) {
  width: 100%;
}

.app-date-picker-input {
  width: 100% !important;
}
</style>