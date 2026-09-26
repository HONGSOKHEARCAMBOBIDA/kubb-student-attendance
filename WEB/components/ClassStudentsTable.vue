<template>
  <div>
    <el-divider content-position="left">
      <el-row :gutter="20">
        <el-col :span="12">
          <AppInput v-model.trim="search" placeholder="ស្វែងរកសិស្ស (ឈ្មោះ ឬ កូដ)" clearable />
        </el-col>
      </el-row>
    </el-divider>

    <AppTable show-index :data="filteredStudents" :columns="studentcolumn" :show-pagination="false">
      <template #gender="{ row }">
        <el-text>{{ row.gender === "1" ? "ប្រុស" : "ស្រី" }}</el-text>
      </template>

      <template #status="{ row }">
        <el-tooltip content="ចុចដើម្បីផ្លាស់ប្តូរស្ថានភាព" placement="top">
          <el-tag
            :type="statusTagType(row.status)"
            size="large"
            style="cursor: pointer"
            @click="canEditClass && $emit('edit-status', row)"
          >
            {{ getUserClassStatusLabel(row.status) }}
          </el-tag>
        </el-tooltip>
      </template>

      <template #actions="{ row: student }">
        <el-tooltip content="បញ្ចូលពិន្ទុ" placement="top">
          <AppButton size="small" icon="EditPen" type="primary" circle @click="$emit('score-student', student)" />
        </el-tooltip>
        <el-tooltip content="កែប្រែសិស្ស" placement="top">
          <AppButton
            v-if="canEditClass"
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="$emit('edit-student', student)"
          />
        </el-tooltip>
      </template>
    </AppTable>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import AppTable from "./AppTable.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";

const props = defineProps({
  classRow: { type: Object, required: true },
  canEditClass: { type: Boolean, default: false },
});
defineEmits(["edit-student", "edit-status", "score-student"]);

const search = ref("");

const filteredStudents = computed(() => {
  const keyword = search.value.trim().toLowerCase();
  const students = props.classRow.students || [];
  if (!keyword) return students;
  return students.filter(
    (s) =>
      s.name_kh?.toLowerCase().includes(keyword) ||
      s.name_en?.toLowerCase().includes(keyword) ||
      s.code?.toLowerCase().includes(keyword),
  );
});

const studentcolumn = [
  { prop: "name_kh", label: "ឈ្មោះខ្មែរ", minwidth: 100 },
  { prop: "name_en", label: "ឈ្មោះអង់គ្លេស", minwidth: 100 },
  { slot: "gender", label: "ភេទ", minwidth: 100 },
  { prop: "code", label: "អត្តលេខ", minwidth: 100 },
  { slot: "status", label: "ស្ថានភាព", minwidth: 100 },
];

const userClassStatus = [
  { value: "STUDY", label: "កំពុងសិក្សា" },
  { value: "SUSPEND", label: "ព្យួរការសិក្សា" },
  { value: "TRANSFER", label: "ដូរ/ផ្ទេរជំនាញ" },
  { value: "DROPPED", label: "បោះបង់ការសិក្សា" },
];

function getUserClassStatusLabel(status) {
  return userClassStatus.find((item) => item.value === status)?.label || status;
}

function statusTagType(status) {
  switch (status) {
    case "STUDY":
      return "success";
    case "SUSPEND":
      return "warning";
    case "TRANSFER":
      return "info";
    case "DROPPED":
      return "danger";
    default:
      return "info";
  }
}
</script>
