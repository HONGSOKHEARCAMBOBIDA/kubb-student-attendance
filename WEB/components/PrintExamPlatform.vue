<template>
  <div id="print-exam-area">
    <el-row :gutter="0" class="doc-header">
      <el-col :span="8">
        <div class="header-logo">
          <img :src="url" fit="contain" class="university-logo" />
          <h4 class="label">សាកលវិទ្យាល័យ ខេមរាវិទូ</h4>
          <h4 class="labelsr">KHEMARAK UNIVERSITY</h4>
          <h4 class="label">សាខាខេត្តបាត់ដំបង</h4>
          <h4 class="labelsr">លេខៈ..........សខ.</h4>
        </div>
      </el-col>
      <el-col :span="8">
        <h3 class="moul big">ព្រះរាជាណាចក្រកម្ពុជា</h3>
        <h4 class="moul big">ជាតិ សាសនា ព្រះមហាក្សត្រ</h4>
        <img :src="taktieng" class="taktieng-logo" alt=""/>
      </el-col>
      <el-col :span="8">

      </el-col>
    </el-row>

    <div class="titles">
      <h3 class="moul">{{ title }}</h3>
      <h3 class="moul">{{ secondtitle }}</h3>
      <h3 class="moul">{{ thirdtitle }}</h3>
    </div>

    <table>
      <thead>
        <tr>
          <th style="width: 40px">ល.រ</th>
          <th>មុខវិជ្ជាប្រឡង</th>
          <th>កាលបរិច្ឆេទប្រឡង</th>
          <th>ពេលវេលាប្រឡង</th>
          <th>ឈ្មោះសាស្ត្រាចារ្យបង្រៀន</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, i) in rows" :key="i">
          <td>{{ toKh(i + 1) }}</td>
          <td class="left">{{ r.subject }}</td>
          <td>{{ r.date }}</td>
          <td>{{ r.time }}</td>
          <td class="left">{{ r.teacher }}</td>
        </tr>
      </tbody>
    </table>

    <div class="notes">
      <u class="moul-sm">បញ្ជាក់៖</u>
      <ul>
        <li  v-for="(n, i) in notes" :key="i">{{ n }}</li>
      </ul>
    </div>

    <div class="signs">
      <div class="sign-left">
        <div class="space"></div>
        <p class="">បានឃើញ និងឯកភាព</p>
        <p class="moul-sm">ជ.នាយិកាប្រតិបត្តិ</p>
        <p class="moul-sm">នាយករង</p>
      </div>
      <div class="sign-right">
        <p>បាត់ដំបង {{ issuedDate }}</p>
        <p class="moul-sm">ប្រធានការិយាល័យសិក្សា</p>
        <div class="space"></div>
        <p class="moul-sm">{{ headOfStudy }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const url = "/logo.png";
const taktieng = "/image.png";
const toKh = (v) => String(v ?? "").replace(/\d/g, (d) => "០១២៣៤៥៦៧៨៩"[d]);

defineProps({
  title: { type: String, default: "" },
  secondtitle: { type: String, default: "" },
  thirdtitle: { type: String, default: "" },
  rows: { type: Array, default: () => [] },
  notes: { type: Array, default: () => [] },
  issuedDate: { type: String, default: "" },
  approver: { type: String, default: "" },
  headOfStudy: { type: String, default: "" },
});

async function print() {
  const images = document.querySelectorAll("#print-exam-area img");
  for (const img of images) {
    if (!img.complete) {
      await new Promise(resolve => {
        img.onload = resolve;
        img.onerror = resolve;
      });
    }
  }
  await document.fonts.ready;
  window.print();
}

defineExpose({ print });
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Moul&display=swap");

#print-exam-area { padding: 1px; font-size: 13px; }
.moul { font-family: "Moul", serif; font-weight: 200; text-align: center; margin: 3px 0; font-size: 18px; }
.moul.big { font-size: 22px; }
.moul-sm { font-family: "Moul", serif; font-weight: 200; font-size: 15px; margin: 2px 0; }
.label { font-family: "Moul", serif; font-weight: 200; font-size: 12px; margin: 2px 0; }
.labelsr { font-weight: 600; font-size: 10px; margin: 2px 0; }
.header-logo { display: flex; flex-direction: column; align-items: center; text-align: center; padding-top: 40px;}
.university-logo { width: 100px; height: 100px; }
.taktieng-logo { width: 110px; height: 25px; display: block; margin: 6px auto; }
.titles { margin: 14px 0 10px; }
table { width: 100%; border-collapse: collapse; }
th, td { border: 1px solid #000; padding: 6px; text-align: center; font-size: 18px; }
th { font-family: "Moul", serif; font-weight: 200; font-size: 16px; }
td.left { text-align: left; }
.notes { margin-top: 12px; }
.notes ul { list-style: none; padding-left: 20px; margin: 6px 0; }
.notes li::before { content: "- "; }
.notes li { margin: 4px 0; font-size: 17px; }
.signs { display: flex; justify-content: space-between; margin-top: 20px; text-align: center; }
.signs p { margin: 2px 0; }
.space { height: 80px; }
</style>

<style>
@media print {
  @page { size: A4 portrait; margin: 10mm; }
  body * { background: transparent; visibility: hidden; }
  #print-exam-area, #print-exam-area * { visibility: visible; }
  #print-exam-area { position: absolute; left: 0; top: 0; width: 100%; }
}
</style>