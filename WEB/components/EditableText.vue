<template>
  <p
    ref="el"
    class="editable"
    contenteditable="true"
    spellcheck="false"
    @input="emitValue"
    @keydown.enter.prevent="$event.target.blur()"
    @paste="onPaste"
  ></p>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";

const props = defineProps({
  modelValue: { type: String, default: "" },
});

const emit = defineEmits(["update:modelValue"]);
const el = ref(null);

function sync() {
  if (el.value && el.value.innerText !== props.modelValue) {
    el.value.innerText = props.modelValue;
  }
}

onMounted(sync);

// Don't touch the DOM while typing, otherwise the cursor jumps
watch(
  () => props.modelValue,
  () => {
    if (document.activeElement !== el.value) sync();
  }
);

function emitValue() {
  emit("update:modelValue", el.value.innerText);
}

// Paste plain text only
function onPaste(e) {
  e.preventDefault();
  const text = e.clipboardData.getData("text/plain").replace(/\n/g, " ");
  document.execCommand("insertText", false, text);
}
</script>

<style scoped>
.editable {
  outline: none;
  cursor: text;
  border-radius: 3px;
  min-height: 1.5em;
  transition: background 0.15s, outline-color 0.15s;
}
.editable:hover {
  outline: 1px dashed #409eff;
  background: rgba(64, 158, 255, 0.06);
}
.editable:focus {
  outline: 1px solid #409eff;
  background: rgba(64, 158, 255, 0.1);
}
.editable:empty::before {
  content: "ចុចដើម្បីកែប្រែ...";
  color: #aaa;
}

@media print {
  .editable,
  .editable:hover,
  .editable:focus {
    outline: none !important;
    background: none !important;
  }
  .editable:empty::before {
    content: none;
  }
}
</style>