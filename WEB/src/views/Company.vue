<template>
  <div>
    <div class="page-header">
      <AppButton
        v-if="canAddCompany"
        type="primary"
        @click="openCreate"
        :block="false"
      >
        បន្ថែមក្រុមហ៑ុន
      </AppButton>
    </div>

    <el-card class="table-card">
      <AppTable
        :data="companies"
        :loading="loading"
        show-index
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchCompanies"
        :columns="[
          { prop: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'user_count', label: 'បុគ្គលិកសរុប', minWidth: 110 },
          { prop: 'latitude', label: 'Latitude', width: 120 },
          { prop: 'longitude', label: 'Longitude', width: 120 },
          { prop: 'radius', label: 'ចម្ងាយអាចស្កែន (m)', width: 150 },
          { prop: 'currency', label: 'រូបិយប័ណ្ណ', width: 100 },
          { prop: 'late_penalty', label: 'ពិន័យយឺត', width: 120 },
          { prop: 'left_early_penalty', label: 'ពិន័យចេញមុនម៉ោង', width: 150 },
          { label: 'អាចស្កែនក្រៅតំបន់', slot: 'outsize', width: 150 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]"
      >
        <template #outsize="{ row }">
          <el-tag :type="row.can_scan_outsize ? 'success' : 'danger'">
            {{ row.can_scan_outsize ? "បាន" : "មិនបាន" }}
          </el-tag>
        </template>

        <template #status="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
            {{ row.is_active ? "Active" : "Inactive" }}
          </el-tag>
        </template>

        <template #actions="{ row }">

          <AppButton
             v-if="canEditCompany"
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="openEdit(row)"         
          >

          </AppButton>
          <AppButton
            v-if="canEditCompany"
            size="small"
            icon="Promotion"
            type="primary"
            circle
            @click="openEditTelegram(row)"          
          >

          </AppButton>
        </template>
      </AppTable>
    </el-card>
    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែក្រុមហ៊ុន' : 'បន្ថែមក្រុមហ៊ុន'"
      width="640px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <AppTabs
          v-model="activeTab"
          :tabs="[
            { name: 'general', label: 'ព័ត៌មានទូទៅ' },
            { name: 'penalty', label: 'ប្រាក់ពិន័យ' },
            {name: 'telegram',label: 'Telegram'}
          ]"
          tab-position="top"
          stretch="true"
        >
          <template #general>
            <el-form-item label="ឈ្មោះក្រុមហ៊ុន" prop="name">
              <el-input
                v-model.trim="form.name"
                size="large"
                placeholder="បញ្ចូលឈ្មោះក្រុមហ៊ុន"
              />
            </el-form-item>

            <el-form-item label="រូបិយប័ណ្ណសម្រាប់បើកប្រាក់ខែ" prop="currency">
              <el-radio-group v-model="form.currency" size="large">
                <el-radio-button value="USD">USD</el-radio-button>
                <el-radio-button value="KHR">KHR</el-radio-button>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="Map Link" prop="map_link">
              <el-input
                v-model.trim="form.map_link"
                size="large"
                placeholder="https://maps.google.com/..."
              />
            </el-form-item>


            <div class="form-row">
              <el-form-item label="Latitude" prop="latitude">
                <el-input
                  v-model.trim="form.latitude"
                  placeholder="e.g. 11.5564"
                  size="large"
                />
              </el-form-item>
              <el-form-item label="Longitude" prop="longitude">
                <el-input
                  v-model.trim="form.longitude"
                  placeholder="e.g. 104.9282"
                  size="large"
                />
              </el-form-item>
            </div>

            <div class="">
              <el-form-item label="ចម្ងាយអាចស្កែនបាន (ម៉េត្រ)" prop="radius">
                <el-input
                  v-model.trim="form.radius"
                  placeholder="e.g. 100"
                  size="large"
                />
              </el-form-item>
              <el-form-item label="អាចស្កែនក្រៅតំបន់" prop="can_scan_outsize">
                <el-radio-group v-model="form.can_scan_outsize" size="large">
                  <el-radio-button :value="1">អាចស្កែនបាន</el-radio-button>
                  <el-radio-button :value="0">មិនអាចស្កែនបាន</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="Color">
                <el-color-picker-panel v-model="form.color" />
              </el-form-item>
              <el-form-item label="ចំនួនថ្ងៃធ្វេីការ">
                <el-input
                  type="number"
                  v-model.number="form.total_work_day"
                  placeholder="30"
                  size="large"
                />
              </el-form-item>
            </div>
          </template>
          <template #penalty>
            <div class="form-row">
              <el-form-item label="ពិន័យពេលមកយឺត" prop="late_penalty">
                <el-input
                  v-model="form.late_penalty"
                  placeholder="ប្រាក់ពិន័យពេលមកយឺត"
                  size="large"
                >
                  <template #append>{{ form.currency }}</template>
                </el-input>
              </el-form-item>
              <el-form-item label="ពិន័យចេញមុនម៉ោង" prop="left_early_penalty">
                <el-input
                  v-model="form.left_early_penalty"
                  placeholder="ប្រាក់ពិន័យពេលចេញមុនម៉ោង"
                  size="large"
                >
                  <template #append>{{ form.currency }}</template>
                </el-input>
              </el-form-item>
            </div>
          </template>
          <template #telegram>
            <el-form-item label="Bot Token" prop="bot_token">
              <el-input
                v-model="form.bot_token"
                placeholder="Telegram bot token"
                size="large"
              />
            </el-form-item>
            <el-form-item label="Group Link" prop="group_link">
              <el-input
                v-model="form.group_link"
                placeholder="Telegram group chat ID"
                size="large"
              />
            </el-form-item>
          </template>
        </AppTabs>
      </el-form>
      <template #footer>
        <AppButton @click="dialogVisible = false" size="large" :block="false" type="warning">
          បោះបង់
        </AppButton>
        <AppButton
          @click="handleSave"
          type="primary"
          :loading="saving"
          size="large"
          :block="false"
        >
          {{ isEdit ? "កែប្រែ" : "បង្កើត" }}
        </AppButton>
      </template>
    </AppDialog>

    <AppDialog v-model="dialogTelegramVisible" title="កែប្រែ" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-divider />
        <p class="section-label">Telegram</p>
        <div class="form-row">
          <el-form-item label="Bot Token" prop="bot_token">
            <el-input
              v-model="form.bot_token"
              placeholder="Telegram bot token"
              size="large"
            />
          </el-form-item>
          <el-form-item label="Group Link" prop="group_link">
            <el-input
              v-model="form.group_link"
              placeholder="Telegram group chat ID"
              size="large"
            />
          </el-form-item>
        </div>
      </el-form>


      <template #footer>
        <AppButton @click="dialogTelegramVisible = false" size="large" :block="false">
          ថតក្រោយ
        </AppButton>
        <AppButton
          type="primary"
          :loading="saving"
          @click="handleUpdateTelegram"
          size="large"
          :block="false"
        >
          កែប្រែ
        </AppButton>
      </template>
    </AppDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage } from "element-plus";
