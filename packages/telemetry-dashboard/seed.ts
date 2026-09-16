import { randomUUID } from 'crypto'
import { client, ingestEvents, type IngestBatch, type IngestEvent } from '@/lib/telemetry-sdk'

const API_KEY = process.env.TELEMETRY_API_KEY || 'your-api-key-here'

client.setConfig({
  baseUrl: 'http://localhost:9999',
  headers: {
    'X-Api-Key': API_KEY,
  },
})

const devicesCount = 40
const batches: IngestBatch[] = []

for (let i = 0; i < devicesCount; i++) {
  const deviceId = randomUUID()
  const deviceEvents: IngestEvent[] = []
  const startDaysAgo = Math.random() * 30
  let currentTime = Math.floor(Date.now() - startDaysAgo * 24 * 60 * 60 * 1000)

  deviceEvents.push({
    name: 'app_opened',
    props: {},
    trackedAt: Math.floor(currentTime),
  })

  currentTime += Math.floor(Math.random() * 5 * 60 * 1000)
  const finishedOnboarding = Math.random() < 0.9
  if (finishedOnboarding) {
    deviceEvents.push({
      name: 'onboarding_complete',
      props: {},
      trackedAt: Math.floor(currentTime),
    })

    currentTime += Math.floor(Math.random() * 10 * 60 * 1000)
    const signedUp = Math.random() < 0.7
    if (signedUp) {
      deviceEvents.push({
        name: 'signup_completed',
        props: {},
        trackedAt: Math.floor(currentTime),
      })

      const activeDays = Math.floor(Math.random() * 5) + 1
      for (let day = 1; day <= activeDays; day++) {
        currentTime += Math.floor((Math.random() * 3 + 1) * 24 * 60 * 60 * 1000)
        if (currentTime > Date.now()) break

        deviceEvents.push({
          name: 'login_completed',
          props: {},
          trackedAt: Math.floor(currentTime),
        })

        if (Math.random() < 0.4) {
          deviceEvents.push({
            name: 'folder_created',
            props: {},
            trackedAt: Math.floor(currentTime + 1000),
          })
        }

        const memoCount = Math.floor(Math.random() * 4) + 1
        for (let m = 0; m < memoCount; m++) {
          const source = Math.random() < 0.65 ? 'manual' : 'smart_scan'
          deviceEvents.push({
            name: 'memo_created',
            props: { source },
            trackedAt: Math.floor(currentTime + (m + 1) * 30 * 1000),
          })
        }

        if (Math.random() < 0.5) {
          deviceEvents.push({
            name: 'review_session_completed',
            props: {},
            trackedAt: Math.floor(currentTime + 5 * 60 * 1000),
          })
        }
      }
    }
  }

  batches.push({
    deviceId,
    sentAt: new Date().toISOString(),
    events: deviceEvents,
  })
}

async function seed() {
  let totalSeeded = 0
  for (let idx = 0; idx < batches.length; idx++) {
    const batch = batches[idx]
    if (!batch) continue
    try {
      const response = await ingestEvents({
        body: batch,
      })

      const result = response.data
      if (result && typeof result === 'object' && 'accepted' in result) {
        const acceptedCount = typeof result.accepted === 'number' ? result.accepted : 0
        totalSeeded += acceptedCount
        console.log(
          `Seeded batch ${idx + 1}/${batches.length} (device: ${batch.deviceId}): accepted ${acceptedCount} events.`,
        )
      } else {
        console.error(
          `Failed to seed batch ${idx + 1}/${batches.length}: unexpected response shape`,
          response,
        )
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err)
      console.error(`Error seeding batch ${idx + 1}:`, message)
    }
  }
  console.log(`Seeding complete! Total events accepted: ${totalSeeded}`)
}

await seed()
