<template>
  <el-container class="layout-container-demo" style="height: 100%">
    <el-aside width="330px">
      <el-scrollbar>
<el-menu :default-openeds="['1']" :default-active="$route.path" router>
  <el-sub-menu index="1">
    <template #title>
      <el-icon><message /></el-icon>របាយការណ៍ផ្សេងៗផ្នែកកិច្ចការសាស្រ្តាចារ្យ
    </template>
    <el-menu-item-group>
      <el-menu-item index="/Report/TeacherSchedule">កាលវិភាគប្រចាំសាប្ដាហ៍</el-menu-item>
      <el-menu-item index="/Report/teaching-hours">របាយការណ៍ម៉ោងបង្រៀនរបស់សាស្រ្តាចារ្យ</el-menu-item>
      <el-menu-item index="/Report/teacher-attendandce">របាយការណ៍វត្តមានសាស្រ្តាចារ្យ</el-menu-item>
    </el-menu-item-group>
  </el-sub-menu>
</el-menu>
      </el-scrollbar>
    </el-aside>

    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { Message } from '@element-plus/icons-vue'
</script>

<style scoped>
.layout-container-demo .el-aside {
  color: var(--el-text-color-primary);
  background: var(--el-color-primary-light-8);
}
.layout-container-demo .el-menu {
  border-right: none;
}
</style>