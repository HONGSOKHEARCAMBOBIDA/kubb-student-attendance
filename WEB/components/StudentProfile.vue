<script setup>
import { ref, reactive, watch } from "vue";
import { ElMessage } from "element-plus";
import AppDialog from "./AppDialog.vue";
import AppButton from "./AppButton.vue";
import AppInput from "./AppInput.vue";
import {
  getprovince,
  getdistrict,
  getcommunce,
  getvillage,
  getstudentcategory,
  getStudentProfile, // GET  /student-profile/by-user/:userId  (assumed name)
  addStudentProfile, // POST /student-profile                  (assumed name)
  updateStudentProfile, // PUT /student-profile/:id            (assumed name)
} from "../src/api/services";

const props = defineProps({
  modelValue: Boolean,
  studentRaw: { type: Object, default: null },
});
const emit = defineEmits(["update:modelValue", "saved"]);

const formRef = ref();
const saving = ref(false);
const loading = ref(false);
const profileId = ref(null); // null => create, otherwise update

const provinces = ref([]);
const districts = ref([]);
const communces = ref([]);
const villages = ref([]);
const studentcategory = ref([]);

const selectedprovinceid = ref(null);
const selecteddistrictid = ref(null);
const selectedcommunceid = ref(null);

const form = reactive({
  student_category_id: null,
  date_of_birth: "",
  nationality: "",
  phone: "",
  village_id: null,
  occupation: "",
  academic_stream: null,
  telegram_username: "",
});

const rules = {
  student_category_id: [
    { required: true, message: "សូមជ្រើសប្រភេទនិស្សិត", trigger: "change" },
  ],
  date_of_birth: [
    { required: true, message: "សូមជ្រើសថ្ងៃកំណើត", trigger: "change" },
  ],
  village_id: [{ required: true, message: "សូមជ្រើសភូមិ", trigger: "change" }],
};

const streams = [
  { label: "វិទ្យាសាស្ត្រ", value: "SCIENCE" },
  { label: "សង្គមវិទ្យា", value: "SOCIAL_SCIENCE" },
];

const toOptions = (list) =>
  (list || []).map((a) => ({
    label: `${a.name_kh} - ${a.name_en}`,
    value: a.id,
  }));

/* ---------- lookups (each one sets its own ref and returns nothing) ---------- */
async function fetchProvince() {
  try {
    provinces.value = toOptions((await getprovince()).data.data);
  } catch {
    ElMessage.error("Failed to load provinces");
  }
}
async function fetchStudentCategory() {
  try {
    const res = await getstudentcategory();
    studentcategory.value = (res.data.data || []).map((s) => ({
      label: `${s.name}`,
      value: s.id,
    }));
  } catch {
    ElMessage.error("Failed to load categories");
  }
}
async function fetchDistrict(id) {
  if (!id) return (districts.value = []);
  try {
    districts.value = toOptions((await getdistrict(id)).data.data);
  } catch {
    ElMessage.error("Failed to load districts");
  }
}
async function fetchCommunce(id) {
  if (!id) return (communces.value = []);
  try {
    communces.value = toOptions((await getcommunce(id)).data.data);
  } catch {
    ElMessage.error("Failed to load communes");
  }
}
async function fetchVillage(id) {
  if (!id) return (villages.value = []);
  try {
    villages.value = toOptions((await getvillage(id)).data.data);
  } catch {
    ElMessage.error("Failed to load villages");
  }
}

/* ---------- cascade handlers (user-driven changes only) ---------- */
async function provincechange() {
  selecteddistrictid.value = null;
  selectedcommunceid.value = null;
  form.village_id = null;
  communces.value = [];
  villages.value = [];
  await fetchDistrict(selectedprovinceid.value);
}
async function districtchange() {
  selectedcommunceid.value = null;
  form.village_id = null;
  villages.value = [];
  await fetchCommunce(selecteddistrictid.value);
}
async function communcechange() {
  form.village_id = null;
  await fetchVillage(selectedcommunceid.value);
}

/* ---------- open / load ---------- */
function resetForm() {
  profileId.value = null;
  Object.assign(form, {
    student_category_id: null,
    date_of_birth: "",
    nationality: "",
    phone: "",
    village_id: null,
    occupation: "",
    academic_stream: null,
    telegram_username: "",
  });
  selectedprovinceid.value = null;
  selecteddistrictid.value = null;
  selectedcommunceid.value = null;
  districts.value = [];
  communces.value = [];
  villages.value = [];
}

