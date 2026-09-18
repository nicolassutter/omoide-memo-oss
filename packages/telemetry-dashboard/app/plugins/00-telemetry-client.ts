import { client } from '@/lib/telemetry-sdk'
import { defineNuxtPlugin } from '#imports'
import ky from 'ky'

export default defineNuxtPlugin(async () => {
  const { apiBaseUrl, apiKey } = await $fetch<{ apiBaseUrl: string; apiKey: string }>(
    '/bff/_telemetry-config',
  )

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
    baseUrl: apiBaseUrl,
    headers: {
      'X-Api-Key': apiKey,
    },
    ky: instance,
  })
})
