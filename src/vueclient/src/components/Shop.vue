<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const goods = ref([]);
  const shop = ref({});

  const isOwner = ref(false);

  onMounted(async () => {

    try {
      let res = await fetch(authStore.rbUrl() + "/api/shops/" + route.params.id);
      if (res.ok) {
        shop.value = await res.json();
        isOwner.value = shop.value.createdBy.email == authStore.loggedEmail;
      } else router.push("/");
    } catch (err) {
      console.log(err);
      router.push("/");
    }

    try {
      let res = await fetch(authStore.rbUrl() + "/api/goods?shopid="+route.params.id, { signal: AbortSignal.timeout(5000), });
      if (res.ok) {
        goods.value = await res.json();
      }
    } catch (err) { console.log(err); }
  });
    
</script>

<template>

  <h1>{{ shop.caption }}</h1>

  <div class="row">
    <table class="table table-sm">
      <thead>
        <tr>
          <td>Товар или Услуга</td>
          <td class="text-end">Цена</td>
        </tr>
      </thead>
      <tbody>
        <tr v-for="good in goods">
          <td>{{ good.caption }}</td>
          <td class="text-end">{{ good.price }}</td>
          <td><RouterLink v-if="isOwner" v-bind:to="`/edit-good/${shop.id}/${good.id}`">[edit]</RouterLink></td>
        </tr>
      </tbody>
    </table>
  </div>

  <br />
  <button class="btn btn-primary" v-if="isOwner" @click="router.push(`/edit-good/${route.params.id}`);">Добавить товар / услугу</button>

</template>
