<script setup lang="ts">
// eslint-disable-next-line no-unused-vars
const props = defineProps(["userInfo"]);
</script>

<template>
  <RouterLink :to="`/profile/${userInfo.Ref}`">{{
    userInfo.Username
  }}</RouterLink>
</template>
