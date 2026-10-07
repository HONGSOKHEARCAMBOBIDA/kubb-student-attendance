<script setup>
import { ref, reactive, watch, computed } from "vue";
import {
  getMajorPrice,
  getClass,
  addFee,
  getFeeschedule,
  getUserClasss,
  getSchoolarship,
} from "../src/api/services";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import AppFilterBar from "./AppFilterBar.vue";
import { useNotification } from "../composables/useNotification.js";
import AppTable from "./AppTable.vue";

const props = defineProps({
  studentRaw: { type: Object, default: () => ({}) },
  modelValue: { type: Boolean, default: false },
  generations: { type: Array, default: () => [] },
  programmes: { type: Array, default: () => [] },
  majors: { type: Array, default: () => [] },
});
const emit = defineEmits(["update:modelValue", "saved"]);
const notify = useNotification();
const title = computed(() =>
  props.studentRaw?.name_kh
    ? `និស្សិតឈ្មោះ ${props.studentRaw.name_kh} | ភេទ ${
        props.studentRaw.gender === 1 ? "ប្រុស" : "ស្រី"
      } | អត្តលេខ ${props.studentRaw.code}`
    : "បង់លុយនិស្សិត",
);
const installmentColumns = [
  { prop: 'sequence_no', label: 'លើកទី', minWidth: 120, align: 'center' },
  { prop: 'due_date', label: 'ថ្ងៃ-ខែ ត្រូវបង់', minWidth: 120, align: 'center' },
  { slot: 'amount', label: 'ចំនួនត្រូវបង់', minWidth: 120, align: 'center' },
  { prop: 'status', label: 'ស្ថានភាព', minWidth: 120, align: 'center' },
]
const saving = ref(false);
const loading = ref(false);
const majorprices = ref([]);
const classes = ref([]);
const feeschedules = ref([]);
const userClasses = ref([]);
const scholarships = ref([]);
const loadingUserClass = ref(false);
async function fetchUserClass() {
  userClasses.value = [];
  if (!props.studentRaw?.id) return;
  loadingUserClass.value = true;
  try {
    const res = await getUserClasss(props.studentRaw.id);
    userClasses.value = res.data.data || [];
  } catch (e) {
    notify.error(errMsg(e));
  } finally {
    loadingUserClass.value = false;
  }
}
const filters = reactive({
  major_id: null,
  generation_id: null,
  programme_id: null,
  year: null,
});

const form = reactive({
  class_id: null,
  major_price_id: null,
  fee_schedule_id: null,
  scholarship_id: null,
  date: new Date().toISOString().slice(0, 10),
});

const yearOptions = [1, 2, 3, 4].map((y) => ({
  label: `ឆ្នាំទី ${y}`,
  value: y,
}));

const classOptions = computed(
  () =>
    classes.value.map((c) => ({
      label: `ថ្នាក់${c.name} - ជំនាញ${c.major_name} - ${c.generation_name} - កម្រិត${c.programme_name} - វេន${c.shift_name} - ឆ្នាំ${c.year} - ឆមាស${c.semester}`,
      value: c.id,
    })), // adjust field if not `name`
);

const majorPriceOptions = computed(() =>
  majorprices.value.map((p) => ({
    label: `${p.year_fee}/១ឆ្នាំ | ${p.semester_fee}/១ឆមាស | ${p.quarter_fee}/៣ខែ | ${p.monthly_fee}/១ខែ`, // adjust to whatever best describes a price row
    value: p.id,
  })),
);

const selectedClass = computed(() =>
  classes.value.find((c) => c.id === form.class_id),
);

function errMsg(e, fallback = "មានបញ្ហាក្នុងការទាញទិន្នន័យ") {
  return e.response?.data?.error || fallback;
}

async function fetchFeeSchedule() {
  try {
    const res = await getFeeschedule();
    feeschedules.value = (res.data.data || []).map((f) => ({
      label: f.description,
      value: f.id, // was installment_count: backend needs the schedule ID
    }));
  } catch (e) {
    notify.error(errMsg(e));
  }
}

async function fetchSchoolarship() {
  try {
    const res = await getSchoolarship();
    scholarships.value = (res.data.data || []).map((f) => ({
      label: f.description,
      value: f.id, // was installment_count: backend needs the schedule ID
    }));
  } catch (e) {
    notify.error(errMsg(e));
  }
}

async function fetchClasses() {
  loading.value = true;
  try {
    const res = await getClass({
      major_id: filters.major_id || undefined,
      generation_id: filters.generation_id || undefined,
      programme_id: filters.programme_id || undefined,
      year: filters.year || undefined,
    });
    classes.value = res.data.data || [];
  } catch (e) {
    notify.error(errMsg(e, ""));
  } finally {
    loading.value = false;
  }
}

