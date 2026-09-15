<template>
  <el-card class="checkin-card">
    <template #header>
      <div
        style="
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
        "
      >
        <el-text style="color: black; font-size: 18px; font-weight: bold"
          >{{ draft?.type_string ?? "វត្តមាន" }}</el-text
        >
      </div>
      <el-row justify="space-between" align="middle">
          <el-icon color="#409efc" :size="25">
   <Calendar />
  </el-icon>
        <el-text style="color: black; font-size: 13px; ">
          {{
            new Date().toLocaleDateString("km-KH", {
              weekday: "long",
              year: "numeric",
              month: "long",
              day: "numeric",
            })
          }}
        </el-text>
        <AppButton
          type="success"
          size="small"
          @click="getLocation()"
          icon="MapLocation"
          circle
        >
        </AppButton>
      </el-row>
    </template>
    <el-form :model="attendForm">


      <!-- Draft info row -->
      <el-form-item v-if="draft">
        <div
          style="
            display: flex;
            justify-content: space-between;
            width: 100%;
            color: #606266;
            font-size: 14px;
          "
        >
        <el-row justify="space-between" align="middle">

   <el-icon color="#409efc" :size="25" style="padding-right: 5px;"><AlarmClock /></el-icon>

            <el-text tag="b" type="primary">
            
            {{ draft.type_string }}
          
          </el-text>
        </el-row>

          <el-row justify="space-between" align="middle">
             <el-icon color="brown" :size="25" style="padding-right: 5px;"><AlarmClock /></el-icon>
            <el-text style="color: brown;"
            >ម៉ោងកំណត់: {{ draft.scheduled_time }}</el-text
          >
          </el-row>
          <!-- <el-switch v-model="attendForm.is_permission" size="large">

          </el-switch> -->
        </div>
      </el-form-item>

      <!-- Late / early-leave warning -->
      <el-form-item v-if="reasonRequired">
        <el-alert
          :title="
            isLate
              ? 'អ្នកមកយឺត សូមបញ្ចូលមូលហេតុ'
              : 'អ្នកចេញមុនម៉ោង សូមបញ្ចូលមូលហេតុ'
          "
          type="warning"
          show-icon
          :closable="false"
        />
      </el-form-item>
            <el-form-item label="មូលហេតុ" :required="reasonRequired" v-show="reasonRequired">
        <el-input
          type="textarea"
          v-model="attendForm.reason"
          :placeholder="reasonPlaceholder"
          style="width: 100%"
        />
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
        <el-avatar
          :size="40"
          
          icon="OfficeBuilding"
        />

        <div class="company-info">
          <div class="company-name">
            {{ c.name }}
          </div>
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
</el-form-item>

      <el-form-item>
        <el-button
          :type="isCheckInType ? 'primary' : 'warning'"
          :loading="loading || draftLoading"
          :disabled="isButtonDisabled"
          @click="handleCheckIn"
          size="large"
          style="width: 100%; height: 80px"
        >
          <div
            style="
              display: flex;
              flex-direction: column;
              align-items: center;
              gap: 4px;
            "
          >
            <div
              style="
                padding-top: 1px;
                padding-bottom: 10px;
                padding-left: 16px;
                padding-right: 16px;
              "
            >
              <span style="font-size: 22px">
                 {{ isCheckInType ? 'ចូល' : 'ចេញ' }}
              </span>
            </div>
            <span style="font-size: 13px; opacity: 0.95">{{
              currentTime
            }}</span>
          </div>
        </el-button>
      </el-form-item>

      <!-- Error state -->
      <el-alert
        v-if="draftError"
        :title="draftError"
        type="warning"
        show-icon
        :closable="false"
      />
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from "vue";
import { ElMessage } from "element-plus";
import { createAttendance, getAttendanceDraft } from "../api/services";
import { ElNotification } from "element-plus";
import AppButton from "../../components/AppButton.vue";
import { viewcompanyscan } from "../api/services";
import { useUserDataStore } from "../stores/user_data";
import { Calendar } from "@element-plus/icons-vue";
import {
  CircleCheckFilled,
  ArrowRight,
  OfficeBuilding
} from '@element-plus/icons-vue'
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
const userDataStore = useUserDataStore();
function selectCompany(id) {
  attendForm.company_id = id;
}

