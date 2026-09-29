<template>
  <el-card class="checkin-card">
    <template #header>
      <div style="
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 4px;
  ">
        <div class="login-logo">
          <el-image src="/logo.png" alt="University Logo" fit="contain" />
        </div>

        <el-tag
  size="large"
  :type="draft?.subject_name ? 'success' : 'warning'"
>
  {{ draft?.subject_name ?? "វត្តមាន" }}
</el-tag>
 
      </div>
      <!-- <el-row justify="space-between" align="middle">
        <el-icon color="#409efc" :size="25">
          <Calendar />
        </el-icon>
        <el-text style="color:black;font-size:13px">
          {{ new Date().toLocaleDateString("km-KH", { weekday: "long", year: "numeric", month: "long", day: "numeric" })
          }}
        </el-text>
        <AppButton type="success" size="small" @click="getLocation()" icon="MapLocation" circle />
      </el-row> -->
    </template>

    <el-form :model="attendForm">
      <el-form-item v-if="draft">
<div class="schedule-section">
  <div class="section-title">
    <span>កាលវិភាគ</span>
  </div>

  <div class="schedule-card">
    <div class="schedule-icon">
      <el-icon :size="22">
        <Sunny />
      </el-icon>
    </div>

    <div class="schedule-info">
      <div class="schedule-time">
        {{ draft.scheduled_time }}
      </div>

      <div class="schedule-type">
        {{ draft.type_string }}
      </div>
    </div>
  </div>
</div>

      </el-form-item>

      <!-- <el-form-item v-if="companies.length">
        <div class="company-list">
          <div v-for="c in companies" :key="c.id" class="company-card"
            :class="{ active: attendForm.company_id === c.id }" @click="selectCompany(c.id)">
            <div class="company-left">
              <el-avatar :size="40" :icon="OfficeBuilding" />
              <div class="company-info">
                <div class="company-name">{{ c.name }}</div>
              </div>
            </div>
            <div class="company-right">
              <el-icon v-if="attendForm.company_id === c.id" class="selected">
                <CircleCheckFilled />
              </el-icon>
              <el-icon v-else>
                <ArrowRight />
              </el-icon>
            </div>
          </div>
        </div>
      </el-form-item> -->

        <div class="action-buttons">
  <AppButton
    type="default"
    native-type="button"
    plain
    :loading="loading"
    :disabled="isBeforeSchedule"
    @click="$router.push('/leaverequest')"
    block
  >
    <span class="login-button-content">
      <el-icon>
        <Eleme />
      </el-icon>
      <span>សុំច្បាប់</span>
    </span>
  </AppButton>

  <AppButton
    color="#626aef"
    native-type="button"
    plain
    :loading="loading"
    :disabled="isBeforeSchedule"
    @click="handleCheckIn"
    block
  >
    <span class="login-button-content">
      <el-icon>
        <CircleCheck />
      </el-icon>
      <span>ចុះវត្តមានចូល</span>
    </span>
  </AppButton>
</div>

      <el-alert v-if="draftError" :title="draftError" type="error" show-icon :closable="false" />
      <el-alert v-else-if="isBeforeSchedule" :title="scheduleWaitMessage" type="error" show-icon :closable="false" />
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from "vue";
import { ElMessage, ElNotification } from "element-plus";
import { Calendar, CircleCheckFilled, ArrowRight, OfficeBuilding } from "@element-plus/icons-vue";
import { createAttendance, getAttendanceDraft, viewcompanyscan } from "../api/services";
import AppButton from "../../components/AppButton.vue";
import { useUserDataStore } from "../stores/user_data";

const userDataStore = useUserDataStore();

const now = ref(new Date());
const currentTime = ref("");
const loading = ref(false);

const attendForm = reactive({
  latitude: "",
  longitude: "",
  reason: "",
  company_id: null,
});

const companies = ref([]);
const draft = ref(null);
const draftLoading = ref(false);
const draftError = ref("");

const defaultcompanyid = computed(() => userDataStore.classid || null);

const scheduledDateTime = computed(() => {
  const t = draft.value?.scheduled_time;

  if (!t) return null;

  const [start] = t.split("-");
  const [h, m] = start.split(":").map(Number);

  const d = new Date();
  d.setHours(h || 0, m || 0, 0, 0);

  return d;
});

const isBeforeSchedule = computed(() => {
  return !!scheduledDateTime.value && now.value < scheduledDateTime.value;
});

