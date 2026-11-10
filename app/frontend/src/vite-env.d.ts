/// <reference types="vite/client" />

import "vue-router"

declare module "vue-router" {
  interface RouteMeta {
    breadcrumb?: Array<{ label: string; to?: string }>
  }
}
