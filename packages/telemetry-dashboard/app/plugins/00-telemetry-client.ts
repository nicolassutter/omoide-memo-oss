import { client } from '~/lib/client/client.gen'
import { defineNuxtPlugin, useRuntimeConfig } from '#imports'
import ky from 'ky'

export default defineNuxtPlugin(() => {
  const runtimeConfiguration = useRuntimeConfig()
  const telemetryServerBaseUrl = runtimeConfiguration.public.telemetryServerBaseUrl
  const telemetryApiKey = runtimeConfiguration.public.telemetryApiKey

  const instance = ky.create({
    hooks: {
      beforeRequest: [
        ({ request }) => {
          console.log('[Telemetry Client] requesting: ', request.url)
        },
      ],
    },
  })

  client.setConfig({
    baseUrl: telemetryServerBaseUrl,
    headers: {
      'X-Api-Key': telemetryApiKey,
    },
    ky: instance,
  })
})
