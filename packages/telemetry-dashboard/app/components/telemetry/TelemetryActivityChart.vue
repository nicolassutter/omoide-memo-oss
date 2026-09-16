<script setup lang="ts">
import { computed } from 'vue'
import type { Event as TelemetryEvent } from '@/lib/telemetry-sdk'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import { Activity } from '@lucide/vue'
import { VisXYContainer, VisArea, VisAxis } from '@unovis/vue'
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
  type ChartConfig,
} from '~/components/ui/chart'

const props = defineProps<{
  events: TelemetryEvent[]
  daysFilter: string
}>()

const timelineSeries = computed(() => {
  if (props.events.length === 0) {
    return []
  }

  const counts: Record<string, number> = {}
  const now = new Date()
  const startDate = new Date()
  const isLast24Hours = props.daysFilter === '1'

  if (isLast24Hours) {
    startDate.setHours(now.getHours() - 24)
  } else {
    startDate.setDate(now.getDate() - parseInt(props.daysFilter, 10))
  }

  props.events.forEach((event) => {
    const date = new Date(event.TrackedAt)
    if (date >= startDate && date <= now) {
      const key = isLast24Hours
        ? `${date.toISOString().split(':')[0] || ''}:00`
        : date.toISOString().split('T')[0] || ''
      counts[key] = (counts[key] || 0) + 1
    }
  })

  const result: Array<{ date: number; events: number }> = []
  const current = new Date(startDate)

  while (current <= now) {
    const key = isLast24Hours
      ? `${current.toISOString().split(':')[0] || ''}:00`
      : current.toISOString().split('T')[0] || ''

    result.push({
      date: current.getTime(),
      events: counts[key] || 0,
    })

    if (isLast24Hours) {
      current.setHours(current.getHours() + 1)
    } else {
      current.setDate(current.getDate() + 1)
    }
  }

  return result
})

function formatChartX(tick: number | Date) {
  const date = new Date(tick)
  if (props.daysFilter === '1') {
    return date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })
  }
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function formatTooltipLabel(tick: number | Date) {
  const date = new Date(tick)
  if (props.daysFilter === '1') {
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    })
  }
  return date.toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })
}

const timelineChartConfig: ChartConfig = {
  events: {
    label: 'Telemetry Events Ingested',
    color: 'var(--chart-1)',
  },
}
</script>

<template>
  <Card class="shadow-sm border-border/40">
    <CardHeader class="flex flex-row items-center justify-between pb-4 space-y-0">
      <div class="space-y-1">
        <CardTitle class="text-base font-semibold flex items-center space-x-2">
          <Activity class="w-4 h-4 text-primary" />
          <span>Ingestion Activity Timeline</span>
        </CardTitle>
        <CardDescription
          >Chronological volume distribution of client telemetry events</CardDescription
        >
      </div>
    </CardHeader>
    <CardContent class="pb-4">
      <div class="h-80 w-full">
        <ChartContainer
          v-if="timelineSeries.length > 0"
          :data="timelineSeries"
          :config="timelineChartConfig"
        >
          <VisXYContainer
            :data="timelineSeries"
            :padding="{ left: 16, right: 16, top: 16, bottom: 16 }"
            class="h-full w-full"
          >
            <VisArea
              :x="(data: { date: number }) => data.date"
              :y="[(data: { events: number }) => data.events]"
              color="var(--chart-1)"
              :curveType="'monotoneX'"
            />
            <VisAxis
              type="x"
              :tickFormat="formatChartX"
              :numTicks="props.daysFilter === '1' ? 8 : 6"
              class="text-xs text-muted-foreground"
            />
            <VisAxis type="y" :numTicks="5" class="text-xs text-muted-foreground" />
            <ChartCrosshair />
            <ChartTooltip
              selectionMode="single"
              :render-content="
                (props: { items: Array<{ data: { date: number; events: number } }> }) => {
                  if (!props.items || props.items.length === 0) return ''
                  const item = props.items[0]
                  if (!item) return ''
                  return componentToString(timelineChartConfig, ChartTooltipContent, {
                    labelFormatter: formatTooltipLabel,
                    payload: [
                      {
                        name: 'Events',
                        value: item.data.events,
                        color: 'var(--chart-1)',
                      },
                    ],
                  })
                }
              "
            />
          </VisXYContainer>
        </ChartContainer>
        <div v-else class="h-full flex items-center justify-center text-muted-foreground">
          No data available for the selected timeframe
        </div>
      </div>
    </CardContent>
  </Card>
</template>