import { useCompanyStore } from '../stores/company'
import { useUserDataStore } from '../stores/user_data'
import {
  getCompany,
  createCompany,
  updateCompany,
  updateTelegram,
} from "../api/services";
import { useAuthStore } from "../stores/auth";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppTabs from "../../components/AppTabs.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
const notify = useNotification()
const userDataStore = useUserDataStore()
const companyStore = useCompanyStore()
const companies = ref([]);
const loading = ref(false);
const useloading = useLoading();
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const dialogVisible = ref(false);
const isEdit = ref(false);
const isEditTelegram = ref(false);
const dialogTelegramVisible = ref(false);
const editId = ref(null);
const formRef = ref();
const auth = useAuthStore();
const activeTab = ref("general");

const form = reactive({
  map_link: "",
  name: "",
  latitude: "",
  longitude: "",
  radius: "",
  currency: "",
  late_penalty: "",
  left_early_penalty: "",
  bot_token: "",
  group_link: "",
  can_scan_outsize: null,
  color: "",
  total_work_day:null
});

const canAddCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.company"),
);

const canEditCompany = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.company"),
);

const rules = {
  name: [{ required: true, message: "Company name is required" }],
  radius: [{ required: true, message: "Radius is required" }],
  group_link: [{ required: true, message: "Group Link is required" }],
  bot_token: [{ required: true, message: "BotToken is required" }],
  currency: [{ required: true, message: "Currency is required" }],
  late_penalty: [{ required: true, message: "Late Penalty required" }],
  left_early_penalty: [
    { required: true, message: "Left Early Penalty required" },
  ],
  can_scan_outsize: [
    { required: true, message: "Employee Can scan outsize or not" },
  ],
};

