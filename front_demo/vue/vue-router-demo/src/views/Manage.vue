<template>
  <a-layout style="height: 100vh; width: 100%;" id="components-layout-demo-side">
    <a-layout-sider v-model:collapsed="collapsed" collapsible style="height: 100%;">
      <div class="logo" style="height: 100%; margin: 0 auto;"/>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline">
        <!-- ... Menu items ... -->
      </a-menu>
    </a-layout-sider>
    <a-layout style="height: 100%;">
      <a-layout-header style="background: #fff; padding: 0; height: 64px;"/>
      <a-layout-content style="margin: 0; padding: 24px 16px; background: #fff; min-height: 0;">
        <a-breadcrumb style="margin: 0;">
          <a-breadcrumb-item>User</a-breadcrumb-item>
          <a-breadcrumb-item>Bill</a-breadcrumb-item>
        </a-breadcrumb>
        <div :style="{ padding: '24px', background: '#fff', minHeight: 'auto' }">
          Bill is a cat.
        </div>
      </a-layout-content>
      <a-layout-footer style="text-align: center; height: 64px;">
        Ant Design ©2018 Created by Ant UED
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>
<script setup>
import {ref} from 'vue';

const collapsed = ref(false);
const selectedKeys = ref(['1']);
</script>
<style scoped>
#components-layout-demo-side .logo {
  height: 100%;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.3);
}

.site-layout .site-layout-background {
  background: #fff;
}

[data-theme='dark'] .site-layout .site-layout-background {
  background: #141414;
}
</style>
