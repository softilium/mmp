<!-- eslint-disable vue/multi-word-component-names -->
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
          LoadImages();
        } else router.push("/shop/" + route.params.shopid);
      } catch (err) {
        console.log(err);
        router.push("/shop/" + route.params.shopid);
      }
    }
  });

  const isImageLoading = ref(true);
  const maxImagesCnt = ref(3);
  const imageSrc = ref([]);
  let curImgIndex = ref(0);

  const LoadImages = async () => {
    for (let i = 0; i < maxImagesCnt.value; i++) {
      let res = await fetch(`${authStore.rbUrl()}/api/goods/images/${route.params.id}/${i}`, { method: "GET" });
      if (res.status == 200) { // status 204 means no image
        let b = await res.blob();
        const src = URL.createObjectURL(b);
        imageSrc.value.push(src);
      }
    }
    isImageLoading.value = false;
  }

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
  <RouterLink :to="`/shop/${good.ownerShop.id}`">Витрина "{{good.ownerShop.caption}}"</RouterLink>
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

  <div class="row mb-3" v-if="imageSrc.length>1">
    <div class="col">
      <span v-for="(src, index) in imageSrc" :key="index">
        <button :class="`${index==curImgIndex ? 'btn btn-secondary btn-sm' : 'btn btn-outline-secondary btn-sm'}`" @click="curImgIndex=index">{{index+1}}</button>&nbsp;
      </span>
    </div>
  </div>

  <div v-if="!isImageLoading && imageSrc.length>0" class="row mb-3">
    <div class="col">
      <img :src="imageSrc[curImgIndex]" class="d-block w-100">
    </div>
  </div>

</template>
