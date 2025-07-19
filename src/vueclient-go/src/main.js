import { createApp } from 'vue';
import { createWebHistory, createRouter } from 'vue-router';
import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import EditShop from "./components/EditShop.vue";
import EditGood from "./components/EditGood.vue";
import Shop from "./components/Shop.vue";
import CheckOut from "./components/CheckOut.vue";
import Orders from "./components/Orders.vue";
import Order from "./components/Order.vue";
import Profile from "./components/Profile.vue";
import SetRoles from "./components/SetRoles.vue";
import IncOrders from "./components/IncOrders.vue";
import Good from "./components/Good.vue";

import App from './App.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: Login },
  { path: '/edit-shop/', component: EditShop },
  { path: '/edit-shop/:id', component: EditShop },
  { path: '/shop/:id', component: Shop },
  { path: '/edit-good/:shopid', component: EditGood },
  { path: '/edit-good/:shopid/:id', component: EditGood },
  { path: '/checkout', component: CheckOut },
  { path: '/orders', component: Orders },
  { path: '/order/:id', component: Order },
  { path: '/myprofile', component: Profile },
  { path: '/profile/:id', component: Profile },
  { path: '/set-roles', component: SetRoles },
  { path: '/inc-orders', component: IncOrders },
  { path: '/good/:id', component: Good },
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

createApp(App)
  .use(router)
  .mount('#app');
