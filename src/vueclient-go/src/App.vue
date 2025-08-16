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
  if (await ctx.SendMsg("", msgtext.value)) {
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
        <RouterLink class="btn" to="/">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="32"
            height="32"
            fill="currentColor"
            class="bi bi-house-heart-fill"
            viewBox="0 0 16 16"
          >
            <path
              d="M7.293 1.5a1 1 0 0 1 1.414 0L11 3.793V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5v3.293l2.354 2.353a.5.5 0 0 1-.708.707L8 2.207 1.354 8.853a.5.5 0 1 1-.708-.707z"
            />
            <path
              d="m14 9.293-6-6-6 6V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5zm-6-.811c1.664-1.673 5.825 1.254 0 5.018-5.825-3.764-1.664-6.691 0-5.018"
            />
          </svg> </RouterLink
        >&nbsp;&nbsp;
        <RouterLink class="btn" to="/search">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="32"
            height="32"
            fill="currentColor"
            class="bi bi-search-heart"
            viewBox="0 0 16 16"
          >
            <path
              d="M6.5 4.482c1.664-1.673 5.825 1.254 0 5.018-5.825-3.764-1.664-6.69 0-5.018"
            />
            <path
              d="M13 6.5a6.47 6.47 0 0 1-1.258 3.844q.06.044.115.098l3.85 3.85a1 1 0 0 1-1.414 1.415l-3.85-3.85a1 1 0 0 1-.1-.115h.002A6.5 6.5 0 1 1 13 6.5M6.5 12a5.5 5.5 0 1 0 0-11 5.5 5.5 0 0 0 0 11"
            />
          </svg>
        </RouterLink>
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
        <span v-if="!ctx.userInfo.id && !ctx.isTg()">
          <RouterLink class="btn btn-outline-secondary btn-sm" to="/login"
            >Войти</RouterLink
          >&nbsp;
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
