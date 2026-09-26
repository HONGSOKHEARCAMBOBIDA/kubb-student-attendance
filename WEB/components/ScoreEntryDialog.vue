<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="បញ្ចូលពិន្ទុ"
    width="900px"
  >
    <el-form label-position="top">
      <div class="form-row">
        <el-form-item label="ថ្នាក់">
          <el-input :model-value="scoreForm.class_name" disabled size="large" />
        </el-form-item>
        <el-form-item label="សិស្ស">
          <el-input :model-value="scoreForm.student_name" disabled size="large" />
        </el-form-item>
      </div>

      <div class="form-row">
        <AppSelect
          v-model="scoreForm.subject_id"
          :options="subjects"
          label="មុខវិជ្ជា"
          placeholder="ជ្រើសរើសមុខវិជ្ជា"
          size="large"
        />
        <el-form-item label="ឆ្នាំ" required>
          <el-input v-model.number="scoreForm.year" type="number" size="large" />
        </el-form-item>
        <el-form-item label="ឆមាស" required>
          <el-input v-model.number="scoreForm.semester" type="number" size="large" />
        </el-form-item>
      </div>

      <el-divider content-position="left"> ពិន្ទុតាមផ្នែក </el-divider>

      <el-table :data="scoreForm.details" border stripe size="default">
        <el-table-column type="index" label="#" width="60" />
        <el-table-column prop="grade_component_name" label="ផ្នែកពិន្ទុ" min-width="220" />
        <el-table-column prop="max_score" label="ពិន្ទុអតិបរមា" width="130" />
        <el-table-column label="ពិន្ទុ" width="180">
          <template #default="{ row }">
            <el-input-number
              v-model="row.score"
              :min="0"
              :max="row.max_score || 100"
              :precision="2"
              :step="0.5"
              size="large"
              style="width: 100%"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false" type="warning">
        បោះបង់
      </AppButton>
      <AppButton @click="handleCreateScore" type="primary" :loading="saving" size="large" :block="false">
        រក្សាទុកពិន្ទុ
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { addScore, getClassAvailableSubjects } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import { useNotification } from "../composables/useNotification.js";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classRow: { type: Object, default: null },
  student: { type: Object, default: null },
  gradeComponents: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const saving = ref(false);
const subjects = ref([]);

const scoreForm = reactive({
  class_id: null,
  class_name: "",
  user_id: null,
  student_name: "",
  major_id: null,
  generation_id: null,
  programme_id: null,
  subject_id: null,
  year: null,
  semester: null,
  details: [],
});

async function fetchSubjectOptions(classId) {
  try {
    const res = await getClassAvailableSubjects(classId);
    subjects.value = (res.data.data || []).map((s) => ({
      label: `${s.code} — ${s.name_kh}`,
      value: s.id,
    }));
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to load subjects");
  }
}

watch(
  () => props.modelValue,
  (open) => {
    if (!open || !props.classRow || !props.student) return;
    const classRow = props.classRow;
    const student = props.student;

    scoreForm.class_id = classRow.id;
    scoreForm.class_name = classRow.name;
    scoreForm.user_id = student.id;
    scoreForm.student_name = `${student.name_kh} (${student.code})`;
    scoreForm.major_id = classRow.major_id;
    scoreForm.generation_id = classRow.generation_id;
    scoreForm.programme_id = classRow.programme_id;
    scoreForm.subject_id = null;
    scoreForm.year = classRow.year;
    scoreForm.semester = classRow.semester;
    scoreForm.details = props.gradeComponents.map((item) => ({
      grade_component_id: item.ID,
      grade_component_name: item.Name,
      max_score: Number(item.WeightPercentage || 100),
      score: 0,
    }));

    fetchSubjectOptions(classRow.id);
  },
);

async function handleCreateScore() {
  if (!scoreForm.class_id) {
    notify.error("មិនមានថ្នាក់");
    return;
  }
  if (!scoreForm.user_id) {
    notify.error("មិនមានសិស្ស");
    return;
  }
  if (!scoreForm.subject_id) {
    notify.error("សូមជ្រើសរើសមុខវិជ្ជា");
    return;
  }

  saving.value = true;
  try {
    const payload = {
      class_id: scoreForm.class_id,
      major_id: scoreForm.major_id,
      generation_id: scoreForm.generation_id,
      programme_id: scoreForm.programme_id,
      subject_id: scoreForm.subject_id,
      year: scoreForm.year,
      semester: scoreForm.semester,
      students: [
        {
          user_id: scoreForm.user_id,
          details: scoreForm.details.map((item) => ({
            grade_component_id: item.grade_component_id,
            score: Number(item.score || 0),
          })),
        },
      ],
    };

    await addScore(payload);
    notify.success("បញ្ចូលពិន្ទុបានជោគជ័យ");
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "បញ្ចូលពិន្ទុមិនបានជោគជ័យ");
  } finally {
    saving.value = false;
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
