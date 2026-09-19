<template>
  <div>
<AppFilterBar
  :fields="[
    { slot: 'name', span: 4 },
    { slot: 'major', span: 4 },
    { slot: 'shift', span: 4 },
    { slot: 'generation', span: 4 },
    { slot: 'programm', span: 4 },
    { slot: 'create', span: 4 },
  ]"
>
  <template #name>
    <el-input
      v-model.trim="filters.name"
      placeholder="ស្វែងរកតាមឈ្មោះថ្នាក់"
      clearable
      size="large"
    />
  </template>

  <template #major>
    <el-select
      v-model="filters.major_id"
      placeholder="ជំនាញ"
      clearable
      filterable
      size="large"
    >
      <el-option
        v-for="item in majors"
        :key="item.id"
        :label="item.name_kh"
        :value="item.id"
      />
    </el-select>
  </template>

  <template #shift>
    <el-select
      v-model="filters.shift_id"
      placeholder="វេន"
      clearable
      size="large"
    >
      <el-option
        v-for="item in shifts"
        :key="item.id"
        :label="item.name"
        :value="item.id"
      />
    </el-select>
  </template>

  <template #generation>
    <el-select
      v-model="filters.generation_id"
      placeholder="ជំនាន់"
      clearable
      filterable
      size="large"
    >
      <el-option
        v-for="item in generations"
        :key="item.id"
        :label="item.name_kh"
        :value="item.id"
      />
    </el-select>
  </template>

  <template #programm>
    <el-select
      v-model="filters.programme_id"
      placeholder="កម្មវិធីសិក្សា"
      clearable
      filterable
      size="large"
    >
      <el-option
        v-for="item in programmes"
        :key="item.id"
        :label="item.name"
        :value="item.id"
      />
    </el-select>
  </template>

  <template #create>
    <AppButton
      v-if="canAddClass"
      type="primary"
      @click="openCreate"
      :block="false"
    >
      បន្ថែមថ្នាក់
    </AppButton>
  </template>
</AppFilterBar>



    <el-card class="table-card">
      <AppTable
        expandable
        :data="classes"
        :loading="loading"
        show-index
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="fetchClasses"
        actions-width="300px"
        :columns="[
          { prop: 'name', label: 'ឈ្មោះថ្នាក់', minWidth: 120 },
          { prop: 'major_name', label: 'ជំនាញ', minWidth: 110 },
          { prop: 'shift_name', label: 'វេន', width: 100 },
          { prop: 'generation_name', label: 'ជំនាន់', width: 110 },
          { prop: 'programme_name', label: 'កម្មវិធីសិក្សា', minWidth: 130 },
          { prop: 'year', label: 'ឆ្នាំ', width: 90 },
          { prop: 'semester', label: 'ឆមាស', width: 90 },
          { prop: 'group', label: 'ក្រុម', width: 90 },
          { prop: 'term', label: 'ឆមាស/Term', width: 100 },
          { prop: 'radius', label: 'ចម្ងាយអាចស្កែន (m)', width: 150 },
          { label: 'អាចស្កែនក្រៅតំបន់', slot: 'outsize', width: 150 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
          { label: 'សិស្សសរុប', slot: 'total', width: 100 },
        ]"
      >
        <template #outsize="{ row }">
          <el-tag :type="row.can_scan_outsize ? 'success' : 'danger'">
            {{ row.can_scan_outsize ? "បាន" : "មិនបាន" }}
          </el-tag>
        </template>

        <template #status="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
            {{ row.is_active ? "Active" : "Inactive" }}
          </el-tag>
        </template>
        <template #total="{row}">
          <el-text>{{ row.students?.length ?? 0 }} នាក់</el-text>
        </template>

        <template #actions="{ row }">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
            :disabled="row.is_active === false"
            v-if="canEditClass"
            size="small"
            icon="Edit"
            type="warning"
            circle
            @click="openEdit(row)"
          >
          </AppButton>
          </el-tooltip>
          <AppButton
            :disabled="row.is_active === false"
            v-if="canEditClass"
            size="small"
            icon="Promotion"
            type="primary"
            circle
            @click="openEditTelegram(row)"
          >
          </AppButton>
          <el-tooltip content="បិទ/បើក ថ្នាក់" placement="top">
            <AppButton
            v-if="canEditClass"
            size="small"
            icon="Switch"
            :type="row.is_active ? 'danger' : 'success'"
            circle
            @click="handleToggleStatus(row)"
          >
          </AppButton>
          </el-tooltip>
