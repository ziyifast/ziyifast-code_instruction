import {defineStore} from 'pinia'
import {computed, ref} from "vue";

// 方式一：actions写法
// export const useMoneyStore = defineStore('money', {
//     state: ()=>{
//         return {
//             money: 100,
//         }
//     },
//     getters: {
//         rmb: (state) => {
//             return state.money
//         },
//         usd: (state) => {
//             return state.money * 0.14
//         },
//         eur: (state) => {
//             return state.money * 0.13
//         }
//     },
//     actions:{
//         win(arg) {
//             this.money += arg
//         },
//         pay(arg){
//             this.money -= arg
//         }
//     }
// })


//方式二：setup写法（推荐）
export const useMoneyStore = defineStore('money', () => {
    const rmb = ref(100)
    //computed替换getters
    const usd = computed(() => {
        return rmb.value * 0.14
    })
    const eur = computed(() => {
        return rmb.value * 0.12
    })
    //function替换actions
    const pay = (arg) => {
        rmb.value -= arg
    }
    const win = (arg) => {
        rmb.value += arg
    }
    return {rmb, usd, eur, pay, win}
})