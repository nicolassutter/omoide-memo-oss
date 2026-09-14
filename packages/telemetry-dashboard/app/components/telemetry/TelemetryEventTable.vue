<script setup lang="ts">
import { computed } from 'vue'
import type { Event as TelemetryEvent } from '@omoide-memo-oss/telemetry-sdk'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from '~/components/ui/table'
import { ClipboardList } from '@lucide/vue'

const props = defineProps<{
  events: TelemetryEvent[]
}>()

const recentEvents = computed(() => {
  return [...props.events]
    .sort((a, b) => new Date(b.TrackedAt).getTime() - new Date(a.TrackedAt).getTime())
    .slice(0, 100)
})

function formatTimestamp(timestamp: string) {
  const date = new Date(timestamp)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    second: '2-digit',
  })
}

function parseProperties(propertiesRaw: unknown) {
  if (!propertiesRaw) {
    return []
  }

  try {
    const parsed = typeof propertiesRaw === 'string' ? JSON.parse(propertiesRaw) : propertiesRaw
    if (parsed && typeof parsed === 'object') {
      return Object.entries(parsed).map(([key, value]) => ({
        key,
        value: typeof value === 'object' ? JSON.stringify(value) : String(value),
      }))
    }
  } catch {
    return []
  }

  return []
}
</script>

<template>
  <Card class="shadow-sm border-border/40">
    <CardHeader class="pb-4">
      <CardTitle class="text-base font-semibold flex items-center space-x-2">
        <ClipboardList class="w-4 h-4 text-primary" />
        <span>Recent Telemetry Event Log</span>
      </CardTitle>
      <CardDescription>Real-time list of the last 100 ingested client events</CardDescription>
    </CardHeader>
    <CardContent class="p-0">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader class="bg-muted/30">
            <TableRow class="border-border/40">
              <TableHead class="w-44 font-semibold">Timestamp</TableHead>
              <TableHead class="w-32 font-semibold">Event Type</TableHead>
              <TableHead class="w-80 font-semibold">Device Identifier</TableHead>
              <TableHead class="font-semibold">Associated Properties</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="event in recentEvents"
              :key="event.ID"
              class="border-border/40 hover:bg-muted/10"
            >
              <TableCell class="font-mono text-xs text-muted-foreground">
                {{ formatTimestamp(event.TrackedAt) }}
              </TableCell>
              <TableCell>
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20"
                >
                  {{ event.Name }}
                </span>
              </TableCell>
              <TableCell class="font-mono text-xs text-muted-foreground truncate max-w-[20rem]">
                {{ event.DeviceID }}
              </TableCell>
              <TableCell>
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="item in parseProperties(event.Props)"
                    :key="item.key"
                    class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] bg-muted text-muted-foreground border border-border/60"
                  >
                    <span class="font-bold text-foreground/80">{{ item.key }}:</span>
                    <span>{{ item.value }}</span>
                  </span>
                  <span
                    v-if="parseProperties(event.Props).length === 0"
                    class="text-xs text-muted-foreground/50 italic"
                  >
                    none
                  </span>
                </div>
              </TableCell>
            </TableRow>
            <TableRow v-if="recentEvents.length === 0">
              <TableCell colspan="4" class="h-32 text-center text-muted-foreground">
                No telemetry events found for this filter
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </CardContent>
  </Card>
</template>
