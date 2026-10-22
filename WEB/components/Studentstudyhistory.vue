<template>
  <AppDialog v-model="visible" title="លំអិតពីការសិក្សា" width="80%">
    <div v-loading="loading" class="history">
      <el-alert
        v-if="errorMsg"
        :title="errorMsg"
        type="error"
        show-icon
        :closable="false"
      />

      <template v-else-if="data">
        <!-- Student profile -->
        <el-descriptions
          title="ព័ត៌មាននិស្សិត"
          :column="3"
          border
          size="default"
        >
          <el-descriptions-item label="ឈ្មោះខ្មែរ">
            {{ dash(data.name_kh) }}
          </el-descriptions-item>
          <el-descriptions-item label="ឈ្មោះឡាតាំង">
            {{ dash(data.name_en) }}
          </el-descriptions-item>
          <el-descriptions-item label="អត្តលេខ">
            {{ dash(data.code) }}
          </el-descriptions-item>
          <el-descriptions-item label="ភេទ">
            {{ data.gender === 1 ? "ប្រុស" : data.gender === 2 ? "ស្រី" : "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="តួនាទី">
            {{ dash(data.role_name) }}
          </el-descriptions-item>
          <el-descriptions-item label="ប្រភេទនិស្សិត">
            {{ dash(data.student_category_name) }}
          </el-descriptions-item>
          <el-descriptions-item label="ថ្ងៃខែឆ្នាំកំណើត">
            {{ dash(data.date_of_birth) }}
          </el-descriptions-item>
          <el-descriptions-item label="សញ្ជាតិ">
            {{ dash(data.nationality) }}
          </el-descriptions-item>
          <el-descriptions-item label="លេខទូរស័ព្ទ">
            {{ dash(data.phone) }}
          </el-descriptions-item>
          <el-descriptions-item label="មុខរបរ">
            {{ dash(data.occupation) }}
          </el-descriptions-item>
          <el-descriptions-item label="ផ្នែកសិក្សា">
            {{ streamLabel(data.academic_stream) }}
          </el-descriptions-item>
          <el-descriptions-item label="Telegram">
            {{ dash(data.telegram_username) }}
          </el-descriptions-item>
          <el-descriptions-item label="ប្រលងចូល">
            <el-tag :type="data.exam_in ? 'success' : 'info'" size="small">
              {{ data.exam_in ? "បាទ/ចាស" : "ទេ" }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="ប្រលងចេញ">
            <el-tag :type="data.exam_out ? 'success' : 'info'" size="small">
              {{ data.exam_out ? "បាទ/ចាស" : "ទេ" }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="ភូមិ">
            {{ dash(data.village) }}
          </el-descriptions-item>
          <el-descriptions-item label="ឃុំ/សង្កាត់">
            {{ dash(data.communce) }}
          </el-descriptions-item>
          <el-descriptions-item label="ស្រុក/ខណ្ឌ">
            {{ dash(data.district) }}
          </el-descriptions-item>
          <el-descriptions-item label="ខេត្ត/រាជធានី">
            {{ dash(data.province) }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- Class history -->
        <div class="section-head">
          <h4>ប្រវត្តិថ្នាក់សិក្សា</h4>
          <el-text type="info">
            សរុប {{ classes.length }} ថ្នាក់ / {{ totalCredits }} ក្រេឌីត
          </el-text>
        </div>

        <el-empty
          v-if="!classes.length"
          description="មិនទាន់មានប្រវត្តិថ្នាក់សិក្សា"
        />

        <el-collapse v-else v-model="openedClasses">
          <el-collapse-item
            v-for="c in classes"
            :key="c.user_class_id"
            :name="c.user_class_id"
          >
            <template #title>
              <div class="class-title">
                <strong>{{ dash(c.class_name) }}</strong>
                <el-text v-if="c.year" type="info">
                  ឆ្នាំទី {{ c.year }} ឆមាសទី {{ c.semester }}
                </el-text>
                <el-tag v-if="c.type" size="large" type="info">
                  {{ c.type === "onclass" ? "សិក្សាផ្ទាល់" : "អនឡាញ" }}
                </el-tag>
                <el-tag
                  size="large"
                  :type="c.active ? 'success' : 'info'"
                  effect="plain"
                >
                  {{ c.active ? "កំពុងសិក្សា" : "បានបញ្ចប់" }}
                </el-tag>
              </div>
            </template>

            <el-descriptions :column="3" size="default" border>
              <el-descriptions-item label="ជំនាញ">
                {{ dash(c.major_name) }}
              </el-descriptions-item>
              <el-descriptions-item label="កម្មវិធីសិក្សា">
                {{ dash(c.programme_name) }}
              </el-descriptions-item>
              <el-descriptions-item label="វេនសិក្សា">
                {{ dash(c.shift_name) }}
              </el-descriptions-item>
              <el-descriptions-item label="ជំនាន់">
                {{ generationLabel(c) }}
              </el-descriptions-item>
              <el-descriptions-item label="ក្រុម">
                {{ dash(c.group_name) }}
              </el-descriptions-item>
              <el-descriptions-item label="វគ្គ">
                {{ dash(c.term) }}
              </el-descriptions-item>
            </el-descriptions>

            <el-table
              :data="c.subject_res || []"
              size="default"
              stripe
              class="subject-table"
              empty-text="មិនមានមុខវិជ្ជា"
              show-summary
              :summary-method="summaryMethod"
            >
              <el-table-column type="index" label="ល.រ" width="60" />
              <el-table-column prop="subject_code" label="កូដ" width="130" />
              <el-table-column
                prop="subject_name_kh"
                label="មុខវិជ្ជា (ខ្មែរ)"
                min-width="180"
              />
              <el-table-column
                prop="subject_name_en"
                label="មុខវិជ្ជា (ឡាតាំង)"
                min-width="180"
              />
              <el-table-column
                prop="credit_hour"
                label="ក្រេឌីត"
                width="100"
                align="center"
              />
            </el-table>
          </el-collapse-item>
        </el-collapse>
      </template>
    </div>

    <template #footer>
      <AppButton
        @click="visible = false"
        size="large"
        :block="false"
        type="warning"
      >
        បិទ
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
// adjust this path to wherever your api/services.js lives
import { getStudentStudyHistory } from "../src/api/services";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  studentRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue"]);

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});

const loading = ref(false);
const errorMsg = ref("");
const data = ref(null);
const openedClasses = ref([]);

const classes = computed(() => data.value?.user_class_res || []);

const totalCredits = computed(() =>
  classes.value.reduce(
    (sum, c) =>
      sum +
      (c.subject_res || []).reduce((s, x) => s + (x.credit_hour || 0), 0),
    0,
  ),
);

const dash = (v) => (v === null || v === undefined || v === "" ? "-" : v);

const streamLabel = (v) =>
  ({ SCIENCE: "វិទ្យាសាស្ត្រពិត", SOCIAL_SCIENCE: "វិទ្យាសាស្ត្រសង្គម" })[v] ||
  "-";

function generationLabel(c) {
  if (!c.generation_name) return "-";
  const range =
    c.generation_start || c.generation_end
      ? ` (${c.generation_start || "?"} - ${c.generation_end || "?"})`
      : "";
  return `${c.generation_name}${range}`;
}

// Total credits row at the bottom of each subject table
function summaryMethod({ columns, data: rows }) {
  return columns.map((col, i) => {
    if (i === 0) return "សរុប";
    if (col.property === "credit_hour") {
      return rows.reduce((s, r) => s + (r.credit_hour || 0), 0);
    }
    return "";
  });
}

async function load() {
  const id = props.studentRaw?.id;
  if (!id) return;
  loading.value = true;
  errorMsg.value = "";
  data.value = null;
  try {
    const res = await getStudentStudyHistory(id);
    data.value = res.data?.data ?? res.data;
    // open every class by default
    openedClasses.value = (data.value?.user_class_res || []).map(
      (c) => c.user_class_id,
    );
  } catch (e) {
    errorMsg.value =
      e.response?.data?.message ||
      e.response?.data?.error ||
      "មិនអាចទាញយកទិន្នន័យបានទេ";
  } finally {
    loading.value = false;
  }
}

// reload every time the dialog opens
watch(
  () => props.modelValue,
  (open) => {
    if (open) load();
  },
);
</script>

<style scoped>
.history {
  min-height: 200px;
  max-height: 70vh;
  overflow-y: auto;
  padding-right: 4px;
}
.section-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin: 24px 0 8px;
}
.section-head h4 {
  margin: 0;
  font-size: 16px;
}
.class-title {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
.subject-table {
  margin-top: 12px;
}
</style>