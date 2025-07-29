<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ctx } from "./ctx.js";
import ProfileLink from "./ProfileLink.vue";

let shops = ref({ Data: [], PagesCount: 0 });

onMounted(async () => {
  try {
    let res = await fetch(ctx.rbUrl() + "/api/shops", {
      signal: AbortSignal.timeout(5000),
    });
    if (res.ok) {
      shops.value = await res.json();
    }
  } catch (err) {
    console.log(err);
  }
});
</script>

<template>
  <nav>
    <button
      class="btn btn-info btn-sm"
      v-if="ctx.userInfo.id && ctx.userInfo.shopManage"
      @click="$router.push('/edit-shop')"
    >
      Добавить витрину
    </button>
  </nav>

  <h1>Витрины</h1>

  <div class="row">
    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Заголовок</th>
          <th>Владелец</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in shops.Data" v-bind:key="item.Ref">
          <td class="col-9">
            <RouterLink v-bind:to="`/shop/${item.Ref}`">{{
              item.Caption
            }}</RouterLink>
          </td>
          <td class="col-2">
            <ProfileLink :userInfo="item.CreatedBy"></ProfileLink>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <br />
</template>
