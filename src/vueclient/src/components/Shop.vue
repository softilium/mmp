<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const goods = ref([]);
  const shop = ref({});

  const isOwner = ref(false);

  onMounted(async () => {

    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops/" + route.params.id);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.createdBy.email == authStore.loggedEmail;
      } else router.push("/");
    } catch (err) {
      console.log(err);
      router.push("/");
    }

    try {
      let res = await fetch(authStore.rbUrl() + "/api/goods?shopid="+route.params.id, { signal: AbortSignal.timeout(5000), });
      if (res.ok) {
        goods.value = await res.json();
      }
    } catch (err) { console.log(err); }
  });
    
</script>

<template>

  <h1>{{ shop.caption }}</h1>

  <table>
    <tr v-for="good in goods">
      <td>{{ good.caption }}</td>
      <td>{{ good.description }}</td>
    </tr>
  </table>

  <button v-if="isOwner" @click="router.push(`/edit-good/${route.params.id}`);">Add good</button>

</template>
