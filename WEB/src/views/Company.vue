<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 4 },
        { slot: 'major', span: 4 },
        { slot: 'shift', span: 4 },
        { slot: 'generation', span: 4 },
        { slot: 'programm', span: 4 },
      ]"
    >
      <template #name>
        <AppInput v-model.trim="filters.name" placeholder="ស្វែងរកតាមឈ្មោះថ្នាក់" clearable size="large"></AppInput>
      </template>

      <template #major>
        <AppSelect v-model="filters.major_id" placeholder="ជំនាញ" clearable size="large" :options="majors"></AppSelect>
      </template>

      <template #shift>
        <AppSelect v-model="filters.shift_id" placeholder="វេន" clearable size="large" :options="shifts"></AppSelect>
      </template>

      <template #generation>
        <AppSelect v-model="filters.generation_id" placeholder="ជំនាន់" clearable size="large" :options="generations"></AppSelect>
      </template>

      <template #programm>
        <AppSelect v-model="filters.programme_id" placeholder="កម្មវិធីសិក្សា" clearable size="large" :options="programmes"></AppSelect>
      </template>

      <template #actions>
        <AppButton v-if="canAddClass" type="primary" @click="openCreate" :block="false">
          បន្ថែមថ្នាក់
        </AppButton>
      </template>
    </AppFilterBar>

    <ClassesMainTable
      :classes="classes"
      :loading="loading"
      v-model:page="page"
      v-model:page-size="pageSize"
      :total="total"
      :can-edit-class="canEditClass"
      :admin-level="adminLevel"
      @page-change="fetchClasses"
      @edit="openEdit"
      @edit-telegram="openEditTelegram"
      @toggle-status="handleToggleStatus"
      @import="openImport"
      @copy="openCopy"
      @schedule="openSchedule"
      @score-import="openScoreImport"
      @score-view="openScoreView"
      @edit-student="openEditStudent"
      @edit-status="openUserClassStatusEditor"
      @score-student="openScoreDialog"
    />

    <ClassFormDialog
      v-model="dialogVisible"
      :is-edit="isEdit"
      :edit-row="editRow"
      :majors="majors"
      :shifts="shifts"
      :generations="generations"
      :programmes="programmes"
      @saved="fetchClasses"
    />

    <ClassTelegramDialog
      v-model="dialogTelegramVisible"
      :class-id="editRow?.id"
      :edit-row="editRow"
      @saved="fetchClasses"
    />

    <ImportStudentsDialog
      v-model="importDialog"
      :class-id="importClassId"
      :class-name="importClassName"
      :generation-name="importGenerationName"
      @saved="fetchClasses"
    />

    <CopyStudentsDialog
      v-model="copyDialog"
      :source-class="copySourceClass"
      :target-options="copyTargetOptions"
      @saved="fetchClasses"
    />

    <ClassScheduleDialog v-model="scheduleDialog" :class-row="scheduleClass" />

    <StudentEditDialog v-model="studentDialogVisible" :student="editStudent" @saved="fetchClasses" />

    <UserClassStatusDialog v-model="userClassStatusDialogVisible" :student="editStatusStudent" @saved="fetchClasses" />

    <ScoreEntryDialog
      v-model="scoreDialogVisible"
      :class-row="scoreClassRow"
      :student="scoreStudent"
      :grade-components="gradeComponents"
      @saved="fetchClasses"
    />

    <ScoreImportDialog v-model="scoreImportDialog" :class-row="scoreImportClassRow" @saved="fetchClasses" />

    <ScoreViewDialog v-model="scoreViewDialog" :class-row="scoreViewClassRow" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getClass,
  changeStatusClass,
  getMajor,
  getShift,
  getGeneration,
  getProgramme,
  getGradecomponent,
} from "../api/services";
import AppButton from "../../components/AppButton.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import ClassScheduleDialog from "../../components/ClassScheduleDialog.vue";
import { useNotification } from "../../composables/useNotification.js";

import ClassesMainTable from "../../components/ClassesMainTable.vue";
import ClassFormDialog from "../../components/ClassFormDialog.vue";
import ClassTelegramDialog from "../../components/ClassTelegramDialog.vue";
import ImportStudentsDialog from "../../components/ImportStudentsDialog.vue";
import CopyStudentsDialog from "../../components/CopyStudentsDialog.vue";
import StudentEditDialog from "../../components/StudentEditDialog.vue";
import UserClassStatusDialog from "../../components/UserClassStatusDialog.vue";
import ScoreEntryDialog from "../../components/ScoreEntryDialog.vue";
import ScoreImportDialog from "../../components/ScoreImportDialog.vue";
import ScoreViewDialog from "../../components/ScoreViewDialog.vue";
import AppInput from "../../components/AppInput.vue";
import AppSelect from "../../components/AppSelect.vue";

const notify = useNotification();
const userDataStore = useUserDataStore();

/* ---------------- lookups & list ---------------- */

const classes = ref([]);
const majors = ref([]);
const shifts = ref([]);
const generations = ref([]);
const programmes = ref([]);
const gradeComponents = ref([]);

