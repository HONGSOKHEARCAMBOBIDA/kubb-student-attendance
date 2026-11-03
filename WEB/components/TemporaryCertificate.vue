<template>
  <div style="text-align: right; margin-bottom: 8px">
    <AppButton icon="Printer" type="primary" @click="print">
      បោះពុម្ព
    </AppButton>
  </div>

  <div class="confirm-page" ref="pageEl">
    <!-- Photo box (left margin) -->


    <!-- Title -->
    <h2 class="title-kh">វិញ្ញាបនបត្របណ្ដោះអាសន្នថ្នាក់បរិញ្ញាបត្រ</h2>

    <!-- Authority -->
    <p class="authority">
      សាកលវិទ្យាធិការ នៃសាកលវិទ្យាល័យ ខេមរៈ បញ្ជាក់ថាៈ
    </p>

    <!-- Data rows -->
    <div class="rows" v-if="d">
    <div class="photo-box">
      <span>4x6</span>
    </div>
      <!-- Name + Latin name -->
      <div class="row">
        <span class="lbl w1">និស្សិតឈ្មោះ:</span>
        <span class="val half bold">{{ d.name_kh }}</span>
        <span class="lbl w2">ជាអក្សរឡាតាំង:</span>
        <span class="val grow latin">{{ d.name_en }}</span>
      </div>

      <!-- Sex + Nationality -->
      <div class="row">
        <span class="lbl w1">ភេទ</span>
        <span class="val half">{{ d.gender_kh }}</span>
        <span class="lbl w2">សញ្ជាតិ</span>
        <span class="val grow">{{ d.national_kh }}</span>
      </div>

      <!-- DOB -->
      <div class="row">
        <span class="lbl w1">ថ្ងៃ ខែ ឆ្នាំកំណើត</span>
        <span class="val grow left">{{ d.dob_kh }}</span>
      </div>

      <!-- Place of birth -->
      <div class="row">
        <span class="lbl w1">ទីកន្លែងកំណើត</span>
        <EditableText class="val grow left" v-model="provincebirthkh" />
      </div>

      <!-- Father -->
      <div class="row">
        <span class="lbl w1">ឪពុកឈ្មោះ:</span>
        <span class="val grow left">{{ d.father_name }}</span>
      </div>

      <!-- Mother -->
      <div class="row">
        <span class="lbl w1">ម្ដាយឈ្មោះ:</span>
        <span class="val grow left">{{ d.mother_name }}</span>
      </div>

      <!-- Branch -->
      <div class="row">
        <span class="lbl">បានបំពេញគ្រប់លក្ខខណ្ឌរបស់</span>
        <span class="val grow left bold">
          សាខាសាកលវិទ្យាល័យខេមរៈ ខេត្តបាត់ដំបង
        </span>
      </div>

      <!-- Programme -->
      <div class="row">
        <span class="lbl">ក្នុងការប្រឡងបញ្ចប់ការសិក្សាថ្នាក់:</span>
        <span class="val grow bold">{{ d.programme_kh }}</span>
      </div>

      <!-- Major -->
      <div class="row">
        <span class="lbl w1">ឯកទេស</span>
        <span class="val grow left bold">{{ d.major_kh }}</span>
      </div>

      <!-- Exam date -->
      <div class="row">
        <span class="lbl w1">សម័យប្រឡងថ្ងៃទី</span>
        <span class="val grow left">{{ d.exam_out_date }}</span>
      </div>
    </div>

    <!-- Purpose -->
    <p class="purpose">
      វិញ្ញាបនបត្រនេះ ចេញឱ្យសាមីជន យកទៅប្រើប្រាស់តាមការដែលអាចប្រើបាន
      និងមានសុពលភាពរយៈពេលពីរឆ្នាំ។
    </p>

    <!-- Signature -->
    <div class="sign">
<EditableText v-draggable class="sign-line"  v-model="lunar" />
<EditableText v-draggable class="sign-line"  v-model="issuedKhText" />
<EditableText v-draggable class="sign-title" v-model="titleText" />
<EditableText v-draggable class="sign-title" v-model="signerTitleText" />
<div class="space"></div>
<EditableText v-draggable class="sign-title" v-model="signerNameText" />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import AppButton from "./AppButton.vue";
import EditableText from "./EditableText.vue";

