import { reactive } from 'vue'
import moment from 'moment';
import 'moment/dist/locale/ru';
import linkifyHtml from 'linkify-html';

function newUserInfo() {
  return { userName: null, shopManage: false, admin: false, id: 0 }
}

export const authStore = reactive({

  basket: { sum: 0 },
  userInfo: newUserInfo(),
  accessToken: "",
  refreshToken: "",

  async loadBasket() {
    this.basket.sum = 0;
    let res = await fetch(`${authStore.rbUrl()}/api/baskets`, { headers: authStore.authHeaders() });
    if (res.ok) {
      res = await res.json();
      res.forEach((_) => {
        this.basket.sum += _.sum;
      });
    }
  },

  fmtDate(date) {
    moment.locale('ru');
    return moment(date).fromNow();
  },

  linkify(src) {
    const opt = {
      target: {
        url: "_blank",
        email: null,
      },
    };
    return linkifyHtml(src, opt);
  },

  tgInitData() {
    if (window.Telegram == undefined || window.Telegram.WebApp == undefined) return null;
    let r = window.Telegram.WebApp.initData;
    if (typeof (r) != "string" || r == "") return null;
    return r;
  },

  tgInitDataStruct() {
    if (window.Telegram == undefined || window.Telegram.WebApp == undefined) return null;
    return window.Telegram.WebApp.initDataUnsafe;
  },

  isTg() {
    var r = this.tgInitData();
    return r != null;
  },

  rbUrl() {
    return import.meta.env.VITE_API_URL;
  },

  tgAuthToken() {
    let u = this.tgInitDataStruct().user;
    let un = "";
    let uid = u.id;
    if (u.username) un = u.username;
    return `tg ${this.tgInitData()}~~${uid}~~${un}`;
  },

  authHeaders() {
    if (this.isTg())
      return { "Authorization": this.tgAuthToken() }
    else
      return { "Authorization": "Bearer " + this.accessToken };
  },

  authHeadersAppJson() {
    if (this.isTg())
      return {
        "Authorization": this.tgAuthToken(),
        "Content-Type": "application/json",
      };
    else
      return {
        "Authorization": "Bearer " + this.accessToken,
        "Content-Type": "application/json",
      };
  },

  SetAccessToken(newToken) {
    this.accessToken = newToken;
    localStorage.setItem("accessToken", newToken);
  },

  SetRefreshToken(newToken) {
    this.refreshToken = newToken;
    localStorage.setItem("refreshToken", newToken);
  },

  async RefreshToken() {

    if (this.isTg()) return;

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
          this.SetAccessToken(res.accessToken);
          this.SetRefreshToken(res.refreshToken);
          console.log('tokens was refreshed');
          this.CheckLogged();
        } else {
          this.SetAccessToken("");
          this.SetRefreshToken("");
          console.log("unauthorized tokens was cleared");
        }
      }
      catch (err) {
        console.log(err);
      };
    }
  },

  async CheckLogged() {

    if (this.isTg()) {
      let res = await fetch(`${this.rbUrl()}/api/profiles/public?email=0`, { headers: this.authHeaders() });
      if (res.ok) {
        res = await res.json();
        this.userInfo = res;
      }
      return;
    }

    if (this.accessToken != "") {
      try {
        let res = await fetch(this.rbUrl() + "/identity/manage/info", {
          method: "GET",
          signal: AbortSignal.timeout(5000),
          headers: this.authHeaders()
        });
        if (res.ok) {
          let obj = await res.json();
          let email = obj.email;

          res = await fetch(this.rbUrl() + "/api/profiles/public?email=" + encodeURIComponent(email))
          if (res.ok) {
            res = await res.json();
            this.userInfo = res;
          }
        }
      } catch (err) {
        console.log(err);
        this.RefreshToken();
      }
    } else await this.RefreshToken();
  },

  async Login(emailString, passwordString) {

    if (this.isTg()) return;

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
      } else {
        let r = await res.json();
        return this.errFormat(r);
      }
    } catch (err) { return err; }
    return "";
  },

  async Logout() {

    if (this.isTg()) return;

    let response = await fetch(this.rbUrl() + "/identity/logout", {
      method: "POST",
      signal: AbortSignal.timeout(5000),
      headers: { "Content-Type": "application/json" }
    });
    if (response.ok) {
      this.userInfo = newUserInfo();
      this.SetAccessToken("");
      this.SetRefreshToken("");
      await this.CheckLogged();
    }

  },

  errFormat(obj) {
    let r = `${obj.title}: `;
    if (obj.errors)
      Object.values(obj.errors).forEach(_ => {
        r += _ + " "
      });
    return r;
  },

  async Register(emailString, passwordString) {

    if (this.isTg()) return;

    let request = { email: emailString, password: passwordString };
    try {
      let response = await fetch(this.rbUrl() + "/identity/register", {
        method: "POST",
        body: JSON.stringify(request),
        headers: {
          "Content-Type": "application/json"
        },
      });
      let r = "";
      if (response.ok)
        r = await this.Login(emailString, passwordString);
      else
        r = this.errFormat(await response.json());
      return r;
    } catch (err) {
      return err;
    }
  }

})
