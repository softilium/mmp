<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';

  const route = useRoute()
  const router = useRouter()

  const order = ref({
    customerComment: "",
    senderComment: "",
    expectedDeliveryDate: "",
    createdOn: "",
    shop: { caption: "", createdByID: 0 },
    status: 0,
    id: 0,
    lines: [],
    qty: 0,
    sum: 0,
    createdByID: 0
  });
  const statuses = ref([]);

  const isSender = ref(false);
  const isCustomer = ref(false);

  onMounted(async () => {

    let res = await fetch(authStore.rbUrl() + "/api/orders/statuses");
    if (res.ok) {
      statuses.value = await res.json();
    }

    res = await fetch(authStore.rbUrl() + "/api/orders/" + route.params.id, {
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      }
    });

    if (res.ok) {
      order.value = await res.json();
      isCustomer.value = order.value.createdByID == authStore.userInfo.id;
      isSender.value = order.value.shop.createdByID == authStore.userInfo.id;
    }
  });

  const Save = async () => {

    if (isCustomer.value) {
      let req = await fetch(`${authStore.rbUrl()}/api/orders/outbox/${route.params.id}`,
        {
          method: "PUT",
          body: JSON.stringify(order.value),
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (req.ok) {
        router.push("/orders");
      }
    }
    if (isSender.value) {
      let req = await fetch(`${authStore.rbUrl()}/api/orders/inbox/${route.params.id}`,
        {
          method: "PUT",
          body: JSON.stringify(order.value),
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (req.ok) {
        router.push("/inc-orders");
      }
    }
  }

</script>

<template>

  <div v-if="order">
    <h1>Заказ {{ order.id }}</h1>
    <h6>{{ glob.fmtDate(order.createdOn) }}</h6>
    <h6>Витрина {{ order.shop.caption }}</h6>

    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Товар или услуга</th>
          <th>Количество</th>
          <th>Цена</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="line in order.lines" v-bind:key="line.id">
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

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Статус</label>
      <div class="col-9">
        <select :disabled="!isSender" class="form-select form-control-sm" v-model="order.status">
          <template v-for="(value,key) in statuses">
            <option v-bind:value="Number(key)">{{ value }}</option>
          </template>
        </select>
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Примечание покупателя</label>
      <div class="col-9">
        <textarea :disabled="!isCustomer" class="form-control" v-model="order.customerComment" rows="5" />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Примечание продавца</label>
      <div class="col-9">
        <textarea :disabled="!isSender" class="form-control" v-model="order.senderComment" rows="5" />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Ожидаемая дата доставки</label>
      <div class="col-9">
        <input :disabled="!isSender" type="datetime-local" class="form-control" v-model="order.expectedDeliveryDate" />
      </div>
    </div>
  </div>
  <button :disabled="!isCustomer && !isSender" class="btn btn-primary" @click="Save()">Сохранить</button>

</template>
