<script setup>

  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router'
  import { ctx } from './ctx.js';
  import ProfileLink from './ProfileLink.vue';
  import { localBasket } from '../services/localBasket';

  const router = useRouter();

  const basket = ref([]);

  onMounted(async () => {
    await Load();
  });

  const Load = async () => {
    basket.value = [];

    if (ctx.userInfo.userName) {
      let res = await fetch(ctx.rbUrl() + "/api/baskets/", {
        mathod: "GET",
        headers: ctx.authHeadersAppJson()
      });
      if (res.ok) {
        res = await res.json();
        const senderMap = new Map();
        res.forEach(_ => {
          senderMap.set(_.good.createdByID, 1);
        });
        basket.value = [];
        for (let [k] of senderMap.entries()) {
          let senderRec = { senderID: k, lines: [], qtyTotal: 0, sumTotal: 0, customerComment: "" };
          basket.value.push(senderRec);
          res.forEach(r => {
            if (r.good.createdByID == k) {
              senderRec.lines.push(r);
              senderRec.qtyTotal += r.qty;
              senderRec.sumTotal += r.sum;
              senderRec.senderInfo = r.good.createdByInfo;
            }
          });
        }
      }
    } else {
      // Load basket for anonymous users from localBasket
      const items = localBasket.getItems();
      const senderMap = new Map();
      for (const item of items) {
        const res = await fetch(`${ctx.rbUrl()}/api/goods/${item.goodId}`);
        if (res.ok) {
          const good = await res.json();
          senderMap.set(item.senderId, 1); // Use senderId from local basket
          let senderRec = basket.value.find(b => b.senderID === item.senderId);
          if (!senderRec) {
            senderRec = { senderID: item.senderId, lines: [], qtyTotal: 0, sumTotal: 0, customerComment: "" };
            basket.value.push(senderRec);
          }
          senderRec.lines.push({
            good: {
              id: good.id,
              caption: good.caption,
              createdByInfo: good.createdByInfo
            },
            qty: item.quantity,
            sum: item.quantity * good.price
          });
          senderRec.qtyTotal += item.quantity;
          senderRec.sumTotal += item.quantity * good.price;
          senderRec.senderInfo = good.createdByInfo;
        }
      }
    }

    if (basket.value.length == 0) {
      router.push("/");
    }
  };

  const Checkout = async (senderRec) => {
    let res = await fetch(ctx.rbUrl() + "/api/orders/outbox/" + senderRec.senderID, {
      method: "POST",
      headers: ctx.authHeadersAppJson(),
      body: senderRec.customerComment
    });
    if (res.ok) {
      await ctx.loadBasket();
      res = await res.json();
      router.push(`/order/${res.id}`);
    }
  };

  const Inc = async (good) => {
    if (!ctx.userInfo.id) {
      localBasket.addItem({
        goodId: good.id,
        quantity: 1,
        price: good.price,
        title: good.caption,
        shopTitle: good.createdByInfo && good.createdByInfo.userName ? good.createdByInfo.userName : '',
        senderId: good.createdByInfo && good.createdByInfo.id ? good.createdByInfo.id : null,
        shopId: good.shopId || null
      });
      await Load();
      await ctx.loadBasket();
      return;
    }
    let res = await fetch(ctx.rbUrl() + "/api/baskets/increase/" + good.id, {
      method: "POST",
      headers: ctx.authHeadersAppJson()
    });
    if (res.ok) {
      await Load();
      ctx.loadBasket();
    }
  }

  const Dec = async (good) => {
    if (!ctx.userInfo.id) {
      localBasket.decItem(good.id);
      await Load();
      await ctx.loadBasket();
      return;
    }
    let res = await fetch(ctx.rbUrl() + "/api/baskets/decrease/" + good.id, {
      method: "POST",
      headers: ctx.authHeadersAppJson()
    });
    if (res.ok) {
      await Load();
      ctx.loadBasket();
    }
  }

</script>

<template>

  <div v-for="sender in basket">
    <h1>
      Отправитель <ProfileLink :userInfo="sender.senderInfo"></ProfileLink>
    </h1>

    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th class="col-8">Товар или услуга</th>
          <th class="text-end">Количество</th>
          <th class="text-end">Стоимость</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="line in sender.lines" v-bind:key="line.id">
          <td class="col-8"><RouterLink :to="`/good/${line.good.id}`">{{ line.good.caption }}</RouterLink></td>
          <td class="text-end">
            {{ line.qty }}
            &nbsp;<button class="btn btn-primary btn-sm" @click="Inc(line.good)">+</button>
            &nbsp;<button class="btn btn-primary btn-sm" @click="Dec(line.good)">-</button>
          </td>
          <td class="text-end">{{ line.sum }}</td>
        </tr>
      </tbody>
      <tfoot class="table-success">
        <tr>
          <td></td>
          <td class="text-end">{{ sender.qtyTotal }}</td>
          <td class="text-end">{{ sender.sumTotal }}</td>
        </tr>
      </tfoot>
    </table>

    <div class="row mb-3" v-if="ctx.userInfo.id">
      <div class="form-group form-group-sm row">
        <label class="col-4 form-label">Примечание к заказу</label>
        <div class="col-8">
          <textarea class="form-control" v-model="sender.customerComment" rows="5" />
        </div>
      </div>
    </div>

    <div class="row mb-3" v-if="ctx.userInfo.id">
      <button class="btn btn-primary btn-sm" @click="Checkout(sender)">Оформить заказ для {{sender.senderInfo.userName}}</button>
    </div>
    <div class="row mb-3" v-if="!ctx.userInfo.id">
      <RouterLink class="btn btn-primary btn-sm" to="/login" >
        Для оформления заказа необходимо войти в систему или зарегистрироваться
      </RouterLink>
    </div>

  </div>

</template>
