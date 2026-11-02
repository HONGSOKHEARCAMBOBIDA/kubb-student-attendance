
<template>
  <div style="text-align: right; margin-bottom: 8px">
    <AppButton icon="Printer" type="primary" @click="print">
      បោះពុម្ព
    </AppButton>
  </div>

  <div class="confirm-page" ref="pageEl">
    <!-- Photo -->
    <div class="photo-box">
      <span>Photo</span>
    </div>

    <!-- Title -->
    <div class="title-block">
      <h2 class="moul title-kh">លិខិតបញ្ជាក់ការសិក្សា</h2>
      <h2 class="title-en">ACADEMIC CONFIRMATION</h2>
      <img :src="taktieng" class="taktieng-logo" alt="" />
    </div>

    <!-- Authority -->
    <div class="authority">
      <div class="auth-kh moul">
        <div>សាកលវិទ្យាធិការ</div>
        <div>នៃសាកលវិទ្យាល័យខេមរៈ</div>
        <div>បញ្ជាក់ថា</div>
      </div>

      <div class="auth-en">
        <div>Under Authority of the rector</div>
        <div>of Khemarak University</div>
        <div>Certifies that</div>
      </div>
    </div>

    <!-- Data rows -->
    <div class="rows" v-if="d">
      <!-- Name -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ឈ្មោះ :</span>
          <span class="val grow bold">{{ d.name_kh }}</span>
          <span class="lbl">ភេទ:</span>
          <span class="val">{{ d.gender_kh }}</span>
          <span class="lbl">សញ្ជាតិ</span>
          <span class="val">{{ d.national_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Name</span>
          <span class="val grow">{{ d.name_en }}</span>
          <span class="lbl">Sex:</span>
          <span class="val">{{ d.gender_en }}</span>
          <span class="lbl">Nationality:</span>
          <span class="val">{{ d.national_en }}</span>
        </div>
      </div>

      <!-- ID + Campus -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">អត្តលេខ:</span>
          <span class="val grow center">{{ d.code_kh }}</span>
          <span class="lbl">{{ d.campus_kh }}</span>
          <span class="lbl">សាខាខេត្តបាត់ដំបង</span>
        </div>

        <div class="side en">
          <span class="lbl">ID:</span>
          <span class="val grow center">{{ d.code_en }}</span>
          <span class="lbl">{{ d.campus_en }}</span>
          <span class="lbl">Battambang Branch</span>
        </div>
      </div>

      <!-- DOB -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ថ្ងៃ-ខែ-ឆ្នាំកំណើត :</span>
          <span class="val grow center">{{ d.dob_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Date of Birth :</span>
          <span class="val grow center">{{ d.dob_en }}</span>
        </div>
      </div>

      <!-- Place of birth -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ទីកន្លែងកំណើត :</span>
          <EditableText class="kh-line val grow" v-model="provincebirthkh" />
        </div>

        <div class="side en">
          <span class="lbl">Place of Birth :</span>
          <EditableText class="en-line val grow" v-model="provincebirthen" />
        </div>
      </div>

      <!-- Address line 1 -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">អាសយដ្ឋានបច្ចុប្បន្ន :ភូមិ</span>
          <span class="val grow center">{{ d.village_kh }}</span>
          <span class="lbl">{{ d.is_capital === true ? 'សង្កាត' : 'ឃុំ' }}</span>
          <span class="val grow center">{{ d.communce_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Current Address :</span>
          <span class="val grow center">{{ d.village_en }}</span>
          <span class="lbl">Village</span>
          <span class="val grow center">{{ d.communce_en }}</span>
        </div>
      </div>

      <!-- Address line 2 -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">{{ d.is_capital === true ? 'ខណ្ឌ' : 'ស្រុក'}}</span>
          <span class="val grow center">{{ d.district_kh }}</span>
          <span class="lbl">{{ d.is_capital === true ? 'រាជធានី' : 'ខេត្ត'}}</span>
          <span class="val grow center">{{ d.province_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">{{ d.is_capital === true ? 'Sangkat' : 'Commune' }}</span>
          <span class="val grow center">{{ d.district_en }}</span>
          <span class="lbl">{{ d.is_capital === true ? 'Khan' : 'District' }}</span>
          <span class="val grow center">{{ d.province_en }}</span>
          <span class="lbl">{{ d.is_capital === true ? 'Capital' : 'Province' }}</span>        
        </div>
      </div>

      <!-- Level / Year / Semester -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ពិតជានិស្សិតថ្នាក់</span>
          <span class="val bold">{{ d.programme_kh }}</span>
          <span class="lbl">ឆ្នាំទី</span>
          <span class="val num">{{ d.year_kh }}</span>
          <span class="lbl">ឆមាសទី</span>
          <span class="val num">{{ d.semester_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Is The</span>
          <span class="val num">
            {{ d.semester_en }}<sup>{{ ordinal(d.semester_en) }}</sup>
          </span>

          <span class="lbl">Semester,</span>
          <span class="val num">
            {{ d.year_en }}<sup>{{ ordinal(d.year_en) }}</sup>
          </span>

          <span class="lbl">Year of</span>
          <span class="val bold simple">{{ d.programme_en }}</span>
          <span class="lbl">'s Student</span>
        </div>
      </div>

      <!-- Faculty -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">មហាវិទ្យាល័យ</span>
          <span class="val grow center">{{ d.faculty_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Faculty of</span>
          <span class="val grow center">{{ d.faculty_en }}</span>
        </div>
      </div>

      <!-- Major -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ឯកទេស</span>
          <span class="val grow center">{{ d.major_kh }}</span>
        </div>

        <div class="side en">
          <span class="lbl">Speciality in</span>
          <span class="val grow center">{{ d.major_en }}</span>
        </div>
      </div>

      <!-- Academic year -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">ក្នុងឆ្នាំសិក្សា</span>
          <span class="val grow center">
            {{ d.start_kh }}-{{ d.end_kh }}
          </span>
          <span class="lbl">ប្រាកដមែន ។</span>
        </div>

        <div class="side en">
          <span class="lbl">In the academic year</span>
          <span class="val grow center">
            {{ d.start_en }}-{{ d.end_en }}
          </span>
        </div>
      </div>
    </div>

    <!-- Purpose -->
    <div class="purpose">
      <p class="purpose-kh">
        លិខិតនេះ ចេញឱ្យសាមីខ្លួនប្រើប្រាស់តាមការដែលអាចប្រើបាន ។
      </p>

      <p class="purpose-en">
        This confirmation is issued for any official purposes it may serve
      </p>
    </div>

    <!-- Signature -->
<!-- Signature -->
<div class="sign">
  <EditableText class="kh-line" v-model="lunar" />
  <EditableText class="kh-line" v-model="issuedKhText" />
  <EditableText class="en-line" v-model="issuedEnText" />
  <EditableText class="moul sign-title" v-model="titleText" />
  <EditableText class="moul sign-title" v-model="signerTitleText" />

  <div class="space"></div>

  <EditableText class="moul sign-title" v-model="signerNameText" />
</div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import AppButton from "./AppButton.vue";
import EditableText from "./EditableText.vue";
const taktieng = "/image.png";
const pageEl = ref(null);

const props = defineProps({
  data: { type: Object, default: null },
  issuedKhLunar: { type: String, default: "" },
  issuedKh: { type: String, default: "" },
  issuedEn: { type: String, default: "" },
  signerTitle: { type: String, default: "សាកលវិទ្យាធិការរង" },
  signerName: { type: String, default: "" },
});

const lunar = ref(props.issuedKhLunar);
const issuedKhText = ref(props.issuedKh);
const issuedEnText = ref(props.issuedEn);
const titleText = ref("ជ.សាកលវិទ្យាធិការ");
const signerTitleText = ref(props.signerTitle);
const signerNameText = ref(props.signerName);
const provincebirthkh = ref("")
const provincebirthen = ref("")

watch(() => props.issuedKhLunar, (v) => (lunar.value = v));
watch(() => props.issuedKh, (v) => (issuedKhText.value = v));
watch(() => props.issuedEn, (v) => (issuedEnText.value = v));
watch(() => props.signerTitle, (v) => (signerTitleText.value = v));
watch(() => props.signerName, (v) => (signerNameText.value = v));

const d = computed(() => props.data);

// Ordinal: 1st, 2nd, 3rd, 4th...
function ordinal(n) {
  const v = Number(n);

  if (!v) return "";

  const mod100 = v % 100;

  if (mod100 >= 11 && mod100 <= 13) return "th";

  switch (v % 10) {
    case 1:
      return "st";
    case 2:
      return "nd";
    case 3:
      return "rd";
    default:
      return "th";
  }
}

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

// Wait for browser layout
function waitForLayout() {
  return new Promise((resolve) => {
    requestAnimationFrame(() => {
      requestAnimationFrame(resolve);
    });
  });
}

// Print
async function print() {
  const src = pageEl.value;

  if (!src) {
    console.error("Print element not found");
    return;
  }

  let root = null;
  let cleaned = false;

  // Remove previous print copy
  document.getElementById("confirm-page")?.remove();

  try {
    // Wait for original images and fonts
    await waitForImages(src);
    await document.fonts.ready;

    // Create temporary print element
    root = document.createElement("div");
    root.id = "confirm-page";

    // Clone original content
    const clone = src.cloneNode(true);

    root.appendChild(clone);
    document.body.appendChild(root);

    // Activate print mode
    document.body.classList.add("confirm-page");

    // Wait for cloned images and layout
    await waitForImages(clone);
    await document.fonts.ready;
    await waitForLayout();

    // Check print content
    if (!root.querySelector(".confirm-page")) {
      throw new Error("Print content not found");
    }

    // Cleanup after printing
    const cleanup = () => {
      if (cleaned) return;

      cleaned = true;

      window.removeEventListener("afterprint", cleanup);

      document.body.classList.remove("confirm-page");

      root?.remove();
    };

    window.addEventListener("afterprint", cleanup);

    // Open print preview
    window.print();

    // Fallback cleanup
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
@import url("https://fonts.googleapis.com/css2?family=Moul&family=Battambang:wght@400;700&display=swap");

.confirm-page {
  position: relative;
  box-sizing: border-box;
  width: 100%;
  background: #fff;
  color: #000;
  padding: 0 4px 8px;
  font-family: "Battambang", "Khmer OS Battambang", Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
}

.moul {
  font-family: "Moul", serif;
  font-weight: 500;
}

/* Photo */
.photo-box {
  position: absolute;
  top: 0;
  right: 4px;
  width: 120px;
  height: 128px;
  box-sizing: border-box;
  border: 2px solid #000;
  border-radius: 6px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 6px;
  font-family: Arial, sans-serif;
  font-weight: 700;
  font-size: 12px;
  letter-spacing: 2px;
}

/* Title */
.title-block {
  text-align: center;
  padding-top: 60px;
}

.title-kh {
  font-size: 20px;
  margin: 0;
}

.title-en {
  font-family: Arial, sans-serif;
  font-weight: 700;
  font-size: 20px;
  margin: 2px 0 0;
}

.taktieng-logo {
  width: 110px;
  height: 25px;
  display: block;
  margin: 6px auto;
}

/* Authority */
.authority {
  display: grid;
  grid-template-columns: 1fr 1fr;
  margin: 10px 0 8px;
  text-align: center;
}

.auth-kh {
  font-size: 16px;
  line-height: 1.9;
}

.auth-en {
  font-family: Arial, sans-serif;
  font-size: 16px;
  line-height: 1.9;
}

/* Rows */
.rows {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  column-gap: 12px;
  break-inside: avoid;
}

.side {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  min-width: 0;
}

.side.en {
  font-family: Arial, sans-serif;
  font-size: 14px;
}

.side.kh {
  font-size: 16px;
}

.lbl {
  font-size: 17px;
  white-space: nowrap;
  flex-shrink: 0;
}

.val {
  border-bottom: 1px dotted #000;
  padding: 0 3px;
  min-width: 18px;
  min-height: 1.5em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: clip;
  flex-shrink: 1;
  text-align: center;
}

.val.grow {
  flex: 1 1 0;
  min-width: 0;
  align-items: center;
}

.val.center {
  text-align: center;
}

.val.bold {
  font-family: "Moul", serif;
  font-weight: 400;
  align-items: center;
}

.val.bold.simple {
  font-family: Arial, sans-serif;
  font-weight: 600;
  align-items: center;
}

.val.num {
  text-align: center;
  min-width: 24px;
}

sup {
  font-size: 9px;
}

/* Purpose */
.purpose {
  text-align: center;
  margin-top: 18px;
}

.purpose p {
  margin: 4px 0;
}

.purpose-kh {
  font-size: 18px;
}

.purpose-en {
  font-family: Arial, sans-serif;
  font-size: 18px;
}

/* Signature */
.sign {
  width: 50%;
  margin: 18px 0 0 auto;
  text-align: center;
  break-inside: avoid;
}

.sign p {
  margin: 3px 0;
}

.kh-line {
  font-size: 17px;
}

.en-line {
  font-family: Arial, sans-serif;
  font-size: 17px;
}

.sign-title {
  font-size: 17px;
}

.space {
  height: 60px;
}
</style>

<!-- Global print styles -->
<style>
/* Hide temporary print element on screen */
#confirm-page {
  display: none;
}

@media print {
  @page {
    size: A4 portrait;
    margin: 0 0.2mm;
   padding: 80px 0px 0px;
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

  /* Hide all original application content */
  body.confirm-page > *:not(#confirm-page) {
    display: none !important;
  }

  /* Show temporary print element */
  body.confirm-page > #confirm-page {
    display: block !important;
    position: static !important;
    width: 210mm !important;
    margin: 0 !important;
    padding: 0 !important;
    visibility: visible !important;
  }

  /* Print document */
  #confirm-page .confirm-page {
    display: block !important;
    position: relative !important;
    width: 210mm !important;
    min-height: 207mm;
    margin: 0 !important;
    padding: 0 0px 0px !important;
    box-sizing: border-box !important;
    background: #fff !important;
    color: #000 !important;
    visibility: visible !important;
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
