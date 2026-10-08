<script setup>
import { computed, nextTick } from "vue";
const props = defineProps({
  data: { type: Object, default: () => ({}) }, 
  receiver: { type: String, default: "" },
  logo: { type: String, default: "/logo.png" }, 
});
const d = computed(() => props.data ?? {});
const money = computed(() => `$${Number(d.value.total || 0).toFixed(2)}`);
/* ---------- print ---------- */
async function print() {
  await nextTick(); // make sure the latest data is rendered
  window.print();
}
defineExpose({ print });
</script>

<template>
  <Teleport to="body">
    <div class="invoice-print-root">
      <div class="sheet">
        <template v-for="copy in 2" :key="copy">
          <section class="receipt">
            <!-- ===== header ===== -->
            <header class="head">
              <img class="logo" :src="logo" alt="" />
              <div class="title">
                <div class="moul-regular">សាកលវិទ្យាល័យ ខេមរវិទូ</div>
                <div class="en-title">KHEMARAK UNIVERSITY</div>
                <div class="moul-regular branch">ខេត្តបាត់ដំបង</div>
                <div class="moul-regular small">បង្កាន់ដៃទទួលប្រាក់</div>
                <div class="receipt-en">Receipt</div>
              </div>
              <div class="meta box">
                <div><span>No:</span><b class="red">{{ d.reference }}</b></div>
                <div><span>Date:</span><b class="red">{{ d.date }}</b></div>
              </div>
            </header>

            <!-- ===== student + payment info ===== -->
            <div class="info">
              <div class="box student">
                <div class="row">
                  <label>ឈ្មោះ</label>
                  <span class="line grow">{{ d.name_en || d.name_kh }}</span>
                  <label>លេខសម្គាល់</label>
                  <span class="line w-code">{{ d.code }}</span>
                </div>
                <div class="row">
                  <label>កម្រិតសិក្សា</label>
                  <span class="line grow">{{ d.programme_name }}</span>
                  <label>ជំនាន់ទី</label>
                  <span class="line w-sm">{{ d.generation_name }}</span>
                </div>
                <div class="row">
                  <label>ជំនាញ</label>
                  <span class="line grow">{{ d.major_name }}</span>
                  <label>ក្រុម</label>
                  <span class="line grow">{{ d.term }}</span>
                </div>
              </div>

              <div class="box payment">
                <div class="row">
                  <label>ចំនួនទឹកប្រាក់</label>
                  <span class="line grow">{{ money }}</span>
                </div>
                <div class="row">
                  <span class="check-item">
                    <i class="radio" ></i>
                    <b>ទូទាត់តាមរយៈ {{ d.method }}</b>
                  </span>
                  <label class="ref">Payment:</label>
                  <span class="">{{ d.transacntion_code ?? "" }}</span>
                </div>
                <div class="row">
                  <label>បង់</label>
                  <span class="check-item">
                    <i class="radio on sq"></i><label>ឆ្នាំទី</label>
                  </span>
                  <span class="line w-xs">{{ d.year ?? "" }}</span>
                  <span class="check-item">
                    <i class="radio on sq"></i><label>បង់លើកទី</label>
                  </span>
                  <span class="line w-xs">{{ d.sequence_no }}</span>
                </div>
              </div>
            </div>

            <!-- ===== signatures ===== -->
            <div class="signs">
              <div class="sign">
                <div class="khmer">អ្នកបង់ប្រាក់</div>
                <div class="en">Paid By</div>
                <div class="box sign-box">
                  <span class="signer">{{ d.name_en || d.name_kh }}</span>
                </div>
              </div>
              <div class="sign">
                <div class="khmer">អ្នកទទួល</div>
                <div class="en">Receiver</div>
                <div class="box sign-box">
                  <span class="signer">{{ receiver }}</span>
                </div>
              </div>
              <div class="sign">
                <div class="khmer">ហរញ្ញិកៈ</div>
                <div class="en">Treasure</div>
                <div class="box sign-box"></div>
              </div>
            </div>

            <!-- ===== footer note ===== -->
            <footer class="box note">
              <div>
                <b>អាសយដ្ឋាន :</b> ភូមិកម្មករ សង្កាត់ស្វាយប៉ោ ក្រុងបាត់ដំបង
                ខេត្តបាត់ដំបង ទូរស័ព្ទលេខ ០១២ ៨២៥ ២៥៦ / ០៥៣ ៥៤២ ៣០៦
              </div>
              <div class="red">
                <b>បញ្ជាក់ :</b> ទឹកប្រាក់ដែលបានបង់រួច មិនអាចដកវិញ
                ឬបន្ទិលនៅឱ្យអ្នកផ្សេងបានទេ ។ សូមរក្សាបង្កាន់ដៃនេះឱ្យបានល្អ ។
              </div>
            </footer>
          </section>

          <div v-if="copy === 1" class="cut" aria-hidden="true"></div>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<style>
