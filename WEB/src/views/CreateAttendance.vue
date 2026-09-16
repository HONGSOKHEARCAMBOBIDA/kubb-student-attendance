<template>
  <el-card class="checkin-card">
    <template #header>
      <div style="display:flex;align-items:center;justify-content:center;gap:8px">
        <el-text style="color:black;font-size:18px;font-weight:bold">
          {{ draft?.type_string ?? "វត្តមាន" }}
        </el-text>
      </div>
      <el-row justify="space-between" align="middle">
        <el-icon color="#409efc" :size="25"><Calendar /></el-icon>
        <el-text style="color:black;font-size:13px">
          {{ new Date().toLocaleDateString("km-KH", { weekday: "long", year: "numeric", month: "long", day: "numeric" }) }}
        </el-text>
        <AppButton type="success" size="small" @click="getLocation()" icon="MapLocation" circle />
      </el-row>
    </template>

    <el-form :model="attendForm">
      <el-form-item v-if="draft">
        <div style="display:flex;justify-content:space-between;width:100%;color:#606266;font-size:14px">
          <el-row justify="space-between" align="middle">
            <el-icon color="#409efc" :size="25" style="padding-right:5px"><AlarmClock /></el-icon>
            <el-text tag="b" type="primary">{{ draft.type_string }}</el-text>
          </el-row>
          <el-row justify="space-between" align="middle">
            <el-icon color="brown" :size="25" style="padding-right:5px"><AlarmClock /></el-icon>
            <el-text style="color:brown">ម៉ោងកំណត់: {{ draft.scheduled_time }}</el-text>
          </el-row>
        </div>
      </el-form-item>

      <el-form-item v-if="companies.length">
        <div class="company-list">
          <div
            v-for="c in companies"
            :key="c.id"
            class="company-card"
            :class="{ active: attendForm.company_id === c.id }"
            @click="selectCompany(c.id)"
          >
            <div class="company-left">
              <el-avatar :size="40" :icon="OfficeBuilding" />
              <div class="company-info">
                <div class="company-name">{{ c.name }}</div>
              </div>
            </div>
            <div class="company-right">
              <el-icon v-if="attendForm.company_id === c.id" class="selected"><CircleCheckFilled /></el-icon>
              <el-icon v-else><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          :loading="loading || draftLoading"
          :disabled="isButtonDisabled"
          @click="handleCheckIn"
          size="large"
          style="width:100%;height:80px"
        >
          <div style="display:flex;flex-direction:column;align-items:center;gap:4px">
            <span style="font-size:22px">ចូល</span>
            <span style="font-size:13px;opacity:0.95">{{ currentTime }}</span>
          </div>
        </el-button>
      </el-form-item>

      <el-alert v-if="draftError" :title="draftError" type="warning" show-icon :closable="false" />
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
const isButtonDisabled = computed(() => !draft.value || !!draftError.value);

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
    draftError.value = e.response?.data?.message || "គ្មានព័ត៌មានវត្តមានបាន";
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

function updateTime() {
  now.value = new Date();
  currentTime.value = now.value.toLocaleTimeString("km-KH", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

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
  updateTime();
  timer = setInterval(updateTime, 1000);
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