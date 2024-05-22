import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ProductView from "@/views/ProductView.vue";
import UserInfoView from "@/views/UserInfoView.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/about',
            name: 'about',
            // route level code-splitting
            // this generates a separate chunk (About.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: () => import('../views/AboutView.vue')
        },
        {
            path:'/product',
            name: 'product',
            component:ProductView,
        },
        {
            path:'/userInfo',
            name: 'userInfo',
            component:UserInfoView,
        }
    ]
})

export default router
