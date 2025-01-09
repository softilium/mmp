<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { authStore } from './authStore.js';

  const users = ref([]);

  onMounted(async () => {
    let res = await fetch(authStore.rbUrl() + "/api/admin/allusers",
      {
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });
    if (res.ok) {
      users.value = await res.json();
    }
  });

  const Save = async (id) => {
    users.value.forEach(async (_) => {
      if (_.id == id) {
        await fetch(authStore.rbUrl() + "/api/admin/allusers/" + _.id,
          {
            method: "PUT",
            body: JSON.stringify(_),
            headers: {
              "Content-Type": "application/json",
              "Authorization": "Bearer " + authStore.accessToken
            }
          });
      }
    });
  }

</script>

<template>

  <h1>Admin</h1>

  <table v-if="users.length > 0" class="table table-sm">

    <thead>
      <tr>
        <th>Имя</th>
        <th>Email</th>
        <th>Админ</th>
        <th>Менеджер витрин</th>
      </tr>
    </thead>

    <tbody>

      <tr v-for="u in users">
        <td>{{ u.userName }}</td>
        <td>{{ u.email }}</td>
        <td> <input v-model="u.admin" type="checkbox" /></td>
        <td><input v-model="u.shopManage" type="checkbox" /></td>
        <td><button class="btn btn-primary" @click="Save(u.id)">save</button></td>
      </tr>

    </tbody>

  </table>

</template>
