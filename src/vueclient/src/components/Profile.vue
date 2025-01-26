<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">

  import { onMounted, ref } from 'vue';
  import { useRoute } from 'vue-router'
  import { authStore } from './authStore.js';
  import { glob } from './globals.js';

  const route = useRoute();

  const me = ref(false);
  const user = ref({ userName: "", email: "", telegramUserName: "", telegramVerified: false, telegramCheckCode: "", botChatId: 0 });
  const newTelegramUserName = ref("");
  const telegramVerifyCode = ref("");
  const result = ref("");
  const userDescription = ref("");

  const NewTelegramCode = async () => {

    if (me.value != true) return;

    let res = await fetch(authStore.rbUrl() + "/api/profiles/newtelegramcode", {
      method: "POST",
      headers: authStore.authHeadersAppJson()
    });
    if (res.ok) Load();

  }
  const CheckTelegramCode = async () => {

    if (me.value != true) return;

    let res = await fetch(authStore.rbUrl() + "/api/profiles/checktelegramcode/" + telegramVerifyCode.value, {
      method: "POST",
      headers: authStore.authHeadersAppJson()
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
          headers: authStore.authHeadersAppJson()
        });
      if (res.ok) {
        user.value = await res.json();
        userDescription.value = glob.linkify(user.value.description);
        me.value = !route.params.id;
      }
    }
    else {
      let res = await fetch(authStore.rbUrl() + "/api/profiles/my",
        {
          headers: authStore.authHeadersAppJson()
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
        headers: authStore.authHeadersAppJson()
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
    Профиль пользователя {{ user.userName }}&nbsp;<br />
    <button v-if="me && !authStore.isTg()" class="btn btn-outline-secondary btn-sm" @click="authStore.Logout(); $router.push('/');">Выйти</button>&nbsp;
    <RouterLink v-if="me" class="btn btn-outline-secondary btn-sm" to="/orders">История заказов</RouterLink>
  </h1>
  <br />

  <div class="row mb-3">
    <label class="col-4 form-label">Имя пользователя</label>
    <div class="col-7">
      <input class="form-control" v-model="user.userName" required :disabled="!me" />
    </div>
  </div>
  <div class="row mb-3" v-if="me">
    <label class="col-4 form-label">Email/логин</label>
    <div class="col-7">
      <input type="email" class="form-control" v-model="user.email" />
    </div>
  </div>

  <div class="row mb-3">
    <label class="col-4 form-label">Обо мне</label>
    <div class="col-7">
      <span v-if="!me" v-html="userDescription"></span>
      <textarea v-if="me" class="form-control" v-model="user.description" rows="7" maxlength="300" />
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
      <label class="col-4 form-label">Пользователь телеграм</label>
      <div class="col-7">
        <input class="form-control" v-model="newTelegramUserName" :readonly="authStore.isTg()" />
      </div>
      <div class="col-1">
        <span class="text-success"><i class="bi bi-star-fill"></i></span>
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
      <button class="btn btn-secondary btn-sm" @click="Save">Сохранить</button>
    </div>
    <div v-if="result" class="alert alert-primary">
      {{ result }}
    </div>

  </div>

</template>