// Prices are driven by the selected class, so price and class always match
async function fetchMajorPrice(cls) {
  majorprices.value = [];
  form.major_price_id = null;
  if (!cls) return;
  try {
    const res = await getMajorPrice({
      major_id: cls.major_id || undefined,
      year: cls.year || undefined,
      programme_id: cls.programme_id || undefined,
      generation_id: cls.generation_id || undefined,
    });
    majorprices.value = res.data.data || [];
    if (majorprices.value.length === 1) {
      form.major_price_id = majorprices.value[0].id; // auto-select the only option
    }
  } catch (e) {
    notify.error(errMsg(e));
  }
}

function resetForm() {
  Object.assign(form, {
    class_id: null,
    major_price_id: null,
    fee_schedule_id: null,
    scholarship_id: null,
    date: new Date().toISOString().slice(0, 10),
  });
  majorprices.value = [];
}

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    Object.assign(filters, {
      year: null,
      generation_id: null,
      programme_id: null,
      major_id: null,
    });
    resetForm();
    fetchFeeSchedule();
    fetchUserClass();
    fetchSchoolarship();
  },
);

// Filters changed -> reload classes, clear dependent selections
watch(
  filters,
  () => {
    if (!props.modelValue) return;
    form.class_id = null;
    majorprices.value = [];
    form.major_price_id = null;
    fetchClasses();
  },
  { deep: true },
);

// Class changed -> load matching prices
watch(
  () => form.class_id,
  () => fetchMajorPrice(selectedClass.value),
);

function validate() {
  if (!props.studentRaw?.id) return "មិនមាននិស្សិត";
  if (!form.class_id) return "សូមជ្រើសរើសថ្នាក់";
  if (!form.major_price_id) return "សូមជ្រើសរើសតម្លៃមុខជំនាញ";
  if (!form.fee_schedule_id) return "សូមជ្រើសរើសរបៀបបង់ប្រាក់";
  if (!form.date) return "សូមជ្រើសរើសកាលបរិច្ឆេទ";
  return null;
}

async function submit() {
  const msg = validate();
  if (msg) return notify.error(msg);

  saving.value = true;
  try {
    await addFee({
      student_id: props.studentRaw.id,
      class_id: form.class_id,
      major_price_id: form.major_price_id,
      fee_schedule_id: form.fee_schedule_id,
      scholarship_id: form.scholarship_id || 0, // 0 = no scholarship (backend must handle)
      date: form.date,
    });
    notify.success?.("រក្សាទុកបានជោគជ័យ");
    emit("saved");
    emit("update:modelValue", false);
  } catch (e) {
    notify.error(errMsg(e, "មានបញ្ហាក្នុងការរក្សាទុក"));
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="title"
    width="95%"
  >
    <AppTable
      expandable
      :show-pagination="false"
      :data="userClasses"
      :loading="loadingUserClass"
      :columns="[
        { slot: 'class_name', label: 'ថ្នាក់', minWidth: 150, align: 'center' },
        { prop: 'major_name', label: 'ជំនាញ', minWidth: 120, align: 'center' },
        { prop: 'shift_name', label: 'វេន', minWidth: 120, align: 'center' },
        {
          slot: 'generation_name',
          label: 'ជំនាន់',
          minWidth: 120,
          align: 'center',
        },
        { slot: 'year', label: 'ឆ្នាំ', minWidth: 40, align: 'center' },
        { slot: 'semester', label: 'ឆមាស', minWidth: 50, align: 'center' },
        { slot: 'group_name', label: 'ក្រុម', minWidth: 40, align: 'center' },
        { slot: 'term', label: 'វគ្គ', minWidth: 40, align: 'center' },
        { slot: 'active', label: 'ស្ថានភាព', minWidth: 120, align: 'center' },
        {
          prop: 'schoolarship',
          label: 'អាហារូបករណ៍',
          minWidth: 120,
          align: 'center',
        },
        {
          prop: 'fee_schedule',
          label: 'របៀបបង់ប្រាក់',
          minWidth: 120,
          align: 'center',
        },
        { slot: 'amount', label: 'សរុប', minWidth: 120, align: 'center' },
        {
          slot: 'discount',
          label: 'បញ្ចុះតម្លៃ',
          minWidth: 120,
          align: 'center',
        },
        {
          slot: 'total',
          label: 'ប្រាក់មិនទាន់បង់',
          minWidth: 120,
          align: 'center',
        },
      ]"
    >
    <template #class_name="{row}">
      <el-text size="large">{{ row.class_name }} | {{ row.type === "onclass" ? "ផ្ទាល់" : "អនឡាញ" }} | <el-text type="danger">{{ row.programme_name }}</el-text></el-text>
    </template>
      <template #generation_name="{ row }">
        {{ row.generation_name }}
        <span v-if="row.generation_start || row.generation_end">
          ({{ row.generation_start }} - {{ row.generation_end }})
        </span>
      </template>
      <template #active="{ row }">
        <el-tag :type="row.active ? 'success' : 'info'" size="small">
          {{ row.active ? "កំពុងសិក្សា" : "បានបញ្ចប់" }}
        </el-tag>
      </template>
      <template #year="{row}">
         <el-statistic
          :value="row.year"
          :value-style="{ fontSize: '20px', color: '#000000' }"
        />
      </template>
            <template #semester="{row}">
         <el-statistic
          :value="row.semester"
          :value-style="{ fontSize: '20px', color: '#000000' }"
        />
      </template>
          <template #group_name="{row}">
         <el-statistic
          :value="row.group_name"
          :value-style="{ fontSize: '20px', color: '#000000' }"
        />
      </template>
          <template #term="{row}">
         <el-statistic
          :value="row.term"
          :value-style="{ fontSize: '20px', color: '#000000' }"
        />
      </template>
       <template #amount="{row}">
            <el-statistic
          :value="row.amount"
          :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
            <template #discount="{row}">
            <el-statistic
          :value="row.discount"
          :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '16px', color: '#000000' }"
        />
      </template>
      <template #total="{row}">
            <el-statistic
          :value="row.total"
          :formatter="(val) => `${Number(val).toFixed(2)}$`"
          :value-style="{ fontSize: '26px', color: '#FF0000' }"
        />
      </template>
