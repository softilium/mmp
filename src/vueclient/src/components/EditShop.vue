<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();
  const router = useRouter();

  let title = "New shop";
  if (route.params.id) title = "Edit shop";

  let captionField = ref("");
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
          captionField.value = res.caption;
        }
      } catch (err) { console.log(err); };
    }
  });

  const Save = async () => {

    if (!authStore.loggedEmail) return;

    if (shopId == null)
      try {
        let res = await fetch(authStore.rbUrl() + "/api/shops", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          },
          body: JSON.stringify({ caption: captionField.value, createdBy: {} })
        });

        if (res.ok) router.push("/");
      } catch (err) { console.log(err); }
    else
      try {
        let res = await fetch(authStore.rbUrl() + "/api/shops/" + shopId, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          },
          body: JSON.stringify({ caption: captionField.value, createdBy: {} })
        });

        if (res.ok) router.push("/");
      } catch (err) { console.log(err); }

  }

</script>

<template>

  <h1>{{ title }}</h1>

  <div class="row">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Caption</label>
      <div class="col-7">
        <input class="form-control" v-model="captionField" required />
      </div>
    </div>
  </div>

  <div>&nbsp;</div>
  <button class="btn btn-primary" @click="Save">Сохранить</button>

</template>
