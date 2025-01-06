<script setup lang="ts">

  import { onMounted, ref, computed } from 'vue';
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
      <RouterLink class="navbar-brand" to="/">Витрины</RouterLink>
      <span v-if="authStore.userName">
        <RouterLink class="navbar-text" to="/orders">Заказы</RouterLink>&nbsp;&nbsp;
      </span>
      <span v-if="!authStore.userName">
        <RouterLink class="navbar-text" to="/login">Войти</RouterLink>&nbsp;&nbsp;
      </span>
      <span v-if="authStore.userName">
        <RouterLink class="navbar-text" to="/logout">Выйти</RouterLink>&nbsp;&nbsp;
      </span>
      <span v-if="!authStore.userName">
        <RouterLink class="navbar-text" to="/register">Зарегистрироваться</RouterLink>&nbsp;&nbsp;
      </span>
      <span v-if="authStore.userName">
        <RouterLink class="navbar-text" to="/profile">{{ authStore.userName }}</RouterLink>
      </span>
    </nav>
  </header>
  <br />

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