<template #expand="{ row: userClass }">
  <AppTable
    :show-pagination="false"
    :data="userClass.installments ?? []"
    :columns="installmentColumns"
  >
    <template #amount="{ row: installment }">
      <el-statistic
        :value="installment.amount"
        :formatter="(val) => `${Number(val).toFixed(2)}$`"
        :value-style="{ fontSize: '20px', color: '#FF0000' }"
      />
    </template>
  </AppTable>
</template>
    </AppTable>
    <el-form label-position="top">
      <AppFilterBar
        :fields="[
          { slot: 'major', span: 6 },
          { slot: 'generation', span: 6 },
          { slot: 'programme', span: 6 },
          { slot: 'year', span: 6 },
        ]"
      >
        <template #major>
          <AppSelect
            v-model="filters.major_id"
            label="ជំនាញ"
            placeholder="ជំនាញ"
            clearable
            size="large"
            :options="majors"
            filterable
          />
        </template>
        <template #generation>
          <AppSelect
            v-model="filters.generation_id"
            label="ជំនាន់"
            placeholder="ជំនាន់"
            clearable
            size="large"
            :options="generations"
          />
        </template>
        <template #programme>
          <AppSelect
            v-model="filters.programme_id"
            label="កម្មវិធីសិក្សា"
            placeholder="កម្មវិធីសិក្សា"
            clearable
            size="large"
            :options="programmes"
          />
        </template>
        <template #year>
          <AppSelect
            v-model="filters.year"
            label="ឆ្នាំ"
            placeholder="ឆ្នាំ"
            clearable
            size="large"
            :options="yearOptions"
          />
        </template>
      </AppFilterBar>
    </el-form>

    <div class="form-grid">
      <AppSelect
        v-model="form.class_id"
        placeholder="ថ្នាក់"
        size="large"
        :options="classOptions"
      />
      <AppSelect
        v-model="form.major_price_id"
        placeholder="តម្លៃមុខជំនាញ"
        size="large"
        :options="majorPriceOptions"
        :disabled="!form.class_id"
      />
      <AppSelect
        v-model="form.fee_schedule_id"
        placeholder="របៀបបង់ប្រាក់"
        size="large"
        :options="feeschedules"
      />
      <AppSelect
        v-model="form.scholarship_id"
        placeholder="អាហារូបករណ៍ (បើមាន)"
        clearable
        size="large"
        :options="scholarships"
      />
      <el-date-picker
        v-model="form.date"
        type="date"
        value-format="YYYY-MM-DD"
        placeholder="កាលបរិច្ឆេទ"
        size="large"
        style="width: 100%"
      />
    </div>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)">បោះបង់</AppButton>
      <AppButton type="primary" :loading="saving" @click="submit"
        >រក្សាទុក</AppButton
      >
    </template>
  </AppDialog>
</template>

<style scoped>
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-top: 16px;
}
</style>
