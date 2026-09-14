import type { UserConfig } from '@hey-api/openapi-ts'
import { resolve } from 'node:path'

const dirname = new URL('.', import.meta.url).pathname
const openApiFile = resolve(dirname, '../telemetry-api/doc/openapi.yaml')
const dist = resolve(dirname, './src')

export default {
  input: openApiFile,
  output: dist,
  plugins: [
    '@hey-api/client-ky',
    {
      name: '@hey-api/sdk',
      operations: {
        strategy: 'flat',
      },
    },
    {
      name: 'zod',
      compatibilityVersion: 4,
      definitions: true,
      requests: true,
      responses: true,
    },
    '@tanstack/vue-query',
  ],
} satisfies UserConfig
