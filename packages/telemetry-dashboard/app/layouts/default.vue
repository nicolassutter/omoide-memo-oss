<script setup lang="ts">
import { Sun, Moon } from '@lucide/vue'
import { Button } from '~/components/ui/button'

useHead({
  titleTemplate: (titleChunk) => {
    return titleChunk ? `${titleChunk} - Omoide Memo Telemetry` : 'Omoide Memo Telemetry'
  },
})

useSeoMeta({
  title: 'Omoide Memo Telemetry',
  description: 'Self-agnostic operational analytics platform',
  ogTitle: 'Omoide Memo Telemetry',
  ogDescription: 'Self-agnostic operational analytics platform',
  ogType: 'website',
})

const colorMode = useColorMode()
function toggleColorMode() {
  colorMode.preference = colorMode.preference === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <div class="min-h-screen bg-background text-foreground transition-colors duration-300">
    <header class="border-b border-border/40 bg-background/95 backdrop-blur-md sticky top-0 z-50">
      <div
        class="container mx-auto flex flex-col md:flex-row items-center justify-between gap-4 py-4 px-4"
      >
        <div class="space-y-1 text-center md:text-left">
          <h1 class="text-xl font-bold tracking-tight">Omoide Memo Telemetry</h1>
          <p class="text-xs text-muted-foreground">Self-agnostic operational analytics platform</p>
        </div>

        <div class="flex items-center gap-3">
          <slot name="header-actions" />

          <Button variant="outline" size="icon" class="border-border/40" @click="toggleColorMode">
            <Sun v-if="colorMode.value === 'dark'" class="w-4 h-4 text-muted-foreground" />
            <Moon v-else class="w-4 h-4 text-muted-foreground" />
          </Button>
        </div>
      </div>
    </header>

    <slot />
  </div>
</template>
