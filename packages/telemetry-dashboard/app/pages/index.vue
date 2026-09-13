<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { Button } from '~/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { RefreshCw } from '@lucide/vue'
import { useTelemetryEvents } from '~/composables/useTelemetryEvents'

const { daysFilter, events, isLoading, error, refetch, isFetching } = useTelemetryEvents()
</script>

<template>
  <NuxtLayout>
    <template #header-actions>
      <Select v-model="daysFilter">
        <SelectTrigger class="w-36 bg-background border-border/40">
          <SelectValue placeholder="Timeframe" />
        </SelectTrigger>
        <SelectContent class="bg-background border-border/40">
          <SelectItem value="1">Last 24 Hours</SelectItem>
          <SelectItem value="7">Last 7 Days</SelectItem>
          <SelectItem value="30">Last 30 Days</SelectItem>
          <SelectItem value="0">All Time</SelectItem>
        </SelectContent>
      </Select>

      <Button
        variant="outline"
        size="icon"
        class="border-border/40"
        :disabled="isFetching"
        @click="() => refetch()"
      >
        <RefreshCw :class="['w-4 h-4 text-muted-foreground', isFetching ? 'animate-spin' : '']" />
      </Button>
    </template>

    <main class="container mx-auto px-4 py-8">
      <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 space-y-4">
        <RefreshCw class="w-8 h-8 text-primary animate-spin" />
        <p class="text-sm text-muted-foreground font-semibold">Synchronizing telemetry data...</p>
      </div>

      <div
        v-else-if="error"
        class="bg-destructive/10 border border-destructive/20 rounded-lg p-6 max-w-lg mx-auto text-center space-y-4"
      >
        <h2 class="text-lg font-semibold text-destructive">Data Sync Failure</h2>
        <p class="text-sm text-muted-foreground">
          Failed to connect to the telemetry backend server
        </p>
        <Button variant="outline" class="border-destructive/20" @click="() => refetch()">
          Retry Synchronization
        </Button>
      </div>

      <div v-else class="space-y-8">
        <TelemetryMetrics :events="events" />

        <Tabs default-value="overview" class="space-y-6">
          <div class="flex items-center justify-between border-b border-border/40 pb-2">
            <TabsList class="bg-muted/40 p-1 border border-border/40 rounded-lg">
              <TabsTrigger value="overview" class="rounded-md px-4 py-1.5 text-sm font-medium">
                Overview & Flows
              </TabsTrigger>
              <TabsTrigger value="logs" class="rounded-md px-4 py-1.5 text-sm font-medium">
                Live Ingest Log
              </TabsTrigger>
            </TabsList>
            <span class="text-xs text-muted-foreground font-semibold">
              Loaded {{ events.length }} records
            </span>
          </div>

          <TabsContent value="overview" class="space-y-6">
            <TelemetryActivityChart :events="events" :days-filter="daysFilter" />
            <div class="grid gap-6 md:grid-cols-2">
              <TelemetryEventDistribution :events="events" />
              <TelemetryFunnel :events="events" />
            </div>
          </TabsContent>

          <TabsContent value="logs">
            <TelemetryEventTable :events="events" />
          </TabsContent>
        </Tabs>
      </div>
    </main>
  </NuxtLayout>
</template>
