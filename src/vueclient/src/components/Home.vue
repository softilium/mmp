<script setup lang="ts">

  import { onMounted, ref, computed } from 'vue';
  import { authStore } from './authStore.js';

  let shops = ref([]);

  onMounted(async () => {
    await fetch(authStore.rbUrl() + "/api/shops")
      .then(async res => { shops.value = await res.json(); })
      .catch(err => { console.log(err); });
  });
      
</script>

<template>

  <h1>Home</h1>

  <h2>Shops</h2>
  <table>
    <tr v-for="item in shops">
      <td>{{ item.name }}</td><td>{{ item.caption }}</td>
    </tr>
  </table>

  <button v-if="authStore.loggedEmail" @click="$router.push('/edit-shop');">Add a shop</button>

</template>
