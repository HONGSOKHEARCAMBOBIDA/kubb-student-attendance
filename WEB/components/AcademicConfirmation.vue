<template>
  <div style="text-align: right; margin-bottom: 8px">
    <AppButton icon="Printer" type="primary" @click="print">
      បោះពុម្ព
    </AppButton>
  </div>


  <div class="confirm-page" ref="pageEl">
    <!-- Photo box -->
    <div class="photo-box"><span>Photo</span></div>

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

    <!-- Data rows: left = Khmer, right = English -->
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

      <!-- ID + campus -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">អត្តលេខ:</span>
          <span class="val grow center">{{ d.code_kh }}</span>
          <span class="lbl">{{ d.campus_kh }}</span>
        </div>
        <div class="side en">
          <span class="lbl">ID:</span>
          <span class="val grow center">{{ d.code_en }}</span>
          <span class="lbl">{{ d.campus_en }}</span>
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
          <span class="val grow">{{ d.pob_kh }}</span>
        </div>
        <div class="side en">
          <span class="lbl">Place of Birth :</span>
          <span class="val grow">{{ d.pob_en }}</span>
        </div>
      </div>

      <!-- Address line 1 -->
      <div class="row">
        <div class="side kh">
          <span class="lbl">អាសយដ្ឋានបច្ចុប្បន្ន :ភូមិ</span>
          <span class="val grow center">{{ d.village_kh }}</span>
          <span class="lbl">សង្កាត់</span>
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
          <span class="lbl">ក្រុង</span>
          <span class="val grow center">{{ d.district_kh }}</span>
          <span class="lbl">ខេត្ត</span>
          <span class="val grow center">{{ d.province_kh }}</span>
        </div>
        <div class="side en">
          <span class="lbl">Sangkat</span>
          <span class="val grow center">{{ d.communce_en }}</span>
          <span class="lbl">City</span>
          <span class="val grow center">{{ d.district_en }}</span>
          <span class="lbl">Province</span>
        </div>
      </div>

      <!-- Level / year / semester -->
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
          <span class="val bold">{{ d.programme_en }}</span>
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
          <span class="val grow center">{{ d.start_kh }}-{{ d.end_kh }}</span>
          <span class="lbl">ប្រាកដមែន ។</span>
        </div>
        <div class="side en">
          <span class="lbl">In the academic year</span>
          <span class="val grow center">{{ d.start_en }}-{{ d.end_en }}</span>
        </div>
      </div>
    </div>

    <!-- Purpose -->
    <div class="purpose">
      <p class="purpose-kh">លិខិតនេះ ចេញឱ្យសាមីខ្លួនប្រើប្រាស់តាមការដែលអាចប្រើបាន ។</p>
      <p class="purpose-en">This confirmation is issued for any official purposes it may serve</p>
    </div>

    <!-- Signature -->
    <div class="sign">
      <p class="kh-line">{{ issuedKhLunar }}</p>
      <p class="kh-line">{{ issuedKh }}</p>
      <p class="en-line">{{ issuedEn }}</p>
      <p class="moul sign-title">ជ.សាកលវិទ្យាធិការ</p>
      <p class="moul sign-title">{{ signerTitle }}</p>
      <div class="space"></div>
      <p class="moul sign-title">{{ signerName }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import AppButton from "./AppButton.vue";

const taktieng = "/image.png";
const pageEl = ref(null);

const props = defineProps({
  // the object returned by the backend (AcademicConfirmation)
  data: { type: Object, default: null },
  // e.g. "ថ្ងៃអង្គារ ១រោច ខែអស្សុជ ឆ្នាំមមី អដ្ឋស័ក ព.ស.២៥៧០"
  issuedKhLunar: { type: String, default: "" },
  // e.g. "រាជធានីភ្នំពេញ ថ្ងៃទី២៧ ខែតុលា ឆ្នាំ២០២៦"
  issuedKh: { type: String, default: "" },
  // e.g. "Phnom Penh, date: October 27, 2026"
  issuedEn: { type: String, default: "" },
  signerTitle: { type: String, default: "សាកលវិទ្យាធិការរង" },
  signerName: { type: String, default: "" },
});

const d = computed(() => props.data);

// 1 -> st, 2 -> nd, 3 -> rd, others -> th (accepts "1" or 1)
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

/*
 * Why the old version was cut off:
 * the form lived inside the dialog, which has fixed width/height and
 * overflow hidden, so the print output was clipped at the dialog's edges.
 * Fix: copy the form into a temporary element directly under <body>,
 * print only that element, then remove it.
 */
async function print() {


  const src = pageEl.value;
  if (!src) return;

  // Wait for images
  const images = src.querySelectorAll("img");

  for (const img of images) {
    if (!img.complete) {
      await new Promise((resolve) => {
        img.onload = resolve;
        img.onerror = resolve;
      });
    }
  }

  await document.fonts.ready;

  // Remove previous print element
  document.getElementById("confirm-page")?.remove();

  // Clone original element
  const root = document.createElement("div");
  root.id = "confirm-page";

  const clone = src.cloneNode(true);
  root.appendChild(clone);

  document.body.appendChild(root);
  document.body.classList.add("confirm-page");

  // Cleanup
  let cleaned = false;

  const cleanup = () => {
    if (cleaned) return;
    cleaned = true;

    window.removeEventListener("afterprint", cleanup);
    document.body.classList.remove("confirm-page");
    root.remove();
  };

  window.addEventListener("afterprint", cleanup);

  // Wait for browser layout
  await new Promise((resolve) =>
    requestAnimationFrame(() =>
      requestAnimationFrame(resolve)
    )
  );

  await document.fonts.ready;

  try {
    window.print();
  } catch (e) {
    cleanup();
    throw e;
  }

  setTimeout(cleanup, 60000);
}

defineExpose({ print });
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Moul&display=swap');

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
  font-weight: 400;
}