<el-tooltip content="បញ្ជូលសិស្សតាមExcell" placement="top">
  <AppButton
    :disabled="row.is_active === false"
    size="small"
    icon="Download"
    type="success"
    circle
    @click="openImport(row)"
  >
  </AppButton>
</el-tooltip>
<el-tooltip content="Copy ទិន្ន័យសិស្សទៅថ្នាក់ផ្សេង" placement="top">
  <AppButton
  :disabled="row.is_active === false"
    size="small"
    icon="CopyDocument"
    type="primary"
    circle
    @click="openCopy(row)"
  >
  </AppButton>
</el-tooltip>
<el-tooltip content="កាលវិភាគ" placement="top">
  <AppButton
  :disabled="row.is_active === false"
    size="small"
    icon="Calendar"
    type="primary"
    circle
    @click="openSchedule(row)"
  >
  </AppButton>
</el-tooltip>
        </template>

          <template #expand="{ row: students }">
          <el-divider content-position="left">
            <el-text> សិស្ស </el-text>
          </el-divider>
          <AppTable
            show-index
            :data="students.students"
            :columns="studentcolumn"
            :show-pagination="false"
          >
            <template #gender="{ row: students }">
              <el-text>{{ students.gender === '1' ? 'ប្រុស' : 'ស្រី' }}</el-text>
            </template>
            <template #actions>
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton
            v-if="canEditClass"
            size="small"
            icon="Edit"
            type="warning"
            circle
          
          >
          </AppButton>
          </el-tooltip>
            </template>
          </AppTable>
        </template>
      </AppTable>
    </el-card>

    <AppDialog
      v-model="dialogVisible"
      :title="isEdit ? 'កែប្រែថ្នាក់' : 'បន្ថែមថ្នាក់'"
      width="640px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="ឈ្មោះថ្នាក់" prop="name">
          <el-input
            v-model.trim="form.name"
            size="large"
            placeholder="បញ្ចូលឈ្មោះថ្នាក់"
          />
        </el-form-item>

        <div class="form-row">
          <el-form-item label="ជំនាញ" prop="major_id">
            <el-select v-model="form.major_id" size="large" placeholder="ជ្រើសរើសជំនាញ">
              <el-option
                v-for="item in majors"
                :key="item.id"
                :label="item.name_kh"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="វេន" prop="shift_id">
            <el-select v-model="form.shift_id" size="large" placeholder="ជ្រើសរើសវេន">
              <el-option
                v-for="item in shifts"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ជំនាន់" prop="generation_id">
            <el-select v-model="form.generation_id" size="large" placeholder="ជ្រើសរើសជំនាន់">
              <el-option
                v-for="item in generations"
                :key="item.id"
                :label="item.name_kh"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="កម្មវិធីសិក្សា" prop="programme_id">
            <el-select v-model="form.programme_id" size="large" placeholder="ជ្រើសរើសកម្មវិធីសិក្សា">
              <el-option
                v-for="item in programmes"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ឆ្នាំ" prop="year">
            <el-input type="number" v-model.number="form.year" size="large" placeholder="e.g. 2026" />
          </el-form-item>
          <el-form-item label="ឆមាស" prop="semester">
            <el-input type="number" v-model.number="form.semester" size="large" placeholder="e.g. 1" />
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="ក្រុម" prop="group">
            <el-input type="number" v-model.number="form.group" size="large" placeholder="e.g. A" />
          </el-form-item>
          <el-form-item label="Term" prop="term">
            <el-input type="number" v-model.number="form.term" size="large" placeholder="e.g. 1" />
          </el-form-item>
        </div>

        <template v-if="!isEdit">
          <el-form-item label="Map Link" prop="map_link">
            <el-input
              v-model.trim="form.map_link"
              size="large"
              placeholder="https://maps.google.com/..."
            />
          </el-form-item>

          <div class="form-row">
            <el-form-item label="ចម្ងាយអាចស្កែនបាន (ម៉េត្រ)" prop="radius">
              <el-input
                v-model.trim="form.radius"
                placeholder="e.g. 100"
                size="large"
              />
            </el-form-item>
            <el-form-item label="អាចស្កែនក្រៅតំបន់" prop="can_scan_outsize">
              <el-radio-group v-model="form.can_scan_outsize" size="large">
                <el-radio-button :value="true">អាចស្កែនបាន</el-radio-button>
                <el-radio-button :value="false">មិនអាចស្កែនបាន</el-radio-button>
              </el-radio-group>
            </el-form-item>
          </div>
        </template>
      </el-form>

      <template #footer>
        <AppButton @click="dialogVisible = false" size="large" :block="false" type="warning">
          បោះបង់
        </AppButton>
        <AppButton
          @click="handleSave"
          type="primary"
          :loading="saving"
          size="large"
          :block="false"
        >
          {{ isEdit ? "កែប្រែ" : "បង្កើត" }}
        </AppButton>
      </template>
    </AppDialog>

    <AppDialog v-model="dialogTelegramVisible" title="កែប្រែ Telegram" width="600px">
      <el-form :model="telegramForm" ref="telegramFormRef" label-position="top">
        <el-divider />
        <p class="section-label">Telegram</p>
        <div class="form-row">
          <el-form-item label="Bot Token" prop="bot_token">
            <el-input
              v-model="telegramForm.bot_token"
              placeholder="Telegram bot token"
              size="large"
            />
          </el-form-item>
          <el-form-item label="Group Link" prop="group_link">
            <el-input
              v-model="telegramForm.group_link"
              placeholder="Telegram group chat ID"
              size="large"
            />
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <AppButton @click="dialogTelegramVisible = false" size="large" :block="false">
          ថតក្រោយ
        </AppButton>
        <AppButton
          type="primary"
          :loading="saving"
          @click="handleUpdateTelegram"
          size="large"
          :block="false"
        >
          កែប្រែ
        </AppButton>
      </template>
    </AppDialog>

    <AppDialog
  v-model="importDialog"
  :title="`បញ្ចូលសិស្សតាម Excel — ${importClassName}`"
  width="900px"
  @closed="resetImport"
