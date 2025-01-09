<script setup lang="ts">

  import { ref } from 'vue';
  import { authStore } from './authStore.js';

  let emailString = ref("");
  let passwordString = ref("");

  let emailStringreg = ref("");
  let passwordStringreg = ref("");

</script>

<template>

  <div class="row">

    <div class="col-5">

      <h1>Вход</h1>

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
      <button class="btn btn-primary" @click="authStore.Login(emailString, passwordString); $router.push('/');">Вход</button>

    </div>

    <div class="col-2"></div>

    <div class="col-5">
      <h1>Регистрация</h1>

      <div class="row mb-3">
        <label class="col-3 form-label">Email</label>
        <div class="col-7">
          <input class="form-control" v-model="emailStringreg" required />
        </div>
      </div>
      <div class="row mb-3">
        <label class="col-3 form-label">Пароль</label>
        <div class="col-7">
          <input class="form-control" v-model="passwordStringreg" type="password" />
        </div>
      </div>
      <button class="btn btn-primary" @click="authStore.Register(emailStringreg, passwordStringreg); $router.push('/');">Регистрация</button>

    </div>

  </div>


</template>
