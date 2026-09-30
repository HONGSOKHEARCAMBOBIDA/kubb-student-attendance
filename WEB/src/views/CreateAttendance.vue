<template>
  <el-card class="checkin-card">
    <template #header>
      <div class="header">
        <el-image class="logo" src="/logo.png" alt="University Logo" fit="contain" />
        <el-tag size="large" :type="draft?.subject_name ? 'success' : 'warning'">
          {{ draft?.subject_name ?? "វត្តមាន" }}
        </el-tag>
      </div>
    </template>

    <section v-if="draft" class="section">
      <h3 class="title">កាលវិភាគ</h3>
      <div class="card schedule">
        <div class="icon-box">
          <el-icon :size="22"><Sunny /></el-icon>
        </div>
        <div>
          <div class="schedule-time">{{ draft.scheduled_time }}</div>
          <div class="schedule-type">{{ draft.type_string }}</div>
        </div>
      </div>
    </section>

    <section v-if="companies.length" class="section">
      <div
        v-for="c in companies"
        :key="c.id"
        class="card company"
        :class="{ active: attendForm.company_id === c.id }"
        @click="attendForm.company_id = c.id"
      >
        <div class="icon-box1">
          <el-icon :size="22" color="#626aef"><Sunny /></el-icon>
        </div>
        <span class="company-name">{{ c.name }}</span>
        <el-icon class="company-check">
          <CircleCheckFilled v-if="attendForm.company_id === c.id" />
          <ArrowRight v-else />
        </el-icon>
      </div>
    </section>

    <div class="action-buttons">
      <AppButton
        type="default"
        plain
        block
        native-type="button"
        :disabled="isBeforeSchedule"
        @click="$router.push('/leaverequest')"
      >
        <span class="btn-content">
          <el-icon color="#ffb301"><Eleme /></el-icon>
          <span>សុំច្បាប់</span>
        </span>
      </AppButton>

      <AppButton
        plain
        block
        color="#626aef"
        native-type="button"
        :loading="loading"
        :disabled="isButtonDisabled"
        @click="handleCheckIn"
      >
        <span class="btn-content">
          <el-icon><CircleCheck /></el-icon>
          <span>ចុះវត្តមានចូល</span>
        </span>
      </AppButton>
    </div>

    <el-alert v-if="errorMessage" :title="errorMessage" type="warning" show-icon :closable="false" />
  </el-card>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onUnmounted } from "vue";
import { ElNotification } from "element-plus";
import {
  Sunny,
  Eleme,
  CircleCheck,
  CircleCheckFilled,
  ArrowRight,
  OfficeBuilding,
} from "@element-plus/icons-vue";
import { createAttendance, getAttendanceDraft, viewcompanyscan } from "../api/services";
import AppButton from "../../components/AppButton.vue";
import { useUserDataStore } from "../stores/user_data";

const userDataStore = useUserDataStore();

const now = ref(new Date());
const loading = ref(false);
const companies = ref([]);
const draft = ref(null);
const draftError = ref("");

const attendForm = reactive({
  latitude: "",
  longitude: "",
  reason: "",
  company_id: null,
});

// Schedule
const scheduledDateTime = computed(() => {
  const start = draft.value?.scheduled_time?.split("-")[0];
  if (!start) return null;
  const [h, m] = start.split(":").map(Number);
  const d = new Date();
  d.setHours(h || 0, m || 0, 0, 0);
  return d;
});

const isBeforeSchedule = computed(
  () => !!scheduledDateTime.value && now.value < scheduledDateTime.value
);

const isButtonDisabled = computed(
  () => !draft.value || !!draftError.value || isBeforeSchedule.value
);

const errorMessage = computed(() => {
  if (draftError.value) return draftError.value;
  if (!isBeforeSchedule.value) return "";
  const time = scheduledDateTime.value.toLocaleTimeString("km-KH", {
    hour: "2-digit",
    minute: "2-digit",
  });
  return `មិនទាន់ដល់ម៉ោងកំណត់ទេ សូមរង់ចាំដល់ម៉ោង ${time}`;
});

// Data
async function fetchCompanies() {
  try {
    const res = await viewcompanyscan({});
    companies.value = res.data.data || [];
  } catch {
    // list stays empty
  }
}

async function fetchDraft() {
  draftError.value = "";
  draft.value = null;
  try {
    const res = await getAttendanceDraft();
    draft.value = res.data.data || null;
  } catch (e) {
    draftError.value = e.response?.data?.message || "គ្មានព័ត៌មានវត្តមាន";
  }
}

function notifyLocationError(message) {
  ElNotification({ title: "មានបញ្ហាទីតាំង", message, type: "error" });
}

function getLocation() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      notifyLocationError("Geolocation not supported");
      return reject();
    }
    navigator.geolocation.getCurrentPosition(
      ({ coords }) => {
        attendForm.latitude = String(coords.latitude);
        attendForm.longitude = String(coords.longitude);
        resolve();
      },
      () => {
        notifyLocationError("ចាប់ទីតាំងមិនបាន");
        reject();
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
    );
  });
}

async function handleCheckIn() {
  if (!attendForm.latitude || !attendForm.longitude) {
    return notifyLocationError("មិនអាចទទួលបានទីតាំង សូមបើកការអនុញ្ញាត GPS");
  }
  loading.value = true;
  try {
    await createAttendance(attendForm);
    ElNotification({ title: "ជោគជ័យ", message: "ចុះវត្តមានបានជោគជ័យ", type: "success" });
    attendForm.reason = "";
    await fetchDraft(); // refresh so the button reflects the next session
  } catch (e) {
    ElNotification({ title: "មានបញ្ហា", message: e.response?.data?.error || "", type: "error" });
  } finally {
    loading.value = false;
  }
}

// Select the user's default company once the list is loaded
watch(
  [companies, () => userDataStore.classid],
  ([list, defaultId]) => {
    if (!attendForm.company_id && list.some((c) => c.id === defaultId)) {
      attendForm.company_id = defaultId;
    }
  },
  { immediate: true }
);

let timer;
onMounted(() => {
  timer = setInterval(() => (now.value = new Date()), 30000);
  getLocation().catch(() => {});
  fetchDraft();
  fetchCompanies();
});
onUnmounted(() => clearInterval(timer));
</script>

<style scoped>
.checkin-card {
  border-radius: 6px;
}
.action-buttons :deep(.el-button + .el-button) {
  margin-left: 0;
}
.header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.logo {
  width: 80px;
  height: 80px;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #122133;
}

.card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  background: #f8fafc;
}

/* Schedule */
.schedule {
  border-color: #ffffff;
}

.icon-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  flex: 0 0 50px;
  color: #52c41a;
  background: #dbf5c5;
  border-radius: 6px;
}

.icon-box1 {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  flex: 0 0 50px;
  color: #626aef;
  background: #cbcdf1;
  border-radius: 6px;
}

.schedule-time {
  font-size: 14px;
  font-weight: 600;
  color: #122133;
}

.schedule-type {
  font-size: 12px;
  color: #7b8794;
}

/* Companies */
.company {
  background: #fff;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.company:hover {
  border-color: #409eff;
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.12);
}

.company.active {
  border: 2px solid #626aef;
  background: #f5f9ff;
}

.company-name {
  flex: 1;
  font-size: 12px;
  font-weight: 600;
  color: #303133;
}

.company-check {
  font-size: 22px;
  color: #409eff;
}

/* Buttons */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

.btn-content {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
}
</style>