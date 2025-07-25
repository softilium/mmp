<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { ctx } from './ctx.js';
  import { localBasket } from '../services/localBasket';

  const route = useRoute()
  const router = useRouter()

  const good = ref({
    OwnerShop: { Ref: "", Caption: "" },
    Caption: "",
    Description: "",
    Basked: null,
    Ref: "",
    Price: 0,
    Article: "",
    Url: "",
    CreatedBy: { Ref:"", Username: "" }
  });
  const basketQty = ref(0);
  const isOwner = ref(false);

  onMounted(async () => {
    if (route.params.id) {
      try {
        let res = await fetch(ctx.rbUrl() + "/api/goods?ref=" + route.params.id);
        if (res.ok) {
          good.value = await res.json();
          isOwner.value = (good.value.CreatedBy.Ref == ctx.userInfo.id);
          good.value.Basked = null;
          basketQty.value = null;
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
       let res = await fetch(`${ctx.rbUrl()}/api/goods/images?ref=${route.params.id}&n=${i}`, { method: "GET" });
       if (res.status == 200) { // status 204 means no image
         let b = await res.blob();
         const src = URL.createObjectURL(b);
         imageSrc.value.push(src);
       }
     }
    isImageLoading.value = false;
  }

  const LoadBasket = async () => {
    basketQty.value = null;
    good.value.Basked = null;
    if (ctx.userInfo.id) {
      let res = await fetch(ctx.rbUrl() + "/api/basket?shop=" + good.value.OwnerShop.Ref, {
        method: "GET",
        headers: await ctx.authHeadersAppJson()
      });
      if (res.ok) {
        res = await res.json();
        res.Data.forEach(l => {
          if (l.Good.Ref == good.value.Ref) {
            good.value.Basked = l.Qty;
            basketQty.value = l.Qty;
          }
        });
      }
    } else {
      // Check localBasket for anonymous users
      const items = localBasket.getItems();
      const found = items.find(i => i.goodId === good.value.Ref);
      if (found) {
        good.value.Basked = found.quantity;
        basketQty.value = found.quantity;
      }
    }
  }

  const Inc = async () => {
    if (!ctx.userInfo.id) {
      localBasket.addItem({
        goodId: good.value.Ref,
        quantity: 1,
        price: good.value.Price,
        title: good.value.Caption,
        shopTitle: good.value.OwnerShop ? good.value.OwnerShop.Caption : '',
        senderId: good.value.CreatedBy.Ref,
        shopId: good.value.OwnerShop ? good.value.OwnerShop.Ref : null
      });
      await LoadBasket();
      await ctx.loadBasket();
    }
    let res = await fetch(ctx.rbUrl() + "/api/basket/increase?goodref=" + good.value.Ref, {
      method: "POST",
      headers: await ctx.authHeadersAppJson()
    });
    if (res.ok) {
      LoadBasket();
      await ctx.loadBasket();
    }
  }

  const Dec = async () => {
    if (!ctx.userInfo.id) {
      localBasket.decItem(good.value.Ref);
      await LoadBasket();
      await ctx.loadBasket();
      return;
    }
    let res = await fetch(ctx.rbUrl() + "/api/basket/decrease?goodref=" + good.value.Ref, {
      method: "POST",
      headers: await ctx.authHeadersAppJson()
    });
    if (res.ok) {
      LoadBasket();
      await ctx.loadBasket();
    }
  }

  const DeleteGood = async () => {

    if (!confirm('Удалить товар, вы уверены?')) return;

    let res = await fetch(`${ctx.rbUrl()}/api/goods?ref=${good.value.Ref}`, {
      method: "DELETE",
      headers: await ctx.authHeadersAppJson()
    });
    if (res.ok) {
      router.push("/shop/" + route.params.shopid);
    }
  }

</script>

<template>
  <nav>
    <RouterLink v-if="isOwner" v-bind:to="`/edit-good/${good.OwnerShop.Ref}/${good.Ref}`">
      <span class="btn btn-info btn-sm">Редактировать товар</span>
    </RouterLink>
    &nbsp;
    <button class="btn btn-info btn-sm" v-if="isOwner" @click="DeleteGood();">Удалить товар</button>
  </nav>
  <RouterLink :to="`/shop/${good.OwnerShop.Ref}`">Витрина "{{good.OwnerShop.Caption}}"</RouterLink>
  <div>&nbsp;</div>
  <h1>{{good.Caption}}</h1>

  <div class="row">
    <div class="col-6 col-md-3">
      Положить в корзину
    </div>
    <div class="col-6 col-md-9">
      <button class="btn btn-primary btn-sm" @click="Inc()">+</button>&nbsp;
      <span v-if="good.Basked">
        <button class="btn btn-primary btn-sm" @click="Dec()">-</button>&nbsp; {{ basketQty }}
      </span>
    </div>
  </div>
  <div class="row">
    <div class="col-6 col-md-3">Цена</div>
    <div class="col-6 col-md-9">{{good.Price}}</div>
  </div>
  <div class="row mb-3">
    <div class="col-6 col-md-3">Артикул</div>
    <div class="col-6 col-md-9">{{good.Article}}</div>
  </div>
  <div v-if="good.Url" class="row mb-3">
    <div class="col"><a :href="good.Url" target="_blank">{{good.Url}}</a></div>
  </div>

  <div class="row-mb3">
    <div class="col">{{good.Description}}</div>
  </div>

  <!-- <div class="row mb-3" v-if="imageSrc.length>1">
    <div class="col">
      <span v-for="(src, index) in imageSrc" :key="index">
        <button :class="`${index==curImgIndex ? 'btn btn-secondary btn-sm' : 'btn btn-outline-secondary btn-sm'}`" @click="curImgIndex=index">{{index+1}}</button>&nbsp;
      </span>
    </div>
  </div> -->

  <div v-if="!isImageLoading && imageSrc.length>0" class="row mb-3">
    <div class="col">
      <img :src="imageSrc[curImgIndex]" class="d-block w-100">
    </div>
  </div>

</template>
