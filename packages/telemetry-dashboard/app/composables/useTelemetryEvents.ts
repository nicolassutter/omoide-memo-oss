import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { fetchEvents, type Event as TelemetryEvent } from '@omoide-memo-oss/telemetry-sdk'

function getStartTimeForDays(daysString: string) {
  const days = parseInt(daysString, 10)
  if (days === 0) {
    return undefined
  }

  const date = new Date()
  if (days === 1) {
    date.setHours(date.getHours() - 24)
  } else {
    date.setDate(date.getDate() - days)
  }
  return date.toISOString()
}

export function useTelemetryEvents() {
  const daysFilter = ref<string>('30')

  const query = useQuery({
    queryKey: ['telemetry-events', daysFilter],
    queryFn: async () => {
      const startTime = getStartTimeForDays(daysFilter.value)
      const response = await fetchEvents({
        query: {
          startTime,
        },
      })
      return response.data
    },
  })

  const events = computed<TelemetryEvent[]>(() => {
    return query.data.value?.events ?? []
  })

  return {
    daysFilter,
    events,
    isLoading: query.isLoading,
    error: query.error,
    refetch: query.refetch,
    isFetching: query.isFetching,
  }
}
