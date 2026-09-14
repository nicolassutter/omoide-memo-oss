import { heyApiPlugin } from '@hey-api/vite-plugin'
import config from './openapi-ts.config'

export const vitePlugin = () =>
  heyApiPlugin({
    config,
    vite: {
      apply: 'serve',
    },
  })
