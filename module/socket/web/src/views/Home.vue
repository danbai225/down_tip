<template>
  <el-container>
    <el-main>
      <el-tabs type="border-card">
        <el-tab-pane label="TCP">
          <TCP />
        </el-tab-pane>
        <el-tab-pane label="UDP">UDP</el-tab-pane>
        <el-tab-pane label="HTTP">HTTP</el-tab-pane>
        <el-tab-pane label="WebSocket">WebSocket</el-tab-pane>
      </el-tabs>
    </el-main>
  </el-container>
</template>

<script>
import TCP from "../components/TCP.vue";
export default {
  name: "Home",
  components: {
    TCP,
  },
};
</script>
