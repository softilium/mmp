<script setup lang="ts">
  import Home from './components/Home.vue';
  import Login from './components/Login.vue';
  import Register from './components/Register.vue';
  import { onMounted, ref, computed } from 'vue';
  import { authStore } from './components/authStore.js';

  // ### navigation (begin)
  const routes = {
    '/': Home,
    '/login': Login,
    '/register': Register
  }
  const currentPath = ref(window.location.hash)
  window.addEventListener('hashchange', () => {
    currentPath.value = window.location.hash
  })
  const currentView = computed(() => {
    return routes[currentPath.value.slice(1) || '/'] || Home
  })

  // ### navigation (end)

  onMounted(() => {
    authStore.SetAccessToken(localStorage.getItem("accessToken"));
    authStore.SetRefreshToken(localStorage.getItem("refreshToken"));
    authStore.CheckLogged();
  });

  //let emailString = "";
  //let passwordString = "";

</script>

<template>

  <div>
    <a href="#/">Home</a> |
    <span v-if="!authStore.loggedEmail">
    <a href="#/login">Login</a> |
    </span>
    <span v-if="!authStore.loggedEmail">
      <a href="#/register">Register</a> |
    </span>
    <span v-if="authStore.loggedEmail">Logged as {{ authStore.loggedEmail }}</span>
  </div>
  <br />
  <component :is="currentView" />

</template>
