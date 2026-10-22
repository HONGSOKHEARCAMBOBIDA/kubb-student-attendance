<template>
  <div>
    <div id="print-area">
      <el-row :gutter="0" align="middle" class="doc-header">
        <el-col :span="4">
          <div class="header-logo">
            <el-image :src="url" fit="contain" class="university-logo" />
            <h4 class="label">សាកលវិទ្យាល័យ ខេមរាវិទូ</h4>
            <h4 class="labelsr">KHEMARAK UNIVERSITY</h4>
            <h4 class="label">ខេត្តបាត់ដំបង</h4>
          </div>
        </el-col>

        <el-col :span="15">
          <h3 class="subtitle-kh">ព្រះរាជាណាចក្រកម្ពុជា</h3>
          <h4 class="subtitle-kh">ជាតិ សាសនា ព្រះមហាក្សត្រ</h4>
          <el-image :src="taktieng" class="taktieng-logo" />
          <h3 class="title-kh">{{ title }}</h3>
          <h3 class="title-kh">{{ secondtitle }}</h3>
          <h3 class="title-khsr">{{ thirdtitle }}</h3>
        </el-col>

        <el-col :span="5">
          <h4 class="labelsr">មុខវិជ្ជាៈ {{ subject }}</h4>
          <h4 class="labelsr">កាលបរិច្ឆេទៈ {{ examdate }}</h4>
          <h4 class="labelsr">វេលាម៉ោងៈ {{ examtime }}</h4>
          <h4 class="labelsr">បន្ទប់លេខៈ {{ room }}</h4>
        </el-col>
      </el-row>

      <table>
        <thead>
          <tr>
            <th v-if="showIndex">ល.រ</th>
            <th v-for="col in columns" :key="col.key">
              {{ col.label }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="(row, index) in rows" :key="row.id ?? index">
            <td v-if="showIndex">{{ index + 1 }}</td>
            <td v-for="col in columns" :key="col.key">
              <!-- empty value + strike flag → show a strikethrough line -->
              <el-text
                v-if="col.strikeIfEmpty && isEmpty(row[col.key])"
                type="danger"
                class="strike-empty"
                size="small"
                tag="del"
              >
                គ្មានសិទ្ធប្រឡង
              </el-text>
              <template v-else>
                {{ col.format ? col.format(row[col.key], row) : row[col.key] }}
              </template>
            </td>
          </tr>
        </tbody>
      </table>
      <el-row :gutter="0" align="middle" class="doc-footer">
        <el-col :span="12">
          <h4 class="labelsr">
            បញ្ឃប់បញ្ជីឈ្មោះត្រឹម {{ totalstudent }}នាក់ ស្រីចំនួន
            {{ totalstudentgirld }} នាក់
          </h4>
          <h4 class="labelsr">វត្តមានចំនួន........នាក់ ស្រី.......នាក់</h4>
          <h4 class="labelsr">អវត្តមានចំនួន.......នាក់ ស្រី.......នាក់</h4>
        </el-col>

        <el-col :span="12">
          <div class="">
            <h4 class="right">
              បាត់ដំបង ថ្ងៃទី...................ខែ.................ឆ្នាំ202
            </h4>
            <h4 class="rightmoul">ហត្ថាលេខា និងឈ្មោះអនុរក្ស</h4>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
const isEmpty = (v) => v === null || v === undefined || v === "" || v === 0;
const url = "/logo.png";
const taktieng = "/image.png";
defineProps({
  title: {
    type: String,
    default: "",
  },
  secondtitle: {
    type: String,
    default: "",
  },
  thirdtitle: {
    type: String,
    default: "",
  },
  columns: {
    type: Array,
    required: true,
  },
  rows: {
    type: Array,
    default: () => [],
  },
  showIndex: {
    type: Boolean,
    default: true,
  },
  subject: {
    type: String,
    default: "",
  },
  examdate: {
    type: String,
    default: "",
  },
  examtime: {
    type: String,
    default: "",
  },
  room: {
    type: String,
    default: "",
  },
  totalstudent: {
    type: Number,
    default: 0,
  },
  totalstudentgirld: {
    type: Number,
    default: 0,
  },
});

function print() {
  window.print();
}

defineExpose({ print });
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Moul&display=swap");

.doc-header {
  margin-bottom: 6px;
}

.doc-footer {
  margin-top: 15px;
}

.header-logo {
  padding-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.university-logo {
  width: 70px;
  height: 70px;
}

.taktieng-logo {
  width: 110px;
  height: 25px;
  display: block;
  margin: 4px auto;
}

#print-area {
  padding: 1px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  font-size: 12px;
  text-align: center;
  border: 1px solid #000;
  padding: 5px;
}

th {
  font-weight: 200;
  font-family: "Moul", serif;
  font-size: 12px;
  text-align: center;
}

.title-kh,
.subtitle-kh,
.title-doc {
  font-family: "Moul", serif;
  font-weight: 200;
  text-align: center;
  margin: 2px 0;
}

.title-khsr {
  font-weight: 200;
  text-align: center;
  margin: 2px 0;
  font-size: 12px;
}

.title-kh {
  font-size: 12px;
}

.subtitle-kh {
  font-size: 14px;
}

.label {
  font-family: "Moul", serif;
  font-weight: 200;
  text-align: center;
  font-size: 8px;
}

.labelsr {
  font-weight: 200;
  text-align: left;
  font-size: 10px;
}

.labelsrright {
  font-weight: 200;
  text-align: center;
  font-size: 10px;
}

.right {
  font-size: 10px;
  padding-top: 15px;
  padding-right: 10px;
  text-align: right;
}

.rightmoul {
  font-family: "Moul", serif;
  font-weight: 200;
  font-size: 10px;
  padding-right: 30px;
  text-align: right;
}
.subtitle-en {
  font-family: Arial, sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-align: center;
  margin: 2px 0;
  letter-spacing: 0.5px;
}

.title-doc {
  font-size: 18px;
  margin-top: 8px;
}
</style>

<style>
/* Not scoped — needs to affect body and elements outside this component */
@media print {

  @page {
    size: A4 portrait;
    margin: 10mm;
  }

  body * {
    background-color: transparent;
    visibility: hidden;
  }

  #print-area,
  #print-area * {
    visibility: visible;
  }

  #print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    padding: 0;
  }

  .print-btn {
    display: none;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    border: 1px solid #000;
    padding: 5px;
  }
}
</style>
