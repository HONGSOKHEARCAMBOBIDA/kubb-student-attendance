<template>
  <div class="filter-wrapper">
    <div v-if="showBack" class="back-wrapper">
      <slot name="back" :goback="goback">
        <app-button circle icon="ArrowLeft" @click="goback()" />
      </slot>
    </div>

    <el-row :gutter="2" class="filter-row">
      <el-col
        v-for="(field, i) in normalizedFields"
        :key="field.slot || i"
        :xs="24"
        :sm="12"
        :md="field.span || 6"
      >
        <slot :name="field.slot" />
      </el-col>

      <el-col :xs="24" :sm="12" :md="actionSpan" class="filter-actions">
        <slot name="actions" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useBack } from '../composables/useback.ts';
import AppButton from './AppButton.vue'

const props = defineProps({
  fields: { type: Array, default: () => [] },
  actionSpan: { type: Number, default: 3 },
  showBack: { type: Boolean, default: false },
})

const { goback } = useBack()
const normalizedFields = computed(() => props.fields)
</script>

<style scoped>
.filter-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.back-wrapper {
  flex: 0 0 auto;
}

.filter-row {
  flex: 1;
  min-width: 0;
  row-gap: 6px;
}

.filter-actions {
  display: flex;
  justify-content: flex-end;
  margin-left: auto;
}

@media (max-width: 768px) {
  .filter-wrapper {
    flex-direction: column;
    gap: 6px;
  }

  .filter-row {
    width: 100%;
    row-gap: 1px;
  }

  .filter-actions {
    justify-content: stretch;
  }

  .filter-actions :deep(.el-button) {
    width: 100%;
  }
}
</style>