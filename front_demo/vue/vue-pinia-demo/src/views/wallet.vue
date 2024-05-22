<script setup>
import {useMoneyStore} from "../stores/money.js";

const moneyStore = useMoneyStore()
const addMoney = () => {
  moneyStore.win(100)
}

const reduceMoney = () => {
  moneyStore.pay(50)
}

</script>

<template>
  <div>
    rmb: {{ moneyStore.rmb }} <br/>
    usd: {{ moneyStore.usd }} <br/>
    eur: {{ moneyStore.eur }} <br/>
    <button @click="addMoney">add</button>
    <button @click="reduceMoney">reduce</button>
  </div>
</template>

<style scoped>
div {
  padding: 1.5em;
  height: 100%;
  width: 100%;
  background-color: #213547;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>