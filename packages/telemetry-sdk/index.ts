import { type UserConfig, createClient } from '@hey-api/openapi-ts'

export type Config = {
  spec: string
  output: string
}

export async function createSdk(config: Config) {
  const userConfig: UserConfig = {
    input: config.spec,
    output: config.output,
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

  await createClient(userConfig)
}
