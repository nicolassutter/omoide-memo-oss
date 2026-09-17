import { type UserConfig, createClient } from '@hey-api/openapi-ts'

export type ClientType = 'ky' | 'fetch'

export type Config = {
  spec: string
  output: string
  client?: ClientType
}

export async function createSdk(config: Config) {
  const clientPlugin = config.client === 'fetch' ? '@hey-api/client-fetch' : '@hey-api/client-ky'

  const userConfig: UserConfig = {
    input: config.spec,
    output: config.output,
    plugins: [
      clientPlugin,
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
