```vue
<template>
  <div class="login-page">
    <div class="login-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        @submit.prevent="handleLogin"
        label-position="top"
      >
        <AppInput
          v-model="form.code"
          label="លេខកូដ"
          prop="code"
          placeholder="បញ្ចូលលេខកូដ"
          prefix-icon="User"
          clearable
          autofocus
        />

        <AppInput
          v-model="form.password"
          label="ពាក្យសម្ងាត់"
          prop="password"
          type="password"
          placeholder="បញ្ចូលពាក្យសម្ងាត់"
          prefix-icon="Lock"
          show-password
          @enter="handleLogin"
        />

        <AppButton
          native-type="submit"
          :loading="loading"
          type="primary"
          block
        >
          ចូលប្រព័ន្ធ
        </AppButton>
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
      message: "សូមបញ្ចូលលេខកូដ",
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

    notify.success("ចូលប្រព័ន្ធបានជោគជ័យ");
  } catch (e) {
    if (e?.response) {
      notify.error(
        e.response?.data?.message || "ការចូលប្រព័ន្ធបរាជ័យ"
      );
    }
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-card {
  width: 420px;
  padding: 48px 40px;
  border: 1px solid #8aaff5;
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
```
