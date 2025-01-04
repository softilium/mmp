<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const good = ref({ ownerShop: { id: 0 }, caption: "", description: "" });

  onMounted(async () => {
    if (route.params.id) {
      try {
        let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id);
        if (res.ok) {
          good.value = await res.json();
        } else router.push("/shop/" + route.params.shopid);
      } catch (err) {
        console.log(err);
        router.push("/shop/" + route.params.shopid);
      }
    }
  });

  const Save = async () => {

    if (route.params.id) {
      let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        },
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        router.push("/shop/" + route.params.shopid);
      }

    } else {
      good.value.ownerShop.id = route.params.shopid;

      let res = await fetch(authStore.rbUrl() + "/api/goods", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        },
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        router.push("/shop/" + route.params.shopid);
      }
    }
  };
    
</script>

<template>

  <div>
    Caption: <input v-model="good.caption" />
    Description: <input v-model="good.description" />
  </div>

  <button @click="Save">Save and return to shop</button>

</template>
