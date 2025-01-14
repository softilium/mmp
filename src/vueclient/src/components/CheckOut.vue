<script setup>

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();
  const router = useRouter();

  const shop = ref({});
  const isOwner = ref(false);

  const basket = ref([]);
  const qtyTotal = ref(0);
  const sumTotal = ref(0);
  const customerComment = ref("");

  onMounted(async () => {

    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops/" + route.params.shopid);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.createdByInfo.id == authStore.userInfo.id;
      } else router.push("/");
    } catch (err) {
      console.log(err);
      router.push("/");
    }

    let res = await fetch(authStore.rbUrl() + "/api/baskets/" + shop.value.id, {
      mathod: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      }
    });

    try {
      if (res.ok) {
        basket.value = await res.json();
        qtyTotal.value = 0;
        sumTotal.value = 0;
        basket.value.forEach(_ => {
          qtyTotal.value += _.qty;
          sumTotal.value += _.sum;
        });
      }
    } catch (err) { console.log(err); }

  });

  const Checkout = async () => {
    let res = await fetch(authStore.rbUrl() + "/api/orders/outbox/" + route.params.shopid, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      },
      body: customerComment.value
    });
    if (res.ok) {
      router.push("/orders");
    }
  };

</script>

<template>

  <h1>
    Сформировать заказ для витрины "{{ shop.caption }}"
  </h1>

  <table class="table">
    <thead class="table-primary">
      <tr>
        <th>Товар или услуга</th>
        <th class="text-end">Количество</th>
        <th class="text-end">Стоимость</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="line in basket" v-bind:key="line.id">
        <td>{{ line.good.caption }}</td>
        <td class="text-end">{{ line.qty }}</td>
        <td class="text-end">{{ line.sum }}</td>
      </tr>
    </tbody>
    <tfoot class="table-success">
      <tr>
        <td>ИТОГО</td>
        <td class="text-end">{{ qtyTotal }}</td>
        <td class="text-end">{{ sumTotal }}</td>
      </tr>
    </tfoot>
  </table>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Примечание к заказу</label>
      <div class="col-9">
        <textarea class="form-control" v-model="customerComment" rows="5" />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <button class="btn btn-primary" @click="Checkout()">Оформить заказ</button>
  </div>
  <RouterLink :to="`/shop/${shop.id}`">Витрина "{{shop.caption}}"</RouterLink>

</template>
