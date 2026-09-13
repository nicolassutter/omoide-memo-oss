<script setup lang="ts">
import { VueQueryDevtools } from '@tanstack/vue-query-devtools'
</script>

<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtPage />
    <VueQueryDevtools />
  </div>
</template>
