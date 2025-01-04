<script setup lang="ts">

  import { ref } from 'vue';
  import { authStore } from './authStore.js';

  const loggedEmail = defineModel();

  let emailString = ref<string>("");
  let passwordString = ref<string>("");

</script>

<template>

  <h1>Войти</h1>

  <div class="row">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Email</label>
      <div class="col-7">
        <input class="form-control" v-model="emailString" required />
      </div>
    </div>
    <div>&nbsp;</div>
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Пароль</label>
      <div class="col-7">
        <input class="form-control" v-model="passwordString" type="password" />
      </div>
    </div>
  </div>
  <div>&nbsp;</div>
  <button class="btn btn-primary" @click="authStore.Login(emailString, passwordString); $router.push('/');">Вход</button>

</template>
