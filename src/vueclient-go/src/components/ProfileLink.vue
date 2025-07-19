<script setup lang="ts">

  // eslint-disable-next-line no-unused-vars
  const props = defineProps(['userInfo']);

</script>

<template>

  <RouterLink :to="`/profile/${userInfo.id}`">{{userInfo.userName}}</RouterLink>

</template>
