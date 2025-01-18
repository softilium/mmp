<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();
  const router = useRouter();

  const shop = ref({ id: 0, caption: "", createdByInfo: { id: 0 } });

  const goods = ref([]);

  const isOwner = ref(false);

  onMounted(async () => {

    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops/" + route.params.id);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.createdByInfo.id == authStore.userInfo.id && authStore.userInfo.shopManage;
      } else router.push("/");
    } catch (err) {
      console.log(err);
      router.push("/");
    }

    try {
      let res = await fetch(authStore.rbUrl() + "/api/goods?shopid=" + route.params.id, { signal: AbortSignal.timeout(5000), });
      if (res.ok) {
        goods.value = await res.json();
      }
    } catch (err) { console.log(err); }

    await nextTick();

    LoadBasket();

  });

  const basketSum = ref(0);

  const LoadBasket = async () => {
    if (authStore.userInfo.id) {

      try {

        basketSum.value = 0;

        let res = await fetch(authStore.rbUrl() + "/api/baskets/" + shop.value.id, {
          method: "GET",
          headers: authStore.authHeadersAppJson()
        });
        if (res.ok) {

          res = await res.json();

          if (res.length > 0) {

            let goodmap = new Map();
            goods.value.forEach(_ => { goodmap.set(_.id, _); });

            res.forEach(_ => {
              let goodId = _.good.id;
              let gObj = goodmap.get(goodId);
              if (gObj != null) {
                gObj.basked = _.qty;
                basketSum.value += gObj.price * _.qty;
              }
            });
          }
        }
      } catch (err) { console.log(err); }
    }
  }

  const Inc = async (good) => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/increase/" + good.id, {
      method: "POST",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) {
      if (!good.basked) good.basked = 0;
      good.basked++;
      LoadBasket();
    }
  }

  const Dec = async (good) => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/decrease/" + good.id, {
      method: "POST",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) {
      if (good.basked && good.basked > 0) good.basked--;
      if (good.basked == 0) good.basked = null;
      LoadBasket();
    }
  }

</script>

<template>

  <h1>{{ shop.caption }}</h1>

  <div class="row">
    <table class="table">
      <thead class="table-primary">
        <tr>
          <th>Товар или Услуга</th>
          <th class="text-end">Цена</th>
          <th>
            <button class="btn btn-outline-success btn-sm" v-if="basketSum" @click="$router.push(`/checkout/${shop.id}`);">
              <i class="bi bi-basket2-fill"></i> {{ basketSum }}
            </button>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="good in goods" v-bind:key="good.id">
          <td class="col-7">
            <RouterLink v-bind:to="`/good/${good.id}`"> {{ good.caption }} </RouterLink>
          </td>
          <td class="col-2 text-end">{{ good.price }}</td>
          <td v-if="authStore.userInfo.id" class="col-3">
            <button class="btn btn-primary btn-sm" @click="Inc(good)">+</button>&nbsp;
            <span v-if="good.basked"><button class="btn btn-primary btn-sm" @click="Dec(good)">-</button>&nbsp;</span>
            {{ good.basked }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <button class="btn btn-info btn-sm" v-if="isOwner" @click="router.push(`/edit-good/${route.params.id}`);">Добавить товар / услугу</button>
  &nbsp;
  <RouterLink class="btn btn-info btn-sm" v-if="authStore.userInfo.shopManage && shop.createdByInfo.id==authStore.userInfo.id" v-bind:to="`/edit-shop/${shop.id}`">
    Редактировать витрину
  </RouterLink>


</template>
