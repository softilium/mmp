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
    localStorage.setItem("refreshToken", newToken);
  },

  async CheckLogged() {
    if (this.accessToken != "") {
      try {
        let res = await fetch(this.rbUrl() + "/identity/manage/info", {
          method: "GET",
          signal: AbortSignal.timeout(5000),
          headers: { "Authorization": "Bearer " + this.accessToken }
        });
        if (res.ok) {
          let obj = await res.json();
          this.loggedEmail = obj.email;
        }
      } catch (err) {
        console.log(err);
        //try to use refresh token
        if (this.refreshToken != "") {
          try {
            let res = await fetch(this.rbUrl() + "/identity/refresh", {
              method: "POST",
              signal: AbortSignal.timeout(5000),
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({ "refreshToken": this.refreshToken })
            });
            if (res.ok) {
              res = await res.json();
              this.loggedEmail = res.email;
              this.SetAccessToken(res.accessToken);
              this.SetRefreshToken(res.refreshToken);
              console.log('tokens was refreshed');
            }
          }
          catch (err) {
            console.log(err);
            this.SetAccessToken("");
            this.SetRefreshToken("");
          };
        }
      }
    }
  },

  async Login(emailString, passwordString) {
    try {
      let res = await fetch(this.rbUrl() + "/identity/login", {
        method: "POST",
        signal: AbortSignal.timeout(5000),
        body: JSON.stringify({ email: emailString, password: passwordString }),
        headers: { "Content-Type": "application/json" }
      });
      if (res.ok) {
        res = await res.json();
        this.SetAccessToken(res.accessToken);
        this.SetRefreshToken(res.refreshToken);
        this.CheckLogged();
      }
    } catch (err) { console.log(err); }
  }, //todo handle errors, don't show

  async Logout() {
    let response = await fetch(this.rbUrl() + "/identity/logout", {
      method: "POST",
      signal: AbortSignal.timeout(5000),
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      this.loggedEmail = "";
      this.SetAccessToken("");
      this.SetRefreshToken("");
      await this.CheckLogged();
    } 

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
      await this.Login(emailString, passwordString);
    }
  }

})
