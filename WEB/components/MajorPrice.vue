<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import {
  addMajorPrice,
  getMajorPrice,
  updateMajorPrice,
} from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppInput from "./AppInput.vue";
import { useNotification } from "../composables/useNotification.js";
import AppFilterBar from "./AppFilterBar.vue";
import AppTable from "./AppTable.vue";

const props = defineProps({
  raw: { type: Object, default: () => ({}) },
  isEdit: { type: Boolean, default: false },
  modelValue: { type: Boolean, default: false }, // list dialog (controlled by parent)
  majorId: { type: Number, default: null },
  generations: { type: Array, default: () => [] },
  programmes: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const formRef = ref();
const saving = ref(false);

/* ---------- list ---------- */
const majorprice = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const filters = reactive({
  year: "",
  generation_id: "",
  programme_id: "",
});

const editId = ref(null);

const addTitle = computed(() =>
  editId.value ? "កែប្រែថ្លៃសិក្សា" : "បន្ថែមថ្លៃសិក្សា",
);

const listTitle = computed(() =>
  props.raw?.name_kh
    ? `ថ្លៃសិក្សា - ជំនាញ${props.raw.name_kh} | មហាវិទ្យាល័យ${props.raw.faculty_name}`
    : "ថ្លៃសិក្សា",
);

async function fetchMajorPrice() {
  if (!props.majorId) return;
  loading.value = true;
  try {
    const res = await getMajorPrice({
      major_id: props.majorId, // only this major's prices
      page: page.value,
      page_size: pageSize.value,
      year: filters.year || undefined,
      programme_id: filters.programme_id || undefined,
      generation_id: filters.generation_id || undefined,
    });
    majorprice.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
    console.log(majorprice.value);
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការទាញទិន្នន័យ");
  } finally {
    loading.value = false;
  }
}

// Fetch every time the list dialog opens (component is mounted once in parent)
watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    page.value = 1;
    Object.assign(filters, { year: "", generation_id: "", programme_id: "" });
    fetchMajorPrice();
  },
);

// Refetch when filters change
watch(filters, () => {
  page.value = 1;
  fetchMajorPrice();
});

function closeList() {
  emit("update:modelValue", false);
}

/* ---------- add dialog ---------- */
const addVisible = ref(false);

const defaultForm = () => ({
  major_id: null,
  generation_id: null,
  programme_id: null,
  year: null,
  monthly_fee: 0,
  quarter_fee: 0,
  semester_fee: 0,
  year_fee: 0,
});
const form = reactive(defaultForm());

const yearOptions = [1, 2, 3, 4].map((y) => ({
  label: `ឆ្នាំទី ${y}`,
  value: y,
}));

const feeRule = (message) => [
  { required: true, message, trigger: "blur" },
  {
    validator: (_, v, cb) =>
      v === "" || v === null || v === undefined || Number(v) < 0
        ? cb(new Error("តម្លៃមិនត្រឹមត្រូវ"))
        : cb(),
    trigger: "blur",
  },
];

const rules = {
  generation_id: [
    { required: true, message: "សូមជ្រើសរើសជំនាន់", trigger: "change" },
  ],
  programme_id: [
    { required: true, message: "សូមជ្រើសរើសកម្មវិធីសិក្សា", trigger: "change" },
  ],
  year: [{ required: true, message: "សូមជ្រើសរើសឆ្នាំ", trigger: "change" }],
  monthly_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំខែ"),
  quarter_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំត្រីមាស"),
  semester_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំឆមាស"),
  year_fee: feeRule("សូមបញ្ចូលថ្លៃប្រចាំឆ្នាំ"),
};

function openAddDialog() {
  editId.value = null;
  Object.assign(form, defaultForm(), { major_id: props.majorId });
  addVisible.value = true;
  formRef.value?.clearValidate();
}

function openEditDialog(row) {
  editId.value = row.id;
  Object.assign(form, {
    major_id: props.majorId,
    generation_id: row.generation_id,
    programme_id: row.programme_id,
    year: row.year,
    monthly_fee: row.monthly_fee,
    quarter_fee: row.quarter_fee,
    semester_fee: row.semester_fee,
    year_fee: row.year_fee,
  });
  addVisible.value = true;
  formRef.value?.clearValidate();
}

