<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';

  const route = useRoute();
  const router = useRouter();

  const orders = ref([]);
  const statuses = ref([]);
  const showAll = ref(false);

  const Load = async () => {

    let url = `${authStore.rbUrl()}/api/orders/outbox`;
    if (showAll.value) url += "?showAll=1";
    let res = await fetch(url,
      {
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });

    if (res.ok) {
      orders.value = await res.json();
    }
  }

  onMounted(async () => {

    let res = await fetch(`${authStore.rbUrl()}/api/orders/statuses`);
    if (res.ok) {
      statuses.value = await res.json();
    }
    Load();

  });

</script>

<template>

  <h1>Ваши заказы</h1>
  Показывать закрытые
  <input v-model="showAll" type="checkbox" @change="Load();" />

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th class="col-3">Витрина</th>
          <th class="col-3">Отправитель</th>
          <th class="col-2">Статус</th>
          <th class="col-2">Создано</th>
          <th class="col-2 text-end">Сумма</th>
        </tr>
      </thead>
      <tr v-for="order in orders">
        <td><RouterLink :to="`/order/${order.id}`">{{ order.shop.caption }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ order.shop.createdByInfo.userName }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ statuses[order.status] }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ glob.fmtDate(order.createdOn) }}</RouterLink></td>
        <td class="text-end"><RouterLink :to="`/order/${order.id}`">{{ order.sum }}</RouterLink></td>
      </tr>
    </table>
  </div>

</template>
