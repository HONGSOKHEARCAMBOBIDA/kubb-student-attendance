<template>
  <el-card class="table-card">
    <AppTable
      expandable
      :data="classes"
      :loading="loading"
      show-index
      :current-page="page"
      @update:current-page="$emit('update:page', $event)"
      :page-size="pageSize"
      @update:page-size="$emit('update:pageSize', $event)"
      :total="total"
      @page-change="$emit('page-change')"
      actions-width="380px"
      :columns="[
        { prop: 'name', label: 'ឈ្មោះថ្នាក់', minWidth: 120 },
        { slot: 'major_name', label: 'ជំនាញ', minWidth: 110 },
        { prop: 'shift_name', label: 'វេន', width: 100 },
        { slot: 'generation_name', label: 'ជំនាន់', width: 200 },
        { slot: 'programme_name', label: 'កម្មវិធីសិក្សា', minWidth: 130 },
        { prop: 'year', label: 'ឆ្នាំ', width: 60 },
        { prop: 'semester', label: 'ឆមាស', width: 70 },
        { prop: 'group', label: 'ក្រុម', width: 60 },
        { prop: 'term', label: 'វគ្គ', width: 60 },
        { slot: 'type', label: 'ប្រភេទ', width: 100 },
        { prop: 'radius', label: 'អាចស្កែន (m)', width: 150 },
        { label: 'អាចស្កែនក្រៅតំបន់', slot: 'outsize', width: 150 },
        { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        { label: 'សិស្សសរុប', slot: 'total', width: 100 },
      ]"
    >
      <template #major_name="{ row }">
        <el-text style="color: red">{{ row.major_name }}</el-text>
      </template>
      <template #generation_name="{ row }">
        <el-text style="color: red">{{ row.generation_name }} | {{ row.generation_start }}-{{ row.generation_end }}</el-text>
      </template>
      <template #programme_name="{ row }">
        <el-text style="color: red">{{ row.programme_name }}</el-text>
      </template>
      <template #outsize="{ row }">
        <el-tag :type="row.can_scan_outsize ? 'success' : 'danger'">
          {{ row.can_scan_outsize ? "បាន" : "មិនបាន" }}
        </el-tag>
      </template>
      <template #status="{ row }">
        <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
          {{ row.is_active ? "Active" : "Inactive" }}
        </el-tag>
      </template>
      <template #total="{ row }">
        <el-text>{{ row.students?.length ?? 0 }} នាក់</el-text>
      </template>
      <template #type="{ row }">
        <el-text>{{ row.type === "onclass" ? "ផ្ទាល់" : "អនឡាញ" }}</el-text>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton
            :disabled="row.is_active === false"
            v-if="canEditClass"
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="$emit('edit', row)"
          />
        </el-tooltip>
        <AppButton
          :disabled="row.is_active === false"
          v-if="canEditClass"
          size="small"
          icon="Promotion"
          type="primary"
          circle
          @click="$emit('edit-telegram', row)"
        />
        <el-tooltip content="បិទ/បើក ថ្នាក់" placement="top">
          <AppButton
            v-if="canEditClass"
            size="small"
            icon="Switch"
            :type="row.is_active ? 'danger' : 'success'"
            circle
            @click="$emit('toggle-status', row)"
          />
        </el-tooltip>
        <el-tooltip content="បញ្ជូលសិស្សតាមExcell" placement="top">
          <AppButton
            v-if="adminLevel"
            :disabled="row.is_active === false"
            size="small"
            icon="Download"
            type="success"
            circle
            @click="$emit('import', row)"
          />
        </el-tooltip>
        <el-tooltip content="Copy ទិន្ន័យសិស្សទៅថ្នាក់ផ្សេង" placement="top">
          <AppButton
            v-if="adminLevel"
            :disabled="row.is_active === false"
            size="small"
            icon="CopyDocument"
            type="primary"
            circle
            @click="$emit('copy', row)"
          />
        </el-tooltip>
        <el-tooltip content="កាលវិភាគ" placement="top">
          <AppButton
            v-if="adminLevel"
            :disabled="row.is_active === false"
            size="small"
            icon="Calendar"
            type="primary"
            circle
            @click="$emit('schedule', row)"
          />
        </el-tooltip>
        <el-tooltip content="បញ្ចូលពិន្ទុតាម Excel" placement="top">
          <AppButton
            v-if="adminLevel"
            :disabled="row.is_active === false"
            size="small"
            icon="UploadFilled"
            type="warning"
            circle
            @click="$emit('score-import', row)"
          />
        </el-tooltip>
        <el-tooltip content="មើលពិន្ទុ" placement="top">
          <AppButton
            :disabled="row.is_active === false"
            v-if="adminLevel"
            size="small"
            icon="View"
            type="info"
            circle
            @click="$emit('score-view', row)"
          />
        </el-tooltip>
      </template>

      <template #expand="{ row }">
        <ClassStudentsTable
          :class-row="row"
          :can-edit-class="canEditClass"
          @edit-student="(s) => $emit('edit-student', s)"
          @edit-status="(s) => $emit('edit-status', s)"
          @score-student="(s) => $emit('score-student', row, s)"
        />
      </template>
    </AppTable>
  </el-card>
</template>

<script setup>
import AppTable from "./AppTable.vue";
import AppButton from "./AppButton.vue";
import ClassStudentsTable from "./ClassStudentsTable.vue";

defineProps({
  classes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  total: { type: Number, default: 0 },
  canEditClass: { type: Boolean, default: false },
  adminLevel: { type: Boolean, default: false },
});

defineEmits([
  "update:page",
  "update:pageSize",
  "page-change",
  "edit",
  "edit-telegram",
  "toggle-status",
  "import",
  "copy",
  "schedule",
  "score-import",
  "score-view",
  "edit-student",
  "edit-status",
  "score-student", // emitted as (classRow, student)
]);
</script>

<style scoped>
.table-card {
  border-radius: 6px;
}
</style>
