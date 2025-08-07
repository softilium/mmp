<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { ctx } from "./ctx.js";

const route = useRoute();

const goods = ref([]);
const tagName = ref("");

onMounted(async () => {
  let rt = await fetch(ctx.rbUrl() + "/api/tags?ref=" + route.params.id);
  if (rt.ok) {
    let rtt = await rt.json();
    tagName.value = rtt.Name;
  }

  let res = await fetch(
    ctx.rbUrl() + "/api/goods-by-tag?ref=" + route.params.id
  );
  if (res.ok) {
    goods.value = await res.json();
  }
});
</script>

<template>
  <h1>{{ tagName }}</h1>

  <div class="row">
    <table class="table">
      <tbody>
        <tr v-for="good in goods" v-bind:key="good.Ref">
          <td class="col">
            &nbsp;<RouterLink v-bind:to="`/good/${good.Ref}`">{{
              good.Caption
            }}</RouterLink>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
