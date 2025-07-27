<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { ctx } from './ctx.js';

  const users = ref({Data:[], PagesCount: 0});

  onMounted(async () => {
    let res = await fetch(ctx.rbUrl() + "/api/admin/allusers",
      {
        headers: await ctx.authHeadersAppJson()
      });
    if (await ctx.CheckUnauth(res)) return;
    if (res.ok) {
      users.value = await res.json();
    }
  });

  const Save = async (Ref) => {
    users.value.Data.forEach(async (u) => {
      if (u.Ref == Ref) {
        let res = await fetch(ctx.rbUrl() + "/api/admin/allusers?ref=" + u.Ref,
          {
            method: "PUT",
            body: JSON.stringify(u),
            headers: await ctx.authHeadersAppJson()
          });
        if (await ctx.CheckUnauth(res)) return;
      }
    });
  }

  const Migrate = async () => {

    if (confirm("Start migration?") == false) return;

    let res = await fetch(ctx.rbUrl() + "/api/admin/migrate", {
      method: "POST",
      headers: await ctx.authHeadersAppJson()
    });
    if (await ctx.CheckUnauth(res)) return;
  };

</script>

<template>

  <h1>Admin</h1>

    <button @click="Migrate()">Migrate</button>

  <table v-if="users.Data.length > 0" class="table table-sm">
    <thead>
      <tr>
        <th>Имя</th>
        <th>Email</th>
        <th>Админ</th>
        <th>Менеджер витрин</th>
        <th>Telegram активен</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="u in users.Data" v-bind:key="u.Ref">
        <td>{{ u.Username }}</td>
        <td>{{ u.Email }}</td>
        <td><input v-model="u.Admin" type="checkbox" /></td>
        <td><input v-model="u.ShopManager" type="checkbox" /></td>
        <td><input v-model="u.TelegramVerified" type="checkbox" disabled /></td>
        <td><button class="btn btn-primary" @click="Save(u.Ref)">save</button></td>
      </tr>
    </tbody>
  </table>

</template>
