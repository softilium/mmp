<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
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
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in shops" v-bind:key="item.id">
          <td class="col-9">
            <div class="row">
              <div class="col-11">
                <RouterLink v-bind:to="`/shop/${item.id}`">{{ item.caption }}</RouterLink>
              </div>
              <div class="col-1 text-end">
                <RouterLink v-if="authStore.userInfo.shopManage && item.createdByInfo.userName==authStore.userInfo.userName" v-bind:to="`/edit-shop/${item.id}`">
                  <span class="text-info"><i class="bi bi-pencil-square" /></span>
                </RouterLink>
              </div>
            </div>
          </td>
          <td class="col-2">{{ item.createdByInfo.userName }}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <br />
  <button class="btn btn-info" v-if="authStore.userInfo.userName && authStore.userInfo.shopManage" @click="$router.push('/edit-shop');">Добавить витрину</button>

</template>
