<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { ref } from 'vue';
  import { ctx } from './ctx.js';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  let emailString = ref("");
  let passwordString = ref("");
  let loginError = ref("");

  let emailStringreg = ref("");
  let passwordStringreg = ref("");
  let registerError = ref("");

  const doLogin = async () => {
    loginError.value = await ctx.Login(emailString.value, passwordString.value);
    if (loginError.value == "") router.push('/');
  }

  const doRegister = async () => {
    let q = await ctx.Register(emailStringreg.value, passwordStringreg.value);
    registerError.value = q;
    if (registerError.value == "") router.push('/');
  }

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
      <button class="btn btn-secondary" @click="doLogin();">Вход</button>
      <div v-if="loginError" class="alert alert-danger" role="alert">{{loginError}}</div>

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
      <button class="btn btn-secondary" @click="doRegister();">Регистрация</button>
      <div v-if="registerError" class="alert alert-danger" role="alert">{{registerError}}</div>

    </div>

  </div>


</template>