>
  <el-upload
    drag
    :auto-upload="false"
    :show-file-list="false"
    accept=".xlsx,.xls"
    :on-change="handleFilePicked"
  >
    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
    <div class="el-upload__text">
      អូសឯកសារមកទីនេះ ឬ <em>ចុចដើម្បីជ្រើសរើសឯកសារ</em>
    </div>
    <template #tip>
      <div class="el-upload__tip">
        Column ដែលត្រូវការ (ជួរទី ១): <b>name_kh, name_en, gender, code</b> —
        gender អាចជា 1/2 ឬ ប្រុស/ស្រី
      </div>
    </template>
  </el-upload>

  <el-table
    v-if="previewRows.length"
    :data="previewRows"
    size="small"
    stripe
    border
    style="margin-top: 16px"
    max-height="360"
  >
    <el-table-column type="index" width="50" label="#" />
    <el-table-column prop="name_kh" label="ឈ្មោះខ្មែរ" />
    <el-table-column prop="name_en" label="ឈ្មោះឡាតាំង" />
    <el-table-column prop="code" label="កូដ" />
    <el-table-column label="ភេទ" width="90">
      <template #default="{ row }">{{ row.gender === 1 ? "ប្រុស" : row.gender === 2 ? "ស្រី" : "?" }}</template>
    </el-table-column>
    <el-table-column label="ស្ថានភាព" width="110">
      <template #default="{ row }">
        <el-tag :type="row._valid ? 'success' : 'danger'" size="small">
          {{ row._valid ? "OK" : "ខ្វះទិន្នន័យ" }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>

  <template #footer>
    <AppButton @click="importDialog = false" size="large" :block="false" type="warning">
      បោះបង់
    </AppButton>
    <AppButton
      @click="handleImportSubmit"
      type="primary"
      :loading="importing"
      size="large"
      :block="false"
      :disabled="!previewRows.length"
    >
      នាំចូល ({{ validRowCount }})
    </AppButton>
  </template>
</AppDialog>

<AppDialog
  v-model="copyDialog"
  :title="`Copy សិស្សពី ${copySourceClass?.name || ''}`"
  width="700px"
>
  <el-form label-position="top">
    <el-form-item label="ថ្នាក់គោលដៅ" required>
      <el-select
        v-model="copyTargetClassId"
        placeholder="ជ្រើសរើសថ្នាក់"
        filterable
        size="large"
        style="width: 100%"
      >
        <el-option
          v-for="c in copyTargetOptions"
          :key="c.id"
          :label="c.name"
          :value="c.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="ជ្រើសរើសសិស្ស">
      <el-checkbox
        v-model="copySelectAll"
        @change="handleCopySelectAll"
        style="margin-bottom: 8px"
      >
        ជ្រើសរើសទាំងអស់
      </el-checkbox>
      <el-table
        :data="copyStudents"
        size="small"
        border
        max-height="320"
      >
        <el-table-column width="50">
          <template #default="{ row }">
            <el-checkbox v-model="row._checked" />
          </template>
        </el-table-column>
        <el-table-column prop="name_kh" label="ឈ្មោះខ្មែរ" />
        <el-table-column prop="name_en" label="ឈ្មោះឡាតាំង" />
        <el-table-column prop="code" label="អត្តលេខ" width="100" />
      </el-table>
    </el-form-item>
  </el-form>

  <template #footer>
    <AppButton @click="copyDialog = false" size="large" :block="false" type="warning">
      បោះបង់
    </AppButton>
    <AppButton
      @click="handleCopySubmit"
      type="primary"
      :loading="copying"
      size="large"
      :block="false"
      :disabled="!copyTargetClassId || !copySelectedCount"
    >
      Copy ({{ copySelectedCount }})
    </AppButton>
  </template>
</AppDialog>
<ClassScheduleDialog v-model="scheduleDialog" :class-row="scheduleClass" />
  </div>
</template>

<script setup>
import * as XLSX from "xlsx";
import { registerUsersExcel } from "../api/services"; // add to your existing services import
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useUserDataStore } from "../stores/user_data";
import {
  getClass,
  createClass,
  updateClass,
  changeStatusClass,
  updateClassTelegram,
  getMajor,
  getShift,
  getGeneration,
  getProgramme,
  adduserclass
} from "../api/services";
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import { useNotification } from "../../composables/useNotification.js";
import { useLoading } from "../../composables/useLoading.js";
import AppFilterBar from "../../components/AppFilterBar.vue";
import ClassScheduleDialog from "../../components/ClassScheduleDialog.vue"; // adjust path
const notify = useNotification();
const userDataStore = useUserDataStore();
const useloading = useLoading();

