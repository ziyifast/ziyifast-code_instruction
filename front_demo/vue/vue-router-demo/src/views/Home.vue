<script setup>
import {useRoute} from "vue-router";

const route = useRoute()
</script>

<template>
  <h1>Home首页 params:{{ route.params }} query:{{ route.query }}</h1>
</template>

<style scoped>

</style>