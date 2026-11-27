<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from "vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import {
  addschoolarship,
  updateschoolarship,
  getSchoolarship,
} from "../api/services";
import type {
  DiscountType,
  Schoolarhisp,
  SchoolarshipRequestCreate,
  SchoolarshipRequestUpdate,
} from "../types/schoolarship";

const schoolarships = ref<Schoolarhisp[]>([]);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const dialogvisible = ref<boolean>(false);
const isEdit = ref<boolean>(false);
const editID = ref<number | null>(null);
const schoolarshipRef = ref<FormInstance>();

const discounttype: { value: DiscountType; label: string }[] = [
  { value: "percentage", label: "ភាគរយ" },
  { value: "amount", label: "ចំនួនទឹកប្រាក់" },
];

const createschoolarshiptform = reactive<SchoolarshipRequestCreate>({
  code: "",
  name: "",
  discount_type: null,
  discount_percentage: null,
  discount_amount: null,
  description: null,
});

const updateschoolarshipform = reactive<SchoolarshipRequestUpdate>({
  code: "",
  name: "",
  discount_type: null,
  discount_percentage: null,
  discount_amount: null,
  description: null,
  active: null,
});

const currentForm = computed(() =>
  isEdit.value ? updateschoolarshipform : createschoolarshiptform,
);

const rules = computed<FormRules>(() => ({
  code: [{ required: true, message: "សូមបញ្ចូលលេខកូដ", trigger: "blur" }],
  name: [{ required: true, message: "សូមបញ្ចូលឈ្មោះ", trigger: "blur" }],
  discount_type: [
    { required: true, message: "សូមជ្រើសរើសប្រភេទ", trigger: "change" },
  ],
  discount_percentage: [
    {
      required: currentForm.value.discount_type === "percentage",
      message: "សូមបញ្ចូលភាគរយ",
      trigger: "blur",
    },
    {
      type: "number",
      min: 0,
      max: 100,
      message: "ភាគរយត្រូវនៅចន្លោះ 0 ដល់ 100",
      trigger: "blur",
    },
  ],
  discount_amount: [
    {
      required: currentForm.value.discount_type === "amount",
      message: "សូមបញ្ចូលចំនួនទឹកប្រាក់",
      trigger: "blur",
    },
    { type: "number", min: 0, message: "ចំនួនត្រូវធំជាង 0", trigger: "blur" },
  ],
}));


watch(
  () => currentForm.value.discount_type,
  (type) => {
    if (type === "percentage") currentForm.value.discount_amount = null;
    else if (type === "amount") currentForm.value.discount_percentage = null;
  },
);

const columns = [
  { prop: "code", label: "លេខកូដ" },
  { prop: "name", label: "ឈ្មោះ" },
  { prop: "discount_type", label: "ប្រភេទ" },
  { prop: "discount_percentage", label: "បញ្ចុះភាគរយ" },
  { prop: "discount_amount", label: "បញ្ចុះចំនួន" },
  { prop: "description", label: "ផ្សេងៗ" },
];

const fetchSchoolarship = async (): Promise<void> => {
  loading.value = true;
  try {
    const res = await getSchoolarship();
    schoolarships.value = res.data.data;
  } catch (e) {
    ElMessage.error(`${e}`);
  } finally {
    loading.value = false;
  }
};

const createschoolarshipdialog = () => {
  isEdit.value = false;
  editID.value = null;
  Object.assign(createschoolarshiptform, {
    code: "",
    name: "",
    discount_type: null,
    discount_percentage: null,
    discount_amount: null,
    description: null,
  });
  dialogvisible.value = true;
};

const updateschoolarshipdialog = (row: Schoolarhisp) => {
  isEdit.value = true;
  editID.value = row.id;
  Object.assign(updateschoolarshipform, {
    code: row.code,
    name: row.name,
    discount_type: row.discount_type,
    discount_percentage: row.discount_percentage,
    discount_amount: row.discount_amount,
    description: row.description,
    active: row.active, // FIX: was missing
  });
  dialogvisible.value = true;
};

const resetForm = () => {
  schoolarshipRef.value?.clearValidate();
};

const handlecreateschoolarshipt = async (): Promise<void> => {
  if (!schoolarshipRef.value) return;
  const valid = await schoolarshipRef.value.validate().catch(() => false);
  if (!valid) return;
  submitting.value = true;
  try {
    await addschoolarship({
      code: createschoolarshiptform.code.trim(),
      name: createschoolarshiptform.name.trim(),
      discount_type: createschoolarshiptform.discount_type,
      discount_percentage: createschoolarshiptform.discount_percentage,
      discount_amount: createschoolarshiptform.discount_amount,
      description: createschoolarshiptform.description,
    });
    ElMessage.success("បង្កើតអាហារូបករណ៍បានជោគជ័យ");
    dialogvisible.value = false;
    await fetchSchoolarship();
  } catch (e) {
    ElMessage.error(`${e}`);
  } finally {
    submitting.value = false;
  }
};

