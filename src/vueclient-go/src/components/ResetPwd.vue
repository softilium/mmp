<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ctx } from "./ctx.js";

const route = useRoute();
const router = useRouter();

const emailString = ref("");
const newPasswordString = ref("");
const confirmPasswordString = ref("");
const resetCode = ref("");

onMounted(() => {
  resetCode.value = route.query.code as string;
  if (!resetCode.value) {
    router.push("/");
  }
});

const setNewPwd = async () => {
  if (newPasswordString.value != confirmPasswordString.value) {
    alert("Пароли не совпадают");
    return;
  }
  if (newPasswordString.value.length < 4) {
    alert("Пароль должен быть не менее 6 символов");
    return;
  }
  if (resetCode.value == "") {
    alert("Код сброса пароля не указан");
    return;
  }
  let q = await fetch(
    `${ctx.rbUrl()}/api/users/resetpwd?email=${emailString.value}&pwd=${
      newPasswordString.value
    }&code=${resetCode.value}`,
    { method: "POST" }
  );
  if (q.ok) {
    alert("Пароль успешно изменен");
    router.push("/login");
  } else {
    alert("Ошибка при смене пароля");
  }
};
</script>

<template>
  <div class="row">
    <div class="col">
      <h1>Сброс пароля</h1>
    </div>
  </div>

  <div class="row mb-3">
    <div class="col-4">
      <label class="form-label">Email</label>
    </div>
    <div class="col-4">
      <input class="form-control" v-model="emailString" required type="email" />
    </div>
  </div>

  <div class="row mb-3">
    <div class="col-4">
      <label class="form-label">Новый пароль</label>
    </div>
    <div class="col-4">
      <input
        class="form-control"
        type="password"
        v-model="newPasswordString"
        required
      />
    </div>
  </div>

  <div class="row mb-3">
    <div class="col-4">
      <label class="form-label">Подтверждение пароля</label>
    </div>
    <div class="col-4">
      <input
        class="form-control"
        type="password"
        v-model="confirmPasswordString"
        required
      />
    </div>
  </div>

  <div class="row mb-3">
    <div class="col-8">
      <button class="btn btn-primary btn-sm" @click="setNewPwd">
        Указать новый пароль
      </button>
    </div>
  </div>
</template>
