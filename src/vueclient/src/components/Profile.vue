<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();

  const user = ref({userName: "", email: "", phoneNumber: ""});

  const result = ref("");

  onMounted(async () => {
    if (route.params.id) {
      let res = await fetch(authStore.rbUrl() + "/api/profiles/" + route.params.id,
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (res.ok) { user.value = await res.json(); }
    }
    else {
      let res = await fetch(authStore.rbUrl() + "/api/profiles/my",
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (res.ok) { user.value = await res.json(); }
    }
  });

  const Save = async () => {
    let res = await fetch(authStore.rbUrl() + "/api/profiles",
      {
        method: "PUT",
        body: JSON.stringify(user.value),
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });
    if (res.ok)
      result.value = "Изменения записаны";
    else
      result.value = await res.text();
  }

</script>

<template>

  <h1>Профиль пользователя {{ user.userName }}&nbsp;<button class="btn btn-outline-secondary" @click="authStore.Logout(); $router.push('/');">Выйти</button></h1>
  <br /><br />

  <div class="row mb-3">
    <label class="col-3 form-label">Имя пользователя</label>
    <div class="col-7">
      <input class="form-control" v-model="user.userName" required />
    </div>
  </div>
  <div class="row mb-3">
    <label class="col-3 form-label">Email</label>
    <div class="col-7">
      <input class="form-control" v-model="user.email" required />
    </div>
  </div>
  <div class="row mb-3">
    <label class="col-3 form-label">Номер телефона</label>
    <div class="col-7">
      <input class="form-control" v-model="user.phoneNumber" />
    </div>
  </div>

  <div v-if="!route.params.id">
    <button class="btn btn-primary" @click="Save">Сохранить</button>
  </div>
  <div v-if="result" class="alert alert-primary">
    {{ result }}
  </div>

</template>
