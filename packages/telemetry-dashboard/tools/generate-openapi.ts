import { createSdk } from '@omoide-memo-oss/telemetry-sdk'
import { join } from 'node:path'

const dirname = new URL('.', import.meta.url).pathname

await createSdk({
  spec: join(dirname, '../../telemetry-api/doc/openapi.yaml'),
  output: join(dirname, '../app/lib/telemetry-sdk'),
})
