<script setup lang="ts">

  import { onMounted } from 'vue';
  import { authStore } from './components/authStore.js';
  import { useRoute } from 'vue-router'

  const route = useRoute();

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
        <RouterLink class="btn btn-outline-secondary btn-sm" to="/">Витрины</RouterLink>&nbsp;&nbsp;
        <span v-if="authStore.userInfo.id">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/orders">Заказы</RouterLink>&nbsp;
        </span>
        <span v-if="!authStore.userInfo.id">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/login">Войти</RouterLink>&nbsp;
        </span>
        <span v-if="authStore.userInfo.id">
          <RouterLink class="btn btn-info btn-sm" to="/inc-orders">Заказы для обработки</RouterLink>&nbsp;
        </span>
        <span v-if="authStore.userInfo.admin">
          <RouterLink class="btn btn-info btn-sm" to="/set-roles">admin</RouterLink>&nbsp;
        </span>
        <ul class="navbar-nav flex-grow-1">&nbsp;</ul>
        <span v-if="authStore.userInfo.id">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/myprofile">{{ authStore.userInfo.userName }}</RouterLink>&nbsp;
        </span>
      </div>
    </nav>
  </header>

  <br />
  <div class="container-fluid">
    <RouterView :key="route.fullPath" />
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
