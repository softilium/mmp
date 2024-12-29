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

  <div>
    <RouterLink to="/">Home</RouterLink> |
    <span v-if="!authStore.loggedEmail">
      <RouterLink to="/login">Login</RouterLink> |
    </span>
    <span v-if="authStore.loggedEmail">
      <RouterLink to="/logout">Logout</RouterLink> |
    </span>
    <span v-if="!authStore.loggedEmail">
      <RouterLink to="/register">Register</RouterLink> |
    </span>
    <span v-if="authStore.loggedEmail">Logged as {{ authStore.loggedEmail }}</span>
  </div>
  <br />
  <RouterView />

</template>
