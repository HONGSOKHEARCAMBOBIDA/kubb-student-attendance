<script setup>
import { ref, reactive, onMounted, watch, markRaw } from "vue";
import { ElMessage ,ElMessageBox, ElNotification} from "element-plus";
import { getincome,deleteincome } from "../api/services";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppInput from "../../components/AppInput.vue";
import { View } from "@element-plus/icons-vue";
import { Delete } from "@element-plus/icons-vue";
const incomes = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const filters = reactive({ name: "", income_date: "" });

async function fetchincome() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.income_date) params.income_date = filters.income_date;
    const res = await getincome(params);
    incomes.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load income");
  } finally {
    loading.value = false;
  }
}

async function handleDeleteIncome(row) {
  try {
    await ElMessageBox.confirm(
      `តើអ្នកពិតជាចង់លុប ${row.total}$ មែនទេ?`,
      "បញ្ជាក់ការលុប",
      {
        confirmButtonText: "លុប",
        cancelButtonText: "បោះបង់",
        type: "error",
        icon: markRaw(Delete),
        confirmButtonClass: "el-button--danger",
        distinguishCancelAndClose: true,
      }
    );

    await deleteincome(row.id);

    ElMessage.success("លុបបានជោគជ័យ");
    fetchincome();
  } catch (e) {
    // អ្នកប្រើចុចបោះបង់ មិនចាំបាច់បង្ហាញ Error ទេ
    if (e === "cancel" || e === "close") return;

    ElNotification.error({
      title: "កំហុស",
      message: e.response?.data?.error || "មិនអាចលុបថ្លៃសិក្សាបានទេ",
      offset: 100,
    });
  }
}

// Debounce so typing in the name field doesn't fire a request per keystroke
let timer = null;
watch(
  () => [filters.name, filters.income_date],
  () => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      page.value = 1; // go back to first page when filter changes
      fetchincome();
    }, 400);
  }
);

onMounted(() => {
  fetchincome();
});
</script>

<template>
  <AppFilterBar :fields="[
    { slot: 'name', span: 5 },
    { slot: 'invoice_date', span: 5 },
  ]">
    <template #name>
      <AppInput v-model="filters.name" placeholder="ស្វែងរកតាមឈ្មោះ" clearable />
    </template>
    <template #invoice_date>
      <el-date-picker v-model="filters.income_date" type="date" value-format="YYYY-MM-DD" placeholder="ជ្រើសរើសថ្ងៃ"
        clearable style="width: 100%" size="large" />
    </template>
  </AppFilterBar>

  <el-card>
    <AppTable :data="incomes" :loading="loading" show-index v-model:current-page="page" v-model:page-size="pageSize"
      :total="total" @page-change="fetchincome" actions-width="90px" :columns="[
        { prop: 'income_code', label: 'កូដ', minWidth: 130 },
        { prop: 'customer_name', label: 'អតិថិជន', minWidth: 130 },
        { prop: 'income_date', label: 'ថ្ងៃបង់ប្រាក់', minWidth: 130 },
        { prop: 'due_date', label: 'ថ្ងៃត្រូវបង់', minWidth: 130 },
        { prop: 'subtotal', label: 'សរុបរង', minWidth: 130, align: 'center' },
        { prop: 'tax', label: 'ពន្ធ', minWidth: 130, align: 'center' },
        { prop: 'discount', label: 'បញ្ចុះតម្លៃ', minWidth: 130, align: 'center' },
        { slot: 'total', label: 'សរុបត្រូវបង់', minWidth: 130, align: 'center' },
        { slot: 'paid', label: 'បានបង់', minWidth: 130, align: 'center' },
        { prop: 'status', label: 'ស្ថានភាព', minWidth: 70, align: 'center' },
        { prop: 'description', label: 'ផ្សេងៗ', minWidth: 130, align: 'center' },
        { slot: 'details', label: 'ព័ត៌មាន', width: 100, align: 'center' },
      ]">
      <template #actions="{ row }">
        <el-tooltip content="លុប" placement="top">
          <AppButton size="small" icon="Delete" type="danger" circle @click="handleDeleteIncome(row)"/>
        </el-tooltip>
      </template>
      <template #total="{ row }">
        <el-statistic :value="row.total" :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '20px', color: '#000000' }" />
      </template>
      <template #paid="{ row }">
        <el-statistic :value="row.paid" :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '20px', color: '#75c406' }" />
      </template>
      <template #details="{ row }">
        <el-popover placement="left" :width="320" trigger="click">
          <template #reference>
            <AppButton circle size="small" icon="View" type="success" plain>
            </AppButton>
          </template>

          <div class="income-details">
            <strong>សេវាកម្ម</strong>
            <div v-for="(item, i) in row.income_items || []" :key="i" class="detail-row">
              <span>{{ item.income_category_name || item.description || '—' }}</span>
              <b>{{ (item.quantity * item.unit_price).toFixed(2) }} $</b>
            </div>

            <el-divider />

            <strong>ប្រវត្តិបង់ប្រាក់</strong>
            <div v-for="(p, i) in row.income_payments || []" :key="i" class="detail-row">
              <span>{{ p.payment_code }} · បង់តាមរយៈ {{ p.method }}</span>
              <b class="paid">${{ Number(p.amount).toFixed(2) }}</b>
            </div>

            <el-empty v-if="!row.income_items?.length && !row.income_payments?.length" description="មិនមានព័ត៌មានលម្អិត"
              :image-size="55" />
          </div>
        </el-popover>
      </template>
    </AppTable>
  </el-card>
</template>

<style scoped>
.income-details {
  font-size: 13px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.detail-row span {
  overflow-wrap: anywhere;
}

.paid {
  color: #389e0d;
  white-space: nowrap;
}
</style>