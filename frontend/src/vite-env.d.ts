/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Base URL of the calculator API. Defaults to "/api" (dev/prod proxy). */
  readonly VITE_API_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
