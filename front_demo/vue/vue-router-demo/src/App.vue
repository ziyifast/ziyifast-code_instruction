<script setup>
import http from "./utils/http.js";
import Manage from "./views/Manage.vue";

const sendGet = () => {
  http.get("/get", {}).then(res => {
    console.log(res.data)
  }).catch(err => {
    console.log(err)
  })
}

const sendGetWithParams = () => {
  http.get("/get", {
    params: {
      'username': 'zhangsan',
      'age': 23
    }
  }).then(res => {
    console.log(res.data)
  })
}

const sendPost = () => {
  http.post("/post", {})
      .then(res => {
        console.log(res)
      })
}

const sendPostWithParams = () => {
  http.post("/post", {
    'username': 'zhangsan',
    'age': 23
  }).then(res => {
    console.log(res.data)
  })
}
</script>

<template>
  <!--  <router-link to="/"></router-link>-->
  <!--  <router-link to="/userinfo/zhangsan/28">用户中心</router-link>-->
  <!--  <router-view/>-->
  <!--  <router-link to="/home">首页</router-link>-->
  <!--  <br/>-->
  <!--  <router-link to="{-->
  <!--  name:'userinfo',-->
  <!--  params:{-->
  <!--    'username':'zhangsan',-->
  <!--    'age':23-->
  <!--  }-->
  <!--  }">-->
  <!--    用户中心 params by obj-->
  <!--  </router-link>-->
  <!--  <br/>-->
  <!--  <router-link to="/userinfo/zhangsan/28">用户中心 params by path</router-link>-->
  <!--  <br/>-->
  <!--  <router-link to="/userinfo/zhangsan/20?id=55">用户中心 query by path</router-link>-->
  <!--  <br/>-->
  <!--  <router-link to="{-->
  <!--      name:'userinfo',-->
  <!--  query:{-->
  <!--    'id':'99'-->
  <!--  }-->
  <!--  }">用户中心 query by obj-->
  <!--  </router-link>-->
  <!--  <router-view/>-->


  <!--  测试请求-->
<!--  <button @click="sendGet">get请求</button>-->
<!--  <button @click="sendGetWithParams">get请求带参数</button>-->
<!--  <button @click="sendPost">post请求</button>-->
<!--  <button @click="sendPostWithParams">post请求带参数</button>-->
  <Manage/>
</template>

<style scoped>
</style>