/* Not scoped: must affect <body> children while printing */
@page {
  size: A4 portrait;
  margin: 0;
}
.invoice-print-root {
  display: none;
}
@media print {
  html,
  body {
    margin: 0 !important;
    padding: 0 !important;
    background: #fff !important;
  }
  body > *:not(.invoice-print-root) {
    display: none !important;
  }
  .invoice-print-root {
    display: block !important;
  }
}
</style>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Moul&display=swap');
.sheet {
  width: 210mm;
  height: 297mm;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  background: #fff;
  color: #000000;
  font-size: 10pt;
  -webkit-print-color-adjust: exact;
  print-color-adjust: exact;
  page-break-after: avoid;
  overflow: hidden;
}

/* each receipt = exactly half of A4 (cut line takes 0) */
.receipt {
  flex: 1 1 0;
  min-height: 0;
  box-sizing: border-box;
  padding: 7mm 9mm 6mm;
  display: flex;
  flex-direction: column;
  gap: 4mm;
}
.cut {
  flex: 0 0 0;
  border-top: 0.3mm dashed #777;
  position: relative;
}

.box {
  border: 0.1mm solid #888888;
  border-radius: 2.5mm;
  padding: 2mm 3mm;
}
.red {
  color: #e00000;
}

/* ---------- header ---------- */
.head {
  display: grid;
  grid-template-columns: 30mm 1fr 42mm;
  align-items: start;
  gap: 4mm;
}
.logo {
  width: 30mm;
  height: 30mm;
  object-fit: contain;
}
.title {
  text-align: center;
  line-height: 1.35;
}
.moul-regular {
  font-family: "Moul", serif;
  font-weight: 400;
  color: #1f2f86;
  font-size: 14pt;
}
.moul-regular.branch {
  color: #000;
  font-size: 10pt;
}
.moul-regular.small {
  color: #000;
  font-size: 8pt;
  margin-top: 1mm;
}
.en-title {
  font-family: Arial, sans-serif;
  font-weight: 700;
  color: #1f2f86;
  font-size: 10pt;
  letter-spacing: 0.2pt;
}
.receipt-en {
  font-family: Arial, sans-serif;
  font-size: 9pt;
}
.meta {
  font-family: Arial, sans-serif;
  font-size: 8.5pt;
  line-height: 1.7;
  margin-top: 2mm;
}
.meta span {
  display: inline-block;
  width: 11mm;
}

/* ---------- info boxes ---------- */
.info {
  display: grid;
  grid-template-columns: 1fr 1.1fr;
  gap: 4mm;
}
.row {
  display: flex;
  align-items: flex-end;
  gap: 2mm;
  margin: 2.4mm 0;
  font-size: 9pt;
}
.row label {
  font-weight: 700;
  white-space: nowrap;
}
.line {
  border-bottom: 0.25mm solid #444;
  min-height: 5mm;
  padding: 0 1mm;
  font-family: Arial, "Noto Sans Khmer", sans-serif;
  font-size: 8.5pt;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.line.grow {
  flex: 1 1 0;
  min-width: 0;
}
.line.words {
  text-transform: capitalize;
}
.w-code {
  width: 26mm;
}
.w-sm {
  width: 14mm;
}
.w-xs {
  width: 11mm;
  text-align: center;
}

/* ---------- checks ---------- */
.check-item {
  display: inline-flex;
  align-items: center;
  gap: 1.2mm;
  font-family: Arial, sans-serif;
}
.radio {
  width: 4.2mm;
  height: 4.2mm;
  border: 0.3mm solid #888;
  border-radius: 50%;
  box-sizing: border-box;
  position: relative;
  flex: none;
}
.radio.sq {
  border-radius: 50%;
}
.radio.on {
  background: #1d4ed8;
  border-color: #1d4ed8;
}
.radio.on::after {
  content: "";
  position: absolute;
  left: 1.3mm;
  top: 0.55mm;
  width: 1mm;
  height: 2mm;
  border: solid #fff;
  border-width: 0 0.4mm 0.4mm 0;
  transform: rotate(45deg);
}
.ref {
  margin-left: 2mm;
}

/* ---------- signatures ---------- */
.signs {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6mm;
  text-align: center;
}
.sign .khmer {
  font-weight: 700;
  font-size: 8pt;
}
.sign .en {
  font-family: Arial, sans-serif;
  font-weight: 700;
  font-size: 8.5pt;
  margin-bottom: 1.5mm;
}
.sign-box {
  height: 20mm;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 3mm;
}
.signer {
  font-family: Arial, "Noto Sans Khmer", sans-serif;
  font-size: 8.5pt;
}

/* ---------- note ---------- */
.note {
  font-size: 8.5pt;
  line-height: 1.7;
  margin-top: auto;
}
</style>