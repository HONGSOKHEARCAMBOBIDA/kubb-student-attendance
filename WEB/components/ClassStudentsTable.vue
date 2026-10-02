<template>

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
               <el-tooltip content="Transcript" placement="top">
  <AppButton
    v-if="canEditClass"
    size="small"
    icon="View"
    type="success"
    circle
    :loading="loadingTranscript"
    @click="viewTranscript(student)"
  />
</el-tooltip>
      </template>
    </AppTable>
  <AppDialog append-to-body v-model="transcriptDialog" title="Transcript" width="80%">
  <template v-if="transcript">
    <!-- Student info -->
<el-descriptions :column="isMobile ? 1 : 6" border size="small">
  <!-- Row 1 -->
  <el-descriptions-item label="ឈ្មោះ" :span="s(3)">
    {{ transcript.Student.name_kh }} | {{ transcript.Student.name_en }}
  </el-descriptions-item>
  <el-descriptions-item label="ភេទ" :span="s(1)">
    {{ Number(transcript.Student.gender) === 1 ? "ប្រុស" : "ស្រី" }}
  </el-descriptions-item>
  <el-descriptions-item label="លេខចុះឈ្មោះ" :span="s(2)">
    {{ transcript.Student.registration_no }}
  </el-descriptions-item>

  <!-- Row 2 -->
  <el-descriptions-item label="សញ្ជាតិ" :span="s(4)">
    {{ transcript.Student.nationality }}
    <b v-if="transcript.Student.campus">{{ transcript.Student.campus }}</b>
  </el-descriptions-item>
  <el-descriptions-item label="មហាវិទ្យាល័យ" :span="s(2)">
    {{ transcript.Header.faculty_name }}
  </el-descriptions-item>

  <!-- Row 3 -->
  <el-descriptions-item label="ថ្ងៃ-ខែ-ឆ្នាំ កំណើត" :span="s(4)">
    {{ transcript.Student.date_of_birth }}
  </el-descriptions-item>
  <el-descriptions-item label="កម្រិត" :span="s(2)">
    {{ transcript.Header.degree_title }}
  </el-descriptions-item>

  <!-- Row 4 -->
  <el-descriptions-item label="ទីកន្លែងកំណើត" :span="s(4)">
    {{ transcript.Student.place_of_birth }}
  </el-descriptions-item>
  <el-descriptions-item label="ជំនាញ" :span="s(2)">
    {{ transcript.Header.major_name }}
  </el-descriptions-item>

  <!-- Row 5 -->
  <el-descriptions-item label="ឆ្នាំសិក្សា" :span="s(6)">
    {{ transcript.Header.start_year }} - {{ transcript.Header.end_year }}
  </el-descriptions-item>
</el-descriptions>

    <!-- Results grouped by year / semester -->
    <div v-for="group in groupedResults" :key="group.key" style="margin-top: 20px">
      <el-divider content-position="left">
        ឆ្នាំទី {{ group.year }} - ឆមាសទី {{ group.semester }}
      </el-divider>

      <AppTable
        show-index
        :data="group.items"
        :columns="resultColumns"
        :show-pagination="false"
      />

   
      <el-text tag="b" size="large" style="display: block; text-align: right; margin-top: 8px;">
  ក្រេឌីតសរុប: {{ group.totalCredits }} | GPA: {{ group.gpa.toFixed(2) }}
</el-text>
    </div>
<el-divider />
<el-descriptions :column="isMobile ? 1 : 2" border size="small">
  <el-descriptions-item label="ក្រេឌីតសរុបទាំងអស់">
    {{ gpax.totalCredits }}
  </el-descriptions-item>
  <el-descriptions-item label="GPAX/CGPA">
    <el-tag type="success" size="large">{{ gpax.gpax.toFixed(2) }}</el-tag>
  </el-descriptions-item>
