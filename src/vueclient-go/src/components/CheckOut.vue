<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ctx } from "./ctx.js";
import ProfileLink from "./ProfileLink.vue";
import { localBasket } from "../services/localBasket";

const router = useRouter();

const basket = ref([]);

onMounted(async () => {
  await Load();
});

const Load = async () => {
  basket.value = [];

  if (ctx.userInfo.userName) {
    let res = await fetch(ctx.rbUrl() + "/api/basket", {
      method: "GET",
      headers: await ctx.authHeadersAppJson(),
    });
    if (await ctx.CheckUnauth(res)) return;
    if (res.ok) {
      let rj = await res.json();
      let rj2 = rj.Data;
      const senderMap = new Map();
      rj2.forEach((_) => {
        senderMap.set(_.Good.CreatedBy.Ref, 1);
      });
      basket.value = [];
      for (let [k] of senderMap.entries()) {
        let senderRec = {
          senderID: k,
          lines: [],
          qtyTotal: 0,
          sumTotal: 0,
          customerComment: "",
          senderInfo: {},
        };
        basket.value.push(senderRec);
        rj2.forEach((r) => {
          r.good = {
            id: r.Good.Ref,
            caption: r.Good.Caption,
            createdByInfo: r.Good.CreatedBy,
          };
          if (r.Good.CreatedBy.Ref == k) {
            senderRec.lines.push(r);
            senderRec.qtyTotal += r.Qty;
            senderRec.sumTotal += r.Sum;
            senderRec.senderInfo = r.Good.CreatedBy;
          }
        });
      }
    }
  } else {
    // Load basket for anonymous users from localBasket
    const items = localBasket.getItems();
    const senderMap = new Map();
    for (const item of items) {
      const res = await fetch(`${ctx.rbUrl()}/api/goods?ref=${item.goodId}`);
      if (res.ok) {
        const good = await res.json();
        senderMap.set(item.senderId, 1); // Use senderId from local basket
        let senderRec = basket.value.find((b) => b.senderID === item.senderId);
        if (!senderRec) {
          senderRec = {
            senderID: item.senderId,
            lines: [],
            qtyTotal: 0,
            sumTotal: 0,
            customerComment: "",
          };
          basket.value.push(senderRec);
        }
        senderRec.lines.push({
          good: {
            id: good.Ref,
            caption: good.Caption,
            createdByInfo: good.CreatedBy,
          },
          qty: item.quantity,
          sum: item.quantity * good.Price,
        });
        senderRec.qtyTotal += item.quantity;
        senderRec.sumTotal += item.quantity * good.Price;
        senderRec.senderInfo = good.CreatedBy;
      }
    }
  }

  if (basket.value.length == 0) {
    router.push("/");
  }
};

const Checkout = async (senderRec) => {
  let res = await fetch(
    ctx.rbUrl() + "/api/orders/checkout?sender=" + senderRec.senderInfo.Ref,
    {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
      body: senderRec.customerComment,
    }
  );
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    await ctx.loadBasket();
    let rj = await res.json();
    router.push(`/order/${rj.Ref}`);
  }
};

const Inc = async (good) => {
  if (!ctx.userInfo.id) {
    localBasket.addItem({
      goodId: good.id,
      quantity: 1,
      price: good.price,
      title: good.caption,
      shopTitle:
        good.createdByInfo && good.createdByInfo.userName
          ? good.createdByInfo.userName
          : "",
      senderId:
        good.createdByInfo && good.createdByInfo.id
          ? good.createdByInfo.id
          : null,
      shopId: good.shopId || null,
    });
    await Load();
    await ctx.loadBasket();
    return;
  }
  let res = await fetch(
    ctx.rbUrl() + "/api/basket/increase?goodref=" + good.id,
    {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
    }
  );
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    await Load();
    ctx.loadBasket();
  }
};

const Dec = async (good) => {
  if (!ctx.userInfo.id) {
    localBasket.decItem(good.id);
    await Load();
    await ctx.loadBasket();
    return;
  }
  let res = await fetch(
    ctx.rbUrl() + "/api/basket/decrease?goodref=" + good.id,
    {
      method: "POST",
      headers: await ctx.authHeadersAppJson(),
    }
  );
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    await Load();
    ctx.loadBasket();
  }
};
</script>

<template>
  <div v-for="sender in basket" v-bind:key="sender.senderID">
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
          <td class="col-8">
            <RouterLink :to="`/good/${line.good.id}`">{{
              line.good.caption
            }}</RouterLink>
          </td>
          <td>
            &nbsp;<button
              class="btn btn-primary btn-sm"
              @click="Inc(line.good)"
            >
              +
            </button>
            &nbsp;<button
              class="btn btn-primary btn-sm"
              @click="Dec(line.good)"
            >
              -
            </button>
            {{ line.Qty }}
          </td>
          <td class="text-end">{{ line.Sum }}</td>
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
          <textarea
            class="form-control"
            v-model="sender.customerComment"
            rows="5"
          />
        </div>
      </div>
    </div>

    <div class="row mb-3" v-if="ctx.userInfo.id">
      <button class="btn btn-primary btn-sm" @click="Checkout(sender)">
        Оформить заказ для {{ sender.senderInfo.Username }}
      </button>
    </div>
    <div class="row mb-3" v-if="!ctx.userInfo.id">
      <RouterLink class="btn btn-primary btn-sm" to="/login">
        Для оформления заказа необходимо войти в систему или зарегистрироваться
      </RouterLink>
    </div>
  </div>
</template>
