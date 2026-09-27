<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="`ពិន្ទុថ្នាក់ — ${classRow?.name || ''} ${classRow?.programme_name || ''} ${classRow?.generation_name || ''} ឆ្នាំ${classRow?.year || ''} ឆមាស${classRow?.semester || ''} ជំនាញ${classRow?.major_name || ''}`"
    width="95%"
  >
    <div class="form-row">
      <AppSelect
        v-model="filters.subject_id"
        :options="subjectOptions"
        label="មុខវិជ្ជា"
        placeholder="ជ្រើសរើសមុខវិជ្ជា (ទាំងអស់)"
        clearable
        size="large"
        @change="
          () => {
            page = 1;
            fetchScoreView();
          }
        "
      />
      <AppInput
        v-model.trim="filters.name"
        label="ស្វែងរកឈ្មោះ/អត្តលេខ"
        placeholder="ឈ្មោះ ឬ អត្តលេខ"
        clearable
        size="large"
        @input="debouncedFetch"
      />
    </div>

    <AppTable
      :data="scoreData"
      v-loading="loading"
      :columns="columns"
      :show-pagination="false"
      actions-width="150"
    >
      <template #gender="{ row }">
        <el-text>{{ row.gender === 1 ? "ប្រុស" : "ស្រី" }}</el-text>
      </template>
      <template #total="{ row }">
        <el-text style="color: red;" tab="b" size="large">
          {{
            row.total
          }}
        </el-text>
      </template>
      <template #rank="{row}">
        <el-text >{{ row.rank }}</el-text>
      </template>
    <template #actions="{ row }">
  <el-tooltip content="កែប្រែ" placement="top">
    <AppButton
      type="warning"
      circle
      icon="Edit"
      size="small"
      @click="openEdit(row)"
    ></AppButton>
  </el-tooltip>
  <el-tooltip content="លុប" placement="top">
    <AppButton type="danger" circle icon="Delete" size="small"></AppButton>
  </el-tooltip>
</template>
    </AppTable>
    <el-form label-position="top">
      <AppDialog v-model="editDialogVisible" title="កែប្រែពិន្ទុ" width="40%">

    <el-row :gutter="20">
      <el-col :span="12">
        <AppInput label="វត្តមាននិស្សិត" v-model.number="editForm.attendance" type="number"></AppInput>
      </el-col>
      <el-col :span="12">
        <AppInput label="កិច្ចការស្រាវជ្រាវ" v-model.number="editForm.research" type="number"></AppInput>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <AppInput label="ប្រឡងពាក់កណ្តាលឆមាស" v-model.number="editForm.midterm" type="number"></AppInput>
      </el-col>
      <el-col :span="12">
        <AppInput label="ប្រឡងបញ្ចប់ឆមាស" v-model.number="editForm.final" type="number"></AppInput>
      </el-col>
    </el-row>

  <template #footer>
    <AppButton @click="editDialogVisible = false">បោះបង់</AppButton>
    <AppButton type="primary" :loading="editSaving" @click="submitEdit">
      រក្សាទុក
    </AppButton>
  </template>
</AppDialog>
    </el-form>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
// import { getClassAvailableSubjects, getScore } from "../api/services";
import AppDialog from "./AppDialog.vue";
import AppTable from "./AppTable.vue";
import AppSelect from "./AppSelect.vue";
import AppInput from "./AppInput.vue";
import { getClassAvailableSubjects, getScore,updateScore } from "../src/api/services.js";
import { useNotification } from "../composables/useNotification.js";
import AppButton from "./AppButton.vue";
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classRow: { type: Object, default: null },
});
defineEmits(["update:modelValue"]);

const notify = useNotification();
const loading = ref(false);
const scoreData = ref([]);
const subjectOptions = ref([]);
const page = ref(1);
const pageSize = ref(400);
const total = ref(0);
const filters = reactive({ subject_id: null, name: "" });

const columns = [
  { prop: "name_kh", label: "ឈ្មោះខ្មែរ", minwidth: 100 },
  { slot: "gender", label: "ភេទ", width: 70 },
  { prop: "code", label: "អត្តលេខ", width: 130 },
  { prop: "ProgrammeName", label: "កម្រិត", width: 100 },
  { prop: "GenerationName", label: "ជំនាន់", width: 100 },
  { prop: "MajorName", label: "ជំនាញ", minwidth: 100 },
  { prop: "SubjectName", label: "មុខវិជ្ជា", minwidth: 100 },
  { prop: "year", label: "ឆ្នាំ", width: 70, align: "center" },
  { prop: "semester", label: "ឆមាស", width: 80, align: "center" },
  {
    prop: "attendance",
    label: "វត្តមាននិស្សិត",
    width: 150,
    align: "center",
  },
  {
    prop: "research",
    label: "កិច្ចការស្រាវជ្រាវ",
    width: 150,
    align: "center",
  },
  {
    prop: "midterm",
    label: "ប្រឡងពាក់កណ្តាលឆមាស",
    width: 210,
    align: "center",
  },
  { prop: "final", label: "ប្រឡងបញ្ចប់ឆមាស", width: 150, align: "center" },
  { slot: "total", label: "ពិន្ទុសរុប", width: 100, align: "center" },
  { slot: "rank", label: "ចំណាត់ថ្នាក់", width: 100, align: "center" },
];

let debounceTimer = null;
function debouncedFetch() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    page.value = 1;
    fetchScoreView();
  }, 400);
}

async function fetchScoreView() {
  const row = props.classRow;
  if (!row) return;

  loading.value = true;
  try {
    const res = await getScore({
      page: page.value,
      page_size: pageSize.value,
      class_id: row.id, // class's own primary key = score.class_id
      generation_id: row.generation_id,
      major_id: row.major_id,
      programme_id: row.programme_id,
      year: row.year,
      semester: row.semester,
      subject_id: filters.subject_id || undefined,
      name: filters.name || undefined,
    });
    scoreData.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
    console.log(scoreData.value);
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចទាញយកពិន្ទុបានទេ");
  } finally {
    loading.value = false;
  }
}

const editDialogVisible = ref(false);
const editSaving = ref(false);
const editForm = reactive({
  score_id: null,
  attendance: 0,
  research: 0,
  midterm: 0,
  final: 0,
});

function openEdit(row) {
  editForm.score_id = row.id;
  editForm.attendance = row.attendance || 0;
  editForm.research = row.research || 0;
  editForm.midterm = row.midterm || 0;
  editForm.final = row.final || 0;
  editDialogVisible.value = true;
}


async function submitEdit() {
  editSaving.value = true;
  try {
    await updateScore({
      score_id: editForm.score_id,
      attendance: editForm.attendance,
      research: editForm.research,
      midterm: editForm.midterm,
      final: editForm.final,
    });
    notify.success("កែប្រែពិន្ទុបានជោគជ័យ");
    editDialogVisible.value = false;
    await fetchScoreView();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចកែប្រែពិន្ទុបានទេ");
  } finally {
    editSaving.value = false;
  }
}

watch(
  () => props.modelValue,
  async (open) => {
    if (!open || !props.classRow) return;
    filters.subject_id = null;
    filters.name = "";
    page.value = 1;

    try {
      const res = await getClassAvailableSubjects(props.classRow.id);
      subjectOptions.value = (res.data.data || []).map((s) => ({
        label: `${s.code} — ${s.name_kh}`,
        value: s.id,
      }));
    } catch (e) {
      notify.error(e.response?.data?.error || "Failed to load subjects");
    }

    fetchScoreView();
  },
);
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
