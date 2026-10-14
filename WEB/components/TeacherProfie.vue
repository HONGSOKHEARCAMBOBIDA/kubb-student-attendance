<script setup>
import { ref, reactive, watch } from "vue";
import { ElMessage } from "element-plus";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";
import {
  getTeacherProfile, // GET  /teacher-profile/by-user/:userId (assumed name)
  addTeacherProfile, // POST /teacher-profile                  (assumed name)
  updateTeacherProfile, // PUT /teacher-profile                (assumed name)
} from "../src/api/services";

const props = defineProps({
  modelValue: Boolean,
  teacherRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const formRef = ref();
const saving = ref(false);
const loading = ref(false);
const exists = ref(false); // false => create, true => update

const form = reactive({
  date_of_birth: "",
  place_of_birth: "",
  nationality: "",
  address: "",
  phone: "",
});

const rules = {
  date_of_birth: [
    { required: true, message: "សូមជ្រើសថ្ងៃកំណើត", trigger: "change" },
  ],
};

function resetForm() {
  exists.value = false;
  Object.assign(form, {
    date_of_birth: "",
    place_of_birth: "",
    nationality: "",
    address: "",
    phone: "",
  });
}

async function loadProfile() {
  resetForm();
  if (!props.teacherRaw) return;
  loading.value = true;
  try {
    let p = null;
    try {
      const res = await getTeacherProfile(props.teacherRaw.id);
      p = res.data.data;
    } catch (e) {
      if (e.response?.status !== 404) throw e; // 404 => no profile yet
    }
    if (!p || !p.id) return;

    exists.value = true;
    Object.assign(form, {
      date_of_birth: p.date_of_birth ?? "",
      place_of_birth: p.place_of_birth ?? "",
      nationality: p.nationality ?? "",
      address: p.address ?? "",
      phone: p.phone ?? "",
    });
  } catch {
    ElMessage.error("Failed to load teacher profile");
  } finally {
    loading.value = false;
  }
}

watch(
  () => props.modelValue,
  (open) => open && loadProfile(),
);

async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  try {
    const payload = {
      user_id: props.teacherRaw.id,
      date_of_birth: form.date_of_birth || null,
      place_of_birth: form.place_of_birth || null,
      nationality: form.nationality || null,
      address: form.address || null,
      phone: form.phone || null,
    };
    if (exists.value) {
      await updateTeacherProfile(payload);
    } else {
      await addTeacherProfile(payload);
    }
    ElMessage.success("រក្សាទុកបានជោគជ័យ");
    emit("saved");
    emit("update:modelValue", false);
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
  <AppDialog
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :title="`ព័ត៌មានលំអិតសាស្ត្រាចារ្យ${teacherRaw ? ' - ' + teacherRaw.name_kh : ''}`"
    width="40%"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      v-loading="loading"
    >
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item label="ថ្ងៃកំណើត" prop="date_of_birth">
            <el-date-picker
              v-model="form.date_of_birth"
              type="date"
              value-format="YYYY-MM-DD"
              size="large"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <AppInput
            label="ទីកន្លែងកំណើត"
            prop="place_of_birth"
            clearable
            v-model="form.place_of_birth"
          />
        </el-col>
        <el-col :xs="24" :sm="12">
          <AppInput
            label="សញ្ជាតិ"
            prop="nationality"
            clearable
            v-model="form.nationality"
          />
        </el-col>
        <el-col :xs="24" :sm="12">
          <AppInput
            label="លេខទូរស័ព្ទ"
            prop="phone"
            clearable
            v-model="form.phone"
          />
        </el-col>
        <el-col :span="24">
          <AppInput
            label="អាសយដ្ឋាន"
            prop="address"
            clearable
            v-model="form.address"
          />
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <AppButton
        @click="emit('update:modelValue', false)"
        size="large"
        :block="false"
        type="warning"
      >
        បោះបង់
      </AppButton>
      <AppButton
        @click="handleSave"
        type="primary"
        :loading="saving"
        size="large"
        :block="false"
      >
        {{ exists ? "កែប្រែ" : "បង្កេីត" }}
      </AppButton>
    </template>
  </AppDialog>
</template>