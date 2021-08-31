<template>
  <el-row :gutter="20">
    <el-col :span="12">
      <div class="grid-content bg-purple">
        <el-form :inline="true" :model="formInline">
          <el-form-item label="数据显示">
            <el-select v-model="showValue" placeholder="请选择">
              <el-option
                v-for="item in showOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="Host">
            <el-input
              v-model="formInline.host"
              placeholder="Host"
              clearable
            ></el-input>
          </el-form-item>
          <el-form-item label="Port">
            <el-input
              v-model.number="formInline.port"
              placeholder="Port"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="conn">
              {{ !connok ? "连接" : "断开" }}
            </el-button>
          </el-form-item>
        </el-form>
        <el-input
          type="textarea"
          :rows="16"
          placeholder
          readonly
          v-model="textarea"
        ></el-input>
        <el-row :gutter="1">
          <el-col :span="18">
            <div class="grid-content bg-purple">
              <el-input
                v-model="input"
                placeholder="请输入发送内容"
                clearable
              ></el-input>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="grid-content bg-purple">
              <el-select v-model="inputValue" placeholder="请选择">
                <el-option
                  v-for="item in showOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </div>
          </el-col>
          <el-col :span="2">
            <div class="grid-content bg-purple">
              <el-button type="primary" @click="bton" :disabled="!connok"
                >发送</el-button
              >
            </div>
          </el-col>
        </el-row>
      </div>
    </el-col>
    <el-col :span="12">
      <div class="grid-content bg-purple"></div>
    </el-col>
  </el-row>
</template>

<script>
import { defineComponent, ref } from "vue";
import { sendMsg, setTagMsg } from "../api/socket";
import { str2hex, str2Binary, hex2str, binary2Str } from "../comm_utils/data";
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
        host: "127.0.0.1",
        port: 2252,
      },
      connok: false,
      showOptions: [
        {
          value: "文本",
          label: "文本",
        },
        {
          value: "16进制",
          label: "16进制",
        },
        {
          value: "二进制",
          label: "二进制",
        },
      ],
      showValue: "text",
      inputValue: "text",
    };
  },
  methods: {
    bton: function () {
      if (this.input == "") {
        return;
      }
      sendMsg({
        type: 11002,
        data: this.fmtData2(this.inputValue,this.input),
      });
      this.textarea +=
        "↑" + new Date().Format("[hh:mm:ss]:") + this.input + "\n";
      this.input = "";
    },
    onMsg: function (msg) {
      var showData = "";
      switch (msg.type) {
        case 10001:
          this.connok = true;
          return;
        case 20001:
          this.connok = false;
          showData = msg.data;
          break;
        default:
          showData = this.fmtData(this.showValue, msg.data);
      }
      this.textarea += "↓" + new Date().Format("[hh:mm:ss]:") + showData + "\n";
    },
    fmtData: function (type, data) {
      switch (type) {
        case "text":
          return data;
        case "16进制":
          return str2hex(data);
        case "二进制":
          return str2Binary(data);
      }
    },
      fmtData2: function (type, data) {
      switch (type) {
        case "text":
          return data;
        case "16进制":
          return hex2str(data);
        case "二进制":
          return binary2Str(data);
      }},
    conn() {
      if (this.connok) {
        //断开连接
        sendMsg({
          type: 11004,
        });
        this.connok = false;
      } else {
        //建立连接
        sendMsg({
          type: 11001,
          host: this.formInline.host,
          port: this.formInline.port,
        });
      }
    },
  },
  mounted() {
    setTagMsg("tcp", this.onMsg);
  },
});
</script>

<style scoped>
</style>
