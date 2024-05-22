<script setup>
//属性绑定响应式对象
/*
数据的动态变化需要反馈到页面；
Vue通过ref()和reactive()包装数据，将会生成一个数据的代理对象。vue内部的 基于依赖追踪的响应式系统 就会追踪感知数据变化，并触发页面的重新渲染。

使用步骤：
1. 使用 ref() 包装原始类型、对象类型数据，生成 代理对象
2. 任何方法、js代码中，使用 代理对象.value 的形式读取和修改值
3. 页面组件中，直接使用 代理对象
注意：推荐使用 const（常量） 声明代理对象。代表代理对象不可变，但是内部值变化会被追踪。

响应式reactive:
使用步骤：
1. 使用 reactive() 包装对象类型数据，生成 代理对象
2. 任何方法、js代码中，使用 代理对象.属性的形式读取和修改值
3. 页面组件中，直接使用 代理对象.属性

总结：
- 基本类型用ref()、对象类型用reactive()
- ref获取值需.value.属性，reactive直接.属性
- ref()可以将所有类型都变为响应式
也可以 ref 一把梭，大不了 天天 .value
 */

import {reactive, ref} from "vue";

//ref将基本数据类型变为响应式，js中获取ref的响应式对象需要.value
let count = ref(0)
let add = () => {
  count.value++
}

let player = reactive({
  name: 'curry',
  age: 18,
  hobbies: ['篮球', '吃爆米花']
})
let getOlder = () => {
  player.age++
}
</script>

<template>
  <button @click="add">点击升职加薪</button>
  <p style="color: red">月薪：{{ count }} K</p>
  <div>
    球员信息：
    <p>姓名：{{ player.name }}</p>
    <p>年龄：{{ player.age }}</p>
    爱好：
    <p v-for="(item, i) in player.hobbies">
      <li>{{ item }}</li>
    </p>
    <button v-on:click="getOlder">让小学生长大</button>
  </div>
</template>

<style scoped>

</style>