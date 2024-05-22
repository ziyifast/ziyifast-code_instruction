<script setup>
//computed：计算属性=》总价=数量*单价
import {computed, ref} from "vue";

let count = ref(0)
let price = ref(10)
let total = computed(() => {
  return count.value * price.value
})

let addCount = () => {
  count.value++
}
let addPrice = () => {
  price.value++
}

</script>

<template>
  <div>
    <p>单价:{{ price }}</p>
    <p>数量:{{ count }}</p>
  </div>
  <button @click="addCount">加量</button>
  <button @click="addPrice">加价</button>
  <p>总花费:{{ total }}</p>

</template>

<style scoped>

</style>