const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const filters = reactive({
  name: "",
  major_id: "",
  shift_id: "",
  generation_id: "",
  programme_id: "",
});

const adminLevel = computed(() => userDataStore.level === 7);
const canAddClass = computed(() => userDataStore.permissions?.some((p) => p.name === "add.company"));
const canEditClass = computed(() => userDataStore.permissions?.some((p) => p.name === "edit.company"));

async function fetchClasses() {
  loading.value = true;
  try {
    const res = await getClass({
      page: page.value,
      page_size: pageSize.value,
      name: filters.name || undefined,
      major_id: filters.major_id || undefined,
      shift_id: filters.shift_id || undefined,
      generation_id: filters.generation_id || undefined,
      programme_id: filters.programme_id || undefined,
    });
    classes.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

async function fetchLookups() {
  try {
    const [majorRes, shiftRes, generationRes, programmeRes] = await Promise.all([
      getMajor(),
      getShift(),
      getGeneration(),
      getProgramme(),
    ]);
    majors.value = (majorRes.data.data || []).map((s)=> ({
      label: `${s.name_kh}`,
      value: s.id
    }))
    shifts.value = (shiftRes.data.data || []).map((s)=> ({
      label: `${s.name}`,
      value: s.id
    }))
    generations.value = (generationRes.data.data || []).map((s)=> ({
      label: `${s.name_kh}`,
      value: s.id
    }))
    programmes.value = (programmeRes.data.data || []).map((s)=> ({
      label: `${s.name}`,
      value: s.id
    }))
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

async function fetchGradeComponents() {
  try {
    const res = await getGradecomponent();
    gradeComponents.value = res.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចទាញយក Grade Component បានទេ");
  }
}

watch(
  () => filters.name,
  () => {
    page.value = 1;
    debounce(() => fetchClasses());
  },
);

watch(
  [() => filters.major_id, () => filters.shift_id, () => filters.generation_id, () => filters.programme_id],
  () => {
    page.value = 1;
    fetchClasses();
  },
);

let debounceTimer = null;
function debounce(fn, delay = 400) {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fn, delay);
}

/* ---------------- class create / edit ---------------- */

const dialogVisible = ref(false);
const isEdit = ref(false);
const editRow = ref(null);

function openCreate() {
  isEdit.value = false;
  editRow.value = null;
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editRow.value = row;
  dialogVisible.value = true;
}

/* ---------------- telegram ---------------- */

const dialogTelegramVisible = ref(false);

function openEditTelegram(row) {
  editRow.value = row;
  dialogTelegramVisible.value = true;
}

/* ---------------- toggle status ---------------- */

async function handleToggleStatus(row) {
  try {
    await changeStatusClass(row.id);
    notify.success("ផ្លាស់ប្តូរស្ថានភាពបានជោគជ័យ");
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

/* ---------------- excel import (students) ---------------- */

const importDialog = ref(false);
const importClassId = ref(null);
const importClassName = ref("");
const importGenerationName = ref("")
function openImport(row) {
  importClassId.value = row.id;
  importClassName.value = row.name;
  importGenerationName.value = row.generation_name;
  importDialog.value = true;
}

/* ---------------- copy students ---------------- */

const copyDialog = ref(false);
const copySourceClass = ref(null);
const copyTargetOptions = computed(() => classes.value.filter((c) => c.id !== copySourceClass.value?.id));

function openCopy(row) {
  copySourceClass.value = row;
  copyDialog.value = true;
}

/* ---------------- schedule ---------------- */

const scheduleDialog = ref(false);
const scheduleClass = ref(null);

function openSchedule(row) {
  scheduleClass.value = row;
  scheduleDialog.value = true;
}

/* ---------------- student edit ---------------- */

const studentDialogVisible = ref(false);
const editStudent = ref(null);

function openEditStudent(row) {
  editStudent.value = row;
  studentDialogVisible.value = true;
}

/* ---------------- user-class status ---------------- */

const userClassStatusDialogVisible = ref(false);
const editStatusStudent = ref(null);

function openUserClassStatusEditor(row) {
  editStatusStudent.value = row;
  userClassStatusDialogVisible.value = true;
}

/* ---------------- score entry ---------------- */

const scoreDialogVisible = ref(false);
const scoreClassRow = ref(null);
const scoreStudent = ref(null);

function openScoreDialog(classRow, student) {
  scoreClassRow.value = classRow;
  scoreStudent.value = student;
  scoreDialogVisible.value = true;
}

/* ---------------- score import ---------------- */

const scoreImportDialog = ref(false);
const scoreImportClassRow = ref(null);

function openScoreImport(row) {
  scoreImportClassRow.value = row;
  scoreImportDialog.value = true;
}

/* ---------------- score view ---------------- */

const scoreViewDialog = ref(false);
const scoreViewClassRow = ref(null);

function openScoreView(row) {
  scoreViewClassRow.value = row;
  scoreViewDialog.value = true;
}

onMounted(() => {
  fetchClasses();
  fetchLookups();
  fetchGradeComponents();
});
</script>
