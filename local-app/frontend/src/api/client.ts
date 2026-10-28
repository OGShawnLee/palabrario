import * as WailsAppBridge from '../../bindings/palabrario/local-app/app';
import { WailsDesktopAdapter } from './wails.adapter';
import { CapacitorMobileAdapter } from './capacitor.adapter';

export type IPalabrarioAPI = typeof WailsAppBridge;

/**
 * Determina si el entorno actual es el runtime de escritorio Wails v3.
 */
export function isWailsRuntime(): boolean {
  if (typeof window === 'undefined') return false;
  return (window as unknown as { wails?: unknown }).wails !== undefined;
}

/**
 * Cliente API Polimórfico Único (Singleton).
 * Despacha llamadas dinámicamente hacia WailsDesktopAdapter en escritorio Wails v3,
 * o hacia CapacitorMobileAdapter en móviles Capacitor o navegador Web SPA.
 */
export const api: IPalabrarioAPI = new Proxy({} as IPalabrarioAPI, {
  get(_target, prop: string | symbol) {
    const adapter = isWailsRuntime() ? WailsDesktopAdapter : CapacitorMobileAdapter;
    const value = (adapter as any)[prop];
    if (typeof value === 'function') {
      return value.bind(adapter);
    }
    return value;
  },
});
