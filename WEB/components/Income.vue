<script setup>
import { ref, reactive, watch, computed } from "vue";
import { ElMessage } from "element-plus";
import { getincomecategory, addincome } from "../src/api/services";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppInput from "./AppInput.vue";

const props = defineProps({
  studentRaw: { type: Object, default: () => ({}) },
  modelValue: { type: Boolean, default: false },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const incomecategories = ref([]);
const saving = ref(false);

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});

const newItem = () => ({
  income_category_id: null,
  description: "",
  qty: 1,
  unit_price: 0,
});

const form = reactive({
  due_date: "",
  tax: 0,
  discount: 0,
  description: "",
  items: [newItem()],
  payNow: false,
  payment: { amount: 0, method: "សាច់ប្រាក់", reference: "" },
});

const subtotal = computed(() =>
  form.items.reduce(
    (s, i) => s + (Number(i.qty) || 0) * (Number(i.unit_price) || 0),
    0,
  ),
);

const total = computed(
  () => subtotal.value + (Number(form.tax) || 0) - (Number(form.discount) || 0),
);
const balance = computed(
  () => total.value - (form.payNow ? Number(form.payment.amount) || 0 : 0),
);
const money = (v) => `$${(Number(v) || 0).toFixed(2)}`;

function addRow() {
  form.items.push(newItem());
}

function removeRow(i) {
  if (form.items.length > 1) form.items.splice(i, 1);
}

function reset() {
  Object.assign(form, {
    due_date: "",
    tax: 0,
    discount: 0,
    description: "",
    items: [newItem()],
    payNow: false,
    payment: { amount: 0, method: "សាច់ប្រាក់", reference: "" },
  });
}

async function loadCategories() {
  try {
    const res = await getincomecategory();
    incomecategories.value = (res.data.data || []).map((f) => ({
      label: f.name,
      value: f.id,
    }));
  } catch {
    ElMessage.error("Failed to load categories");
  }
}

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    reset();
    if (!incomecategories.value.length) loadCategories();
  },
);

// when "pay now" is switched on, default the amount to the full total
watch(
  () => form.payNow,
  (on) => {
    form.payment.amount = on ? Number(total.value.toFixed(2)) : 0;
  },
);

function validate() {
  if (!props.studentRaw?.id) return "មិនមានអតិថិជន";
  if (!form.items.length) return "សូមបន្ថែមមុខទំនិញ";
  for (const it of form.items) {
    if (!it.income_category_id) return "សូមជ្រើសប្រភេទ";
    if (!(Number(it.qty) > 0)) return "ចំនួនត្រូវធំជាង 0";
    if (Number(it.unit_price) < 0) return "តម្លៃមិនត្រឹមត្រូវ";
  }
  if (total.value < 0) return "សរុបមិនអាចតិចជាង 0";
  if (form.payNow) {
    const amt = Number(form.payment.amount);
    if (!(amt > 0)) return "ចំនួនទឹកប្រាក់បង់ត្រូវធំជាង 0";
    if (amt > total.value) return "ទឹកប្រាក់បង់លើសចំនួនសរុប";
  }
  return null;
}

