<script setup>

  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const router = useRouter();

  const basket = ref([]);

  onMounted(async () => {
    await Load();
  });

  const Load = async () => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/", {
      mathod: "GET",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) {
      res = await res.json();
      const senderMap = new Map();
      res.forEach(_ => {
        senderMap.set(_.good.createdByID, 1);
      });
      basket.value = [];
      for (let [k] of senderMap.entries()) {
        let senderRec = { senderID: k, lines: [], qtyTotal: 0, sumTotal: 0, customerComment: "" };
        basket.value.push(senderRec);
        res.forEach(r => {
          if (r.good.createdByID == k) {
            senderRec.lines.push(r);
            senderRec.qtyTotal += r.qty;
            senderRec.sumTotal += r.sum;
            senderRec.senderInfo = r.good.createdByInfo;
          }
        });
      }
    }
    if (basket.value.length == 0)
      router.push("/");
  }

  const Checkout = async (senderRec) => {
    let res = await fetch(authStore.rbUrl() + "/api/orders/outbox/" + senderRec.senderID, {
      method: "POST",
      headers: authStore.authHeadersAppJson(),
      body: senderRec.customerComment
    });
    if (res.ok) {
      await authStore.loadBasket();
      res = await res.json();
      router.push(`/order/${res.id}`);      
    }
  };

  const Inc = async (good) => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/increase/" + good.id, {
      method: "POST",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) {
      await Load();
      authStore.loadBasket();
    }
  }

  const Dec = async (good) => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/decrease/" + good.id, {
      method: "POST",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) {
      await Load();
      authStore.loadBasket();
    }
  }

</script>

<template>

  <!--<h1>Сформировать заказ для витрины "{{ shop.caption }}"</h1>
  <RouterLink :to="`/shop/${shop.id}`">Витрина "{{shop.caption}}"</RouterLink>
  <div>&nbsp;</div>-->

  <div v-for="sender in basket">
    <h1>
      Отправитель <RouterLink :to="`/profile/${sender.senderInfo.id}`">{{sender.senderInfo.userName}}</RouterLink>
    </h1>

    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Товар или услуга</th>
          <th class="text-end">Количество</th>
          <th class="text-end">Стоимость</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="line in sender.lines" v-bind:key="line.id">
          <td><RouterLink :to="`/good/${line.good.id}`">{{ line.good.caption }}</RouterLink></td>
          <td class="text-end">
            {{ line.qty }}
            &nbsp;<button class="btn btn-primary btn-sm" @click="Inc(line.good)">+</button>
            &nbsp;<button class="btn btn-primary btn-sm" @click="Dec(line.good)">-</button>
          </td>
          <td class="text-end">{{ line.sum }}</td>
        </tr>
      </tbody>
      <tfoot class="table-success">
        <tr>
          <td></td>
          <td class="text-end">{{ sender.qtyTotal }}</td>
          <td class="text-end">{{ sender.sumTotal }}</td>
        </tr>
      </tfoot>
    </table>

    <div class="row mb-3">
      <div class="form-group form-group-sm row">
        <label class="col-4 form-label">Примечание к заказу</label>
        <div class="col-8">
          <textarea class="form-control" v-model="sender.customerComment" rows="5" />
        </div>
      </div>
    </div>

    <div class="row mb-3">
      <button class="btn btn-primary btn-sm" @click="Checkout(sender)">Оформить заказ для {{sender.senderInfo.userName}}</button>
    </div>

  </div>

  <!--<div v-if="shop.deliveryConditions" class="text-center">
    <div class="row mb-3">
      <div class="col">
        Условия доставки, важно прочитать перед заказом:<br />
        <strong>
          <span v-html="shop.deliveryConditions"></span>
        </strong>
        <br />
        <br />
      </div>
      <div class="row mb-3">
        <div class="col">
          Условия понятны <input v-model="okConditions" type="checkbox" />
        </div>
      </div>
    </div>
  </div>-->

</template>