const scheduleDialog = ref(false);
const scheduleClass = ref(null);
function openSchedule(row) {
  scheduleClass.value = row;
  scheduleDialog.value = true;
}

const classes = ref([]);
const majors = ref([]);
const shifts = ref([]);
const generations = ref([]);
const programmes = ref([]);

const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogVisible = ref(false);
const dialogTelegramVisible = ref(false);
const isEdit = ref(false);
const editId = ref(null);
const formRef = ref();
const telegramFormRef = ref();

const copyDialog = ref(false);
const copying = ref(false);
const copySourceClass = ref(null);
const copyTargetClassId = ref(null);
const copyStudents = ref([]);
const copySelectAll = ref(true);

const copyTargetOptions = computed(() =>
  classes.value.filter((c) => c.id !== copySourceClass.value?.id),
);

const copySelectedCount = computed(
  () => copyStudents.value.filter((s) => s._checked).length,
);

function openCopy(row) {
  copySourceClass.value = row;
  copyStudents.value = (row.students || []).map((s) => ({
    ...s,
    _checked: true,
  }));
  copyTargetClassId.value = null;
  copySelectAll.value = true;
  copyDialog.value = true;
}

function handleCopySelectAll(val) {
  copyStudents.value.forEach((s) => (s._checked = val));
}

async function handleCopySubmit() {
  const selected = copyStudents.value.filter((s) => s._checked);
  if (!selected.length || !copyTargetClassId.value) return;

  copying.value = true;
  try {
    await adduserclass({
      class_id: copyTargetClassId.value,
      user_id: selected.map((s) => ({ user_id: s.id })),
    });
    notify.success("Copy ទិន្ន័យសិស្សបានជោគជ័យ");
    copyDialog.value = false;
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "Copy បរាជ័យ");
  } finally {
    copying.value = false;
  }
}

