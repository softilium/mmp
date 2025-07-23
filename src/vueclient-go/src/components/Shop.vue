<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { ctx } from './ctx.js';
  import { localBasket } from '../services/localBasket';

  const route = useRoute();
  const router = useRouter();

  const shop = ref({ Ref: 0, Caption: "", Description: "", CreatedBy: { Ref: "", Username: "" } });
  const shopDescription = ref("");

  const goods = ref({Data: [], PagesCount: 0});

  const isOwner = ref(false);

  onMounted(async () => {

    try {
      let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + route.params.id);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.CreatedBy.Ref == ctx.userInfo.Ref && ctx.userInfo.ShopManager;
        shopDescription.value = ctx.linkify(shop.value.Description);
      } else router.push("/");
    } catch (err) {
      console.log(err);
      router.push("/");
    }

    try {
      let res = await fetch(ctx.rbUrl() + "/api/goods?shopref=" + route.params.id, { signal: AbortSignal.timeout(5000), });
      if (res.ok) {
        goods.value = await res.json();

        goods.value.Data.forEach(async g => {
          // let res = await fetch(`${ctx.rbUrl()}/api/goods/thumbs/${g.id}/0`, { method: "GET" });
          // if (res.status == 200) { // status 204 means no image
          //   let b = await res.blob();
          //   const src = URL.createObjectURL(b);
          //   g.thumb = src;
          // }
        });

      }
    } catch (err) { console.log(err); }

    await nextTick();

    LoadBasket();

  });

  const LoadBasket = async () => {
    if (ctx.userInfo.Ref) {

      try {
        let res = await fetch(ctx.rbUrl() + "/api/baskets/" + shop.value.Ref, {
          method: "GET",
          headers: ctx.authHeadersAppJson()
        });
        if (res.ok) {

          res = await res.json();

          if (res.formData.length > 0) {

            let goodmap = new Map();
            goods.value.Data.forEach(_ => { goodmap.set(_.id, _); });

            res.Data.forEach(_ => {
              let goodId = _.good.id;
              let gObj = goodmap.get(goodId);
              if (gObj != null) {
                gObj.basked = _.qty;
              }
            });
          }
        }
      } catch (err) { console.log(err); }
    } else {
      // Anonymous users: use localBasket
      const items = localBasket.getItems();
      goods.value.Data.forEach(good => {
        const found = items.find(i => i.goodId === good.id);
        if (found) {
          good.basked = found.quantity;
        } else {
          good.basked = null;
        }
      });
    }
  }

  const Inc = async (good) => {
    if (!ctx.userInfo.Ref) {
      localBasket.addItem({
        goodId: good.id,
        quantity: 1,
        price: good.price,
        title: good.caption,
        shopTitle: shop.value.Caption,
        senderId: shop.value.CreatedBy.Ref,
        shopId: shop.value.Ref
      });
      await ctx.loadBasket();
      LoadBasket();
      return;
    }
    let res = await fetch(ctx.rbUrl() + "/api/baskets/increase/" + good.id, {
      method: "POST",
      headers: ctx.authHeadersAppJson()
    });
    if (res.ok) {
      if (!good.basked) good.basked = 0;
      good.basked++;
      await ctx.loadBasket();
      LoadBasket();
    }
  }

  const Dec = async (good) => {
    if (!ctx.userInfo.Ref) {
      localBasket.decItem(good.id);
      LoadBasket();
      await ctx.loadBasket();
      return;
    }
    let res = await fetch(ctx.rbUrl() + "/api/baskets/decrease/" + good.id, {
      method: "POST",
      headers: ctx.authHeadersAppJson()
    });
    if (res.ok) {
      if (good.basked && good.basked > 0) good.basked--;
      if (good.basked == 0) good.basked = null;
      LoadBasket();
      await ctx.loadBasket();
    }
  }

  const DeleteShop = async () => {

    if (!confirm("Удалить витрину, вы уверены?")) return;

    let res = await fetch(`${ctx.rbUrl()}/api/shops/${route.params.id}`,
      {
        method: "DELETE",
        headers: ctx.authHeaders()
      }
    );
    if (res.ok) {
      router.push("/");
    } else {
      var err = await res.text();
      alert(err);
    }
  }

</script>

<template>
  <nav>
    <button class="btn btn-info btn-sm" v-if="isOwner" @click="router.push(`/edit-good/${route.params.id}`);">Добавить товар / услугу</button>
    &nbsp;
    <RouterLink class="btn btn-info btn-sm" v-if="isOwner" v-bind:to="`/edit-shop/${shop.Ref}`">
      Редактировать витрину
    </RouterLink>
    &nbsp;
    <button class="btn btn-info btn-sm" v-if="isOwner" @click="DeleteShop();">Удалить витрину</button>
  </nav>

  <h1>{{ shop.Caption }}</h1>
  <div class="col">&nbsp;</div>
  <div v-if="shop.Description" class="row-mb3">
    <div class="col"><span v-html="shopDescription"></span></div>
    <div class="col">&nbsp;</div>
  </div>
  <div class="row">
    <table class="table">
      <tbody>
        <tr v-for="good in goods.Data" v-bind:key="good.Ref">
          <td class="col-9">
            <div class="row">
              <div class="col-4">
                <RouterLink v-if="good.thumb" v-bind:to="`/good/${good.Ref}`"><img :src="good.thumb" class="img-fluid img-thumbnail" height="60" width="60"></RouterLink>
              </div>
              <div class="col-8">
                &nbsp;<RouterLink v-bind:to="`/good/${good.Ref}`">{{ good.Caption }}</RouterLink>
              </div>
            </div>
          </td>
          <td class="col-3">
            {{ good.Price }}<br />
              <button class="btn btn-primary btn-sm" @click="Inc(good)">+</button>&nbsp;
              <span v-if="good.basked"><button class="btn btn-primary btn-sm" @click="Dec(good)">-</button>&nbsp;{{ good.basked }}</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

</template>
