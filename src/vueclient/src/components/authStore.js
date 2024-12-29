import { reactive, ref } from 'vue'

const rbUrl = "http://localhost:5078";

export const authStore = reactive({
  loggedEmail: "",
  accessToken: "",
  refreshToken: "",

  SetAccessToken(newToken) {
    this.accessToken = newToken;
  },

  SetRefreshToken(newToken) {
    this.refreshToken = newToken;
  },

  async CheckLogged() {
    if (this.accessToken != "") {
      let response = await fetch(rbUrl + "/identity/manage/info", {
        method: "GET",
        headers: { "Authorization": "Bearer " + this.accessToken }
      });
      if (response.ok) {
        let res = await response.json();
        this.loggedEmail = res.email;
      }
    }
  },

  async Login (emailString, passwordString) {
    let request = { email: emailString, password: passwordString };
    let response = await fetch(rbUrl + "/identity/login", {
      method: "POST",
      body: JSON.stringify(request),
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      let res = await response.json();
      this.SetAccessToken(res.accessToken);
      this.SetRefreshToken(res.refreshToken);
      this.loggedEmail = emailString;
      localStorage.setItem("accessToken", this.accessToken);
      localStorage.setItem("refreshToken", this.refreshToken);
      this.CheckLogged();
    } else alert(await response.text()); //todo handle errors, don't show

  },

  async Register (emailString, passwordString) {
    let request = { email: emailString, password: passwordString };
    let response = await fetch(rbUrl + "/identity/register", {
      method: "POST",
      body: JSON.stringify(request),
      headers: {
        "Content-Type": "application/json"
      },
    });
    if (response.ok) {
      let res = await response.json();
      this.accessToken = res.accessToken;
      this.refreshToken = res.RefreshToken;
      localStorage.setItem("accessToken", this.accessToken);
      localStorage.setItem("refreshToken", this.refreshToken);

    } else
      alert(await response.text()); //todo handle errors, don't show
  }

})
