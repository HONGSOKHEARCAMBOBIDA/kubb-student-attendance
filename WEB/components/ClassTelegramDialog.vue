<template>
  <AppDialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="កែប្រែ Telegram"
    width="600px"
  >
    <el-form :model="telegramForm" ref="telegramFormRef" label-position="top">
      <el-divider />
      <p class="section-label">Telegram</p>
      <div class="form-row">
        <el-form-item label="Bot Token" prop="bot_token">
          <el-input v-model="telegramForm.bot_token" placeholder="Telegram bot token" size="large" />
        </el-form-item>
        <el-form-item label="Group Link" prop="group_link">
          <el-input v-model="telegramForm.group_link" placeholder="Telegram group chat ID" size="large" />
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <AppButton @click="$emit('update:modelValue', false)" size="large" :block="false">
        ថតក្រោយ
      </AppButton>
      <AppButton type="primary" :loading="saving" @click="handleUpdate" size="large" :block="false">
        កែប្រែ
      </AppButton>
    </template>
  </AppDialog>
</template>

<script setup>
import { reactive, ref, watch } from "vue";
import { updateClassTelegram } from "../src/api/services.js";
import AppButton from "./AppButton.vue";
import AppDialog from "./AppDialog.vue";
import { useNotification } from "../composables/useNotification.js";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  classId: { type: [Number, String], default: null },
  editRow: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const notify = useNotification();
const telegramFormRef = ref();
const saving = ref(false);

const telegramForm = reactive({
  bot_token: "",
  group_link: "",
});

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return;
    const row = props.editRow || {};
    telegramForm.bot_token = row.bot_token || "";
    telegramForm.group_link = row.group_chatID || row.group_link || "";
  },
);

async function handleUpdate() {
  try {
    await updateClassTelegram(props.classId, {
      bot_token: telegramForm.bot_token,
      group_link: telegramForm.group_link,
    });
    notify.success("កែប្រែបានជោគជ័យ");
    emit("update:modelValue", false);
    emit("saved");
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
  }
}
</script>

<style scoped>
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
