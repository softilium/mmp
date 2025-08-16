<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ctx } from "./ctx.js";
import GoodList from "./GoodList.vue";

const route = useRoute();
const router = useRouter();

const shop = ref({
  Ref: "",
  Caption: "",
  Description: "",
  CreatedBy: { Ref: "", Username: "" },
});
const shopDescription = ref("");
const tags = ref([{ tagRef: "", tagName: "", tagged: false, tagColor: "" }]);
const goods = ref({ Data: [], PagesCount: 0 });
const isOwner = ref(false);

onMounted(async () => {
  try {
    let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + route.params.id);
    if (res.ok) {
      shop.value = await res.json();
      isOwner.value =
        shop.value.CreatedBy.Ref == ctx.userInfo.id && ctx.userInfo.shopManage;
      shopDescription.value = ctx.linkify(shop.value.Description);

      let r2 = await fetch(
        ctx.rbUrl() + "/api/tags-by-shop?ref=" + shop.value.Ref
      );
      if (r2.ok) {
        tags.value = await r2.json();
      }
    } else router.push("/");
  } catch (err) {
    console.log(err);
    router.push("/");
  }

  try {
    let res = await fetch(
      ctx.rbUrl() + "/api/goods?shopref=" + route.params.id,
      { signal: AbortSignal.timeout(5000) }
    );
    if (res.ok) {
      goods.value = await res.json();
    }
  } catch (err) {
    console.log(err);
  }
});

const DeleteShop = async () => {
  if (!confirm("Удалить витрину, вы уверены?")) return;

  let res = await fetch(`${ctx.rbUrl()}/api/shops?ref=${route.params.id}`, {
    method: "DELETE",
    headers: await ctx.authHeaders(),
  });
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    router.push("/");
  } else {
    var err = await res.text();
    alert(err);
  }
};
</script>

<template>
  <nav>
    <button
      class="btn btn-info btn-sm"
      v-if="isOwner"
      @click="router.push(`/edit-good/${route.params.id}`)"
    >
      Добавить товар / услугу
    </button>
    &nbsp;
    <RouterLink
      class="btn btn-info btn-sm"
      v-if="isOwner"
      v-bind:to="`/edit-shop/${shop.Ref}`"
    >
      Редактировать витрину
    </RouterLink>
    &nbsp;
    <button class="btn btn-info btn-sm" v-if="isOwner" @click="DeleteShop()">
      Удалить витрину
    </button>
  </nav>

  <h1>{{ shop.Caption }}</h1>
  <div class="col">&nbsp;</div>
  <div v-if="shop.Description" class="row-mb3">
    <div class="col"><span v-html="shopDescription"></span></div>
    <div class="col">&nbsp;</div>
  </div>

  <div class="row">
    <div class="col text-center">
      <span v-for="(tag, idx) in tags" v-bind:key="tag.tagRef">
        <span :style="`padding: 3px; font-size: ${1.4 / (idx * 0.02 + 1)}em`">
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
  <GoodList :goods="goods.Data" />
</template>
