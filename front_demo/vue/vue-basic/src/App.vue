<script setup>
import Father from "./components/Father.vue";
// import WatchNum from './components/6_WatchNum.vue'
</script>

<template>
  <div>
    <!--    插值语法-->
    <!--    <InsertGrammer/>-->

    <!--    指令用法-->
    <!--    <CommandDemo/>-->

    <!--绑定响应式对象-->
    <!--    <AttributeRefBond/>-->

    <!--    表单绑定-->
    <!--    <FormBind/>-->

    <!--    计算属性-->
    <!--    <ComputedAttribute/>-->

    <!--    watch 监听属性-->
    <!--    <WatchNum/>-->

<!--    钩子函数：Vue生命周期函数-->
<!--    <HookFunc/>-->

<!--    子父组件传值-->
<!--<Father/>-->

    <Father/>

  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
