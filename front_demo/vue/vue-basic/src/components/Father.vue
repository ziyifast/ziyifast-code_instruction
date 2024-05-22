<!-- 父组件 -->
<script setup>
import Son from "./Son.vue";
import {ref} from "vue";

let data = ref({
  books: ['哈利波特', '三体'],
  money: 500
})

let moneyMinus = (v) => {
  console.log('father....value=', v)
  data.value.money += v
  console.log(data.value.money)
}
</script>

<template>
  <Son :books="data.books" :money="data.money" @buy="moneyMinus">
<!--    <button>按钮</button>-->
  </Son>
</template>

<style scoped>
</style>
