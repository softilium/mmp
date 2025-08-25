<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ctx } from "./ctx.js";

const route = useRoute();
const router = useRouter();

let title = "Новая витрина";
if (route.params.id) title = "Редактировать витрину";

let shop = ref({
  Ref: "",
  Caption: "",
  Description: "",
  DeliveryConditions: "",
  DiscountPercent: 0,
  CreatedBy: { Ref: "", Username: "", Description: "" },
});

let shopId = null;
let imageSrc = ref("");

const isImageLoading = ref(true);

onMounted(async () => {
  shopId = route.params.id;
  if (shopId != null) {
    try {
      let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + shopId, {
        method: "GET",
      });
      if (res.ok) {
        let rj = await res.json();
        shop.value = rj;
      }
    } catch (err) {
      console.log(err);
    }
    let res = await fetch(
      `${ctx.rbUrl()}/api/shops/images?ref=${route.params.id}&n=0`,
      { method: "GET" }
    );
    if (res.status == 200) {
      // status 204 also ok but it means no image
      let b = await res.blob();
      const src = URL.createObjectURL(b);
      imageSrc.value = src;
    }
  }
  isImageLoading.value = false;
});

const Save = async () => {
  if (!ctx.userInfo.id) return;
  if (shopId == null)
    try {
      let res = await fetch(ctx.rbUrl() + "/api/shops", {
        method: "POST",
        headers: await ctx.authHeadersAppJson(),
        body: JSON.stringify(shop.value),
      });
      if (await ctx.CheckUnauth(res)) return;
      if (res.ok) {
        let rj = await res.json();
        shop.value.Ref = rj.Ref;
        router.push(`/shop/${shop.value.Ref}`);
      }
    } catch (err) {
      console.log(err);
    }
  else
    try {
      let res = await fetch(ctx.rbUrl() + "/api/shops?ref=" + shopId, {
        method: "PUT",
        headers: await ctx.authHeadersAppJson(),
        body: JSON.stringify(shop.value),
      });
      if (await ctx.CheckUnauth(res)) return;
      if (res.ok) router.push(`/shop/${shop.value.Ref}`);
    } catch (err) {
      console.log(err);
    }

  if (!imageSrc.value) {
    let res = await fetch(`${ctx.rbUrl()}/api/shops/images?ref=${shopId}&n=0`, {
      method: "DELETE",
      headers: await ctx.authHeaders(),
    });
    if (await ctx.CheckUnauth(res)) return;
    if (!res.ok) console.log(res);
  } else {
    let blob = await fetch(imageSrc.value).then((r) => r.blob()); // load image from blob url
    let data = new FormData();
    data.append("image", blob);
    let res = await fetch(`${ctx.rbUrl()}/api/shops/images?ref=${shopId}&n=0`, {
      method: "POST",
      headers: await ctx.authHeaders(),
      body: data,
    });
    if (await ctx.CheckUnauth(res)) return;
    if (!res.ok) console.log(res);
  }
};

const handelFileUpload = (e) => {
  var files = e.target.files || e.dataTransfer.files;
  if (!files.length) return;

  for (let i = 0; i < files.length; i++) {
    const src = URL.createObjectURL(files[i]);
    imageSrc.value = src;
  }
};

const removeImage = () => {
  imageSrc.value = "";
};
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
        <textarea
          class="form-control"
          v-model="shop.Description"
          rows="5"
          placeholder="Это описание будет показано вверху под заголовком витрины."
        />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Скидка, %</label>
      <div class="col-7">
        <input
          type="number"
          class="form-control"
          v-model="shop.DiscountPercent"
        />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label"
        >Условия доставки, возврата и проч. (показывается при формировании
        заказа)</label
      >
      <div class="col-7">
        <textarea
          class="form-control"
          v-model="shop.DeliveryConditions"
          rows="5"
          placeholder="Этот текст с условиями, сроками доставки, условиями возврата и другой важной информацией будет показан на странице перед заказом"
        />
      </div>
    </div>
  </div>

  <div class="row">
    <div class="col">
      <input
        type="file"
        accept="image/*"
        @change="(event) => handelFileUpload(event)"
      />
      <button
        v-if="imageSrc"
        @click="removeImage()"
        class="btn btn-secondary btn-sm"
      >
        Удалить картинку
      </button>
    </div>
  </div>

  <div v-if="!isImageLoading" class="row mb-3">
    <div class="col-11">
      <img :src="imageSrc" class="d-block w-100" />
    </div>
  </div>
</template>
