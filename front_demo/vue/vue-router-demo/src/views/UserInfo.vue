<script setup>
import {useRoute, useRouter} from 'vue-router'
import axios from "axios";
import instance from '../utils/http.js'

const route = useRoute()
const router = useRouter()
//回到上一页
const back = () => {
  //前进-1=》后退一步
  router.go(-1)
}

const home = () => {
  router.push({path: "/home", query: {username: "张三", age: 18}})
}

const getUserInfo = () => {
  instance.get('/user/info').then(res => {
    console.log(res)
  })
}
</script>

<template>
  <h1>用户中心</h1>
  <button @click="getUserInfo">获取用户信息</button>
  <!--  获取到的数据：-->
  <!--  <div>-->
  <!--    params: {{ route.params }}-->
  <!--    query: {{ route.query }}-->
  <!--  </div>-->
  <!--  <button @click="back">回到上一页</button>-->
  <!--  <button v-on:click="home">回到首页</button>-->
  <!--  <hr/>-->
  <!--  <router-link to="/userinfo/333/profile">个人头像</router-link>-->
  <!--  <br/>-->
  <!--  <router-link to="/userinfo/333/email">个人邮件</router-link>-->
  <!--  <router-view/>-->
</template>

<style scoped>

</style>