async function handleSave() {
  if (!formRef.value) return;
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  const payload = {
    major_id: props.majorId,
    generation_id: form.generation_id,
    programme_id: form.programme_id,
    year: form.year,
    monthly_fee: Number(form.monthly_fee),
    quarter_fee: Number(form.quarter_fee),
    semester_fee: Number(form.semester_fee),
    year_fee: Number(form.year_fee),
  };

  saving.value = true;
  try {
    if (editId.value) {
      await updateMajorPrice(editId.value, payload);
      notify.success("កែប្រែបានជោគជ័យ");
    } else {
      await addMajorPrice(payload);
      notify.success("បង្កើតបានជោគជ័យ");
    }
    addVisible.value = false;
    await fetchMajorPrice();
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការរក្សាទុក");
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <!-- 1) List dialog: shown first -->
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="listTitle"
    width="75%"
  >
    <AppFilterBar
      :fields="[
        { slot: 'generation', span: 5 },
        { slot: 'programme', span: 5 },
        { slot: 'year', span: 5 },
      ]"
      :action-span="4"
    >
      <template #generation>
        <AppSelect
          v-model="filters.generation_id"
          placeholder="ជំនាន់"
          clearable
          size="large"
          :options="generations"
        />
      </template>
      <template #programme>
        <AppSelect
          v-model="filters.programme_id"
          placeholder="កម្មវិធីសិក្សា"
          clearable
          size="large"
          :options="programmes"
        />
      </template>
      <template #year>
        <AppSelect
          v-model="filters.year"
          placeholder="ឆ្នាំ"
          clearable
          size="large"
          :options="yearOptions"
        />
      </template>
      <template #actions>
        <AppButton type="primary" @click="openAddDialog"
          >បន្ថែមថ្លៃសិក្សា</AppButton
        >
      </template>
    </AppFilterBar>

    <AppTable
      show-index
      :data="majorprice"
      :loading="loading"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      @page-change="fetchMajorPrice"
      actions-width="90px"
      :columns="[
        { prop: 'major_name', label: 'ជំនាញ', minWidth: 120, align: 'center' },
        {
          prop: 'generation_name',
          label: 'ជំនាន់',
          minWidth: 120,
          align: 'center',
        },
        {
          prop: 'programme_name',
          label: 'កម្រិត',
          minWidth: 120,
          align: 'center',
        },
        { prop: 'year', label: 'ឆ្នាំ', minWidth: 120, align: 'center' },
        {
          slot: 'monthly_fee',
          label: '១ខែម្ដង',
          minWidth: 120,
          align: 'center',
        },
        {
          slot: 'quarter_fee',
          label: '៣ខែម្ដង',
          minWidth: 120,
          align: 'center',
        },
        {
          slot: 'semester_fee',
          label: '៦ខែម្ដង',
          minWidth: 120,
          align: 'center',
        },
        {
          slot: 'year_fee',
          label: '១ឆ្នាំម្ដង',
          minWidth: 120,
          align: 'center',
        },
      ]"
    >
      <template #monthly_fee="{ row }">
        <el-statistic
          :value="row.monthly_fee"
          :formatter="(val) => `${Number(val).toFixed(2)}$ / ១ខែ`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
      <template #quarter_fee="{ row }">
        <el-statistic
          :value="row.quarter_fee"
          :formatter="(val) => `${Number(val).toFixed(2)}$ / ៣ខែ`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
      <template #semester_fee="{ row }">
        <el-statistic
          :value="row.semester_fee"
          :formatter="(val) => `${Number(val).toFixed(2)}$ / ៦ខែ`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
      <template #year_fee="{ row }">
        <el-statistic
          :value="row.year_fee"
          :formatter="(val) => `${Number(val).toFixed(2)}$ / ១ឆ្នាំ`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="openEditDialog(row)"
          />
        </el-tooltip>
      </template>
    </AppTable>
    <template #footer>
      <AppButton @click="closeList">បិទ</AppButton>
    </template>
  </AppDialog>

  <!-- 2) Add dialog: only opens when user clicks "Add" -->
  <AppDialog v-model="addVisible" :title="addTitle" width="640px">
    <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
      <div class="form-row">
        <AppSelect
          label="ជំនាន់"
          prop="generation_id"
          v-model="form.generation_id"
          :options="generations"
          placeholder="ជ្រើសរើសជំនាន់"
          size="large"
          filterable
          clearable
        />
        <AppSelect
          label="កម្មវិធីសិក្សា"
          prop="programme_id"
          v-model="form.programme_id"
          :options="programmes"
          placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
          size="large"
          filterable
          clearable
        />
      </div>

      <div class="form-row">
        <AppSelect
          label="ឆ្នាំ"
          prop="year"
          v-model="form.year"
          :options="yearOptions"
          placeholder="ជ្រើសរើសឆ្នាំ"
          size="large"
        />
        <AppInput
          label="ថ្លៃប្រចាំខែ"
          prop="monthly_fee"
          v-model.number="form.monthly_fee"
          type="number"
          size="large"
        />
      </div>

      <div class="form-row">
        <AppInput
          label="ថ្លៃប្រចាំត្រីមាស"
          prop="quarter_fee"
          v-model.number="form.quarter_fee"
          type="number"
          size="large"
        />
        <AppInput
          label="ថ្លៃប្រចាំឆមាស"
          prop="semester_fee"
          v-model.number="form.semester_fee"
          type="number"
          size="large"
        />
      </div>

      <div class="form-row">
        <AppInput
          label="ថ្លៃប្រចាំឆ្នាំ"
          prop="year_fee"
          v-model.number="form.year_fee"
          type="number"
          size="large"
        />
        <div class="form-spacer" />
      </div>
    </el-form>

    <template #footer>
      <AppButton @click="addVisible = false">បោះបង់</AppButton>
      <AppButton type="primary" :loading="saving" @click="handleSave"
        >រក្សាទុក</AppButton
      >
    </template>
  </AppDialog>
</template>

<style scoped>
.form-row {
  display: flex;
  gap: 20px;
}
.form-row .el-form-item,
.form-spacer {
  flex: 1;
  min-width: 0;
}
</style>
