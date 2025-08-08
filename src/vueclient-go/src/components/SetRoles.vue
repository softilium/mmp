<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ctx } from "./ctx.js";

const users = ref({ Data: [], PagesCount: 0 });

const tags = ref([
  { Ref: "1", Name: "Tag1", Color: "#FF0000" },
  { Ref: "2", Name: "Tag2", Color: "#00FF00" },
]);

const newTagName = ref("");

onMounted(async () => {
  await LoadTags();

  const r1 = await fetch(ctx.rbUrl() + "/api/admin/allusers", {
    headers: await ctx.authHeadersAppJson(),
  });
  if (await ctx.CheckUnauth(r1)) return;
  if (r1.ok) {
    users.value = await r1.json();
  }
});

const LoadTags = async () => {
  const r2 = await fetch(ctx.rbUrl() + "/api/tags", {
    headers: await ctx.authHeadersAppJson(),
  });
  if (r2.ok) {
    const r2t = await r2.json();
    tags.value = r2t.Data || [];
  } else {
    console.error("Failed to fetch tags:", r2.statusText);
  }
};

const SaveUser = async (Ref) => {
  users.value.Data.forEach(async (u) => {
    if (u.Ref == Ref) {
      let res = await fetch(ctx.rbUrl() + "/api/admin/allusers?ref=" + u.Ref, {
        method: "PUT",
        body: JSON.stringify(u),
        headers: await ctx.authHeadersAppJson(),
      });
      if (await ctx.CheckUnauth(res)) return;
    }
  });
};

const SaveTag = async (tag) => {
  const response = await fetch(ctx.rbUrl() + "/api/tags?ref=" + tag.Ref, {
    method: "PUT",
    body: JSON.stringify(tag),
    headers: await ctx.authHeadersAppJson(),
  });
  if (await ctx.CheckUnauth(response)) return;

  if (response.ok) {
    await LoadTags(); // Refresh tags after saving
  } else {
    console.error("Failed to save tag:", response.statusText);
  }
};

const AddTag = async () => {
  if (!newTagName.value) {
    alert("Please fill tag name.");
    return;
  }
  const newTag = { Name: newTagName.value };
  const r1 = await fetch(ctx.rbUrl() + "/api/tags", {
    method: "POST",
    body: JSON.stringify(newTag),
    headers: await ctx.authHeadersAppJson(),
  });
  if (await ctx.CheckUnauth(r1)) return;

  if (r1.ok) {
    await LoadTags(); // Refresh tags after adding a new one
  } else {
    console.error("Failed to add tag:", r1.statusText);
  }
};

const DeleteTag = async (tag) => {
  if (!tag || !tag.Ref) {
    alert("No tag selected for deletion.");
    return;
  }
  const r1 = await fetch(ctx.rbUrl() + "/api/tags?ref=" + tag.Ref, {
    method: "DELETE",
    headers: await ctx.authHeadersAppJson(),
  });
  if (await ctx.CheckUnauth(r1)) return;

  if (r1.ok) {
    await LoadTags(); // Refresh tags after deletion
  } else {
    console.error("Failed to delete tag:", r1.statusText);
  }
};
</script>

<template>
  <h1>Admin</h1>

  <h2>Users</h2>

  <table class="table table-sm">
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
        <td>
          <button class="btn btn-primary" @click="SaveUser(u.Ref)">save</button>
        </td>
      </tr>
    </tbody>
  </table>

  <h1>Tags</h1>

  <table>
    <thead>
      <tr>
        <th>Tag</th>
        <th>Color</th>
        <th></th>
        <th>Preview</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="tag in tags" v-bind:key="tag.Ref">
        <td><input v-model="tag.Name" /></td>
        <td>
          <select class="form-select form-control-sm" v-model="tag.Color">
            <option value="text-bg-primary">primary</option>
            <option value="text-bg-secondary">secondary</option>
            <option value="text-bg-success">success</option>
            <option value="text-bg-danger">danger</option>
            <option value="text-bg-warning">warning</option>
            <option value="text-bg-info">info</option>
            <option value="text-bg-light">light</option>
            <option value="text-bg-dark">dark</option>
          </select>
        </td>
        <td>
          <button class="btn btn-primary btn-sm" @click="SaveTag(tag)">
            Save</button
          >&nbsp;
          <button class="btn btn-primary btn-sm" @click="DeleteTag(tag)">
            Delete
          </button>
        </td>
        <td>
          &nbsp;
          <RouterLink :to="`/goods-by-tag/${tag.Ref}`"
            ><span :class="['badge', tag.Color]">
              {{ tag.Name }}
            </span>
          </RouterLink>
        </td>
      </tr>
    </tbody>
  </table>

  <h2>Add new tag</h2>
  <input v-model="newTagName" placeholder="Tag Name" required />
  <button class="btn btn-primary" @click="AddTag()">Add Tag</button>
</template>
