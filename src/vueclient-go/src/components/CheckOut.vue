<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ctx } from "./ctx.js";
import ProfileLink from "./ProfileLink.vue";
import { localBasket } from "../services/localBasket";
import { priceWithDiscount } from "../services/localBasket";

const router = useRouter();

const basket = ref([]);

interface Shop {
  Ref: string;
  DiscountPercent: number;
}

interface Good {
  Ref: string;
  OwnerShop: Shop;
  Caption: string;
  Price: number;
  IsDeleted: boolean;
}

interface BasketLine {
  Ref: string;
  Good: Good;
  Qty: number;
  Sum: number;
}

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
      senderMap.forEach((_, k) => {
        let senderRec = {
          senderID: k,
          lines: Array<BasketLine>(),
          qtyTotal: 0,
          sumTotal: 0,
          customerComment: "",
          senderInfo: {},
        };
        basket.value.push(senderRec);
        rj2.forEach((r) => {
          if (r.Good.CreatedBy.Ref == k) {
            senderRec.lines.push(r);
            r.Sum = r.Qty * priceWithDiscount(r.Good);
            senderRec.qtyTotal += r.Qty;
            senderRec.sumTotal += r.Sum;
            senderRec.senderInfo = r.Good.CreatedBy;
          }
        });
      });
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
            lines: Array<BasketLine>(),
            qtyTotal: 0,
            sumTotal: 0,
            customerComment: "",
          };
          basket.value.push(senderRec);
        }
        senderRec.lines.push({
          Good: good,
          Qty: item.quantity,
          Sum: item.quantity * priceWithDiscount(good),
        });
        senderRec.qtyTotal += item.quantity;
        senderRec.sumTotal += item.quantity * priceWithDiscount(good);
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
      goodId: good.Ref,
      quantity: 1,
      senderId: good.CreatedBy.Ref,
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
    localBasket.decItem(good.Ref);
    await Load();
    await ctx.loadBasket();
    return;
  }
  let res = await fetch(
    ctx.rbUrl() + "/api/basket/decrease?goodref=" + good.Ref,
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
            <span v-if="!line.Good.IsDeleted">
              <RouterLink :to="`/good/${line.Good.Ref}`">{{
                line.Good.Caption
              }}</RouterLink>
            </span>
            <span v-if="line.Good.IsDeleted">
              {{ line.Good.Caption }} [Удален]
            </span>
          </td>
          <td>
            &nbsp;<button
              class="btn btn-primary btn-sm"
              @click="Inc(line.Good)"
            >
              +
            </button>
            &nbsp;<button
              class="btn btn-primary btn-sm"
              @click="Dec(line.Good)"
            >
              -
            </button>
            {{ line.Qty }}
          </td>
          <td class="text-end">{{ line.Sum }}</td>
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
