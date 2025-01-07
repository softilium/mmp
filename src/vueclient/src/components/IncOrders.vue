<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';

  const route = useRoute();
  const router = useRouter();

  const orders = ref([]);
  const statuses = ref([]);

  onMounted(async () => {

    let res = await fetch(authStore.rbUrl() + "/api/orders/statuses");
    if (res.ok) {
      statuses.value = await res.json();
    }

    res = await fetch(authStore.rbUrl() + "/api/orders/inbox",
      {
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });

    if (res.ok) {
      orders.value = await res.json();
    }
  });

</script>

<template>

  <h1>Ваши заказы</h1>

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Витрина</th>
          <th>Заказчик</th>
          <th>Статус</th>
          <th>Создано</th>
          <th class="text-end">Сумма</th>
        </tr>
      </thead>
      <tr v-for="order in orders">
        <td><RouterLink :to="`/order/${order.id}`">{{ order.shop.caption }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ order.createdByInfo.userName }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ statuses[order.status] }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ glob.fmtDate(order.createdOn) }}</RouterLink></td>
        <td class="text-end"><RouterLink :to="`/order/${order.id}`">{{ order.sum }}</RouterLink></td>
      </tr>
    </table>
  </div>

</template>
