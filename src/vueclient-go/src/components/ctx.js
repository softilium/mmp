import { reactive } from "vue";
import moment from "moment";
import "moment/dist/locale/ru";
import linkifyHtml from "linkify-html";
import { localBasket } from "../services/localBasket";

function newUserInfo() {
  return { userName: null, shopManage: false, admin: false, id: "" };
}

export const ctx = reactive({
  basket: { sum: 0 },
  userInfo: newUserInfo(),
  accessToken: "",
  refreshToken: "",
  accessTokenExpiresAt: 0,
  refreshTokenExpiresAt: 0,

  async loadBasket() {
    this.basket.sum = 0;
    if (!this.userInfo.id) {
      // Anonymous user: calculate sum from localBasket
      const items = localBasket.getItems();
      items.forEach((item) => {
        this.basket.sum += item.price * item.quantity;
      });
    } else {
      let res = await fetch(`${ctx.rbUrl()}/api/basket`, {
        headers: await ctx.authHeaders(),
      });
      if (await this.CheckUnauth(res)) return;
      if (res.ok) {
        res = await res.json();
        res.Data.forEach((_) => {
          this.basket.sum += _.Sum;
        });
      }
    }
  },

  fmtDate(date) {
    moment.locale("ru");
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
    if (window.Telegram == undefined || window.Telegram.WebApp == undefined)
      return null;
    let r = window.Telegram.WebApp.initData;
    if (typeof r != "string" || r == "") return null;
    return r;
  },

  tgInitDataStruct() {
    if (window.Telegram == undefined || window.Telegram.WebApp == undefined)
      return null;
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

  async authHeaders() {
    if (this.isTg()) return { Authorization: this.tgAuthToken() };
    else {
      if (this.accessTokenExpiresAt < Date.now()) {
        await this.RefreshToken();
      }
    }

    return { Authorization: "Bearer " + this.accessToken };
  },

  async authHeadersAppJson() {
    if (this.isTg())
      return {
        Authorization: this.tgAuthToken(),
        "Content-Type": "application/json",
      };
    else {
      if (this.accessTokenExpiresAt < Date.now()) {
        await this.RefreshToken();
      }
      return {
        Authorization: "Bearer " + this.accessToken,
        "Content-Type": "application/json",
      };
    }
  },

  SetAccessToken(newToken, expires) {
    this.accessToken = newToken;
    this.accessTokenExpiresAt = expires;
    localStorage.setItem("accessToken", newToken);
    localStorage.setItem("accessTokenExpiresAt", expires);
  },

  SetRefreshToken(newToken, expires) {
    this.refreshToken = newToken;
    this.refreshTokenExpiresAt = expires;
    localStorage.setItem("refreshToken", newToken);
    localStorage.setItem("refreshTokenExpiresAt", expires);
  },

  async RefreshToken() {
    if (this.isTg()) return;
    if (this.refreshToken != "") {
      let res = await fetch(this.rbUrl() + "/identity/refresh", {
        method: "POST",
        signal: AbortSignal.timeout(5000),
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refreshToken: this.refreshToken }),
      });
      if (await this.CheckUnauth(res)) return;
      if (res.ok) {
        res = await res.json();
        this.SetAccessToken(res.accessToken, res.accessTokenExpiresAt);
        this.SetRefreshToken(res.refreshToken, res.refreshTokenExpiresAt);
      }
    }
  },

  async CheckLogged() {
    if (this.isTg()) {
      let res = await fetch(`${this.rbUrl()}/api/profiles/public?email=0`, {
        headers: await this.authHeaders(),
      });
      if (res.ok) {
        res = await res.json();
        this.userInfo = res;
      }
      return;
    }
    if (this.accessToken != "") {
      let res = await fetch(this.rbUrl() + "/identity/myprofile", {
        method: "GET",
        signal: AbortSignal.timeout(5000),
        headers: await this.authHeaders(),
      });
      if (await this.CheckUnauth(res)) return;
      if (res.ok) {
        res = await res.json();
        this.userInfo = res;
      }
    }
  },

  async CheckUnauth(res) {
    if (res.status == 401) {
      let responseText = await res.text();
      if (responseText.includes("token expired") && this.accessToken != "") {
        this.accessToken = "";
        await this.RefreshToken();
        if (this.accessToken != "") {
          this.CheckLogged();
          return;
        }
      }
      this.userInfo = newUserInfo();
      this.SetAccessToken("");
      this.SetRefreshToken("");
      return true;
    }
    return false;
  },

  async Login(emailString, passwordString) {
    if (this.isTg()) return;

    try {
      let res = await fetch(this.rbUrl() + "/identity/login", {
        method: "POST",
        signal: AbortSignal.timeout(5000),
        body: JSON.stringify({ email: emailString, password: passwordString }),
        headers: { "Content-Type": "application/json" },
      });
      if (res.ok) {
        res = await res.json();
        this.SetAccessToken(res.accessToken, res.accessTokenExpiresAt);
        this.SetRefreshToken(res.refreshToken, res.refreshTokenExpiresAt);
        await this.mergeAnonymousBasket();
        await this.CheckLogged();
        await this.loadBasket();
      } else {
        return await res.text();
      }
    } catch (err) {
      return err;
    }
    return "";
  },

  async Logout() {
    if (this.isTg()) return;

    let response = await fetch(this.rbUrl() + "/identity/logout", {
      method: "POST",
      signal: AbortSignal.timeout(5000),
      headers: await this.authHeaders(),
    });
    if (response.ok) {
      this.userInfo = newUserInfo();
      this.SetAccessToken("");
      this.SetRefreshToken("");
      await this.loadBasket();
      await this.CheckLogged();
    }
  },

  errFormat(obj) {
    let r = `${obj.title}: `;
    if (obj.errors)
      Object.values(obj.errors).forEach((_) => {
        r += _ + " ";
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
          "Content-Type": "application/json",
        },
      });
      let r = "";
      if (response.ok) {
        r = await this.Login(emailString, passwordString);
        await this.mergeAnonymousBasket();
      } else r = this.errFormat(await response.json());
      return r;
    } catch (err) {
      return err;
    }
  },

  async SendMsg(userid, msgtext) {
    if (!msgtext) return false;

    let res = await fetch(`${ctx.rbUrl()}/api/profiles/sendmsg/${userid}`, {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
      body: msgtext,
    });

    return res.ok;
  },

  async mergeAnonymousBasket() {
    const items = localBasket.getItems();
    if (!items.length) return;
    // Prepare payload for API
    const payload = items.map((item) => ({
      goodId: item.goodId,
      quantity: item.quantity,
    }));
    await fetch(`${this.rbUrl()}/api/basket/merge`, {
      method: "POST",
      headers: await this.authHeadersAppJson(),
      body: JSON.stringify(payload),
    });
    localBasket.clear();
  },
});