async function loadProfile() {
  resetForm();
  if (!props.studentRaw) return;
  loading.value = true;
  try {
    if (!provinces.value.length) fetchProvince();
    if (!studentcategory.value.length) fetchStudentCategory();

    let p = null;
    try {
      const res = await getStudentProfile(props.studentRaw.id);
      p = res.data.data;
    } catch (e) {
      if (e.response?.status !== 404) throw e; // 404 => no profile yet, create mode
    }
    if (!p) return;

    profileId.value = p.id;
    Object.assign(form, {
      student_category_id: p.student_category_id,
      date_of_birth: p.date_of_birth,
      nationality: p.nationality ?? "",
      phone: p.phone ?? "",
      village_id: p.village_id,
      occupation: p.occupation ?? "",
      academic_stream: p.academic_stream || null,
      telegram_username: p.telegram_username ?? "",
    });
    if (p.province_id) {
      selectedprovinceid.value = p.province_id;
      await fetchDistrict(p.province_id);
    }
    if (p.district_id) {
      selecteddistrictid.value = p.district_id;
      await fetchCommunce(p.district_id);
    }
    if (p.commune_id) {
      selectedcommunceid.value = p.commune_id;
      await fetchVillage(p.commune_id);
    }
  } catch {
    ElMessage.error("Failed to load student profile");
  } finally {
    loading.value = false;
  }
}

watch(
  () => props.modelValue,
  (open) => open && loadProfile(),
);

/* ---------- save ---------- */
async function handleSave() {
  await formRef.value.validate();
  saving.value = true;
  try {
    const payload = {
      user_id: props.studentRaw.id,
      student_category_id: form.student_category_id,
      date_of_birth: form.date_of_birth,
      nationality: form.nationality || null,
      phone: form.phone || null,
      village_id: form.village_id,
      occupation: form.occupation || null,
      academic_stream: form.academic_stream || "",
      telegram_username: form.telegram_username || null,
    };
    if (profileId.value) {
      await updateStudentProfile(profileId.value, payload);
    } else {
      await addStudentProfile(payload);
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
    :title="`ពត៍មានលំអិតនិស្សិត${studentRaw ? ' - ' + studentRaw.name_kh : ''}`"
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
          <el-form-item label="ប្រភេទនិស្សិត" prop="student_category_id">
            <el-select
              v-model="form.student_category_id"
              size="large"
              style="width: 100%"
              clearable
            >
              <el-option
                v-for="o in studentcategory"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
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

        <el-col :xs="24" :sm="12">
          <el-form-item label="ខេត្ត/រាជធានី">
            <el-select
              v-model="selectedprovinceid"
              size="large"
              style="width: 100%"
              clearable
              filterable
              @change="provincechange"
            >
              <el-option
                v-for="o in provinces"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ស្រុក/ខណ្ឌ">
            <el-select
              v-model="selecteddistrictid"
              size="large"
              style="width: 100%"
              clearable
              filterable
              :disabled="!selectedprovinceid"
              @change="districtchange"
            >
              <el-option
                v-for="o in districts"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ឃុំ/សង្កាត់">
            <el-select
              v-model="selectedcommunceid"
              size="large"
              style="width: 100%"
              clearable
              filterable
              :disabled="!selecteddistrictid"
              @change="communcechange"
            >
              <el-option
                v-for="o in communces"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ភូមិ" prop="village_id">
            <el-select
              v-model="form.village_id"
              size="large"
              style="width: 100%"
              clearable
              filterable
              :disabled="!selectedcommunceid"
            >
              <el-option
                v-for="o in villages"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :xs="24" :sm="12">
          <AppInput
            label="មុខរបរ"
            prop="occupation"
            clearable
            v-model="form.occupation"
          />
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="ផ្នែកសិក្សា" prop="academic_stream">
            <el-select
              v-model="form.academic_stream"
              size="large"
              style="width: 100%"
              clearable
            >
              <el-option
                v-for="o in streams"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col>
          <AppInput
            label="Telegram"
            prop="telegram_username"
            clearable
            v-model="form.telegram_username"
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
        {{ profileId ? "កែប្រែ" : "បង្កេីត" }}
      </AppButton>
    </template>
  </AppDialog>
</template>
