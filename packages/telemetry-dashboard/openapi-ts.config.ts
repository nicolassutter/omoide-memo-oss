import type { UserConfig } from '@hey-api/openapi-ts'

export default {
  input: '../telemetry-api/doc/openapi.yaml',
  output: 'app/lib/client',
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
