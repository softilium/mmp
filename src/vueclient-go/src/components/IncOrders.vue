<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { ctx } from './ctx.js';

  const orders = ref([]);
  const showAll = ref(false);
  const statuses = ref([]);

  const Load = async () => {

    let url = `${ctx.rbUrl()}/api/orders/inbox`;
    if (showAll.value) url += "?showAll=1"
    let res = await fetch(url,
      {
        headers: await ctx.authHeadersAppJson()
      });
    if (await ctx.CheckUnauth(res)) return;
    if (res.ok) {
      res = await res.json();
      orders.value = res.Data;
    } else console.log(res);
  }

  onMounted(async () => {
    let res = await fetch(ctx.rbUrl() + "/api/orders/statuses");
    if (res.ok) {
      statuses.value = await res.json();
    }
    Load();
  });

</script>

<template>

  <h1>Заказы для обработки</h1>

  Показывать закрытые
  <input v-model="showAll" type="checkbox" @change="Load();" />

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th class="col-2">Заказчик</th>
          <th class="col-2">Статус</th>
          <th class="col-2">Создано</th>
          <th class="col-2 text-end">Сумма</th>
        </tr>
      </thead>
      <tr v-for="order in orders" :key="order.Ref">
        <td>
          <RouterLink :to="`/order/${order.Ref}`">{{ order.CreatedBy.Username }}</RouterLink>
        </td>
        <td>
          <RouterLink :to="`/order/${order.Ref}`">{{ statuses[order.Status] }}</RouterLink>
        </td>
        <td>
          <RouterLink :to="`/order/${order.Ref}`">{{ ctx.fmtDate(order.CreatedAt) }}</RouterLink>
        </td>
        <td class="text-end">
          <RouterLink :to="`/order/${order.Ref}`">{{ order.Sum }}</RouterLink>
        </td>
      </tr>
    </table>
  </div>

</template>
