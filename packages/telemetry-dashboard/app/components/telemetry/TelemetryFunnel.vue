<script setup lang="ts">
import { ref, computed, watch, watchEffect } from 'vue'
import type { Event as TelemetryEvent } from '~/lib/client/types.gen'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectItem,
} from '~/components/ui/select'
import { Layers } from '@lucide/vue'

const props = defineProps<{
  events: TelemetryEvent[]
}>()

const uniqueEventNames = computed(() => {
  const names = new Set(props.events.map((event) => event.Name))
  return Array.from(names).sort()
})

const stageOneEventName = ref<string>('')
const stageTwoEventName = ref<string>('')
const stageThreeEventName = ref<string>('')

watchEffect(() => {
  const names = uniqueEventNames.value
  if (names.length > 0) {
    if (!stageOneEventName.value || !names.includes(stageOneEventName.value)) {
      stageOneEventName.value = names[0] || ''
    }
    if (!stageTwoEventName.value || !names.includes(stageTwoEventName.value)) {
      stageTwoEventName.value = names[1] || names[0] || ''
    }
    if (!stageThreeEventName.value || !names.includes(stageThreeEventName.value)) {
      stageThreeEventName.value = names[2] || names[1] || names[0] || ''
    }
  }
})

const funnelStages = computed(() => {
  if (props.events.length === 0 || !stageOneEventName.value) {
    return []
  }

  const deviceStageOneTimestamps: Record<string, number> = {}
  props.events.forEach((event) => {
    if (event.Name === stageOneEventName.value) {
      const timestamp = new Date(event.TrackedAt).getTime()
      const existing = deviceStageOneTimestamps[event.DeviceID]
      if (existing === undefined || timestamp < existing) {
        deviceStageOneTimestamps[event.DeviceID] = timestamp
      }
    }
  })

  const stageOneCount = Object.keys(deviceStageOneTimestamps).length
  if (stageOneCount === 0) {
    return [
      { stage: 'Stage 1', eventName: stageOneEventName.value, count: 0, percentage: 0 },
      { stage: 'Stage 2', eventName: stageTwoEventName.value, count: 0, percentage: 0 },
      { stage: 'Stage 3', eventName: stageThreeEventName.value, count: 0, percentage: 0 },
    ]
  }

  const deviceStageTwoTimestamps: Record<string, number> = {}
  if (stageTwoEventName.value) {
    props.events.forEach((event) => {
      if (event.Name === stageTwoEventName.value) {
        const stageOneTime = deviceStageOneTimestamps[event.DeviceID]
        if (stageOneTime !== undefined) {
          const timestamp = new Date(event.TrackedAt).getTime()
          if (timestamp >= stageOneTime) {
            const existing = deviceStageTwoTimestamps[event.DeviceID]
            if (existing === undefined || timestamp < existing) {
              deviceStageTwoTimestamps[event.DeviceID] = timestamp
            }
          }
        }
      }
    })
  }

  const stageTwoCount = Object.keys(deviceStageTwoTimestamps).length

  const deviceStageThreeTimestamps: Record<string, number> = {}
  if (stageThreeEventName.value) {
    props.events.forEach((event) => {
      if (event.Name === stageThreeEventName.value) {
        const stageTwoTime = deviceStageTwoTimestamps[event.DeviceID]
        if (stageTwoTime !== undefined) {
          const timestamp = new Date(event.TrackedAt).getTime()
          if (timestamp >= stageTwoTime) {
            const existing = deviceStageThreeTimestamps[event.DeviceID]
            if (existing === undefined || timestamp < existing) {
              deviceStageThreeTimestamps[event.DeviceID] = timestamp
            }
          }
        }
      }
    })
  }

  const stageThreeCount = Object.keys(deviceStageThreeTimestamps).length

  return [
    {
      stage: 'Stage 1',
      eventName: stageOneEventName.value,
      count: stageOneCount,
      percentage: 100,
    },
    {
      stage: 'Stage 2',
      eventName: stageTwoEventName.value,
      count: stageTwoCount,
      percentage: parseFloat(((stageTwoCount / stageOneCount) * 100).toFixed(1)),
    },
    {
      stage: 'Stage 3',
      eventName: stageThreeEventName.value,
      count: stageThreeCount,
      percentage: parseFloat(((stageThreeCount / stageOneCount) * 100).toFixed(1)),
    },
  ]
})
</script>

<template>
  <Card class="shadow-sm border-border/40">
    <CardHeader class="pb-4">
      <CardTitle class="text-base font-semibold flex items-center space-x-2">
        <Layers class="w-4 h-4 text-primary" />
        <span>Custom Interactive Funnel</span>
      </CardTitle>
      <CardDescription
        >Track conversion, progression, and dropoff across custom selected sequential
        events</CardDescription
      >
    </CardHeader>
    <CardContent class="space-y-6 pb-8">
      <div
        class="grid grid-cols-1 md:grid-cols-3 gap-4 bg-muted/30 p-4 rounded-lg border border-border/50"
      >
        <div class="space-y-2">
          <label class="text-xs font-semibold text-muted-foreground uppercase tracking-wider"
            >Stage 1 Event</label
          >
          <Select v-model="stageOneEventName">
            <SelectTrigger class="w-full bg-background border-border/40">
              <SelectValue placeholder="Select event..." />
            </SelectTrigger>
            <SelectContent class="bg-background border-border/40">
              <SelectGroup>
                <SelectItem v-for="name in uniqueEventNames" :key="name" :value="name">
                  {{ name }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-semibold text-muted-foreground uppercase tracking-wider"
            >Stage 2 Event</label
          >
          <Select v-model="stageTwoEventName">
            <SelectTrigger class="w-full bg-background border-border/40">
              <SelectValue placeholder="Select event..." />
            </SelectTrigger>
            <SelectContent class="bg-background border-border/40">
              <SelectGroup>
                <SelectItem v-for="name in uniqueEventNames" :key="name" :value="name">
                  {{ name }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-semibold text-muted-foreground uppercase tracking-wider"
            >Stage 3 Event</label
          >
          <Select v-model="stageThreeEventName">
            <SelectTrigger class="w-full bg-background border-border/40">
              <SelectValue placeholder="Select event..." />
            </SelectTrigger>
            <SelectContent class="bg-background border-border/40">
              <SelectGroup>
                <SelectItem v-for="name in uniqueEventNames" :key="name" :value="name">
                  {{ name }}
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div v-if="funnelStages.length > 0" class="space-y-6 pt-4">
        <div v-for="(item, index) in funnelStages" :key="item.stage" class="relative">
          <div class="space-y-2">
            <div class="flex items-center justify-between text-sm">
              <div class="flex items-center space-x-2">
                <span class="font-bold text-primary">{{ item.stage }}:</span>
                <span class="font-semibold">{{ item.eventName }}</span>
              </div>
              <span class="text-muted-foreground font-semibold"
                >{{ item.count }} devices ({{ item.percentage }}%)</span
              >
            </div>
            <div class="h-3 w-full bg-muted rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500 bg-chart-1"
                :style="{ width: `${item.percentage}%` }"
              ></div>
            </div>
          </div>

          <div v-if="index < funnelStages.length - 1" class="flex justify-center py-2">
            <div class="text-[10px] text-muted-foreground/60 flex items-center space-x-1 font-bold">
              <span>Dropoff: {{ 100 - (funnelStages[index + 1]?.percentage ?? 0) }}%</span>
            </div>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
