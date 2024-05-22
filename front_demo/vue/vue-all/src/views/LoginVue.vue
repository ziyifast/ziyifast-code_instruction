<script>
import {login} from "@/api/user.js";
export default {
  data() {
    return {
      loginUser: {
        username: 'admin',
        password: 'admin',
      },
    };
  },
  methods: {
    userLogin(){
      login(this.loginUser).then(res => {
        console.log(res.status)
        if(res.status === 200){
          alert('登录成功')
          this.$router.push({path: `/userInfo/`})
        }
      }, error => {
        alert('登录失败')
      })
    }


  },
}
</script>

<template>
  <div class="loginContainer">
    <div class="loginHeader">
      <span style="color: #2b92e4">ziyi后台管理系统</span>
    </div>
    <div class="loginMainWrapper">
      <div class="loginWrapper">
        <div class="loginTipsWrapper" style="color: #ffffff">
          <span class="siteSummaryTitle">ziy后台管理系统</span>
          <div class="siteSummary">
            <ul>
              <li>
                <a-icon type="check-square" style="margin-right: 10px"/>
                <span>统一方便的资源管理</span></li>
              <li>
                <a-icon type="check-square" style="margin-right: 10px"/>
                <span>友好的界面展示，基于AntDesign组件库</span></li>
              <li>
                <a-icon type="check-square" style="margin-right: 10px"/>
                <span>前后端分离，便于维护</span></li>
            </ul>
          </div>
        </div>
        <div class="loginBoxWrapper" align="center">
          <a-form-model :model="loginUser" style="width: 100%">
            <a-form-model-item style="color: palevioletred;text-align: center">
              <h1>欢迎登录</h1>
            </a-form-model-item>
            <a-form-model-item>
              <a-input v-model:value="loginUser.username" placeholder="账号">
                <a-icon slot="prefix" type="user" style="color:rgba(0,0,0,.25)"/>
              </a-input>
            </a-form-model-item>
            <a-form-model-item>
              <a-input v-model:value="loginUser.password" type="password" placeholder="密码">
                <a-icon slot="prefix" type="lock" style="color:rgba(0,0,0,.25)"/>
              </a-input>
            </a-form-model-item>
            <a-form-model-item>
            </a-form-model-item>
            <a-form-model-item>
              <a-button type="primary" html-type="submit" style="width: 100%" @submit.prevent
                        @click="userLogin">立即登录
              </a-button>
              <a-button style="width: 100%;margin-top: 10px">注册账号</a-button>
            </a-form-model-item>
            <a-form-model-item style="text-align: center">
              <a-space size="large">
                <a-icon type="weibo" style="font-size: 2em;cursor: pointer;color:#f50"/>
                <a-icon type="qq" style="font-size: 2em;cursor: pointer;color:#2b92e4;"/>
                <a-icon type="github" style="font-size: 2em;cursor: pointer;color:#333;"/>
              </a-space>
            </a-form-model-item>
          </a-form-model>
        </div>
      </div>
      <div class="loginFooter">
        © 2020-2024 CoreCmd 版权所有. 川ICP备18063315号-1 川公网安备 52158202001416号
      </div>
    </div>
  </div>
</template>

<style scoped>
.loginContainer {
  width: 100%;
  position: relative;
  background-attachment: fixed;
  background-repeat: no-repeat !important;
  background-size: cover;
  background-position: center;
  padding-top: 105px;
  min-height: 100%;
}

.loginHeader {
  width: 100%;
  height: 64px;
  background-color: #ffffff;
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  z-index: 1000;
  box-shadow: 0 2px 15px rgba(0, 0, 0, .2);
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-direction: row;
  padding: 10px;
}

.loginHeader span {
  font-size: 25px;
  font-weight: 700;
}

.loginMainWrapper {
  width: 100%;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  align-items: center;
}

.loginWrapper {
  width: 90%;
  display: flex;
  align-items: flex-start;
  justify-content: right;
  flex-direction: row;
  height: 100%;
}

.loginFooter {
  width: 100%;
  min-height: 64px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  flex-direction: row;
  text-align: center;
  color: #ffffff;
  margin-top: 40px;
}

.siteSummary ul {
  list-style: none;
  padding: 0;
}

.siteSummary ul li {
  margin-top: 10px;
  list-style: none;
}

@media screen and (min-width: 1200px) {
  .loginTipsWrapper {
    padding: 20px 30px;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    flex-direction: column;
    min-height: 100%;
  }

  .loginBoxWrapper {
    background-color: #ffffff;
    padding: 20px;
    width: 400px;
    height: 100%;
    border-radius: 5px;
  }
}

@media screen and (min-width: 769px) and (max-width: 1200px) {
  .loginTipsWrapper {
    padding: 20px 30px;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    flex-direction: column;
    height: 100%;
  }

  .loginBoxWrapper {
    background-color: #ffffff;
    padding: 20px;
    width: 400px;
    height: 100%;
    border-radius: 5px;
  }
}

@media screen and (max-width: 768px) {
  .loginTipsWrapper {
    display: none;
  }

  .loginBoxWrapper {
    background-color: #ffffff;
    padding: 20px;
    width: 100%;
    border-radius: 5px;
  }
}

.siteSummaryTitle {
  color: #ffffff;
  font-size: 1.5rem;
  font-weight: 700;
}

.siteSummary {
  margin-top: 25px;
  color: #ffffff;
  font-size: 1.2rem;
  font-weight: 300;
}
</style>

