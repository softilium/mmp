import { createApp } from 'vue'
import { createWebHistory, createRouter } from 'vue-router';
import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import Logout from "./components/Logout.vue";
import Register from "./components/Register.vue";
import EditShop from "./components/EditShop.vue";
import EditGood from "./components/EditGood.vue";
import Shop from "./components/Shop.vue";
import CheckOut from "./components/CheckOut.vue";

import App from './App.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: Login },
  { path: '/logout', component: Logout },
  { path: '/register', component: Register },
  { path: '/edit-shop/', component: EditShop },
  { path: '/edit-shop/:id', component: EditShop },
  { path: '/shop/:id', component: Shop },
  { path: '/edit-good/:shopid', component: EditGood },
  { path: '/edit-good/:shopid/:id', component: EditGood },
  { path: '/checkout/:shopid', component: CheckOut },
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

createApp(App)
  .use(router)
  .mount('#app');
