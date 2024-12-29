<script setup>

  import { authStore } from './authStore.js';

</script>

<template>

  <button @click="authStore.Logout(); $router.push('/');">Logout</button>

</template>
