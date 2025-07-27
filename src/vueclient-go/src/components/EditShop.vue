<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { ctx } from './ctx.js';

  const route = useRoute();
  const router = useRouter();

  let title = "Новая витрина";
  if (route.params.id) title = "Редактировать витрину";

  let shop = ref({ Ref: "", Caption: "", Description: "", DeliveryConditions: "", CreatedBy: { Ref: "", Username: "", Description: "" } });

  let shopId = null;

  onMounted(async () => {
    shopId = route.params.id;
    if (shopId != null) {
      try {
        let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + shopId, {
          method: "GET"
        });
        if (res.ok) {
          res = await res.json();
          shop.value = res;
        }
      } catch (err) { console.log(err); };
    }
  });

  const Save = async () => {
    if (!ctx.userInfo.id) return;
    if (shopId == null)
      try {
        let res = await fetch(ctx.rbUrl() + "/api/shops", {
          method: "POST",
          headers: await ctx.authHeadersAppJson(),
          body: JSON.stringify(shop.value)
        });
        if (await ctx.CheckUnauth(res)) return;
        if (res.ok) {
          res = await res.json();
          shop.value.Ref = res.Ref;
          router.push(`/shop/${shop.value.Ref}`);
        }
      } catch (err) { console.log(err); }
    else
      try {
        let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + shopId, {
          method: "PUT",
          headers: await ctx.authHeadersAppJson(),
          body: JSON.stringify(shop.value)
        });
        if (await ctx.CheckUnauth(res)) return;
        if (res.ok)
          router.push(`/shop/${shop.value.Ref}`);

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
        <input class="form-control" v-model="shop.Caption" required />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Подробное описание витрины</label>
      <div class="col-7">
        <textarea class="form-control" v-model="shop.Description" rows="5" placeholder="Это описание будет показано вверху под заголовком витрины." />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Условия доставки, возврата и проч. (показывается при формировании заказа)</label>
      <div class="col-7">
        <textarea class="form-control" v-model="shop.DeliveryConditions" rows="5" placeholder="Этот текст с условиями, сроками доставки, условиями возврата и другой важной информацией будет показан на странице перед заказом" />
      </div>
    </div>
  </div>

</template>
