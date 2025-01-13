<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();
  const router = useRouter();

  const shop = ref({ id: 0, caption: "", createdByInfo: { userName: "" } });

  const goods = ref([]);

  const isOwner = ref(false);

  onMounted(async () => {

    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops/" + route.params.id);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.createdByInfo.userName == authStore.userInfo.userName && authStore.userInfo.shopManage;
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
    if (authStore.userInfo.userName) {

      try {

        basketSum.value = 0;

        let res = await fetch(authStore.rbUrl() + "/api/baskets/" + shop.value.id, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
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
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      },
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
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      },
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
          <td class="col-9">
            <div class="row">
              <div class="col-11">
                <RouterLink v-bind:to="`/good/${good.id}`"> {{ good.caption }} </RouterLink>
              </div>
              <div class="col-1 text-end">
                <RouterLink v-if="isOwner" v-bind:to="`/edit-good/${shop.id}/${good.id}`">
                  <span class="text-info"><i class="bi bi-pencil-square" /></span>
                </RouterLink>
              </div>
            </div>
          </td>
          <td class="col-1 text-end">{{ good.price }}</td>
          <td v-if="authStore.userInfo.userName" class="col-2">
            <button class="btn btn-primary btn-sm" @click="Inc(good)">+</button>&nbsp;
            <span v-if="good.basked"><button class="btn btn-primary btn-sm" @click="Dec(good)">-</button>&nbsp;</span>
            {{ good.basked }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <button class="btn btn-primary" v-if="isOwner" @click="router.push(`/edit-good/${route.params.id}`);">Добавить товар / услугу</button>

</template>
