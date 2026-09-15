<!-- <template>
  <div>

    <AppFilterBar
      :fields="[
        { slot: 'name', span: 10 },
        { slot: 'date', span: 6 },
      ]"
      :action-span="2"
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
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'gender_string', label: 'ភេទ', minWidth: 90 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 110 },
          { prop: 'company_name', label: 'ក្រុមហ៑ុន', minWidth: 130 },
          { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]"
      >
        <template #status="{ row }">
          <el-tag
            :type="row.status === 'COMPLETE' ? 'success' : 'warning'"
            size="small"
          >
            {{ row.status }}
          </el-tag>
        </template>

        <template #actions="{ row }">
          <div class="detail-action">
            <el-badge
              :value="row.attendance_record?.length || 0"
              :type="row.status === 'COMPLETE' ? 'success' : 'warning'"
            >
              <el-button
                size="large"
                icon="List"
                circle
                @click="viewRecords(row)"
              />
            </el-badge>
          </div>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="recordsDialog"
      :title="`លម្អិត — ${selectedRow?.name}`"
      width="60%"
    >
      <AppTable
        :data="selectedRow?.attendance_record || []"
        :total="selectedRow?.attendance_record.length || []"
        :columns="[
          { prop: 'day_string', label: 'ថ្ងៃ', width: 100 },
          { prop: 'type_string', label: 'ប្រភេទ', width: 150 },
          { prop: 'check_time', label: 'ម៉ោងបានស្កែន', width: 150 },
          { prop: 'scheduled_time', label: 'ម៉ោងត្រូវស្កែន', width: 110 },
          { prop: 'time_diff', label: 'យឺត/មុនម៉ោង', width: 150 },
          { label: 'ស្ថានភាព', slot: 'attendance_type_name', width: 100 },
          { label: 'ទីតាំងស្កែន', slot: 'check_location' },
          { label: 'មេីលទីតាំង', slot: 'view_location' },
          { prop: 'reason', label: 'មូលហេតុ', width: 120 },
        ]"
      >
        <template #attendance_type_name="{ row }">
          <el-tag size="small" :type="getAttendTypeTag(row.attendance_type)">
            {{ row.attendance_type_name || "—" }}
          </el-tag>
        </template>
        <template #check_location="{ row }">
          <el-icon size="large" :color="row.inzone ? '#67c23a' : '#f56c6c'">
            <component
              :is="row.inzone ? 'CircleCheckFilled' : 'CircleCloseFilled'"
            />
          </el-icon>
        </template>
        <template #view_location="{ row }">
          <el-icon
            size="large"
            :color="row.inzone ? '#67c23a' : '#f56c6c'"
            style="cursor: pointer"
            @click="openLocation(row)"
          >
            <LocationInformation />
          </el-icon>
        </template>
      </AppTable>
    </AppDialog>
  </div>
</template> -->
<template>
  <div>
    <!-- Filters -->
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 10 },
        { slot: 'date', span: 6 },
        {slot: 'company',span:5}
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
      <template #company>
        <el-select v-model="filters.company_id" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%" size="large"
          @change="fetchAttendance">
          <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
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
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'gender_string', label: 'ភេទ', minWidth: 80 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 100 },
          { prop: 'company_name', label: 'ក្រុមហ៑ុន', minWidth: 120 },
          { prop: 'check_date', label: 'ថ្ងៃស្កែន', minWidth: 110 },
          { label: 'ចូលព្រឹក', slot: 'check_in1', minWidth: 110 },
          { label: 'ចេញថ្ងៃត្រង់', slot: 'check_out1', minWidth: 110 },
          { label: 'ចូលថ្ងៃត្រង់', slot: 'check_in2', minWidth: 110 },
          { label: 'ចេញល្ងាច', slot: 'check_out2', minWidth: 110 },
          { prop: 'reason', label: 'មូលហេតុ', minWidth: 130 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]"
      >
        <template #status="{ row }">
          <el-tag
            :type="row.status === 'COMPLETE' ? 'success' : 'warning'"
            size="small"
          >
            {{ row.status }}
          </el-tag>
        </template>

        <template #check_in1="{ row }">
          <CheckCell :time="row.check_in1" :diff="row.check_in1_diff" />
        </template>
        <template #check_out1="{ row }">
          <CheckCell :time="row.check_out1" :diff="row.check_out1_diff" />
        </template>
        <template #check_in2="{ row }">
          <CheckCell :time="row.check_in2" :diff="row.check_in2_diff" />
        </template>
        <template #check_out2="{ row }">
          <CheckCell :time="row.check_out2" :diff="row.check_out2_diff" />
        </template>
        <template #actions="{row}">
           <AppButton
                v-show="false"
                v-if=candeleteattendance
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
<!-- <script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { createAttendance, getAttendance } from "../api/services";
import { LocationInformation } from "@element-plus/icons-vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";
import AppTable from "../../components/AppTable.vue";
import AppDialog from "../../components/AppDialog.vue";
const attendance = ref([]);
const loading = ref(false);
const checking = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const filters = reactive({ 
name: "", 
check_date: new Date().toISOString().split("T")[0],  
company_id: "" });
const recordsDialog = ref(false);
const selectedRow = ref(null);

