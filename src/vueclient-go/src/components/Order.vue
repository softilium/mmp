<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ctx } from "./ctx.js";
import ProfileLink from "./ProfileLink.vue";

const route = useRoute();
const router = useRouter();

const order = ref({
  CustomerComment: "",
  SenderComment: "",
  ExpectedDeliveryDate: "",
  CreatedAt: "",
  Sender: { Ref: "", Username: "", Description: "" },
  CreatedBy: { Ref: "", Username: "", Description: "" },
  Status: 0,
  Ref: "",
  Nbr: "",
  lines: [],
  Qty: 0,
  Sum: 0,
});
const statuses = ref([]);

const isSender = ref(false);
const isCustomer = ref(false);

onMounted(async () => {
  let res = await fetch(ctx.rbUrl() + "/api/orders/statuses");
  if (res.ok) {
    statuses.value = await res.json();
  }

  res = await fetch(ctx.rbUrl() + "/api/orders?ref=" + route.params.id, {
    headers: await ctx.authHeadersAppJson(),
  });
  if (await ctx.CheckUnauth(res)) return;
  if (res.ok) {
    order.value = await res.json();

    order.value.Nbr = order.value.Ref.substring(7, 12);
    isCustomer.value = order.value.CreatedBy.Ref == ctx.userInfo.id;
    isSender.value = order.value.Sender.Ref == ctx.userInfo.id;

    if (isCustomer.value || isSender.value) {
      res = await fetch(
        ctx.rbUrl() + "/api/orderlines?orderref=" + order.value.Ref,
        {
          headers: await ctx.authHeadersAppJson(),
        }
      );
      if (await ctx.CheckUnauth(res)) return;
      if (res.ok) {
        const lines = await res.json();
        order.value.lines = lines.Data;
      }
    }
  }
});

const Save = async () => {
  let ok = false;
  if (isCustomer.value || isSender.value) {
    let req = await fetch(`${ctx.rbUrl()}/api/orders?ref=${route.params.id}`, {
      method: "PUT",
      body: JSON.stringify(order.value),
      headers: await ctx.authHeadersAppJson(),
    });
    if (await ctx.CheckUnauth(req)) return;
    ok = req.ok;
  }
  if (ok && isCustomer.value) {
    router.push("/orders");
  }
  if (ok && isSender.value) {
    router.push("/inc-orders");
  }
};
</script>

<template>
  <div v-if="order">
    <button
      :disabled="!isCustomer && !isSender"
      class="btn btn-primary"
      @click="Save()"
    >
      Сохранить
    </button>
    <h1>Заказ {{ order.Nbr }}</h1>
    <h6>{{ ctx.fmtDate(order.CreatedAt) }}</h6>
    <h6>Отправитель <ProfileLink :userInfo="order.Sender"></ProfileLink></h6>
    <h6>Заказчик <ProfileLink :userInfo="order.CreatedBy"></ProfileLink></h6>
    <div>&nbsp;</div>

    <table class="table table-sm">
      <thead class="table-primary">
        <tr>
          <th>Товар или услуга</th>
          <th>Количество</th>
          <th>Цена</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="line in order.lines" v-bind:key="line.Ref">
          <td>
            <RouterLink :to="`/good/${line.Good.Ref}`">{{
              line.Good.Caption
            }}</RouterLink>
          </td>
          <td>{{ line.Qty }}</td>
          <td>{{ line.Sum }}</td>
        </tr>
      </tbody>
      <tfoot class="table-primary">
        <tr>
          <td>ИТОГО</td>
          <td>{{ order.Qty }}</td>
          <td>{{ order.Sum }}</td>
        </tr>
      </tfoot>
    </table>
  </div>

  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Статус</label>
      <div class="col-9">
        <select
          :disabled="!isSender"
          class="form-select form-control-sm"
          v-model="order.Status"
        >
          <template v-for="(value, key) in statuses">
            <option v-bind:value="Number(key)">{{ value }}</option>
          </template>
        </select>
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Примечание покупателя</label>
      <div class="col-9">
        <textarea
          :disabled="!isCustomer"
          class="form-control"
          v-model="order.CustomerComment"
          rows="5"
        />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Примечание продавца</label>
      <div class="col-9">
        <textarea
          :disabled="!isSender"
          class="form-control"
          v-model="order.SenderComment"
          rows="5"
        />
      </div>
    </div>
  </div>
  <div class="row mb-3">
    <div class="form-group form-group-sm row">
      <label class="col-3 form-label">Ожидаемая дата доставки</label>
      <div class="col-9">
        <input
          :disabled="!isSender"
          type="date"
          v-model="order.ExpectedDeliveryDate"
        />
      </div>
    </div>
  </div>
</template>
