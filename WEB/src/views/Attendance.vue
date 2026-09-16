<template>
  <div>
    <!-- Filters -->
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 10 },
        { slot: 'date', span: 6 },
        { slot: 'class', span: 5 }
      ]"
      :action-span="3"
    >
      <template #name>
        <el-input
          v-model="filters.name"
          placeholder="ស្វែងរក"
          prefix-icon="Search"
          clearable
          @change="fetchAttendance"
          size="large"
        />
      </template>
      <template #date>
        <el-date-picker
          v-model="filters.check_date"
          type="date"
          placeholder="ជ្រេីសរេីសថ្ងៃទី"
          value-format="YYYY-MM-DD"
          clearable
          @change="fetchAttendance"
          style="width: 100%"
          size="large"
        />
      </template>
      <template #class>
        <el-select
          v-model="filters.class_id"
          placeholder="ថ្នាក់"
          clearable
          style="width: 100%"
          size="large"
          @change="fetchAttendance"
        >
          <el-option
            v-for="cls in classes"
            :key="cls.id"
            :label="cls.name"
            :value="cls.id"
          />
        </el-select>
      </template>
      <template #actions>
        <AppButton type="primary" @click="fetchAttendance"> ស្វែងរក </AppButton>
      </template>
    </AppFilterBar>

    <el-card>
      <AppTable
        :data="attendance"
        :loading="loading"
        show-index
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchAttendance"
        :columns="[
          { prop: 'name_kh', label: 'ឈ្មោះ', minWidth: 130 },
          { prop: 'code', label: 'កូដ', minWidth: 90 },
          { label: 'ភេទ', slot: 'gender', minWidth: 80 },
          { prop: 'class_name', label: 'ថ្នាក់', minWidth: 110 },
          { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 110 },
          { label: 'ម៉ោងទី ១', slot: 'session1', minWidth: 100 },
          { label: 'ម៉ោងទី ២', slot: 'session2', minWidth: 100 },
          { label: 'ម៉ោងទី ៣', slot: 'session3', minWidth: 100 },
          { label: 'ម៉ោងទី ៤', slot: 'session4', minWidth: 100 },
          { label: 'ម៉ោងទី ៥', slot: 'session5', minWidth: 100 },
          { prop: 'reason', label: 'មូលហេតុ', minWidth: 100 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]"
      >
        <template #gender="{ row }">
          {{ genderLabel(row.gender) }}
        </template>

        <template #status="{ row }">
          <el-tag
            :type="row.status === 'COMPLETE' ? 'success' : 'warning'"
            size="small"
          >
            {{ row.status === 'COMPLETE' ? 'ចេញពីរឿន' : 'កំពុងរៀន'  }}
          </el-tag>
        </template>

        <template #session1="{ row }"><CheckCell :time="row.session1" /></template>
        <template #session2="{ row }"><CheckCell :time="row.session2" /></template>
        <template #session3="{ row }"><CheckCell :time="row.session3" /></template>
        <template #session4="{ row }"><CheckCell :time="row.session4" /></template>
        <template #session5="{ row }"><CheckCell :time="row.session5" /></template>

        <template #actions="{ row }">
          <AppButton
            v-if="candeleteattendance"
            size="small"
            icon="Delete"
            type="danger"
            circle
            @click="deleteattendancev1(row)"
          />
        </template>
      </AppTable>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h, computed } from "vue";
import { ElMessage } from "element-plus";
import { getAttendance, exportAttendancePDF, getClass, deleteattendance } from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";
import AppTable from "../../components/AppTable.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useUserDataStore } from "../stores/user_data.js";

const userDataStore = useUserDataStore();
const notify = useNotification();

const attendance = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const classes = ref([]);

const filters = reactive({
  name: "",
  check_date: new Date().toISOString().split("T")[0],
  class_id: "",
});

const candeleteattendance = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "delete.backup")
);

// gender comes back as an int from the backend — adjust mapping to match your DB convention
function genderLabel(gender) {
  if (gender === 1) return "ប្រុស";
  if (gender === 2) return "ស្រី";
  return "—";
}

// backend has no per-session "diff" field, so this just renders the time or a dash
const CheckCell = (props) => {
  if (!props.time) {
    return h("span", { style: "color:#c0c4cc" }, "—");
  }
  return h("div", { style: "font-weight:600" }, props.time);
};
CheckCell.props = ["time"];

async function fetchClasses() {
  try {
    const res = await getClass();
    classes.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load classes");
  }
}

async function deleteattendancev1(row) {
  try {
    await deleteattendance(row.id);
    notify.success("Delete success");
    await fetchAttendance();
  } catch (e) {
    notify.error(e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ");
  }
}

async function fetchAttendance() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.class_id) params.class_id = filters.class_id;

    const res = await exportAttendancePDF(params);
    attendance.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
    console.log(attendance.value)
  } catch {
    ElMessage.error("Failed to load attendance");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchAttendance();
  fetchClasses();
});
</script>

<style scoped>
.checkin-card { border-radius: 6px; }
.filter-card { border-radius: 6px; }
.card-title { font-weight: 600; font-size: 15px; }
.pagination-wrap { margin-top: 16px; display: flex; justify-content: flex-end; }
.detail-action { padding: 10px 0; }
</style>