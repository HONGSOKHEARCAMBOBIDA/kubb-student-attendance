<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 5 },
        { slot: 'payroll_type', span: 6 },
        {slot: 'company',span: 5},
        { slot: 'payroll_date', span: 5 },
        {slot: 'delete_payroll',span:3}
      ]"
      :action-span="5"
    >
      <template #name>
        <el-input
          size="large"
          v-model="filters.name"
          placeholder="ស្វែងរក"
          prefix-icon="Search"
          clearable
        />
      </template>
      <template #payroll_type>
        <el-select
          size="large"
          v-model="filters.payroll_type"
          @change="fetchPayroll"
          clearable
        >
          <el-option label="ទាំងអស់" value="" />
          <el-option label="បេីកពេញ១ខែ" :value="1" />
          <el-option label="បេីកកន្លះខែ" :value="2" />
        </el-select>
      </template>
      <template #company>
          <el-select v-model="filters.company_id" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%" size="large"
          @change="fetchPayroll">
          <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
        </el-select>      
      </template>
      <template #payroll_date>
        <el-date-picker
          size="large"
          v-model="filters.payroll_date"
          type="month"
          placeholder="ជ្រើសរើសខែ និង ឆ្នាំ"
          format="MM/YYYY"
          value-format="YYYY-MM"
          @change="fetchPayroll"
          style="width: 100%;"
        />
      </template>
      <template #delete_payroll>
        <AppButton
        v-if="candeletepayroll"
        type="danger"
        @click="confirmDeletePayroll"
        :disabled="!selectedRows.length"
        :loading="loading"
        :block="true"
        >
        លុប
        </AppButton>
      </template>
    </AppFilterBar>
    <el-card class="top-card">
      <AppTable
      :data="payrolls"
      :loading="loading"
      show-index
      selectable
      :total="totalCount"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchPayroll"
      :columns="[
        {prop: 'user_name',label: 'ឈ្មោះ',minWidth:100},
        {label: 'ភេទ',slot: 'gender',width: 80},
        {prop: 'role_name',label: 'តួនាទី',width: 110},
        {prop: 'company_name',label: 'ក្រុមហ៊ុន',width:160},
        {label:'ប្រភេទប្រាក់ខែ',slot: 'payroll_type',width:130},
        {label: 'ប្រាក់ខែគោល',slot: 'salary',width: 130},
        {label: 'ពាក់កណ្ដាល',slot: 'half_salary',width: 135},
        {prop: 'total_work_day',label: 'ចំនួនថ្ងៃធ្វើការ',width: 120},
        {label: 'កាត់លុយ',slot: 'total_deduction',width: 130},
        {label: 'ប្រាក់ខែទទួលបាន',slot: 'net_salary',width: 160},
        {label: 'ផ្សេងៗ',slot: 'other',width: 110},
        {label: 'ស្ថានភាព',slot: 'status',width: 100},
        {prop: 'payroll_date',label:'កាលបរិច្ឆេទ',width: 120},
        {label: 'សម្គាល់',slot: 'note',width: 120},
      ]"
      @selection-change="(rows) => (selectedRows = rows)"
      >
      <template #gender="{ row }">
            <el-tag
              :type="row.user_gender === 1 ? 'primary' : 'danger'"
              size="small"
            >
              {{ row.user_gender === 1 ? "ប្រុស" : "ស្រី" }}
            </el-tag>
          </template>
      <template #payroll_type="{ row }">
            <el-tag
              :type="row.payroll_type === 1 ? 'primary' : 'warning'"
              size="small"
            >
              {{ row.payroll_type === 1 ? "ពេញ១ខែ" : "កន្លះខែ" }}
            </el-tag>
          </template>
      <template #salary="{ row }"
            >{{ row.basic_salary }} {{ row.currency }}
          </template
          >
        <template #half_salary="{ row }"
            >{{ row.half_salary }} {{ row.currency }}</template
          >
         <template #total_deduction="{ row }"
            >{{ row.total_deduction || "0.00" }} {{ row.currency }}</template
          >
          <template #net_salary="{ row }">
            <strong>{{ row.net_salary }} {{ row.currency }}</strong>
          </template>
          <template #other="{ row }">
            <span v-if="row.other">{{ row.other }}</span>
            <span v-else>—</span>
          </template>
          <template #status="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">
              {{ row.status === 0 ? "បានបេីក" : "មិនទាន់" }}
            </el-tag>
          </template>
          <template #note="{ row }">
            <span v-if="row.note">{{ row.note }}</span>
            <span v-else>—</span>
          </template>
      </AppTable>
    </el-card>
    <AppDialog v-model="deleteDialog" title="ពិនិត្យប្រាក់ខែម្ដងទៀត" width="720px">
         <el-alert
        title="សូមពិនិត្យព័ត៌មានប្រាក់បៀវត្សរ៍ឱ្យបានត្រឹមត្រូវ មុនពេលលុប!"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 16px"
      />   
      <el-table :data="selectedRows" size="large" max-height="400">
        <el-table-column prop="user_name" label="ឈ្មោះបុគ្គលិក" />
        <el-table-column prop="net_salary" label="ប្រាក់ខែទទួលបាន" width="150">
          <template #default="{ row }"
            >{{ row.net_salary }} {{ row.currency }}</template
          >
        </el-table-column>
      </el-table>
          <template #footer>
        <el-button @click="deleteDialog = false">Cancel</el-button>
        <el-button type="danger" :loading="loading" @click="deletepayroll"
          >Delete</el-button
        >
      </template>
    </AppDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed } from "vue";