async function handleSave() {
  const err = validate();
  if (err) return ElMessage.warning(err);

  const payload = {
    customer_id: props.studentRaw.id,
    due_date: form.due_date || "",
    tax: Number(form.tax) || 0,
    discount: Number(form.discount) || 0,
    description: form.description || null,
    income_items: form.items.map((i) => ({
      income_category_id: i.income_category_id,
      description: i.description,
      qty: Number(i.qty),
      unit_price: Number(i.unit_price),
    })),
  };
  if (form.payNow) {
    payload.payment = {
      amount: Number(form.payment.amount),
      method: form.payment.method,
      reference: form.payment.reference || null,
      description: null,
    };
  }

  saving.value = true;
  try {
    await addincome(payload);
    ElMessage.success("រក្សាទុកបានជោគជ័យ");
    emit("saved");
    visible.value = false;
  } catch (e) {
    ElMessage.error(
      e.response?.data?.message || e.response?.data?.error || "Failed to save",
    );
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <AppDialog v-model="visible" title="បង់ថ្លៃសេវាផ្សេងៗ" width="900px">
    <div class="student-info">
      <el-text tag="b">{{ studentRaw?.name_kh }}</el-text>
      <el-text type="info"> ({{ studentRaw?.code }})</el-text>
      <el-text tag="b">
        ភេទ: {{ studentRaw?.gender === 1 ? "ប្រុស" : "ស្រី" }}</el-text
      >
    </div>

    <el-form label-position="top">
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item label="ថ្ងៃត្រូវបង់">
            <el-date-picker
              v-model="form.due_date"
              type="date"
              value-format="YYYY-MM-DD"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ពិពណ៌នា">
            <el-input v-model="form.description" size="large" clearable />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <el-form label-position="top">
      <el-table :data="form.items" border size="default">
        <el-table-column label="ប្រភេទ" min-width="190">
          <template #default="{ row }">
            <AppSelect
              v-model="row.income_category_id"
              placeholder="ជ្រើសប្រភេទ"
              label="ប្រភេទចំណូល"
              size="large"
              filterable
              :options="incomecategories"
            >
            </AppSelect>
          </template>
        </el-table-column>
        <el-table-column label="ពិពណ៌នា" min-width="170">
          <template #default="{ row }">
            <AppInput v-model="row.description" label="Note"> </AppInput>
          </template>
        </el-table-column>
        <el-table-column label="ចំនួន" width="120">
          <template #default="{ row }">
            <el-input-number
              v-model="row.qty"
              :min="0.01"
              :precision="2"
              :controls="false"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column label="តម្លៃ" width="130">
          <template #default="{ row }">
            <el-input-number
              v-model="row.unit_price"
              :min="0"
              :precision="2"
              :controls="false"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column label="សរុប" width="100" align="right">
          <template #default="{ row }">{{
            money(row.qty * row.unit_price)
          }}</template>
        </el-table-column>
        <el-table-column width="60" align="center">
          <template #default="{ $index }">
            <AppButton
              size="small"
              icon="Delete"
              type="danger"
              circle
              :block="false"
              :disabled="form.items.length === 1"
              @click="removeRow($index)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-form>

    <AppButton
      class="add-row"
      size="small"
      type="primary"
      :block="true"
      plain
      @click="addRow"
    >
      បន្ថែមមុខទំនិញ
    </AppButton>

    <!-- totals -->
    <el-row :gutter="16" class="mt">
      <el-col :xs="24" :sm="12">
        <el-form label-position="top">
          <AppInput
            label="ពន្ធ"
            v-model="form.tax"
            type="number"
            :controls="false"
            placeholder="ពន្ធ"
          >
          </AppInput>
          <AppInput
            label="បញ្ចុះតម្លៃ"
            v-model="form.discount"
            type="number"
            :controls="false"
            placeholder="ពន្ធ"
          >
          </AppInput>
        </el-form>
      </el-col>
      <el-col :xs="24" :sm="12">
        <div class="summary">
          <div>
            <span>សរុបរង</span><b>{{ money(subtotal) }}</b>
          </div>
          <div>
            <span>ពន្ធ</span><b>{{ money(form.tax) }}</b>
          </div>
          <div>
            <span>បញ្ចុះតម្លៃ</span><b>-{{ money(form.discount) }}</b>
          </div>
          <div class="grand">
            <span>សរុប</span><b>{{ money(total) }}</b>
          </div>
          <div v-if="form.payNow">
            <span>នៅខ្វះ</span><b>{{ money(balance) }}</b>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- payment -->
    <el-divider />
    <el-checkbox size="large" border v-model="form.payNow"
      >បង់ឥឡូវនេះ?</el-checkbox
    >

    <el-row v-if="form.payNow" :gutter="10" class="mt">
      <el-col :span="8">
        <AppInput
          label="ទឹកប្រាក់បង់"
          label-position="top"
          v-model="form.payment.amount"
          type="number"
          :controls="false"
        >
        </AppInput>
      </el-col>
      <el-col :span="8">
        <AppInput
          label="វិធីបង់"
          label-position="top"
          v-model="form.payment.method"
        >
        </AppInput>
      </el-col>
      <el-col :span="8">
        <AppInput
          label="លេខយោង"
          label-position="top"
          v-model="form.payment.reference"
          clearable
        >
        </AppInput>
      </el-col>
    </el-row>

    <template #footer>
      <AppButton
        @click="visible = false"
        size="large"
        :block="true"
        type="warning"
        plain
        >បោះបង់</AppButton
      >
      <AppButton
        @click="handleSave"
        type="primary"
        :loading="saving"
        size="large"
        :block="true"
      >
        រក្សាទុក
      </AppButton>
    </template>
  </AppDialog>
</template>

<style scoped>
.student-info {
  margin-bottom: 12px;
}
.add-row {
  margin-top: 10px;
}
.mt {
  margin-top: 16px;
}
.summary {
  background: var(--el-fill-color-light);
  border-radius: 8px;
  padding: 12px 16px;
}
.summary > div {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
}
.summary .grand {
  border-top: 1px dashed var(--el-border-color);
  margin-top: 4px;
  padding-top: 8px;
  font-size: 24px;
  color: darkblue;
}
</style>
