import { defineConfig } from 'orval'

export default defineConfig({
  taskmanager: {
    output: 'src/lib/api/index.ts',
    input: '../backendGoVanilaTaskmanager/swagger/openapi.yaml',
  },
})
