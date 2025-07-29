<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { ctx } from "./components/ctx.js";
import { useRoute } from "vue-router";

const route = useRoute();

const myurl = ref("");
const msgtext = ref("");

const updateMyUrl = () => {
  // telegram hash param delete from URL
  const url = new URL(window.location.href);
  url.hash = "";
  myurl.value = url.href;
};

onMounted(async () => {
  ctx.basket.sum = 0;
  ctx.SetAccessToken(
    localStorage.getItem("accessToken"),
    localStorage.getItem("accessTokenExpiresAt")
  );
  ctx.SetRefreshToken(
    localStorage.getItem("refreshToken"),
    localStorage.getItem("refreshTokenExpiresAt")
  );
  await ctx.CheckLogged();
  await ctx.loadBasket();
  updateMyUrl();
});

// update myurl when user navigate within SPA
watch(
  () => route.fullPath,
  () => {
    updateMyUrl();
  }
);

const SendMsg = async () => {
  if (await ctx.SendMsg(-1, msgtext.value)) {
    msgtext.value = "";
    alert("Сообщение отослано");
  }
};
</script>

<template>
  <header>
    <nav
      class="navbar fixed-top navbar-expand-sm navbar-toggleable-sm navbar-light bg-white border-bottom box-shadow.mb-3"
    >
      <div class="container-fluid">
        <RouterLink class="btn btn-outline-secondary btn-sm" to="/"
          ><img
            src="/src/assets/bkg.png"
            width="24px"
            height="24px"
          />&nbsp;Витрины</RouterLink
        >&nbsp;&nbsp;
        <span v-if="!ctx.userInfo.id && !ctx.isTg()">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/login"
            >Войти</RouterLink
          >&nbsp;
        </span>
        <span v-if="!ctx.userInfo.id && ctx.isTg()" class="text-danger">
          Напишите в чат бота любое сообщение для авто-регистрации&nbsp;
        </span>
        <ul class="navbar-nav flex-grow-1">
          &nbsp;
        </ul>
        <span v-if="ctx.basket.sum">
          <RouterLink class="btn btn-success btn-sm" to="/checkout"
            ><i class="bi bi-basket2-fill"></i>&nbsp;{{
              ctx.basket.sum
            }}</RouterLink
          >&nbsp;&nbsp;
        </span>
        <span v-if="ctx.userInfo.id">
          <RouterLink
            class="btn btn-outline-secondary btn-sm"
            to="/myprofile"
            >{{ ctx.userInfo.userName }}</RouterLink
          >&nbsp;
        </span>
      </div>
    </nav>
  </header>

  <br />
  <br />
  <br />
  <div class="container-fluid">
    <RouterView :key="route.fullPath" />
  </div>

  <br />

  <nav
    class="navbar navbar-expand-sm navbar-toggleable-sm navbar-light bg-white border-bottom box-shadow.mb-3"
    v-if="ctx.userInfo.admin || ctx.userInfo.shopManage"
  >
    <div class="container-fluid">
      <span v-if="ctx.userInfo.shopManage">
        <RouterLink class="btn btn-info btn-sm" to="/inc-orders"
          >Заказы для обработки</RouterLink
        >&nbsp;
      </span>
      <span v-if="ctx.userInfo.admin">
        <RouterLink class="btn btn-info btn-sm" to="/set-roles"
          >admin</RouterLink
        >&nbsp;
      </span>
    </div>
  </nav>
  <div class="container-fluid">
    <div class="row">
      <div class="col">
        <span>&copy;2024-2025,</span>&nbsp;&nbsp;<button
          class="btn btn-warning btn-sm"
          data-bs-toggle="modal"
          data-bs-target="#sendAdminMsgModal"
        >
          <i class="bi bi-chat-quote"></i>&nbsp;Вопросы, проблемы, предложения?
        </button>
        <div v-if="ctx.isTg() && ctx.userInfo.shopManage">
          Адрес страницы для браузера:
          <input class="form-control sm" v-model="myurl" />
        </div>
      </div>
    </div>
  </div>

  <div
    class="modal fade"
    id="sendAdminMsgModal"
    tabindex="-1"
    aria-labelledby="sendAdminMsgModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="sendAdminMsgModalLabel">
            Написать администратору сервиса
          </h5>
        </div>
        <div class="modal-body">
          <textarea class="form-control" rows="5" v-model="msgtext"></textarea>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-warning btn-sm"
            data-bs-dismiss="modal"
            @click="SendMsg()"
          >
            Послать
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