const scheduleWaitMessage = computed(() => {
  if (!isBeforeSchedule.value || !scheduledDateTime.value) return "";
  const timeStr = scheduledDateTime.value.toLocaleTimeString("km-KH", {
    hour: "2-digit",
    minute: "2-digit",
  });
  return `មិនទាន់ដល់ម៉ោងកំណត់ទេ សូមរង់ចាំដល់ម៉ោង ${timeStr}`;
});

const isButtonDisabled = computed(
  () => !draft.value || !!draftError.value || isBeforeSchedule.value
);

function selectCompany(id) {
  attendForm.company_id = id;
}

async function fetchCompanies() {
  loading.value = true;
  try {
    const res = await viewcompanyscan({});
    companies.value = res.data.data || [];
  } catch (e) {
    // ignore, list stays empty
  } finally {
    loading.value = false;
  }
}

async function fetchDraft() {
  draftLoading.value = true;
  draftError.value = "";
  draft.value = null;
  try {
    const res = await getAttendanceDraft();
    draft.value = res.data.data || null;
  } catch (e) {
    draftError.value = e.response?.data?.message || "គ្មានព័ត៌មានវត្តមាន";
  } finally {
    draftLoading.value = false;
  }
}

function getLocation() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      ElNotification({ title: "មានបញ្ហាទីតាំង", message: "Geolocation not supported", type: "error" });
      return reject();
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        attendForm.latitude = String(pos.coords.latitude);
        attendForm.longitude = String(pos.coords.longitude);
        resolve();
      },
      () => {
        ElNotification({ title: "មានបញ្ហាទីតាំង", message: "ចាប់ទីតាំងមិនបាន", type: "error" });
        reject();
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
    );
  });
}

// function updateTime() {
//   now.value = new Date();
//   currentTime.value = now.value.toLocaleTimeString("km-KH", {
//     hour: "2-digit",
//     minute: "2-digit",
//     second: "2-digit",
//   });
// }

async function handleCheckIn() {
  if (!attendForm.latitude || !attendForm.longitude) {
    return ElNotification({
      title: "មានបញ្ហាទីតាំង",
      message: "មិនអាចទទួលបានទីតាំង សូមបើកការអនុញ្ញាត GPS",
      type: "error",
    });
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

let timer;
onMounted(() => {
  // updateTime();
  // timer = setInterval(updateTime, 1000);
  getLocation();
  fetchDraft();
  fetchCompanies();
});
onUnmounted(() => clearInterval(timer));

watch(
  () => [companies.value, defaultcompanyid.value],
  ([list, defaultId]) => {
    if (!list.length) return;
    const exists = list.some((c) => c.id === defaultId);
    if (exists && !attendForm.company_id) {
      attendForm.company_id = defaultId;
    }
  },
  { immediate: true }
);
</script>

<style scoped>

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.login-button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
}

.login-button-content {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
.schedule-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 7px;

  font-size: 18px;
  font-weight: 600;
  color: #122133;
}

.section-title .el-icon {
  font-size: 17px;
  color: #122133;
}

.schedule-card {
  display: flex;
  align-items: center;
  gap: 12px;

  width: 100%;
  box-sizing: border-box;

  padding: 12px 14px;

  background: #f8fafc;
  border: 1px solid #0862e9;
  border-radius: 10px;
}

.schedule-icon {
  width: 50px;
  height: 50px;
  flex: 0 0 50px;

  display: flex;
  align-items: center;
  justify-content: center;

  color: #52c41a;
  background: #f1fbe9;
  border-radius: 6px;
}

.schedule-info {
  min-width: 0;
  flex-direction: column;

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

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 5px;
}

.login-logo .el-image {
  width: 80px;
  height: 80px;
}

.checkin-card {
  border-radius: 6px;
}

.company-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.company-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 18px;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  background: #fff;
  cursor: pointer;
  transition: all 0.25s ease;
}

.company-card:hover {
  border-color: #409eff;
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.12);
}

.company-card.active {
  border: 2px solid #409eff;
  background: #f5f9ff;
}

.company-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.company-info {
  display: flex;
  flex-direction: column;
}

.company-name {
  font-size: 12px;
  font-weight: 600;
  color: #303133;
}

.company-right {
  font-size: 22px;
  color: #409eff;
}

.selected {
  color: #409eff;
}
</style>