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

  <h1>Витрины</h1>

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Заголовок</th>
          <th>Владелец</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in shops">
          <td class="col-9"><RouterLink v-bind:to="`/shop/${item.id}`">{{ item.caption }}</RouterLink></td>
          <td class="col-2">{{ item.createdBy.userName }}</td>
          <td class="col-1">
            <RouterLink v-if="item.createdBy.email==authStore.loggedEmail" v-bind:to="`/edit-shop/${item.id}`">
              <i class="bi bi-pencil-square" />
            </RouterLink>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <br />
  <button class="btn btn-primary" v-if="authStore.loggedEmail" @click="$router.push('/edit-shop');">Добавить витрину</button>

</template>
