<template>
  <el-button
    v-bind="$attrs"
    :type="type"
    :size="size"
    :loading="loading"
    :disabled="disabled || loading"
    :native-type="nativeType"
    :circle="circle"
    :class="{ 'app-button--block': block, 'app-button--circle': circle }"
    @click="$emit('click')"
  >
    <el-icon v-if="icon" color="#273b91"><component :is="icon" /></el-icon>
    <slot v-if="!circle" />
  </el-button>
</template>

<script setup>
defineProps({
  type: { type: String, default: "default" },
  size: { type: String, default: "small" },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  nativeType: { type: String, default: "button" },
  block: { type: Boolean, default: false },
  icon: { type: [String, Object], default: null },
  circle: { type: Boolean, default: false },
});
defineEmits(["click"]);
</script>

<style scoped>
.app-button--block {
  width: 100%;
}
</style>
