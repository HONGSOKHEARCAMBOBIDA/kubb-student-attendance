<template>
  <div>
    <AppFilterBar :fields="[
      { slot: 'name', span: 10 },
      { slot: 'role', span: 6 },
      { slot: 'company', span: 5 },
      { slot: 'add', span: 3 }
    ]" :action-span="3">
      <template #name>
        <!-- <el-input v-model="filters.name" placeholder="ស្វែងរក" prefix-icon="Search" clearable @input="fetchUsers" /> -->
        <AppInput v-model="filters.name" placeholder="ស្វែងរក" prefix-icon="Search" clearable @input="fetchUsers">

        </AppInput>
      </template>

      <template #role>
        <el-select v-model="filters.role_id" placeholder="តួនាទី" clearable style="width: 100%" @change="fetchUsers"
          size="large">
          <el-option v-for="role in roles" :key="role.id" :label="role.display_name" :value="role.id" />
        </el-select>
      </template>

      <template #company>
        <el-select v-model="filters.company_id" placeholder="ក្រុមហ៑ុន" clearable style="width: 100%" size="large"
          @change="fetchUsers">
          <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
        </el-select>
      </template>
      <template #add>
        <AppButton v-if="candadd" type="primary" @click="openCreate">
          បន្ថែមបុគ្គលិក
        </AppButton>
      </template>
    </AppFilterBar>
    <el-card>
      <AppTable :data="users" :loading="loading" show-index v-model:current-page="page" v-model:page-size="pageSize"
        :total="total" @page-change="fetchUsers" :columns="[
          { slot: 'name', label: 'ឈ្មោះ', minWidth: 110 },
          { prop: 'gender_string', label: 'ភេទ', minWidth: 90 },
          { prop: 'phone_hash', label: 'លេខទូរសព្ទ', minWidth: 130 },
          { prop: 'role_name', label: 'តួនាទី', minWidth: 110 },
          { prop: 'company_name', label: 'ក្រុមហ៑ុន', minWidth: 130 },
          { label: 'QR', slot: 'qrtoken', width: 150 },
          { slot: 'base_salary', label: 'ប្រាក់ខែ', width: 120 },
          { label: 'ស្ថានភាព', slot: 'status', width: 100 },
        ]" actionsWidth="200">
        <template #name="{ row }">
          <span class="name-cell">
            {{ row.name }}
            <el-tooltip v-if="row.is_verify" content="Verified" placement="top">
              <el-icon class="verified-badge" :size="20">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M12 2L9.5 4.5 6 4l-1 3.5L1.5 9 3 12l-1.5 3L6 16.5l-1 3.5 3.5.5L12 22l2.5-2.5 3.5-.5-1-3.5L20.5 15 19 12l1.5-3-3.5-1.5L16 4l-3.5.5L12 2z" />
                  <path d="M9 12l2 2 4-4" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round"
                    stroke-linejoin="round" />
                </svg>
              </el-icon>
            </el-tooltip>
          </span>
        </template>
        <template #base_salary="{ row }">
          <el-text tag="b">{{ row.base_salary }} {{ row.currency }}</el-text>
        </template>

        <template #qrtoken="{ row }">
          <el-button type="primary" link @click="openQR(row.qr_token)">
            <el-text tag="ins">View QR</el-text>
          </el-button>
        </template>
        <template #status="{ row }">
          <el-text :style="{ color: row.is_active ? 'black' : 'red' }">
            {{ row.is_active ? "កំពុងធ្វើការ" : "ឈប់ធ្វើការ" }}
          </el-text>

        </template>
        <template #actions="{ row }" v-if="canedit">
          <el-tooltip content="កែប្រែ" placement="top">
            <AppButton size="small" icon="Edit" type="warning" circle @click="openEdit(row)">
            </AppButton>
          </el-tooltip>
          <el-tooltip content="បិទ" placement="top">
            <AppButton size="small" :icon="row.is_active ? 'CircleClose' : 'CircleCheck'"
              :type="row.is_active ? 'danger' : 'success'" circle @click="toggleStatus(row)">
            </AppButton>
          </el-tooltip>
          <el-tooltip content="ម៉ោងធ្វេីការ" placement="top">
            <AppButton size="small" icon="View" type="primary" circle @click="openShifts(row)">
            </AppButton>
          </el-tooltip>
          <AppButton v-show="false" size="small" icon="Delete" type="danger" circle @click="DeleteUser(row)">
          </AppButton>
          <el-tooltip content="Verify" placement="top">
            <AppButton size="small" icon="CircleCheck" type="info" circle @click="VerifyUser(row)">

            </AppButton>
          </el-tooltip>
        </template>
      </AppTable>
    </el-card>
    <AppDialog v-model="createDialog" :title="isEdit ? 'កែប្រែបុគ្គលិក' : 'បន្ថែមបុគ្គលិក'" width="983px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-position="top">
        <AppTabs v-model="activeTab" :tabs="[
          { name: 'general', label: 'ព័ត៌មានទូទៅ' },
          { name: 'shifts', label: 'ម៉ោងធ្វើការ' },
        ]" tab-position="top" stretch="true" :lazy="true">
          <template #general>
            <el-row :gutter="16">
              <el-col :xs="24" :sm="12">
                <AppInput label="ឈ្មោះ" prop="name" clearable v-model="createForm.name">
                </AppInput>
              </el-col>
              <el-col :xs="24" :sm="12">
                <AppInput label="លេខទូរសព្ទ" prop="phone_hash" clearable type="number" v-model="createForm.phone_hash"
                  size="large" trim>
                </AppInput>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="តួនាទី" prop="role_id">
                  <el-select v-model="createForm.role_id" placeholder="ជ្រើសតួនាទី" clearable size="large"
                    style="width: 100%">
                    <el-option v-for="role in roles" :key="role.id" :label="role.display_name" :value="role.id" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="ភេទ" prop="gender">
                  <el-select v-model="createForm.gender" size="large" style="width: 100%">
                    <el-option label="ប្រុស" :value="1" />
                    <el-option label="ស្រី" :value="2" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <AppInput label="ប្រាក់ខែ" clearable type="number" v-model="createForm.base_salary" trim>
                </AppInput>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="ក្រុមហ៑ុន" prop="company_id">
                  <el-select v-model="createForm.company_id" placeholder="ជ្រើសក្រុមហ៑ុន" clearable size="large"
                    style="width: 100%">
                    <el-option v-for="company in companys" :key="company.id" :label="company.name"
                      :value="company.id" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="សិទ្ធ">
                  <el-select v-model="createForm.manage_company" placeholder="ជ្រើសការគ្រប់គ្រងក្រុមហ៊ុន" clearable
                    size="large" style="width: 100%">
                    <el-option v-for="managecp in managecompany" :key="managecp.id" :label="managecp.name"
                      :value="managecp.id" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="សាខាដែលអាចមើលទិន្នន័យ">
                  <el-select :disabled="createForm.manage_company != 2" v-model="createForm.company_ids"
                    placeholder="ជ្រើសរើស" multiple clearable size="large" style="width: 100%">
                    <el-option v-for="cp in companys" :key="cp.id" :label="cp.name" :value="cp.id" />
                  </el-select>
                </el-form-item>
              </el-col>

            </el-row>
          </template>

          <template #shifts>
            <div v-for="(shift, i) in createForm.shifts" :key="i" class="shift-row">
              <el-row :gutter="8" align="middle">
                <el-col :xs="12" :sm="4" :gap="5">
                  <el-select v-model="shift.day" placeholder="Day" size="large"
                    style="width: 100%; margin-bottom: 10px">
                    <el-option v-for="(d, idx) in days" :key="idx" :label="d" :value="idx + 1" />
                  </el-select>
                </el-col>
                <el-col :xs="12" :sm="4">
                  <el-select v-model="shift.shift_type" placeholder="Type" size="large" clearable
                    style="width: 100%; margin-bottom: 10px">
                    <el-option label="ធ្វេីការពេញម៉ោង" :value="1" />
                    <el-option label="ធ្វេីការតែមួយព្រឹក" :value="2" />
                    <el-option label="ធ្វេីការតែមួយរសៀល" :value="3" />
                  </el-select>
                </el-col>
                <el-col :xs="12" :sm="3">
                  <el-time-picker v-model="shift.check_in1" placeholder="In 1" size="large" value-format="HH:mm"
                    format="HH:mm" style="width: 100%; margin-bottom: 10px" />
                </el-col>
                <el-col :xs="12" :sm="3">
                  <el-time-picker v-model="shift.check_out1" placeholder="In 1" size="large" value-format="HH:mm"
                    format="HH:mm" style="width: 100%; margin-bottom: 10px" />
                </el-col>
                <el-col :xs="12" :sm="3">
                  <el-time-picker v-model="shift.check_in2" placeholder="In 1" size="large" value-format="HH:mm"
                    format="HH:mm" style="width: 100%; margin-bottom: 10px" />
                </el-col>
                <el-col :xs="12" :sm="3">
                  <el-time-picker v-model="shift.check_out2" placeholder="In 1" size="large" value-format="HH:mm"
                    format="HH:mm" style="width: 100%; margin-bottom: 10px" />
                </el-col>
                <el-col :xs="16" :sm="2">
                  <el-checkbox v-model="shift.is_dayoff" size="large">
                    សម្រាក
                  </el-checkbox>
                </el-col>
                <el-col :xs="8" :sm="2" style="text-align: right">
                  <AppButton size="small" icon="Delete" circle type="danger" @click="createForm.shifts.splice(i, 1)" />
                </el-col>
              </el-row>
            </div>

            <AppButton size="small" icon="Plus" type="default" :block="false" @click="addShift"
              :disabled="createForm.shifts.length >= 7">
              បន្ថែមវេនការងារ
            </AppButton>
          </template>
        </AppTabs>
      </el-form>

      <template #footer>
        <AppButton @click="createDialog = false" size="large" :block="false" type="warning">
          បោះបង់
        </AppButton>
        <AppButton @click="handleCreate" type="primary" :loading="saving" size="large" :block="false">
          {{ isEdit ? "កែប្រែ" : "បង្កេីត" }}
        </AppButton>
      </template>
    </AppDialog>

    <AppDialog v-model="editDialog" title="កែប្រែបុគ្គលិក" width="800px">
      <el-form :model="editForm" :rules="rules" ref="editFormRef" label-position="top">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item label="ឈ្មោះ" prop="name">
              <el-input v-model="editForm.name" placeholder="បញ្ចូលឈ្មោះ" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="លេខទូរសព្ទ" prop="phone_hash">
              <el-input v-model="editForm.phone_hash" placeholder="0xx xxx xxx" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="ភេទ" prop="gender">
              <el-select v-model="editForm.gender" placeholder="ជ្រើសភេទ" style="width: 100%">
                <el-option label="ប្រុស" :value="1" />
                <el-option label="ស្រី" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="ប្រាក់ខែ" prop="base_salary">
              <el-input v-model="editForm.base_salary" :min="0" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="តួនាទី" prop="role_id">
              <el-select v-model="editForm.role_id" placeholder="ជ្រើសតួនាទី" clearable style="width: 100%">
                <el-option v-for="role in roles" :key="role.id" :label="role.display_name" :value="role.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="ក្រុនហ៑ុន" prop="company_id">
              <el-select v-model="editForm.company_id" placeholder="ជ្រើសក្រុនហ៑ុន" clearable style="width: 100%">
                <el-option v-for="company in companys" :key="company.id" :label="company.name" :value="company.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="សិទ្ធ">
              <el-select v-model="editForm.manage_company" placeholder="ជ្រើសការគ្រប់គ្រងក្រុមហ៊ុន" clearable
                size="large" style="width: 100%">
                <el-option v-for="managecp in managecompany" :key="managecp.id" :label="managecp.name"
                  :value="managecp.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="សាខាដែលអាចមើលទិន្នន័យ">
              <el-select :disabled="editForm.manage_company != 2" v-model="editForm.company_ids" placeholder="ជ្រើសរើស"
                multiple clearable size="large" style="width: 100%">
                <el-option v-for="cp in companys" :key="cp.id" :label="cp.name" :value="cp.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <AppButton @click="editDialog = false" type="warning">
            បោះបង់
          </AppButton>
          <AppButton type="primary" :loading="saving" @click="handleUpdate">
            រក្សាទុក
          </AppButton>
        </div>
      </template>
    </AppDialog>

    <AppDialog v-model="shiftsDialog" :title="`វេនធ្វើការ របស់ ${selectedUser?.name}`" width="50%" class="shift-dialog"
      @closed="shiftsEditMode = false">
      <div style="display: flex; justify-content: flex-end; margin-bottom: 12px">
        <AppButton icon="Plus" size="large" type="success" :block="false"
          :disabled="selectedUser?.shift_response?.length >= 7" @click="startShiftCreate">
          បន្ថែមវេន
        </AppButton>
        <AppButton v-if="!shiftsEditMode" icon="Edit" size="large" type="default" :block="false"
          @click="startShiftEdit">
          កែប្រែ
        </AppButton>
      </div>

      <el-table :data="editableShifts" size="large" stripe border empty-text="មិនមានទិន្នន័យវេនការងារ">
        <el-table-column prop="day_name" label="ថ្ងៃ" width="80" fixed />

        <el-table-column label="ប្រភេទ" width="190">
          <template #default="{ row }">
            <el-select v-if="shiftsEditMode" v-model="row.shift_type" size="large" style="width: 100%" clearable>
              <el-option label="ធ្វេីការពេញម៉ោង" :value="1" />
              <el-option label="ធ្វេីការតែមួយព្រឹក" :value="2" />
              <el-option label="ធ្វេីការតែមួយរសៀល" :value="3" />
            </el-select>
            <span v-else>{{ row.shift_type_string }}</span>
          </template>
        </el-table-column>

        <el-table-column label="ចូលវេន១" width="165">
          <template #default="{ row }">
            <el-time-picker v-if="shiftsEditMode" v-model="row.check_in1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px">
            </el-time-picker>
            <span v-else>{{ row.check_in1 || "—" }}</span>
          </template>
        </el-table-column>

        <el-table-column label="ចេញវេន១" width="165">
          <template #default="{ row }">
            <el-time-picker v-if="shiftsEditMode" v-model="row.check_out1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px">
            </el-time-picker>
            <span v-else>{{ row.check_out1 || "—" }}</span>
          </template>
        </el-table-column>

        <el-table-column label="ចូលវេន២" width="165">
          <template #default="{ row }">
            <el-time-picker v-if="shiftsEditMode" v-model="row.check_in2" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px">
            </el-time-picker>
            <span v-else>{{ row.check_in2 || "—" }}</span>
          </template>
        </el-table-column>

        <el-table-column label="ចេញវេន២" width="165">
          <template #default="{ row }">
            <el-time-picker v-if="shiftsEditMode" v-model="row.check_out2" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px">
            </el-time-picker>
            <span v-else>{{ row.check_out2 || "—" }}</span>
          </template>
        </el-table-column>

        <el-table-column label="ថ្ងៃសម្រាក" width="100" align="center">
          <template #default="{ row }">
            <el-checkbox v-if="shiftsEditMode" v-model="row.is_dayoff" />
            <el-tag v-else :type="row.is_dayoff ? 'danger' : 'success'">
              {{ row.is_dayoff ? "សម្រាក" : "ធ្វើការ" }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <template v-if="shiftsEditMode">
          <AppButton @click="cancelShiftEdit" size="large" :block="false" type="warning">
            បោះបង់
          </AppButton>
          <AppButton type="primary" :loading="savingShifts" @click="handleShiftUpdate" size="large" :block="false">
            រក្សាទុក
          </AppButton>
        </template>
      </template>
    </AppDialog>

    <el-dialog v-model="qrDialog" title="ស្កែន QR" width="220px" center :show-close="false">
      <div class="qr-box">
        <img v-if="qrImage" :src="qrImage" />
      </div>
    </el-dialog>

    <AppDialog v-model="shiftsCreateDialog" title="បន្ថែមវេនការងារ" width="1200px">
      <div v-for="(shift, i) in newShifts" :key="i" class="shift-row">
        <el-row :gutter="8" align="middle">
          <el-col :xs="12" :sm="4" :gap="5">
            <el-select v-model="shift.day" placeholder="Day" size="large" style="width: 100%; margin-bottom: 10px">
              <el-option v-for="(d, idx) in days" :key="idx" :label="d" :value="idx + 1" />
            </el-select>
          </el-col>
          <el-col :xs="12" :sm="4">
            <el-select v-model="shift.shift_type" placeholder="Type" size="large" clearable
              style="width: 100%; margin-bottom: 10px">
              <el-option label="ធ្វេីការពេញម៉ោង" :value="1" />
              <el-option label="ធ្វេីការតែមួយព្រឹក" :value="2" />
              <el-option label="ធ្វេីការតែមួយរសៀល" :value="3" />
            </el-select>
          </el-col>
          <el-col :xs="12" :sm="3">

            <el-time-picker v-model="shift.check_in1" placeholder="In 1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px" />
          </el-col>
          <el-col :xs="12" :sm="3">
            <el-time-picker v-model="shift.check_out1" placeholder="In 1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px" />
          </el-col>
          <el-col :xs="12" :sm="3">
            <el-time-picker v-model="shift.check_in2" placeholder="In 1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px" />
          </el-col>
          <el-col :xs="12" :sm="3">
            <el-time-picker v-model="shift.check_out2" placeholder="In 1" size="large" value-format="HH:mm"
              format="HH:mm" style="width: 100%; margin-bottom: 10px" />
          </el-col>
          <el-col :xs="16" :sm="2">
            <el-checkbox v-model="shift.is_dayoff" size="large">
              សម្រាក
            </el-checkbox>
          </el-col>
          <el-col :xs="8" :sm="2" style="text-align: right">
            <AppButton size="small" icon="Delete" circle type="danger" @click="newShifts.splice(i, 1)" />
          </el-col>
        </el-row>
      </div>
      <el-button size="small" icon="Plus" style="margin-top: 8px" @click="addNewShiftRow"
        :disabled="newShifts.length >= 7">
        បន្ថែមជួរ
      </el-button>

      <template #footer>
        <el-button @click="shiftsCreateDialog = false">បោះបង់</el-button>
        <el-button type="primary" :loading="savingShifts" @click="handleShiftCreate">
          រក្សាទុក
        </el-button>
      </template>
    </AppDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  getUsers,
  createUser,
  updateUser,
  toggleUserStatus,
  getrole,
  updateShift,
  createShift,
  getCompany,
  viewcompanyscan,
  deleteuser,
  viewmanagecompany,
  verifyuser
} from "../api/services";
import { watch } from "vue";
import { debounce } from "lodash-es";
import QRCode from "qrcode";
import { ElNotification } from "element-plus";
// import { useAuthStore } from "../stores/auth";
import { useUserDataStore } from '../stores/user_data'
import AppTable from "../../components/AppTable.vue";
import AppButton from "../../components/AppButton.vue";
import AppDialog from "../../components/AppDialog.vue";
import AppFilterBar from "../../components/AppFilterBar.vue";
import AppTabs from "../../components/AppTabs.vue";
import AppInput from "../../components/AppInput.vue";
import { CircleCheckFilled } from '@element-plus/icons-vue'
const shiftsCreateDialog = ref(false);
const newShifts = ref([]);
const qrDialog = ref(false);
const qrImage = ref("");
const users = ref([]);
const roles = ref([]);
const managecompany = ref([]);
const companys = ref([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const filters = reactive({ name: "", role_id: null, company_id: null });
const createDialog = ref(false);
const editDialog = ref(false);
const shiftsDialog = ref(false);
const selectedUser = ref(null);
const editId = ref(null);
const createFormRef = ref();
const editFormRef = ref();
const day = ref(1);
const shiftsEditMode = ref(false);
const savingShifts = ref(false);
const editableShifts = ref([]);
const userDataStore = useUserDataStore()
const canedit = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "edit.user"),
);
const candadd = computed(() =>
  userDataStore.permissions?.some((p) => p.name === "add.user")
)
const activeTab = ref("general");
const days = [
  "ចន្ទ",
  "អង្គារ",
  "ពុធ",
  "ព្រហស្បតិ៍",
  "សុក្រ",
  "សៅរ៍",
  "អាទិត្យ",
];