const filters = reactive({
  name: "",
  major_id: "",
  shift_id: "",
  generation_id: "",
  programme_id: "",
});

let debounceTimer = null;
function debounce(fn, delay = 400) {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fn, delay);
}

const studentcolumn = [
  { prop: "name_kh", label: "ឈ្មោះខ្មែរ", minwidth: 100 },
  { prop: "name_en", label: "ឈ្មោះអង់គ្លេស", minwidth: 100 },
   { slot: "gender", label: "ភេទ", minwidth: 100 },
  { prop: "code", label: "អត្តលេខ", minwidth: 100 },
];

const form = reactive({
  name: "",
  major_id: null,
  shift_id: null,
  generation_id: null,
  programme_id: null,
  year: null,
  semester: null,
  group: null,
  term: null,
  map_link: "",
  radius: "",
  can_scan_outsize: false,
});

const telegramForm = reactive({
  bot_token: "",
  group_link: "",
});

const canAddClass = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.company"),
);

const canEditClass = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.company"),
);

const rules = {
  name: [{ required: true, message: "Class name is required" }],
  major_id: [{ required: true, message: "Major is required" }],
  shift_id: [{ required: true, message: "Shift is required" }],
  generation_id: [{ required: true, message: "Generation is required" }],
  programme_id: [{ required: true, message: "Programme is required" }],
  year: [{ required: true, message: "Year is required" }],
  semester: [{ required: true, message: "Semester is required" }],
  map_link: [{ required: true, message: "Map link is required" }],
  radius: [{ required: true, message: "Radius is required" }],
  can_scan_outsize: [
    { required: true, message: "Employee can scan outside or not" },
  ],
};

async function fetchClasses() {
  loading.value = true;
  try {
    const res = await getClass({
      page: page.value,
      page_size: pageSize.value,
      name: filters.name || undefined,
      major_id: filters.major_id || undefined,
      shift_id: filters.shift_id || undefined,
      generation_id: filters.generation_id || undefined,
      programme_id: filters.programme_id || undefined,
    });
    classes.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
    console.log(classes.value)
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    loading.value = false;
  }
}

async function fetchLookups() {
  try {
    const [majorRes, shiftRes, generationRes, programmeRes] = await Promise.all([
      getMajor(),
      getShift(),
      getGeneration(),
      getProgramme(),
    ]);
    majors.value = majorRes.data.data || [];
    shifts.value = shiftRes.data.data || [];
    generations.value = generationRes.data.data || [];
    programmes.value = programmeRes.data.data || [];
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

function resetForm() {
  form.name = "";
  form.major_id = null;
  form.shift_id = null;
  form.generation_id = null;
  form.programme_id = null;
  form.year = "";
  form.semester = "";
  form.group = "";
  form.term = "";
  form.map_link = "";
  form.radius = "";
  form.can_scan_outsize = null;
}

function openCreate() {
  isEdit.value = false;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, {
    name: row.name || "",
    major_id: row.major_id ?? null,
    shift_id: row.shift_id ?? null,
    generation_id: row.generation_id ?? null,
    programme_id: row.programme_id ?? null,
    year: row.year || "",
    semester: row.semester || "",
    group: row.group || "",
    term: row.term || "",
    map_link: "",
    radius: row.radius || "",
    can_scan_outsize: row.can_scan_outsize ?? null,
  });
  dialogVisible.value = true;
}

function openEditTelegram(row) {
  editId.value = row.id;
  telegramForm.bot_token = row.bot_token || "";
  telegramForm.group_link = row.group_chatID || row.group_link || "";
  dialogTelegramVisible.value = true;
}

async function handleSave() {
  await formRef.value.validate();
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    if (isEdit.value) {
      const payload = {
        name: form.name,
        major_id: form.major_id,
        shift_id: form.shift_id,
        generation_id: form.generation_id,
        programme_id: form.programme_id,
        year: form.year,
        semester: form.semester,
        group: form.group,
        term: form.term,
      };
      await updateClass(editId.value, payload);
      notify.success("កែប្រែបានជោគជ័យ");
    } else {
      await createClass(form);
      notify.success("បង្កើតថ្នាក់បានជោគជ័យ");
    }
    dialogVisible.value = false;
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    useloading.hide();
  }
}

