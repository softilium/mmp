<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const good = ref({ ownerShop: { id: 0 }, caption: "", description: "" });

  onMounted(async () => {
    if (route.params.id) {
      try {
        let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id);
        if (res.ok) {
          good.value = await res.json();
        } else router.push("/shop/" + route.params.shopid);
      } catch (err) {
        console.log(err);
        router.push("/shop/" + route.params.shopid);
      }
    }
  });

  const Save = async () => {

    if (route.params.id) {
      let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        },
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        router.push("/shop/" + route.params.shopid);
      }

    } else {
      good.value.ownerShop.id = route.params.shopid;

      let res = await fetch(authStore.rbUrl() + "/api/goods", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        },
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        router.push("/shop/" + route.params.shopid);
      }
    }
  };

</script>

<template>
  <h1>Редактирование товара / услуги</h1>
  <div>&nbsp;</div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Название</label>
      <div class="col-7">
        <input class="form-control" v-model="good.caption" required />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Описание</label>
      <div class="col-7">
        <textarea class="form-control" v-model="good.description" rows="5" />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Цена</label>
      <div class="col-7">
        <input class="form-control" v-model="good.price" />
      </div>
    </div>
  </div>
  <button class="btn btn-primary" @click="Save">Сохранить</button>
</template>
