import { client } from '~/lib/client/client.gen'
import { defineNuxtPlugin, useRuntimeConfig } from '#imports'

export default defineNuxtPlugin(() => {
  const runtimeConfiguration = useRuntimeConfig()
  const telemetryServerBaseUrl = runtimeConfiguration.public.telemetryServerBaseUrl

  client.setConfig({
    baseUrl: telemetryServerBaseUrl,
  })
})
