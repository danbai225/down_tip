<template>
  <el-row :gutter="20">
    <el-col :span="12"
      ><div class="grid-content bg-purple">
        <el-form :inline="true" :model="formInline" class="demo-form-inline">
          <el-form-item label="Host">
            <el-input v-model="formInline.host" placeholder="Host"></el-input>
          </el-form-item>
          <el-form-item label="Port">
            <el-input
              v-model.number="formInline.port"
              placeholder="Port"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="conn">{{
              !connok ? "连接" : "断开"
            }}</el-button>
          </el-form-item>
        </el-form>
        <el-input type="textarea" :rows="16" placeholder="" v-model="textarea">
        </el-input>
        <el-row :gutter="1">
          <el-col :span="22"
            ><div class="grid-content bg-purple">
              <el-input
                v-model="input"
                placeholder="请输入发送内容"
              ></el-input></div
          ></el-col>
          <el-col :span="2"
            ><div class="grid-content bg-purple">
              <el-button type="primary" @click="bton" :disabled="!connok"
                >发送</el-button
              >
            </div></el-col
          >
        </el-row>
      </div></el-col
    >
    <el-col :span="12"><div class="grid-content bg-purple"></div></el-col>
  </el-row>
</template>

<script>
import { defineComponent, ref } from "vue";
import { sendMsg, setTagMsg } from "../api/socket";
export default defineComponent({
  name: "TCP",
  setup() {
    return {
      input: ref(""),
      textarea: ref(""),
    };
  },
  data() {
    return {
      formInline: {
        host: "",
        port: 0,
      },
      connok: false,
    };
  },
  methods: {
    bton: function () {
      sendMsg({
        type: 10002,
        data: this.input,
      });
      this.textarea += "↑" + this.input + "\n";
      this.input = "";
    },
    onMsg: function (msg) {
      switch (msg.type) {
        case 10004:
          this.connok = true;
          return;
      }
      this.textarea += "↓" + msg.data + "\n";
    },
    conn() {
      sendMsg({
        type: 10001,
        host: this.formInline.host,
        port: this.formInline.port,
      });
    },
  },
  mounted() {
    setTagMsg("tcp", this.onMsg);
  },
});
</script>

<style scoped></style>
