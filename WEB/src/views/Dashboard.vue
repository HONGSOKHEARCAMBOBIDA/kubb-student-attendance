
<template>
  <div class="dashboard">
    <div class="dash-grid">
      <el-card class="dash-card" shadow="never">
        <template #header>
          <div>
            <div class="card-title">សកម្មភាពរហ័ស</div>
            <div class="card-subtitle">
              ជ្រើសរើសមុខងារដែលអ្នកចង់ប្រើ
            </div>
          </div>
        </template>

        <div class="action-grid">

          <div
            class="action-btn"
            @click="$router.push('/attendance')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Calendar />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">មើលវត្តមានសង្ខែប</span>
              <span class="action-description">
                ពិនិត្យវត្តមានសង្ខែប
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
            class="action-btn"
            @click="$router.push('/AttendanceReport')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <DataAnalysis />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">របាយការណ៍វត្តមាន</span>
              <span class="action-description">
                របាយការណ៍វត្តមានសរុប
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
            class="action-btn"
            @click="$router.push('/leaverequest')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <DocumentChecked />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">សុំច្បាប់</span>
              <span class="action-description">
                គ្រប់គ្រងការសុំច្បាប់
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
            v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/Notpermissionleave')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <CircleClose />
              </el-icon>
            </div>

            <div class="action-content">

              <span class="action-title">ឈប់មិនសុំច្បាប់</span>
              <span class="action-description">
                កត់ត្រាអវត្តមាន
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
            class="action-btn"
            @click="$router.push('/createattendance')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Camera />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">កត់ត្រាវត្តមាន</span>
              <span class="action-description">
                ស្កែន និងកត់ត្រាវត្តមាន
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
          v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/company')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <School />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">ថ្នាក់រៀន</span>
              <span class="action-description">
                គ្រប់គ្រងថ្នាក់រៀន
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
          v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/subject')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Reading />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">មុខវិជ្ជា</span>
              <span class="action-description">
                គ្រប់គ្រងមុខវិជ្ជា
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

          <div
          v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/Major')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Collection />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">ជំនាញ</span>
              <span class="action-description">
                គ្រប់គ្រងជំនាញសិក្សា
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

                    <div
          v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/AttendanceForUpdate')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Collection />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">កែវត្តមាន</span>
              <span class="action-description">
                កែវត្តមាន
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

           <div
          v-if="adminLevel"
            class="action-btn"
            @click="$router.push('/Generation')"
          >
            <div class="icon-box">
              <el-icon :size="23">
                <Collection />
              </el-icon>
            </div>

            <div class="action-content">
              <span class="action-title">ជំនាន់</span>
              <span class="action-description">
                ជំនាន់
              </span>
            </div>

            <el-icon class="arrow-icon">
              <ArrowRight />
            </el-icon>
          </div>

        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useUserDataStore } from "../stores/user_data";
const userDataStore = useUserDataStore();
const adminLevel = computed(() => userDataStore.level === 7)
</script>

<style scoped>
.dashboard {
  width: 100%;
}

.dash-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.dash-card {
  border: 1px solid #e5e9ef;
  border-radius: 10px;
}

/* Header */
.card-title {
  font-size: 15px;
  font-weight: 700;
  color: #122133;
}

.card-subtitle {
  margin-top: 3px;
  font-size: 12px;
  color: #7b8794;
}

/* Actions */
.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 12px;

  min-height: 72px;
  padding: 12px 14px;

  background: #fff;
  border: 1px solid #e5e9ef;
  border-radius: 9px;

  cursor: pointer;

  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.action-btn:hover {
  border-color: #122133;
  box-shadow: 0 4px 12px rgba(18, 33, 51, 0.08);
  transform: translateY(-1px);
}

/* Icon */
.icon-box {
  width: 42px;
  height: 42px;
  flex: 0 0 42px;

  display: flex;
  align-items: center;
  justify-content: center;

  color: #122133;
  background: #f2f4f7;
  border-radius: 9px;
}

/* Content */
.action-content {
  flex: 1;
  min-width: 0;

  display: flex;
  flex-direction: column;
  gap: 3px;
}

.action-title {
  font-size: 13px;
  font-weight: 600;
  color: #122133;
}

.action-description {
  font-size: 11px;
  color: #7b8794;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Arrow */
.arrow-icon {
  color: #a3adb8;
  transition: transform 0.18s ease;
}

.action-btn:hover .arrow-icon {
  color: #122133;
  transform: translateX(3px);
}

/* Mobile */
@media (max-width: 767px) {
  .action-grid {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .action-btn {
    min-height: 68px;
  }
}

@media (max-width: 400px) {
  .action-title {
    font-size: 12px;
  }

  .action-description {
    font-size: 10px;
  }
}

</style>