</el-descriptions>
    <el-empty v-if="!groupedResults.length" description="គ្មានទិន្នន័យ" />
  </template>
  </AppDialog>
</template>

<script setup>
import { ref, computed } from "vue";
import AppTable from "./AppTable.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";
import { getTranscript } from "../src/api/services.js";
import AppDialog from "./AppDialog.vue";
const transcriptDialog = ref(false);
const transcript = ref(null);
const isMobile = computed(() => window.innerWidth <= 768);
const s = (n) => (isMobile.value ? 1 : n);
const loadingTranscript = ref(false);
const resultColumns = [
  { prop: "code", label: "កូដ", "min-width": 90 },
  { prop: "name_kh", label: "មុខវិជ្ជា (ខ្មែរ)", "min-width": 180 },
  { prop: "name_en", label: "Subject (EN)", "min-width": 180 },
  { prop: "credits", label: "ក្រេឌីត", width: 90, align: "center" },
  { prop: "point_sum", label: "ពិន្ទុ", width: 90, align: "center" },
  { prop: "percent_sum", label: "ភាគរយ", width: 100, align: "center" },
  { prop: "grade", label: "Grade", width: 100, align: "center" },
  { prop: "gp", label: "Grade Point", width: 120, align: "center" },
  { prop: "remark", label: "Remark", width: 100, align: "center" },
];
async function viewTranscript(student) {
  try {
    loadingTranscript.value = true;

    const res = await getTranscript(student.id);

    transcript.value = res.data.data;
    transcriptDialog.value = true;
    console.log(transcript.value)
  } catch (error) {
    console.error(error);
  } finally {
    loadingTranscript.value = false;
  }
}


const getGrade = (score) => {
  score = Number(score);

  if (score >= 85) return { grade: "A", gp: 4.0, remark: "Excellent" };
  if (score >= 80) return { grade: "B+", gp: 3.5, remark: "Very Good" };
  if (score >= 70) return { grade: "B", gp: 3.0, remark: "Good" };
  if (score >= 65) return { grade: "C+", gp: 2.5, remark: "Fairly Good" };
  if (score >= 50) return { grade: "C", gp: 2.0, remark: "Fair" };
  if (score >= 45) return { grade: "D", gp: 1.5, remark: "Poor" };
  if (score >= 40) return { grade: "E", gp: 1.0, remark: "Very Poor" };

  return { grade: "F", gp: 0.0, remark: "Failure" };
};

const groupedResults = computed(() => {
  const results = transcript.value?.Results || [];
  const map = {};

  results.forEach((r) => {
    const key = `${r.year}-${r.semester}`;
    if (!map[key]) {
      map[key] = {
        key,
        year: r.year,
        semester: r.semester,
        items: [],
        totalCredits: 0,
        totalPoints: 0, // Σ (gp × credit)
      };
    }

    const grade = getGrade(r.point_sum);
    const credits = Number(r.credits) || 0;

    map[key].items.push({
      ...r,
      grade: grade.grade,
      gp: grade.gp,
      remark: grade.remark,
    });
    map[key].totalCredits += credits;
    map[key].totalPoints += grade.gp * credits;
  });

  return Object.values(map)
    .map((g) => ({
      ...g,
      gpa: g.totalCredits ? g.totalPoints / g.totalCredits : 0,
    }))
    .sort((a, b) => a.year - b.year || a.semester - b.semester);
});

const gpax = computed(() => {
  const totalCredits = groupedResults.value.reduce((sum, g) => sum + g.totalCredits, 0);
  const totalPoints = groupedResults.value.reduce((sum, g) => sum + g.totalPoints, 0);
  // reduce() = យក Array មកបង្រួមជាតម្លៃមួយ
  return {
    totalCredits,
    gpax: totalCredits ? totalPoints / totalCredits : 0,
  };
});


// Object.values() គឺជា JavaScript Method សម្រាប់យក Values ទាំងអស់ចេញពី Object មកជា Array

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
