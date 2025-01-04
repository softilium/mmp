<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute()
  const router = useRouter()

  const captionField = ref("");
  const descriptionField = ref("");
  const good = ref({ ownerShop: { id: 0 }, caption: "", description: "" });

  const Save = async () => {

    good.value.ownerShop.id = route.params.shopid;
    good.value.caption = captionField;
    good.value.description = descriptionField;

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
  };
    
</script>

<template>

  <div>
    Caption: <input v-model="captionField" />
    Description: <input v-model="descriptionField" />
  </div>

  <button @click="Save">Save and return to shop</button>

</template>
