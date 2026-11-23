
<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-logo">
        <el-image src="/logo.png" alt="University Logo" fit="contain" />
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        @submit.prevent="handleLogin"
        label-position="top"
      >
        <AppInput
          v-model.trim="form.code"
          label="អត្តលេខ"
          prop="code"
          placeholder="បញ្ចូលលេខកូដ"
          prefix-icon="User"
          clearable
          autofocus
        />

        <AppInput
          v-model.trim="form.password"
          label="ពាក្យសម្ងាត់"
          prop="password"
          type="password"
          placeholder="បញ្ចូលពាក្យសម្ងាត់"
          prefix-icon="Lock"
          show-password
          @enter="handleLogin"
        />
        <AppButton :loading="loading" type="primary" block @click="handleLogin">ចូលប្រព័ន្ធ</AppButton>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { login } from "../api/services";
import AppButton from "../../components/AppButton.vue";
import AppInput from "../../components/AppInput.vue";
import { useNotification } from "../../composables/useNotification.js";
import { ElMessage } from "element-plus";
const notify = useNotification();
const router = useRouter();
const auth = useAuthStore();

const formRef = ref();
const loading = ref(false);

const form = reactive({
  code: "",
  password: "",
});

const rules = {
  code: [
    {
      required: true,
      message: "សូមបញ្ចូលអត្តលេខ",
      trigger: "blur",
    },
  ],
  password: [
    {
      required: true,
      message: "សូមបញ្ចូលពាក្យសម្ងាត់",
      trigger: "blur",
    },
  ],
};

async function handleLogin() {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    loading.value = true;

    const res = await login({
      code: form.code,
      password: form.password,
    });

    auth.setAuth(res.data.data);

    await router.push("/Dashboard");

    ElMessage.success("ចូលប្រព័ន្ធបានជោគជ័យ");
  } catch (e) {
    if (e?.response) {
      notify.error(e.response?.data?.message || "ការចូលប្រព័ន្ធបរាជ័យ");
    }
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 5px;
}
.login-logo .el-image {
  width: 80px;
  height: 80px;
}
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
}

.login-card {
  width: 300px;
  padding: 38px 30px;
  border: 1px solid #b9b9b9;
  border-radius: 2px;
}

@media (max-width: 768px) {
  .login-page {
    padding: 25px;
  }

  .login-card {
    width: 100%;
    max-width: 420px;
    padding: 24px 20px;
    border-radius: 8px;
  }
}
</style>