const attendForm = reactive({ latitude: "", longitude: "", reason: "" });

function getAttendTypeTag(type) {
  const map = { 1: "success", 2: "info", 3: "warning", 4: "danger" };
  return map[type] || "info";
}

function openLocation(row) {
  if (!row.latitdude || !row.longitude) {
    return ElMessage.warning("មិនមានទីតាំងស្កែន");
  }
  const url = `https://www.google.com/maps?q=${row.latitdude},${row.longitude}`;
  window.open(url, "_blank");
}

function getLocation() {
  if (!navigator.geolocation)
    return ElMessage.warning("Geolocation not supported");
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      attendForm.latitude = String(pos.coords.latitude);
      attendForm.longitude = String(pos.coords.longitude);
      ElMessage.success("Location detected");
    },
    () => ElMessage.error("Failed to get location"),
  );
}

async function handleCheckIn() {
  if (!attendForm.latitude || !attendForm.longitude) {
    return ElMessage.warning("Please provide latitude and longitude");
  }
  checking.value = true;
  try {
    await createAttendance(attendForm);
    ElMessage.success("Attendance recorded successfully");
    fetchAttendance();
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to record attendance");
  } finally {
    checking.value = false;
  }
}

async function fetchAttendance() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.company_id) params.company_id = filters.company_id;
    const res = await getAttendance(params);
    attendance.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load attendance");
  } finally {
    loading.value = false;
  }
}

function viewRecords(row) {
  selectedRow.value = row;
  recordsDialog.value = true;
}

onMounted(fetchAttendance);
</script> -->

<script setup>
import { ref, reactive, onMounted, h, computed } from "vue";
import { ElMessage } from "element-plus";
import { getAttendance,exportAttendancePDF, getCompany,deleteattendance } from "../api/services";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppButton from "../../components/AppButton.vue";
import AppTable from "../../components/AppTable.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useUserDataStore } from "../stores/user_data.js";
const userDataStore = useUserDataStore();
const notify = useNotification()
const attendance = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const companys = ref([]);
const filters = reactive({
  name: "",
  check_date: new Date().toISOString().split("T")[0],
  company_id: "",
});

const candeleteattendance = computed(()=>userDataStore.permissions?.some((p)=> p.name === "delete.backup"))

// small functional component to render time + colored diff badge underneath
const CheckCell = (props) => {
  if (!props.time) {
    return h("span", { style: "color:#c0c4cc" }, "—");
  }
let color = "#909399"; 
if (props.diff?.includes("យឺត")) {
  color = "#f56c6c"; // Red
} else if (props.diff?.includes("មុនម៉ោង")) {
  color = "#e6a23c"; // Orange
} else if (props.diff?.includes("បន្ថែម")) {
  color = "#67c23a"; // Green
}

  return h("div", [
    h("div", { style: "font-weight:600" }, props.time),
    props.diff
      ? h("div", { style: `font-size:14px;color:${color}` }, props.diff)
      : null,
  ]);
};
CheckCell.props = ["time", "diff"];

async function fetchCompany() {
  loading.value = true;
  try {
    const res = await getCompany();
    companys.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

async function deleteattendancev1(row) {
  try {
    await deleteattendance(row.id);
    notify.success("Delete success");
    await fetchAttendance()
  }catch(e){
 notify.error(
      e.response?.data?.error || "មិនអាចធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានទេ",
    );
  }
}

async function fetchAttendance() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.check_date) params.check_date = filters.check_date;
    if (filters.company_id) params.company_id = filters.company_id;
    const res = await exportAttendancePDF(params);
    attendance.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load attendance");
  } finally {
    loading.value = false;
  }
}

onMounted(()=>{
  fetchAttendance();
  fetchCompany()
});
</script>

<style scoped>
.checkin-card {
  border-radius: 6px;
}
.filter-card {
  border-radius: 6px;
}
.card-title {
  font-weight: 600;
  font-size: 15px;
}
.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.detail-action {
  padding: 10px 0;
}
</style>
