<script setup lang="ts">
import { nextTick, watch } from "vue";
import { ctx } from "./ctx.js";
import { localBasket } from "../services/localBasket";

interface Good {
  Ref: string;
  Caption: string;
  Price: number;
  thumb: string;
  basked?: number | null;
}

const props = defineProps({
  goods: {
    type: Array<Good>,
    required: true,
  },
});

watch(
  () => props.goods,
  async (newValue) => {
    if (newValue.length > 0) {
      props.goods.forEach(async (g) => {
        let res = await fetch(
          `${ctx.rbUrl()}/api/goods/thumbs?ref=${g.Ref}&n=0`,
          { method: "GET" }
        );
        if (res.status == 200) {
          // status 204 means no image
          let b = await res.blob();
          const src = URL.createObjectURL(b);
          g.thumb = src;
        }
      });
      await nextTick();
      await LoadBasket();
    }
  },
  { immediate: true }
);

const LoadBasket = async () => {
  if (ctx.userInfo.id) {
    try {
      let res = await fetch(ctx.rbUrl() + "/api/basket", {
        method: "GET",
        headers: await ctx.authHeadersAppJson(),
      });
      if (await ctx.CheckUnauth(res)) return;
      if (res.ok) {
        let rj = await res.json();

        if (rj.Data.length > 0) {
          let goodmap = new Map();
          props.goods.forEach((_) => {
            goodmap.set(_.Ref, _);
          });

          rj.Data.forEach((_) => {
            let goodId = _.Good.Ref;
            let gObj = goodmap.get(goodId);
            if (gObj != null) {
              gObj.basked = _.Qty;
            }
          });
        }
      } else {
        if (res.status == 401) {
          await ctx.CheckLogged();
        }
      }
    } catch (err) {
      console.log(err);
    }
  } else {
    // Anonymous users: use localBasket
    const items = localBasket.getItems();
    props.goods.forEach((good) => {
      const found = items.find((i) => i.goodId === good.Ref);
      if (found) {
        good.basked = found.quantity;
      } else {
        good.basked = null;
      }
    });
  }
};

const Inc = async (good) => {
  if (!ctx.userInfo.id) {
    localBasket.addItem({
      goodId: good.Ref,
      quantity: 1,
      price: good.Price,
    });
    await ctx.loadBasket();
    LoadBasket();
    return;
  }
  let res = await fetch(
    ctx.rbUrl() + "/api/basket/increase?goodref=" + good.Ref,
    {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
    }
  );
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    if (!good.basked) good.basked = 0;
    good.basked++;
    await ctx.loadBasket();
    LoadBasket();
  }
};

const Dec = async (good) => {
  if (!ctx.userInfo.id) {
    localBasket.decItem(good.Ref);
    LoadBasket();
    await ctx.loadBasket();
    return;
  }
  let res = await fetch(
    ctx.rbUrl() + "/api/basket/decrease?goodref=" + good.Ref,
    {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
    }
  );
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    if (good.basked && good.basked > 0) good.basked--;
    if (good.basked == 0) good.basked = null;
    LoadBasket();
    await ctx.loadBasket();
  }
};
</script>

<template>
  <div class="row">
    <table class="table">
      <tbody>
        <tr v-for="good in goods" :key="good.Ref">
          <td class="col-9">
            <div class="row">
              <div class="col-4">
                <RouterLink v-if="good.thumb" :to="`/good/${good.Ref}`">
                  <img
                    :src="good.thumb"
                    class="img-fluid img-thumbnail"
                    height="60"
                    width="60"
                  />
                </RouterLink>
              </div>
              <div class="col-8">
                &nbsp;<RouterLink :to="`/good/${good.Ref}`">{{
                  good.Caption
                }}</RouterLink>
              </div>
            </div>
          </td>
          <td class="col-3">
            {{ good.Price }}<br />
            <button class="btn btn-primary btn-sm" @click="Inc(good)">+</button
            >&nbsp;
            <span v-if="good.basked">
              <button class="btn btn-primary btn-sm" @click="Dec(good)">
                -</button
              >&nbsp;{{ good.basked }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
