<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const good = ref({ ownerShop: { id: 0 }, caption: "", description: "", article: "", url: "", price: 0, orderInShop: 100 });
  const isImageLoading = ref(true);

  onMounted(async () => {
    if (route.params.id) {
      try {
        let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id);
        if (res.ok) {
          good.value = await res.json();
          await LoadImages();
        } else router.push("/shop/" + route.params.shopid);
      } catch (err) {
        console.log(err);
        router.push("/shop/" + route.params.shopid);
      }
    } else isImageLoading.value = false;
  });

  const LoadImages = async () => {
    for (let i = 0; i < maxImagesCnt.value; i++) {
      let res = await fetch(`${authStore.rbUrl()}/api/goods/images/${route.params.id}/${i}`, { method: "GET" });
      if (res.status == 200) { // status 204 also ok but it means no image
        let b = await res.blob();
        const src = URL.createObjectURL(b);
        imageSrc.value.push(src);
      }
    }
    isImageLoading.value = false;
  }

  const SaveImages = async (gid) => {
    for (let i = 0; i <= maxImagesCnt.value - 1; i++) {
      if (!imageSrc.value[i]) {
        let res = await fetch(`${authStore.rbUrl()}/api/goods/images/${gid}/${i}`, {
          method: "DELETE",
          headers: authStore.authHeaders()
        });
        if (!res.ok) console.log(res);
      } else {
        let blob = await fetch(imageSrc.value[i]).then(r => r.blob()); // load image from blob url
        let data = new FormData();
        data.append("image", blob);
        let res = await fetch(`${authStore.rbUrl()}/api/goods/images/${gid}/${i}`, {
          method: "POST",
          headers: authStore.authHeaders(),
          body: data
        });
        if (!res.ok) console.log(res);
      }
    }
  }

  const Save = async () => {
    if (route.params.id) {
      let res = await fetch(authStore.rbUrl() + "/api/goods/" + route.params.id, {
        method: "PUT",
        headers: authStore.authHeadersAppJson(),
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        await SaveImages(route.params.id);
        router.push("/shop/" + route.params.shopid);
      }
    } else {
      good.value.ownerShop.id = route.params.shopid;
      let res = await fetch(authStore.rbUrl() + "/api/goods", {
        method: "POST",
        headers: authStore.authHeadersAppJson(),
        body: JSON.stringify(good.value)
      });
      if (res.ok) {
        res = await res.json();
        await SaveImages(res.id);
        router.push("/shop/" + route.params.shopid);
      }
    }
  };

  //images handling

  const maxImagesCnt = ref(3);
  const imageSrc = ref([]);
  let curImgIndex = ref(0);

  const handelFileUpload = (e) => {
    curImgIndex.value = 0;
    var files = e.target.files || e.dataTransfer.files;
    if (!files.length) return;

    for (let i = 0; i < files.length; i++) {
      const src = URL.createObjectURL(files[i]);
      imageSrc.value.push(src);

    }
    curImgIndex.value = imageSrc.value.length - 1;
  };

  const removeItem = (index) => {
    imageSrc.value.splice(index, 1);
    curImgIndex.value = 0;
  };


</script>

<template>
  <h1>Редактирование товара / услуги</h1>
  <div>&nbsp;</div>
  <button class="btn btn-primary btn-sm" @click="Save">Сохранить</button>
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
      <label class="col-3 form-label">Артикул</label>
      <div class="col-7">
        <input class="form-control" v-model="good.article" />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Ссылка</label>
      <div class="col-7">
        <input class="form-control" v-model="good.url" />
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

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Порядок на витрине (меньше = выше)</label>
      <div class="col-7">
        <input class="form-control" v-model="good.orderInShop" />
      </div>
    </div>
  </div>

  <div class="row mb-3">
    <div class="col-1">
      <span v-for="(src, index) in imageSrc" :key="index">
        <button :class="`${index == curImgIndex ? 'btn btn-secondary btn-sm' : 'btn btn-outline-secondary btn-sm'}`"
                @click="curImgIndex = index">
          {{ index + 1 }}
        </button>&nbsp;
      </span>
    </div>
    <div class="col">
      <input v-if="imageSrc.length < maxImagesCnt" type="file" accept="image/*" @change="(event) => handelFileUpload(event)" />
    </div>
  </div>

  <div v-if="!isImageLoading && imageSrc.length > 0" class="row mb-3">
    <div class="col-1">
      <button class="btn btn-secondary btn-sm" @click="removeItem(curImgIndex)">Удалить</button>
    </div>
    <div class="col-11">
      <img :src="imageSrc[curImgIndex]" class="d-block w-100">
    </div>
  </div>

</template>
