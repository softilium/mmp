<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute } from 'vue-router'
  import { authStore } from './authStore.js';

  const route = useRoute();

  const me = ref(false);
  const user = ref({ userName: "", email: "", telegramUserName: "", telegramVerified: false, telegramCheckCode: "", botChatId: 0 });
  const newTelegramUserName = ref("");
  const telegramVerifyCode = ref("");
  const result = ref("");

  const NewTelegramCode = async () => {

    if (me.value != true) return;

    let res = await fetch(authStore.rbUrl() + "/api/profiles/newtelegramcode", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      }
    });
    if (res.ok) Load();

  }
  const CheckTelegramCode = async () => {

    if (me.value != true) return;

    let res = await fetch(authStore.rbUrl() + "/api/profiles/checktelegramcode/" + telegramVerifyCode.value, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + authStore.accessToken
      }
    });
    if (res.ok) Load();

  }

  onMounted(async () => {
    me.value = !route.params.id;
    Load();
  });

  const Load = async () => {
    if (!me.value) {
      let res = await fetch(authStore.rbUrl() + "/api/profiles/" + route.params.id,
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (res.ok) {
        user.value = await res.json();
        me.value = !route.params.id || authStore.userInfo.userName == user.value.userName;
      }
    }
    else {
      let res = await fetch(authStore.rbUrl() + "/api/profiles/my",
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + authStore.accessToken
          }
        });
      if (res.ok) {
        user.value = await res.json();
        newTelegramUserName.value = user.value.telegramUserName;
      }
    }
  }

  const Save = async () => {
    if (!me.value) return;
    if (user.value.telegramUserName != newTelegramUserName.value) {
      user.value.telegramVerified = false;
      user.value.telegramCheckCode = "";
    }
    user.value.telegramUserName = newTelegramUserName.value;
    let res = await fetch(authStore.rbUrl() + "/api/profiles",
      {
        method: "PUT",
        body: JSON.stringify(user.value),
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + authStore.accessToken
        }
      });
    if (res.ok) {
      result.value = "Изменения записаны";
      await Load();
    }
    else
      result.value = await res.text();
  }

</script>

<template>

  <h1>
    Профиль пользователя {{ user.userName }}&nbsp;<button v-if="me" class="btn btn-outline-secondary"
      @click="authStore.Logout(); $router.push('/');">Выйти</button>
  </h1>
  <br />

  <div class="row mb-3">
    <label class="col-3 form-label">Имя пользователя</label>
    <div class="col-7">
      <input class="form-control" v-model="user.userName" required />
    </div>
  </div>
  <div class="row mb-3">
    <label class="col-3 form-label">Email/логин</label>
    <div class="col-7">
      <input class="form-control" v-model="user.email" readonly />
    </div>
  </div>

  <div v-if="me">
    <br />
    <h4 v-if="!user.telegramVerified">
      Мы используем Telegam для уведомлений. Ваше имя пользователя не передается другим пользователям сайта. Для настройки уведомлений вам нужно
      добавить чат для
      <a href="https://t.me/RiverStoresBot" target="_blank">нашего бота</a>. После этого укажите ваше имя пользователя из Telegram ниже.
    </h4>
    <div class="row mb-3">
      <label class="col-3 form-label">Пользователь телеграм</label>
      <div class="col-7">
        <input class="form-control" v-model="newTelegramUserName" />
      </div>
      <div class="col-1" v-if="user.telegramVerified">
        <span class="text-success"><i class="bi bi-star-fill"></i> Активно</span>
      </div>
    </div>
    <div class="row mb-3 text-danger" v-if="user.telegramUserName && !user.botChatId">
      <p>Это имя пользователя не найдено. Добавьте чат, напишите в чат любое сообщение и перезагрузите эту страницу для продолжения.</p>
    </div>
    <div class="row mb-3" v-if="user.botChatId && !user.telegramVerified && user.telegramUserName == newTelegramUserName">
      <div class="col">
        <button class="btn btn-primary" @click="NewTelegramCode">Отослать проверочный код в чат для подтверждения</button>
      </div>
    </div>
    <div class="row mb-3" v-if="user.botChatId && user.telegramCheckCode && !user.telegramVerified">
      <div class="row mb-3">
        <label class="col-3 form-label">Введите новый код из чата</label>
        <div class="col-7">
          <input class="form-control" v-model="telegramVerifyCode" />
        </div>
        <div class="col-1">
          <button class="btn btn-primary" @click="CheckTelegramCode">Проверить</button>
        </div>
      </div>
    </div>

    <div v-if="!route.params.id">
      <button class="btn btn-primary" @click="Save">Сохранить</button>
    </div>
    <div v-if="result" class="alert alert-primary">
      {{ result }}
    </div>

  </div>

</template>
