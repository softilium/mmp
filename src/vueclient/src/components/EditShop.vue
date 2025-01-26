<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();
  const router = useRouter();

  let title = "Новая витрина";
  if (route.params.id) title = "Редактировать витрину";

  let shop = ref({ id: 0, caption: "", description: "", deliveryConditions: "", createdByInfo: { id: 0, userName: "" } })

  let shopId = null;

  onMounted(async () => {
    shopId = route.params.id;
    if (shopId != null) {
      try {
        let res = await fetch(authStore.rbUrl() + "/api/shops/" + shopId, {
          method: "GET"
        });
        if (res.ok) {
          res = await res.json();
          shop.value = res;
          shop.value.createdByInfo = null;
        }
      } catch (err) { console.log(err); };
    }
  });

  const Save = async () => {

    if (!authStore.userInfo.id) return;

    if (shopId == null)
      try {
        let res = await fetch(authStore.rbUrl() + "/api/shops", {
          method: "POST",
          headers: authStore.authHeadersAppJson(),
          body: JSON.stringify(shop.value)
        });

        if (res.ok) {
          res = await res.json();
          shop.value.id = res.id;
          router.push(`/shop/${shop.value.id}`);
        }
      } catch (err) { console.log(err); }
    else
      try {
        let res = await fetch(authStore.rbUrl() + "/api/shops/" + shopId, {
          method: "PUT",
          headers: authStore.authHeadersAppJson(),
          body: JSON.stringify(shop.value)
        });

        if (res.ok)
          router.push(`/shop/${shop.value.id}`);

      } catch (err) { console.log(err); }


  }

</script>

<template>

  <h1>{{ title }}</h1>
  <div>&nbsp;</div>
  <button class="btn btn-primary btn-sm" @click="Save">Сохранить</button>
  <div>&nbsp;</div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Название</label>
      <div class="col-7">
        <input class="form-control" v-model="shop.caption" required />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Подробное описание витрины</label>
      <div class="col-7">
        <textarea class="form-control" v-model="shop.description" rows="5" placeholder="Это описание будет показано вверху под заголовком витрины." />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Условия доставки, возврата и проч. (показывается при формировании заказа)</label>
      <div class="col-7">
        <textarea class="form-control" v-model="shop.deliveryConditions" rows="5" placeholder="Этот текст с условиями, сроками доставки, условиями возврата и другой важной информацией будет показан на странице перед заказом" />
      </div>
    </div>
  </div>

</template>