async function fetchCompanies() {
  loading.value = true;
  try {
    const res = await getCompany({
      page: page.value,
      page_size: pageSize.value,
    });
    companies.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  isEdit.value = false;
  Object.keys(form).forEach((k) => (form[k] = ""));
  form.can_scan_outsize = null;
  form.late_penalty = 0;
  form.left_early_penalty = 0;
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    map_link: "",
    name: row.name || "",
    latitude: row.latitude || "",
    longitude: row.longitude || "",
    radius: row.radius || "",
    currency: row.currency || "",
    bot_token: row.bot_token || "",
    group_link: row.group_link || "",
    late_penalty: row.late_penalty || "",
    left_early_penalty: row.left_early_penalty || "",
    can_scan_outsize: row.can_scan_outsize ?? null,
    color: row.color || "",
    total_work_day: row.total_work_day || "",
  });
  dialogVisible.value = true;
}

function openEditTelegram(row) {
  isEditTelegram.value = true;
  editId.value = row.id;
  Object.assign(form, {
    bot_token: row.bot_token || "",
    group_link: row.group_link || "",
  });
  dialogTelegramVisible.value = true;
}


async function handleSave() {
  await formRef.value.validate();
  useloading.show({
    text: "កំពុងដំណេីរការ..."
  })
  try {
    if (isEdit.value) {
      const payload = {};
      Object.entries(form).forEach(([k, v]) => {
        if (v !== "") payload[k] = v;
      });
      if (form.can_scan_outsize !== null && form.can_scan_outsize !== "") {
        payload.can_scan_outsize = form.can_scan_outsize;
      }
      await updateCompany(editId.value, payload);
            if (form.color) {
        companyStore.setColor(form.color);
      }
      notify.success("កែប្រែបានជោគជ័យ");
    } else {
      await createCompany(form);
      notify.success("បង្កេីតក្រុមហ៑ុនបានជោគជ័យ");
    }
    dialogVisible.value = false;
    fetchCompanies();
  } catch (e) {
    notify.error(e.response?.data?.error || "")
  } finally {
    useloading.hide()
  }
}

async function handleUpdateTelegram() {
  await formRef.value.validate();
  useloading.show({
    text: "កំពុងដំណេីរការ..."
  })
  try {
    if (isEditTelegram.value) {
      const payload = {};
      Object.entries(form).forEach(([k, v]) => {
        if (v !== "") payload[k] = v;
      });
      await updateTelegram(editId.value, payload);
      notify.success("កែប្រែបានជោគជ័យ");
    }
    dialogTelegramVisible.value = false;
    fetchCompanies();
  } catch (e) {
    notify.error(e.response?.data?.error || "")
  } finally {
    useloading.hide()
  }
}

onMounted(fetchCompanies);
</script>


<style scoped>
.page-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}

.section-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
  margin-bottom: 12px;
}

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
