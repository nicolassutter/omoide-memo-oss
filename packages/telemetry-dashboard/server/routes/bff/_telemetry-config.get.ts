export default defineEventHandler(() => {
  const config = useRuntimeConfig()
  return {
    apiBaseUrl: '/api',
    apiKey: config.telemetryApiKey,
  }
})
