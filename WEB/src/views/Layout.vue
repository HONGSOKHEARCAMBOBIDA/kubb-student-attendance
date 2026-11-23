
<template>
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed }">
      <div class="sidebar-header">
        <el-image src="/logo.png" fit="contain" class="logo" />

        <span v-if="!collapsed" class="brand">
          សាកលវិទ្យាល័យខេមរាវិទូ
        </span>
      </div>

      <el-menu
        :default-active="route.path"
        :collapse="collapsed"
        router
      >
        <el-menu-item v-for="item in nav" :key="item.path" :index="item.path">
          <el-icon>
            <component :is="item.icon" />
          </el-icon>

          <template #title>
            {{ item.label }}
          </template>
        </el-menu-item>
      </el-menu>
    </aside>

    <!-- Main -->
    <div class="wrap">
      <header class="topbar">
        <div class="left">
          <el-button
            class="desktop"
            :icon="collapsed ? 'Expand' : 'Fold'"
            circle
            plain
            @click="collapsed = !collapsed"
          />

          <b>{{ pageTitle }}</b>
        </div>

        <!-- User menu -->
        <el-dropdown @command="handleCommand">
          <div class="user">
            <b>{{ userName }}</b>
          </div>

          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                ប្រវត្តរូប
              </el-dropdown-item>

              <el-dropdown-item command="logout" divided>
                ចាកចេញ
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </header>

      <main class="content">
        <router-view />
      </main>
    </div>

    <!-- Mobile navigation -->
    <nav class="bottom-nav">
      <router-link
        v-for="item in nav"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: route.path === item.path }"
      >
        <el-icon size="20">
          <component :is="item.icon" />
        </el-icon>

        <span>{{ item.short }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useUserDataStore } from '../stores/user_data'

const route = useRoute()
const router = useRouter()

const auth = useAuthStore()
const userdataStore = useUserDataStore()

const collapsed = ref(true)

const userName = computed(() => userdataStore.name || 'User')
const pageTitle = computed(() => route.meta.title || '')

// Get menu from router
const nav = computed(() =>
  router.getRoutes()
    .filter(item => item.meta.showInNav)
    .filter(item =>
      !item.meta.permission ||
      userdataStore.permissions?.some(
        permission => permission.name === item.meta.permission
      )
    )
    .map(item => ({
      path: item.path,
      icon: item.meta.icon,
      label: item.meta.title,
      short: item.meta.short
    }))
)

// Dropdown menu
function handleCommand(command) {
  if (command === 'logout') {
    auth.logout()
    router.push('/login')
    return
  }

  router.push('/profile')
}

onMounted(() => {
  userdataStore.getuserdata()
})
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: #f1f1f1;
}

/* Sidebar */
.sidebar {
  width: 220px;
  background: #e5e5e5;
  border-right: 1px solid #ccc;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 48px;
  padding: 0 12px;
  background: #d5d5d5;
}

.logo {
  width: 30px;
  height: 30px;
}

.brand {
  font-size: 14px;
  font-weight: bold;
  white-space: nowrap;
}

:deep(.el-menu) {
  border: none;
  background: #e5e5e5;
}

:deep(.el-menu-item) {
  height: 36px;
  border-bottom: 1px solid #ddd;
  color: #000;
}

:deep(.el-menu-item:hover) {
  background: #d5d5d5;
}

:deep(.el-menu-item.is-active) {
  background: #c5d9ed;
  color: #123d63;
}

/* Main */
.wrap {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 45px;
  padding: 0 12px;
  background: #dedede;
}

.left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user {
  padding: 5px 8px;
  cursor: pointer;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

/* Mobile */
.bottom-nav {
  display: none;
}

@media (max-width: 768px) {
  .sidebar,
  .desktop {
    display: none;
  }

  .content {
    padding-bottom: 70px;
  }

  .topbar {
    height: 50px;
  }

  .bottom-nav {
    display: flex;
    position: fixed;
    bottom: 0;
    width: 100%;
    height: 58px;
    background: white;
    border-top: 1px solid #ddd;
  }

  .nav-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    color: #999;
    font-size: 11px;
    text-decoration: none;
  }

  .nav-item.active {
    color: #409eff;
  }
}
</style>
