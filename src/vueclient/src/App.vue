<script setup lang="ts">

  import { onMounted, ref, watch } from 'vue';
  import { authStore } from './components/authStore.js';
  import { useRoute } from 'vue-router';

  const route = useRoute();

  const myurl = ref("");

  const updateMyUrl = () => {
    // telegram hash param delete from URL
    const url = new URL(window.location.href);
    url.hash = "";
    myurl.value = url.href;
  }

  onMounted(async () => {
    authStore.basket.sum = 0;
    authStore.SetAccessToken(localStorage.getItem("accessToken"));
    authStore.SetRefreshToken(localStorage.getItem("refreshToken"));
    authStore.CheckLogged();
    updateMyUrl();
    await authStore.loadBasket();
  });

  // update myurl when user navigate within SPA
  watch(
    () => route.fullPath,
    async newId => {
      updateMyUrl();
    }
  )

</script>

<template>

  <header>
    <nav class="navbar navbar-expand-sm navbar-toggleable-sm navbar-light bg-white border-bottom box-shadow.mb-3">
      <div class="container-fluid">
        <RouterLink class="btn btn-outline-secondary btn-sm" to="/">Витрины</RouterLink>&nbsp;&nbsp;
        <span v-if="!authStore.userInfo.id">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/login">Войти</RouterLink>&nbsp;
        </span>
        <ul class="navbar-nav flex-grow-1">&nbsp;</ul>
        <span v-if="authStore.userInfo.id && authStore.basket.sum">
          <RouterLink class="btn btn-outline-success btn-sm" to="/checkout"><i class="bi bi-basket2-fill"></i>&nbsp;{{authStore.basket.sum}}</RouterLink>&nbsp;&nbsp;
        </span>
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
  <footer class="border-top footer text-muted">
    <nav class="navbar navbar-expand-sm navbar-toggleable-sm navbar-light bg-white border-bottom box-shadow.mb-3">
      <div class="container-fluid">
        <span v-if="authStore.userInfo.id">
          <RouterLink class="btn btn-info btn-sm" to="/inc-orders">Заказы для обработки</RouterLink>&nbsp;
        </span>
        <span v-if="authStore.userInfo.admin">
          <RouterLink class="btn btn-info btn-sm" to="/set-roles">admin</RouterLink>&nbsp;
        </span>
      </div>
    </nav>
    <div class="container-fluid">
      <div class="row">
        <div class="col">
          <span>&copy;2024-2025, </span>
          <div v-if="authStore.isTg() && authStore.userInfo.shopManage">Адрес страницы для браузера: <input class="form-control sm" v-model="myurl" /></div>
        </div>
      </div>
    </div>
  </footer>

</template>