const pageEl = ref(null);

const props = defineProps({
  data: { type: Object, default: null },
  issuedKhLunar: { type: String, default: "" },
  issuedKh: { type: String, default: "" },
  signerTitle: { type: String, default: "សាកលវិទ្យាធិការរង" },
  signerName: { type: String, default: "" },
});

const lunar = ref(props.issuedKhLunar);
const issuedKhText = ref(props.issuedKh);
const titleText = ref("ជ.សាកលវិទ្យាធិការ");
const signerTitleText = ref(props.signerTitle);
const signerNameText = ref(props.signerName);
const provincebirthkh = ref("ខេត្ត បាត់ដំបង");

watch(() => props.issuedKhLunar, (v) => (lunar.value = v));
watch(() => props.issuedKh, (v) => (issuedKhText.value = v));
watch(() => props.signerTitle, (v) => (signerTitleText.value = v));
watch(() => props.signerName, (v) => (signerNameText.value = v));

const d = computed(() => props.data);

// Wait until images are loaded
async function waitForImages(element) {
  const images = element.querySelectorAll("img");
  await Promise.all(
    [...images].map((img) => {
      if (img.complete) return Promise.resolve();
      return new Promise((resolve) => {
        img.addEventListener("load", resolve, { once: true });
        img.addEventListener("error", resolve, { once: true });
      });
    })
  );
}

const vDraggable = {
  mounted(el) {
    let x = 0, y = 0;

    // stop text selection / focus when Alt is held
    el.addEventListener("mousedown", (e) => {
      if (e.altKey) e.preventDefault();
    });

    el.addEventListener("pointerdown", (e) => {
      if (!e.altKey) return;
      e.preventDefault();
      el.setPointerCapture(e.pointerId);

      const sx = e.clientX - x, sy = e.clientY - y;
      const move = (ev) => {
        x = ev.clientX - sx;
        y = ev.clientY - sy;
        el.style.translate = `${x}px ${y}px`;
      };
      const up = () => {
        el.removeEventListener("pointermove", move);
        el.removeEventListener("pointerup", up);
      };
      el.addEventListener("pointermove", move);
      el.addEventListener("pointerup", up);
    });
  },
};

function waitForLayout() {
  return new Promise((resolve) => {
    requestAnimationFrame(() => requestAnimationFrame(resolve));
  });
}

async function print() {
  const src = pageEl.value;
  if (!src) {
    console.error("Print element not found");
    return;
  }

  let root = null;
  let cleaned = false;

  document.getElementById("confirm-page")?.remove();

  try {
    await waitForImages(src);
    await document.fonts.ready;

    root = document.createElement("div");
    root.id = "confirm-page";

    const clone = src.cloneNode(true);
    root.appendChild(clone);
    document.body.appendChild(root);

    document.body.classList.add("confirm-page");

    await waitForImages(clone);
    await document.fonts.ready;
    await waitForLayout();

    if (!root.querySelector(".confirm-page")) {
      throw new Error("Print content not found");
    }

    const cleanup = () => {
      if (cleaned) return;
      cleaned = true;
      window.removeEventListener("afterprint", cleanup);
      document.body.classList.remove("confirm-page");
      root?.remove();
    };

    window.addEventListener("afterprint", cleanup);
    window.print();
    setTimeout(cleanup, 60000);
  } catch (error) {
    console.error("Print error:", error);
    document.body.classList.remove("confirm-page");
    root?.remove();
    throw error;
  }
}

defineExpose({ print });
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Moul&display=swap');

/*
  Layout is measured from the scanned form (A4, 210 x 297 mm):
  content starts 40 mm from the left, title sits ~64 mm from the top,
  photo box is in the left margin at ~139 mm from the top.
*/
.confirm-page {
  position: relative;
  box-sizing: border-box;
  width: 210mm;
  min-height: 296mm;
  margin: 0 auto;
  padding: 62mm 14mm 0 40mm;
  background: #fff;
  color: #000;
  font-family: "Battambang", "Khmer OS Battambang", Arial, sans-serif;
  font-size: 12pt;
  line-height: 1.5;
}

