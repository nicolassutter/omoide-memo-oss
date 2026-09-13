import { defineConfig } from 'oxlint'

export default defineConfig({
  options: {
    typeAware: true,
  },
  plugins: ['eslint', 'promise', 'node', 'jsdoc', 'oxc', 'import', 'typescript', 'unicorn', 'vue'],
})
