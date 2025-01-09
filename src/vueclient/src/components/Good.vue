<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const good = ref({ ownerShop: { id: 0, caption: "" }, caption: "", description: "", basked: null, id:0, price:0, article:"", url:"" });
  const basketSum = ref(0);

  onMounted(async () => {
    if (route.params.id) {
      try {
        let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id);
        if (res.ok) {
          good.value = await res.json();
          good.value.basked = null;
          basketSum.value = null;
          LoadBasket();
        } else router.push("/shop/" + route.params.shopid);
      } catch (err) {
        console.log(err);
        router.push("/shop/" + route.params.shopid);
      }
    }
  });

  const LoadBasket = async () => {
    basketSum.value = null;
    good.value.basked = null;
    if (authStore.userInfo.userName) {
      let res = await fetch(authStore.rbUrl() + "/api/baskets/" + good.value.ownerShop.id, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });
      if (res.ok) {
        res = await res.json();
        res.forEach(l => {
          if (l.good.id == good.value.id) {
            good.value.basked = l.qty;
            basketSum.value = good.value.price * l.qty;
          }
        });
      }
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
    if (res.ok) LoadBasket();
  }

  const Dec = async (good) => {
    let res = await fetch(authStore.rbUrl() + "/api/baskets/decrease/" + good.id, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      },
    });
    if (res.ok) LoadBasket();
  }


</script>

<template>
  <h1>{{good.caption}}</h1>
  <div>&nbsp;</div>

  <table class="table">
    <tbody>
      <tr v-if="authStore.userInfo.userName">
        <td>В корзине</td>
        <td>
          <button class="btn btn-primary btn-sm" @click="Inc(good)">+</button>&nbsp;
          <span v-if="good.basked">
            <button class="btn btn-primary btn-sm" @click="Dec(good)">-</button>&nbsp; {{ basketSum }}
            <button class="btn btn-outline-success btn-sm" @click="$router.push(`/checkout/${good.ownerShop.id}`);">
              <i class="bi bi-basket2-fill"></i>
            </button>
          </span>
        </td>
      </tr>
      <tr>
        <td class="col-2">Артикул</td>
        <td>{{good.article}}</td>
      </tr>
      <tr>
        <td class="col-2">Ссылка</td>
        <td>{{good.url}}</td>
      </tr>
      <tr>
        <td class="col-2">Описание</td>
        <td>{{good.description}}</td>
      </tr>
      <tr>
        <td class="col-2">Цена</td>
        <td>{{good.price}}</td>
      </tr>
    </tbody>
  </table>

  <RouterLink :to="`/shop/${good.ownerShop.id}`">Витрина "{{good.ownerShop.caption}}"</RouterLink>

</template>
