<script setup lang="ts">
  //import HelloWorld from './components/HelloWorld.vue'
  //import TheWelcome from './components/TheWelcome.vue'
  //import NavBar from './components/NavBar.vue'
  import { onMounted, ref } from 'vue'

  const rbUrl = "http://localhost:5078";

  let accessToken = null;
  let refreshToken = null;
  const loggedEmail = ref("");

  onMounted(() => {

    accessToken = localStorage.getItem("accessToken");
    refreshToken = localStorage.getItem("refreshToken");

    checkLogged();

  });

  let emailString = "";
  let passwordString = "";

  const checkLogged = async () => {
    if (accessToken != null) {
      let response = await fetch(rbUrl + "/identity/manage/info", {
        method: "GET",
        headers: { "Authorization": "Bearer " + accessToken }
      });
      if (response.ok) {
        let res = await response.json();
        loggedEmail.value = res.email;
      }
    }
  };


  const doLogin = async () => {

    let request = { email: emailString, password: passwordString };

    let response = await fetch(rbUrl + "/identity/login", {
      method: "POST",
      body: JSON.stringify(request),
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      let res = await response.json();
      accessToken = res.accessToken;
      refreshToken = res.refreshToken;
      localStorage.setItem("accessToken", accessToken);
      localStorage.setItem("refreshToken", refreshToken);
      checkLogged();
    } else alert(await response.text()); //todo handle errors, don't show

  };

  const doReg = async () => {

    var request = { email: emailString, password: passwordString };

    let response = await fetch(rbUrl + "/identity/register", {
      method: "POST",
      body: JSON.stringify(request),
      headers: {
        "Content-Type": "application/json"
      },
    });
    if (response.ok) {
      let res = await response.json();
      accessToken = res.accessToken;
      refreshToken = res.RefreshToken;

      localStorage.setItem("accessToken", accessToken);
      localStorage.setItem("refreshToken", refreshToken);

    } else
      alert(await response.text()); //todo handle errors, don't show

  };

</script>

<template>

  Logged Email: {{ loggedEmail }}

  <a href="/">Home</a>
  <a v-if="!loggedEmail" href="/login">Login</a>
  <a v-if="loggedEmail" href="/logout">Logout</a>

  

  <div>
    Email: <input v-model="emailString" />
    Password: <input v-model="passwordString" />
    <button @click="doLogin">Login</button>
    <button @click="doReg">Register</button>
  </div>

</template>
