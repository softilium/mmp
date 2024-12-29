<template>
  <a href="/">Home</a>
  {{greeting}}
</template>

<script>
  export default {
    data() {
      return {
        greeting: 'Привет всем!'
      }
    },
    methods:
    {
      $backendUrl: () => "http://localhost:5078"
    },
    async mounted() {

      let response = await fetch(this.$backendUrl() + "/identity/refresh", { method: 'POST' });
      this.greeting = "123123123";
    }
  }
</script>
