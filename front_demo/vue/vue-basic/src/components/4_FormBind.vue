<script setup>
import {reactive, ref} from "vue";
//定义表单数据
// v-model：在表单输入元素或组件上创建双向绑定
let data = reactive({
  username: 'jackson',
  agree: false,
  hobbies: [],
  sex: '',
  degree: '',
  course: [],
})
</script>

<template>
  <div style="display: flex;">
    <div style="border: 1px solid black;width: 300px">
      <form>
        <h1>表单绑定</h1>
        <p style="background-color: darkgray"><label>姓名(文本框)：</label><input v-model="data.username"/></p>
        <p style="background-color: darkgray"><label>同意协议(checkbox)：</label>
          <input type="checkbox" v-model="data.agree"/>
        </p>
        <p style="background-color: darkgray">
          <label>兴趣(多选框)：</label><br/>
          <label><input type="checkbox" value="足球" v-model="data.hobbies"/>足球</label>
          <label><input type="checkbox" value="篮球" v-model="data.hobbies"/>篮球</label>
          <label><input type="checkbox" value="羽毛球" v-model="data.hobbies"/>羽毛球</label>
          <label><input type="checkbox" value="乒乓球" v-model="data.hobbies"/>乒乓球</label>
        </p>
        <p style="background-color: darkgray">
          <label>性别(单选框)：</label>
          <label><input type="radio" name="sex" value="男" :model="data.sex">男</label>
          <label><input type="radio" name="sex" value="女" :model="data.sex">女</label>
        </p>
        <p style="background-color: darkgray">
          <label>学历(单选下拉列表)：</label>
          <select v-model="data.degree">
            <option disabled value="">选择学历</option>
            <option>小学</option>
            <option>初中</option>
            <option>高中</option>
            <option>大学</option>
          </select>
        </p>
        <p style="background-color: darkgray">
          <label>课程(多选下拉列表)：</label>
          <br/>
          <select multiple v-model="data.course">
            <option disabled value="">选择课程</option>
            <option>语文</option>
            <option>数学</option>
            <option>英语</option>
            <option>道法</option>
          </select>
        </p>
      </form>
    </div>
    <div style="border: 1px solid blue;width: 200px">
      <h1>结果预览</h1>
      <p style="background-color: darkgray"><label>姓名：{{ data.username }}</label></p>
      <p style="background-color: darkgray"><label>同意协议：{{ data.agree }}</label>
      </p>
      <p style="background-color: darkgray">
        <label>兴趣：
          <ul v-for="item in data.hobbies">{{ item }}</ul>
        </label>
      </p>
      <p style="background-color: darkgray">
        <label>性别：{{ data.sex }}</label>
      </p>
      <p style="background-color: darkgray">
        <label>学历：{{ data.degree }}</label>
      </p>
      <p style="background-color: darkgray">
        <label>课程：
          <li v-for="c in data.course">
            {{ c }}
          </li>
        </label>
      </p>
    </div>
  </div>
</template>

<style scoped>

</style>