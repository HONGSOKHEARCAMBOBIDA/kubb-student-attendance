<template>
  <div>
    <AppFilterBar
      :fields="[
        { slot: 'payrollType', span: 6 },
        {slot: 'company',span: 6},
        { slot: 'getdraft', span: 6 },
        { slot: 'createpayroll', span: 6 },
      ]"
      :action-span="4"
    >
      <template #payrollType>
        <el-select v-model="payrollType" @change="fetchDraft" size="large">
          <el-option label="បេីកពេញ១ខែ" :value="1" />
          <el-option label="បេីកកន្លះខែ" :value="2" />
        </el-select>
      </template >

      <template #company>
        <el-select v-model="selectcompany" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%" size="large"
          @change="fetchDraft">
          <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
        </el-select>
      </template>
      <template #getdraft>
        <AppButton
          type="primary"
          @click="fetchDraft"
          :loading="draftLoading"
          :block="true"
        >
          ទាញទិន្ន័យម្ដងទៀត
        </AppButton>
      </template>
      <template #createpayroll>
        <AppButton
          type="success"
          @click="confirmPayroll"
          :disabled="!selectedRows.length"
          :loading="paying"
          :block="true"
        >
          បេីកប្រាក់ខែ ({{ selectedRows.length }}) នាក់
        </AppButton>
      </template>
    </AppFilterBar>
    <el-card class="top-card">
      <AppTable
        :data="draft"
        :loading="draftLoading"
        show-index
        selectable
        :total="draft.length"
        :columns="[
          { prop: 'user_name', label: 'ឈ្មោះបុគ្គលិក', min: 50 },
          { prop: 'role_name', label: 'តួនាទី', min: 110 },
          { label: 'ប្រាក់ខែគោល', slot: 'basic_salary', min: 120 },
          { label: 'ប្រាក់ខែពាក់កណ្ដាល', slot: 'half_salary', min: 155 },
          { label: 'ចំនួនថ្ងៃធ្វើការ', prop: 'total_work_day', min: 100 },
          { label: 'ចំនួនយឺត', slot: 'total_late', min: 105 },
          { label: 'ចេញមុនម៉ោង', slot: 'total_left_early', min: 110 },
          {slot:'leave_deduction',label:'កាត់លុយច្បាប់',min:110},
          { label: 'កាត់លុយសរុប', slot: 'total_deduction', min: 80 },
          { label: 'ប្រាក់ខែទទួលបាន', slot: 'net_salary', min: 180 },
        ]"
        @selection-change="(rows) => (selectedRows = rows)"
      >
        <template #basic_salary="{ row }">
          {{ row.basic_salary }} {{ row.currency }}
        </template>
        <template #half_salary="{ row }">
          {{ row.half_salary }} {{ row.currency }}
        </template>
        <template #total_late="{row}">
          {{ row.total_late }}
          <el-tag type="danger">
           កាត់លុយ {{ row.total_penalty_late }} {{ row.currency }}
          </el-tag> 
        </template>
        <template #leave_deduction="{row}">
          <el-text tag="mark">{{ row.leave_deduction }} {{ row.currency }}</el-text>
        </template>
        <template #total_left_early="{row}">
          {{ row.total_left_early }}
                    <el-tag type="danger">
           កាត់លុយ {{ row.total_left_early_penalty }} {{ row.currency }}
          </el-tag> 
        </template>
        <template #total_deduction="{ row }">
          <el-input v-model="row.total_deduction" :row="1">
            <!-- <template #append>{{ row.currency }}</template> -->
          </el-input>
        </template>
        <template #net_salary="{ row }">
          <el-input v-model="row.net_salary">
            <!-- <template #append>{{ row.currency }}</template> -->
          </el-input>
        </template>
      </AppTable>
    </el-card>
    <AppDialog v-model="payDialog" title="ពិនិត្យប្រាក់ខែម្ដងទៀត" width="720px">
      <el-alert
        title="សូមពិនិត្យព័ត៌មានប្រាក់បៀវត្សរ៍ឱ្យបានត្រឹមត្រូវ មុនពេលបញ្ជាក់។ សកម្មភាពនេះមិនអាចត្រឡប់ក្រោយវិញបានទេ"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 16px"
      />
      <AppFilterBar
        :fields="[
          { slot: 'payroll_date', span: 10 },
          { slot: 'payroll_type', span: 6 },
          { slot: 'note', span: 5 },
        ]"
        :action-span="3"
      >
        <template #payroll_date>
          <el-date-picker
            v-model="payForm.payroll_date"
            type="date"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </template>
        <template #payroll_type>
          <el-select v-model="payForm.payroll_type" style="width: 100%">
            <el-option label="បេីកពេញ១ខែ" :value="1" />
            <el-option label="បេីកកន្លះខែ" :value="2" />
          </el-select>
        </template>
        <template #note>
          <el-input v-model="payForm.note" placeholder="Optional note" />
        </template>
      </AppFilterBar>
      <el-table :data="selectedRows" size="large" max-height="400">
        <el-table-column prop="user_name" label="ឈ្មោះបុគ្គលិក" />
        <el-table-column prop="net_salary" label="ប្រាក់ខែទទួលបាន" width="150">
          <template #default="{ row }"
            >{{ row.net_salary }} {{ row.currency }}</template
          >
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="payDialog = false">Cancel</el-button>
        <el-button type="primary" :loading="paying" @click="submitPayroll"
          >Confirm & Pay</el-button
        >
      </template>
    </AppDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getPayrollDraft, createPayroll,getCompany } from "../api/services";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
const draft = ref([]);
const draftLoading = ref(false);
const paying = ref(false);
const payrollType = ref(1);
const selectcompany = ref()
const selectedRows = ref([]);
const payDialog = ref(false);
const companys = ref([]);
const loading = ref(false);
const payForm = reactive({
  payroll_date: new Date().toISOString().split("T")[0],
  payroll_type: 1,
  note: "",
});

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


async function fetchDraft() {
  draftLoading.value = true;
  try {
    const res = await getPayrollDraft({ 
      payroll_type: payrollType.value,
      company_id : selectcompany.value
    });
    draft.value = res.data.data || [];
    console.log(res)
  } catch (e) {
    ElMessage.error(
      e.response?.data?.message || "Failed to load payroll draft",
    );
  } finally {
    draftLoading.value = false;
  }
}

function confirmPayroll() {
  if (!selectedRows.value.length)
    return ElMessage.warning("Select employees to pay");
  payForm.payroll_type = payrollType.value;
  payDialog.value = true;
}

async function submitPayroll() {
  paying.value = true;
  try {
    const payrolls = selectedRows.value.map((row) => ({
      user_id: row.user_id,
      basic_salary: row.basic_salary,
      half_salary: row.half_salary,
      other: "0",
      total_work_day: row.total_work_day,
      total_deduction: row.total_deduction,
      net_salary: row.net_salary,
      payroll_type: payForm.payroll_type,
      payroll_date: payForm.payroll_date,
      status: 0,
      note: payForm.note,
      attendance_id: row.unpaid_attendance || [],
      leave_id: row.leave_id || []
    }));
    await createPayroll({ payrolls });
    ElMessage.success("ប្រាក់ខែបេីកបានជោគជ័យ");
    payDialog.value = false;
    selectedRows.value = [];
    fetchDraft();
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to process payroll");
  } finally {
    paying.value = false;
  }
}

onMounted(()=>{
  fetchCompany();
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
:deep(.el-input-group__append) {
  background-color: #409eff;
  color: white;
  border-color: #409eff;
}

</style>
