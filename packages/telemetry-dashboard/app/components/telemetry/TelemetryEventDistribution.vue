<script setup lang="ts">
import { computed } from 'vue'
import type { Event as TelemetryEvent } from '~/lib/client/types.gen'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import { BarChart3 } from '@lucide/vue'

const props = defineProps<{
  events: TelemetryEvent[]
}>()

const eventDistributionSeries = computed(() => {
  const counts: Record<string, number> = {}
  props.events.forEach((event) => {
    counts[event.Name] = (counts[event.Name] || 0) + 1
  })

  const sortedEntries = Object.entries(counts)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)

  const maximumCount = sortedEntries[0]?.count || 0

  return sortedEntries.map((entry) => ({
    name: entry.name,
    count: entry.count,
    percentage: maximumCount > 0 ? parseFloat(((entry.count / maximumCount) * 100).toFixed(1)) : 0,
  }))
})
</script>

<template>
  <Card class="shadow-sm border-border/40 flex flex-col justify-between">
    <CardHeader class="pb-4">
      <CardTitle class="text-base font-semibold flex items-center space-x-2">
        <BarChart3 class="w-4 h-4 text-primary" />
        <span>Event Types Distribution</span>
      </CardTitle>
      <CardDescription>Relative frequency ranking of unique client event names</CardDescription>
    </CardHeader>
    <CardContent class="space-y-6 pb-8">
      <div v-if="eventDistributionSeries.length > 0" class="space-y-4">
        <div v-for="item in eventDistributionSeries" :key="item.name" class="space-y-2">
          <div class="flex items-center justify-between text-sm">
            <span class="font-medium truncate mr-2">{{ item.name }}</span>
            <span class="text-muted-foreground font-semibold flex-shrink-0"
              >{{ item.count }} occurrences</span
            >
          </div>
          <div class="h-3 w-full bg-muted rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500 bg-chart-1"
              :style="{ width: `${item.percentage}%` }"
            ></div>
          </div>
        </div>
      </div>
      <div v-else class="h-40 flex items-center justify-center text-muted-foreground">
        No event types found in the current dataset
      </div>
    </CardContent>
  </Card>
</template>
