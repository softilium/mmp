<script setup lang="ts">

  import { onMounted, ref, computed } from 'vue';
  import { authStore } from './authStore.js';

  let shops = ref([]);

  onMounted(async () => {
    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops", { signal: AbortSignal.timeout(5000), });
      if (res.ok) {
        shops.value = await res.json();
      }
    } catch (err) { console.log(err); }
  });

</script>

<template>

  <h1>Shops</h1>

  <table>
    <tr v-for="item in shops">
      <td><RouterLink v-bind:to="`/shop/${item.id}`">{{ item.caption }}</RouterLink></td>
      <td>{{ item.createdBy.userName }}</td>
      <td><RouterLink v-if="item.createdBy.email==authStore.loggedEmail" v-bind:to="`/edit-shop/${item.id}`">[edit]</RouterLink></td>
    </tr>
  </table>

  <button v-if="authStore.loggedEmail" @click="$router.push('/edit-shop');">Add a shop</button>

</template>
