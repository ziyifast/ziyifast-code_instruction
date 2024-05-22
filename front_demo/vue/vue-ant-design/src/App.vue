<template>
  <a-layout>
    <a-layout-sider
        breakpoint="lg"
        collapsed-width="0"
        @collapse="onCollapse"
        @breakpoint="onBreakpoint"
        id="components-layout-demo-responsive"
        style="height: 100vh"
    >
      <div class="logo"/>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline">
        <a-menu-item key="/product">
          <user-outlined/>
          <span class="nav-text">商品列表</span>
        </a-menu-item>
        <a-menu-item key="/userInfo">
          <video-camera-outlined/>
          <span class="nav-text">用户中心</span>
        </a-menu-item>
        <a-menu-item key="3">
          <upload-outlined/>
          <span class="nav-text">nav 3</span>
        </a-menu-item>
        <a-menu-item key="4">
          <user-outlined/>
          <span class="nav-text">nav 4</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header :style="{ background: '#fff', padding: 0 }"/>
      <a-layout-content :style="{ margin: '24px 16px 0' }">
        <router-view/>
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        Ant Design ©2018 Created by Ant UED
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>
<script lang="ts" setup>
import {ref, watch} from 'vue';
import {useRouter} from "vue-router";
import {UploadOutlined, UserOutlined, VideoCameraOutlined} from '@ant-design/icons-vue';

const router = useRouter();
const onCollapse = (collapsed: boolean, type: string) => {
  console.log(collapsed, type);
};

const onBreakpoint = (broken: boolean) => {
  console.log(broken);
};

const selectedKeys = ref<string[]>(['4']);

//根据选中model条件不同路由页面
watch(selectedKeys, (newValue) => {
  console.log(newValue[0])
  router.push({path: newValue[0]})
})
</script>
<style scoped>
#components-layout-demo-responsive .logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.2);
  margin: 16px;
}

.site-layout-sub-header-background {
  background: #fff;
}

.site-layout-background {
  background: #fff;
}

[data-theme='dark'] .site-layout-sub-header-background {
  background: #141414;
}
</style>
