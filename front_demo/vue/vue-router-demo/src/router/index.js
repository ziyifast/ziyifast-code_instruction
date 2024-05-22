import Home from '../views/Home.vue'
import {createMemoryHistory, createRouter} from "vue-router";
import UserInfo from "../views/UserInfo.vue";
import UserEmail from "../views/UserEmail.vue";
import UserProfile from "../views/UserProfile.vue";
import Register from "../views/Register.vue";
import Manage from "../views/Manage.vue";

const routes = [
    {
        path: "/",
        component: Register,
    },
    {
        path: '/home',
        component: Home
    },
    {
        path: '/userinfo/:username/:age',
        component: UserInfo,
        name: 'userinfo',
        children: [
            {
                path: 'email',
                component: UserEmail
            },
            {
                path: 'profile',
                component: UserProfile
            }
        ]
    },
    {
        path: '/manage',
        component: Manage,
        name: 'manage'
    }
]

//注册路由
const router = createRouter({
    routes: routes,
    history: createMemoryHistory()
})

//添加路由守卫，跳转之前判断
router.beforeEach(async (to, from) => {
    // console.log('beforeEach', to, from)
    if (from.fullPath == '/' && to.fullPath == '/home') {
        console.log('params', from.params)
    }
})
export default router;