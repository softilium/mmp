<script setup lang="ts">

  import { onMounted } from 'vue';
  import { authStore } from './components/authStore.js';

  onMounted(() => {
    authStore.SetAccessToken(localStorage.getItem("accessToken"));
    authStore.SetRefreshToken(localStorage.getItem("refreshToken"));
    authStore.CheckLogged();
  });

</script>

<template>

  <header>
    <nav class="navbar navbar-expand-sm navbar-toggleable-sm navbar-light bg-white border-bottom box-shadow.mb-3">
      <div class="container-fluid">
        <RouterLink class="navbar-brand" to="/">Витрины</RouterLink>
        <span v-if="authStore.userInfo.userName">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/orders">Заказы</RouterLink>&nbsp;&nbsp;
        </span>
        <span v-if="!authStore.userInfo.userName">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/login">Войти</RouterLink>&nbsp;&nbsp;
        </span>
        <span v-if="authStore.userInfo.userName">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/logout">Выйти</RouterLink>&nbsp;&nbsp;
        </span>
        <span v-if="authStore.userInfo.userName">
          <RouterLink class="btn btn-info btn-sm" to="/inc-orders">Заказы для обработки</RouterLink>&nbsp;&nbsp;
        </span>
        <span v-if="authStore.userInfo.admin">
          <RouterLink class="btn btn-info btn-sm" to="/set-roles">admin</RouterLink>
        </span>
        <div class="navbar-collapse collapse d-sm-inline-flex justify-content-between">
          <ul class="navbar-nav flex-grow-1">&nbsp;</ul>
          <span v-if="authStore.userInfo.userName">
            <RouterLink class="btn btn-outline-secondary btn-sm" to="/profile">{{ authStore.userInfo.userName }}</RouterLink>&nbsp;&nbsp;
          </span>
        </div>
      </div>
    </nav>
  </header>

  <div class="container-fluid">
    <RouterView />
  </div>

  <br />
  <br />
  <footer class="border-top footer text-muted">
    <div class="row">
      <div class="col">
        <p>
          &copy; 2024-2025,
        </p>
      </div>
    </div>
  </footer>

</template>