function startShiftCreate() {
  day.value = 1;
  newShifts.value = [
    // បន្ថែម row ១ ទុកជាដំបូង
    {
      day: day.value,
      shift_type: 1,
      check_in1: "08:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
  ];
  shiftsCreateDialog.value = true;
}

function addNewShiftRow() {
  const lastDay = newShifts.value.at(-1)?.day ?? 0;

  if (lastDay >= 7) {
    ElMessage.warning("មិនអាចបន្ថែមបានទៀតទេ មានតែ ៧ ថ្ងៃក្នុងមួយសប្តាហ៍");
    return;
  }

  newShifts.value.push({
    day: lastDay + 1,
    shift_type: 1,
    check_in1: "08:00",
    check_out1: "11:30",
    check_in2: "13:30",
    check_out2: "17:00",
    is_dayoff: false,
  });
}

async function handleShiftCreate() {
  if (newShifts.value.length === 0) {
    ElMessage.warning("សូមបន្ថែមវេនមួយយ៉ាងហោចណាស់");
    return;
  }

  savingShifts.value = true;
  try {
    const payload = {
      shifts: newShifts.value.map((s) => ({
        user_id: selectedUser.value.id,
        day: s.day,
        shift_type: s.shift_type,
        check_in1: s.check_in1 || "",
        check_out1: s.check_out1 || "",
        check_in2: s.check_in2 || "",
        check_out2: s.check_out2 || "",
        is_dayoff: s.is_dayoff,
      })),
    };
    await createShift(payload);
    //  selectedUser.value.shift_response = JSON.parse(JSON.stringify(newShifts.value))
    ElMessage.success("បន្ថែមវេនការងារបានជោគជ័យ");
    await fetchUsers();
    shiftsCreateDialog.value = false;
    shiftsDialog.value = false;
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to create shifts");
  } finally {
    savingShifts.value = false;
  }
}

const openQR = async (token) => {
  qrImage.value = await QRCode.toDataURL(token);
  qrDialog.value = true;
};

const debouncedFetch = debounce(() => {
  page.value = 1;
  fetchUsers();
}, 500);

watch(() => filters.name, debouncedFetch);
const createForm = reactive({
  company_ids: [],
  manage_company: null,
  company_id: null,
  name: "",
  phone_hash: "",
  role_id: null,
  gender: null,
  base_salary: "",
  shifts: [],
});
const createRules = {
  name: [{ required: true, message: "Name is required" }],
  phone_hash: [{ required: true, message: "Phone is required" }],
};

const rules = {
  name: [{ required: true, message: "សូមបញ្ចូលឈ្មោះ", trigger: "blur" }],
  phone_hash: [
    { required: true, message: "សូមបញ្ចូលលេខទូរសព្ទ", trigger: "blur" },
  ],
  gender: [{ required: true, message: "សូមជ្រើសភេទ", trigger: "change" }],
  base_salary: [
    { required: true, message: "សូមបញ្ចូលប្រាក់ខែ", trigger: "blur" },
  ],
  role_id: [{ required: true, message: "សូមជ្រើសតួនាទី", trigger: "change" }],
};

const editForm = reactive({
  name: "",
  phone_hash: "",
  gender: 1,
  base_salary: "",
  role_id: null,
  company_id: null,
  manage_company: null,
  company_ids: []
});

function addShift() {
  createForm.shifts.push({
    day: 1,
    shift_type: 1,
    check_in1: "",
    check_out1: "",
    check_in2: "",
    check_out2: "",
    is_dayoff: false,
  });
}

async function fetchUsers() {
  loading.value = true;
  try {
    const params = { page: page.value, page_size: pageSize.value };
    if (filters.name) params.name = filters.name;
    if (filters.role_id) params.role_id = filters.role_id;
    if (filters.company_id) params.company_id = filters.company_id;
    const res = await getUsers(params);
    users.value = res.data.data || [];
    total.value = res.data.pagination?.totalCount || 0;
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

async function fetchRole() {
  loading.value = true;
  try {
    const res = await getrole();
    roles.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

async function fetchmanagecompany() {
  loading.value = true;
  try {
    const res = await viewmanagecompany()
    managecompany.value = res.data.data || [];
  } catch (e) {
    ElMessage.error("Failed to load managecompany");
  } finally {
    loading.value = false;
  }
}

async function fetchCompany() {
  loading.value = true;
  try {
    const res = await viewcompanyscan();
    companys.value = res.data.data || [];
  } catch {
    ElMessage.error("Failed to load employees");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  createForm.shifts = [
    {
      day: 1,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 2,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 3,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 4,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 5,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 6,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
    {
      day: 7,
      shift_type: 1,
      check_in1: "07:00",
      check_out1: "11:30",
      check_in2: "13:30",
      check_out2: "17:00",
      is_dayoff: false,
    },
  ];
  Object.assign(createForm, {
    name: "",
    phone_hash: "",
    role_id: null,
    gender: null,
    base_salary: "",
    company_id: null,
  });
  createDialog.value = true;
}

async function handleCreate() {
  await createFormRef.value.validate();
  saving.value = true;
  try {
    const payload = {
      name: createForm.name,
      phone_hash: createForm.phone_hash,
      role_id: createForm.role_id,
      gender: createForm.gender,
      base_salary: createForm.base_salary,
      company_id: createForm.company_id,
      manage_company: createForm.manage_company,
      company_ids: createForm.company_ids,
      day: createForm.shifts.map((s) => s.day),
      shift_type: createForm.shifts.map((s) => s.shift_type),
      check_in1: createForm.shifts.map((s) => s.check_in1 || null),
      check_out1: createForm.shifts.map((s) => s.check_out1 || null),
      check_in2: createForm.shifts.map((s) => s.check_in2 || null),
      check_out2: createForm.shifts.map((s) => s.check_out2 || null),
      is_dayoff: createForm.shifts.map((s) => s.is_dayoff),
    };
    await createUser(payload);
    ElMessage.success("បង្កេីតបុគ្គលិកបានជោគជ័យ");
    createDialog.value = false;
    fetchUsers();
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to create");
  } finally {
    saving.value = false;
  }
}

function openEdit(row) {
  editId.value = row.id;
  Object.assign(editForm, {
    name: row.name,
    phone_hash: row.phone_hash,
    gender: row.gender,
    base_salary: row.base_salary,
    role_id: row.role_id,
    company_id: row.company_id,
    manage_company: row.manage_company,
    company_ids: row.company_ids
  });
  editDialog.value = true;
}

async function handleUpdate() {
  saving.value = true;
  try {
    await updateUser(editId.value, editForm);
    ElMessage.success("កែប្រែបុគ្គលិកបានជោគជ័យ");
    editDialog.value = false;
    fetchUsers();
  } catch (e) {
    ElMessage.error(
      e.response?.data?.message ||
      e.response?.data?.error ||
      e.message ||
      "Failed to update",
    );
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row) {
  await ElMessageBox.confirm(
    `${row.is_active ? "បិទ" : "បេីក"} ${row.name}?`,
    "សូមបញ្ជាក់",
    { type: "warning" },
  );
  try {
    await toggleUserStatus(row.id);
    ElMessage.success("កែប្រែស្ថានភាពបានជោគជ័យ");
    fetchUsers();
  } catch (e) {
    ElNotification.error({
      title: "Error",
      message: e.response?.data?.error,
      offset: 100,
    });
  }
}

async function VerifyUser(row) {
  try {
    await verifyuser(row.id);
    ElNotification.success({
      title: "Success",
      message: "Verify Success",
      offset: 10,
    });
    fetchUsers();
  } catch (e) {
    ElNotification.error({
      title: "Error",
      message: e.response?.data?.error,
      offset: 100,
    });
  }
}

async function DeleteUser(row) {
  try {
    await ElMessageBox.confirm(
      `តើអ្នកពិតជាចង់លុប ${row.name} មែនទេ? សកម្មភាពនេះមិនអាចត្រឡប់វិញបានទេ។`,
      "សូមបញ្ជាក់",
      { type: "warning" },
    );
  } catch {
    // user cancelled — do nothing
    return;
  }

  try {
    await deleteuser(row.id);
    ElMessage.success("លុបបុគ្គលិកបានជោគជ័យ");
    fetchUsers();
  } catch (e) {
    ElNotification.error({
      title: "Error",
      message: e.response?.data?.error,
      offset: 100,
    });
  }
}

function openShifts(row) {
  selectedUser.value = row;
  editableShifts.value = JSON.parse(JSON.stringify(row.shift_response || []));
  shiftsEditMode.value = false;
  shiftsDialog.value = true;
}

function startShiftEdit() {
  editableShifts.value = JSON.parse(
    JSON.stringify(selectedUser.value?.shift_response || []),
  );
  shiftsEditMode.value = true;
}

function cancelShiftEdit() {
  editableShifts.value = JSON.parse(
    JSON.stringify(selectedUser.value?.shift_response || []),
  );
  shiftsEditMode.value = false;
}

async function handleShiftUpdate() {
  savingShifts.value = true;
  try {
    const payload = {
      shifts: editableShifts.value.map((s) => ({
        id: s.id,
        check_in1: s.check_in1 || null,
        check_out1: s.check_out1 || null,
        check_in2: s.check_in2 || null,
        check_out2: s.check_out2 || null,
        shift_type: s.shift_type,
        is_dayoff: s.is_dayoff,
      })),
    };
    await updateShift(payload);
    selectedUser.value.shift_response = JSON.parse(
      JSON.stringify(editableShifts.value),
    );
    ElMessage.success("កែប្រែវេនការងារបានជោគជ័យ");
    shiftsEditMode.value = false;
  } catch (e) {
    ElMessage.error(e.response?.data?.message || "Failed to update shifts");
  } finally {
    savingShifts.value = false;
  }
}

onMounted(() => {
  fetchCompany();
  fetchUsers();
  fetchRole();
  fetchmanagecompany()
});
</script>

<style scoped>
.filter-card {
  border-radius: 6px;
  padding: 4px 0;
}

.table-card {
  border-radius: 6px;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.shift-row {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 8px;
}



.shift-dialog .el-dialog__header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.shift-dialog .el-dialog__body {
  padding: 24px;
}

.dialog-content {
  margin-top: 4px;
}

.qr-box {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 6px;
  background: #ffffff;
}

.qr-box img {
  width: 170px;
  height: 170px;
}

.name-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.verified-badge {
  color: #1d9bf0;
  display: inline-flex;
}

@media (max-width: 768px) {
  .shift-dialog {
    width: 95% !important;
  }

  .shift-dialog .el-dialog__body {
    padding: 16px;
  }
}
</style>
