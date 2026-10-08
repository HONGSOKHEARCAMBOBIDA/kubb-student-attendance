<script setup>
import { ref, reactive, computed, watch } from "vue";
import { addFeeTransaction } from "../src/api/services";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import AppSelect from "./AppSelect.vue";
import { useNotification } from "../composables/useNotification.js";
import AppInput from "./AppInput.vue";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  installment: { type: Object, default: null },
});

const emit = defineEmits(["update:modelValue", "saved"]);
const notify = useNotification();
const saving = ref(false);

const today = () => new Date().toISOString().slice(0, 10);

const form = reactive({
  date: today(),
  amount: 0,
  discount: 0,
  tax: 0,
  reference: "",
  method: null,
  message: "",
  description: "",
});

const total = computed(() =>
  Math.max(0, +(form.amount - form.discount + form.tax).toFixed(2)),
);

watch(
  () => props.modelValue,
  (open) => {
    if (!open || !props.installment) return;
    Object.assign(form, {
      date: today(),
      amount: Number(props.installment.amount) || 0,
      discount: 0,
      tax: 0,
      reference: "",
      method: null,
      message: "",
      description: "",
    });
  },
);

function validate() {
  if (!props.installment?.id) return "មិនមានលើកបង់ប្រាក់";
  if (!form.date) return "សូមជ្រើសរើសកាលបរិច្ឆេទ";
  if (!form.method) return "សូមជ្រើសរើសវិធីបង់ប្រាក់";
  if (form.amount <= 0) return "ចំនួនទឹកប្រាក់មិនត្រឹមត្រូវ";
  if (form.discount < 0 || form.tax < 0) return "តម្លៃមិនអាចអវិជ្ជមាន";
  if (form.discount > form.amount) return "បញ្ចុះតម្លៃលើសចំនួនទឹកប្រាក់";
  return null;
}

async function submit() {
  const msg = validate();
  if (msg) return notify.error(msg);

  saving.value = true;
  try {
    await addFeeTransaction({
      installment_id: props.installment.id,
      date: form.date,
      amount: form.amount,
      discount: form.discount,
      tax: form.tax,
      total: total.value,
      reference: form.reference || null,
      method: form.method,
      message: form.message || null,
      description: form.description || null,
    });
    notify.success?.("បង់ប្រាក់បានជោគជ័យ");
    emit("saved");
    emit("update:modelValue", false);
  } catch (e) {
    notify.error(e.response?.data?.error || "មានបញ្ហាក្នុងការបង់ប្រាក់");
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="បង់ប្រាក់"
    width="600px"
    append-to-body
  >
    <el-descriptions v-if="installment" :column="2" border size="small">
      <el-descriptions-item label="លើកទី">
        {{ installment.sequence_no }}
      </el-descriptions-item>
      <el-descriptions-item label="ថ្ងៃ-ខែ ត្រូវបង់">
        {{ installment.due_date }}
      </el-descriptions-item>
    </el-descriptions>

    <el-form label-position="top" class="pay-form">
      <el-form-item label="កាលបរិច្ឆេទ">
        <el-date-picker
          v-model="form.date"
          type="date"
          value-format="YYYY-MM-DD"
          size="large"
          style="width: 100%"
        />
      </el-form-item>

      <AppInput
        label="វិធីបង់ប្រាក់"
        v-model="form.method"
        placeholder="វិធីបង់ប្រាក់"
        size="large"
      >
      </AppInput>
      <AppInput
        label="ចំនួនទឹកប្រាក់ ($)"
        v-model="form.amount"
        :min="0"
        :precision="2"
        :controls="false"
        size="large"
        style="width: 100%"
        disabled
      >
      </AppInput>
      <AppInput
        label="បញ្ចុះតម្លៃ ($)"
        v-model="form.discount"
        :min="0"
        :precision="2"
        :controls="false"
        size="large"
        style="width: 100%"
      >
      </AppInput>
      <AppInput
        label="ពន្ធ ($)"
        v-model="form.tax"
        :min="0"
        :precision="2"
        :controls="false"
        size="large"
        class="full"
      >
      </AppInput>
      <AppInput
        label="សារ"
        v-model="form.message"
        type="textarea"
        :rows="2"
        class="full"
      >
      </AppInput>
      <AppInput
        label="ពិពណ៌នា"
        class="full"
        v-model="form.description"
        type="textarea"
        :rows="2"
      >
      </AppInput>
    </el-form>

    <div class="total-box">
      <span>សរុបត្រូវបង់</span>
      <el-statistic
        :value="total"
        :formatter="(v) => `${Number(v).toFixed(2)}$`"
        :value-style="{ fontSize: '28px', color: '#FF0000' }"
      />
    </div>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)">បោះបង់</AppButton>
      <AppButton type="success" :loading="saving" @click="submit">
        បង់ប្រាក់
      </AppButton>
    </template>
  </AppDialog>
</template>

<style scoped>
.pay-form {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  column-gap: 16px;
  margin-top: 16px;
}
.pay-form .full {
  grid-column: 1 / -1;
}
.total-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  padding: 12px 16px;
  background: #fff5f5;
  border-radius: 8px;
}
</style>