/* Photo */
.photo-box {
  position: absolute;
  top: 0;
  right: 4px;
  width: 100px;
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
  font-size: 22px;
  margin: 0;
}

.title-en {
  font-family: Arial, sans-serif;
  font-weight: 700;
  font-size: 22px;
  margin: 2px 0 0;
}

.taktieng-logo {
  width: 110px;
  height: 25px;
  display: block;
  margin: 6px auto;
}

/* Authority block */
.authority {
  display: grid;
  grid-template-columns: 1fr 1fr;
  margin: 14px 0 12px;
  text-align: center;
}

.auth-kh {
  font-size: 14px;
  line-height: 1.9;
}

.auth-en {
  font-family: Arial, sans-serif;
  font-size: 14px;
  line-height: 1.9;
}

/* Rows: minmax(0, 1fr) lets both columns shrink instead of overflowing */
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
  font-size: 12px;
}

.side.kh {
  font-size: 13px;
}

.lbl {
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
}

.val.grow {
  flex: 1 1 0;
  min-width: 0;
}

.val.center {
  text-align: center;
}

.val.bold {
  font-weight: 700;
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
  font-size: 16px;
}

.purpose-en {
  font-family: Arial, sans-serif;
  font-size: 14px;
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
  font-size: 15px;
}

.en-line {
  font-family: Arial, sans-serif;
  font-size: 14px;
}

.sign-title {
  font-size: 14px;
}

.space {
  height: 60px;
}
</style>

<style>
/* Hidden on screen; only shown while printing */

/* Hide print copy on screen */
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

  /* Hide everything except print copy */
  body.confirm-page > *:not(#confirm-page) {
    display: none !important;
  }

  /* Show print copy */
  body.confirm-page #confirm-page {
    display: block !important;
    position: static !important;
    width: 210mm !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  #confirm-page .confirm-page {
    display: block !important;
    position: relative !important;
    width: 210mm !important;
    min-height: 297mm;
    margin: 0 !important;
    padding: 0 4px 8px !important;
    box-sizing: border-box;
    background: #fff !important;
    color: #000 !important;
  }
}

</style>