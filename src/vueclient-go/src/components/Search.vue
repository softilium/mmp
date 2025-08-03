<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { ref, watch } from "vue";
import { ctx } from "./ctx.js";

const q = ref("");
const goods = ref([]);

watch(
  () => q.value,
  async () => {
    setTimeout(async () => {
      if (q.value.length < 2) {
        goods.value = [];
        return;
      }
      let res = await fetch(ctx.rbUrl() + "/api/goods/search?q=" + q.value);
      if (res.ok) {
        goods.value = await res.json();
      } else {
        goods.value = [];
      }
    }, 1000);
  }
);
</script>

<template>
  <h1>Поиск</h1>

  <div class="row mb-3">
    <div class="col">
      <input type="text" class="form-control" v-model="q" />
    </div>
  </div>
  <table class="table table-striped">
    <tbody>
      <tr v-for="good in goods" :key="good.Ref">
        <td>
          <RouterLink :to="'/good/' + good.Ref">{{ good.Caption }}</RouterLink>
        </td>
      </tr>
    </tbody>
  </table>
</template>
