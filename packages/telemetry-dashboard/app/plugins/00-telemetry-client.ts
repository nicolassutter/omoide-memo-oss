import { client } from '@/lib/telemetry-sdk'
import { defineNuxtPlugin } from '#imports'
import ky from 'ky'

export default defineNuxtPlugin(() => {
  const instance = ky.create({
    hooks: {
      beforeRequest: [
        ({ request }) => {
          console.log('[Telemetry Client] requesting: ', request.url)
        },
      ],
    },
  })

  const baseUrl = import.meta.env.DEV ? 'http://localhost:9999' : '/api'

  client.setConfig({
    baseUrl,
    ky: instance,
  })
})
