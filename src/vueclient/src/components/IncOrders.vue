<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';


  const orders = ref([]);
  const statuses = ref([]);

  const showAll = ref(false);

  const Load = async () => {

    let url = `${authStore.rbUrl()}/api/orders/inbox`;
    if (showAll.value) url += "?showAll=1"
    let res = await fetch(url,
      {
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });

    if (res.ok) {
      res = await res.json();
      orders.value = res;
    } else console.log(res);
  }

  onMounted(async () => {

    let res = await fetch(authStore.rbUrl() + "/api/orders/statuses");
    if (res.ok) {
      statuses.value = await res.json();
    }

    Load();

  });

  const Save = async () => {
    orders.value.forEach((order) => {
      if (order.newStatus) {
        try {
          fetch(`${authStore.rbUrl()}/api/orders/inbox?orderid=${order.id}&newstatus=${order.newStatus}`,
            {
              method: "PUT",
              headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + authStore.accessToken
              }
            });
        } catch (err) { console.log(err); }
      }
    });
    window.location.reload();
    Load();
  }

</script>

<template>

  <h1>Ваши заказы</h1>

  Показывать закрытые
  <input v-model="showAll" type="checkbox" @change="Load();" />

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th class="col-2">Витрина</th>
          <th class="col-2">Заказчик</th>
          <th class="col-2">Статус</th>
          <th class="col-2">Создано</th>
          <th class="col-2 text-end">Сумма</th>
          <th class="col-2">Новый статус</th>
        </tr>
      </thead>
      <tr v-for="order in orders" :key="order.id">
        <td><RouterLink :to="`/order/${order.id}`">{{ order.shop.caption }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ order.createdByInfo.userName }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ statuses[order.status] }}</RouterLink></td>
        <td><RouterLink :to="`/order/${order.id}`">{{ glob.fmtDate(order.createdOn) }}</RouterLink></td>
        <td class="text-end"><RouterLink :to="`/order/${order.id}`">{{ order.sum }}</RouterLink></td>
        <td>
          <select class="form-select form-control-sm" v-model="order.newStatus">
            <option></option>
            <template v-for="(value,key) in statuses">
              <option v-if="key != order.status" v-bind:value="key">{{ value }}</option>
            </template>
          </select>
        </td>
      </tr>
    </table>
  </div>

  <button class="btn btn-primary" @click="Save">Сохранить</button>

</template>