async function handleUpdateTelegram() {
  useloading.show({ text: "កំពុងដំណេីរការ..." });
  try {
    await updateClassTelegram(editId.value, {
      bot_token: telegramForm.bot_token,
      group_link: telegramForm.group_link,
    });
    notify.success("កែប្រែបានជោគជ័យ");
    dialogTelegramVisible.value = false;
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  } finally {
    useloading.hide();
  }
}

async function handleToggleStatus(row) {
  try {
    await changeStatusClass(row.id);
    notify.success("ផ្លាស់ប្តូរស្ថានភាពបានជោគជ័យ");
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "");
  }
}

/* ---------------- Excel import (scoped to a class) ---------------- */

const importDialog = ref(false);
const importing = ref(false);
const previewRows = ref([]);
const pickedFile = ref(null);
const importClassId = ref(null);
const importClassName = ref("");

const validRowCount = computed(
  () => previewRows.value.filter((r) => r._valid).length,
);

function openImport(row) {
  importClassId.value = row.id;
  importClassName.value = row.name;
  resetImport();
  importDialog.value = true;
}

function resetImport() {
  previewRows.value = [];
  pickedFile.value = null;
}

function normalizeGender(val) {
  if (val === 1 || val === 2) return val;
  const s = String(val || "").trim();
  if (s === "1" || s === "ប្រុស" || /^m(ale)?$/i.test(s)) return 1;
  if (s === "2" || s === "ស្រី" || /^f(emale)?$/i.test(s)) return 2;
  return null;
}

// Client-side parse is only for an instant preview so the user can catch
// mistakes before uploading. The authoritative parse happens server-side.
function handleFilePicked(uploadFile) {
  const file = uploadFile.raw;
  pickedFile.value = file;

  const reader = new FileReader();
  reader.onload = (e) => {
    const wb = XLSX.read(e.target.result, { type: "array" });
    const sheet = wb.Sheets[wb.SheetNames[0]];
    const rows = XLSX.utils.sheet_to_json(sheet, { defval: "" });

    previewRows.value = rows.map((r) => {
      const name_kh = String(r.name_kh ?? r["ឈ្មោះខ្មែរ"] ?? "").trim();
      const name_en = String(r.name_en ?? r["ឈ្មោះឡាតាំង"] ?? "").trim();
      const code = String(r.code ?? r["កូដ"] ?? "").trim();
      const gender = normalizeGender(r.gender ?? r["ភេទ"]);
      return {
        name_kh,
        name_en,
        code,
        gender,
        _valid: !!(name_kh && name_en && code && gender),
      };
    });

    if (!previewRows.value.length) {
      notify.error("រកមិនឃើញទិន្នន័យក្នុងឯកសារនេះទេ");
    }
  };
  reader.readAsArrayBuffer(file);
}

async function handleImportSubmit() {
  if (!pickedFile.value) return;
  if (validRowCount.value === 0) {
    notify.error("គ្មានជួរដេលត្រឹមត្រូវសម្រាប់នាំចូលទេ");
    return;
  }

  importing.value = true;
  try {
    const formData = new FormData();
    formData.append("file", pickedFile.value);
    formData.append("class_id", importClassId.value); // scope import to this class
    const res = await registerUsersExcel(formData);
    const created = res.data?.created ?? validRowCount.value;
    notify.success(`នាំចូលសិស្សបានជោគជ័យ (${created})`);
    importDialog.value = false;
    fetchClasses();
  } catch (e) {
    notify.error(e.response?.data?.error || "Failed to import file");
  } finally {
    importing.value = false;
  }
}

watch(
  () => filters.name,
  () => {
    page.value = 1;
    debounce(() => fetchClasses());
  }
);

watch(
  [
    () => filters.major_id,
    () => filters.shift_id,
    () => filters.generation_id,
    () => filters.programme_id,
  ],
  () => {
    page.value = 1;
    fetchClasses();
  }
);

onMounted(() => {
  fetchClasses();
  fetchLookups();
});
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.table-card {
  border-radius: 6px;
}

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