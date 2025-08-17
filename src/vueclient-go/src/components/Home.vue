<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ctx } from "./ctx.js";

let shops = ref({ Data: [], PagesCount: 0, thumb: "" });
const tags = ref([{ tagRef: "", tagName: "", tagged: false, tagColor: "" }]);

onMounted(async () => {
  try {
    let res = await fetch(ctx.rbUrl() + "/api/shops", {
      signal: AbortSignal.timeout(5000),
    });
    if (res.ok) {
      shops.value = await res.json();

      shops.value.Data.forEach(async (shop) => {
        let res = await fetch(
          `${ctx.rbUrl()}/api/shops/thumbs?ref=${shop.Ref}&n=0`,
          { method: "GET" }
        );
        if (res.status == 200) {
          // status 204 means no image
          let b = await res.blob();
          const src = URL.createObjectURL(b);
          shop.thumb = src;
        }
      });
    }

    let r2 = await fetch(ctx.rbUrl() + "/api/tags-by-all");
    if (r2.ok) {
      tags.value = await r2.json();
    }
  } catch (err) {
    console.log(err);
  }
});
</script>

<template>
  <div class="row">
    <div class="col text-center" style="opacity: 0.8">
      <span v-for="tag in tags" v-bind:key="tag.tagRef">
        <span :style="`padding: 3px; font-size: ${1}em`">
          <RouterLink :to="`/goods-by-tag/${tag.tagRef}`"
            ><span :class="['badge', tag.tagColor]">
              {{ tag.tagName }}
            </span></RouterLink
          >
        </span>
      </span>
    </div>
  </div>

  <div class="col">&nbsp;</div>

  <nav>
    <button
      class="btn btn-info btn-sm"
      v-if="ctx.userInfo.id && ctx.userInfo.shopManage"
      @click="$router.push('/edit-shop')"
    >
      Добавить витрину
    </button>
  </nav>
  <h1>&nbsp;</h1>

  <div class="row">
    <table class="table table-sm">
      <tbody>
        <tr v-for="item in shops.Data" v-bind:key="item.Ref">
          <td>
            <RouterLink v-if="item.thumb" :to="`/shop/${item.Ref}`">
              <img
                :src="item.thumb"
                class="img-fluid img-thumbnail"
                height="90"
                width="90"
              />
            </RouterLink>
          </td>
          <td class="col-9">
            <h5>
              <RouterLink v-bind:to="`/shop/${item.Ref}`">{{
                item.Caption
              }}</RouterLink>
            </h5>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <br />
</template>
