<script setup lang="ts">

  import { ref } from 'vue';
  import { authStore } from './authStore.js';

  let emailString = ref<string>("");
  let passwordString = ref<string>("");

</script>

<template>

  <h1>Регистрация</h1>

  <div class="row mb-3">
    <label class="col-3 form-label">Email</label>
    <div class="col-7">
      <input class="form-control" v-model="emailString" required />
    </div>
  </div>
  <div class="row mb-3">
    <label class="col-3 form-label">Пароль</label>
    <div class="col-7">
      <input class="form-control" v-model="passwordString" type="password" />
    </div>
  </div>
  <button class="btn btn-primary" @click="authStore.Register(emailString, passwordString); $router.push('/');">Регистрация</button>

</template>
