<script setup lang="ts">

  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  let title = "New shop";
  if (route.params.id) title = "Edit shop";

  let nameField = "";
  let captionField = "";

  const Save = async () => {

    if (!authStore.loggedEmail) return;

    await fetch(authStore.rbUrl() + "/api/shops", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Autorization": "Bearer " + authStore.refreshToken
      },
      body: JSON.stringify({ name: nameField, caption: captionField })
    }).then(async res => {
      if (res.ok) {
        router.push("/");
      }
    }).catch(err => {
      console.log(err);
    });
  }

</script>

<template>

  <h1>{{ title }}</h1>

  Name:
  <input v-model="nameField" />

  Caption:
  <input v-model="captionField" />

  <button @click="Save">Save</button>

</template>