import { ElMessage } from "element-plus";
import { getPayroll,getCompany,deletePayroll } from "../api/services";
import { debounce } from "lodash-es";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppTable from "../../components/AppTable.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import { ElNotification } from 'element-plus'
import { useAuthStore } from "../stores/auth.js";
import { useUserDataStore } from '../stores/user_data'
const userDataStore = useUserDataStore()
const candeletepayroll = computed(()=>userDataStore.permissions?.some((p)=> p.name === "edit.payroll"));
const payrolls = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const totalCount = ref(0);
const companys = ref([])
const selectcompany = ref()
const selectedRows = ref([]);
const deleteDialog = ref(false)
const filters = reactive({
  name: "",
  payroll_date: new Date().toISOString().slice(0, 7),
  payroll_type: "",
  company_id: null
});

function confirmDeletePayroll(){
  if(!selectedRows.value.length)
  return   ElNotification.error({
    title: 'មានបញ្ហា',
    message: 'ជ្រេីសរេីសតេីមួយណាដែលអ្នកចង់លុប',
    offset: 100,
  })
  deleteDialog.value = true;
}

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

async function deletepayroll(){
  loading.value = true;
  try {
const payroll_ids = selectedRows.value.map(r => r.id)
    await deletePayroll({ payroll_ids })
    ElNotification.success({
      title: 'ជោគជ័យ',
      message: 'លុបប្រាក់ខែបានជោគជ័យ',
      offset: 100
    });
    deleteDialog.value = false;
    selectedRows.value = [];
    await fetchPayroll();
  } catch (e) {
    ElNotification.error({
      title: 'មានបញ្ហា',
      message: e.response?.data?.error || "",
      offset: 100,
    });
  } finally {
    loading.value = false;
  }
}


async function fetchPayroll() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.payroll_date) params.payroll_date = filters.payroll_date;
    if (filters.payroll_type) params.payroll_type = filters.payroll_type;
    if(filters.company_id) params.company_id = filters.company_id
    const res = await getPayroll(params);
    payrolls.value = res.data.data || [];
    totalCount.value = res.data.pagination?.totalCount || 0;
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to load payroll");
  } finally {
    loading.value = false;
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1;
  fetchPayroll();
}, 500);
watch(() => filters.name, debouncedFetch);
onMounted(()=>{
  fetchPayroll()
  fetchCompany()
});
</script>

<style scoped>
.top-card {
  border-radius: 6px;
}
.card-title {
  font-weight: 600;
  font-size: 15px;
}

</style>
