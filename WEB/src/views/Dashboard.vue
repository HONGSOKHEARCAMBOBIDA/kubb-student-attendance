
<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserDataStore } from '../stores/user_data'
import { getcountleave } from '../api/services'

const router = useRouter()
const userDataStore = useUserDataStore()

// User level
const adminLevel = computed(() => userDataStore.level >= 2)
const staffLevel = computed(() => userDataStore.level === 2)
const managerLevel = computed(() => userDataStore.level === 3)

// Permission
const viewIncome = computed(() =>
  userDataStore.permissions?.some(p => p.name === 'view.income')
)

const addAttendance = computed(() =>
  userDataStore.permissions?.some(p => p.name === 'add.attendance')
)

// Leave count
const leaves = ref([])

const leaveCount = computed(() =>
  Array.isArray(leaves.value)
    ? leaves.value.length
    : Number(leaves.value || 0)
)

async function fetchCountLeave() {
  try {
    const res = await getcountleave()
    leaves.value = res.data.data || []
  } catch (error) {
    console.error(error)
  }
}

// Menu items
const actions = computed(() => [
  {
    label: 'មើលវត្តមានសង្ខេប',
    description: 'ពិនិត្យវត្តមានសង្ខេប',
    path: '/attendance',
    icon: 'Calendar',
    show: true
  },
  {
    label: 'របាយការណ៍វត្តមាន',
    description: 'របាយការណ៍វត្តមានសរុប',
    path: '/AttendanceReport',
    icon: 'DataAnalysis',
    show: true
  },
  {
    label: 'សុំច្បាប់',
    description: 'គ្រប់គ្រងការសុំច្បាប់',
    path: '/leaverequest',
    icon: 'DocumentChecked',
    show: true
  },
  {
    label: 'ឈប់មិនសុំច្បាប់',
    description: 'កត់ត្រាអវត្តមាន',
    path: '/Notpermissionleave',
    icon: 'CircleClose',
    show: adminLevel.value
  },
  {
    label: 'កត់ត្រាវត្តមាន',
    description: 'ស្កែន និងកត់ត្រាវត្តមាន',
    path: '/createattendance',
    icon: 'Camera',
    show: addAttendance.value
  },
  {
    label: 'ថ្នាក់រៀន',
    description: 'គ្រប់គ្រងថ្នាក់រៀន',
    path: '/company',
    icon: 'School',
    show: true
  },
  {
    label: 'មុខវិជ្ជា',
    description: 'គ្រប់គ្រងមុខវិជ្ជា',
    path: '/subject',
    icon: 'Reading',
    show: adminLevel.value
  },
  {
    label: 'ជំនាញ',
    description: 'គ្រប់គ្រងជំនាញសិក្សា',
    path: '/Major',
    icon: 'Collection',
    show: adminLevel.value
  },
  {
    label: 'កែវត្តមាន',
    description: 'កែវត្តមាន',
    path: '/AttendanceForUpdate',
    icon: 'Edit',
    show: adminLevel.value
  },
  {
    label: 'ជំនាន់',
    description: 'គ្រប់គ្រងជំនាន់',
    path: '/Generation',
    icon: 'Collection',
    show: adminLevel.value
  },
  {
    label: 'ម៉ោងសិក្សា',
    description: 'គ្រប់គ្រងម៉ោងសិក្សា',
    path: '/Shift',
    icon: 'Clock',
    show: adminLevel.value
  },
  {
    label: 'អ្នកប្រើប្រាស់/និស្សិត/សាស្រ្តាចារ្យ',
    description: 'គ្រប់គ្រងអ្នកប្រើប្រាស់',
    path: '/users',
    icon: 'User',
    show: adminLevel.value || staffLevel.value
  },
  {
    label: 'ចំណូលផ្សេងៗ',
    description: 'ចំណូលក្រៅពីថ្លៃសិក្សា',
    path: '/Income',
    icon: 'Money',
    show: viewIncome.value
  },
  {
    label: 'អាហារូបករណ៍',
    description: 'ចំណូលប្រភេទអាហារូបករណ៍សំខាន់ៗក្រៅពីថ្លៃសិក្សា',
    path: '/Schoolarship',
    icon: 'School',
    show: true
  },
  {
    label: 'របាយការណ៍ផ្សេងៗ',
    description: 'របាយការណ៍ផ្សេងៗ',
    path: '/Report',
    icon: 'Eleme',
    show: adminLevel.value || managerLevel.value
  }
].filter(item => item.show))

onMounted(() => {
  fetchCountLeave()
})
</script>


<template>
  <div class="dashboard">
    <el-card shadow="never" class="dash-card">
      <div class="action-grid">
        <div
          v-for="item in actions"
          :key="item.path"
          class="action-btn"
          @click="router.push(item.path)"
        >
          <div class="icon-box">
            <el-badge
              v-if="item.path === '/leaverequest' && leaveCount > 0"
              :value="leaveCount"
              :max="99"
            >
              <el-icon :size="23">
                <component :is="item.icon" />
              </el-icon>
            </el-badge>

            <el-icon v-else :size="23">
              <component :is="item.icon" />
            </el-icon>
          </div>

          <div class="action-content">
            <span class="action-title">{{ item.label }}</span>
            <span class="action-description">{{ item.description }}</span>
          </div>

          <el-icon class="arrow-icon">
            <ArrowRight />
          </el-icon>
        </div>
      </div>
    </el-card>
  </div>
</template>


<style scoped>
.dashboard {
  width: 100%;
}

.dash-card {
  background: #f1f1f1;
}

:deep(.el-card__body) {
  padding: 10px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 5px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 9px;
  min-height: 55px;
  padding: 7px 9px;
  background: #f8f8f8;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
}

.action-btn:hover {
  background: #dce9f5;
}

.icon-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  color: #315f88;
  background: #e6e6e6;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.action-content {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
  gap: 3px;
}

.action-title {
  font-size: 12px;
  font-weight: bold;
}

.action-description {
  font-size: 10px;
  color: #666;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.arrow-icon {
  color: #888;
}

@media (max-width: 1024px) {
  .action-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>