/* Photo box */
.photo-box {
  position: absolute;
  top: 50mm;    
  left: -33mm;         
  width: 24mm;
  height: 32mm;
  box-sizing: border-box;
  border: 0.3mm solid #000;
  border-radius: 2mm;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 1mm;
  font-family: Arial, sans-serif;
  font-size: 7pt;
  color: #444;
}

/* Title */
.title-kh {
  font-family: "Moul", serif;
  font-weight: 200;
  font-style: normal;
  margin: 0;
  font-size: 15pt;
  line-height: 1.6;
  text-decoration: underline;
  text-underline-offset: 1.2mm;
}

.authority {
  font-family: "Moul", serif;
  font-weight: 200;
  font-style: normal;
  margin: 4mm 0 6mm;
  font-size: 12.5pt;
}

/* Rows */
.rows {
  position: relative;   
  display: flex;
  flex-direction: column;
  gap: 1.2mm;
}

.row {
  display: flex;
  align-items: flex-end;
  gap: 2mm;
  min-height: 7.3mm;
  break-inside: avoid;
}

.lbl {
  font-size: 14pt;
  flex-shrink: 0;
  white-space: nowrap;
}

/* fixed label columns so values line up like the form */
.lbl.w1 {
  width: 32mm;
}
.lbl.w2 {
  width: 26mm;
}

.val {
  font-size: 14pt;
  min-width: 0;
  padding: 0 1mm;
  border-bottom: 0.3mm dotted #000;
  white-space: nowrap;
  overflow: hidden;
  text-align: center;
}

.val.grow {
  flex: 1 1 0;
}

.val.half {
  flex: 0 0 48mm;
}

.val.left {
  text-align: left;
}

.val.bold {
  font-family: "Moul", serif;
  font-weight: 200;
  font-style: normal; 
  font-size: 12pt; 
  padding: 1mm 1mm 0;
}

.val.latin {
  font-family: Arial, sans-serif;
  font-size: 11pt;
  font-weight: 700;
}

/* Purpose */
.purpose {
  font-family: "Moul", serif;
  font-weight: 200;
  margin: 6mm 0 0;
  font-size: 12.5pt;
  line-height: 1.75;
}

/* Signature */
.sign {
  width: 62%;
  margin: 8mm 0 0 auto;
  text-align: center;
  break-inside: avoid;
  cursor: pointer;
}

.sign-line {
  font-size: 12pt;
  line-height: 1.7;
}

.sign-title {
  font-family: "Moul", serif;
  font-weight: 200;  
  font-size: 13t;
  line-height: 1.7;
}

.space {
  height: 18mm;
}
</style>

<!-- Global print styles -->
<style>
#confirm-page {
  display: none;
}

@media print {
  @page {
    size: A4 portrait;
    margin: 0;
  }

  html,
  body {
    width: 100% !important;
    height: auto !important;
    margin: 0 !important;
    padding: 0 !important;
    overflow: visible !important;
    background: #fff !important;
  }

  body.confirm-page > *:not(#confirm-page) {
    display: none !important;
  }

  body.confirm-page > #confirm-page {
    display: block !important;
    position: static !important;
    width: 210mm !important;
    margin: 0 !important;
    padding: 0 !important;
    visibility: visible !important;
  }

  #confirm-page .confirm-page {
    display: block !important;
    position: relative !important;
    width: 210mm !important;
    min-height: 296mm;
    margin: 0 !important;
    box-sizing: border-box !important;
    background: #fff !important;
    color: #000 !important;
    visibility: visible !important;
     padding: 62mm 14mm 0 40mm !important;
  }

  #confirm-page .confirm-page * {
    visibility: visible !important;
  }

  #confirm-page .row,
  #confirm-page .sign {
    break-inside: avoid;
  }

  * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
  }
}
</style>