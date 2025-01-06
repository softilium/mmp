<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';

  const route = useRoute();
  const router = useRouter();

  const order = ref(null);

  onMounted(async () => {

    let res = await fetch(authStore.rbUrl() + "/api/orders/outbox/" + route.params.id, {
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      }
    });

    if (res.ok) {
      order.value = await res.json();
    }
  });

</script>

<template>

  <div v-if="order">
    <h1>Заказ {{ order.id }}</h1>
    <h6>{{ glob.fmtDate(order.createdOn) }}</h6>
    <h6>Витрина {{ order.shop.caption }}</h6>
    <h6>Статус {{ order.status }}</h6>

    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Товар или услуга</th>
          <th>Количество</th>
          <th>Цена</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="line in order.lines">
          <td>{{line.good.caption}}</td>
          <td>{{line.qty}}</td>
          <td>{{line.sum}}</td>
        </tr>
      </tbody>
      <tfoot class="table-primary">
        <tr>
          <td>ИТОГО</td>
          <td>{{ order.qty}}</td>
          <td>{{ order.sum}}</td>
        </tr>
      </tfoot>
    </table>
  </div>

</template>