const handleupdateschoolarship = async (): Promise<void> => {
  if (!schoolarshipRef.value || editID.value === null) return;
  const valid = await schoolarshipRef.value.validate().catch(() => false);
  if (!valid) return;
  submitting.value = true;
  try {
    await updateschoolarship(editID.value, {
      code: updateschoolarshipform.code,
      name: updateschoolarshipform.name,
      discount_type: updateschoolarshipform.discount_type,
      discount_percentage: updateschoolarshipform.discount_percentage,
      discount_amount: updateschoolarshipform.discount_amount,
      description: updateschoolarshipform.description,
      active: updateschoolarshipform.active, // FIX: was `isactive`
    });
    ElMessage.success("កែប្រែអាហារូបករណ៍បានជោគជ័យ");
    dialogvisible.value = false;
    await fetchSchoolarship();
  } catch (e) {
    ElMessage.error(`${e}`);
  } finally {
    submitting.value = false;
  }
};

const handleSubmit = () =>
  isEdit.value ? handleupdateschoolarship() : handlecreateschoolarshipt();

onMounted(() => {
  fetchSchoolarship();
});
</script>

<template>
  <AppFilterBar show-back :fields="[]">
    <template #actions>
        <AppButton @click="createschoolarshipdialog">បង្កើតថ្មី</AppButton>
    </template>
  </AppFilterBar>

  <AppTable
    :data="schoolarships"
    :columns="columns"
    :loading="loading"
    :show-pagination="false"
    show-index
    actions-width="70"
  >

    <template #actions="{ row }">
      <AppButton @click="updateschoolarshipdialog(row as Schoolarhisp)">កែប្រែ</AppButton>
    </template>
  </AppTable>

  <AppDialog
    v-model="dialogvisible"
    :title="isEdit ? 'កែប្រែ' : 'បង្កើត'"
    :loading="submitting"
    @closed="resetForm"
    @confirm="handleSubmit"
    @cancel="dialogvisible = false"
    width="500px"
  >
    <!-- FIX: `res` -> `ref` -->
     
    <el-form
      ref="schoolarshipRef"
      label-position="top"
      :model="currentForm"
      :rules="rules"
    >
      <el-form-item label="លេខកូដ" prop="code">
        <el-input v-model="currentForm.code" placeholder="បញ្ចូលលេខកូដ" />
      </el-form-item>

      <el-form-item label="ឈ្មោះ" prop="name">
        <el-input v-model="currentForm.name" placeholder="បញ្ចូលឈ្មោះ" />
      </el-form-item>

      <el-form-item label="ប្រភេទ" prop="discount_type">
        <el-select
          v-model="currentForm.discount_type"
          placeholder="ជ្រើសរើសប្រភេទ"
          style="width: 100%"
        >
          <el-option
            v-for="t in discounttype"
            :key="t.value"
            :label="t.label"
            :value="t.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item
        v-if="currentForm.discount_type === 'percentage'"
        label="បញ្ចុះភាគរយ (%)"
        prop="discount_percentage"
      >
        <el-input-number
          v-model="currentForm.discount_percentage"
          :min="0"
          :max="100"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item
        v-if="currentForm.discount_type === 'amount'"
        label="បញ្ចុះចំនួន"
        prop="discount_amount"
      >
        <el-input-number
          v-model="currentForm.discount_amount"
          :min="0"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item label="ផ្សេងៗ" prop="description">
        <el-input
          v-model="currentForm.description"
          type="textarea"
          :rows="3"
          placeholder="បញ្ចូលការពិពណ៌នា"
        />
      </el-form-item>

      <el-form-item v-if="isEdit" label="សកម្ម" prop="active">
        <el-switch v-model="updateschoolarshipform.active" />
      </el-form-item>

    <el-row :gutter="20">
      <el-col :span="12">
        <AppButton block @click="dialogvisible = false">ចាកចេញ</AppButton>
      </el-col>
      <el-col :span="12">
        <AppButton block :loading="submitting" @click="handleSubmit">
          បញ្ជូន
        </AppButton>
      </el-col>
    </el-row>
    </el-form>
  </AppDialog>
</template>

<style scoped></style>