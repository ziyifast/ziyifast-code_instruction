<script setup>
//基本类型
let msg = "hello"
//对象类型
let player = {
  name: "jack",
  age: 18
}
</script>

<template>
  <div>
    <!--   插值语法-->
    {{ msg }}
    <p>玩家姓名:{{ player.name }}</p>
    <p>玩家年龄:{{ player.age }}</p>
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
