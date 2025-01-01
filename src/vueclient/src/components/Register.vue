<script setup lang="ts" >

  import { ref } from 'vue';
  import { authStore } from './authStore.js';

  let emailString = ref<string>("");
  let passwordString = ref<string>("");

</script>

<template>

  <h1>Register</h1>

  <div>

    User email: <input v-model="emailString" />
    Password: <input v-model="passwordString" />

    <button @click="authStore.Register(emailString, passwordString); $router.push('/');">Register</button>

  </div>


</template>
