import {createRouter, createWebHistory} from 'vue-router'
import UserInfo from "@/views/UserInfo.vue";
import LoginVue from "@/views/LoginVue.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/', //默认展示登录页面
            name: 'login',
            component: LoginVue,
        },
        {
            path: '/userInfo',
            name: 'userInfo',
            component: UserInfo,
        }
    ]
})

export default router
