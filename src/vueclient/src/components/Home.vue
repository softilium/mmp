<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { authStore } from './authStore.js';
  import ProfileLink from './ProfileLink.vue';

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
  <nav>
    <button class="btn btn-info btn-sm" v-if="authStore.userInfo.id && authStore.userInfo.shopManage" @click="$router.push('/edit-shop');">Добавить витрину</button>
  </nav>

  <h1>Витрины</h1>

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Заголовок</th>
          <th>Владелец</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in shops" v-bind:key="item.id">
          <td class="col-9">
            <RouterLink v-bind:to="`/shop/${item.id}`">{{ item.caption }}</RouterLink>
          </td>
          <td class="col-2"><ProfileLink :userInfo="item.createdByInfo"></ProfileLink></td>
        </tr>
      </tbody>
    </table>
  </div>
  <br />

</template>
