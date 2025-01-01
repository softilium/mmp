import { reactive } from 'vue'

export const authStore = reactive({
  rbUrl: () => "http://localhost:5078",
  loggedEmail: "",
  accessToken: "",
  refreshToken: "",

  SetAccessToken(newToken) {
    this.accessToken = newToken;
    localStorage.setItem("accessToken", newToken);
  },

  SetRefreshToken(newToken) {
    this.refreshToken = newToken;
    localStorage.setItem("refreshToken", this.refreshToken);
  },

  async CheckLogged() {
    if (this.accessToken != "") {
      await fetch(this.rbUrl() + "/identity/manage/info", {
        method: "GET",
        headers: { "Authorization": "Bearer " + this.accessToken }
      }).then(async res => {
        let obj = await res.json();
        this.loggedEmail = obj.email;
      }).catch(async () => {
        this.SetAccessToken("");
        //try to use refresh token
        if (this.refreshToken != "") {
          await fetch(this.rbUrl + "/identity/refresh",
            {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({ "refreshToken": this.refreshToken })
            })
            .then(async res => {
              let obj = await res.json();
              this.loggedEmail = res.email;
              this.SetAccessToken(obj.accessToken);
              this.SetRefreshToken(obj.refreshToken);
              console.log('tokens was refreshed');
            })
            .catch(err => console.log(err));
        }
      });
    }
  },

  async Login(emailString, passwordString) {
    let request = { email: emailString, password: passwordString };
    let response = await fetch(this.rbUrl() + "/identity/login", {
      method: "POST",
      body: JSON.stringify(request),
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      let res = await response.json();
      this.SetAccessToken(res.accessToken);
      this.SetRefreshToken(res.refreshToken);
      this.loggedEmail = emailString;
      this.CheckLogged();
    } else alert(await response.text()); //todo handle errors, don't show

  },

  async Logout() {
    let response = await fetch(this.rbUrl() + "/identity/logout", {
      method: "POST",
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      this.SetAccessToken("");
      this.SetRefreshToken("");
      this.loggedEmail = "";
      this.CheckLogged();
    } else alert(await response.text()); //todo handle errors, don't show

  },

  async Register(emailString, passwordString) {
    let request = { email: emailString, password: passwordString };
    let response = await fetch(this.rbUrl() + "/identity/register", {
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
