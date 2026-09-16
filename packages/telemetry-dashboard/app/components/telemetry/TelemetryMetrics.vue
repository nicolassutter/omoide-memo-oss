<script setup lang="ts">
import { computed } from 'vue'
import type { Event as TelemetryEvent } from '@/lib/telemetry-sdk'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import { Database, Users, Activity, TrendingUp } from '@lucide/vue'

const props = defineProps<{
  events: TelemetryEvent[]
}>()

const totalEvents = computed(() => {
  return props.events.length
})

const uniqueDevices = computed(() => {
  const deviceIds = new Set(props.events.map((event) => event.DeviceID))
  return deviceIds.size
})

const averageEventsPerDevice = computed(() => {
  if (uniqueDevices.value === 0) {
    return 0
  }
  return parseFloat((totalEvents.value / uniqueDevices.value).toFixed(1))
})

const dailyVelocity = computed(() => {
  if (props.events.length === 0) {
    return 0
  }

  const timestamps = props.events.map((event) => new Date(event.TrackedAt).getTime())
  const minimumTimestamp = Math.min(...timestamps)
  const maximumTimestamp = Math.max(...timestamps)

  const millisecondsDifference = maximumTimestamp - minimumTimestamp
  const daysDifference = millisecondsDifference / (1000 * 60 * 60 * 24)
  const activeDays = Math.max(1, Math.ceil(daysDifference))

  return parseFloat((totalEvents.value / activeDays).toFixed(1))
})
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
    <Card class="shadow-sm border-border/40">
      <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <div class="space-y-1">
          <CardTitle class="text-sm font-medium text-muted-foreground"
            >Total Events Ingested</CardTitle
          >
          <CardDescription>All tracked client actions</CardDescription>
        </div>
        <Database class="w-4 h-4 text-primary" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold tracking-tight">{{ totalEvents }}</div>
      </CardContent>
    </Card>

    <Card class="shadow-sm border-border/40">
      <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <div class="space-y-1">
          <CardTitle class="text-sm font-medium text-muted-foreground"
            >Unique Active Devices</CardTitle
          >
          <CardDescription>Distinct client identifiers</CardDescription>
        </div>
        <Users class="w-4 h-4 text-primary" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold tracking-tight">{{ uniqueDevices }}</div>
      </CardContent>
    </Card>

    <Card class="shadow-sm border-border/40">
      <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <div class="space-y-1">
          <CardTitle class="text-sm font-medium text-muted-foreground"
            >Average Events / Device</CardTitle
          >
          <CardDescription>Engagement density ratio</CardDescription>
        </div>
        <Activity class="w-4 h-4 text-primary" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold tracking-tight">{{ averageEventsPerDevice }}</div>
      </CardContent>
    </Card>

    <Card class="shadow-sm border-border/40">
      <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
        <div class="space-y-1">
          <CardTitle class="text-sm font-medium text-muted-foreground"
            >Daily Event Velocity</CardTitle
          >
          <CardDescription>Mean events ingested per day</CardDescription>
        </div>
        <TrendingUp class="w-4 h-4 text-primary" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold tracking-tight">{{ dailyVelocity }}/day</div>
      </CardContent>
    </Card>
  </div>
</template>
