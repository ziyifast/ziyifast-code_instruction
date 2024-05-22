<script setup>
/*
指令：
  v-on:事件绑定
  v-if：条件判断
  v-for：循环
  更多指令https://cn.vuejs.org/api/built-in-directives.html；
 */
let buy = () => {
  alert("购买成功")
}
let price = 200
let arr = ['广东宏远', '辽宁本钢', '成都蓉城', '山东泰山']
</script>

<template>
  <div>
    <!--    v-on：绑定点击事件-->
    <button v-on:click="buy">点我购买</button>
    <!--    v-if -->
    <div v-if="price> 100" style="color: red">上升期</div>
    <div v-else style="color: green">冒绿光</div>
    <!--    v-for -->
    <div v-for="(team,index)  in arr">
      <ul>球队：{{ team }} 序号：{{ index }}</ul>
    </div>
  </div>
</template>

<style scoped>

</style>