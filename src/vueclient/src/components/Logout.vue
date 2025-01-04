<script setup>

  import { authStore } from './authStore.js';

</script>

<template>

  <h3>Выйти из системы?</h3>
  <br />

  <button class="btn btn-primary" @click="authStore.Logout(); $router.push('/');">Выйти</button>

</template>