const defaultcompanyid = computed(() => userDataStore.companyid || null);

async function fetchCompanies() {
  loading.value = true;
  try {
    const res = await viewcompanyscan({});
    companies.value = res.data.data || [];
  } catch (e) {
  } finally {
    loading.value = false;
  }
}
const draft = ref(null);
const draftLoading = ref(false);
const draftError = ref("");

const isCheckInType = computed(
  () => draft.value?.type === 1 || draft.value?.type === 3,
);
const isCheckOutType = computed(
  () => draft.value?.type === 2 || draft.value?.type === 4,
);

const scheduledDateTime = computed(() => {
  const t = draft.value?.scheduled_time;
  if (!t) return null;
  const [h, m, s] = t.split(":").map(Number);
  const d = new Date();
  d.setHours(h || 0, m || 0, s || 0, 0);
  return d;
});

const isLate = computed(
  () =>
    isCheckInType.value &&
    scheduledDateTime.value &&
    now.value > scheduledDateTime.value,
);

const isEarlyLeave = computed(
  () =>
    isCheckOutType.value &&
    scheduledDateTime.value &&
    now.value < scheduledDateTime.value,
);

const reasonRequired = computed(() => isLate.value || isEarlyLeave.value);

const reasonPlaceholder = computed(() => {
  if (isLate.value) return "សូមបញ្ចូលមូលហេតុមកយឺត";
  if (isEarlyLeave.value) return "សូមបញ្ចូលមូលហេតុចេញមុនម៉ោង";
  return "បញ្ចូលមូលហេតុ";
});

const reasonMissing = computed(
  () => reasonRequired.value && !attendForm.reason.trim(),
);

const isButtonDisabled = computed(() => !draft.value || reasonMissing.value);

async function fetchDraft() {
  draftLoading.value = true;
  draftError.value = "";
  draft.value = null;
  try {
    const res = await getAttendanceDraft();
    draft.value = res.data.data || [];
  } catch (e) {
    draftError.value = e.response?.data?.message || "គ្មានព័ត៌មានវត្តមានបាន";
  } finally {
    draftLoading.value = false;
  }
}

// function getLocation() {
//   if (!navigator.geolocation) return ElMessage.warning('Geolocation not supported')
//   navigator.geolocation.getCurrentPosition(
//     (pos) => {
//       attendForm.latitude = String(pos.coords.latitude)
//       attendForm.longitude = String(pos.coords.longitude)
//     },
//     () => ElMessage.error('Failed to get location')
//   )
// }

function getLocation() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      ElNotification({
        title: "មានបញ្ហាទីតាំង",
        message: "Geolocation not supported",
        type: "error",
      });
      return reject();
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        attendForm.latitude = String(pos.coords.latitude);
        attendForm.longitude = String(pos.coords.longitude);
        resolve();
      },
      () => {
        ElNotification({
          title: "មានបញ្ហាទីតាំង",
          message: "ចាប់ទីតាំងមិនបាន",
          type: "error",
        });
        reject();
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 },
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
  if (reasonMissing.value) {
    return ElMessage.warning("សូមបញ្ចូលមូលហេតុជាមុនសិន");
  }
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
    ElNotification({
      title: "ជោគជ័យ",
      message: "ចុះវត្តមានបានជោគជ័យ",
      type: "success",
    });
    ((attendForm.reason = ""));
    // Refresh draft so the button updates to the next session
    await fetchDraft();
  } catch (e) {
    ElNotification({
      title: "មានបញ្ហា",
      message: e.response?.data?.error || "",
      type: "error",
    });
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
  { immediate: true },
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
  transition: all .25s ease;
}

.company-card:hover {
  border-color: #409EFF;
  box-shadow: 0 6px 16px rgba(64,158,255,.12);
}

.company-card.active {
  border: 2px solid #409EFF;
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

.company-sub {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.company-right {
  font-size: 22px;
  color: #409EFF;
}

.selected {
  color: #409EFF;
}
</style>
