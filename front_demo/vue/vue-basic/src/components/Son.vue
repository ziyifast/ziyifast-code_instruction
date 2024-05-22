<script setup>

//子组件接收父组件传过来的属性
let props = defineProps({
  books: Array,
  money: {
    type: Number,
    required: true,
    default: 100
  }
})

//子组件向父组件传值：子组件定义发生的事件
let emits = defineEmits((['buy']));

function buy() {
  console.log('son.....buy')
  emits('buy', -5)
}
</script>

<template>
  <div>books:{{ props.books }}</div>
  <div>money:{{ props.money }}</div>
  <button @click="buy">buy it!</button>
  <slot>哈哈</slot>
</template>

<style scoped>

</style>