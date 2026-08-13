<template>
  <a-app>
    <AppRegistrar />
    <router-view />
  </a-app>
</template>

<script setup lang="ts">
import AppRegistrar from '@/components/AppRegistrar.vue'
</script>
