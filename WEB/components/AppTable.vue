<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import type { TableColumnCtx } from "element-plus";

// Row data
interface RowData {
  [key: string]: any;
}

// Table column
interface Column extends Partial<TableColumnCtx<RowData>> {
  prop?: string;
  label?: string;
  slot?: string;
}

// Props
interface Props {
  data?: RowData[];
  columns: Column[];
  loading?: boolean;
  showIndex?: boolean;
  indexLabel?: string;
  actionsLabel?: string;
  actionsWidth?: string | number;
  mobileBreakpoint?: number;
  showPagination?: boolean;
  currentPage?: number;
  pageSize?: number;
  total?: number;
  tableProps?: Record<string, any>;
  selectable?: boolean;
  expandable?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  loading: false,
  showIndex: false,
  indexLabel: "ល.រ",
  actionsLabel: "សកម្មភាព",
  actionsWidth: 190,
  mobileBreakpoint: 768,
  showPagination: true,
  currentPage: 1,
  pageSize: 10,
  total: 0,
  tableProps: () => ({}),
  selectable: false,
  expandable: false,
});

// Events
const emit = defineEmits<{
  "update:currentPage": [page: number];
  "update:pageSize": [size: number];
  "page-change": [];
  "selection-change": [rows: RowData[]];
}>();

// Pagination
const currentPageProxy = computed({
  get: () => props.currentPage,
  set: (page) => emit("update:currentPage", page),
});

const pageSizeProxy = computed({
  get: () => props.pageSize,
  set: (size) => emit("update:pageSize", size),
});

// Mobile screen
const isMobile = ref(window.innerWidth <= props.mobileBreakpoint);

function checkScreen() {
  isMobile.value = window.innerWidth <= props.mobileBreakpoint;
}

onMounted(() => {
  window.addEventListener("resize", checkScreen);
});

onUnmounted(() => {
  window.removeEventListener("resize", checkScreen);
});

// Selected rows
const selectedRows = ref<RowData[]>([]);

function onSelectionChange(rows: RowData[]) {
  selectedRows.value = rows;
  emit("selection-change", rows);
}

function isSelected(row: RowData) {
  return selectedRows.value.includes(row);
}

function toggleSelection(row: RowData, checked: boolean | string | number) {
  const isChecked = Boolean(checked);

  if (isChecked && !isSelected(row)) {
    selectedRows.value.push(row);
  } else if (!isChecked) {
    selectedRows.value = selectedRows.value.filter((item) => item !== row);
  }

  emit("selection-change", selectedRows.value);
}

// Reset selection when data changes
watch(
  () => props.data,
  () => {
    selectedRows.value = [];
    emit("selection-change", []);
  },
);
</script>

<template>
  <div class="app-table">
    <!-- filter -->
     <div v-if="$slots.filter" class="filter-container">

      <slot name="filter" />
     </div>
    <!-- Pagination -->
    <div v-if="showPagination" class="pagination">
      <span>សរុប {{ total }}</span>

      <el-pagination
        v-model:current-page="currentPageProxy"
        v-model:page-size="pageSizeProxy"
        :total="total"
        :layout="isMobile ? 'prev, next' : 'prev, pager, next'"
        @change="emit('page-change')"
      />
    </div>

    <!-- Desktop table -->
    <el-table
      v-if="!isMobile"
      :data="data"
      :loading="loading"
      border
      size="small"
      v-bind="tableProps"
      @selection-change="onSelectionChange"
    >
      <el-table-column v-if="selectable" type="selection" width="50" />

      <el-table-column
        v-if="showIndex"
        type="index"
        :label="indexLabel"
        width="50"
        align="center"
      />

      <!-- Expand row -->
      <el-table-column v-if="expandable" type="expand">
        <template #default="scope">
          <slot name="expand" v-bind="scope" />
        </template>
      </el-table-column>

      <!-- Dynamic columns -->
      <el-table-column
        v-for="column in columns"
        :key="column.prop || column.label"
        v-bind="column"
      >
        <template v-if="column.slot" #default="scope">
          <slot :name="column.slot" v-bind="scope" />
        </template>
      </el-table-column>

      <!-- Actions -->
      <el-table-column
        v-if="$slots.actions"
        :label="actionsLabel"
        :width="actionsWidth"
      >
        <template #default="scope">
          <slot name="actions" v-bind="scope" />
        </template>
      </el-table-column>
    </el-table>

    <!-- Mobile cards -->
    <div v-else class="mobile-cards" v-loading="loading">
      <el-empty
        v-if="data.length === 0 && !loading"
        description="គ្មានទិន្នន័យ"
      />

      <div
        v-for="(row, index) in data"
        :key="index"
        class="mobile-card"
        :class="{ selected: isSelected(row) }"
      >
        <!-- Checkbox -->
        <div v-if="selectable" class="card-checkbox">
          <el-checkbox
            :model-value="isSelected(row)"
            @change="toggleSelection(row, $event)"
          />
        </div>

        <!-- Card columns -->
        <div
          v-for="column in columns"
          :key="column.prop || column.label"
          class="card-row"
        >
          <span class="card-label">
            {{ column.label }}
          </span>

          <span class="card-value">
            <slot
              v-if="column.slot"
              :name="column.slot"
              :row="row"
              :$index="index"
            />

            <template v-else>
              {{ column.prop ? row[column.prop] : "" }}
            </template>
          </span>
        </div>

        <!-- Actions -->
        <div v-if="$slots.actions" class="card-actions">
          <slot name="actions" :row="row" :$index="index" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.filter-container {
  background: #e5e7e9;
  width: 100%;
}

.app-table {
  width: 100%;
  font-size: 12px;
}

/* Table */
:deep(.el-table) {
  font-size: 12px;
  --el-table-border-color: #d1d5d9;
  --el-table-header-bg-color: #e5e7e9;
  --el-table-row-hover-bg-color: #e8f1fa;
}

:deep(.el-table th.el-table__cell) {
  padding: 0px 0;
  color: #111;
  background: #e5e7e9;
}

:deep(.el-table td.el-table__cell) {
  padding: 0px 0;
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #e5e7e9;
  border: 1px solid #ccc;
  font-size: 12px;
}

/* Mobile cards */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.mobile-card {
  padding: 10px;
  background: white;
  border: 1px solid #ccc;
}

.mobile-card.selected {
  background: #e8f1fa;
  border-color: #7198bc;
}

.card-checkbox {
  margin-bottom: 5px;
}

.card-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  padding: 5px 0;
  border-bottom: 1px solid #eee;
}

.card-label {
  font-weight: bold;
  color: #555;
}

.card-value {
  text-align: right;
  word-break: break-word;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 5px;
  margin-top: 8px;
}
</style>
