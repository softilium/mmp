<script setup lang="ts">

  import { ref } from 'vue';
  import { authStore } from './authStore.js';

  const loggedEmail = defineModel();

  let emailString = ref<string>("");
  let passwordString = ref<string>("");

</script>

<template>

  <h1>Login</h1>

  <div>

    User email: <input v-model="emailString" />
    Password: <input v-model="passwordString" />

    <button @click="authStore.Login(emailString, passwordString)">Login</button>

  </div>

</template>
